package router

// apiResponse is a generic json response for all APIs
type apiResponse struct {
	StatusCode int    `json:"status_code"`
	Message    any    `json:"message"`
	RequestId  string `json:"request_id"`
}

// loginResponse represents the json body that is sent back when
// a user successfully logs in.
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// newTokenResponse represents the json body that is sent back when
// a new access token is generated.
type newTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// prepBentoResponse represents the json body that is sent back when a new bento gets prep.
type prepBentoResponse struct {
	BentoId string `json:"bento_id"`
}

// getChallengeResponse represents the json body that is sent back when getting a new challenge.
type getChallengeResponse struct {
	ChallengeId string `json:"challenge_id"`
	Challenge   string `json:"challenge"`
}
