package response

type MessageResponse struct {
	Message string `json:"message,omitempty"`
}

func Error(msg string) MessageResponse {
	return MessageResponse{
		Message: msg,
	}
}

func Success(msg string) MessageResponse {
	return MessageResponse{
		Message: msg,
	}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}
