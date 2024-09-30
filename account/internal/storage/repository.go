package storage

import (
	"crypto/sha256"

	models "github.com/gnom48/hospital-api-lib"
)

type Repository struct {
	storage *Storage
}

func (r *Repository) AddUser(user *models.User) (*models.User, error) {
	encryptor := sha256.New()
	encryptor.Write([]byte(user.Password))
	hashed_password := encryptor.Sum(nil)
	if err := r.storage.db.QueryRow(
		"INSERT INTO users (last_name, first_name, username, password) VALUES ($1, $2, $3, $4) RETURNING id",
		user.LastName, user.FirstName, user.Username, hashed_password,
	).Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserByUsernamePassword(username string, password string) (*models.User, error) {
	encryptor := sha256.New()
	encryptor.Write([]byte(password))
	hashed_password := encryptor.Sum(nil)

	user := models.User{}

	if err := r.storage.db.QueryRow(
		"SELECT * FROM users WHERE username = $1 AND password = $2",
		username, hashed_password,
	).Scan(&user.Id, &user.LastName, &user.FirstName, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserById(userId int) (*models.User, error) {
	user := models.User{}

	if err := r.storage.db.QueryRow(
		"SELECT * FROM users WHERE id = $1",
		userId,
	).Scan(&user.Id, &user.LastName, &user.FirstName, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
