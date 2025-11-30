package infrastructure

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type UserRepository struct {
	db Repository
}

func NewUserRepository(db Repository) *UserRepository {
	return &UserRepository{db: db}
}

// SignUp - Create a new user account
func (r *UserRepository) SignUp(email, password, username string) (*User, error) {
	// Check if user already exists
	var existingID int
	err := r.db.QueryRow("SELECT  ID FROM User WHERE email = ?", email).Scan(&existingID)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Insert user
	result, err := r.db.Exec("INSERT INTO User (email, password, username) VALUES (?, ?, ?)",
		email, string(hashedPassword), username)
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       int(userID),
		Email:    email,
		Username: username,
	}, nil
}

// Login - Authenticate user and return user data
func (r *UserRepository) Login(email, password string) (*User, error) {
	var user User

	err := r.db.QueryRow("SELECT  ID, email, password, Username FROM User WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Username)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid email or password")
	}
	if err != nil {
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	user.Password = "" // Clear password before returning
	return &user, nil
}

// GetUserByID - Retrieve user by ID
func (r *UserRepository) GetUserByID(id int) (*User, error) {
	var user User

	err := r.db.QueryRow("SELECT  ID, email, Username FROM User WHERE id = ?", id).
		Scan(&user.ID, &user.Email, &user.Username)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail - Retrieve user by email
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User

	err := r.db.QueryRow("SELECT  ID, email, Username FROM User WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Username)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
