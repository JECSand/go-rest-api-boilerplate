package database

import (
	"fmt"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"reflect"
	"testing"
)

func Test_GroupCreate(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string        // The name of the test
		want    *models.Group // What out instance we want our function to return.
		wantErr bool          // whether we want an error.
		group   *models.Group // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.Group{Id: "000000000000000000000001", Name: "test", RootAdmin: false},
			false,
			&models.Group{Id: "000000000000000000000001", Name: "test", RootAdmin: false},
		},
		{
			"missing name",
			nil,
			true,
			&models.Group{Id: "000000000000000000000002", RootAdmin: false},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := initTestGroupService()
			//fmt.Println("\n\nPRE CREATE: ", tt.group)
			got, err := testService.GroupCreate(tt.group)
			//fmt.Println("\nPOST CREATE: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupService.GroupCreate() error = %v, wantErr %v", err, tt.wantErr)
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
				if !reflect.DeepEqual(got, tt.want) { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("GroupService.GroupCreate() = %v, want %v", got, tt.want)
				}
			case "missing name":
				if got != tt.want { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("GroupService.GroupUpdate() = %v, want %v", got.Name, tt.want.Name)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_GroupsFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string // The name of the test
		want    int    // What out instance we want our function to return.
		wantErr bool   // whether we want an error.
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			2,
			false,
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestGroups()
			got, err := testService.GroupsFind()
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupService.GroupsFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want { // Asserting whether we get the correct wanted value
				t.Errorf("GroupService.GroupsFind() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_GroupFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string        // The name of the test
		want    *models.Group // What out instance we want our function to return.
		wantErr bool          // whether we want an error.
		group   *models.Group
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"find by id",
			&models.Group{Id: "000000000000000000000002", Name: "test2", RootAdmin: false},
			false,
			&models.Group{Id: "000000000000000000000002"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestGroups()
			got, err := testService.GroupFind(tt.group)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupService.GroupFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "find by id":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("GroupService.GroupFind() = %v, want %v", got.Id, tt.want.Id)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_GroupUpdate(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string        // The name of the test
		want    *models.Group // What out instance we want our function to return.
		wantErr bool          // whether we want an error.
		group   *models.Group
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"update group name",
			&models.Group{Id: "000000000000000000000002", Name: "test4"},
			false,
			&models.Group{Id: "000000000000000000000002", Name: "test4"},
		},
		{
			"name taken",
			nil,
			true,
			&models.Group{Id: "000000000000000000000002", Name: "test3"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestGroups()
			got, err := testService.GroupUpdate(tt.group)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupService.GroupUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Name != tt.want.Name { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("GroupService.GroupUpdate() = %v, want %v", got.Name, tt.want.Name)
				}
			case "name taken":
				if got != tt.want { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("GroupService.GroupUpdate() = %v, want %v", got.Name, tt.want.Name)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_GroupDelete(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string        // The name of the test
		want    *models.Group // What out instance we want our function to return.
		wantErr bool          // whether we want an error.
		group   *models.Group
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.Group{Id: "000000000000000000000002", Name: "test4"},
			false,
			&models.Group{Id: "000000000000000000000002"},
		},
		{
			"group not found",
			&models.Group{Id: "000000000000000000000004", Name: "test3"},
			true,
			&models.Group{Id: "000000000000000000000004"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestGroups()
			got, err := testService.GroupDelete(tt.group)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupService.GroupDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("GroupService.GroupDelete() = %v, want %v", got.Id, tt.want.Id)
				}
			case "name taken":
				if got != tt.want { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("GroupService.GroupUpdate() = %v, want %v", got.Name, tt.want.Name)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}
