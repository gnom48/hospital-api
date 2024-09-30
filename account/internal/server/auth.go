package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/gnom48/hospital-api-lib"
)

func (s *ApiServer) HandleAuthenticationSignUp() http.HandlerFunc {
	type RequestBody struct {
		LastName  string `json:"last_name"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := &RequestBody{}
		if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		user := &models.User{
			FirstName: requestBody.FirstName,
			LastName:  requestBody.LastName,
			Username:  requestBody.Username,
			Password:  requestBody.Password,
		}
		if returning, err := s.storage.Repository().AddUser(user); err != nil {
			s.ErrorRespond(w, r, http.StatusUnprocessableEntity, err)
		} else {
			s.Respond(w, r, http.StatusCreated, returning)
		}
	}
}

func (s *ApiServer) HandleAuthenticationSignIn() http.HandlerFunc {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := &RequestBody{}
		if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
			s.ErrorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		if user, err := s.storage.Repository().GetUserByUsernamePassword(requestBody.Username, requestBody.Password); err != nil {
			s.ErrorRespond(w, r, http.StatusUnprocessableEntity, err)
		} else {
			creationToken, cte := s.tokenSigner.GenerateCreationToken(user)
			regularToken, rte := s.tokenSigner.GenerateRegularToken(user)
			if cte != nil || rte != nil {
				s.ErrorRespond(w, r, http.StatusUnprocessableEntity, fmt.Errorf("Errors: %v", cte, rte))
				return
			}

			s.Respond(w, r, http.StatusCreated, struct {
				CreationToken string `json:"creation_token"`
				RegularToken  string `json:"regular_token"`
			}{
				CreationToken: creationToken,
				RegularToken:  regularToken,
			})
		}
	}
}

func (s *ApiServer) HandleAuthenticationSignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *ApiServer) HandleAuthenticationValidate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		accessToken := params.Get("accessToken")
		if data, err := s.tokenSigner.ValidateRegularToken(accessToken); err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, err)
		} else {
			s.Respond(w, r, http.StatusOK, data)
		}
	}
}
