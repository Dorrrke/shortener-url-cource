package db

import (
	"context"

	"github.com/Dorrrke/shortener-url-cource/internal/domain/models"
)

func (s *Storage) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err := s.pool.Exec(ctx,
		"INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)",
		user.ID, user.Name, user.Email, user.Password,
	)

	return err
}

func (s *Storage) GetUser(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var user models.User

	err := s.pool.QueryRow(
		ctx,
		"SELECT id, name, email, password FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *Storage) GetUserInfo(userID string) (models.User, []models.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.User{}, nil, err
	}
	defer tx.Rollback(ctx)

	var user models.User
	err = tx.QueryRow(
		ctx,
		"SELECT id, name, email, created_at FROM users WHERE id = $1",
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return models.User{}, nil, err
	}

	var links []models.Link
	rows, err := tx.Query(ctx, "SELECT short, original, clicks, updated_at, created_at FROM links WHERE user_id = $1", userID)
	if err != nil {
		return models.User{}, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var link models.Link
		if err := rows.Scan(&link.Short, &link.Original, &link.Clicks, &link.UpdatedAt, &link.CratedAt); err != nil {
			return models.User{}, nil, err
		}
		link.UserID = userID
		links = append(links, link)
	}

	if err = tx.Commit(ctx); err != nil {
		return models.User{}, nil, err
	}

	return user, links, nil
}
