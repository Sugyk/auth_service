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
