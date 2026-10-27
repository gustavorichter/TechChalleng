package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techchalleng/oficina/internal/application/port"
	"github.com/techchalleng/oficina/internal/domain/entity"
	"github.com/techchalleng/oficina/internal/domain/valobj"
)

const queryBuscarOrdemServico = `
	SELECT id, cliente_id, veiculo_id, status, valor_total, observacoes,
		orcamento_aprovado, motivo_recusa, criado_em, entregue_em, excluida_em, atualizado_em
	FROM ordens_servico
	WHERE id = $1`

const queryListarOrdensAtivas = `
	SELECT id, cliente_id, veiculo_id, status, valor_total, observacoes,
		orcamento_aprovado, motivo_recusa, criado_em, entregue_em, excluida_em, atualizado_em
	FROM ordens_servico
	WHERE excluida_em IS NULL
	  AND status NOT IN ('Finalizada', 'Entregue')
	ORDER BY
		CASE status
			WHEN 'Em execução' THEN 1
			WHEN 'Aguardando aprovação' THEN 2
			WHEN 'Em diagnóstico' THEN 3
			WHEN 'Recebida' THEN 4
			ELSE 5
		END,
		criado_em ASC`

var _ port.OrdemServicoRepository = (*OrdemServicoRepo)(nil)

type OrdemServicoRepo struct {
	pool *pgxpool.Pool
}

func NewOrdemServicoRepo(pool *pgxpool.Pool) port.OrdemServicoRepository {
	return &OrdemServicoRepo{pool: pool}
}

func (r *OrdemServicoRepo) Criar(ctx context.Context, os *entity.OrdemServico) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		INSERT INTO ordens_servico (
			id, cliente_id, veiculo_id, status, valor_total, observacoes,
			orcamento_aprovado, motivo_recusa, criado_em, entregue_em, excluida_em, atualizado_em
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		os.ID, os.ClienteID, os.VeiculoID, string(os.Status), os.ValorTotal, os.Observacoes,
		os.OrcamentoAprovado, motivoRecusaDB(os.MotivoRecusa), os.CriadoEm, os.EntregueEm, os.ExcluidaEm, os.AtualizadoEm)
	if err != nil {
		return err
	}
	if err := insertItensOS(ctx, tx, os); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrdemServicoRepo) BuscarPorID(ctx context.Context, id uuid.UUID) (*entity.OrdemServico, error) {
	row := r.pool.QueryRow(ctx, queryBuscarOrdemServico, id)

	os, err := scanOrdemServico(row)
	if err != nil {
		return nil, err
	}

	servicos, err := r.loadServicos(ctx, id)
	if err != nil {
		return nil, err
	}
	pecas, err := r.loadPecas(ctx, id)
	if err != nil {
		return nil, err
	}

	os.Servicos = servicos
	os.Pecas = pecas
	return os, nil
}

func (r *OrdemServicoRepo) Listar(ctx context.Context) ([]*entity.OrdemServico, error) {
	rows, err := r.pool.Query(ctx, queryListarOrdensAtivas)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.OrdemServico
	for rows.Next() {
		os, err := scanOrdemServicoFromRows(rows)
		if err != nil {
			return nil, err
		}
		servicos, err := r.loadServicos(ctx, os.ID)
		if err != nil {
			return nil, err
		}
		pecas, err := r.loadPecas(ctx, os.ID)
		if err != nil {
			return nil, err
		}
		os.Servicos = servicos
		os.Pecas = pecas
		result = append(result, os)
	}
	return result, rows.Err()
}

