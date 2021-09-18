package redis

import (
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const Nil = redis.Nil

func New(dsn string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to ping redis server")
	}
	return client, nil
}

// 実装の参考のため保存
func SetValue(savePath string, value string) error {
	redisPath := os.Getenv("REDIS_HOST")
	client, err := New(redisPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer client.Close()

	err = client.Set(savePath, value, -1).Err()
	if err != nil {
		return errors.Wrap(err, "Failed to save item")
	}
	return nil
}

func GetValue(savePath string) (string, error) {
	redisPath := os.Getenv("REDIS_HOST")
	log.Println(redisPath)
	client, err := New(redisPath)
	if err != nil {
		log.Println(err)
		return "error", err
	}
	defer client.Close()
	//セーブパスが存在するかチェックする
	err = client.Get(savePath).Err()
	if err == redis.Nil {
		err = client.Set(savePath, "init", -1).Err()
		if err != nil {
			return "error", errors.Wrap(err, "Failed to get redis client")
		}
	} else if err != nil {
		return "error", errors.Wrapf(err, "Failed to get %s", savePath)
	} else {
		value, err := client.Get(savePath).Result()
		if err != nil {
			return "error", errors.Wrap(err, "Failed to save item")
		}
		return value, nil
	}
	return "error", errors.New("an unexpected error has occurred...")
}
