package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techchalleng/oficina/internal/domain/entity"
)

type mockPecaRepo struct {
	pecas map[uuid.UUID]*entity.Peca
}

func newMockPecaRepo() *mockPecaRepo {
	return &mockPecaRepo{pecas: make(map[uuid.UUID]*entity.Peca)}
}

func (m *mockPecaRepo) Criar(ctx context.Context, peca *entity.Peca) error {
	m.pecas[peca.ID] = peca
	return nil
}

func (m *mockPecaRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.Peca, error) {
	p, ok := m.pecas[id]
	if !ok {
		return nil, assert.AnError
	}
	return p, nil
}

func (m *mockPecaRepo) Listar(ctx context.Context) ([]*entity.Peca, error) { return nil, nil }
func (m *mockPecaRepo) Atualizar(ctx context.Context, peca *entity.Peca) error {
	m.pecas[peca.ID] = peca
	return nil
}
func (m *mockPecaRepo) Excluir(ctx context.Context, id uuid.UUID) error { return nil }

func TestGestorEstoque_VerificarDisponibilidade_Sucesso(t *testing.T) {
	repo := newMockPecaRepo()
	pecaID := uuid.New()
	repo.pecas[pecaID] = entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 10)

	gestor := NewGestorEstoque(repo)
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 5}}

	err := gestor.VerificarDisponibilidade(context.Background(), os)
	require.NoError(t, err)
}

func TestGestorEstoque_VerificarDisponibilidade_Insuficiente(t *testing.T) {
	repo := newMockPecaRepo()
	pecaID := uuid.New()
	repo.pecas[pecaID] = entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 2)

	gestor := NewGestorEstoque(repo)
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 5}}

	err := gestor.VerificarDisponibilidade(context.Background(), os)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "estoque insuficiente")
}

func TestGestorEstoque_BaixarEstoque(t *testing.T) {
	repo := newMockPecaRepo()
	pecaID := uuid.New()
	repo.pecas[pecaID] = entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 10)

	gestor := NewGestorEstoque(repo)
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 3}}

	err := gestor.BaixarEstoque(context.Background(), os)
	require.NoError(t, err)
	assert.Equal(t, 7, repo.pecas[pecaID].QuantidadeEstoque)
}

func TestGestorEstoque_BaixarEstoque_Insuficiente(t *testing.T) {
	repo := newMockPecaRepo()
	pecaID := uuid.New()
	repo.pecas[pecaID] = entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 1)

	gestor := NewGestorEstoque(repo)
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 5}}

	err := gestor.BaixarEstoque(context.Background(), os)
	require.Error(t, err)
}
