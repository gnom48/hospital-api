package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	models "github.com/gnom48/hospital-api-lib"
)

type timetableInfoRequestBody struct {
	HospitalId string    `json:"hospital_id"`
	DoctorId   string    `json:"doctor_id"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	Room       string    `json:"room"`
}

var timetableContextKey StringContextKey = "timetable"

func (s *ApiServer) TimetableInfoRequestBodyValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var timetableInfo timetableInfoRequestBody
		if err := json.NewDecoder(r.Body).Decode(&timetableInfo); err != nil {
			s.ErrorRespond(w, r, http.StatusUnprocessableEntity, fmt.Errorf("invalid request body: %w", err))
			return
		}
		s.logger.Debug(timetableInfo.From.String())

		if !IsValidTime(timetableInfo.From, timetableInfo.To) {
			s.ErrorRespond(w, r, http.StatusBadRequest, fmt.Errorf("invalid time range: from and to must be in increments of 30 minutes and to must be greater than from"))
			return
		}

		ctx := context.WithValue(r.Context(), timetableContextKey, timetableInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// @Summary Create a new timetable entry
// @Description Create a new entry in the timetable
// @Tags Timetable
// @Accept json
// @Produce json
// @Param timetable body timetableInfoRequestBody true "Timetable object"
// @Router /api/Timetable [post]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleCreateTimetable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}
		timetableInfo, ok := r.Context().Value(timetableContextKey).(timetableInfoRequestBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("No timetable info"))
			return
		}

		returning, err := s.storage.Repository().AddTimetable(models.Timetable{
			HospitalId: timetableInfo.HospitalId,
			DoctorId:   timetableInfo.DoctorId,
			Room:       timetableInfo.Room,
			TimeFrom:   timetableInfo.From,
			TimeTo:     timetableInfo.To,
		})
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusCreated, returning)
	}
}

// @Summary Update a timetable entry
// @Description Update an existing entry in the timetable
// @Tags Timetable
// @Accept json
// @Produce json
// @Param id path string true "Timetable ID"
// @Param timetable body timetableInfoRequestBody true "Timetable object"
// @Router /api/Timetable/{id} [put]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleUpdateTimetableById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}
		timetableInfo, ok := r.Context().Value(timetableContextKey).(timetableInfoRequestBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("No timetable info"))
			return
		}

		id := r.URL.Path[len("/api/Timetable/"):]

		err := s.storage.Repository().UpdateTimetable(models.Timetable{
			Id:         id,
			HospitalId: timetableInfo.HospitalId,
			DoctorId:   timetableInfo.DoctorId,
			TimeFrom:   timetableInfo.From,
			TimeTo:     timetableInfo.To,
			Room:       timetableInfo.Room,
		})
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusNoContent, nil)
	}
}

// @Summary Delete a timetable entry
// @Description Delete an existing entry in the timetable
// @Tags Timetable
// @Param id path string true "Timetable ID"
// @Router /api/Timetable/{id} [delete]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleDeleteTimetableById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		id := r.URL.Path[len("/api/Timetable/"):]

		err := s.storage.Repository().DeleteTimetable(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.ErrorRespond(w, r, http.StatusNotFound, fmt.Errorf("Nothing to delete"))
			} else {
				s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Respond(w, r, http.StatusNoContent, nil)
	}
}

// @Summary Delete all timetable entries for a doctor
// @Description Delete all entries in the timetable for a specific doctor
// @Tags Timetable
// @Param id path string true "Doctor ID"
// @Router /api/Timetable/Doctor/{id} [delete]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleDeleteTimetableByDoctorId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		doctorId := r.URL.Path[len("/api/Timetable/Doctor/"):]

		err := s.storage.Repository().DeleteTimetableByDoctorId(doctorId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.ErrorRespond(w, r, http.StatusNotFound, fmt.Errorf("Nothing to delete"))
			} else {
				s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Respond(w, r, http.StatusNoContent, nil)
	}
}

// @Summary Delete all timetable entries for a hospital
// @Description Delete all entries in the timetable for a specific hospital
// @Tags Timetable
// @Param id path string true "Hospital ID"
// @Router /api/Timetable/Hospital/{id} [delete]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleDeleteTimetableByHospitalId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		hospitalId := r.URL.Path[len("/api/Timetable/Hospital/"):]

		err := s.storage.Repository().DeleteTimetableByHospitalId(hospitalId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.ErrorRespond(w, r, http.StatusNotFound, fmt.Errorf("Nothing to delete"))
			} else {
				s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Respond(w, r, http.StatusNoContent, nil)
	}
}

// @Summary Get hospital timetable by ID
// @Description Retrieve the timetable for a specific hospital
// @Tags Timetable
// @Param id path string true "Hospital ID"
// @Param from query string false "From date (ISO8601)"
// @Param to query string false "To date (ISO8601)"
// @Router /api/Timetable/Hospital/{id} [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetTimetableByHospitalId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		hospitalId := r.URL.Path[len("/api/Timetable/Hospital/"):]

		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		from, err := time.Parse(fromStr, time.RFC3339)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}
		to, err := time.Parse(toStr, time.RFC3339)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		timetable, err := s.storage.Repository().GetTimetableByHospitalId(hospitalId, from, to)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, timetable)
	}
}

// @Summary Get doctor's timetable by ID
// @Description Retrieve the timetable for a specific doctor
// @Tags Timetable
// @Param id path string true "Doctor ID"
// @Param from query string false "From date (ISO8601)"
// @Param to query string false "To date (ISO8601)"
// @Router /api/Timetable/Doctor/{id} [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetTimetableByDoctorId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		doctorId := r.URL.Path[len("/api/Timetable/Doctor/"):]

		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		from, err := time.Parse(fromStr, time.RFC3339)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}
		to, err := time.Parse(toStr, time.RFC3339)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		timetable, err := s.storage.Repository().GetTimetableByDoctorId(doctorId, from, to)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, timetable)
	}
}

// @Summary Get hospital room timetable
// @Description Retrieve the timetable for a specific room in a hospital
// @Tags Timetable
// @Param id path string true "Hospital ID"
// @Param room path string true "Room Name"
// @Param from query string false "From date (ISO8601)"
// @Param to query string false "To date (ISO8601)"
// @Router /api/Timetable/Hospital/{id}/Room/{room} [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetTimetableByRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1") || IsUserInRole(userInfo.Roles, "2")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		// hospitalId := r.URL.Path[len("/api/Timetable/Hospital/") : len(r.URL.Path)-len("/Room/")-1]
		room := r.URL.Path[len(r.URL.Path)-len("/Room/"):]

		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		from, err := time.Parse(fromStr, time.RFC3339)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}
		to, err := time.Parse(toStr, time.RFC3339)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		timetable, err := s.storage.Repository().GetTimetableByHospitalRoom(room, from, to)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, timetable)
	}
}
