package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash sql.NullString
	GoogleID     sql.NullString
}

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: sql.NullString{String: passwordHash, Valid: true},
		GoogleID:     sql.NullString{String: "", Valid: false},
	}

	row := us.DB.QueryRow(`
	  INSERT INTO users (email, password_hash)
	  VALUES ($1, $2) RETURNING id`, email, passwordHash)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`
		SELECT id, password_hash
		FROM users WHERE email=$1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return &user, nil
}

func (us *UserService) GetUser(email string) (*User, error) {
	email = strings.ToLower(email)
	row := us.DB.QueryRow(`
		SELECT id, google_id 
		FROM users 
		where email = $1;
	`, email)

	user := User{
		Email: email,
	}

	err := row.Scan(&user.ID, &user.GoogleID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) CreateGoogleUser(email, googleID string) (*User, error) {
	email = strings.ToLower(email)

	user := User{
		Email:    email,
		GoogleID: sql.NullString{String: googleID, Valid: true},
	}

	row := us.DB.QueryRow(`
	  INSERT INTO users (email, google_id)
	  VALUES ($1, $2) RETURNING id;`, email, googleID)

	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}
