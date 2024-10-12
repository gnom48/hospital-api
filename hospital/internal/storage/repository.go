package storage

import (
	"errors"
	"fmt"
	"strings"

	models "github.com/gnom48/hospital-api-lib"
)

type Repository struct {
	storage *Storage
}

func (r *Repository) GetHospitals(from int, count int) ([]models.Hospital, error) {
	rows, err := r.storage.db.Query(
		"SELECT id, name, address, contact_phone, created_at, is_active FROM hospitals ORDER BY created_at LIMIT $1 OFFSET $2;",
		count, from,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hospitals []models.Hospital
	for rows.Next() {
		var hospital models.Hospital
		if err := rows.Scan(&hospital.Id, &hospital.Name, &hospital.Address, &hospital.ContactPhone, &hospital.CreatedAt, &hospital.IsActive); err != nil {
			return nil, err
		}
		hospitals = append(hospitals, hospital)
	}

	return hospitals, nil
}

func (r *Repository) GetHospitalById(id string) (models.Hospital, error) {
	var hospital models.Hospital

	err := r.storage.db.QueryRow(
		"SELECT id, name, address, contact_phone, created_at, is_active FROM hospitals WHERE id = $1;",
		id,
	).Scan(&hospital.Id, &hospital.Name, &hospital.Address, &hospital.ContactPhone, &hospital.CreatedAt, &hospital.IsActive)

	if err != nil {
		return models.Hospital{}, err
	}

	return hospital, nil
}

func (r *Repository) GetRoomsByHospitalId(id string) ([]models.Room, error) {
	rows, err := r.storage.db.Query(
		"SELECT id, name, hospital_id FROM rooms WHERE hospital_id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	rooms := make([]models.Room, 0)
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.Id, &room.HospitalId, &room.HospitalId); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *Repository) AddHospital(name string, address string, contactNAme string, rooms []string) (string, error) {
	hospital := models.Hospital{
		Name:         name,
		Address:      address,
		ContactPhone: contactNAme,
	}
	hospital.Id, _ = models.GenerateUuid32()
	_, err := r.storage.db.Exec(
		"INSERT INTO hospitals (id, name, address, contact_phone) VALUES ($1, $2, $3, $4);",
		hospital.Id, hospital.Name, hospital.Address, hospital.ContactPhone,
	)

	var roomsErrors []string
	for _, roomName := range rooms {
		room := models.Room{
			HospitalId: hospital.Id,
			Name:       roomName,
		}
		room.Id, _ = models.GenerateUuid32()
		if _, e := r.addRoom(room); e != nil {
			roomsErrors = append(roomsErrors, e.Error())
		}
	}
	if len(roomsErrors) == 0 {
		err = nil
	} else {
		err = errors.New(strings.Join(roomsErrors, "; "))
	}

	return hospital.Id, err
}

func (r *Repository) addRoom(room models.Room) (string, error) {
	room.Id, _ = models.GenerateUuid32()
	_, err := r.storage.db.Exec(
		"INSERT INTO rooms (id, name, hospital_id) VALUES ($1, $2, $3);",
		room.Id, room.Name, room.HospitalId,
	)
	return room.Id, err
}

func (r *Repository) UpdateHospital(hospital models.Hospital, rooms []string) error {
	_, err := r.storage.db.Exec(
		"UPDATE hospitals SET name = $1, address = $2, contact_phone = $3 WHERE id = $4;",
		hospital.Name, hospital.Address, hospital.ContactPhone, hospital.Id,
	)

	var roomsErrors []string
	for _, roomName := range rooms {
		room := models.Room{
			HospitalId: hospital.Id,
			Name:       roomName,
		}
		room.Id, _ = models.GenerateUuid32()
		if _, e := r.addRoom(room); e != nil {
			roomsErrors = append(roomsErrors, e.Error())
		}
	}
	if len(roomsErrors) == 0 {
		err = nil
	} else {
		err = errors.New(strings.Join(roomsErrors, "; "))
	}

	return err
}

func (r *Repository) DeleteHospital(id string) error {
	isActive := true
	err := r.storage.db.QueryRow(
		"SELECT is_active FROM hospitals WHERE id = $1",
		id,
	).Scan(&isActive)
	if err != nil {
		return err
	}
	if !isActive {
		return fmt.Errorf("Hospital already had been deleted")
	}
	_, err = r.storage.db.Exec(
		"UPDATE hospitals SET is_active = false WHERE id = $1",
		id,
	)
	return err
}
