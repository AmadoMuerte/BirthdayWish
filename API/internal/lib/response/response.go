package response

type UserResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type ErrResponse struct {
	Message string `json:"message,omitempty"`
}

func Error(msg string) ErrResponse {
	return ErrResponse{
		Message: msg,
	}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}
