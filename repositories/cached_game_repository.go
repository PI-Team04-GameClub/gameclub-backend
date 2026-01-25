package repositories

import (
	"context"
	"log"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/redis"
	goredis "github.com/redis/go-redis/v9"
)

type CachedGameRepository struct {
	base  GameRepository
	cache redis.Cache
	ttl   time.Duration
}

func NewCachedGameRepository(base GameRepository, c redis.Cache, ttl time.Duration) *CachedGameRepository {
	return &CachedGameRepository{
		base:  base,
		cache: c,
		ttl:   ttl,
	}
}

func (r *CachedGameRepository) FindAll(ctx context.Context) ([]models.Game, error) {
	var games []models.Game

	err := r.cache.Get(ctx, redis.KeyGameAll, &games)
	if err == nil {
		return games, nil
	}

	if err != goredis.Nil {
		log.Printf("Cache get error for %s: %v", redis.KeyGameAll, err)
	}

	games, err = r.base.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if cacheErr := r.cache.Set(ctx, redis.KeyGameAll, games, r.ttl); cacheErr != nil {
		log.Printf("Cache set error for %s: %v", redis.KeyGameAll, cacheErr)
	}

	return games, nil
}

func (r *CachedGameRepository) FindByID(ctx context.Context, id string) (*models.Game, error) {
	var game models.Game
	key := redis.GameByIDKey(id)

	err := r.cache.Get(ctx, key, &game)
	if err == nil {
		return &game, nil
	}

	if err != goredis.Nil {
		log.Printf("Cache get error for %s: %v", key, err)
	}

	result, err := r.base.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cacheErr := r.cache.Set(ctx, key, result, r.ttl); cacheErr != nil {
		log.Printf("Cache set error for %s: %v", key, cacheErr)
	}

	return result, nil
}

func (r *CachedGameRepository) Create(ctx context.Context, game *models.Game) error {
	if err := r.base.Create(ctx, game); err != nil {
		return err
	}

	if err := r.cache.Delete(ctx, redis.KeyGameAll); err != nil {
		log.Printf("Cache invalidation error for %s: %v", redis.KeyGameAll, err)
	}

	return nil
}

func (r *CachedGameRepository) Update(ctx context.Context, game *models.Game) error {
	if err := r.base.Update(ctx, game); err != nil {
		return err
	}

	key := redis.GameByIDKeyUint(game.ID)
	if err := r.cache.Delete(ctx, key, redis.KeyGameAll); err != nil {
		log.Printf("Cache invalidation error: %v", err)
	}

	return nil
}

func (r *CachedGameRepository) Delete(ctx context.Context, id uint) error {
	if err := r.base.Delete(ctx, id); err != nil {
		return err
	}

	key := redis.GameByIDKeyUint(id)
	if err := r.cache.Delete(ctx, key, redis.KeyGameAll); err != nil {
		log.Printf("Cache invalidation error: %v", err)
	}

	return nil
}
