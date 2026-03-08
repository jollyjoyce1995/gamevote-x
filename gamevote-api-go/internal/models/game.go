package models

type Game struct {
	ID       string `json:"id,omitempty" surreal:"id,omitempty"`
	AppID    int    `json:"appId" surreal:"appId"`
	Name     string `json:"name" surreal:"name"`
	ImageURL string `json:"imageUrl" surreal:"imageUrl"`
}
