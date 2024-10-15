package storage

import (
	"database/sql"
	"fmt"
	"time"

	models "github.com/gnom48/hospital-api-lib"
)

type Repository struct {
	storage *Storage
}

func (r *Repository) AddTimetable(timetable models.Timetable) (string, error) {
	timetable.Id, _ = models.GenerateUuid32()
	err := r.storage.db.QueryRow(
		"INSERT INTO timetable (id, hospital_id, doctor_id, time_from, time_to, room) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		timetable.Id, timetable.HospitalId, timetable.DoctorId, timetable.TimeFrom, timetable.TimeTo, timetable.Room,
	).Scan(&timetable.Id)
	if err != nil {
		return "", err
	}

	return timetable.Id, nil
}

func (r *Repository) UpdateTimetable(timetable models.Timetable) error {
	_, err := r.storage.db.Exec(
		"UPDATE timetable SET hospital_id = $1, doctor_id = $2, time_from = $3, time_to = $4, room = $5 WHERE id = $6",
		timetable.HospitalId, timetable.DoctorId, timetable.TimeFrom, timetable.TimeTo, timetable.Room, timetable.Id,
	)
	return err
}

func (r *Repository) DeleteTimetable(id string) error {
	result, err := r.storage.db.Exec("DELETE FROM timetable WHERE id = $1", id)
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

func (r *Repository) DeleteTimetableByDoctorId(doctorId string) error {
	result, err := r.storage.db.Exec("DELETE FROM timetable WHERE doctor_id = $1;", doctorId)
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

func (r *Repository) DeleteTimetableByHospitalId(hospitalId string) error {
	result, err := r.storage.db.Exec("DELETE FROM timetable WHERE hospital_id = $1;", hospitalId)
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

func (r *Repository) GetTimetableByHospitalId(hospitalId string, from time.Time, to time.Time) ([]models.Timetable, error) {
	var timetables []models.Timetable
	rows, err := r.storage.db.Query(
		"SELECT id, hospital_id, doctor_id, time_from, time_to, room FROM timetable WHERE hospital_id = $1 AND time_from >= $2 AND time_to <= $3",
		hospitalId, from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timetable models.Timetable
		if err := rows.Scan(&timetable.Id, &timetable.HospitalId, &timetable.DoctorId, &timetable.TimeFrom, &timetable.TimeTo, &timetable.Room); err != nil {
			return nil, err
		}
		timetables = append(timetables, timetable)
	}

	return timetables, nil
}

func (r *Repository) GetTimetableByDoctorId(doctorId string, from time.Time, to time.Time) ([]models.Timetable, error) {
	var timetables []models.Timetable
	rows, err := r.storage.db.Query(
		"SELECT id, hospital_id, doctor_id, time_from, time_to, room FROM timetable WHERE doctor_id = $1 AND time_from >= $2 AND time_to <= $3;",
		doctorId, from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timetable models.Timetable
		if err := rows.Scan(&timetable.Id, &timetable.HospitalId, &timetable.DoctorId, &timetable.TimeFrom, &timetable.TimeTo, &timetable.Room); err != nil {
			return nil, err
		}
		timetables = append(timetables, timetable)
	}

	return timetables, nil
}

func (r *Repository) GetTimetableById(id string) (models.Timetable, error) {
	var timetable models.Timetable
	if err := r.storage.db.QueryRow(
		"SELECT id, hospital_id, doctor_id, time_from, time_to, room FROM timetable WHERE id = $1",
		id,
	).Scan(&timetable); err != nil {
		return models.Timetable{}, err
	}
	return timetable, nil
}

func (r *Repository) GetTimetableByHospitalRoom(roomId string, from time.Time, to time.Time) ([]models.Timetable, error) {
	var timetables []models.Timetable
	rows, err := r.storage.db.Query(
		"SELECT id, hospital_id, doctor_id, time_from, time_to, room FROM timetable WHERE room = $1 AND time_from >= $2 AND time_to <= $3",
		roomId, from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timetable models.Timetable
		if err := rows.Scan(&timetable.Id, &timetable.HospitalId, &timetable.DoctorId, &timetable.TimeFrom, &timetable.TimeTo, &timetable.Room); err != nil {
			return nil, err
		}
		timetables = append(timetables, timetable)
	}

	return timetables, nil
}

func (r *Repository) GetAvailableAppointments(timetableId string) ([]time.Time, error) {
	var availableSlots []time.Time

	timetable, err := r.GetTimetableById(timetableId)
	if err != nil {
		return nil, err
	}

	for t := timetable.TimeFrom; t.Before(timetable.TimeTo); t = t.Add(30 * time.Minute) {
		availableSlots = append(availableSlots, t)
	}

	var bookedAppointments []models.Appointment
	rows, err := r.storage.db.Query(
		"SELECT id, timetable_id, user_id, appointment_time FROM appointments WHERE timetable_id = $1 LIMIT 1",
		timetable.Id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bookedAppointment models.Appointment
		if err := rows.Scan(&bookedAppointment.Id, &bookedAppointment.TimetableId, &bookedAppointment.UserId, &bookedAppointment.AppointmentTime); err != nil {
			return nil, err
		}
		bookedAppointments = append(bookedAppointments, bookedAppointment)
	}

	var filteredSlots []time.Time

	for _, availableSlot := range availableSlots {
		found := false
		for _, bookedAppointment := range bookedAppointments {
			if availableSlot == bookedAppointment.AppointmentTime {
				found = true
				break
			}
		}
		if !found {
			filteredSlots = append(filteredSlots, availableSlot)
		}
	}

	return filteredSlots, nil
}

func (r *Repository) BookAppointment(appointment models.Appointment) (string, error) {
	isValid := false

	timetable, err := r.GetTimetableById(appointment.TimetableId)
	if err != nil {
		return "", fmt.Errorf("Timetable not found")
	}

	availableAppointments, err := r.GetAvailableAppointments(timetable.Id)
	for _, appointmentTime := range availableAppointments {
		if appointmentTime == appointment.AppointmentTime {
			isValid = true
			break
		}
	}
	if !isValid {
		return "", fmt.Errorf("Invalid appointment time")
	}

	appointment.Id, _ = models.GenerateUuid32()
	_, err = r.storage.db.Exec(
		"INSERT INTO appointments (id, timetable_id, user_id, appointment_time) VALUES ($1, $2, $3, $4)",
		appointment.Id, appointment.TimetableId, appointment.UserId, appointment.AppointmentTime,
	)
	if err != nil {
		return "", err
	} else {
		return appointment.Id, nil
	}
}

func (r *Repository) CancelAppointment(appointmentId string) error {
	result, err := r.storage.db.Exec("DELETE FROM appointments WHERE id = $1;", appointmentId)
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

func (r *Repository) HasUserBookedAppointment(appointmentId string, userId string) (bool, error) {
	var appointment models.Appointment
	err := r.storage.db.QueryRow(
		"SELECT id, timetable_id, user_id, appointment_time FROM appointments WHERE id = $1 LIMIT 1",
		appointmentId,
	).Scan(&appointment.Id, &appointment.TimetableId, &appointment.UserId, &appointment.AppointmentTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if appointment.UserId != userId {
		return false, fmt.Errorf("Another user has booked this appointment")
	}

	return true, nil
}
