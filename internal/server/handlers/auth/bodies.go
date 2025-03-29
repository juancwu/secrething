package auth

import "github.com/juancwu/go-valkit/v2/validator"

type createUserBody struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,password"`
	Name     *string `json:"name" validate:"omitnil,omitempty,max=50"`
}

func (b createUserBody) GetMessages() validator.ValidationMessages {
	msgs := validator.NewValidationMessages()
	msgs.SetMessage("email", "required", "Email is required.")
	msgs.SetMessage("email", "email", "'{2}' is not a valid email.")
	msgs.SetMessage("password", "required", "Password is required.")
	msgs.SetMessage("name", "max", "Name must not be longer than {2} characters.")
	return msgs
}

type signinBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func (b signinBody) GetMessages() validator.ValidationMessages {
	msgs := validator.NewValidationMessages()
	msgs.SetMessage("email", "required", "Email is required.")
	msgs.SetMessage("email", "email", "'{2}' is not a valid email.")
	msgs.SetMessage("password", "required", "Password is required.")
	return msgs
}
