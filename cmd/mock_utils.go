package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Execute test an http request
func executeRequest(ta App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ta.server.Router.ServeHTTP(rr, req)
	return rr
}

// Check response code returned from a test http request
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// signIn
func signIn(ta App, email string, password string) *httptest.ResponseRecorder {
	payload := []byte(`{"email":"` + email + `","password":"` + password + `"}`)
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	response := executeRequest(ta, req)
	return response
}

// CreateTestGroup creates a group doc for test setup
func createTestGroup(ta App, groupType int) *models.Group {
	group := models.Group{}
	if groupType == 1 {
		group.Id = "000000000000000000000002"
		group.Name = "test2"
		group.RootAdmin = false
		group.LastModified = time.Now().UTC()
		group.CreatedAt = time.Now().UTC()
	} else {
		group.Id = "000000000000000000000003"
		group.Name = "test3"
		group.RootAdmin = false
		group.LastModified = time.Now().UTC()
		group.CreatedAt = time.Now().UTC()
	}
	_, err := ta.server.GroupService.GroupDocInsert(&group)
	if err != nil {
		panic(err)
	}
	return &group
}

// createTestUser creates a user doc for test setup
func createTestUser(ta App, userType int) *models.User {
	user := models.User{}
	if userType == 1 {
		user.Id = "000000000000000000000012"
		user.Username = "test_user"
		user.Password = "abc123"
		user.FirstName = "Jill"
		user.LastName = "Tester"
		user.Email = "test2@email.com"
		user.Role = "member"
		user.RootAdmin = false
		user.GroupId = "000000000000000000000002"
		user.LastModified = time.Now().UTC()
		user.CreatedAt = time.Now().UTC()
	} else {
		user.Id = "000000000000000000000013"
		user.Username = "test_user2"
		user.Password = "abc123"
		user.FirstName = "Bill"
		user.LastName = "Quality"
		user.Email = "test3@email.com.com"
		user.Role = "member"
		user.RootAdmin = false
		user.GroupId = "000000000000000000000003"
		user.LastModified = time.Now().UTC()
		user.CreatedAt = time.Now().UTC()
	}
	_, err := ta.server.UserService.UserDocInsert(&user)
	if err != nil {
		panic(err)
	}
	return &user
}

// createTestTask creates a task doc for test setup
func createTestTask(ta App, taskType int) *models.Task {
	task := models.Task{}
	now := time.Now()
	if taskType == 1 {
		task.Id = "000000000000000000000021"
		task.Name = "testTask"
		task.Status = models.NOTSTARTED
		task.Due = now.Add(time.Hour * 24).UTC()
		task.Description = "Updated Task to complete"
		task.UserId = "000000000000000000000012"
		task.GroupId = "000000000000000000000002"
		task.LastModified = now.UTC()
		task.CreatedAt = now.UTC()
	} else {
		task.Id = "000000000000000000000022"
		task.Name = "testTask2"
		task.Status = models.NOTSTARTED
		task.Due = now.Add(time.Hour * 48).UTC()
		task.Description = "Updated Task to complete2"
		task.UserId = "000000000000000000000013"
		task.GroupId = "000000000000000000000002"
		task.LastModified = now.UTC()
		task.CreatedAt = now.UTC()
	}
	_, err := ta.server.TaskService.TaskDocInsert(&task)
	if err != nil {
		panic(err)
	}
	return &task
}

// getTestUserPayload
func getTestUserPayload(tCase string) []byte {
	switch tCase {
	case "CREATE":
		return []byte(`{"username":"test_user","password":"abc123","firstname":"test","lastname":"user","email":"test2@email.com","group_id":"000000000000000000000002","role":"member"}`)
	case "UPDATE":
		return []byte(`{"username":"newUserName","password":"newUserPass","email":"new_test@email.com","group_id":"000000000000000000000003","role":"member"}`)
	}
	return nil
}

// getTestPasswordPayload
func getTestPasswordPayload(tCase string) []byte {
	switch tCase {
	case "UPDATE_PASSWORD_ERROR":
		return []byte(`{"current_password":"789test122","new_password":"789test124"}`)
	case "UPDATE_PASSWORD_SUCCESS":
		return []byte(`{"current_password":"abc123","new_password":"789test124"}`)
	}
	return nil
}

// getTestGroupPayload
func getTestGroupPayload(tCase string) []byte {
	switch tCase {
	case "CREATE":
		return []byte(`{"name":"testingGroup"}`)
	case "UPDATE":
		return []byte(`{"name":"newTestingGroup"}`)
	}
	return nil
}

// getTestTaskPayload
func getTestTaskPayload(tCase string) []byte {
	var tTask models.Task
	now := time.Now()
	switch tCase {
	case "CREATE":
		tTask.Name = "testTask"
		tTask.Status = models.NOTSTARTED
		tTask.Due = now.Add(time.Hour * 24).UTC()
		tTask.Description = "Updated Task to complete"
		tTask.UserId = "000000000000000000000012"
		tTask.GroupId = "000000000000000000000002"
		b, _ := json.Marshal(tTask)
		return b
	case "UPDATE":
		tTask.Name = "NewTestTask"
		tTask.Status = models.COMPLETED
		tTask.Description = "Updated Task to complete"
		tTask.UserId = "000000000000000000000012"
		tTask.GroupId = "000000000000000000000002"
		b, _ := json.Marshal(tTask)
		return b
	}
	return nil
}
