package router

import (
	v1 "go-service/api/v1"
	"go-service/internal/config"

	"github.com/gorilla/mux"
)

func SetupRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()

	// Register v1 API routes
	v1.RegisterV1Routes(r, cfg)

	return r
}
