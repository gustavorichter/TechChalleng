package notification

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techchalleng/oficina/internal/application/port"
)

func TestEmailSimulado_RegistraMensagem(t *testing.T) {
	adapter := NewEmailSimulado()
	err := adapter.NotificarAtualizacaoStatus(context.Background(), port.NotificacaoStatus{
		ClienteID:      uuid.New(),
		Email:          "joao@email.com",
		OrdemServicoID: uuid.New(),
		StatusAnterior: "Recebida",
		StatusNovo:     "Em diagnóstico",
	})
	require.NoError(t, err)
	require.Len(t, adapter.Mensagens, 1)
	assert.Contains(t, adapter.Mensagens[0], "joao@email.com")
	assert.Contains(t, adapter.Mensagens[0], "Em diagnóstico")
}
