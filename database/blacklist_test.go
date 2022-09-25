package database

import (
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"testing"
)

func Test_BlacklistAuthToken(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name      string            // The name of the test
		want      *models.Blacklist // What out instance we want our function to return.
		wantErr   bool              // whether we want an error.
		authToken string            // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.Blacklist{
				AuthToken: "123445608654321",
			},
			false,
			"123445608654321",
		},
		{
			"no token",
			&models.Blacklist{
				AuthToken: "",
			},
			true,
			"",
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := initTestBlacklistService()
			//fmt.Println("\n\nPRE CREATE: ", tt.group)
			err := testService.BlacklistAuthToken(tt.authToken)
			//fmt.Println("\nPOST CREATE: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupService.BlacklistAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_CheckTokenBlacklist(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name      string // The name of the test
		want      bool   // What out instance we want our function to return.
		authToken string // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"token found",
			true,
			"123445608654321",
		},
		{
			"token not found",
			false,
			"123445608654300",
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestBlacklists()
			found := testService.CheckTokenBlacklist(tt.authToken)
			// Checking the error
			if found != tt.want { // Asserting whether we get the correct wanted value
				t.Errorf("GroupService.CheckTokenBlacklist() = %v, want %v", found, tt.want)
			}
		})
	}
}
