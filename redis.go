package midimaggot

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

const redisPrefix = "midimaggot"

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func addProgramNumber(pn int, command string) {
	fmt.Println("adding program...")
	err := redisClient.SAdd(redisPrefix+":programNumbers", pn).Err()
	if err != nil {
		panic(err)
	}
	pData := map[string]interface{}{
		"pn":      pn,
		"command": command,
	}
	err = redisClient.HMSet(redisPrefix+":program:"+strconv.FormatInt(int64(pn), 10), pData).Err()
	if err != nil {
		panic(err)
	}
}

func retrieveCommandViaProgramNumber(pn int64) string {
	exists, err := redisClient.SIsMember(redisPrefix+":programNumbers", pn).Result()
	if err != nil {
		panic(err)
	}
	if exists {
		command, err := redisClient.HGet(redisPrefix+":program:"+strconv.FormatInt(pn, 10), "command").Result()
		if err != nil {
			panic(err)
		}
		return command
	} else {
		return ""
	}
}
