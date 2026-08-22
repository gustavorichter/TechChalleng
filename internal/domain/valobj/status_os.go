package valobj

import (
	"errors"
	"fmt"
)

// StatusOS representa o status de uma Ordem de Serviço com ciclo de vida controlado.
type StatusOS string

const (
	StatusRecebida            StatusOS = "Recebida"
	StatusEmDiagnostico       StatusOS = "Em diagnóstico"
	StatusAguardandoAprovacao StatusOS = "Aguardando aprovação"
	StatusEmExecucao          StatusOS = "Em execução"
	StatusFinalizada          StatusOS = "Finalizada"
	StatusEntregue            StatusOS = "Entregue"
)

var transicoesValidas = map[StatusOS][]StatusOS{
	StatusRecebida:            {StatusEmDiagnostico},
	StatusEmDiagnostico:       {StatusAguardandoAprovacao},
	StatusAguardandoAprovacao: {StatusEmExecucao},
	StatusEmExecucao:          {StatusFinalizada},
	StatusFinalizada:          {StatusEntregue},
	StatusEntregue:            {},
}

// NewStatusOS cria um StatusOS a partir de uma string, validando se é um status conhecido.
func NewStatusOS(s string) (StatusOS, error) {
	status := StatusOS(s)
	switch status {
	case StatusRecebida, StatusEmDiagnostico, StatusAguardandoAprovacao,
		StatusEmExecucao, StatusFinalizada, StatusEntregue:
		return status, nil
	default:
		return "", fmt.Errorf("status inválido: %s", s)
	}
}

// TransicionarPara valida e retorna o próximo status se a transição for permitida.
func (s StatusOS) TransicionarPara(destino StatusOS) (StatusOS, error) {
	permitidos, ok := transicoesValidas[s]
	if !ok {
		return s, errors.New("status atual desconhecido")
	}
	for _, p := range permitidos {
		if p == destino {
			return destino, nil
		}
	}
	return s, fmt.Errorf("transição inválida de '%s' para '%s'", s, destino)
}

// ProximoStatus retorna o próximo status na sequência automática, se existir.
func (s StatusOS) ProximoStatus() (StatusOS, bool) {
	permitidos := transicoesValidas[s]
	if len(permitidos) == 0 {
		return s, false
	}
	return permitidos[0], true
}

// PodeTransicionarPara verifica se a transição é válida sem efetivá-la.
func (s StatusOS) PodeTransicionarPara(destino StatusOS) bool {
	_, err := s.TransicionarPara(destino)
	return err == nil
}

// RequerBaixaEstoque indica se o status implica baixa de estoque.
func (s StatusOS) RequerBaixaEstoque() bool {
	return s == StatusFinalizada
}
