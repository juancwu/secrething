# Roadmap for Rebuilding Konbini Server

This document outlines a comprehensive plan to reset and rebuild the `/server` code with improved maintainability and flexibility.

## Phase 1: Analysis and Design (Planning Phase)

### Step 1: Domain Analysis

- Identify core domains from existing code
  - Authentication (users, tokens, TOTP)
  - Bentos (secret containers and their ingredients)
  - Groups (sharing and permissions)
- Document domain relationships and dependencies

### Step 2: Architecture Design

- Design layered domain-oriented architecture:
  - Domain layer (business logic)
  - Application layer (use cases)
  - Infrastructure layer (database, cache, etc.)
  - API layer (HTTP handlers)
- Define clear interfaces between layers
- Create directory structure blueprint

### Step 3: Define Core Interfaces

- Authentication interfaces
- Data access interfaces
- Service interfaces between domains

## Phase 2: Foundation Building

### Step 1: Setup Project Structure

```
server/
├── domain/
│   ├── auth/
│   ├── bento/
│   ├── group/
│   └── common/
├── application/
├── infrastructure/
├── api/
└── cmd/
```

### Step 2: Core Infrastructure

1. Database connectivity

   - Implement DB connector similar to existing `/server/db/connection.go`
   - Create transaction helpers

2. Error handling system

   - Implement structured error types
   - Create domain-specific error factories

3. Validation

   - Port existing validator middleware
   - Add domain-specific validation rules

4. Caching infrastructure
   - Port existing memcache implementation

## Phase 3: Domain Implementation (One at a Time)

### Step 1: Auth Domain

1. Models - User, Token, TOTP
2. Repositories - Data access
3. Services - Business logic
   - Authentication
   - Token management
   - TOTP functionality
4. API handlers
   - Login, Register, TOTP setup

### Step 2: Bento Domain

1. Models - Bento, Ingredient
2. Repositories
3. Services
   - Bento creation/management
   - Ingredient management
4. API handlers
   - CRUD operations for bentos
   - Ingredient management

### Step 3: Group Domain

1. Models - Group, Membership
2. Repositories
3. Services
   - Group management
   - Invitation system
   - Permissions
4. API handlers
   - Group CRUD
   - Invitation functionality

## Phase 4: Cross-Domain Integration

1. Implement domain service interfaces
2. Setup permissions system
3. Implement event system for cross-domain communication

## Phase 5: API Layer and Server Setup

1. Route configuration
2. Middleware setup
   - Authentication
   - Rate limiting
   - Request validation
   - Logging
3. Server configuration

## Phase 6: Testing and Validation

1. Unit tests for domain logic
2. Integration tests for API
3. Manual testing against existing API

## Detailed Implementation Plan

Let's dive deeper into the concrete steps to rebuild each part:

### Step 1: Project Setup

Set up the new directory structure and create initial base files:

```bash
mkdir -p server/domain/{auth,bento,group,common}/{models,repository,service}
mkdir -p server/infrastructure/{db,cache,validator,errors}
mkdir -p server/api/{handlers,middleware,routes}
mkdir -p server/cmd/server
```

### Step 2: Core Infrastructure Implementation

#### Database Connector

```go
// server/infrastructure/db/connector.go
package db

import (
	"database/sql"
	"github.com/tursodatabase/go-libsql"
)

type Connector interface {
	Connect() (*sql.DB, error)
	Close() error
}

type SQLConnector struct {
	dbURL      string
	authToken  string
	connection *sql.DB
}

func NewConnector(dbURL, authToken string) *SQLConnector {
	return &SQLConnector{
		dbURL:     dbURL,
		authToken: authToken,
	}
}

func (c *SQLConnector) Connect() (*sql.DB, error) {
	if c.connection != nil {
		return c.connection, nil
	}

	connector, err := libsql.NewConnector(c.dbURL, libsql.WithAuthToken(c.authToken))
	if err != nil {
		return nil, err
	}

	c.connection = sql.OpenDB(connector)
	return c.connection, nil
}

func (c *SQLConnector) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}
```

#### Error Handling System

```go
// server/infrastructure/errors/errors.go
package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "validation"
	ErrorTypeDatabase     ErrorType = "database"
	ErrorTypeAuthorization ErrorType = "authorization"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeRateLimit    ErrorType = "rate_limit"
	ErrorTypeInternal     ErrorType = "internal"
)

type AppError struct {
	Type           ErrorType  `json:"-"`
	Code           int        `json:"code"`
	PublicMessage  string     `json:"message"`
	Errors         []string   `json:"errors,omitempty"`
	RequestId      string     `json:"request_id,omitempty"`
	PrivateMessage string     `json:"-"`
	InternalError  error      `json:"-"`
}

func (e AppError) Error() string {
	if e.InternalError != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.PublicMessage, e.InternalError.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.PublicMessage)
}

// Error factories
func NewValidationError(message string, details []string, err error) AppError {
	return AppError{
		Type:           ErrorTypeValidation,
		Code:           http.StatusBadRequest,
		PublicMessage:  message,
		Errors:         details,
		PrivateMessage: "Validation error",
		InternalError:  err,
	}
}

// Additional error factory methods...
```

### Step 3: Domain Implementation - Auth Domain

#### Auth Models

```go
// server/domain/auth/models/user.go
package models

import (
	"time"
)

type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Password      string     `json:"-"` // Never expose
	Nickname      string     `json:"nickname"`
	EmailVerified bool       `json:"email_verified"`
	TotpSecret    *string    `json:"-"` // Never expose
	TotpLocked    bool       `json:"totp_enabled"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Additional models for Token, TOTP, etc.
