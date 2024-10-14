package entity

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmSignUpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
