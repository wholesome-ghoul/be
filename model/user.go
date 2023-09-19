package model

import (
	"github/be/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint
	Username string `json:"username"`
	Password string `json:"-"`
	Entries  []Entry
}

func (user *User) Save() (*User, error) {
	if err := user.Prepare(); err != nil {
		return &User{}, err
	}

	lastInsertId := 0
	err := database.Database.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, user.Password).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = uint(lastInsertId)

	return user, nil
}

func (user *User) Prepare() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Password = string(passwordHash)

	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := database.Database.
		QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	var user User
	err := database.Database.
		QueryRow("SELECT id, username, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
