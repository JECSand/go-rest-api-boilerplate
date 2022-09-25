package cmd

import (
	"bytes"
	"net/http"
	"os"
	"testing"
)

var ta App

// Setup Tests
func setup() {
	os.Setenv("ENV", "test")
	ta = App{}
	err := ta.Initialize()
	if err != nil {
		panic(err)
	}
}

/*
AUTH & USER TESTS
*/

// User SignIn Test
func TestSignIn(t *testing.T) {
	setup()
	testResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// Create User Test
func TestCreateUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Create User Test Request
	payload := getTestUserPayload("CREATE")
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestCreateUser() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean test database and check test response
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// User SignOut Test
func TestSignOut(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// SignOut Test
	reqSignOut, err := http.NewRequest("DELETE", "/auth", nil)
	if err != nil {
		t.Errorf("TestSignOut() error = %v", err)
	}
	reqSignOut.Header.Add("Content-Type", "application/json")
	reqSignOut.Header.Add("Auth-Token", authToken)
	signOutTestResponse := executeRequest(ta, reqSignOut)
	checkResponseCode(t, http.StatusOK, signOutTestResponse.Code)
	// Test to ensure the Auth-Token or API-Key is now invalid
	payload := getTestUserPayload("CREATE")
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestSignOut() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean test database and check test response
	checkResponseCode(t, http.StatusUnauthorized, testResponse.Code)
}

// Update Password Test
func TestUpdatePassword(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	user := createTestUser(ta, 1)
	authResponse := signIn(ta, user.Email, "abc123")
	checkResponseCode(t, http.StatusOK, authResponse.Code)
	authToken := authResponse.Header().Get("Auth-Token")
	// Update User Password Test Request with incorrect current password
	payloadErr := getTestPasswordPayload("UPDATE_PASSWORD_ERROR")
	reqErr, err := http.NewRequest("POST", "/auth/password", bytes.NewBuffer(payloadErr))
	if err != nil {
		t.Errorf("TestUpdatePassword() error = %v", err)
	}
	reqErr.Header.Add("Content-Type", "application/json")
	reqErr.Header.Add("Auth-Token", authToken)
	testResponseErr := executeRequest(ta, reqErr)
	checkResponseCode(t, http.StatusUnauthorized, testResponseErr.Code)
	// Test to ensure password did not update due to incorrect current password
	testResponseErrAuth := signIn(ta, user.Email, "789test124")
	checkResponseCode(t, http.StatusUnauthorized, testResponseErrAuth.Code)
	// Update User Password Test Request correctly
	payloadOK := getTestPasswordPayload("UPDATE_PASSWORD_SUCCESS")
	reqOK, err := http.NewRequest("POST", "/auth/password", bytes.NewBuffer(payloadOK))
	if err != nil {
		t.Errorf("TestUpdatePassword() error = %v", err)
	}
	reqOK.Header.Add("Content-Type", "application/json")
	reqOK.Header.Add("Auth-Token", authToken)
	testResponseOK := executeRequest(ta, reqOK)
	checkResponseCode(t, http.StatusAccepted, testResponseOK.Code)
	// Test to ensure user can now log in with new password
	testResponseOKAuth := signIn(ta, user.Email, "789test124")
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponseOKAuth.Code)
}

// TestModifyUser User Test
func TestModifyUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestGroup(ta, 2)
	createTestUser(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Modify User Document Test
	payloadUpdate := getTestUserPayload("UPDATE")
	reqUpdate, err := http.NewRequest("PATCH", "/users/000000000000000000000012", bytes.NewBuffer(payloadUpdate))
	if err != nil {
		t.Errorf("TestModifyUser() error = %v", err)
	}
	reqUpdate.Header.Add("Content-Type", "application/json")
	reqUpdate.Header.Add("Auth-Token", authToken)
	updateTestResponse := executeRequest(ta, reqUpdate)
	checkResponseCode(t, http.StatusAccepted, updateTestResponse.Code)
	// Attempt to test Modified user doc by loggin in with new password/username
	testResponseOKAuth := signIn(ta, "new_test@email.com", "newUserPass")
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponseOKAuth.Code)
}

