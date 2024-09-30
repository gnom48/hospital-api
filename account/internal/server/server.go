package server

import (
	"account/internal/server/tokens"
	"account/internal/storage"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	_ "account/internal/server/docs"

	models "github.com/gnom48/hospital-api-lib"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type ApiServer struct {
	config      *Config
	logger      *logrus.Logger
	router      *mux.Router
	storage     *storage.Storage
	tokenSigner tokens.TokenSigner
}

func New(config *Config) *ApiServer {
	return &ApiServer{
		config:      config,
		logger:      logrus.New(),
		router:      mux.NewRouter(),
		tokenSigner: &tokens.TokenSign{},
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
				s.ErrorRespond(w, r, http.StatusInternalServerError, fmt.Errorf("Error: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// @title Account
// @version 1.0
// @description Account API documentation
// @host localhost:8081
// @BasePath /
// @schemes http
func (s *ApiServer) ConfigureRouter() {
	s.router.Use(s.internalServerErrorMiddleware)
	s.router.Use(s.loggingMiddleware)
	s.router.HandleFunc("/hello", s.AuthMiddleware(s.HandleHello()))

	s.router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}

type StringContextKey string

var UserContextKey StringContextKey = "user"

func (s *ApiServer) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is empty", http.StatusUnauthorized)
			return
		}

		claims, err := s.tokenSigner.ValidateRegularToken(tokenString)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
			return
		}

		user, err := s.storage.Repository().GetUserById(claims.UserId)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, *user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// @Summary Hello endpoint
// @Description Returns a greeting message
// @Accept json
// @Produce json
// @Success 200 {string} string "Successful response"
// @Router /hello [get]
// @Param Authorization header string true "Authorization header"
func (s *ApiServer) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, "+r.Context().Value(UserContextKey).(models.User).Username+"!")
	}
}
