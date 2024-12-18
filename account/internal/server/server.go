package server

import (
	"account/internal/server/elasticsearch"
	"account/internal/server/tokens"
	"account/internal/storage"
	"bytes"
	"context"
	"fmt"
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
	config                  *Config
	logger                  *logrus.Logger
	router                  *mux.Router
	storage                 *storage.Storage
	tokenSigner             tokens.TokenSigner
	elasticsearchConnection *elasticsearch.ElasticsearchConnection
}

func New(config *Config) *ApiServer {
	return &ApiServer{
		config:                  config,
		logger:                  logrus.New(),
		router:                  mux.NewRouter(),
		tokenSigner:             &tokens.TokenSign{},
		elasticsearchConnection: &elasticsearch.ElasticsearchConnection{},
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

	// if err := s.elasticsearchConnection.Repository().CreateElasticsearchIndexes(); err != nil {
	// 	s.logger.Error("Couldn't configure elasticsearch")
	// 	s.logger.Error(err)
	// 	// return err
	// } else {
	// 	s.logger.Info("Elasticsearch configured")
	// }

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
	defer s.storage.Close()
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
				s.ErrorRespond(w, r, http.StatusNotImplemented, fmt.Errorf("Internal Server Error: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// @title Account
// @version 1.0
// @description Account API (Account microservice) documentation. Jтвечает за авторизацию и данные о пользователях. Все остальные сервисы зависят от него, ведьименно он выпускает JWT токен и проводит интроспекцию.
// @host localhost:8081
// @BasePath /
// @schemes http
func (s *ApiServer) ConfigureRouter() {
	s.router.Use(s.internalServerErrorMiddleware)
	s.router.Use(s.loggingMiddleware)

	s.router.HandleFunc("/api/Authentication/SignUp", s.HandleAuthenticationSignUp()).Methods("POST")
	s.router.HandleFunc("/api/Authentication/SignIn", s.HandleAuthenticationSignIn()).Methods("POST")
	s.router.HandleFunc("/api/Authentication/SignOut", s.AuthCreationTokenMiddleware(s.HandleAuthenticationSignOut())).Methods("HEAD")
	s.router.HandleFunc("/api/Authentication/Validate", s.HandleAuthenticationValidate()).Methods("GET")
	s.router.HandleFunc("/api/Authentication/Refresh", s.AuthCreationTokenMiddleware(s.HandleAuthenticationRefresh())).Methods("GET")

	s.router.HandleFunc("/api/Accounts/Me", s.AuthRegularTokenMiddleware(s.userRoleMiddleware(s.HandleGetCurrentAccount()))).Methods("GET")
	s.router.HandleFunc("/api/Accounts/Update", s.AuthRegularTokenMiddleware(s.HandleUpdateAccount())).Methods("PUT")
	s.router.HandleFunc("/api/Accounts", s.AuthRegularTokenMiddleware(s.userRoleMiddleware(s.HandleGetAllAccounts()))).Methods("GET")
	s.router.HandleFunc("/api/Accounts", s.AuthRegularTokenMiddleware(s.userRoleMiddleware(s.HandleCreateAccount()))).Methods("POST")
	s.router.HandleFunc("/api/Accounts/{id}", s.AuthRegularTokenMiddleware(s.userRoleMiddleware(s.HandleUpdateAccountById()))).Methods("PUT")
	s.router.HandleFunc("/api/Accounts/{id}", s.AuthRegularTokenMiddleware(s.userRoleMiddleware(s.HandleSoftDeleteAccountById()))).Methods("DELETE")

	s.router.HandleFunc("/api/Doctors", s.AuthRegularTokenMiddleware(s.userRoleMiddleware(s.HandleGetDoctors()))).Methods("GET")
	s.router.HandleFunc("/api/Doctors/{id}", s.AuthRegularTokenMiddleware(s.userRoleMiddleware(s.HandleGetDoctorById()))).Methods("GET")

	s.router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}

type StringContextKey string

var UserContextKey StringContextKey = "user"

func (s *ApiServer) AuthRegularTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is empty", http.StatusUnauthorized)
			return
		}

		claims, err := s.tokenSigner.ValidateRegularToken(tokenString)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
			return
		}

		defer s.storage.Close()
		if token, err := s.storage.Repository().GetTokenById(claims.ID); err != nil || token == nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
			return
		}

		// var user *models.User
		// if cachedUser, err := s.elasticsearchConnection.Repository().GetUserInfoByIdElasticsearch(claims.UserId); err != nil {
		// 	user, err = s.storage.Repository().GetUserById(claims.UserId)
		// 	if err != nil {
		// 		s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
		// 		return
		// 	}
		// } else {
		// 	user = &models.User{
		// 		Id:       cachedUser.Id,
		// 		Username: cachedUser.Username,
		// 		Password: cachedUser.Password,
		// 	}
		// }

		user, err := s.storage.Repository().GetUserById(claims.UserId)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *ApiServer) AuthCreationTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is empty", http.StatusUnauthorized)
			return
		}

		claims, err := s.tokenSigner.ValidateCreationToken(tokenString)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
			return
		}

		defer s.storage.Close()
		token, err := s.storage.Repository().GetTokenById(claims.ID)
		if err != nil || token == nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
			return
		}

		// var user *models.User
		// if cachedUser, err := s.elasticsearchConnection.Repository().GetUserInfoByLoginPasswordElasticsearch(claims.Username, claims.Password); err != nil {
		// 	user, err = s.storage.Repository().GetUserByUsernamePassword(claims.Username, claims.Password)
		// 	if err != nil {
		// 		s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
		// 		return
		// 	}
		// } else {
		// 	user = &models.User{
		// 		Id:       cachedUser.Id,
		// 		Username: cachedUser.Username,
		// 		Password: cachedUser.Password,
		// 	}
		// }

		user, err := s.storage.Repository().GetUserByUsernamePassword(claims.Username, claims.Password)
		if err != nil {
			s.ErrorRespond(w, r, http.StatusUnauthorized, tokenError)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

var RoleContextKey StringContextKey = "role"

type userRolesResponseBody struct {
	Roles []models.Role `json:"roles"`
}

func (s *ApiServer) userRoleMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			s.ErrorRespond(w, r, http.StatusUnauthorized, fmt.Errorf("User not found"))
			return
		}
		defer s.storage.Close()
		if userRoles, err := s.storage.Repository().GetAllUserRoles(user.Id); err != nil {
			s.ErrorRespond(w, r, http.StatusForbidden, fmt.Errorf("Role not found"))
		} else {
			ctx := context.WithValue(r.Context(), RoleContextKey, *&userRoles)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
