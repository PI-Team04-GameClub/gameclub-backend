package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/mocks"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCachedGameRepository_FindAll_CacheHit(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	expectedGames := []models.Game{
		{Name: "Game 1"},
		{Name: "Game 2"},
	}

	mockCache.On("Get", mock.Anything, redis.KeyGameAll, mock.AnythingOfType("*[]models.Game")).
		Run(func(args mock.Arguments) {
			dest := args.Get(2).(*[]models.Game)
			*dest = expectedGames
		}).
		Return(nil)

	ctx := context.Background()
	games, err := cachedRepo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedGames, games)
	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything)
	mockCache.AssertExpectations(t)
}

func TestCachedGameRepository_FindAll_CacheMiss(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	expectedGames := []models.Game{
		{Name: "Game 1"},
		{Name: "Game 2"},
	}

	mockCache.On("Get", mock.Anything, redis.KeyGameAll, mock.AnythingOfType("*[]models.Game")).
		Return(goredis.Nil)
	mockRepo.On("FindAll", mock.Anything).Return(expectedGames, nil)
	mockCache.On("Set", mock.Anything, redis.KeyGameAll, expectedGames, ttl).Return(nil)

	ctx := context.Background()
	games, err := cachedRepo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedGames, games)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedGameRepository_FindByID_CacheHit(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	expectedGame := models.Game{Name: "Test Game"}
	gameID := "123"
	cacheKey := redis.GameByIDKey(gameID)

	mockCache.On("Get", mock.Anything, cacheKey, mock.AnythingOfType("*models.Game")).
		Run(func(args mock.Arguments) {
			dest := args.Get(2).(*models.Game)
			*dest = expectedGame
		}).
		Return(nil)

	ctx := context.Background()
	game, err := cachedRepo.FindByID(ctx, gameID)

	assert.NoError(t, err)
	assert.Equal(t, &expectedGame, game)
	mockRepo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	mockCache.AssertExpectations(t)
}

func TestCachedGameRepository_FindByID_CacheMiss(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	expectedGame := &models.Game{Name: "Test Game"}
	gameID := "123"
	cacheKey := redis.GameByIDKey(gameID)

	mockCache.On("Get", mock.Anything, cacheKey, mock.AnythingOfType("*models.Game")).
		Return(goredis.Nil)
	mockRepo.On("FindByID", mock.Anything, gameID).Return(expectedGame, nil)
	mockCache.On("Set", mock.Anything, cacheKey, expectedGame, ttl).Return(nil)

	ctx := context.Background()
	game, err := cachedRepo.FindByID(ctx, gameID)

	assert.NoError(t, err)
	assert.Equal(t, expectedGame, game)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedGameRepository_Create_InvalidatesCache(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	game := &models.Game{Name: "New Game"}

	mockRepo.On("Create", mock.Anything, game).Return(nil)
	mockCache.On("Delete", mock.Anything, []string{redis.KeyGameAll}).Return(nil)

	ctx := context.Background()
	err := cachedRepo.Create(ctx, game)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedGameRepository_Update_InvalidatesCache(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	game := &models.Game{Name: "Updated Game"}
	game.ID = 123

	mockRepo.On("Update", mock.Anything, game).Return(nil)
	mockCache.On("Delete", mock.Anything, []string{redis.GameByIDKeyUint(123), redis.KeyGameAll}).Return(nil)

	ctx := context.Background()
	err := cachedRepo.Update(ctx, game)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedGameRepository_Delete_InvalidatesCache(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	gameID := uint(123)

	mockRepo.On("Delete", mock.Anything, gameID).Return(nil)
	mockCache.On("Delete", mock.Anything, []string{redis.GameByIDKeyUint(gameID), redis.KeyGameAll}).Return(nil)

	ctx := context.Background()
	err := cachedRepo.Delete(ctx, gameID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedGameRepository_FailOpen_OnCacheError(t *testing.T) {
	mockRepo := new(mocks.MockGameRepository)
	mockCache := new(mocks.MockCache)
	ttl := 5 * time.Minute

	cachedRepo := NewCachedGameRepository(mockRepo, mockCache, ttl)

	expectedGames := []models.Game{{Name: "Game 1"}}

	mockCache.On("Get", mock.Anything, redis.KeyGameAll, mock.AnythingOfType("*[]models.Game")).
		Return(errors.New("cache connection error"))
	mockRepo.On("FindAll", mock.Anything).Return(expectedGames, nil)
	mockCache.On("Set", mock.Anything, redis.KeyGameAll, expectedGames, ttl).Return(errors.New("cache set error"))

	ctx := context.Background()
	games, err := cachedRepo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedGames, games)
	mockRepo.AssertExpectations(t)
}
