package persistence

import (
	"context"
	"errors"
	"log"
	"os"
	"rinha-backend-q1-2024/internal/entities"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotEnoughBalance = errors.New("not enough balance")
)

type PersistenceLayer interface {
	AddTransaction(ctx context.Context, clienteId int, financialTransaction entities.FinancialTransaction) (entities.TransactionBalance, error)
	GetStatement(ctx context.Context, clienteId int) (entities.Statement, error)
	Close()
}

type postgresBase struct {
	db *pgxpool.Pool
}

func NewPersistence(ctx context.Context) PersistenceLayer {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	poolConfig.MinConns = 30
	poolConfig.MaxConns = 150

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	return &postgresBase{db: db}
}

func (pb *postgresBase) AddTransaction(ctx context.Context, clienteId int, financialTransaction entities.FinancialTransaction) (entities.TransactionBalance, error) {
	tx, err := pb.db.Begin(ctx)
	if err != nil {
		return entities.TransactionBalance{}, err
	}
	defer tx.Rollback(ctx)

	amount := financialTransaction.Valor
	if financialTransaction.Tipo == "d" {
		amount = financialTransaction.Valor * -1
	}

	row := tx.QueryRow(ctx, "UPDATE clientes SET saldo = saldo + ($1) WHERE id = ($2) AND (saldo + ($1)) >= (limite * -1) RETURNING saldo, limite", amount, clienteId)
	var balance int64
	var limit int64
	if err := row.Scan(&balance, &limit); err != nil {
		return entities.TransactionBalance{}, err
	}

	_, err = tx.Exec(ctx, "INSERT INTO transacoes (valor, tipo, descricao, cliente_id) VALUES ($1, $2, $3, $4)", financialTransaction.Valor, financialTransaction.Tipo, financialTransaction.Descricao, clienteId)
	if err != nil {
		return entities.TransactionBalance{}, err
	}

	return entities.TransactionBalance{
		Limite: limit,
		Saldo:  balance,
	}, tx.Commit(ctx)
}

func (pb *postgresBase) GetStatement(ctx context.Context, clienteId int) (entities.Statement, error) {
	statement := entities.Statement{
		LastTransactions: make([]entities.FinancialTransaction, 0),
	}

	row := pb.db.QueryRow(ctx, "SELECT saldo, limite FROM clientes WHERE id = $1", clienteId)
	var balance int64
	var limit int64
	if err := row.Scan(&balance, &limit); err != nil {
		return entities.Statement{}, err
	}
	statement.Balance.Limite = limit
	statement.Balance.Total = balance

	rows, err := pb.db.Query(ctx, "SELECT valor, tipo, descricao, realizado_em FROM transacoes WHERE cliente_id = $1 ORDER BY realizado_em DESC LIMIT 10", clienteId)
	if err != nil {
		return entities.Statement{}, err
	}

	for rows.Next() {
		var valor int64
		var tipo string
		var descricao string
		var realizada_em time.Time
		err = rows.Scan(&valor, &tipo, &descricao, &realizada_em)
		if err != nil {
			if err.Error() == pgx.ErrNoRows.Error() {
				return statement, nil
			}
			return entities.Statement{}, err
		}

		statement.LastTransactions = append(statement.LastTransactions, entities.FinancialTransaction{
			Valor:       valor,
			Tipo:        tipo,
			Descricao:   descricao,
			RealizadaEm: realizada_em.UTC().Format(time.RFC3339Nano),
		})
	}

	return statement, nil
}

func (pb *postgresBase) Close() {
	pb.db.Close()
}