func (r *OrdemServicoRepo) Atualizar(ctx context.Context, os *entity.OrdemServico) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		UPDATE ordens_servico
		SET status=$2, valor_total=$3, observacoes=$4, orcamento_aprovado=$5, motivo_recusa=$6,
		    entregue_em=$7, excluida_em=$8, atualizado_em=$9
		WHERE id=$1`,
		os.ID, string(os.Status), os.ValorTotal, os.Observacoes, os.OrcamentoAprovado, motivoRecusaDB(os.MotivoRecusa),
		os.EntregueEm, os.ExcluidaEm, os.AtualizadoEm)
	if err != nil {
		return err
	}

	if err := replaceItensOS(ctx, tx, os); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrdemServicoRepo) TempoMedioExecucao(ctx context.Context) (float64, error) {
	var media *float64
	err := r.pool.QueryRow(ctx, `
		SELECT AVG(EXTRACT(EPOCH FROM (entregue_em - criado_em)) / 3600.0)
		FROM ordens_servico WHERE entregue_em IS NOT NULL`).Scan(&media)
	if err != nil {
		return 0, err
	}
	if media == nil {
		return 0, nil
	}
	return *media, nil
}

func (r *OrdemServicoRepo) loadServicos(ctx context.Context, osID uuid.UUID) ([]entity.ItemServicoOS, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT servico_id, nome, valor_mao_obra FROM os_servicos WHERE ordem_servico_id = $1`, osID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.ItemServicoOS
	for rows.Next() {
		var item entity.ItemServicoOS
		if err := rows.Scan(&item.ServicoID, &item.Nome, &item.ValorMaoObra); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrdemServicoRepo) loadPecas(ctx context.Context, osID uuid.UUID) ([]entity.ItemPecaOS, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT peca_id, nome, quantidade, valor_unitario FROM os_pecas WHERE ordem_servico_id = $1`, osID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.ItemPecaOS
	for rows.Next() {
		var item entity.ItemPecaOS
		if err := rows.Scan(&item.PecaID, &item.Nome, &item.Quantidade, &item.ValorUnitario); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanOrdemServico(row pgx.Row) (*entity.OrdemServico, error) {
	os, err := scanOrdemServicoValues(row.Scan)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("ordem de serviço não encontrada")
	}
	return os, err
}

func scanOrdemServicoFromRows(rows pgx.Rows) (*entity.OrdemServico, error) {
	return scanOrdemServicoValues(rows.Scan)
}

func scanOrdemServicoValues(scan func(dest ...any) error) (*entity.OrdemServico, error) {
	var os entity.OrdemServico
	var status string
	var motivo *string
	if err := scan(&os.ID, &os.ClienteID, &os.VeiculoID, &status, &os.ValorTotal,
		&os.Observacoes, &os.OrcamentoAprovado, &motivo, &os.CriadoEm, &os.EntregueEm, &os.ExcluidaEm, &os.AtualizadoEm); err != nil {
		return nil, err
	}
	os.Status = valobj.StatusOS(status)
	if motivo != nil {
		os.MotivoRecusa = *motivo
	}
	os.Servicos = []entity.ItemServicoOS{}
	os.Pecas = []entity.ItemPecaOS{}
	return &os, nil
}

func motivoRecusaDB(motivo string) *string {
	if motivo == "" {
		return nil
	}
	return &motivo
}

func insertItensOS(ctx context.Context, tx pgx.Tx, os *entity.OrdemServico) error {
	for _, s := range os.Servicos {
		if _, err := tx.Exec(ctx, `
			INSERT INTO os_servicos (ordem_servico_id, servico_id, nome, valor_mao_obra)
			VALUES ($1, $2, $3, $4)`,
			os.ID, s.ServicoID, s.Nome, s.ValorMaoObra); err != nil {
			return err
		}
	}
	for _, p := range os.Pecas {
		if _, err := tx.Exec(ctx, `
			INSERT INTO os_pecas (ordem_servico_id, peca_id, nome, quantidade, valor_unitario)
			VALUES ($1, $2, $3, $4, $5)`,
			os.ID, p.PecaID, p.Nome, p.Quantidade, p.ValorUnitario); err != nil {
			return err
		}
	}
	return nil
}

func replaceItensOS(ctx context.Context, tx pgx.Tx, os *entity.OrdemServico) error {
	if _, err := tx.Exec(ctx, `DELETE FROM os_servicos WHERE ordem_servico_id = $1`, os.ID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM os_pecas WHERE ordem_servico_id = $1`, os.ID); err != nil {
		return err
	}
	return insertItensOS(ctx, tx, os)
}
