package models

import "gorm.io/gorm"

type GameComplexity string
type GameCategory string

const (
	ComplexityEasy   GameComplexity = "Easy"
	ComplexityMedium GameComplexity = "Medium"
	ComplexityHard   GameComplexity = "Hard"
	ComplexityExpert GameComplexity = "Expert"

	CategoryStrategy    GameCategory = "Strategy"
	CategoryParty       GameCategory = "Party"
	CategoryFamily      GameCategory = "Family"
	CategoryCard        GameCategory = "Card"
	CategoryDice        GameCategory = "Dice"
	CategoryCooperative GameCategory = "Cooperative"
)

type Game struct {
	gorm.Model
	Name            string         `gorm:"not null;unique"`
	Description     string         `gorm:"type:text"`
	NumberOfPlayers int            `gorm:"not null"`
	MinPlayers      int            `gorm:"not null;default:2"`
	MaxPlayers      int            `gorm:"not null;default:4"`
	PlaytimeMinutes int            `gorm:"not null;default:30"`
	MinAge          int            `gorm:"not null;default:8"`
	Complexity      GameComplexity `gorm:"type:varchar(20);default:'Medium'"`
	Category        GameCategory   `gorm:"type:varchar(30);default:'Strategy'"`
	Publisher       string         `gorm:"type:varchar(100)"`
	YearPublished   int            `gorm:"default:2024"`
	Rating          float64        `gorm:"type:decimal(3,2);default:0.0"`

	Tournaments []Tournament `gorm:"foreignKey:GameID"`
}

type GameBuilder struct {
	game *Game
}

func NewGameBuilder() *GameBuilder {
	return &GameBuilder{
		game: &Game{
			MinPlayers:      2,
			MaxPlayers:      4,
			PlaytimeMinutes: 30,
			MinAge:          8,
			Complexity:      ComplexityMedium,
			Category:        CategoryStrategy,
			YearPublished:   2024,
			Rating:          0.0,
		},
	}
}

func (b *GameBuilder) SetName(name string) *GameBuilder {
	b.game.Name = name
	return b
}

func (b *GameBuilder) SetDescription(description string) *GameBuilder {
	b.game.Description = description
	return b
}

func (b *GameBuilder) SetNumberOfPlayers(number int) *GameBuilder {
	b.game.NumberOfPlayers = number
	return b
}

func (b *GameBuilder) SetMinPlayers(min int) *GameBuilder {
	b.game.MinPlayers = min
	return b
}

func (b *GameBuilder) SetMaxPlayers(max int) *GameBuilder {
	b.game.MaxPlayers = max
	return b
}

func (b *GameBuilder) SetPlayerRange(min, max int) *GameBuilder {
	b.game.MinPlayers = min
	b.game.MaxPlayers = max
	b.game.NumberOfPlayers = max
	return b
}

func (b *GameBuilder) SetPlaytimeMinutes(minutes int) *GameBuilder {
	b.game.PlaytimeMinutes = minutes
	return b
}

func (b *GameBuilder) SetMinAge(age int) *GameBuilder {
	b.game.MinAge = age
	return b
}

func (b *GameBuilder) SetComplexity(complexity GameComplexity) *GameBuilder {
	b.game.Complexity = complexity
	return b
}

func (b *GameBuilder) SetCategory(category GameCategory) *GameBuilder {
	b.game.Category = category
	return b
}

func (b *GameBuilder) SetPublisher(publisher string) *GameBuilder {
	b.game.Publisher = publisher
	return b
}

func (b *GameBuilder) SetYearPublished(year int) *GameBuilder {
	b.game.YearPublished = year
	return b
}

func (b *GameBuilder) SetRating(rating float64) *GameBuilder {
	b.game.Rating = rating
	return b
}

func (b *GameBuilder) SetID(id uint) *GameBuilder {
	b.game.ID = id
	return b
}

func (b *GameBuilder) SetModel(model gorm.Model) *GameBuilder {
	b.game.Model = model
	return b
}

func (b *GameBuilder) Build() *Game {
	if b.game.NumberOfPlayers == 0 {
		b.game.NumberOfPlayers = b.game.MaxPlayers
	}
	return b.game
}

func (b *GameBuilder) Reset() *GameBuilder {
	b.game = &Game{
		MinPlayers:      2,
		MaxPlayers:      4,
		PlaytimeMinutes: 30,
		MinAge:          8,
		Complexity:      ComplexityMedium,
		Category:        CategoryStrategy,
		YearPublished:   2024,
		Rating:          0.0,
	}
	return b
}
