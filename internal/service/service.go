package service

import (
	"context"
	"rinha-backend-q1-2024/internal/entities"
	"rinha-backend-q1-2024/internal/persistence"
)

type Service interface {
	DoTransaction(ctx context.Context, financialTransaction entities.FinancialTransaction, clienteId int) (entities.TransactionBalance, error)
	GetStatement(ctx context.Context, clienteId int) (entities.Statement, error)
}

type myService struct {
	persistenceDb persistence.PersistenceLayer
}

func (ms *myService) DoTransaction(ctx context.Context, financialTransaction entities.FinancialTransaction, clienteId int) (entities.TransactionBalance, error) {
	balance, err := ms.persistenceDb.AddTransaction(ctx, clienteId, financialTransaction)
	if err != nil {
		return entities.TransactionBalance{}, err
	}
	return balance, err
}

func (ms *myService) GetStatement(ctx context.Context, clienteId int) (entities.Statement, error) {
	statement, err := ms.persistenceDb.GetStatement(ctx, clienteId)
	if err != nil {
		return entities.Statement{}, err
	}
	return statement, nil
}

func NewService(db persistence.PersistenceLayer) Service {
	return &myService{
		persistenceDb: db,
	}
}
