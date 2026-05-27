package store

import (
	"database/sql"
	"time"

	"github.com/nitin-dixit/RepJot/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateTokens(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `

	insert into tokens (hash, user_id, expiry, scope)
	values ($1,$2,$3,$4)
	`
	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)

	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `
	delete from tokens
	where scope=$1 and user_id=$2
	`

	_, err := t.db.Exec(query, scope, userID)
	return err
}
