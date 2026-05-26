package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err

	}
	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"_"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

func (s *PostgresUserStore) CreateUser(u *User) error {
	query := `
	insert into users (username,email,password_hash,bio)
	values ($1,$2,$3,$4)
	returning id,created_at, updated_at
	`
	err := s.db.QueryRow(query, u.Username, u.Email, u.PasswordHash.hash, u.Bio).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) GetUserByUsername(uname string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `
	select id, username, email, password_hash, bio, created_at, updated_at
	from users
	where username=$1
	`
	err := s.db.QueryRow(query, uname).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(u *User) error {
	query := `
	update users
	set username = $1, email=$2, bio=$3, updated_at=CURRENT_TIMESTAMP
	where id=$4
	returning updated_at
	`

	result, err := s.db.Exec(query, u.Username, u.Email, u.Bio, u.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *PostgresUserStore) GetUserToken(scope, plainTextPassword string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plainTextPassword))

	query := `
	select u.id, u.username, u.email, u.password_hash, u.bio, u.created_at, u.updated_at
	from users u
	inner join tokens t on t.user_id=u.id
	where t.hash=$1 and t.scope=$2 and t.expiry > $3
	`

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
