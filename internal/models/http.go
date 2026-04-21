package models

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if r.Login == "" {
		return NewValidationErr("login can not be empty")
	}
	if len([]rune(r.Password)) < 16 {
		return NewValidationErr("password must have 16 symbols or more")
	}
	return nil
}

type RegisterResponse201 struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Login == "" {
		return NewValidationErr("login can not be empty")
	}
	if len([]rune(r.Password)) < 16 {
		return NewValidationErr("password must have 16 symbols or more")
	}
	return nil
}

type LoginResponse200 struct {
	AccessToken string `json:"access_token"`
}

type LoginResponse401 struct {
	Message string `json:"message"`
}
