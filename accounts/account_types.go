package accounts

import (
	client "github.com/banjoh/fake-api-client"
	"github.com/google/uuid"
)

type Attributes struct {
	Country                 string   `json:"country,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BIC                     string   `json:"bic,omitempty"`
	IBAN                    string   `json:"iban,omitempty"`
	CustomerID              string   `json:"customer_id,omitempty"`
	Name                    []string `json:"name,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	AccountClassification   string   `json:"account_classification,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
	Status                  string   `json:"status,omitempty"`
}

type Account struct {
	Type           string      `json:"type,omitempty"`
	ID             *uuid.UUID  `json:"id,omitempty"`
	Version        *int        `json:"version,omitempty"`
	OrganisationID *uuid.UUID  `json:"organisation_id,omitempty"`
	Attributes     *Attributes `json:"attributes,omitempty"`
	CreatedOn      string      `json:"created_on,omitempty"` // DEBT: Parse string to time struct
	ModifiedOn     string      `json:"modified_on,omitempty"`
}

type AccountDTO struct {
	Data Account `json:"data"`
}

type AccountCreate struct {
	Type           string      `json:"type,omitempty"`
	ID             *uuid.UUID  `json:"id,omitempty"`
	Version        *int        `json:"version,omitempty"`
	OrganisationID *uuid.UUID  `json:"organisation_id,omitempty"`
	Attributes     *Attributes `json:"attributes,omitempty"`
}

type AccountCreateDTO struct {
	Data AccountCreate `json:"data"`
}

type Resource struct {
	BaseURL      string
	client       client.HTTPClient
	retrySleeper client.RetrySleeper
}
