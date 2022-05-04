package types

type User struct {
	Name     string `json:"name"`
	LoggedIn bool   `json:"loggedIn"`
}
