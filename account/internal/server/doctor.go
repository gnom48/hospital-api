package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	models "github.com/gnom48/hospital-api-lib"
)

// @Summary Get doctor information by ID
// @Description Retrieve doctor information by doctor ID
// @Tags Doctors
// @Accept json
// @Produce json
// @Param id path string true "Doctor ID"
// @Router /api/Doctors/{id} [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetDoctorById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		id := r.URL.Path[len("/api/Doctors/"):]

		doctor, err := s.storage.Repository().GetDoctorById(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.ErrorRespond(w, r, http.StatusNotFound, fmt.Errorf("Doctor not found"))
			} else {
				s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Respond(w, r, http.StatusOK, doctor)
	}
}

// @Summary Get list of doctors
// @Description Retrieve a list of doctors with optional name filtering
// @Tags Doctors
// @Accept json
// @Produce json
// @Param nameFilter query string false "Filter by doctor's full name"
// @Param from query int false "Pagination start"
// @Param count query int false "Number of records per page"
// @Router /api/Doctors [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetDoctors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		nameFilter := r.URL.Query().Get("nameFilter")
		from := r.URL.Query().Get("from")
		count := r.URL.Query().Get("count")

		doctors, err := s.storage.Repository().GetDoctors(nameFilter, from, count)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, doctors)
	}
}
