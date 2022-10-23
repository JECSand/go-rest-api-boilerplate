package database

import (
	"fmt"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"testing"
	"time"
)

func Test_TaskCreate(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.Task // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		task    *models.Task // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.Task{Id: "000000000000000000000022", Name: "Task1", Status: models.NOTSTARTED},
			false,
			&models.Task{
				Id:      "000000000000000000000022",
				Name:    "Task1",
				Due:     time.Now().UTC(),
				UserId:  "000000000000000000000012",
				GroupId: "000000000000000000000002",
			},
		},
		{
			"missing name",
			&models.Task{Id: "000000000000000000000022"},
			true,
			&models.Task{
				Id:      "000000000000000000000001",
				Due:     time.Now().UTC(),
				UserId:  "000000000000000000000012",
				GroupId: "000000000000000000000002",
			},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := initTestTaskService()
			fmt.Println("\n\nPRE CREATE: ", tt.task)
			got, err := testService.TaskCreate(tt.task)
			fmt.Println("\nPOST CREATE: ", got)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TaskCreate() error = %v, wantErr %v", err, tt.wantErr)
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
				if got.Id != tt.want.Id || got.CreatedAt.IsZero() || got.Status != tt.want.Status { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskCreate() = %v, want %v", got, tt.want)
				}
			}
			if failMsg != "" { // Asserting whether we get the correct wanted value
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_TasksFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string // The name of the test
		want    int    // What out instance we want our function to return.
		wantErr bool   // whether we want an error.
		task    *models.Task
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"user tasks success",
			1,
			false,
			&models.Task{
				Due:    time.Now().UTC(),
				UserId: "000000000000000000000012",
			},
		},
		{
			"group tasks success",
			2,
			false,
			&models.Task{
				Due:     time.Now().UTC(),
				GroupId: "000000000000000000000002",
			},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestTasks()
			got, err := testService.TasksFind(tt.task)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TasksFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want { // Asserting whether we get the correct wanted value
				t.Errorf("TaskService.TasksFind() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_TaskFind(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.Task // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		task    *models.Task
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"find by id",
			&models.Task{Id: "000000000000000000000022", Name: "Task1"},
			false,
			&models.Task{Id: "000000000000000000000022"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestTasks()
			got, err := testService.TaskFind(tt.task)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TaskFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id || got.Name != tt.want.Name { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskFind() = %v, want %v", got.Id, tt.want.Id)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}

		})
	}
}

func Test_TaskUpdate(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.Task // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		task    *models.Task
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"complete task",
			&models.Task{Id: "000000000000000000000022", Name: "Task1", Status: models.COMPLETED},
			false,
			&models.Task{Id: "000000000000000000000022", Status: models.COMPLETED},
		},
		{
			"in progress task",
			&models.Task{Id: "000000000000000000000022", Name: "Task1", Status: models.INPROGRESS},
			false,
			&models.Task{Id: "000000000000000000000022", Status: models.INPROGRESS},
		},
		{
			"invalid user id",
			&models.Task{Id: "000000000000000000000022", Name: "Task1"},
			true,
			&models.Task{Id: "000000000000000000000022", UserId: "000000000000000000000002", Status: models.COMPLETED},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestTasks()
			got, err := testService.TaskUpdate(tt.task)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TaskUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "complete task":
				if got.Status != models.COMPLETED || got.Name != tt.want.Name { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskUpdate() = %v, want %v", got.Name, tt.want.Name)
				}
			case "in progress task":
				if got.Status != models.INPROGRESS || got.Name != tt.want.Name { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskUpdate() = %v, want %v", got.Name, tt.want.Name)
				}
			}

			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}

func Test_TaskDelete(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string       // The name of the test
		want    *models.Task // What out instance we want our function to return.
		wantErr bool         // whether we want an error.
		task    *models.Task
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&models.Task{Id: "000000000000000000000022", Name: "Task1"},
			false,
			&models.Task{Id: "000000000000000000000022"},
		},
		{
			"task not found",
			nil,
			true,
			&models.Task{Id: "000000000000000000000025"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testService := setupTestTasks()
			got, err := testService.TaskDelete(tt.task)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskService.TaskDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var failMsg string
			switch tt.name {
			case "success":
				if got.Id != tt.want.Id { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskDelete() = %v, want %v", got.Id, tt.want.Id)
				}
			case "task not found":
				if got != tt.want { // Asserting whether we get the correct wanted value
					failMsg = fmt.Sprintf("TaskService.TaskDelete() = %v, want %v", got, tt.want)
				}
			}
			if failMsg != "" {
				t.Errorf(failMsg)
			}
		})
	}
}
