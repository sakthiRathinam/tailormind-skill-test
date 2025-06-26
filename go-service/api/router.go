package router

import (
	"go-service/internal/config"
	"github.com/gorilla/mux"
)

func SetupRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()


	return r
} 