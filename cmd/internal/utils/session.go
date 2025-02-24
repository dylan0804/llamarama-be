package utils

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionStore struct {
	client *redis.Client
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		}),
	}
}

func (s *SessionStore) Get(ctx context.Context, token string) (map[string]string, error) {
	user, err := s.client.HGetAll(ctx, token).Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	return s.client.Del(ctx, sessionID).Err()
}

func (s *SessionStore) CreateToken(ctx context.Context, userID string) (string, error) {
	token := uuid.New().String()

	err := s.client.HSet(ctx, token, map[string]interface{}{
		"user_id": userID,
		"created_at": time.Now(),
	}).Err()
	if err != nil {
		return "", err
	}
	
	s.client.Expire(ctx, token, 24*time.Hour)

	return token, nil
}