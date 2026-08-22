package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

// Cliente representa um cliente da oficina.
type Cliente struct {
	ID        uuid.UUID
	Nome      string
	CPF       *valobj.CPF
	CNPJ      *valobj.CNPJ
	Email     valobj.Email
	Telefone  string
	CriadoEm  time.Time
	AtualizadoEm time.Time
}

func NewCliente(nome string, email valobj.Email, telefone string) *Cliente {
	now := time.Now().UTC()
	return &Cliente{
		ID:           uuid.New(),
		Nome:         nome,
		Email:        email,
		Telefone:     telefone,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
}

func (c *Cliente) DefinirCPF(cpf valobj.CPF) {
	c.CPF = &cpf
	c.CNPJ = nil
}

func (c *Cliente) DefinirCNPJ(cnpj valobj.CNPJ) {
	c.CNPJ = &cnpj
	c.CPF = nil
}
