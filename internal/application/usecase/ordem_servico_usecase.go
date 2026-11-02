package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/application/port"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/service"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

// OrdemServicoUseCase orquestra o ciclo de vida das Ordens de Serviço.
type OrdemServicoUseCase struct {
	osRepo      port.OrdemServicoRepository
	clienteRepo port.ClienteRepository
	veiculoRepo port.VeiculoRepository
	servicoRepo port.ServicoRepository
	pecaRepo    port.PecaRepository
	notificador port.NotificadorStatus
	calculadora *service.CalculadoraOrcamento
	gestor      *service.GestorEstoque
}

func NewOrdemServicoUseCase(
	osRepo port.OrdemServicoRepository,
	clienteRepo port.ClienteRepository,
	veiculoRepo port.VeiculoRepository,
	servicoRepo port.ServicoRepository,
	pecaRepo port.PecaRepository,
	notificador port.NotificadorStatus,
) *OrdemServicoUseCase {
	return &OrdemServicoUseCase{
		osRepo:      osRepo,
		clienteRepo: clienteRepo,
		veiculoRepo: veiculoRepo,
		servicoRepo: servicoRepo,
		pecaRepo:    pecaRepo,
		notificador: notificador,
		calculadora: service.NewCalculadoraOrcamento(),
		gestor:      service.NewGestorEstoque(),
	}
}

func (uc *OrdemServicoUseCase) Criar(ctx context.Context, req dto.CriarOSRequest) (*dto.OrdemServicoResponse, error) {
	cliente, err := uc.clienteRepo.BuscarPorID(ctx, req.ClienteID)
	if err != nil {
		return nil, fmt.Errorf("cliente não encontrado: %w", err)
	}

	veiculo, err := uc.veiculoRepo.BuscarPorID(ctx, req.VeiculoID)
	if err != nil {
		return nil, fmt.Errorf("veículo não encontrado: %w", err)
	}
	if veiculo.ClienteID != req.ClienteID {
		return nil, errors.New("veículo não pertence ao cliente informado")
	}

	os := entity.NewOrdemServico(req.ClienteID, req.VeiculoID, req.Observacoes)
	for _, item := range req.Servicos {
		servico, err := uc.servicoRepo.BuscarPorID(ctx, item.ServicoID)
		if err != nil {
			return nil, fmt.Errorf("serviço não encontrado: %w", err)
		}
		if err := os.AdicionarServico(entity.ItemServicoOS{
			ServicoID:    servico.ID,
			Nome:         servico.Nome,
			ValorMaoObra: servico.ValorMaoObra,
		}); err != nil {
			return nil, err
		}
	}
	for _, item := range req.Pecas {
		peca, err := uc.pecaRepo.BuscarPorID(ctx, item.PecaID)
		if err != nil {
			return nil, fmt.Errorf("peça não encontrada: %w", err)
		}
		if !peca.TemSaldo(item.Quantidade) {
			return nil, fmt.Errorf("estoque insuficiente para peça '%s': disponível %d, solicitado %d",
				peca.Nome, peca.QuantidadeEstoque, item.Quantidade)
		}
		if err := os.AdicionarPeca(entity.ItemPecaOS{
			PecaID:        peca.ID,
			Nome:          peca.Nome,
			Quantidade:    item.Quantidade,
			ValorUnitario: peca.ValorUnitario,
		}); err != nil {
			return nil, err
		}
	}

	uc.calculadora.Aplicar(os)
	if err := uc.osRepo.Criar(ctx, os); err != nil {
		return nil, err
	}
	uc.notificar(ctx, cliente, os, "")

	resp := dto.ToOrdemServicoResponse(os)
	return &resp, nil
}

func (uc *OrdemServicoUseCase) AdicionarServico(ctx context.Context, osID uuid.UUID, req dto.AdicionarServicoOSRequest) (*dto.OrdemServicoResponse, error) {
	os, err := uc.osRepo.BuscarPorID(ctx, osID)
	if err != nil {
		return nil, err
	}

	servico, err := uc.servicoRepo.BuscarPorID(ctx, req.ServicoID)
	if err != nil {
		return nil, fmt.Errorf("serviço não encontrado: %w", err)
	}

	if err := os.AdicionarServico(entity.ItemServicoOS{
		ServicoID:    servico.ID,
		Nome:         servico.Nome,
		ValorMaoObra: servico.ValorMaoObra,
	}); err != nil {
		return nil, err
	}

	uc.calculadora.Aplicar(os)
	if err := uc.osRepo.Atualizar(ctx, os); err != nil {
		return nil, err
	}

	resp := dto.ToOrdemServicoResponse(os)
	return &resp, nil
}

