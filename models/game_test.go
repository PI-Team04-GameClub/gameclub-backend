package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameBuilder(t *testing.T) {
	// Given: Nothing (creating a new builder)

	// When: Creating a new game builder
	builder := NewGameBuilder()

	// Then: The builder should be created with default values
	assert.NotNil(t, builder)
}

func TestGameBuilder_SetName(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the name
	result := builder.SetName("Chess")

	// Then: The builder should be returned for chaining and name should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, "Chess", builder.game.Name)
}

func TestGameBuilder_SetDescription(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the description
	result := builder.SetDescription("A strategy game")

	// Then: The description should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, "A strategy game", builder.game.Description)
}

func TestGameBuilder_SetNumberOfPlayers(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the number of players
	result := builder.SetNumberOfPlayers(6)

	// Then: The number of players should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 6, builder.game.NumberOfPlayers)
}

func TestGameBuilder_SetMinPlayers(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the minimum players
	result := builder.SetMinPlayers(3)

	// Then: The min players should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 3, builder.game.MinPlayers)
}

func TestGameBuilder_SetMaxPlayers(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the maximum players
	result := builder.SetMaxPlayers(8)

	// Then: The max players should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 8, builder.game.MaxPlayers)
}

func TestGameBuilder_SetPlayerRange(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the player range
	result := builder.SetPlayerRange(2, 6)

	// Then: Min, max, and number of players should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 2, builder.game.MinPlayers)
	assert.Equal(t, 6, builder.game.MaxPlayers)
	assert.Equal(t, 6, builder.game.NumberOfPlayers)
}

func TestGameBuilder_SetPlaytimeMinutes(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the playtime
	result := builder.SetPlaytimeMinutes(90)

	// Then: The playtime should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 90, builder.game.PlaytimeMinutes)
}

func TestGameBuilder_SetMinAge(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the minimum age
	result := builder.SetMinAge(12)

	// Then: The min age should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 12, builder.game.MinAge)
}

func TestGameBuilder_SetComplexity(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the complexity
	result := builder.SetComplexity(ComplexityHard)

	// Then: The complexity should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, ComplexityHard, builder.game.Complexity)
}

func TestGameBuilder_SetCategory(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the category
	result := builder.SetCategory(CategoryCooperative)

	// Then: The category should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, CategoryCooperative, builder.game.Category)
}

func TestGameBuilder_SetPublisher(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the publisher
	result := builder.SetPublisher("Hasbro")

	// Then: The publisher should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, "Hasbro", builder.game.Publisher)
}

func TestGameBuilder_SetYearPublished(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the year published
	result := builder.SetYearPublished(2023)

	// Then: The year should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 2023, builder.game.YearPublished)
}

func TestGameBuilder_SetRating(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the rating
	result := builder.SetRating(9.5)

	// Then: The rating should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, 9.5, builder.game.Rating)
}

func TestGameBuilder_SetID(t *testing.T) {
	// Given: A game builder
	builder := NewGameBuilder()

	// When: Setting the ID
	result := builder.SetID(42)

	// Then: The ID should be set
	assert.Equal(t, builder, result)
	assert.Equal(t, uint(42), builder.game.ID)
}

func TestGameBuilder_Build(t *testing.T) {
	// Given: A game builder with all fields set
	builder := NewGameBuilder().
		SetName("Catan").
		SetDescription("Trade and build").
		SetPlayerRange(3, 4).
		SetPlaytimeMinutes(60)

	// When: Building the game
	game := builder.Build()

	// Then: The game should have all the set values
	assert.NotNil(t, game)
	assert.Equal(t, "Catan", game.Name)
	assert.Equal(t, "Trade and build", game.Description)
	assert.Equal(t, 3, game.MinPlayers)
	assert.Equal(t, 4, game.MaxPlayers)
	assert.Equal(t, 60, game.PlaytimeMinutes)
}

