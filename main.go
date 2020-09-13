package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

func NewRdb() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "B5fUo0kpCUoWgFMb", // no password set
		DB:       3,                  // use default DB
	})
	return rdb
}

var rdb = NewRdb()

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {

		var memberName = uuid.New().String()
		var member = redis.Z{
			Member: memberName,
			Score:  cast.ToFloat64(time.Now().Unix()),
		}

		_, err := rdb.ZRank(c, "qps:records", memberName).Result()

		if err == nil {
			member.Score += 1
		}

		rdb.ZAdd(c, "qps:records", &member)

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/records", func(c *gin.Context) {

		var total = 0

		var data = make(map[string]int)
		res, err := rdb.ZRangeWithScores(c, "qps:records", 0, -1).Result()
		if err == nil {
			for _, value := range res {
				total += 1
				var date = time.Unix(cast.ToInt64(value.Score), 0).Format("2006-01-02 15:04:05")
				if _, ok := data[date]; ok {
					data[date] += 1
				} else {
					data[date] = 1
				}
			}
		}

		c.JSON(200, gin.H{
			"data":  data,
			"total": total,
		})
	})
	r.Run()
}
