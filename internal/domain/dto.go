package domain

type GetUser struct {
	Login string `json:"login"`
	Image []byte `json:"image"`
}

type RegisterUser struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
