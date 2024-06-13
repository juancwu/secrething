// This file contains structs that are use as request bodies for all routes.
package router

// signupRequest represents the request body to sign up.
type signupRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	FirstName *string `json:"first_name" validate:"omitnil,alpha"`
	LastName  *string `json:"last_name" validate:"omitnil,alpha"`
	Password  string  `json:"password" validate:"required,min=12"`
}

// loginRequest represents the request body when logging in.
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// prepBentoRequest represents the request body when prepping a new bento.
type prepBentoRequest struct {
	Name string `json:"name"`
	// base64 encoded public key
	PubKey string `json:"pub_key"`
}

// addIngridientRequest represents the request body when adding a new ingridien (secret) to an existing bento.
type addIngridientRequest struct {
	Key             string `json:"key"`
	Value           string `json:"value"`
	BentoId         string `json:"bento_id"`
	SignedChallenge string `json:"signed_challenge"`
}
