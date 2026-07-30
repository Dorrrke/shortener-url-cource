package db

import (
	"context"
	"url-shortener/internal/domain/models"
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
