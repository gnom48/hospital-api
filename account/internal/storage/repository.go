package storage

import models "github.com/gnom48/hospital-api-lib"

type Repository struct {
	storage *Storage
}

func (r *Repository) AddUser(u *models.User) (*models.User, error) {
	if err := r.storage.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		u.Email, u.Password,
	).Scan(&u.Id); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) GetUserByEmailPassword(email string, password string) (*models.User, error) {
	return nil, nil
}
