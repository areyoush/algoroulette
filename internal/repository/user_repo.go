package repository

import (
	"database/sql"
	"log"
	"time"
	
	
	"github.com/areyoush/algoroulette/internal/model"
)


type UserRepository struct {
	db	*sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at",
		user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, email, password, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err

}

func (r *UserRepository) DenylistToken(token string, expires_at time.Time) error {
	query := `INSERT INTO denylist (token, expires_at) VALUES ($1, $2) ON CONFLICT (token) DO NOTHING` 	
	_, err := r.db.Exec(query, token, expires_at)
	
	return err
}

func (r *UserRepository) IsTokenDenylisted(token string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM denylist WHERE token = $1)"
	
	err := r.db.QueryRow(query, token).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking token denylist: %v", err)
		return false
	}
	return exists
}

func (r *UserRepository) CleanupDenylist() error {
	query := "DELETE FROM denylist WHERE expires_at < NOW()"
	
	_, err := r.db.Exec(query)
	if err != nil {
		log.Printf("Failed to cleanup token denylist: %v", err)
	}
	return err
	
}



