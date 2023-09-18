package utils

import (
	"fmt"
	"go-template/pkg"
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
func Subscribe(channel string) (string, error) {
	sub := pkg.Red.Subscribe(channel)
	fmt.Println("Subscribe 。。。。")
	msg, err := sub.ReceiveMessage()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe 。。。。", msg.Payload)
	return msg.Payload, err
}
