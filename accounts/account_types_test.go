package accounts

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixtures() map[string]struct {
	acc  Account
	json string
} {
	version := 0
	id := uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	oID := uuid.MustParse("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")
	jointAcc := false

	tests := map[string]struct {
		acc  Account
		json string
	}{
		"typical values": {
			json: `{
				"data": {
				  "type": "accounts",
				  "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
				  "version": 0,
				  "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				  "modified_on": "2021-05-23T16:05:52.970Z",
				  "attributes": {
					"country": "GB",
					"base_currency": "GBP",
					"account_number": "41426819",
					"bank_id": "400300",
					"bank_id_code": "GBDSC",
					"bic": "NWBKGB22",
					"iban": "GB11NWBK40030041426819",
					"status": "confirmed",
					"joint_account": false
				  }
				}
			  }`,
			acc: Account{
				Type:           "accounts",
				ID:             &id,
				Version:        &version,
				OrganisationID: &oID,
				ModifiedOn:     "2021-05-23T16:05:52.970Z",
				Attributes: &Attributes{
					Country:       "GB",
					BaseCurrency:  "GBP",
					AccountNumber: "41426819",
					BankID:        "400300",
					BankIDCode:    "GBDSC",
					BIC:           "NWBKGB22",
					IBAN:          "GB11NWBK40030041426819",
					Status:        "confirmed",
					JointAccount:  &jointAcc,
				},
			},
		},
		"empty": {
			json: `{
				"data": {}
			  }`,
			acc: Account{},
		},
		"empty attributes": {
			json: `{
				"data": {
					"attributes": {}
				}
			  }`,
			acc: Account{
				Attributes: &Attributes{},
			},
		},
	}

	return tests
}

func TestMarshalAccount(t *testing.T) {
	tests := fixtures()

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(&AccountDTO{Data: tc.acc})

			require.NoError(t, err)
			assert.JSONEq(t, tc.json, string(got))
		})
	}
}

func TestUnmarshalAccount(t *testing.T) {
	tests := fixtures()

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var got AccountDTO
			err := json.Unmarshal(bytes.NewBufferString(tc.json).Bytes(), &got)

			want := AccountDTO{Data: tc.acc}

			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}
