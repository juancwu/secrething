package main

import (
	"konbini/db"
	"konbini/routes"
	serverconfig "konbini/server_config"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := serverconfig.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load server configuration")
	}

	dbUrl, dbAuthToken := cfg.GetDatabaseConfig()
	conn, err := db.NewConnection(dbUrl, dbAuthToken)
	defer conn.Close()
	queries := db.New(conn)

	e := echo.New()

	// v1 routes
	apiV1 := e.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Echo:         apiV1,
		ServerConfig: cfg,
		DatabaseConn: conn,
		Queries:      queries,
	}
	routes.SetupRoutesV1(routeConfig)

	err = e.Start(cfg.GetPort())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server.")
	}
}
