package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

func TestClienteUseCase_Criar_ComCPF(t *testing.T) {
	repo := newMockClienteRepo()
	uc := NewClienteUseCase(repo)

	resp, err := uc.Criar(context.Background(), dto.CriarClienteRequest{
		Nome:     "Maria",
		CPF:      "529.982.247-25",
		Email:    "maria@email.com",
		Telefone: "11988887777",
	})
	require.NoError(t, err)
	assert.Equal(t, "Maria", resp.Nome)
	assert.NotEmpty(t, resp.CPF)
}

func TestClienteUseCase_Criar_SemDocumento(t *testing.T) {
	repo := newMockClienteRepo()
	uc := NewClienteUseCase(repo)

	_, err := uc.Criar(context.Background(), dto.CriarClienteRequest{
		Nome: "Maria", Email: "maria@email.com", Telefone: "11988887777",
	})
	require.Error(t, err)
}

func TestClienteUseCase_Criar_CPFInvalido(t *testing.T) {
	repo := newMockClienteRepo()
	uc := NewClienteUseCase(repo)

	_, err := uc.Criar(context.Background(), dto.CriarClienteRequest{
		Nome: "Maria", CPF: "111.111.111-11", Email: "maria@email.com", Telefone: "11988887777",
	})
	require.Error(t, err)
}

func TestClienteUseCase_Buscar_Listar_Excluir(t *testing.T) {
	repo := newMockClienteRepo()
	uc := NewClienteUseCase(repo)

	created, err := uc.Criar(context.Background(), dto.CriarClienteRequest{
		Nome: "João", CPF: "529.982.247-25", Email: "joao@email.com", Telefone: "11999999999",
	})
	require.NoError(t, err)

	found, err := uc.BuscarPorID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)

	list, err := uc.Listar(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 1)

	require.NoError(t, uc.Excluir(context.Background(), created.ID))
}

func TestPecaUseCase_Criar(t *testing.T) {
	repo := newMockPecaRepoUC()
	uc := NewPecaUseCase(repo)

	resp, err := uc.Criar(context.Background(), dto.CriarPecaRequest{
		Nome: "Pastilha", Codigo: "PST001", ValorUnitario: decimal.NewFromFloat(120), QuantidadeEstoque: 5,
	})
	require.NoError(t, err)
	assert.Equal(t, 5, resp.QuantidadeEstoque)
}

func TestPecaUseCase_ValorInvalido(t *testing.T) {
	repo := newMockPecaRepoUC()
	uc := NewPecaUseCase(repo)

	_, err := uc.Criar(context.Background(), dto.CriarPecaRequest{
		Nome: "Pastilha", Codigo: "PST001", ValorUnitario: decimal.Zero, QuantidadeEstoque: 5,
	})
	require.Error(t, err)
}

func TestVeiculoUseCase_Criar(t *testing.T) {
	clienteID := uuid.New()
	clienteRepo := newMockClienteRepo()
	clienteRepo.clientes[clienteID] = &entity.Cliente{ID: clienteID}
	veiculoRepo := newMockVeiculoRepo()
	uc := NewVeiculoUseCase(veiculoRepo, clienteRepo)

	resp, err := uc.Criar(context.Background(), dto.CriarVeiculoRequest{
		ClienteID: clienteID, Placa: "ABC1D23", Marca: "VW", Modelo: "Gol", Ano: 2022,
	})
	require.NoError(t, err)
	assert.Equal(t, "ABC1D23", resp.Placa)
}

func TestVeiculoUseCase_PlacaInvalida(t *testing.T) {
	clienteID := uuid.New()
	clienteRepo := newMockClienteRepo()
	clienteRepo.clientes[clienteID] = &entity.Cliente{ID: clienteID}
	uc := NewVeiculoUseCase(newMockVeiculoRepo(), clienteRepo)

	_, err := uc.Criar(context.Background(), dto.CriarVeiculoRequest{
		ClienteID: clienteID, Placa: "INVALID", Marca: "VW", Modelo: "Gol", Ano: 2022,
	})
	require.Error(t, err)
}

func TestServicoUseCase_Criar(t *testing.T) {
	repo := newMockServicoRepo()
	uc := NewServicoUseCase(repo)

	resp, err := uc.Criar(context.Background(), dto.CriarServicoRequest{
		Nome: "Alinhamento", Descricao: "Completo", ValorMaoObra: decimal.NewFromFloat(80),
	})
	require.NoError(t, err)
	assert.True(t, resp.Ativo)
}

func TestOrdemServicoUseCase_AcompanharStatus(t *testing.T) {
	uc, clienteID, veiculoID, _ := setupOrdemServicoUC()

	created, err := uc.Criar(context.Background(), dto.CriarOSRequest{ClienteID: clienteID, VeiculoID: veiculoID})
	require.NoError(t, err)

	status, err := uc.AcompanharStatus(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, string(valobj.StatusRecebida), status.Status)
}
