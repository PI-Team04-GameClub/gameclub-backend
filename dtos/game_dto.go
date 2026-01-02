package dtos

type CreateGameRequest struct {
	Name            string  `json:"name" validate:"required"`
	Description     string  `json:"description"`
	NumberOfPlayers int     `json:"numberOfPlayers" validate:"required,min=1"`
	MinPlayers      int     `json:"minPlayers" validate:"min=1"`
	MaxPlayers      int     `json:"maxPlayers" validate:"min=1"`
	PlaytimeMinutes int     `json:"playtimeMinutes" validate:"min=1"`
	MinAge          int     `json:"minAge" validate:"min=3"`
	Complexity      string  `json:"complexity"`
	Category        string  `json:"category"`
	Publisher       string  `json:"publisher"`
	YearPublished   int     `json:"yearPublished" validate:"min=1900,max=2100"`
	Rating          float64 `json:"rating" validate:"min=0,max=10"`
}

type GameResponse struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	NumberOfPlayers int     `json:"numberOfPlayers"`
	MinPlayers      int     `json:"minPlayers"`
	MaxPlayers      int     `json:"maxPlayers"`
	PlaytimeMinutes int     `json:"playtimeMinutes"`
	MinAge          int     `json:"minAge"`
	Complexity      string  `json:"complexity"`
	Category        string  `json:"category"`
	Publisher       string  `json:"publisher"`
	YearPublished   int     `json:"yearPublished"`
	Rating          float64 `json:"rating"`
}
