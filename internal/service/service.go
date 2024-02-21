package service

import (
	"context"
	"rinha-backend-q1-2024/internal/canonical"
	"rinha-backend-q1-2024/internal/persistence"
)

type Service interface {
	DoTransaction(ctx context.Context, financialTransaction canonical.FinancialTransaction, clienteId int) (canonical.TransactionBalance, error)
	GetStatement(ctx context.Context, clienteId int) (canonical.Statement, error)
}

type myService struct {
	persistenceDb persistence.PersistenceLayer
}

func (ms *myService) DoTransaction(ctx context.Context, financialTransaction canonical.FinancialTransaction, clienteId int) (canonical.TransactionBalance, error) {
	balance, err := ms.persistenceDb.AddTransaction(ctx, clienteId, financialTransaction)
	if err != nil {
		return canonical.TransactionBalance{}, err
	}
	return balance, err
}

func (ms *myService) GetStatement(ctx context.Context, clienteId int) (canonical.Statement, error) {
	statement, err := ms.persistenceDb.GetStatement(ctx, clienteId)
	if err != nil {
		return canonical.Statement{}, err
	}
	return statement, nil
}

func NewService(db persistence.PersistenceLayer) Service {
	return &myService{
		persistenceDb: db,
	}
}
