package database

import (
	"context"
	"errors"
	"github.com/JECSand/go-rest-api-boilerplate/models"
	"github.com/JECSand/go-rest-api-boilerplate/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

// UserService is used by the app to manage all user related controllers and functionality
type UserService struct {
	collection   DBCollection
	db           DBClient
	userHandler  *DBHandler[*userModel]
	groupHandler *DBHandler[*groupModel]
}

// NewUserService is an exported function used to initialize a new UserService struct
func NewUserService(db DBClient, uHandler *DBHandler[*userModel], gHandler *DBHandler[*groupModel]) *UserService {
	collection := db.GetCollection("users")
	return &UserService{collection, db, uHandler, gHandler}
}

// checkLinkedRecords ensures the email is unique and groupId valid for a User
func (p *UserService) checkLinkedRecords(g *groupModel, u *userModel, curUser *userModel) error {
	var wg sync.WaitGroup
	uCh := make(chan *userModel)
	uErr := make(chan error)
	gCh := make(chan *groupModel)
	gErr := make(chan error)
	uRoutine := p.userHandler.newRoutine()
	gRoutine := p.groupHandler.newRoutine()
	wg.Add(2)
	go uRoutine.execute(FindOne, uCh, uErr, u, nil)
	go gRoutine.execute(FindOne, gCh, gErr, g, nil)
	go uRoutine.resolve(uCh, uErr, &wg)
	go gRoutine.resolve(gCh, gErr, &wg)
	wg.Wait()
	close(uCh)
	close(uErr)
	close(gCh)
	close(gErr)
	if curUser == nil && uRoutine.err == nil {
		return errors.New("email is taken")
	} else if uRoutine.err == nil && curUser.Email != u.Email {
		return errors.New("email is taken")
	} else if gRoutine.err != nil {
		return errors.New("invalid group id")
	}
	return nil
}

// AuthenticateUser is used to authenticate users that are signing in
func (p *UserService) AuthenticateUser(u *models.User) (*models.User, error) {
	um, err := newUserModel(u)
	if err != nil {
		return nil, err
	}
	checkUser, err := p.userHandler.FindOne(um)
	if err != nil {
		return nil, errors.New("invalid email")
	}
	rootUser := checkUser.toRoot()
	err = rootUser.Authenticate(u.Password)
	if err == nil {
		return rootUser, nil
	}
	return nil, errors.New("invalid password")
}

// UserCreate is used to create a new user
func (p *UserService) UserCreate(u *models.User) (*models.User, error) {
	if u.Id == "" {
		u.Id = utilities.GenerateObjectID()
	}
	um, err := newUserModel(u)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := p.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = p.checkLinkedRecords(&groupModel{Id: um.GroupId}, &userModel{Email: um.Email}, nil)
	if err != nil {
		return nil, err
	}
	err = u.HashPassword()
	if err != nil {
		return nil, err
	}
	u.RootAdmin = false
	if docCount == 0 {
		u.Role = "admin"
		u.RootAdmin = true
	} else if u.Role != "admin" {
		u.Role = "member"
	}
	um, err = newUserModel(u)
	if err != nil {
		return nil, err
	}
	um, err = p.userHandler.InsertOne(um)
	if err != nil {
		return nil, err
	}
	return um.toRoot(), err
}

// UserDelete is used to delete an User
func (p *UserService) UserDelete(u *models.User) (*models.User, error) {
	um, err := newUserModel(u)
	if err != nil {
		return nil, err
	}
	um, err = p.userHandler.DeleteOne(um)
	if err != nil {
		return nil, err
	}
	return um.toRoot(), err
}

// UserDeleteMany is used to delete many Users
func (p *UserService) UserDeleteMany(u *models.User) (*models.User, error) {
	um, err := newUserModel(u)
	if err != nil {
		return nil, err
	}
	um, err = p.userHandler.DeleteMany(um)
	if err != nil {
		return nil, err
	}
	return um.toRoot(), err
}

// UsersFind is used to find all user docs
func (p *UserService) UsersFind(u *models.User) ([]*models.User, error) {
	var users []*models.User
	um, err := newUserModel(u)
	if err != nil {
		return users, err
	}
	ums, err := p.userHandler.FindMany(um)
	if err != nil {
		return users, err
	}
	for _, m := range ums {
		users = append(users, m.toRoot())
	}
	return users, nil
}

// UserFind is used to find a specific user doc
func (p *UserService) UserFind(u *models.User) (*models.User, error) {
	um, err := newUserModel(u)
	if err != nil {
		return nil, err
	}
	um, err = p.userHandler.FindOne(um)
	if err != nil {
		return nil, err
	}
	return um.toRoot(), err
}

// UserUpdate is used to update an existing user doc
func (p *UserService) UserUpdate(u *models.User) (*models.User, error) {
	filter, err := u.BuildFilter()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := p.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if docCount == 0 {
		return nil, errors.New("no users found")
	}
	f, err := newUserModel(filter)
	if err != nil {
		return nil, err
	}
	curUser, err := p.userHandler.FindOne(f)
	if err != nil {
		return u, err
	}
	u.BuildUpdate(curUser.toRoot())
	um, err := newUserModel(u)
	if err != nil {
		return nil, err
	}
	err = p.checkLinkedRecords(&groupModel{Id: um.GroupId}, &userModel{Email: um.Email}, curUser)
	if err != nil {
		return nil, err
	}
	err = u.HashPassword()
	if u.Password != "" {
		if err != nil {
			return nil, err
		}
		um.Password = u.Password
	}
	um, err = p.userHandler.UpdateOne(f, um)
	if err != nil {
		return nil, err
	}
	return um.toRoot(), err
}

// UpdatePassword is used to update the currently logged-in user's password
func (p *UserService) UpdatePassword(u *models.User, currentPassword string, newPassword string) (*models.User, error) {
	um, err := newUserModel(u)
	if err != nil {
		return nil, err
	}
	user, err := p.userHandler.FindOne(um)
	if err != nil {
		return nil, err
	}
	rootUser := user.toRoot()
	err = rootUser.Authenticate(currentPassword)
	if err == nil { // 3. Update doc with new password
		currentTime := time.Now().UTC()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		filter := bson.D{{"_id", user.Id}}
		update := bson.D{{"$set",
			bson.D{
				{"password", string(hashedPassword)},
				{"last_modified", currentTime},
			},
		}}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, err = p.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
		user.Password = ""
		return user.toRoot(), nil
	}
	return nil, errors.New("invalid password")
}

// UserDocInsert is used to insert user doc directly into mongodb for testing purposes
func (p *UserService) UserDocInsert(u *models.User) (*models.User, error) {
	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return u, err
	}
	u.Password = string(hashedPassword)
	insertUser, err := newUserModel(u)
	if err != nil {
		return u, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = p.collection.InsertOne(ctx, insertUser)
	if err != nil {
		return u, err
	}
	return insertUser.toRoot(), nil
}
