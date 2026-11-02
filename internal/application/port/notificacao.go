package port

import (
	"context"

	"github.com/google/uuid"
)

// NotificacaoStatus descreve a mudança comunicada ao cliente.
type NotificacaoStatus struct {
	ClienteID      uuid.UUID
	Email          string
	OrdemServicoID uuid.UUID
	StatusAnterior string
	StatusNovo     string
}

// NotificadorStatus é a porta de saída para avisar o cliente sobre mudança de status.
type NotificadorStatus interface {
	NotificarAtualizacaoStatus(ctx context.Context, notificacao NotificacaoStatus) error
}
