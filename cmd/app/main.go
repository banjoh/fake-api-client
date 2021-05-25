package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/banjoh/fake-api-client/accounts"
	"github.com/google/uuid"
)

// NOTE: This implementation is for my own testing Remove me please
func main() {
	attr := accounts.Attributes{
		Country: "GB",
		Name:    []string{"John Doe"},
	}

	id := uuid.MustParse("ad27e266-9605-4b4b-a0e5-3003ea9cc4dc")
	oId := uuid.MustParse("eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")

	accCreate := accounts.AccountCreate{
		Type:           "accounts",
		ID:             &id,
		OrganisationID: &oId,
		Attributes:     &attr,
	}

	accClient := accounts.New()
	ctx := context.Background()

	acc, err := accClient.Create(ctx, &accCreate)
	if err != nil {
		fmt.Printf("Failed creating acc: err=%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created = %v\n", toJSON(acc))

	acc, err = accClient.Fetch(ctx, id)
	if err != nil {
		fmt.Printf("Failed fetching acc: id=%s, err=%v\n", id, err)
		os.Exit(1)
	}

	fmt.Printf("Fetch = %v\n", toJSON(acc))

	err = accClient.Delete(ctx, id, *acc.Version)
	if err != nil {
		fmt.Printf("Failed deleting acc: id=%s, err=%v\n", id, err)
		os.Exit(1)
	}

	fmt.Printf("Delete = %v\n", id)
}

func toJSON(acc *accounts.Account) string {
	data, _ := json.MarshalIndent(acc, "", "  ")
	return string(data)
}
