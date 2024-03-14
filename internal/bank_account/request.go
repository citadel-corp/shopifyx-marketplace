package bankaccount

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateBankAccountPayload struct {
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

func (p CreateBankAccountPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.BankName, validation.Required, validation.Length(5, 15)),
		validation.Field(&p.BankAccountName, validation.Required, validation.Length(5, 15)),
		validation.Field(&p.BankAccountNumber, validation.Required, validation.Length(5, 15)),
	)
}
