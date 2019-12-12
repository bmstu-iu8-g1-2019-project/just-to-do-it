package models

type Group struct {
	Id          int    `json:"id"`
	CreatorId   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
