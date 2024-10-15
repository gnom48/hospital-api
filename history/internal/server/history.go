package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "github.com/gnom48/hospital-api-lib"
)

// @Summary Get visit history and appointments for an account
// @Description Retrieve records where {pacientId} = {id}
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Router /api/History/Account/{id} [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetAccountHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patientId := r.URL.Path[len("/api/History/Account/"):]

		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "2") || (IsUserInRole(userInfo.Roles, "3") && userInfo.User.Id == patientId)) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		history, err := s.storage.Repository().GetHistoryByPatientId(patientId)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, history)
	}
}

// @Summary Get detailed information about a visit and appointments
// @Description Retrieve visit and appointment details
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "History ID"
// @Router /api/History/{id} [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetHistoryDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		historyId := r.URL.Path[len("/api/History/"):]
		historyDetails, err := s.storage.Repository().GetHistoryDetailsById(historyId)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "2") || (IsUserInRole(userInfo.Roles, "3") && userInfo.User.Id == historyDetails.PatientId)) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		s.Respond(w, r, http.StatusOK, historyDetails)
	}
}

type createHistoryRequestBody struct {
	Date       time.Time `json:"date"`
	PatientId  string    `json:"patient_id"`
	HospitalId string    `json:"hospital_id"`
	DoctorId   string    `json:"doctor_id"`
	Room       string    `json:"room"`
	Data       string    `json:"data"`
}

func (chrb *createHistoryRequestBody) parseIntoVisitHistory() models.VisitHistory {
	return models.VisitHistory{
		Id:         "",
		VisitDate:  chrb.Date,
		PatientId:  chrb.PatientId,
		HospitalId: chrb.HospitalId,
		Room:       chrb.Room,
		Data:       chrb.Data,
		DoctorId:   chrb.DoctorId,
	}
}

// @Summary Create visit history and appointment
// @Description Create visit record
// @Tags History
// @Accept json
// @Produce json
// @Param history body createHistoryRequestBody true "History object"
// @Router /api/History [post]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleCreateHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1") || IsUserInRole(userInfo.Roles, "2")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}
		// TODO: проверка на то, что у patient роль = 3

		var requestBody createHistoryRequestBody
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.storage.Repository().CreateVisitHistory(requestBody.parseIntoVisitHistory())
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusCreated, nil)
	}
}

// @Summary Update visit history and appointment
// @Description Update visit record
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "History ID"
// @Param history body createHistoryRequestBody true "Updated History object"
// @Router /api/History/{id} [put]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleUpdateHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1") || IsUserInRole(userInfo.Roles, "2")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}
		// TODO: проверка на то, что у patient роль = 3

		historyID := r.URL.Path[len("/api/History/"):]

		var requestBody createHistoryRequestBody
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.storage.Repository().UpdateVisitHistory(historyID, requestBody.parseIntoVisitHistory())
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, nil)
	}
}
