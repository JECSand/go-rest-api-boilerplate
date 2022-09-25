package models

import (
	"testing"
	"time"
)

func Test_ValidateUser(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string // The name of the test
		wantErr bool   // whether we want an error.
		user    *User  // The input of the test
		valCase string // What out instance we want our function to return.
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"auth success",
			false,
			&User{
				Id:      "000000000000000000000011",
				GroupId: "000000000000000000000001",
				Role:    "member",
			},
			"auth",
		},
		{
			"auth id error",
			true,
			&User{
				Id:      "",
				GroupId: "000000000000000000000001",
				Role:    "member",
			},
			"auth",
		},
		{
			"create success",
			false,
			&User{
				Username: "testUser",
				Email:    "test3@example.com",
				Password: "abc123",
				GroupId:  "000000000000000000000001",
			},
			"create",
		},
		{
			"create error",
			true,
			&User{
				Username: "testUser",
				Email:    "test3@example.com",
				Password: "abc123",
				GroupId:  "000000000000000000000000",
			},
			"create",
		},
		{
			"update success",
			false,
			&User{
				Id: "000000000000000000000001",
			},
			"update",
		},
		{
			"update error",
			true,
			&User{
				Id: "000000000000000000000000",
			},
			"update",
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.Validate(tt.valCase)
			// Checking the error
			if (got != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", got, tt.wantErr)
				return
			}
		})
	}
}

func Test_ValidateGroup(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string // The name of the test
		wantErr bool   // whether we want an error.
		group   *Group // The input of the test
		valCase string // What out instance we want our function to return.
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"create success",
			false,
			&Group{
				Name: "testUser",
			},
			"create",
		},
		{
			"create error",
			true,
			&Group{
				Name: "",
			},
			"create",
		},
		{
			"update success",
			false,
			&Group{
				Id: "000000000000000000000001",
			},
			"update",
		},
		{
			"update error",
			true,
			&Group{
				Id: "000000000000000000000000",
			},
			"update",
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.group.Validate(tt.valCase)
			// Checking the error
			if (got != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", got, tt.wantErr)
				return
			}
		})
	}
}

func Test_ValidateTask(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string // The name of the test
		wantErr bool   // whether we want an error.
		task    *Task  // The input of the test
		valCase string // What out instance we want our function to return.
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"create success",
			false,
			&Task{
				Name:    "testUser",
				UserId:  "00000000000000000000011",
				GroupId: "00000000000000000000001",
				Due:     time.Now().UTC(),
			},
			"create",
		},
		{
			"create error",
			true,
			&Task{
				Name:    "testUser",
				UserId:  "000000000000000000000000",
				GroupId: "000000000000000000000001",
				Due:     time.Now().UTC(),
			},
			"create",
		},
		{
			"update success",
			false,
			&Task{
				Id: "000000000000000000000001",
			},
			"update",
		},
		{
			"update error",
			true,
			&Task{
				Id: "000000000000000000000000",
			},
			"update",
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task.Validate(tt.valCase)
			// Checking the error
			if (got != nil) != tt.wantErr {
				t.Errorf("Task.Validate() error = %v, wantErr %v", got, tt.wantErr)
				return
			}
		})
	}
}
