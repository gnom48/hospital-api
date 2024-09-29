package storage

import models "github.com/gnom48/hospital-api-lib"

type Repository struct {
	storage *Storage
}

func (r *Repository) AddUser(u *models.User) (*models.User, error) {
	return nil, nil
}

func (r *Repository) GetUserByEmailPassword(email string, password string) (*models.User, error) {
	return nil, nil
}
