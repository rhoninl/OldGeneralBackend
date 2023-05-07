package main

import (
	"context"
	"log"
	"time"

	"github.com/leepala/OldGeneralBackend/pkg/database"
)

func main() {
	rdb := database.GetRDB()
	c := rdb.Set(context.TODO(), "aaa", "bbb", time.Second)
	r, e := c.Result()
	log.Println(r, e)
	sc := rdb.Get(context.TODO(), "aaa")
	d := sc.String()
	log.Println(d)
	ic := rdb.Del(context.TODO(), "aaa")
	log.Println(ic.Result())
	sc = rdb.Get(context.TODO(), "aaa")
	d = sc.Val()
	log.Println(d)
}
