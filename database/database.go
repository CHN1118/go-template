package database

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Red *redis.Client
