package onlineQeueu

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type QueueService struct {
	RedisClient *redis.Client
}

func NewQueueService(redisAddr string) *QueueService {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &QueueService{RedisClient: client}
}

func (qs *QueueService) CreateOfficeQueue(officeID uint) error {
	keyQueue := fmt.Sprintf("queue:%d", officeID)
	keyInService := fmt.Sprintf("in_service:%d", officeID)

	if err := qs.RedisClient.Del(ctx, keyQueue).Err(); err != nil {
		return err
	}
	if err := qs.RedisClient.Del(ctx, keyInService).Err(); err != nil {
		return err
	}
	return nil
}

func (qs *QueueService) AddClientToQueue(officeID uint, clientNumber int) error {
	keyQueue := fmt.Sprintf("queue:%d", officeID)
	return qs.RedisClient.RPush(ctx, keyQueue, clientNumber).Err()
}

func (qs *QueueService) RemoveClientFromQueue(officeID uint, clientNumber int) error {
	keyQueue := fmt.Sprintf("queue:%d", officeID)
	return qs.RedisClient.LRem(ctx, keyQueue, 0, clientNumber).Err()
}

func (qs *QueueService) MoveClientToInService(officeID uint, operatorID uint) (int, error) {
	keyQueue := fmt.Sprintf("queue:%d", officeID)
	keyInService := fmt.Sprintf("in_service:%d", officeID)

	clientStr, err := qs.RedisClient.LPop(ctx, keyQueue).Result()
	if err != nil {
		return 0, err
	}

	err = qs.RedisClient.HSet(ctx, keyInService, fmt.Sprintf("%d", operatorID), clientStr).Err()
	if err != nil {
		return 0, err
	}

	clientNumber, _ := strconv.Atoi(clientStr)
	return clientNumber, nil
}

func (qs *QueueService) GetClientInService(officeID uint, operatorID uint) (int, error) {
	keyInService := fmt.Sprintf("in_service:%d", officeID)
	clientStr, err := qs.RedisClient.HGet(ctx, keyInService, fmt.Sprintf("%d", operatorID)).Result()
	if err != nil {
		return 0, err
	}
	clientNumber, _ := strconv.Atoi(clientStr)
	return clientNumber, nil
}

func (qs *QueueService) FinishService(officeID uint, operatorID uint) error {
	keyInService := fmt.Sprintf("in_service:%d", officeID)
	return qs.RedisClient.HDel(ctx, keyInService, fmt.Sprintf("%d", operatorID)).Err()
}
