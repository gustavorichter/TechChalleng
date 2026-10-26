package entity

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

// ItemServicoOS vincula um serviço a uma ordem de serviço.
type ItemServicoOS struct {
	ServicoID    uuid.UUID
	Nome         string
	ValorMaoObra decimal.Decimal
}

// ItemPecaOS vincula uma peça a uma ordem de serviço com quantidade.
type ItemPecaOS struct {
	PecaID        uuid.UUID
	Nome          string
	Quantidade    int
	ValorUnitario decimal.Decimal
}

// OrdemServico é o Aggregate Root do domínio de atendimento.
type OrdemServico struct {
	ID                uuid.UUID
	ClienteID         uuid.UUID
	VeiculoID         uuid.UUID
	Status            valobj.StatusOS
	ValorTotal        decimal.Decimal
	Observacoes       string
	Servicos          []ItemServicoOS
	Pecas             []ItemPecaOS
	OrcamentoAprovado *bool
	MotivoRecusa      string
	CriadoEm          time.Time
	EntregueEm        *time.Time
	ExcluidaEm        *time.Time
	AtualizadoEm      time.Time
}

// NewOrdemServico cria uma nova OS no status inicial "Recebida".
func NewOrdemServico(clienteID, veiculoID uuid.UUID, observacoes string) *OrdemServico {
	now := time.Now().UTC()
	return &OrdemServico{
		ID:           uuid.New(),
		ClienteID:    clienteID,
		VeiculoID:    veiculoID,
		Status:       valobj.StatusRecebida,
		ValorTotal:   decimal.Zero,
		Observacoes:  observacoes,
		Servicos:     []ItemServicoOS{},
		Pecas:        []ItemPecaOS{},
		CriadoEm:     now,
		AtualizadoEm: now,
	}
}

// AdicionarServico adiciona um serviço à OS.
func (os *OrdemServico) AdicionarServico(item ItemServicoOS) error {
	if os.Status != valobj.StatusRecebida && os.Status != valobj.StatusEmDiagnostico {
		return errors.New("serviços só podem ser adicionados nos status Recebida ou Em diagnóstico")
	}
	os.Servicos = append(os.Servicos, item)
	os.AtualizadoEm = time.Now().UTC()
	return nil
}

// AdicionarPeca adiciona uma peça à OS.
func (os *OrdemServico) AdicionarPeca(item ItemPecaOS) error {
	if os.Status != valobj.StatusRecebida && os.Status != valobj.StatusEmDiagnostico {
		return errors.New("peças só podem ser adicionadas nos status Recebida ou Em diagnóstico")
	}
	if item.Quantidade <= 0 {
		return errors.New("quantidade deve ser positiva")
	}
	os.Pecas = append(os.Pecas, item)
	os.AtualizadoEm = time.Now().UTC()
	return nil
}

// AtualizarStatus transiciona o status da OS respeitando a máquina de estados.
func (os *OrdemServico) AtualizarStatus(novoStatus valobj.StatusOS) error {
	proximo, err := os.Status.TransicionarPara(novoStatus)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	os.Status = proximo
	os.AtualizadoEm = now

	if proximo == valobj.StatusEntregue {
		os.EntregueEm = &now
	}
	if proximo.OcultaDaListagem() && os.ExcluidaEm == nil {
		os.ExcluidaEm = &now
	}
	return nil
}

// AprovarOrcamento confirma o orçamento e inicia a execução.
func (os *OrdemServico) AprovarOrcamento() error {
	if os.Status != valobj.StatusAguardandoAprovacao {
		return errors.New("orçamento só pode ser aprovado quando a OS está aguardando aprovação")
	}
	if len(os.Servicos) == 0 && len(os.Pecas) == 0 {
		return errors.New("não é possível aprovar um orçamento sem serviços ou peças")
	}
	aprovado := true
	os.OrcamentoAprovado = &aprovado
	os.MotivoRecusa = ""
	return os.AtualizarStatus(valobj.StatusEmExecucao)
}

// RecusarOrcamento devolve a OS ao diagnóstico para revisão do orçamento.
func (os *OrdemServico) RecusarOrcamento(motivo string) error {
	if os.Status != valobj.StatusAguardandoAprovacao {
		return errors.New("orçamento só pode ser recusado quando a OS está aguardando aprovação")
	}
	motivo = strings.TrimSpace(motivo)
	if motivo == "" {
		return errors.New("motivo da recusa é obrigatório")
	}
	recusado := false
	os.OrcamentoAprovado = &recusado
	os.MotivoRecusa = motivo
	os.Status = valobj.StatusEmDiagnostico
	os.AtualizadoEm = time.Now().UTC()
	return nil
}

// VisivelNaListagem indica se a OS permanece na fila operacional.
func (os *OrdemServico) VisivelNaListagem() bool {
	if os == nil || os.ExcluidaEm != nil {
		return false
	}
	return !os.Status.OcultaDaListagem()
}

// FiltrarEOrdenarListagem oculta OS finalizadas ou entregues e ordena a fila.
// Prioridade: Em execução, Aguardando aprovação, Em diagnóstico, Recebida.
// Dentro do mesmo status, as mais antigas vêm primeiro.
func FiltrarEOrdenarListagem(ordens []*OrdemServico) []*OrdemServico {
	ativas := make([]*OrdemServico, 0, len(ordens))
	for _, os := range ordens {
		if os.VisivelNaListagem() {
			ativas = append(ativas, os)
		}
	}
	sort.SliceStable(ativas, func(i, j int) bool {
		pi := ativas[i].Status.PrioridadeListagem()
		pj := ativas[j].Status.PrioridadeListagem()
		if pi != pj {
			return pi < pj
		}
		return ativas[i].CriadoEm.Before(ativas[j].CriadoEm)
	})
	return ativas
}

// AvancarStatus avança automaticamente para o próximo status na sequência.
func (os *OrdemServico) AvancarStatus() error {
	proximo, ok := os.Status.ProximoStatus()
	if !ok {
		return errors.New("não há próximo status disponível")
	}
	return os.AtualizarStatus(proximo)
}

// TempoExecucao retorna a duração entre criação e entrega, se entregue.
func (os *OrdemServico) TempoExecucao() (time.Duration, bool) {
	if os.EntregueEm == nil {
		return 0, false
	}
	return os.EntregueEm.Sub(os.CriadoEm), true
}
