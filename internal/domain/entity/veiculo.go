package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

// Veiculo representa um veículo vinculado a um cliente.
type Veiculo struct {
	ID           uuid.UUID
	ClienteID    uuid.UUID
	Placa        valobj.Placa
	Marca        string
	Modelo       string
	Ano          int
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

func NewVeiculo(clienteID uuid.UUID, placa valobj.Placa, marca, modelo string, ano int) *Veiculo {
	now := time.Now().UTC()
	return &Veiculo{
		ID:           uuid.New(),
		ClienteID:    clienteID,
		Placa:        placa,
		Marca:        marca,
		Modelo:       modelo,
		Ano:          ano,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
}
