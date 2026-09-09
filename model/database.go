package model

type UserTable struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Nickname     string `json:"nickname"`
	Admin        bool   `json:"admin"`
	Enable       bool   `json:"enable"`
	TokenVersion string `json:"token_version"`
}

type PortTable struct {
	ID    int    `json:"id"`
	Min   int    `json:"min"`
	Max   int    `json:"max"`
	Token string `json:"token"`
	User  string `json:"user"`
}

type TokenTable struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	User   string `json:"user"`
	Enable int    `json:"enable"`
}
