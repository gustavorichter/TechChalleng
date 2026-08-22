package dto

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

func TestMappers(t *testing.T) {
	email, _ := valobj.NewEmail("test@example.com")
	cpf, _ := valobj.NewCPF("529.982.247-25")
	placa, _ := valobj.NewPlaca("ABC1234")

	cliente := entity.NewCliente("João", email, "11999999999")
	cliente.DefinirCPF(cpf)
	cr := ToClienteResponse(cliente)
	assert.Equal(t, cliente.Nome, cr.Nome)
	assert.Equal(t, cpf.String(), cr.CPF)

	veiculo := entity.NewVeiculo(cliente.ID, placa, "Fiat", "Uno", 2020)
	vr := ToVeiculoResponse(veiculo)
	assert.Equal(t, placa.String(), vr.Placa)

	servico := entity.NewServico("Troca", "desc", decimal.NewFromFloat(100))
	sr := ToServicoResponse(servico)
	assert.True(t, sr.ValorMaoObra.Equal(decimal.NewFromFloat(100)))

	peca := entity.NewPeca("Filtro", "FLT", decimal.NewFromFloat(50), 10)
	pr := ToPecaResponse(peca)
	assert.Equal(t, 10, pr.QuantidadeEstoque)

	now := time.Now().UTC()
	os := entity.NewOrdemServico(cliente.ID, veiculo.ID, "obs")
	os.Servicos = []entity.ItemServicoOS{{ServicoID: servico.ID, Nome: "Troca", ValorMaoObra: decimal.NewFromFloat(100)}}
	os.Pecas = []entity.ItemPecaOS{{PecaID: peca.ID, Nome: "Filtro", Quantidade: 1, ValorUnitario: decimal.NewFromFloat(50)}}
	os.ValorTotal = decimal.NewFromFloat(150)
	os.EntregueEm = &now

	or := ToOrdemServicoResponse(os)
	assert.Len(t, or.Servicos, 1)
	assert.Len(t, or.Pecas, 1)

	sr2 := ToStatusOSResponse(os)
	assert.Equal(t, os.ID, sr2.ID)
	assert.NotNil(t, sr2.EntregueEm)
}
