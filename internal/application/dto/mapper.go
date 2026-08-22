package dto

import (
	"github.com/techchalleng/oficina/internal/domain/entity"
)

func ToClienteResponse(c *entity.Cliente) ClienteResponse {
	resp := ClienteResponse{
		ID:       c.ID,
		Nome:     c.Nome,
		Email:    c.Email.String(),
		Telefone: c.Telefone,
		CriadoEm: c.CriadoEm,
	}
	if c.CPF != nil {
		resp.CPF = c.CPF.String()
	}
	if c.CNPJ != nil {
		resp.CNPJ = c.CNPJ.String()
	}
	return resp
}

func ToVeiculoResponse(v *entity.Veiculo) VeiculoResponse {
	return VeiculoResponse{
		ID:        v.ID,
		ClienteID: v.ClienteID,
		Placa:     v.Placa.String(),
		Marca:     v.Marca,
		Modelo:    v.Modelo,
		Ano:       v.Ano,
		CriadoEm:  v.CriadoEm,
	}
}

func ToServicoResponse(s *entity.Servico) ServicoResponse {
	return ServicoResponse{
		ID:           s.ID,
		Nome:         s.Nome,
		Descricao:    s.Descricao,
		ValorMaoObra: s.ValorMaoObra,
		Ativo:        s.Ativo,
	}
}

func ToPecaResponse(p *entity.Peca) PecaResponse {
	return PecaResponse{
		ID:                p.ID,
		Nome:              p.Nome,
		Codigo:            p.Codigo,
		ValorUnitario:     p.ValorUnitario,
		QuantidadeEstoque: p.QuantidadeEstoque,
		Ativo:             p.Ativo,
	}
}

func ToOrdemServicoResponse(os *entity.OrdemServico) OrdemServicoResponse {
	servicos := make([]ItemServicoOSResponse, len(os.Servicos))
	for i, s := range os.Servicos {
		servicos[i] = ItemServicoOSResponse{
			ServicoID:    s.ServicoID,
			Nome:         s.Nome,
			ValorMaoObra: s.ValorMaoObra,
		}
	}
	pecas := make([]ItemPecaOSResponse, len(os.Pecas))
	for i, p := range os.Pecas {
		pecas[i] = ItemPecaOSResponse{
			PecaID:        p.PecaID,
			Nome:          p.Nome,
			Quantidade:    p.Quantidade,
			ValorUnitario: p.ValorUnitario,
		}
	}
	return OrdemServicoResponse{
		ID:           os.ID,
		ClienteID:    os.ClienteID,
		VeiculoID:    os.VeiculoID,
		Status:       string(os.Status),
		ValorTotal:   os.ValorTotal,
		Observacoes:  os.Observacoes,
		Servicos:     servicos,
		Pecas:        pecas,
		CriadoEm:     os.CriadoEm,
		EntregueEm:   os.EntregueEm,
		AtualizadoEm: os.AtualizadoEm,
	}
}

func ToStatusOSResponse(os *entity.OrdemServico) StatusOSResponse {
	return StatusOSResponse{
		ID:         os.ID,
		Status:     string(os.Status),
		ValorTotal: os.ValorTotal,
		CriadoEm:   os.CriadoEm,
		EntregueEm: os.EntregueEm,
	}
}
