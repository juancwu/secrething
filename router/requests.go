// This file contains structs that are use as request bodies for all routes.
package router

// signupRequest represents the request body to sign up.
type signupRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	FirstName *string `json:"first_name" validate:"omitnil,alpha"`
	LastName  *string `json:"last_name" validate:"omitnil,alpha"`
	Password  string  `json:"password" validate:"required,min=12"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
