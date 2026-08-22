package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/techchalleng/oficina/internal/domain/entity"
)

func TestCalculadoraOrcamento_Calcular(t *testing.T) {
	calc := NewCalculadoraOrcamento()
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")

	os.Servicos = []entity.ItemServicoOS{
		{ServicoID: uuid.New(), Nome: "Troca de óleo", ValorMaoObra: decimal.NewFromFloat(150.00)},
		{ServicoID: uuid.New(), Nome: "Alinhamento", ValorMaoObra: decimal.NewFromFloat(80.00)},
	}
	os.Pecas = []entity.ItemPecaOS{
		{PecaID: uuid.New(), Nome: "Filtro de óleo", Quantidade: 1, ValorUnitario: decimal.NewFromFloat(45.00)},
		{PecaID: uuid.New(), Nome: "Óleo 5W30", Quantidade: 4, ValorUnitario: decimal.NewFromFloat(35.00)},
	}

	total := calc.Calcular(os)
	// 150 + 80 + 45 + (4*35) = 415
	assert.True(t, total.Equal(decimal.NewFromFloat(415.00)))
}

func TestCalculadoraOrcamento_OS_Vazia(t *testing.T) {
	calc := NewCalculadoraOrcamento()
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	assert.True(t, calc.Calcular(os).IsZero())
}

func TestCalculadoraOrcamento_Aplicar(t *testing.T) {
	calc := NewCalculadoraOrcamento()
	os := entity.NewOrdemServico(uuid.New(), uuid.New(), "")
	os.Servicos = []entity.ItemServicoOS{
		{ValorMaoObra: decimal.NewFromFloat(100.00)},
	}
	calc.Aplicar(os)
	assert.True(t, os.ValorTotal.Equal(decimal.NewFromFloat(100.00)))
}
