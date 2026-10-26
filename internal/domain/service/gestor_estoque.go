package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/techchalleng/oficina/internal/domain/entity"
)

// GestorEstoque aplica as regras de saldo sobre as peças já carregadas pelo caso de uso.
// A leitura e a gravação ficam na aplicação, para o domínio não depender de porta de persistência.
type GestorEstoque struct{}

func NewGestorEstoque() *GestorEstoque {
	return &GestorEstoque{}
}

// VerificarDisponibilidade valida se há saldo para todas as peças da OS.
func (g *GestorEstoque) VerificarDisponibilidade(os *entity.OrdemServico, pecas map[uuid.UUID]*entity.Peca) error {
	for _, item := range os.Pecas {
		peca, ok := pecas[item.PecaID]
		if !ok || peca == nil {
			return fmt.Errorf("peça %s não encontrada", item.PecaID)
		}
		if !peca.TemSaldo(item.Quantidade) {
			return fmt.Errorf("estoque insuficiente para peça '%s': disponível %d, solicitado %d",
				peca.Nome, peca.QuantidadeEstoque, item.Quantidade)
		}
	}
	return nil
}

// BaixarEstoque efetua a baixa física em memória. O caso de uso persiste as peças alteradas.
func (g *GestorEstoque) BaixarEstoque(os *entity.OrdemServico, pecas map[uuid.UUID]*entity.Peca) error {
	for _, item := range os.Pecas {
		peca, ok := pecas[item.PecaID]
		if !ok || peca == nil {
			return fmt.Errorf("peça %s não encontrada", item.PecaID)
		}
		if err := peca.BaixarEstoque(item.Quantidade); err != nil {
			return fmt.Errorf("falha ao baixar estoque da peça '%s': %w", peca.Nome, err)
		}
	}
	return nil
}
