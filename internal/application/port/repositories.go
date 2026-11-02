package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/domain/entity"
)

// ClienteRepository é a porta de persistência de clientes.
type ClienteRepository interface {
	Criar(ctx context.Context, cliente *entity.Cliente) error
	BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Cliente, error)
	Listar(ctx context.Context) ([]*entity.Cliente, error)
	Atualizar(ctx context.Context, cliente *entity.Cliente) error
	Excluir(ctx context.Context, id uuid.UUID) error
}

// VeiculoRepository é a porta de persistência de veículos.
type VeiculoRepository interface {
	Criar(ctx context.Context, veiculo *entity.Veiculo) error
	BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Veiculo, error)
	ListarPorCliente(ctx context.Context, clienteID uuid.UUID) ([]*entity.Veiculo, error)
	Listar(ctx context.Context) ([]*entity.Veiculo, error)
	Atualizar(ctx context.Context, veiculo *entity.Veiculo) error
	Excluir(ctx context.Context, id uuid.UUID) error
}

// ServicoRepository é a porta de persistência do catálogo de serviços.
type ServicoRepository interface {
	Criar(ctx context.Context, servico *entity.Servico) error
	BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Servico, error)
	Listar(ctx context.Context) ([]*entity.Servico, error)
	Atualizar(ctx context.Context, servico *entity.Servico) error
	Excluir(ctx context.Context, id uuid.UUID) error
}

// PecaRepository é a porta de persistência do estoque de peças.
type PecaRepository interface {
	Criar(ctx context.Context, peca *entity.Peca) error
	BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Peca, error)
	Listar(ctx context.Context) ([]*entity.Peca, error)
	Atualizar(ctx context.Context, peca *entity.Peca) error
	Excluir(ctx context.Context, id uuid.UUID) error
}

// OrdemServicoRepository é a porta de persistência do agregado Ordem de Serviço.
type OrdemServicoRepository interface {
	Criar(ctx context.Context, os *entity.OrdemServico) error
	BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.OrdemServico, error)
	Listar(ctx context.Context) ([]*entity.OrdemServico, error)
	Atualizar(ctx context.Context, os *entity.OrdemServico) error
	TempoMedioExecucao(ctx context.Context) (float64, error)
}
