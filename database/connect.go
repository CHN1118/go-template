package database

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-template/config"
	"go-template/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

// ConnectDB connect to db
func ConnectDB() (db *gorm.DB, red *redis.Client) {
	allModels := []interface{}{
		&models.UserBasic{},
	}
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic(err)
	}

	sqlLog := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second, //慢SQL阈值
		LogLevel:      logger.Info, //级别
		Colorful:      true,        //彩色
	})

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	fmt.Println(dsn)
	DB, err := gorm.Open(postgres.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
			PrepareStmt:                              true, // 开启自动更新UpdatedAt字段
			Logger:                                   sqlLog,
		})
	if err != nil {
		panic("failed to connect database")
	}

	//创表
	for _, m := range allModels {
		if !DB.Migrator().HasTable(m) {
			if err = DB.AutoMigrate(m); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Database Connected")

	// 初始化redis
	Red := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     30,
		MinIdleConns: 30,
	})
	pong, err := Red.Ping().Result()
	if err != nil {
		fmt.Println("init redis 。。。。", err)
	} else {
		fmt.Println(" Redis inited 。。。。", pong) //
	}

	return DB, Red
}
