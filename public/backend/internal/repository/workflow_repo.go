package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelforge/backend/internal/model"
)

type WorkflowRepository struct {
	db *pgxpool.Pool
}

func NewWorkflowRepository(db *pgxpool.Pool) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) List(ctx context.Context) ([]model.Workflow, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, description, type, category, icon, is_active, version, created_at, updated_at 
		 FROM workflows ORDER BY type DESC, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []model.Workflow
	for rows.Next() {
		var w model.Workflow
		if err := rows.Scan(&w.ID, &w.UserID, &w.Name, &w.Description, &w.Type,
			&w.Category, &w.Icon, &w.IsActive, &w.Version, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workflows = append(workflows, w)
	}
	return workflows, nil
}

func (r *WorkflowRepository) GetByID(ctx context.Context, id int) (*model.Workflow, error) {
	w := &model.Workflow{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, description, type, category, icon, is_active, version, created_at, updated_at 
		 FROM workflows WHERE id = $1`, id,
	).Scan(&w.ID, &w.UserID, &w.Name, &w.Description, &w.Type,
		&w.Category, &w.Icon, &w.IsActive, &w.Version, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// 加载节点
	nodeRows, err := r.db.Query(ctx,
		`SELECT id, workflow_id, node_key, node_type, agent_id, name, config, position_x, position_y, sort_order 
		 FROM workflow_nodes WHERE workflow_id = $1 ORDER BY sort_order`, id)
	if err != nil {
		return nil, err
	}
	defer nodeRows.Close()

	for nodeRows.Next() {
		var n model.WorkflowNode
		if err := nodeRows.Scan(&n.ID, &n.WorkflowID, &n.NodeKey, &n.NodeType, &n.AgentID,
			&n.Name, &n.Config, &n.PositionX, &n.PositionY, &n.SortOrder); err != nil {
			return nil, err
		}
		w.Nodes = append(w.Nodes, n)
	}

	// 加载边
	edgeRows, err := r.db.Query(ctx,
		`SELECT id, workflow_id, from_node_id, to_node_id, edge_type, condition_expr, label, sort_order 
		 FROM workflow_edges WHERE workflow_id = $1 ORDER BY sort_order`, id)
	if err != nil {
		return nil, err
	}
	defer edgeRows.Close()

	for edgeRows.Next() {
		var e model.WorkflowEdge
		if err := edgeRows.Scan(&e.ID, &e.WorkflowID, &e.FromNodeID, &e.ToNodeID,
			&e.EdgeType, &e.ConditionExpr, &e.Label, &e.SortOrder); err != nil {
			return nil, err
		}
		w.Edges = append(w.Edges, e)
	}

	return w, nil
}

func (r *WorkflowRepository) Create(ctx context.Context, w *model.Workflow) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO workflows (user_id, name, description, type, category, icon) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		w.UserID, w.Name, w.Description, w.Type, w.Category, w.Icon,
	).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (r *WorkflowRepository) Update(ctx context.Context, w *model.Workflow) error {
	_, err := r.db.Exec(ctx,
		`UPDATE workflows SET name=$1, description=$2, category=$3, icon=$4, 
			is_active=$5, version=version+1, updated_at=NOW() WHERE id=$6`,
		w.Name, w.Description, w.Category, w.Icon, w.IsActive, w.ID)
	return err
}

func (r *WorkflowRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM workflows WHERE id = $1 AND type = 'custom'", id)
	return err
}

func (r *WorkflowRepository) CreateExecution(ctx context.Context, exec *model.WorkflowExecution) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO workflow_executions (workflow_id, project_id, user_id, status, input_data)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, started_at`,
		exec.WorkflowID, exec.ProjectID, exec.UserID, exec.Status, exec.InputData,
	).Scan(&exec.ID, &exec.StartedAt)
}

func (r *WorkflowRepository) GetExecution(ctx context.Context, id int) (*model.WorkflowExecution, error) {
	exec := &model.WorkflowExecution{}
	err := r.db.QueryRow(ctx,
		`SELECT id, workflow_id, project_id, user_id, status, input_data, output_data, 
				current_node_id, error_message, started_at, completed_at
		 FROM workflow_executions WHERE id = $1`, id,
	).Scan(&exec.ID, &exec.WorkflowID, &exec.ProjectID, &exec.UserID, &exec.Status,
		&exec.InputData, &exec.OutputData, &exec.CurrentNodeID, &exec.ErrorMessage,
		&exec.StartedAt, &exec.CompletedAt)
	return exec, err
}
