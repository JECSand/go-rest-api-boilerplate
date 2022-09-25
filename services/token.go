package services

import (
	"github.com/JECSand/go-rest-api-boilerplate/auth"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"github.com/JECSand/go-rest-api-boilerplate/utilities"
	"net/http"
	"time"
)

// TokenService is used by the app to manage db auth functionality
type TokenService struct {
	uService UserService
	gService GroupService
	bService BlacklistService
}

// NewTokenService is an exported function used to initialize a new authService struct
func NewTokenService(uService UserService, gService GroupService, bService BlacklistService) *TokenService {
	return &TokenService{uService, gService, bService}
}

// verifyTokenUser verifies Token's User
func (a *TokenService) verifyTokenUser(decodedToken *auth.TokenData) (bool, string) {
	tUser := decodedToken.ToUser()
	checkUser, err := a.uService.UserFind(tUser)
	if err != nil {
		return false, err.Error()
	}
	checkGroup, err := a.gService.GroupFind(&models.Group{Id: tUser.GroupId})
	if err != nil {
		return false, err.Error()
	}
	// validate the Group id of the User and the associated User's Group
	if checkUser.GroupId != checkGroup.Id {
		return false, "Incorrect group id"
	}
	return true, "No Error"
}

// tokenVerifyMiddleWare inputs the route handler function along with User roleType to verify User token and permissions
func (a *TokenService) tokenVerifyMiddleWare(roleType string, next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	var errorObject utilities.JWTError
	authToken := r.Header.Get("Auth-Token")
	if a.bService.CheckTokenBlacklist(authToken) {
		errorObject.Message = "Invalid Token"
		utilities.RespondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
	decodedToken, err := auth.DecodeJWT(r.Header.Get("Auth-Token"))
	if err != nil {
		errorObject.Message = err.Error()
		utilities.RespondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
	verified, verifyMsg := a.verifyTokenUser(decodedToken)
	if verified {
		if roleType == "Admin" && decodedToken.Role == "admin" {
			next.ServeHTTP(w, r)
		} else if roleType != "Admin" {
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = "Invalid Token"
			utilities.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	} else {
		errorObject.Message = verifyMsg
		utilities.RespondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
}

// GenerateToken outputs an auth token string for an inputted User
func (a *TokenService) GenerateToken(u *models.User, tType string) (string, error) {
	expDT := time.Now().Add(time.Hour * 1).Unix() // Default 1 hour expiration for session token
	if tType == "api" {
		expDT = time.Now().Add(time.Hour * 4380).Unix() // 6 month expiration for api key
	}
	tData, err := auth.InitUserToken(u)
	if err != nil {
		return "", err
	}
	return tData.CreateToken(expDT)
}

// AdminTokenVerifyMiddleWare is used to verify that the requester is a valid admin
func (a *TokenService) AdminTokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.tokenVerifyMiddleWare("Admin", next, w, r)
		return
	}
}

// MemberTokenVerifyMiddleWare is used to verify that a requester is authenticated
func (a *TokenService) MemberTokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.tokenVerifyMiddleWare("Member", next, w, r)
		return
	}
}

// BlacklistAuthToken is used to blacklist an unexpired token
func (a *TokenService) BlacklistAuthToken(authToken string) error {
	return a.bService.BlacklistAuthToken(authToken)
}
