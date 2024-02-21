package canonical

type FinancialTransaction struct {
	Valor       int64  `json:"valor,omitempty"`
	Tipo        string `json:"tipo,omitempty"`
	Descricao   string `json:"descricao,omitempty"`
	RealizadaEm string `json:"realizada_em,omitempty"`
}

type TransactionBalance struct {
	Limite int64 `json:"limite"`
	Saldo  int64 `json:"saldo"`
}

type StatementBalance struct {
	Limite        int64  `json:"limite"`
	Total         int64  `json:"total"`
	StatementDate string `json:"data_extrato"`
}

type Statement struct {
	Balance          StatementBalance       `json:"saldo,omitempty"`
	LastTransactions []FinancialTransaction `json:"ultimas_transacoes,omitempty"`
}
