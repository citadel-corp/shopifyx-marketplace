package user

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateUserPayload struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (p CreateUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Username, validation.Required, validation.Length(5, 15)),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p LoginPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Username, validation.Required, validation.Length(5, 15)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}
