package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techchalleng/oficina/internal/domain/entity"
)

func TestGestorEstoque_VerificarDisponibilidade_Sucesso(t *testing.T) {
	pecaID := uuid.New()
	peca := entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 10)
	pecas := map[uuid.UUID]*entity.Peca{pecaID: peca}

	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 5}}

	err := NewGestorEstoque().VerificarDisponibilidade(os, pecas)
	require.NoError(t, err)
}

func TestGestorEstoque_VerificarDisponibilidade_Insuficiente(t *testing.T) {
	pecaID := uuid.New()
	peca := entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 2)
	pecas := map[uuid.UUID]*entity.Peca{pecaID: peca}

	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 5}}

	err := NewGestorEstoque().VerificarDisponibilidade(os, pecas)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "estoque insuficiente")
}

func TestGestorEstoque_BaixarEstoque(t *testing.T) {
	pecaID := uuid.New()
	peca := entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 10)
	pecas := map[uuid.UUID]*entity.Peca{pecaID: peca}

	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 3}}

	err := NewGestorEstoque().BaixarEstoque(os, pecas)
	require.NoError(t, err)
	assert.Equal(t, 7, peca.QuantidadeEstoque)
}

func TestGestorEstoque_BaixarEstoque_Insuficiente(t *testing.T) {
	pecaID := uuid.New()
	peca := entity.NewPeca("Filtro", "FLT001", decimal.NewFromFloat(50), 1)
	pecas := map[uuid.UUID]*entity.Peca{pecaID: peca}

	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: pecaID, Quantidade: 5}}

	err := NewGestorEstoque().BaixarEstoque(os, pecas)
	require.Error(t, err)
}

func TestGestorEstoque_PecaAusente(t *testing.T) {
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Pecas = []entity.ItemPecaOS{{PecaID: uuid.New(), Quantidade: 1}}

	err := NewGestorEstoque().VerificarDisponibilidade(os, map[uuid.UUID]*entity.Peca{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "não encontrada")
}
