package api

import (
	"github.com/juancwu/go-valkit/v2/validations"
	"github.com/juancwu/go-valkit/v2/validator"
	"github.com/juancwu/secrething/internal/config"
	"github.com/juancwu/secrething/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// API handles all the routes and controllers
type API struct {
	Echo   *echo.Echo
	Config *config.Config
	DB     *db.Queries
	Valkit *validator.Validator
}

// New creates a new API instance
func New(cfg *config.Config) *API {
	conn, err := db.Connect(cfg)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	if err := conn.Ping(); err != nil {
		panic("Failed to ping database: " + err.Error())
	}
	queries := db.New(conn)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		AllowCredentials: true,
		MaxAge:           300,
	}))

	e.HTTPErrorHandler = errorHandler

	valkit := validator.New()
	valkit.UseJsonTagName()

	valkit.SetDefaultMessage("{field} is invalid.")
	valkit.SetDefaultTagMessage("required", "{field} is required.")
	valkit.SetDefaultTagMessage("min", "{field} length must be at least {param}.")
	valkit.SetDefaultTagMessage("max", "{field} length must be at most {param}.")
	valkit.SetDefaultTagMessage("email", "{value} is not a valid email address.")
	validations.AddPasswordValidation(valkit, validations.DefaultPasswordOptions())

	api := &API{
		Echo:   e,
		Config: cfg,
		DB:     queries,
		Valkit: valkit,
	}

	api.registerRoutes()

	return api
}

// Start starts the server and listens on the address
func (api *API) Start(addr string) error {
	return api.Echo.Start(addr)
}

// registerRoutes registers all API routes
func (api *API) registerRoutes() {
	// Health check endpoint
	api.Echo.GET("/health", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	// Public routes
	api.registerAuthRoutes()

	// Protected routes
	api.registerUserRoutes()
}

type apiResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}
