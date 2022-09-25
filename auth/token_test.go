package auth

import (
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"reflect"
	"testing"
	"time"
)

func Test_initUserToken(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *TokenData   // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&TokenData{UserId: "000000000000000000000001", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
			false,
			&models.User{Id: "000000000000000000000001", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
		},
		{"invalid user",
			&TokenData{},
			true,
			&models.User{Id: "000000000000000000000000", GroupId: "000000000000000000000012", Role: "member", RootAdmin: false},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitUserToken(tt.user)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("InitUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) { // Asserting whether we get the correct wanted value
				t.Errorf("InitUserToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createToken(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name      string     // The name of the test
		exp       int64      // Token expiration time
		want      string     // What out instance we want our function to return.
		wantErr   bool       // whether we want an error.
		tokenData *TokenData // The command arguments used for this test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			time.Now().Add(time.Hour * 1).Unix(),
			"",
			false,
			&TokenData{UserId: "000000000000000000000001", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
		},
		{
			"missing claims",
			time.Now().Add(time.Hour * 1).Unix(),
			"1",
			true,
			&TokenData{UserId: "", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
		},
		{
			"expiration time of 0",
			0,
			"1",
			true,
			&TokenData{UserId: "", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tokenData.CreateToken(tt.exp)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenData.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want { // Asserting whether we got an unwanted value (opposite as usually)
				t.Errorf("TokenData.CreateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeToken(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name      string     // The name of the test
		exp       int64      // Token expiration time
		want      *TokenData // What out instance we want our function to return.
		wantErr   bool       // whether we want an error.
		tokenData *TokenData // The command arguments used for this test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			time.Now().Add(time.Hour * 1).Unix(),
			&TokenData{UserId: "000000000000000000000001", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
			false,
			&TokenData{UserId: "000000000000000000000001", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
		},
		{
			"expired token",
			time.Now().Add(time.Second * 1).Unix(),
			&TokenData{},
			true,
			&TokenData{UserId: "000000000000000000000001", GroupId: "000000000000000000000011", Role: "member", RootAdmin: false},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testToken, _ := tt.tokenData.CreateToken(tt.exp)
			time.Sleep(3 * time.Second)
			got, err := DecodeJWT(testToken)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) { // Asserting whether we get the correct wanted value
				t.Errorf("DecodeJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
