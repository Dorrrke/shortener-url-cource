package db

import (
	"context"
	"url-shortener/internal/domain/models"
)

func (s *Storage) Save(link models.Link) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	//TODO: проверить что у пользователя уже не сохранена эта ссылка
	_, err := s.pool.Exec(ctx,
		"INSERT INTO links (short, original, user_id) VALUES ($1, $2, $3)",
		link.Short, link.Original, link.UserID,
	)

	return err
}

func (s *Storage) Get(short string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var link string

	err = tx.QueryRow(
		ctx,
		"SELECT original FROM links WHERE short = $1",
		short,
	).Scan(&link)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(ctx,
		"UPDATE links SET clicks = clicks + 1, updated_at = NOW() WHERE short = $1",
		short,
	)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return link, nil
}

func (s *Storage) Delete(short string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err := s.pool.Exec(
		ctx,
		"DELETE FROM links WHERE short = $1",
		short,
	)

	return err
}

func (s *Storage) GetAll(uid string) (map[string]models.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	links := make(map[string]models.Link)

	rows, err := s.pool.Query(
		ctx,
		"SELECT short, original, user_id FROM links WHERE user_id = $1",
		uid,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var link models.Link
		if err := rows.Scan(
			&link.Short,
			&link.Original,
			&link.UserID,
		); err != nil {
			return nil, err
		}

		links[link.Short] = link
	}

	return links, nil
}
