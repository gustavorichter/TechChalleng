package service

import (
	"github.com/shopspring/decimal"
	"github.com/techchalleng/oficina/internal/domain/entity"
)

// CalculadoraOrcamento calcula o valor total de uma Ordem de Serviço.
type CalculadoraOrcamento struct{}

func NewCalculadoraOrcamento() *CalculadoraOrcamento {
	return &CalculadoraOrcamento{}
}

// Calcular soma mão de obra dos serviços e valor unitário × quantidade das peças.
func (c *CalculadoraOrcamento) Calcular(os *entity.OrdemServico) decimal.Decimal {
	total := decimal.Zero

	for _, s := range os.Servicos {
		total = total.Add(s.ValorMaoObra)
	}

	for _, p := range os.Pecas {
		subtotal := p.ValorUnitario.Mul(decimal.NewFromInt(int64(p.Quantidade)))
		total = total.Add(subtotal)
	}

	return total
}

// Aplicar recalcula e persiste o valor total na entidade.
func (c *CalculadoraOrcamento) Aplicar(os *entity.OrdemServico) {
	os.ValorTotal = c.Calcular(os)
}
