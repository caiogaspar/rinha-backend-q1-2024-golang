package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"rinha-backend-q1-2024/internal/canonical"
	"rinha-backend-q1-2024/internal/persistence"
	"rinha-backend-q1-2024/internal/service"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx"
)

const basePath = "/clientes"

var (
	mySvc service.Service
)

func main() {
	port := os.Getenv("PORT")
	app := fiber.New()

	ctx := context.Background()

	persistenceLayer := persistence.NewPersistence(ctx)
	mySvc = service.NewService(persistenceLayer)

	app.Post(basePath+"/:id/transacoes", transacaoHandler)
	app.Get(basePath+"/:id/extrato", extratoHandler)

	app.Listen(":" + port)

}

func transacaoHandler(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	clienteId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	var transacao canonical.FinancialTransaction
	if err := json.Unmarshal(ctx.Body(), &transacao); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if transacao.Tipo != "c" && transacao.Tipo != "d" {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if len(transacao.Descricao) < 1 || len(transacao.Descricao) > 10 {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	transactionResult, err := mySvc.DoTransaction(ctx.Context(), transacao, clienteId)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return ctx.SendStatus(fiber.StatusUnprocessableEntity)
		}

		log.Error(err)
		log.Info(err, fmt.Sprintf("%v", err))
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(transactionResult)
}

func extratoHandler(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	clienteId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	statementResult, err := mySvc.GetStatement(ctx.Context(), clienteId)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	statementResult.Balance.StatementDate = time.Now().UTC().Format(time.RFC3339Nano)

	return ctx.JSON(statementResult)
}