```

#### Auth Repository Interface

```go
// server/domain/auth/repository/user_repository.go
package repository

import (
	"context"
	"github.com/juancwu/konbini/server/domain/auth/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (string, error)
	Update(ctx context.Context, user *models.User) error
	SetEmailVerified(ctx context.Context, userID string, verified bool) error
	SetTOTPSecret(ctx context.Context, userID string, secret *string) error
	SetTOTPLocked(ctx context.Context, userID string, locked bool) error
}
```

#### Auth Service Implementation

```go
// server/domain/auth/service/auth_service.go
package service

import (
	"context"
	"github.com/juancwu/konbini/server/domain/auth/models"
	"github.com/juancwu/konbini/server/domain/auth/repository"
	"github.com/juancwu/konbini/server/infrastructure/errors"
	"github.com/juancwu/konbini/server/infrastructure/utils"
)

type AuthService interface {
	Register(ctx context.Context, email, password, nickname string) (*models.User, error)
	Login(ctx context.Context, email, password string, totpCode *string) (*models.AuthToken, error)
	VerifyEmail(ctx context.Context, token string) error
	// Additional methods
}

type authServiceImpl struct {
	userRepo      repository.UserRepository
	tokenRepo     repository.TokenRepository
	emailService  EmailService
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	emailService EmailService,
) AuthService {
	return &authServiceImpl{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		emailService: emailService,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, email, password, nickname string) (*models.User, error) {
	// Check if user exists
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to check user existence", err)
	}

	if exists {
		return nil, errors.NewValidationError("Email already registered", nil, nil)
	}

	// Hash password
	passwordHash, err := utils.GeneratePasswordHash(password)
	if err != nil {
		return nil, errors.NewInternalError("Failed to hash password", err)
	}

	// Create user
	user := &models.User{
		Email:         email,
		Password:      passwordHash,
		Nickname:      nickname,
		EmailVerified: false,
		TotpLocked:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create user", err)
	}
	user.ID = userID

	// Send verification email
	go s.emailService.SendVerificationEmail(email, userID)

	return user, nil
}

// Additional method implementations
```

### Step 4: API Layer Implementation

```go
// server/api/handlers/auth_handler.go
package handlers

import (
	"github.com/juancwu/konbini/server/domain/auth/service"
	"github.com/juancwu/konbini/server/infrastructure/errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=12,max=32"`
		Nickname string `json:"nickname" validate:"required,min=3,max=32"`
	}

	if err := c.Bind(&req); err != nil {
		return errors.NewValidationError("Invalid request format", nil, err)
	}

	if err := c.Validate(req); err != nil {
		return err // Assuming validation middleware handles this
	}

	user, err := h.authService.Register(c.Request().Context(), req.Email, req.Password, req.Nickname)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user_id": user.ID,
	})
}

// Additional handler methods
```

### Step 5: Server Configuration and Wiring

```go
// server/cmd/server/main.go
package main

import (
	"github.com/juancwu/konbini/server/api/handlers"
	"github.com/juancwu/konbini/server/api/middleware"
	"github.com/juancwu/konbini/server/api/routes"
	"github.com/juancwu/konbini/server/domain/auth/repository"
	"github.com/juancwu/konbini/server/domain/auth/service"
	"github.com/juancwu/konbini/server/infrastructure/config"
	"github.com/juancwu/konbini/server/infrastructure/db"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup database connection
	dbConnector := db.NewConnector(cfg.DatabaseURL, cfg.DatabaseAuthToken)
	conn, err := dbConnector.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbConnector.Close()

	queries := db.New(conn)

	// Setup repositories
	userRepo := repository.NewUserRepository(queries)
	tokenRepo := repository.NewTokenRepository(queries)
	// Additional repositories

	// Setup services
	emailService := service.NewEmailService(cfg.EmailConfig)
	authService := service.NewAuthService(userRepo, tokenRepo, emailService)
	// Additional services

	// Setup HTTP handlers
	authHandler := handlers.NewAuthHandler(authService)
	// Additional handlers

	// Setup Echo server
	e := echo.New()
	e.Validator = middleware.NewValidator()
	e.HTTPErrorHandler = middleware.ErrorHandler()

	// Register middleware
	middleware.RegisterMiddleware(e, cfg)

	// Register routes
	routes.RegisterAuthRoutes(e, authHandler)
	// Additional routes

	// Start server
	log.Info().Msgf("Starting server on %s", cfg.Port)
	if err := e.Start(cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
```

## Key Improvements Over Current Implementation

1. **Clear Layer Separation**

   - Domain logic isolated from HTTP concerns
   - Repository pattern abstracts data access
   - Service layer encapsulates business rules

2. **Interface-Based Design**

   - Dependencies flow through interfaces
   - Easy testing with mocks
   - Reduced coupling between components

3. **Consistent Error Handling**

   - Centralized error definitions
   - Domain-specific error types
   - Clear separation between public and private error information

4. **Simplified Handler Logic**

   - Handlers focus on HTTP concerns only
   - Business logic delegated to services
   - Smaller, more focused functions

5. **Improved Testability**

   - Each layer can be tested in isolation
   - Clean interfaces support mocking
   - Core business logic separate from infrastructure

6. **Explicit Dependencies**
   - Dependencies explicitly passed via constructors
   - No global state or singletons
   - Clear ownership of resources

## Implementation Strategy

1. **Incremental Approach**

   - Build the core infrastructure first
   - Implement one domain at a time
   - Start with auth as it's the foundation

2. **Feature Parity**

   - Ensure each domain maintains feature parity with existing code
   - Reference existing implementations for business rules

3. **Parallel Development**
   - Keep existing server running during development
   - Test new implementation against same database
   - Gradual switchover when ready

