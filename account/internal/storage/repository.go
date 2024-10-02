package storage

import (
	"crypto/sha256"
	"encoding/base64"

	models "github.com/gnom48/hospital-api-lib"
)

type Repository struct {
	storage *Storage
}

func (r *Repository) AddUser(user *models.User) (*models.User, error) {
	encryptor := sha256.New()
	encryptor.Write([]byte(user.Password))
	hashed_password := encryptor.Sum(nil)
	hashed_password_base64 := base64.StdEncoding.EncodeToString(hashed_password)
	user.Password = hashed_password_base64
	user.Id, _ = models.GenerateUuid32()

	if err := r.storage.db.QueryRow(
		"INSERT INTO users (id, last_name, first_name, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Id, user.LastName, user.FirstName, user.Username, hashed_password_base64,
	).Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserByUsernamePassword(username string, password string) (*models.User, error) {
	encryptor := sha256.New()
	encryptor.Write([]byte(password))
	hashed_password := encryptor.Sum(nil)
	hashed_password_base64 := base64.StdEncoding.EncodeToString(hashed_password)

	user := models.User{}

	if err := r.storage.db.QueryRow(
		"SELECT * FROM users WHERE username = $1 AND password = $2",
		username, hashed_password_base64,
	).Scan(&user.Id, &user.LastName, &user.FirstName, &user.Username, &user.Password, &user.CreatedAt, &user.IsActive); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserById(userId string) (*models.User, error) {
	user := models.User{}

	if err := r.storage.db.QueryRow(
		"SELECT * FROM users WHERE id = $1",
		userId,
	).Scan(&user.Id, &user.LastName, &user.FirstName, &user.Username, &user.Password, &user.CreatedAt, &user.IsActive); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) addToken(token *models.Token) (string, error) {
	token.Id, _ = models.GenerateUuid32()
	if err := r.storage.db.QueryRow(
		"INSERT INTO tokens (id, token, user_id, is_regular) VALUES ($1, $2, $3, $4) RETURNING id",
		token.Id, token.Token, token.UserId, token.IsRegular,
	).Scan(&token.Id); err != nil {
		return "", err
	}

	return token.Id, nil
}

func (r *Repository) DeleteTokensPair(userId string) (bool, error) {
	if _, err := r.storage.db.Query(
		"DELETE FROM tokens WHERE user_id = $1;",
		userId,
	); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) SyncToken(tokenString string, userId string, isRegular bool) (string, error) {
	if _, err := r.storage.db.Query(
		"DELETE FROM tokens WHERE user_id = $1 AND is_regular = $2",
		userId, isRegular,
	); err != nil {
		return "", err
	}

	id, _ := models.GenerateUuid32()

	res, err := r.addToken(&models.Token{
		Id:        id,
		UserId:    userId,
		Token:     tokenString,
		IsRegular: isRegular,
	})
	if err != nil {
		return "", err
	}

	return res, nil
}
