package pkg

import (
	"github.com/go-redis/redis"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

const LOCAL_USERID_UINT = "user_id_uint"
const LOCAL_USERID_INT64 = "user_id_int64"
const LOCAL_TOKEN = "token"

const MESSAGE_FAIL = -1
const TOKEN_FAIL = -2

// 检测结构体
var Validate = validator.New()

var DB *gorm.DB

var Red *redis.Client
