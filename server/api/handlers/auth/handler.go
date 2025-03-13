package auth

import (
	"github.com/juancwu/konbini/server/config"
	"github.com/juancwu/konbini/server/db"
)

type AuthHandler struct {
	config *config.Config
	db     *db.TursoConnector
}

func NewAuthHandler(cfg *config.Config, db *db.TursoConnector) *AuthHandler {
	return &AuthHandler{
		config: cfg,
		db:     db,
	}
}