// List Users Test
func TestListUsers(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	createTestUser(ta, 2)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// List all users test
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("TestListUsers() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// TestListUser User Test
func TestListUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// List a specific user Test
	req, err := http.NewRequest("GET", "/users/000000000000000000000012", nil)
	if err != nil {
		t.Errorf("TestListUser() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// TestDeleteUser User Test
func TestDeleteUser(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Delete a user Test
	req, err := http.NewRequest("DELETE", "/users/000000000000000000000012", nil)
	if err != nil {
		t.Errorf("TestDeleteUser() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// TestTokenRefresh Auth Token Test
func TestTokenRefresh(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Refresh Auth-Token Test
	reqRefresh, err := http.NewRequest("GET", "/auth", nil)
	if err != nil {
		t.Errorf("TestTokenRefresh() error = %v", err)
	}
	reqRefresh.Header.Add("Content-Type", "application/json")
	reqRefresh.Header.Add("Auth-Token", authToken)
	refreshTestResponse := executeRequest(ta, reqRefresh)
	checkResponseCode(t, http.StatusOK, refreshTestResponse.Code)
	authToken = refreshTestResponse.Header().Get("Auth-Token")
	// Test new token to ensure it works by creating a new group
	payload := getTestUserPayload("CREATE")
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestTokenRefresh() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// TestGenerateAPIKey Test
func TestGenerateAPIKey(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Refresh Auth-Token Test
	reqAPIKey, err := http.NewRequest("GET", "/auth/api-key", nil)
	if err != nil {
		t.Errorf("TestGenerateAPIKey() error = %v", err)
	}
	reqAPIKey.Header.Add("Content-Type", "application/json")
	reqAPIKey.Header.Add("Auth-Token", authToken)
	apiKeyTestResponse := executeRequest(ta, reqAPIKey)
	checkResponseCode(t, http.StatusOK, apiKeyTestResponse.Code)
	apiKey := apiKeyTestResponse.Header().Get("API-Key")
	// Test new token to ensure it works by creating a new group
	payload := getTestUserPayload("CREATE")
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestGenerateAPIKey() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", apiKey)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

/*
GROUP TESTS
*/

// TestCreateGroup Test
func TestCreateGroup(t *testing.T) {
	// Test Setup
	setup()
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Create new group test
	payload := getTestGroupPayload("CREATE")
	req, err := http.NewRequest("POST", "/groups", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestCreateGroup() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// TestModifyGroup Test
func TestModifyGroup(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Modify group test
	payload := getTestGroupPayload("UPDATE")
	req, err := http.NewRequest("PATCH", "/groups/000000000000000000000002", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestModifyGroup() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusAccepted, testResponse.Code)
}

// TestListGroups Test
func TestListGroups(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestGroup(ta, 2)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// List all groups test
	req, err := http.NewRequest("GET", "/groups", nil)
	if err != nil {
		t.Errorf("TestListGroups() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// TestListGroup Test
func TestListGroup(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// List all groups test
	req, err := http.NewRequest("GET", "/groups/000000000000000000000002", nil)
	if err != nil {
		t.Errorf("TestListGroup() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// TestDeleteGroup Test
func TestDeleteGroup(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// List all groups test
	req, err := http.NewRequest("DELETE", "/groups/000000000000000000000002", nil)
	if err != nil {
		t.Errorf("TestDeleteGroup() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

/*
TASKS TESTS
*/

// TestCreateTask Test
func TestCreateTask(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	createTestUser(ta, 1)
	authResponse := signIn(ta, os.Getenv("ROOT_EMAIL"), os.Getenv("ROOT_PASSWORD"))
	authToken := authResponse.Header().Get("Auth-Token")
	// Create new todos test
	payload := getTestTaskPayload("CREATE")
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestCreateTask() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusCreated, testResponse.Code)
}

// TestModifyTask Test
func TestModifyTask(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	user := createTestUser(ta, 1)
	createTestTask(ta, 1)
	authResponse := signIn(ta, user.Email, "abc123")
	checkResponseCode(t, http.StatusOK, authResponse.Code)
	authToken := authResponse.Header().Get("Auth-Token")
	// Modify todos test
	payload := getTestTaskPayload("UPDATE")
	req, err := http.NewRequest("PATCH", "/tasks/000000000000000000000021", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("TestModifyTask() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusAccepted, testResponse.Code)
}

// TestListTasks Test
func TestListTasks(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	user := createTestUser(ta, 1)
	createTestTask(ta, 1)
	createTestTask(ta, 2)
	authResponse := signIn(ta, user.Email, "abc123")
	checkResponseCode(t, http.StatusOK, authResponse.Code)
	authToken := authResponse.Header().Get("Auth-Token")
	// List all todos test
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Errorf("TestListTasks() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// TestListTask Test
func TestListTask(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	user := createTestUser(ta, 1)
	createTestTask(ta, 1)
	authResponse := signIn(ta, user.Email, "abc123")
	checkResponseCode(t, http.StatusOK, authResponse.Code)
	authToken := authResponse.Header().Get("Auth-Token")
	// List a specific todos doc
	req, err := http.NewRequest("GET", "/tasks/000000000000000000000021", nil)
	if err != nil {
		t.Errorf("TestListTask() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}

// Delete Todos Test
func TestDeleteTask(t *testing.T) {
	// Test Setup
	setup()
	createTestGroup(ta, 1)
	user := createTestUser(ta, 1)
	createTestTask(ta, 1)
	authResponse := signIn(ta, user.Email, "abc123")
	checkResponseCode(t, http.StatusOK, authResponse.Code)
	authToken := authResponse.Header().Get("Auth-Token")
	// List a specific todos doc
	req, err := http.NewRequest("DELETE", "/tasks/000000000000000000000021", nil)
	if err != nil {
		t.Errorf("TestDeleteTask() error = %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", authToken)
	testResponse := executeRequest(ta, req)
	// Clean database and do final status check
	checkResponseCode(t, http.StatusOK, testResponse.Code)
}
