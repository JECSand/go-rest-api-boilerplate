package database

import (
	"fmt"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"testing"
)

func Test_UserCreate(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.User // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.User{
				Id:        "000000000000000000000012",
				Email:     "test2@email.com",
				Password:  "abc123",
				GroupId:   "000000000000000000000002",
				RootAdmin: false,
			},
			false,
			&models.User{
				Id:        "000000000000000000000012",
				Email:     "test2@email.com",
				Password:  "abc123",
				GroupId:   "000000000000000000000002",
				RootAdmin: false,
			},
		},
		{
			"success no id",
			&models.User{
				Email:     "test2@email.com",
				Password:  "abc123",
				GroupId:   "000000000000000000000002",
				RootAdmin: false,
			},
			false,
			&models.User{
				Email:     "test2@email.com",
				Password:  "abc123",
				GroupId:   "000000000000000000000002",
				RootAdmin: false,
			},
		},
		{
			"missing email",
			&models.User{
				Id:        "00000000000000000000012",
				Password:  "abc123",
				GroupId:   "00000000000000000000002",
				RootAdmin: false,
			},
			true,
			&models.User{
				Id:        "00000000000000000000012",
				Password:  "abc123",
				GroupId:   "00000000000000000000002",
				RootAdmin: false,
			},
		},
		{
			"missing password",
			&models.User{
				Id:        "00000000000000000000012",
				Email:     "test2@email.com",
				GroupId:   "00000000000000000000002",
				RootAdmin: false,
			},
			true,
			&models.User{
				Id:        "00000000000000000000012",
				Email:     "test2@email.com",
				GroupId:   "00000000000000000000002",
				RootAdmin: false,
			},
		},
		{
			"missing group id",
			&models.User{
				Id:        "0000000000000000000012",
				Email:     "test2@email.com",
				Password:  "abc123",
				RootAdmin: false,
			},
			true,
			&models.User{
				Id:        "00000000000000000000012",
				Email:     "test2@email.com",
				Password:  "abc123",
				RootAdmin: false,
			},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := initTestUserService()
			//fmt.Println("\n\nPRE CREATE: ", tt.user)
			got, err := testService.UserCreate(tt.user)
			//fmt.Println("\nPOST CREATE: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UserCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !got.CreatedAt.IsZero() && !got.LastModified.IsZero() {
					tt.want.CreatedAt = got.CreatedAt
					tt.want.LastModified = got.LastModified
				}
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("UserService.UserCreate() = %v, want %v", got.Id, tt.want.Id)
				}
			case "success no id":
				if got.Email != tt.want.Email { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("UserService.UserCreate() = %v, want %v", got.Email, tt.want.Email)
				} else if got.Id == "000000000000000000000000" || got.Id == "" {
					failMsg = fmt.Sprintf("UserService.UserCreate() = %v", got.Id)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_UsersFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    int          // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"find by id",
			1,
			false,
			&models.User{Id: "000000000000000000000012"},
		},
		{
			"find by group id",
			2,
			false,
			&models.User{GroupId: "000000000000000000000002"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestUsers()
			// fmt.Println("\n\nPRE FIND: ", tt.user)
			got, err := testService.UsersFind(tt.user)
			//fmt.Println("\n\nPOST FIND: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UsersFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want { // Asserting whether we get the correct wanted value
				t.Errorf("UserService.UsersFind() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_UserFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.User // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"find by id",
			&models.User{Id: "000000000000000000000012", Email: "test2@email.com", RootAdmin: false},
			false,
			&models.User{Id: "000000000000000000000012"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestUsers()
			fmt.Println("\nPRE FIND: ", tt.user)
			got, err := testService.UserFind(tt.user)
			fmt.Println("\nPOST FIND: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UserFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Email != tt.want.Email { // Asserting whether we get the correct wanted value
				t.Errorf("UserService.UserFind() = %v, want %v", got.Id, tt.want.Id)
			}
		})
	}
}

func Test_UserUpdate(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.User // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"update user group",
			&models.User{
				Id:      "000000000000000000000012",
				Email:   "test2@email.com",
				GroupId: "000000000000000000000003",
			},
			false,
			&models.User{
				Id:      "000000000000000000000012",
				GroupId: "000000000000000000000002",
			},
		},
		{
			"group not found",
			&models.User{Id: "000000000000000000000012", GroupId: "000000000000000000000004"},
			true,
			&models.User{Id: "000000000000000000000012", GroupId: "000000000000000000000004"},
		},
		{
			"update user email",
			&models.User{Id: "000000000000000000000012", Email: "test4@email.com"},
			false,
			&models.User{Id: "000000000000000000000012", Email: "test4@email.com"},
		},
		{
			"email taken",
			&models.User{Id: "000000000000000000000012", Email: "test3@email.com"},
			true,
			&models.User{Id: "000000000000000000000012", Email: "test3@email.com"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestUsers()
			got, err := testService.UserUpdate(tt.user)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UserUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "update user group", "update user email":
				if got.Email != tt.want.Email { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("UserService.UserUpdate() = %v, want %v", got.Email, tt.want.Email)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_UserDelete(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.User // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.User{Id: "000000000000000000000012"},
			false,
			&models.User{Id: "000000000000000000000012"},
		},
		{
			"user not found",
			nil,
			true,
			&models.User{Id: "000000000000000000000014"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestUsers()
			got, err := testService.UserDelete(tt.user)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UserDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("UserService.UserDelete() = %v, want %v", got.Id, tt.want.Id)
				}
			case "user not found":
				if got != tt.want { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("UserService.UserDelete() = %v, want %v", got, tt.want)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}

		})
	}
}

func Test_AuthenticateUser(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.User // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.User{Id: "000000000000000000000013"},
			false,
			&models.User{Email: "test3@email.com", Password: "abc123"},
		},
		{
			"incorrect password",
			&models.User{Id: "000000000000000000000012"},
			true,
			&models.User{Email: "test2@email.com", Password: "abc12"},
		},
		{
			"invalid email",
			&models.User{Id: "000000000000000000000012"},
			true,
			&models.User{Email: "test5@email.com", Password: "abc123"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestUsers()
			//fmt.Println("\nPRE AUTH: ", tt.user)
			got, err := testService.AuthenticateUser(tt.user)
			//fmt.Println("\nPOST AUTH: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("UserService.AuthenticateUser() = %v, want %v", got.Email, tt.want.Email)
				}
			}
			if failMsg != "" { // Asserting whether we get the correct wanted value
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_UpdatePassword(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.User // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		user    *models.User
		CPW     string
		NPW     string
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.User{Id: "000000000000000000000012"},
			false,
			&models.User{
				Id:        "000000000000000000000012",
				GroupId:   "000000000000000000000002",
				RootAdmin: false,
				Role:      "member",
			},
			"abc123",
			"abc321",
		},
		{
			"incorrect password",
			&models.User{Id: "000000000000000000000012"},
			true,
			&models.User{
				Id:        "000000000000000000000012",
				GroupId:   "000000000000000000000002",
				RootAdmin: false,
				Role:      "member",
			},
			"ab123",
			"abc321",
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestUsers()
			//fmt.Println("\nPRE PW UPDATE: ", tt.user)
			got, err := testService.UpdatePassword(tt.user, tt.CPW, tt.NPW)
			//fmt.Println("\nPOST PW UPDATE: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id || got.LastModified == tt.want.CreatedAt { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("UserService.AuthenticateUser() = %v, want %v", got.CreatedAt, tt.want.LastModified)
				}
			}
			if failMsg != "" { // Asserting whether we get the correct wanted value
				t.Errorf(failMsg)
			}
		})
	}
}