func TestGameBuilder_Build_SetsNumberOfPlayersIfZero(t *testing.T) {
	// Given: A game builder without NumberOfPlayers set
	builder := NewGameBuilder().
		SetName("Test").
		SetMaxPlayers(5)

	// When: Building the game
	game := builder.Build()

	// Then: NumberOfPlayers should be set to MaxPlayers
	assert.Equal(t, 5, game.NumberOfPlayers)
}

func TestGameBuilder_Reset(t *testing.T) {
	// Given: A game builder with custom values
	builder := NewGameBuilder().
		SetName("Custom").
		SetComplexity(ComplexityExpert)

	// When: Resetting the builder
	result := builder.Reset()

	// Then: The builder should have default values
	assert.Equal(t, builder, result)
	assert.Equal(t, "", builder.game.Name)
	assert.Equal(t, ComplexityMedium, builder.game.Complexity)
	assert.Equal(t, 2, builder.game.MinPlayers)
	assert.Equal(t, 4, builder.game.MaxPlayers)
}

func TestGameBuilder_DefaultValues(t *testing.T) {
	// Given: A new game builder

	// When: Creating a new builder
	builder := NewGameBuilder()

	// Then: It should have default values
	assert.Equal(t, 2, builder.game.MinPlayers)
	assert.Equal(t, 4, builder.game.MaxPlayers)
	assert.Equal(t, 30, builder.game.PlaytimeMinutes)
	assert.Equal(t, 8, builder.game.MinAge)
	assert.Equal(t, ComplexityMedium, builder.game.Complexity)
	assert.Equal(t, CategoryStrategy, builder.game.Category)
	assert.Equal(t, 2024, builder.game.YearPublished)
	assert.Equal(t, 0.0, builder.game.Rating)
}

func TestGameBuilder_ChainedCalls(t *testing.T) {
	// Given: A new game builder

	// When: Chaining multiple setter calls
	game := NewGameBuilder().
		SetName("Monopoly").
		SetDescription("Real estate trading").
		SetPlayerRange(2, 8).
		SetPlaytimeMinutes(180).
		SetMinAge(8).
		SetComplexity(ComplexityMedium).
		SetCategory(CategoryFamily).
		SetPublisher("Hasbro").
		SetYearPublished(1935).
		SetRating(7.5).
		Build()

	// Then: All values should be correctly set
	assert.Equal(t, "Monopoly", game.Name)
	assert.Equal(t, "Real estate trading", game.Description)
	assert.Equal(t, 2, game.MinPlayers)
	assert.Equal(t, 8, game.MaxPlayers)
	assert.Equal(t, 180, game.PlaytimeMinutes)
	assert.Equal(t, 8, game.MinAge)
	assert.Equal(t, ComplexityMedium, game.Complexity)
	assert.Equal(t, CategoryFamily, game.Category)
	assert.Equal(t, "Hasbro", game.Publisher)
	assert.Equal(t, 1935, game.YearPublished)
	assert.Equal(t, 7.5, game.Rating)
}

func TestGameComplexity_Values(t *testing.T) {
	// Given: Game complexity constants

	// When: Checking their values

	// Then: They should have the correct string values
	assert.Equal(t, GameComplexity("Easy"), ComplexityEasy)
	assert.Equal(t, GameComplexity("Medium"), ComplexityMedium)
	assert.Equal(t, GameComplexity("Hard"), ComplexityHard)
	assert.Equal(t, GameComplexity("Expert"), ComplexityExpert)
}

func TestGameCategory_Values(t *testing.T) {
	// Given: Game category constants

	// When: Checking their values

	// Then: They should have the correct string values
	assert.Equal(t, GameCategory("Strategy"), CategoryStrategy)
	assert.Equal(t, GameCategory("Party"), CategoryParty)
	assert.Equal(t, GameCategory("Family"), CategoryFamily)
	assert.Equal(t, GameCategory("Card"), CategoryCard)
	assert.Equal(t, GameCategory("Dice"), CategoryDice)
	assert.Equal(t, GameCategory("Cooperative"), CategoryCooperative)
}
