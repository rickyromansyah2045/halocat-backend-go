package config

import (
	"fmt"
	"log"
	"time"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis"
	"github.com/rickyromansyah2045/halocat-backend-go/constant"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func InitDB(isProduction bool) *gorm.DB {
	// Read and Write Connection
	dsnRW := fmt.Sprintf(constant.STRING_DSN, constant.DB_USER, constant.DB_PASS, constant.RW_HOST, constant.DB_PORT, constant.DB_NAME)
	db, err := gorm.Open(mysql.Open(dsnRW), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB.SetMaxOpenConns(constant.DB_MAX_OPEN_CONNECTIONS)
	sqlDB.SetMaxIdleConns(constant.DB_MAX_IDLE_CONNECTIONS)
	sqlDB.SetConnMaxLifetime(time.Minute)

	// Read Only Connection
	if constant.RO_HOST != "" {
		dsnRO := fmt.Sprintf(constant.STRING_DSN, constant.DB_USER, constant.DB_PASS, constant.RO_HOST, constant.DB_PORT, constant.DB_NAME)

		db.Use(
			dbresolver.Register(dbresolver.Config{
				Replicas:          []gorm.Dialector{mysql.Open(dsnRO)},
				Policy:            dbresolver.RandomPolicy{},
				TraceResolverMode: true,
			}).SetMaxOpenConns(constant.DB_MAX_OPEN_CONNECTIONS).SetMaxIdleConns(constant.DB_MAX_IDLE_CONNECTIONS),
		)
	}

	if constant.DB_CACHING {
		var redisClient *redis.Client

		if isProduction {
			redisOptions := &redis.Options{
				Addr:     fmt.Sprintf("%v:%v", constant.REDIS_HOST, constant.REDIS_PORT),
				Password: constant.REDIS_PASS,
			}

			redisClient = redis.NewClient(redisOptions)
		} else {
			redisOptions := &redis.Options{
				Addr: fmt.Sprintf("%v:%v", constant.REDIS_HOST, constant.REDIS_PORT),
			}

			redisClient = redis.NewClient(redisOptions)
		}

		if _, err := redisClient.Ping().Result(); err != nil {
			if err.Error() == "ERR AUTH <password> called without any password configured for the default user. Are you sure your configuration is correct?" {
				redisOptions := &redis.Options{
					Addr: fmt.Sprintf("%v:%v", constant.REDIS_HOST, constant.REDIS_PORT),
				}

				redisClient = redis.NewClient(redisOptions)

				if _, err := redisClient.Ping().Result(); err != nil {
					log.Fatal("error connection to redis server, error: ", err.Error())
				}
			} else {
				log.Fatal("error connection to redis server, error: ", err.Error())
			}
		}

		gormCache, _ := cache.NewGorm2Cache(&config.CacheConfig{
			CacheLevel:           config.CacheLevelAll,
			CacheStorage:         config.CacheStorageRedis,
			RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
			InvalidateWhenUpdate: true,
			CacheTTL:             3000, // 3000 = 3 second
			CacheMaxItemCnt:      0,    // 0 = cache all queries
		})

		db.Use(gormCache)
	}

	if err := db.Raw(helper.ConvertToInLineQuery("SET GLOBAL FOREIGN_KEY_CHECKS = 0;")).Error; err != nil {
		log.Fatal(err.Error())
	}

	return db
}
