package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Servico representa um serviço oferecido pela oficina.
type Servico struct {
	ID            uuid.UUID
	Nome          string
	Descricao     string
	ValorMaoObra  decimal.Decimal
	Ativo         bool
	CriadoEm      time.Time
	AtualizadoEm  time.Time
}

func NewServico(nome, descricao string, valorMaoObra decimal.Decimal) *Servico {
	now := time.Now().UTC()
	return &Servico{
		ID:           uuid.New(),
		Nome:         nome,
		Descricao:    descricao,
		ValorMaoObra: valorMaoObra,
		Ativo:        true,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
}
