package service

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestCreateRedisClient(t *testing.T) {
	c, err := CreateRedisClient()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	var ctx = context.Background()
	val, err := c.Get(ctx, "").Result()

	fmt.Println(val)
}
