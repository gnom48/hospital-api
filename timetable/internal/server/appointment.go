package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "github.com/gnom48/hospital-api-lib"
)

// @Summary Get available appointments for a given timetable entry
// @Description Retrieve available appointment slots based on the timetable entry
// @Tags Timetable
// @Param id path string true "Timetable ID"
// @Router /api/Timetable/{id}/Appointments [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetAvailableAppointments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		timetableId := r.URL.Path[len("/api/Timetable/") : len(r.URL.Path)-len("/Appointments")]

		defer s.storage.Close()
		appointments, err := s.storage.Repository().GetAvailableAppointments(timetableId)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, appointments)
	}
}

type createAppointmentRequestBody struct {
	Time time.Time `json:"time"`
}

// @Summary Book an appointment
// @Description Make an appointment for a specific slot
// @Tags Timetable
// @Param id path string true "Timetable ID"
// @Param appointment body createAppointmentRequestBody true "Appointment request"
// @Router /api/Timetable/{id}/Appointments [post]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleBookAppointment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		timetableId := r.URL.Path[len("/api/Timetable/") : len(r.URL.Path)-len("/Appointments")]

		var appointmentInfo createAppointmentRequestBody
		if err := json.NewDecoder(r.Body).Decode(&appointmentInfo); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		defer s.storage.Close()
		appointmentId, err := s.storage.Repository().BookAppointment(models.Appointment{
			UserId:          userInfo.User.Id,
			TimetableId:     timetableId,
			AppointmentTime: appointmentInfo.Time,
		})
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusCreated, appointmentId)
	}
}

// @Summary Cancel an appointment
// @Description Cancel a previously made appointment
// @Tags Appointment
// @Param id path string true "Appointment ID"
// @Router /api/Appointment/{id} [delete]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleCancelAppointment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(InfoAboutMeResponseBody)
		if !ok || !(IsUserInRole(userInfo.Roles, "0") || IsUserInRole(userInfo.Roles, "1") || IsUserInRole(userInfo.Roles, "3")) {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Access denied"))
			return
		}

		appointmentId := r.URL.Path[len("/api/Appointment/"):]
		defer s.storage.Close()
		if res, err := s.storage.Repository().HasUserBookedAppointment(appointmentId, userInfo.User.Id); err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		} else {
			if !res {
				s.ErrorRespond(w, r, http.StatusBadRequest, fmt.Errorf("Not found appointment"))
				return
			}
		}

		err := s.storage.Repository().CancelAppointment(appointmentId)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusNoContent, nil)
	}
}