func (uc *OrdemServicoUseCase) AdicionarPeca(ctx context.Context, osID uuid.UUID, req dto.AdicionarPecaOSRequest) (*dto.OrdemServicoResponse, error) {
	os, err := uc.osRepo.BuscarPorID(ctx, osID)
	if err != nil {
		return nil, err
	}

	peca, err := uc.pecaRepo.BuscarPorID(ctx, req.PecaID)
	if err != nil {
		return nil, fmt.Errorf("peça não encontrada: %w", err)
	}

	if !peca.TemSaldo(req.Quantidade) {
		return nil, fmt.Errorf("estoque insuficiente para peça '%s': disponível %d, solicitado %d",
			peca.Nome, peca.QuantidadeEstoque, req.Quantidade)
	}

	if err := os.AdicionarPeca(entity.ItemPecaOS{
		PecaID:        peca.ID,
		Nome:          peca.Nome,
		Quantidade:    req.Quantidade,
		ValorUnitario: peca.ValorUnitario,
	}); err != nil {
		return nil, err
	}

	uc.calculadora.Aplicar(os)
	if err := uc.osRepo.Atualizar(ctx, os); err != nil {
		return nil, err
	}

	resp := dto.ToOrdemServicoResponse(os)
	return &resp, nil
}

func (uc *OrdemServicoUseCase) AtualizarStatus(ctx context.Context, osID uuid.UUID, req dto.AtualizarStatusOSRequest) (*dto.OrdemServicoResponse, error) {
	os, err := uc.osRepo.BuscarPorID(ctx, osID)
	if err != nil {
		return nil, err
	}
	if os.Status == valobj.StatusAguardandoAprovacao {
		return nil, errors.New("transição a partir de 'Aguardando aprovação' exige decisão explícita do orçamento")
	}

	novoStatus, err := valobj.NewStatusOS(req.Status)
	if err != nil {
		return nil, err
	}

	if err := uc.aplicarTransicao(ctx, os, novoStatus); err != nil {
		return nil, err
	}

	resp := dto.ToOrdemServicoResponse(os)
	return &resp, nil
}

func (uc *OrdemServicoUseCase) AvancarStatus(ctx context.Context, osID uuid.UUID) (*dto.OrdemServicoResponse, error) {
	os, err := uc.osRepo.BuscarPorID(ctx, osID)
	if err != nil {
		return nil, err
	}
	if os.Status == valobj.StatusAguardandoAprovacao {
		return nil, errors.New("transição a partir de 'Aguardando aprovação' exige decisão explícita do orçamento")
	}

	proximo, ok := os.Status.ProximoStatus()
	if !ok {
		return nil, errors.New("não há próximo status disponível")
	}
	if err := uc.aplicarTransicao(ctx, os, proximo); err != nil {
		return nil, err
	}

	resp := dto.ToOrdemServicoResponse(os)
	return &resp, nil
}

// DecidirOrcamento aprova o orçamento e inicia a execução, ou recusa e devolve a OS ao diagnóstico.
func (uc *OrdemServicoUseCase) DecidirOrcamento(ctx context.Context, osID uuid.UUID, req dto.DecidirOrcamentoRequest) (*dto.OrdemServicoResponse, error) {
	os, err := uc.osRepo.BuscarPorID(ctx, osID)
	if err != nil {
		return nil, err
	}
	anterior := string(os.Status)

	if req.Aprovado {
		if err := os.AprovarOrcamento(); err != nil {
			return nil, err
		}
	} else if err := os.RecusarOrcamento(req.Motivo); err != nil {
		return nil, err
	}

	if err := uc.osRepo.Atualizar(ctx, os); err != nil {
		return nil, err
	}
	uc.notificarOS(ctx, os, anterior)

	resp := dto.ToOrdemServicoResponse(os)
	return &resp, nil
}

