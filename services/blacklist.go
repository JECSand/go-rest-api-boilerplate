package services

// BlacklistService is an interface used to manage the relevant group doc controllers
type BlacklistService interface {
	BlacklistAuthToken(authToken string) error
	CheckTokenBlacklist(authToken string) bool
}
