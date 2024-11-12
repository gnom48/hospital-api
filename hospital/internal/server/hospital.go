package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/gnom48/hospital-api-lib"
)

// @Summary Get list of hospitals
// @Description Retrieve a list of hospitals
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param from query int false "Pagination start"
// @Param count query int false "Number of records per page"
// @Router /api/Hospitals [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetHospitals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		fromParam := r.URL.Query().Get("from")
		countParam := r.URL.Query().Get("count")

		from, err := strconv.Atoi(fromParam)
		if err != nil {
			http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
			return
		}

		count, err := strconv.Atoi(countParam)
		if err != nil {
			http.Error(w, "Invalid 'count' parameter", http.StatusBadRequest)
			return
		}

		defer s.storage.Close()
		hospitals, err := s.storage.Repository().GetHospitals(from, count)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, hospitals)
	}
}

// @Summary Get hospital information by ID
// @Description Retrieve hospital information by hospital ID
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param id path string true "Hospital ID"
// @Router /api/Hospitals/{id} [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetHospitalById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		id := r.URL.Path[len("/api/Hospitals/"):]

		defer s.storage.Close()
		hospital, err := s.storage.Repository().GetHospitalById(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.ErrorRespond(w, r, http.StatusNotFound, fmt.Errorf("Hospital not found"))
			} else {
				s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Respond(w, r, http.StatusOK, hospital)
	}
}

// @Summary Get rooms by hospital ID
// @Description Retrieve a list of rooms in a hospital by hospital ID
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param id path string true "Hospital ID"
// @Router /api/Hospitals/{id}/Rooms [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetRoomsByHospitalId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		id := r.URL.Path[len("/api/Hospitals/") : len(r.URL.Path)-len("/Rooms")]

		defer s.storage.Close()
		rooms, err := s.storage.Repository().GetRoomsByHospitalId(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.ErrorRespond(w, r, http.StatusNotFound, fmt.Errorf("Rooms not found"))
			} else {
				s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Respond(w, r, http.StatusOK, rooms)
	}
}

type createUpdateHospitalRequestBody struct {
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	ContactPhone string   `json:"contact_phone"`
	Rooms        []string `json:"rooms"`
}

// @Summary Create a new hospital
// @Description Create a new hospital record
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param hospital body createUpdateHospitalRequestBody true "Hospital object"
// @Router /api/Hospitals [post]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleCreateHospital() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_info, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !IsUserInRole(user_info.Roles, "0") {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		var hospitalInfo createUpdateHospitalRequestBody
		if err := json.NewDecoder(r.Body).Decode(&hospitalInfo); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		defer s.storage.Close()
		returning, err := s.storage.Repository().AddHospital(hospitalInfo.Name, hospitalInfo.Address, hospitalInfo.ContactPhone, hospitalInfo.Rooms)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusCreated, returning)
	}
}

// @Summary Update hospital information by ID
// @Description Update hospital information
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param id path string true "Hospital ID"
// @Param hospital body createUpdateHospitalRequestBody true "Hospital object"
// @Router /api/Hospitals/{id} [put]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleUpdateHospital() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_info, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !IsUserInRole(user_info.Roles, "0") {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		id := r.URL.Path[len("/api/Hospitals/"):]

		var hospitalInfo createUpdateHospitalRequestBody
		if err := json.NewDecoder(r.Body).Decode(&hospitalInfo); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		editableHospital := models.Hospital{
			Id:           id,
			Name:         hospitalInfo.Name,
			Address:      hospitalInfo.Address,
			ContactPhone: hospitalInfo.ContactPhone,
		}
		defer s.storage.Close()
		err := s.storage.Repository().UpdateHospital(editableHospital, hospitalInfo.Rooms)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, hospitalInfo)
	}
}

// @Summary Soft delete a hospital record
// @Description Soft delete a hospital by ID
// @Tags Hospitals
// @Accept json
// @Produce json
// @Param id path string true "Hospital ID"
// @Router /api/Hospitals/{id} [delete]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleSoftDeleteHospital() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_info, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !IsUserInRole(user_info.Roles, "0") {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		id := r.URL.Path[len("/api/Hospitals/"):]

		defer s.storage.Close()
		err := s.storage.Repository().DeleteHospital(id)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusNoContent, nil)
	}
}
