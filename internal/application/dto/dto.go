package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// --- Auth ---

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// --- Cliente ---

type CriarClienteRequest struct {
	Nome     string `json:"nome" binding:"required"`
	CPF      string `json:"cpf"`
	CNPJ     string `json:"cnpj"`
	Email    string `json:"email" binding:"required,email"`
	Telefone string `json:"telefone" binding:"required"`
}

type ClienteResponse struct {
	ID        uuid.UUID `json:"id"`
	Nome      string    `json:"nome"`
	CPF       string    `json:"cpf,omitempty"`
	CNPJ      string    `json:"cnpj,omitempty"`
	Email     string    `json:"email"`
	Telefone  string    `json:"telefone"`
	CriadoEm  time.Time `json:"criado_em"`
}

// --- Veículo ---

type CriarVeiculoRequest struct {
	ClienteID uuid.UUID `json:"cliente_id" binding:"required"`
	Placa     string    `json:"placa" binding:"required"`
	Marca     string    `json:"marca" binding:"required"`
	Modelo    string    `json:"modelo" binding:"required"`
	Ano       int       `json:"ano" binding:"required,min=1900"`
}

type VeiculoResponse struct {
	ID        uuid.UUID `json:"id"`
	ClienteID uuid.UUID `json:"cliente_id"`
	Placa     string    `json:"placa"`
	Marca     string    `json:"marca"`
	Modelo    string    `json:"modelo"`
	Ano       int       `json:"ano"`
	CriadoEm  time.Time `json:"criado_em"`
}

// --- Serviço ---

type CriarServicoRequest struct {
	Nome         string          `json:"nome" binding:"required"`
	Descricao    string          `json:"descricao"`
	ValorMaoObra decimal.Decimal `json:"valor_mao_obra" binding:"required"`
}

type ServicoResponse struct {
	ID           uuid.UUID       `json:"id"`
	Nome         string          `json:"nome"`
	Descricao    string          `json:"descricao"`
	ValorMaoObra decimal.Decimal `json:"valor_mao_obra"`
	Ativo        bool            `json:"ativo"`
}

// --- Peça ---

type CriarPecaRequest struct {
	Nome              string          `json:"nome" binding:"required"`
	Codigo            string          `json:"codigo" binding:"required"`
	ValorUnitario     decimal.Decimal `json:"valor_unitario" binding:"required"`
	QuantidadeEstoque int             `json:"quantidade_estoque" binding:"required,min=0"`
}

type PecaResponse struct {
	ID                uuid.UUID       `json:"id"`
	Nome              string          `json:"nome"`
	Codigo            string          `json:"codigo"`
	ValorUnitario     decimal.Decimal `json:"valor_unitario"`
	QuantidadeEstoque int             `json:"quantidade_estoque"`
	Ativo             bool            `json:"ativo"`
}

// --- Ordem de Serviço ---

type CriarOSRequest struct {
	ClienteID   uuid.UUID `json:"cliente_id" binding:"required"`
	VeiculoID   uuid.UUID `json:"veiculo_id" binding:"required"`
	Observacoes string    `json:"observacoes"`
}

type AdicionarServicoOSRequest struct {
	ServicoID uuid.UUID `json:"servico_id" binding:"required"`
}

type AdicionarPecaOSRequest struct {
	PecaID     uuid.UUID `json:"peca_id" binding:"required"`
	Quantidade int       `json:"quantidade" binding:"required,min=1"`
}

type AtualizarStatusOSRequest struct {
	Status string `json:"status" binding:"required"`
}

type ItemServicoOSResponse struct {
	ServicoID    uuid.UUID       `json:"servico_id"`
	Nome         string          `json:"nome"`
	ValorMaoObra decimal.Decimal `json:"valor_mao_obra"`
}

type ItemPecaOSResponse struct {
	PecaID        uuid.UUID       `json:"peca_id"`
	Nome          string          `json:"nome"`
	Quantidade    int             `json:"quantidade"`
	ValorUnitario decimal.Decimal `json:"valor_unitario"`
}

type OrdemServicoResponse struct {
	ID           uuid.UUID               `json:"id"`
	ClienteID    uuid.UUID               `json:"cliente_id"`
	VeiculoID    uuid.UUID               `json:"veiculo_id"`
	Status       string                  `json:"status"`
	ValorTotal   decimal.Decimal         `json:"valor_total"`
	Observacoes  string                  `json:"observacoes"`
	Servicos     []ItemServicoOSResponse `json:"servicos"`
	Pecas        []ItemPecaOSResponse    `json:"pecas"`
	CriadoEm     time.Time               `json:"criado_em"`
	EntregueEm   *time.Time              `json:"entregue_em,omitempty"`
	AtualizadoEm time.Time               `json:"atualizado_em"`
}

type StatusOSResponse struct {
	ID         uuid.UUID  `json:"id"`
	Status     string     `json:"status"`
	ValorTotal decimal.Decimal `json:"valor_total"`
	CriadoEm   time.Time  `json:"criado_em"`
	EntregueEm *time.Time `json:"entregue_em,omitempty"`
}

type TempoMedioResponse struct {
	TempoMedioHoras float64 `json:"tempo_medio_horas"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