func (uc *OrdemServicoUseCase) BuscarPorID(ctx context.Context, id uuid.UUID) (*dto.OrdemServicoResponse, error) {
	os, err := uc.osRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToOrdemServicoResponse(os)
	return &resp, nil
}

func (uc *OrdemServicoUseCase) AcompanharStatus(ctx context.Context, id uuid.UUID) (*dto.StatusOSResponse, error) {
	os, err := uc.osRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToStatusOSResponse(os)
	return &resp, nil
}

func (uc *OrdemServicoUseCase) Listar(ctx context.Context) ([]dto.OrdemServicoResponse, error) {
	ordens, err := uc.osRepo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	ordens = entity.FiltrarEOrdenarListagem(ordens)
	result := make([]dto.OrdemServicoResponse, len(ordens))
	for i, os := range ordens {
		result[i] = dto.ToOrdemServicoResponse(os)
	}
	return result, nil
}

func (uc *OrdemServicoUseCase) TempoMedioExecucao(ctx context.Context) (*dto.TempoMedioResponse, error) {
	media, err := uc.osRepo.TempoMedioExecucao(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.TempoMedioResponse{TempoMedioHoras: media}, nil
}

func (uc *OrdemServicoUseCase) aplicarTransicao(ctx context.Context, os *entity.OrdemServico, novoStatus valobj.StatusOS) error {
	anterior := string(os.Status)
	var pecas map[uuid.UUID]*entity.Peca
	if novoStatus == valobj.StatusFinalizada {
		carregadas, err := uc.carregarPecas(ctx, os)
		if err != nil {
			return err
		}
		pecas = carregadas
		if err := uc.gestor.VerificarDisponibilidade(os, pecas); err != nil {
			return err
		}
	}

	if err := os.AtualizarStatus(novoStatus); err != nil {
		return err
	}

	if novoStatus == valobj.StatusFinalizada {
		if err := uc.gestor.BaixarEstoque(os, pecas); err != nil {
			return err
		}
		for _, peca := range pecas {
			if err := uc.pecaRepo.Atualizar(ctx, peca); err != nil {
				return fmt.Errorf("falha ao persistir estoque da peça '%s': %w", peca.Nome, err)
			}
		}
	}

	if err := uc.osRepo.Atualizar(ctx, os); err != nil {
		return err
	}
	uc.notificarOS(ctx, os, anterior)
	return nil
}

func (uc *OrdemServicoUseCase) carregarPecas(ctx context.Context, os *entity.OrdemServico) (map[uuid.UUID]*entity.Peca, error) {
	pecas := make(map[uuid.UUID]*entity.Peca, len(os.Pecas))
	for _, item := range os.Pecas {
		if _, ok := pecas[item.PecaID]; ok {
			continue
		}
		peca, err := uc.pecaRepo.BuscarPorID(ctx, item.PecaID)
		if err != nil {
			return nil, fmt.Errorf("peça %s não encontrada: %w", item.PecaID, err)
		}
		pecas[item.PecaID] = peca
	}
	return pecas, nil
}

func (uc *OrdemServicoUseCase) notificarOS(ctx context.Context, os *entity.OrdemServico, statusAnterior string) {
	cliente, err := uc.clienteRepo.BuscarPorID(ctx, os.ClienteID)
	if err != nil {
		log.Printf("falha ao carregar cliente %s para notificação da OS %s: %v", os.ClienteID, os.ID, err)
		return
	}
	uc.notificar(ctx, cliente, os, statusAnterior)
}

func (uc *OrdemServicoUseCase) notificar(ctx context.Context, cliente *entity.Cliente, os *entity.OrdemServico, statusAnterior string) {
	if uc.notificador == nil || statusAnterior == string(os.Status) {
		return
	}
	email := ""
	if cliente != nil {
		email = cliente.Email.String()
	}
	err := uc.notificador.NotificarAtualizacaoStatus(ctx, port.NotificacaoStatus{
		ClienteID:      os.ClienteID,
		Email:          email,
		OrdemServicoID: os.ID,
		StatusAnterior: statusAnterior,
		StatusNovo:     string(os.Status),
	})
	if err != nil {
		log.Printf("falha ao notificar cliente sobre a OS %s: %v", os.ID, err)
	}
}
