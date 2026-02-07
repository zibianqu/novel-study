package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelforge/backend/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO users (username, email, password_hash, avatar, settings) 
		 VALUES ($1, $2, $3, $4, '{}') 
		 RETURNING id, created_at, updated_at`,
		user.Username, user.Email, user.PasswordHash, user.Avatar,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password_hash, avatar, settings, created_at, updated_at 
		 FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Avatar, &user.Settings, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password_hash, avatar, settings, created_at, updated_at 
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Avatar, &user.Settings, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}
	return user, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	return count > 0, err
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	return count > 0, err
}

func (r *UserRepository) UpdateAPIKey(ctx context.Context, userID int, encryptedKey string) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET api_key_encrypted = $1, updated_at = NOW() WHERE id = $2",
		encryptedKey, userID)
	return err
}
