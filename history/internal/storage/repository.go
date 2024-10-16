package storage

import (
	models "github.com/gnom48/hospital-api-lib"
)

type Repository struct {
	storage *Storage
}

func (r *Repository) GetHistoryByPatientId(patientId string) ([]models.VisitHistory, error) {
	rows, err := r.storage.db.Query(`
		SELECT id, patient_id, hospital_id, doctor_id, room, visit_date, data 
		FROM visit_history 
		WHERE patient_id = $1`, patientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.VisitHistory
	for rows.Next() {
		var record models.VisitHistory
		if err := rows.Scan(&record.Id, &record.PatientId, &record.HospitalId, &record.DoctorId, &record.Room, &record.VisitDate, &record.Data); err != nil {
			return nil, err
		}
		history = append(history, record)
	}

	return history, nil
}

func (r *Repository) GetHistoryDetailsById(historyId string) (models.VisitHistory, error) {
	var history models.VisitHistory

	err := r.storage.db.QueryRow(`
		SELECT id, patient_id, hospital_id, doctor_id, room, visit_date, data 
		FROM visit_history 
		WHERE id = $1`, historyId).Scan(
		&history.Id, &history.PatientId, &history.HospitalId, &history.DoctorId, &history.Room, &history.VisitDate, &history.Data)

	if err != nil {
		return models.VisitHistory{}, err
	}

	return history, nil
}

func (r *Repository) CreateVisitHistory(history models.VisitHistory) error {
	history.Id, _ = models.GenerateUuid32()
	_, err := r.storage.db.Exec(`
		INSERT INTO visit_history (id, patient_id, hospital_id, doctor_id, room, visit_date, data) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		history.Id, history.PatientId, history.HospitalId, history.DoctorId, history.Room, history.VisitDate, history.Data)
	return err
}

func (r *Repository) UpdateVisitHistory(historyId string, updatedHistory models.VisitHistory) error {
	_, err := r.storage.db.Exec(`
		UPDATE visit_history SET 
			patient_id = $1,
			hospital_id = $2,
			doctor_id = $3,
			room = $4,
			visit_date = $5,
			data = $6 
		WHERE id = $7`,
		updatedHistory.PatientId, updatedHistory.HospitalId, updatedHistory.DoctorId, updatedHistory.Room, updatedHistory.VisitDate, updatedHistory.Data, historyId)
	return err
}
