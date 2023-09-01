package tango_test

import (
	"tango"
	"testing"
)

func TestNew(t *testing.T) {
	// Test that Should Pass
	_, err := tango.New("token", "accountIdentifier", true, "production")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Test Token Required - Should Fail
	_, err = tango.New("", "accountIdentifier", true, "production")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test AccountIdentifier Required - Should Fail
	_, err = tango.New("token", "", true, "production")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test Environment Default - Should Pass
	_, err = tango.New("token", "accountIdentifier", true, "")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Test Environment Production - Should Pass
	_, err = tango.New("token", "accountIdentifier", true, "production")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Test Environment Erroneous Input
	_, err = tango.New("token", "accountIdentifier", true, "bad")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	t.Logf("Testing for New() complete")
}
