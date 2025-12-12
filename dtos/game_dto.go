package dtos

type CreateGameRequest struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description"`
	NumberOfPlayers int    `json:"numberOfPlayers" validate:"required,min=1"`
}

type GameResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	NumberOfPlayers int    `json:"numberOfPlayers"`
}
