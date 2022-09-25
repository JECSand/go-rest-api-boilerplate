package database

// BlacklistService is used by the app to manage all group related controllers and functionality
type BlacklistService struct {
	collection DBCollection
	db         DBClient
	handler    *DBHandler[*blacklistModel]
}

// NewBlacklistService is an exported function used to initialize a new GroupService struct
func NewBlacklistService(db DBClient, handler *DBHandler[*blacklistModel]) *BlacklistService {
	collection := db.GetCollection("blacklists")
	return &BlacklistService{collection, db, handler}
}

// BlacklistAuthToken is used during sign-out to add the now invalid auth-token/api key to the blacklist collection
func (a *BlacklistService) BlacklistAuthToken(authToken string) error {
	_, err := a.handler.InsertOne(&blacklistModel{AuthToken: authToken})
	if err != nil {
		return err
	}
	return nil
}

// CheckTokenBlacklist to determine if the submitted Auth-Token or API-Key with what's in the blacklist collection
func (a *BlacklistService) CheckTokenBlacklist(authToken string) bool {
	_, err := a.handler.FindOne(&blacklistModel{AuthToken: authToken})
	if err != nil {
		return false
	}
	return true
}
