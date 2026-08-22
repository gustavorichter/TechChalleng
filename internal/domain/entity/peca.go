package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Peca representa uma peça ou insumo do estoque.
type Peca struct {
	ID               uuid.UUID
	Nome             string
	Codigo           string
	ValorUnitario    decimal.Decimal
	QuantidadeEstoque int
	Ativo            bool
	CriadoEm         time.Time
	AtualizadoEm     time.Time
}

func NewPeca(nome, codigo string, valorUnitario decimal.Decimal, quantidade int) *Peca {
	now := time.Now().UTC()
	return &Peca{
		ID:                uuid.New(),
		Nome:              nome,
		Codigo:            codigo,
		ValorUnitario:     valorUnitario,
		QuantidadeEstoque: quantidade,
		Ativo:             true,
		CriadoEm:          now,
		AtualizadoEm:      now,
	}
}

// Reservar verifica se há saldo suficiente para a quantidade solicitada.
func (p *Peca) TemSaldo(quantidade int) bool {
	return p.QuantidadeEstoque >= quantidade
}

// BaixarEstoque reduz o estoque físico após validação.
func (p *Peca) BaixarEstoque(quantidade int) error {
	if quantidade <= 0 {
		return errors.New("quantidade deve ser positiva")
	}
	if !p.TemSaldo(quantidade) {
		return errors.New("estoque insuficiente")
	}
	p.QuantidadeEstoque -= quantidade
	p.AtualizadoEm = time.Now().UTC()
	return nil
}

// AdicionarEstoque incrementa o estoque.
func (p *Peca) AdicionarEstoque(quantidade int) error {
	if quantidade <= 0 {
		return errors.New("quantidade deve ser positiva")
	}
	p.QuantidadeEstoque += quantidade
	p.AtualizadoEm = time.Now().UTC()
	return nil
}
