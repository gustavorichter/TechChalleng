package notification

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/techchalleng/oficina/internal/application/port"
)

// EmailSimulado é o adaptador de notificação que registra o envio sem integrar um provedor real.
type EmailSimulado struct {
	mu        sync.Mutex
	Mensagens []string
}

func NewEmailSimulado() *EmailSimulado {
	return &EmailSimulado{}
}

func (e *EmailSimulado) NotificarAtualizacaoStatus(_ context.Context, n port.NotificacaoStatus) error {
	destino := n.Email
	if destino == "" {
		destino = n.ClienteID.String()
	}
	anterior := n.StatusAnterior
	if anterior == "" {
		anterior = "(abertura)"
	}
	mensagem := fmt.Sprintf("para %s: OS %s de %s para %s", destino, n.OrdemServicoID, anterior, n.StatusNovo)

	e.mu.Lock()
	e.Mensagens = append(e.Mensagens, mensagem)
	e.mu.Unlock()

	log.Printf("notificacao simulada de e-mail %s", mensagem)
	return nil
}

var _ port.NotificadorStatus = (*EmailSimulado)(nil)
