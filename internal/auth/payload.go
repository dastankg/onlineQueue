package auth

type RegisterRequest struct {
	Name       string `json:"name"`
	Login      string `json:"login"`
	Password1  string `json:"password1"`
	Password2  string `json:"password2"`
	RegisterID *uint  `json:"office_id,omitempty"`
}

type RegisterResponse struct {
	Name         string `json:"name"`
	Login        string `json:"login"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Login        string `json:"login"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
