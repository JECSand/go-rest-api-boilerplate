package server

import (
	"github.com/JECSand/go-rest-api-boilerplate/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// Server is a struct that stores the API Apps high level attributes such as the router, config, and services
type Server struct {
	Router       *mux.Router
	TokenService *services.TokenService
	UserService  services.UserService
	GroupService services.GroupService
	TaskService  services.TaskService
	FileService  services.FileService
}

// NewServer is a function used to initialize a new Server struct
func NewServer(u services.UserService, g services.GroupService, tt services.TaskService, f services.FileService, t *services.TokenService) *Server {
	router := mux.NewRouter().StrictSlash(true)
	router = NewGroupRouter(router, t, g, u)
	router = NewUserRouter(router, t, u, g, f)
	router = NewTaskRouter(router, t, tt)
	return &Server{
		Router:       router,
		TokenService: t,
		UserService:  u,
		GroupService: g,
		TaskService:  tt,
		FileService:  f,
	}
}

// Start starts the initialized Server
func (s *Server) Start() {
	log.Println("Listening on port " + os.Getenv("PORT"))
	go func() {
		if err := http.ListenAndServe(":"+os.Getenv("PORT"), handlers.LoggingHandler(os.Stdout, s.Router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)
}
