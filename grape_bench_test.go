package main

import (
	"github.com/go-redis/redis"
	"strconv"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6001",
		Password: "",
		DB:       0,
	})

	for i := 0; i < 100; i++ {
		err := client.Set("shq"+strconv.Itoa(i), "value", 0).Err()
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 100; i++ {
		err := client.Set("lsy"+strconv.Itoa(i), "value", 0).Err()
		if err != nil {
			panic(err)
		}
	}

}

func BenchmarkGet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6002",
		Password: "",
		DB:       0,
	})

	for i := 0; i < 100; i++ {
		value, err := client.Get("shq" + strconv.Itoa(i)).Result()
		if err != nil && value != "value" {
			panic(err)
		}
	}

	for i := 0; i < 100; i++ {
		value, err := client.Get("lsy" + strconv.Itoa(i)).Result()
		if err != nil && value != "value" {
			panic(err)
		}
	}
}
