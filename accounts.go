package client

import (
	"fmt"

	"github.com/google/uuid"
)

type Account struct {
	Type           string
	ID             uuid.UUID
	OrganisationID uuid.UUID
	Version        int
}

func CreateAccount() (*Account, error) {
	return nil, fmt.Errorf("unimplemented")
}

func FetchAccount() (*Account, error) {
	return nil, fmt.Errorf("unimplemented")
}

func DeleteAccount() error {
	return fmt.Errorf("unimplemented")
}
