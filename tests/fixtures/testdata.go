package fixtures

import "example.com/m/src/domain"

var ValidAccount = &domain.Account{
	Number:         "1111111111111111",
	ExparationDate: "12-34",
	FullName:       "MISHA ANIKUTIN",
	CVV:            123,
}

var InvalidAccountNumber = "1234567891011121"
