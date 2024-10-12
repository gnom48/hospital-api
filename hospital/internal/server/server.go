package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hospital/internal/storage"
	"io/ioutil"
	"net/http"
	"strings"

	_ "hospital/internal/server/docs"

	models "github.com/gnom48/hospital-api-lib"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type ApiServer struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *Config) *ApiServer {
	return &ApiServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *ApiServer) Start() error {
	if err := s.ConfigureLogger(); err != nil {
		return err
	}
	s.logger.Info("Logger configured")

	s.ConfigureRouter()
	s.logger.Info("Router configured")

	if err := s.ConfigureStore(); err != nil {
		s.logger.Error(err)
		return err
	}
	s.logger.Info("Storage configured")

	s.logger.Info("Starting ApiServer")
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *ApiServer) ConfigureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *ApiServer) ConfigureStore() error {
	st := storage.New(s.config.StorageConfig)
	if err := st.Open(); err != nil {
		return err
	}

	s.storage = st

	return nil
}

func (s *ApiServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := &strings.Builder{}
		headers.Write([]byte("["))
		if s.config.LogHeaders {
			for key, values := range r.Header {
				for _, value := range values {
					headers.Write([]byte(key + " = " + value + ", "))
				}
			}
		}
		headers.Write([]byte("]"))

		bodyBytes := make([]byte, 0)
		if s.config.LogBody {
			bodyBytes, _ = ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		queryParams := ""
		if s.config.LogQueryParams {
			queryParams = r.URL.Query().Encode()
		}

		s.logger.Info("Method: " + r.Method + " | Path: " + r.URL.Path + " | Headers: " + headers.String() + " | Body: " + string(bodyBytes) + " | Query: " + queryParams)

		next.ServeHTTP(w, r)
	})
}

func (s *ApiServer) internalServerErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.ErrorRespond(w, r, http.StatusNotImplemented, fmt.Errorf("Error: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// @title Account
// @version 1.0
// @description Account API documentation
// @host localhost:8082
// @BasePath /
// @schemes http
func (s *ApiServer) ConfigureRouter() {
	s.router.Use(s.internalServerErrorMiddleware)
	s.router.Use(s.loggingMiddleware)

	s.router.PathPrefix("/api/Hospitals/{id}/Rooms").Handler(s.AuthByTokenMiddleware(s.HandleGetRoomsByHospitalId())).Methods("GET")
	s.router.PathPrefix("/api/Hospitals/{id}").Handler(s.AuthByTokenMiddleware(s.HandleGetHospitalById())).Methods("GET")
	s.router.PathPrefix("/api/Hospitals/{id}").Handler(s.AuthByTokenMiddleware(s.HandleUpdateHospital())).Methods("PUT")
	s.router.PathPrefix("/api/Hospitals/{id}").Handler(s.AuthByTokenMiddleware(s.HandleSoftDeleteHospital())).Methods("DELETE")
	s.router.PathPrefix("/api/Hospitals").Handler(s.AuthByTokenMiddleware(s.HandleGetHospitals())).Methods("GET")
	s.router.PathPrefix("/api/Hospitals").Handler(s.AuthByTokenMiddleware(s.HandleCreateHospital())).Methods("POST")

	s.router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}

type StringContextKey string

var UserContextKey StringContextKey = "user_role"

type InfoAboutMeResponseBody struct {
	Token string        `json:"token"`
	User  models.User   `json:"user"`
	Roles []models.Role `json:"roles"`
}

func (s *ApiServer) AuthByTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is empty", http.StatusUnauthorized)
			return
		}

		url := "http://account-service:8081/api/Accounts/Me"
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, err)
			return
		}

		req.Header.Set("Authorization", tokenString)

		resp, err := client.Do(req)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, err)
			return
		}

		var responseMap map[string]interface{}
		if err := json.Unmarshal(body, &responseMap); err != nil {
			s.ErrorRespond(w, r, http.StatusInternalServerError, fmt.Errorf("Authorization failed: invalid token"))
			return
		}
		if _, exists := responseMap["server_error"]; exists {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf(responseMap["server_error"].(string)))
			return
		}

		var aboutMeResponse InfoAboutMeResponseBody
		if err := json.Unmarshal(body, &aboutMeResponse); err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, err)
			return
		}

		aboutMeResponse.Token = tokenString
		ctx := context.WithValue(r.Context(), UserContextKey, aboutMeResponse)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
