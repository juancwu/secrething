package router

// reqBodyValidationMsgs is a collection of all request body validation error messages.
var reqBodyValidationMsgs = map[string]string{
	// signupReqBody messages
	"signupReqBody.Email.required":    "Missing email in request body",
	"signupReqBody.Email.email":       "Invalid email provided",
	"signupReqBody.Password.required": "Missing password in request body",
	"signupReqBody.Password.password": "Invalid password format. A password must of at least 12 characters and combines at least one special character, uppercase and lowercase letter, and number",
	"signupReqBody.Name.required":     "Missing name in request body",
	"signupReqBody.Name.min":          "Name must be at least 3 characters long",
	"signupReqBody.Name.max":          "Name must not be longer than 30 characters",
	"signupReqBody.Name.alpha":        "Name must only include alphabet characters",
}

// signupReqBody represents the request body that is expected when handling a signup request.
type signupReqBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Name     string `json:"name" validate:"required,min=3,max=30,alpha"`
}
