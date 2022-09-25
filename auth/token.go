package auth

import (
	"errors"
	"fmt"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"github.com/dgrijalva/jwt-go"
	"os"
)

// TokenData stores the structured data from a session token for use
type TokenData struct {
	UserId    string
	Role      string
	RootAdmin bool
	GroupId   string
}

// InitUserToken inputs a pointer to a user and returns TokenData
func InitUserToken(u *models.User) (*TokenData, error) {
	err := u.Validate("auth")
	if err != nil {
		return &TokenData{}, err
	}
	return &TokenData{
		UserId:    u.Id,
		Role:      u.Role,
		RootAdmin: u.RootAdmin,
		GroupId:   u.GroupId,
	}, nil
}

// ToUser creates a new User struct using the TokenData and returns a pointer to it
func (t *TokenData) ToUser() *models.User {
	return &models.User{
		Id:        t.UserId,
		Role:      t.Role,
		RootAdmin: t.RootAdmin,
		GroupId:   t.GroupId,
	}
}

// AdminRouteRoleCheck checks admin routes JWT tokens to ensure that a group admin does not break scope
func (t *TokenData) AdminRouteRoleCheck() string {
	groupId := ""
	if t.RootAdmin {
		groupId = t.GroupId
	}
	return groupId
}

// CreateToken is used to create a new session JWT token
func (t *TokenData) CreateToken(exp int64) (string, error) {
	if t.UserId == "" || t.GroupId == "" || t.Role == "" {
		return "", errors.New("missing required token claims")
	}
	if exp == 0 {
		return "", errors.New("new token must have a expiration time greater than 0")
	}
	var MySigningKey = []byte(os.Getenv("TOKEN_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = t.UserId
	claims["role"] = t.Role
	claims["root"] = t.RootAdmin
	claims["group_id"] = t.GroupId
	claims["exp"] = exp
	return token.SignedString(MySigningKey)
}

// DecodeJWT is used to decode a JWT token
func DecodeJWT(curToken string) (*TokenData, error) {
	var tokenData TokenData
	if curToken == "" {
		return &tokenData, errors.New("unauthorized")
	}
	var MySigningKey = []byte(os.Getenv("TOKEN_SECRET"))
	// Decode token
	token, err := jwt.Parse(curToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(MySigningKey), nil
	})
	if err != nil {
		return &tokenData, err
	}
	// Determine user based on token
	if token.Valid {
		tokenClaims := token.Claims.(jwt.MapClaims)
		tokenData.UserId = tokenClaims["id"].(string)
		tokenData.Role = tokenClaims["role"].(string)
		tokenData.RootAdmin = tokenClaims["root"].(bool)
		tokenData.GroupId = tokenClaims["group_id"].(string)
		return &tokenData, nil
	}
	return &tokenData, errors.New("invalid token")
}
