package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Vaibhavsahu2810/redis-go/internal/config"
	"github.com/Vaibhavsahu2810/redis-go/internal/email"
	"github.com/go-redis/redis/v8"
)

const (
	emailQueue = "email_queue"
)

type EmailTask struct {
	To           string                 `json:"to"`
	Subject      string                 `json:"subject"`
	TemplateName string                 `json:"templateName"`
	Data         map[string]interface{} `json:"data"`
}

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("fail to connect to Redis: %w", err)
	}

	return client, nil
}

func EnqueueEmail(ctx context.Context, client *redis.Client, task EmailTask) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("error serializing task: %w", err)
	}

	if err := client.RPush(ctx, emailQueue, taskJSON).Err(); err != nil {
		return fmt.Errorf("error adding task to queue: %w", err)
	}

	return nil
}

func StartWorker(ctx context.Context, client *redis.Client, sender *email.Sender) {
	log.Println("Starting worker for email processing...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped")
			return
		default:

			result, err := client.BLPop(ctx, 0, emailQueue).Result()
			if err != nil {
				if err != redis.Nil && err != context.Canceled {
					log.Printf("Error getting task from queue: %v", err)
				}
				continue
			}

			if len(result) < 2 {
				continue
			}

			var task EmailTask
			if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
				log.Printf("Error deserializing task: %v", err)
				continue
			}

			if err := sender.SendEmail(task.To, task.Subject, task.TemplateName, task.Data); err != nil {
				log.Printf("Error sending email: %v", err)
			} else {
				log.Printf("E-mail sent to %s", task.To)
			}
		}
	}
}
