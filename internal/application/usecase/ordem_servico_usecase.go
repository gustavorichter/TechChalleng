package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/application/dto"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/repository"
	"github.com/techchalleng/oficina/internal/domain/service"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

// OrdemServicoUseCase orquestra o ciclo de vida das Ordens de Serviço.
type OrdemServicoUseCase struct {
	osRepo      repository.OrdemServicoRepository
	clienteRepo repository.ClienteRepository
	veiculoRepo repository.VeiculoRepository
	servicoRepo repository.ServicoRepository
	pecaRepo    repository.PecaRepository
	calculadora *service.CalculadoraOrcamento
	gestor      *service.GestorEstoque
}

func NewOrdemServicoUseCase(
	osRepo repository.OrdemServicoRepository,
	clienteRepo repository.ClienteRepository,
	veiculoRepo repository.VeiculoRepository,
	servicoRepo repository.ServicoRepository,
	pecaRepo repository.PecaRepository,
) *OrdemServicoUseCase {
	return &OrdemServicoUseCase{
		osRepo:      osRepo,
		clienteRepo: clienteRepo,
		veiculoRepo: veiculoRepo,
		servicoRepo: servicoRepo,
		pecaRepo:    pecaRepo,
		calculadora: service.NewCalculadoraOrcamento(),
		gestor:      service.NewGestorEstoque(pecaRepo),
	}
}

func (uc *OrdemServicoUseCase) Criar(ctx context.Context, req dto.CriarOSRequest) (*dto.OrdemServicoResponse, error) {
	if _, err := uc.clienteRepo.BuscarPorID(ctx, req.ClienteID); err != nil {
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
	if err := uc.osRepo.Criar(ctx, os); err != nil {
		return nil, err
	}

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

	item := entity.ItemServicoOS{
		ServicoID:    servico.ID,
		Nome:         servico.Nome,
		ValorMaoObra: servico.ValorMaoObra,
	}

	if err := os.AdicionarServico(item); err != nil {
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

	item := entity.ItemPecaOS{
		PecaID:        peca.ID,
		Nome:          peca.Nome,
		Quantidade:    req.Quantidade,
		ValorUnitario: peca.ValorUnitario,
	}

	if err := os.AdicionarPeca(item); err != nil {
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

	novoStatus, err := valobj.NewStatusOS(req.Status)
	if err != nil {
		return nil, err
	}

	// Ao finalizar, verificar estoque e dar baixa
	if novoStatus == valobj.StatusFinalizada {
		if err := uc.gestor.VerificarDisponibilidade(ctx, os); err != nil {
			return nil, err
		}
	}

	if err := os.AtualizarStatus(novoStatus); err != nil {
		return nil, err
	}

	if novoStatus == valobj.StatusFinalizada {
		if err := uc.gestor.BaixarEstoque(ctx, os); err != nil {
			return nil, err
		}
	}

	if err := uc.osRepo.Atualizar(ctx, os); err != nil {
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

	proximo, ok := os.Status.ProximoStatus()
	if !ok {
		return nil, errors.New("não há próximo status disponível")
	}

	return uc.AtualizarStatus(ctx, osID, dto.AtualizarStatusOSRequest{Status: string(proximo)})
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
