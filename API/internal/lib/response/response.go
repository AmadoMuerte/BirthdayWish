package response

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
