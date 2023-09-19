package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-template/pkg"
	"log"
)

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(channel string, msg string) error {
	var err error
	fmt.Println("Publish 。。。。", msg)
	err = pkg.Red.Publish(channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(channel string) (*redis.PubSub, error) {
	log.Println("开始订阅 Redis 消息...")
	sub := pkg.Red.Subscribe(channel) // 订阅 Redis
	_, err := sub.Receive()           // 从 Redis 订阅中读取
	if err != nil {
		log.Printf("订阅 Redis 遇到错误: %v\n", err)
		return nil, err
	}
	log.Println("成功订阅 Redis!")
	return sub, nil
}
