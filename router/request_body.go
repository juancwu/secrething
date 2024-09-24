package router

// signupReqBody represents the request body that is expected when handling a signup request.
type signupReqBody struct {
	Email    string `json:"email" validate:"required,email" errormsg:"required=Missing email in request body,email=Invalid email"`
	Password string `json:"password" validate:"required,password" errormsg:"required=Missing password in request body,password=Invalid password format. A password must of at least 12 characters and combines at least one special character, uppercase and lowercase letter, and number"`
	Name     string `json:"name" validate:"required,min=3,max=30,alpha" errormsg:"required=Missing name in request body,min=Name must be 3 to 30 characters long,max=Name must be 3 to 30 characters long,alpha=Name must only contain letters"`
}

// signinReqBody represents the request body that is expected when handling a signin request.
type signinReqBody struct {
	// Email is the username that is used to signin
	Email string `json:"email" validate:"required,email" errormsg:"required=Missing email,email=Invalid email"`
	// Password no password validation here because there is no need to gives hints about it when signing in.
	Password string `json:"password" validate:"required" errormsg:"required=Missing password"`
}

// resendVerificationEmailReqBody represents the request body that is expected when handling a resend ev request.
type resendVerificationEmailReqBody struct {
	Email string `json:"email" validate:"required,email" errormsg:"Invalid email"`
}

// newBentoReqBody represents the request body that is expected when handling a new bento requets.
type newBentoReqBody struct {
	Name        string       `json:"name" validate:"required,min=3,max=50,printascii" errormsg:"required=Missing name in request body,min=Name must be longer than 3 characters,max=Name must not be longer than 50 characters,printascii=Name must only contain printable ascii characters"`
	PubKey      string       `json:"pub_key" validate:"required" errormsg:"Missing pub_key"`
	Ingridients []Ingridient `json:"ingridients,omitempty" validate:"omitnil,dive"`
}

// addIngridientsReqBody represents the request body that is expected when handling adding a new Ingridient to a prepared bento.
type addIngridientsReqBody struct {
	BentoId     string       `json:"bento_id" validate:"required,uuid4" errormsg:"required=Missing bento id,uuid4=Invalid bento id; Only UUID v4"`
	Ingridients []Ingridient `json:"ingridients" validate:"required,gt=0,dive" errormsg:"required=Missing ingridients,gt=Missing ingridients,__default=Invalid ingridients field"`
	Challenge   string       `json:"challenge" validate:"required" errormsg:"required=Missing challenge"`
	Signature   string       `json:"signature" validate:"required" errormsg:"required=Missing signature"`
}

// Ingridient is used in the addIngridientsReqBody.
type Ingridient struct {
	Name  string `json:"name" validate:"required,printascii" errormsg:"required=Missing ingridient name,printascii=Ingridient name can only contain printable ascii"`
	Value string `json:"value" validate:"required" errormsg:"Missing ingridient value"`
}

// A challenger represents the pair challenge and signature that is used to verify someone
// who is making a request in fact has access to the private key of a target bento. This
// helps in verifying that the user might be a trusful user. More permissions have to be checked
// to have a better understanding on who the server is dealing with and whether said requested action
// is to be performed or denied.
type challenger struct {
	Challenge string `json:"challenge" validate:"required" errormsg:"Missing challenge"`
	Signature string `json:"signature" validate:"required" errormsg:"Missing signature"`
}

// renameBentoReqBody represents the request body that is expected when handling rename bento requests.
type renameBentoReqBody struct {
	BentoId string `json:"bento_id" validate:"required,uuid4" errormsg:"required=Missing bento id,uuid4=Invalid bento id; Only UUID v4"`
	NewName string `json:"new_name" validate:"required,min=3,max=50,printascii" errormsg:"required=Missing new bento name,min|max=New name must be 3 to 50 characters long,printascii=New name can only contain printable ascii"`
}

type shareBentoReqBody struct {
	BentoId      string `json:"bento_id" validate:"required,uuid4" errormsg:"required|uuid4=Invalid bento id; Must be UUID v4"`
	ShareToEmail string `json:"share_to_email" validate:"required,email" errormsg:"required|email=Invalid email"`
	Challenge    string `json:"challenge" validate:"required" errormsg:"Missing challenge"`
	Signature    string `json:"signature" validate:"required" errormsg:"Missing signature"`
	// PermissionLevels represents the up to what level of permission should the target user get.
	// Keep in mind that the max level of permission that can be granted is the max level of the requesting user.
	// The permission to grant share is not included. That needs to be exclusively granted in another route.
	PermissionLevels []string `json:"permission_levels" validate:"omitnil,options=all write delete share rename_bento rename_ingridient write_ingridient delete_ingridient" errormsg:"options=Invalid permission levels. Please refer to 'https://github.com/juancwu/konbini/blob/main/.github/docs/DOCUMENTATION.md' for the available permission levels and its usage."`
}

type updateEmailReqBody struct {
	NewEmail string `json:"new_email" validate:"required,email" errormsg:"required|email=Invalid email"`
}

type renameIngridientReqBody struct {
	BentoId    string     `json:"bento_id" validate:"required,uuid4" errormsg:"Invalid bento id."`
	Challenger challenger `json:"challenger" validate:"required" errormsg:"Missing challenger."`
	OldName    string     `json:"old_name" validate:"required,printascii" errormsg:"required=Missing old ingridient name,printascii=Old name can only consist of printable ASCII characters"`
	NewName    string     `json:"new_name" validate:"required,printascii" errormsg:"required=Missing new ingridient name,printascii=New name can only consist of printable ASCII characters"`
}
