package service

import (
	"context"
	"fmt"

	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/repository"
)

// GestorEstoque gerencia reservas e baixas de estoque vinculadas a Ordens de Serviço.
type GestorEstoque struct {
	pecaRepo repository.PecaRepository
}

func NewGestorEstoque(pecaRepo repository.PecaRepository) *GestorEstoque {
	return &GestorEstoque{pecaRepo: pecaRepo}
}

// VerificarDisponibilidade valida se há saldo para todas as peças da OS.
func (g *GestorEstoque) VerificarDisponibilidade(ctx context.Context, os *entity.OrdemServico) error {
	for _, item := range os.Pecas {
		peca, err := g.pecaRepo.BuscarPorID(ctx, item.PecaID)
		if err != nil {
			return fmt.Errorf("peça %s não encontrada: %w", item.PecaID, err)
		}
		if !peca.TemSaldo(item.Quantidade) {
			return fmt.Errorf("estoque insuficiente para peça '%s': disponível %d, solicitado %d",
				peca.Nome, peca.QuantidadeEstoque, item.Quantidade)
		}
	}
	return nil
}

// BaixarEstoque efetua a baixa física de todas as peças da OS.
func (g *GestorEstoque) BaixarEstoque(ctx context.Context, os *entity.OrdemServico) error {
	for _, item := range os.Pecas {
		peca, err := g.pecaRepo.BuscarPorID(ctx, item.PecaID)
		if err != nil {
			return fmt.Errorf("peça %s não encontrada: %w", item.PecaID, err)
		}
		if err := peca.BaixarEstoque(item.Quantidade); err != nil {
			return fmt.Errorf("falha ao baixar estoque da peça '%s': %w", peca.Nome, err)
		}
		if err := g.pecaRepo.Atualizar(ctx, peca); err != nil {
			return fmt.Errorf("falha ao persistir estoque da peça '%s': %w", peca.Nome, err)
		}
	}
	return nil
}
