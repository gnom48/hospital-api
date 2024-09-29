package server

import (
	"account/internal/storage"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

func (s *ApiServer) ConfigureRouter() {
	s.router.HandleFunc("/hello", s.HandleHello())
}

func (s *ApiServer) ConfigureStore() error {
	st := storage.New(s.config.StorageConfig)
	if err := st.Open(); err != nil {
		return err
	}

	s.storage = st

	return nil
}

func (s *ApiServer) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
