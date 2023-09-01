package tango_test

import (
	tango "tango"
	"testing"
)

func TestOrder(t *testing.T) {
	// Setup Passing Test
	data := tango.CreateOrderData{
		Amount:             1.00,
		CustomerIdentifier: "UID",
		EmailSubject:       "",
		Etid:               "test",

		// Optional
		ExternalRefID: "test",
		Message:       "",
		Recipient: tango.Person{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "John@Doe.com",
			Address: tango.Address{
				StreetLine1:     "123 Main St",
				StreetLine2:     "",
				City:            "Anytown",
				StateOrProvince: "CA",
				PostalCode:      "12345",
				Country:         "USA",
			},
		},
		Sender: tango.Sender{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "John@Doe.com",
		},
		Campaign: "test",
		Notes:    "test",
		Utid:     "test",
	}

	client := tango.TangoClient{
		Environment:       "sandbox",
		Token:             "x123",
		SendEmail:         false,
		AccountIdentifier: "123456",
	}

	// Test that Should Pass
	_, err := client.Order(data)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

}
