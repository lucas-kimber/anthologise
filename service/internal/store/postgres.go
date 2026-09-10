package store

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lucas-kimber/anthologise/service/internal/stremio"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, databaseURL string) (*PostgresStore, error) {

	dbConn, err := pgxpool.Connect(ctx, databaseURL)

	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: dbConn}, nil
}

func (s *PostgresStore) Health(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *PostgresStore) Migrate(ctx context.Context) error {

	q := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			edit_token_hash BYTEA NOT NULL
				CHECK (octet_length(edit_token_hash) = 32)
		);

		CREATE TABLE IF NOT EXISTS anthologies (
			id TEXT PRIMARY KEY,
			preview JSONB NOT NULL,
			contents JSONB NOT NULL
		);

		CREATE TABLE IF NOT EXISTS user_anthologies (
			user_id UUID NOT NULL REFERENCES users(id),
			anthology_id TEXT NOT NULL REFERENCES anthologies(id),

			PRIMARY KEY (user_id, anthology_id)
		);
	`

	_, err := s.db.Exec(ctx, q)
	return err
}

func (s *PostgresStore) Close() {
	s.db.Close()
}

func (s *PostgresStore) CreateAnthology(ctx context.Context, anthology stremio.Anthology) error {

	q := `
		INSERT INTO anthologies (id, preview, contents)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(ctx, q, anthology.ID, anthology.AnthologyPreview, anthology)
	return err
}

func (s *PostgresStore) CreateUser(ctx context.Context, userID string, editTokenHash [32]byte) error {

	q := `
		INSERT INTO users (id, edit_token_hash)
		VALUES ($1, $2)
	`

	_, err := s.db.Exec(ctx, q, userID, editTokenHash[:])
	return err
}

func (s *PostgresStore) GetAnthology(ctx context.Context, anthologyID string) (stremio.Anthology, error) {

	q := `
		SELECT contents 
		FROM anthologies 
		WHERE id = $1
	`

	var a stremio.Anthology
	if err := s.db.QueryRow(ctx, q, anthologyID).Scan(&a); err != nil {
		return stremio.Anthology{}, err
	}

	return a, nil
}

func (s *PostgresStore) GetCatalog(ctx context.Context, userID string) (stremio.Catalog, error) {

	q := `
		SELECT a.preview
		FROM anthologies AS a
		JOIN user_anthologies AS ua
			ON a.id =  ua.anthology_id
		WHERE ua.user_id = $1
		ORDER BY LOWER(a.preview->>'name')
	`

	rows, err := s.db.Query(ctx, q, userID)

	if err != nil {
		return stremio.Catalog{}, err
	}

	defer rows.Close()

	previews := make([]stremio.AnthologyPreview, 0)

	for rows.Next() {

		var p stremio.AnthologyPreview
		if err := rows.Scan(&p); err != nil {
			return stremio.Catalog{}, err
		}

		previews = append(previews, p)
	}

	if err := rows.Err(); err != nil {
		return stremio.Catalog{}, err
	}

	return stremio.Catalog{Metas: previews}, nil
}

func (s *PostgresStore) GetTokenHash(ctx context.Context, userID string) ([32]byte, error) {

	q := `
		SELECT edit_token_hash
		FROM users
		WHERE id = $1
	`

	var bytes []byte
	if err := s.db.QueryRow(ctx, q, userID).Scan(&bytes); err != nil {
		return [32]byte{}, err
	}

	var hash [32]byte
	copy(hash[:], bytes)

	return hash, nil
}

func (s *PostgresStore) SetCatalog(ctx context.Context, userID string, anthologyIDs []string) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var allExist bool

	err = tx.QueryRow(ctx, `
		SELECT NOT EXISTS (
			SELECT 1
			FROM unnest($1::text[]) AS id
			WHERE NOT EXISTS (
				SELECT 1
				FROM anthologies
				WHERE anthologies.id = id
			)
		)
	`, anthologyIDs).Scan(&allExist)

	if err != nil {
		return err
	}

	if !allExist {
		return ErrAnthologyNotFound
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM user_anthologies 
		WHERE user_id = $1
	`, userID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO user_anthologies (user_id, anthology_id)
		SELECT $1, id
		FROM unnest($2::text[]) AS id
	`, userID, anthologyIDs)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
