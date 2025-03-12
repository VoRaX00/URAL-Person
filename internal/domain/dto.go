package domain

type GetPerson struct {
	Login string `json:"login"`
	Image []byte `json:"image"`
}

type RegisterPerson struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
