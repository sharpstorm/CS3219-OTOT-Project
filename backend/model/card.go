package model

type Card struct {
	Id       int    `bun:"card_id" json:"id"`
	UniqueId string `bun:"card_unique_id" json:"uniqueId"`
	Pokemon  string `bun:"card_pokemon" json:"pokemon"`
	ImageUrl string `bun:"card_image" json:"imageUrl"`
}
