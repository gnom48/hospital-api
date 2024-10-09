package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/gnom48/hospital-api-lib"
)

// @Summary Get current account
// @Description Retrieve the current account's data
// @Tags Accounts
// @Accept json
// @Produce json
// @Router /api/Accounts/Me [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetCurrentAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		s.Respond(w, r, http.StatusOK, user)
	}
}

type updateAccountRequestBody struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Password  string `json:"password"`
}

// @Summary Update account
// @Description Update the current account's information
// @Tags Accounts
// @Accept json
// @Produce json
// @Param requestBody body updateAccountRequestBody true "Account Update Data"
// @Router /api/Accounts/Update [put]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleUpdateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}

		requestBody := &updateAccountRequestBody{}
		if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		user.LastName = requestBody.LastName
		user.FirstName = requestBody.FirstName
		user.Password = requestBody.Password

		if err := s.storage.Repository().UpdateUser(&user); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}
		s.Respond(w, r, http.StatusOK, nil)
	}
}

type getAllAccountsResponseBody struct {
	Accounts []models.User `json:"accounts"`
}

// @Summary Get all accounts
// @Description Retrieve a list of all accounts
// @Tags Accounts
// @Accept json
// @Produce json
// @Param from query int false "Start index"
// @Param count query int false "Number of records"
// @Router /api/Accounts [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleGetAllAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}
		roles, ok := r.Context().Value(RoleContextKey).([]models.Role)
		if !ok {
			s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access forbidden"))
			return
		} else {
			if IsUserInRole(roles, "0") {
				s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access only for admin"))
				return
			}
		}

		from := r.URL.Query().Get("from")
		count := r.URL.Query().Get("count")

		accounts, err := s.storage.Repository().GetAllAccounts(from, count)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Respond(w, r, http.StatusOK, getAllAccountsResponseBody{Accounts: accounts})
	}
}

type createAccountRequestBody struct {
	LastName  string   `json:"last_name"`
	FirstName string   `json:"first_name"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
}

// @Summary Create a new account
// @Description Create a new user account by admin
// @Tags Accounts
// @Accept json
// @Produce json
// @Param requestBody body createAccountRequestBody true "Account Creation Data"
// @Router /api/Accounts [post]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleCreateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}
		roles, ok := r.Context().Value(RoleContextKey).([]models.Role)
		if !ok {
			s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access forbidden"))
			return
		} else {
			if IsUserInRole(roles, "0") {
				s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access only for admin"))
				return
			}
		}

		requestBody := &createAccountRequestBody{}
		if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		newUser := &models.User{
			LastName:  requestBody.LastName,
			FirstName: requestBody.FirstName,
			Username:  requestBody.Username,
			Password:  requestBody.Password,
		}

		if returning, err := s.storage.Repository().AddUser(newUser); err != nil {
			s.ErrorRespond(w, r, http.StatusUnprocessableEntity, err)
			return
		} else {
			currentErrors := make([]string, 0)
			for _, r := range requestBody.Roles {
				if e := s.storage.Repository().AddUserRole(r); e != nil {
					currentErrors = append(currentErrors, fmt.Errorf(e).Error())
				}
			}
			s.Respond(w, r, http.StatusCreated, struct {
				NewUserId string   `json:"new_user_id"`
				Errors    []string `json:"errors"`
			}{
				NewUserId: returning.Id,
				Errors:    currentErrors,
			})
		}
	}
}

// @Summary Update account by ID
// @Description Update a user account by ID
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param requestBody body createAccountRequestBody true "Account Details"
// @Router /api/Accounts/{id} [put]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleUpdateAccountById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}
		roles, ok := r.Context().Value(RoleContextKey).([]models.Role)
		if !ok {
			s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access forbidden"))
			return
		} else {
			if IsUserInRole(roles, "0") {
				s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access only for admin"))
				return
			}
		}

		id := r.URL.Path[len("/api/Accounts/"):]
		requestBody := &createAccountRequestBody{}
		if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		editableUser := models.User{
			Id:        id,
			LastName:  requestBody.LastName,
			FirstName: requestBody.FirstName,
			Username:  requestBody.Username,
			Password:  requestBody.Password,
		}

		if err := s.storage.Repository().UpdateUser(&editableUser); err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
		} else {
			currentErrors := make([]string, 0)
			e := s.storage.Repository().DeleteAllUserRoles(editableUser.Id)
			if e != nil {
				currentErrors = append(currentErrors, e)
			}
			for _, r := range requestBody.Roles {
				if e := s.storage.Repository().AddUserRole(r); e != nil {
					currentErrors = append(currentErrors, fmt.Errorf(e).Error())
				}
			}
			s.Respond(w, r, http.StatusCreated, struct {
				UserId string   `json:"user_id"`
				Errors []string `json:"errors"`
			}{
				UserId: editableUser.Id,
				Errors: currentErrors,
			})
		}
	}
}

// @Summary Soft delete an account by ID
// @Description Soft delete a user account by ID
// @Tags Accounts
// @Router /api/Accounts/{id} [delete]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleSoftDeleteAccountById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}
		roles, ok := r.Context().Value(RoleContextKey).([]models.Role)
		if !ok {
			s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access forbidden"))
			return
		} else {
			if IsUserInRole(roles, "0") {
				s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Access only for admin"))
				return
			}
		}

		id := r.URL.Path[len("/api/Accounts/"):]
		if err := s.storage.Repository().SoftDeleteUser(id); err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, err)
			return
		}
		s.Respond(w, r, http.StatusOK, nil)
	}
}
