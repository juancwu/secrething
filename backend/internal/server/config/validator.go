package config

import "github.com/juancwu/go-valkit/v2/validator"

func DefaultValidator() *validator.Validator {
	v := validator.New().UseJsonTagName()

	v.SetDefaultMessage("This field failed to pass validation.")

	v.SetDefaultTagMessage("required", "{0} is required.")
	v.SetDefaultTagMessage("email", "Email is required.")
	v.SetDefaultTagMessage("min", "{0}'s length must be at least {2}.")
	v.SetDefaultTagMessage("max", "{0}'s length must be at most {2}.")
	v.SetDefaultTagMessage("alpha", "{0} must only consist of alphabet characters.")
	v.SetDefaultTagMessage("oneof", "{0} must be one of the following '{2}'")

	return v
}
