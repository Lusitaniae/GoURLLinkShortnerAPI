/*
	Database layer for the url shortner. The GoLinkShortner API uses a redis database as a backend.
	It follows a common connection pattern where a main session is created then other sessions are created by copying the information of the main session while utilizing a different socket from a socket pool
*/
package main

import (
	"errors"
	"gopkg.in/redis.v5"
	"log"
)

var ErrDuplicate = errors.New("duplicate")
var ErrNotFound = errors.New("not found")
var ErrConnection = errors.New("connection issues")

type RedisCache struct {
	*redis.Client
}

func NewDBConnection() *RedisCache {

	return &RedisCache{
		Client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})}

}

func (this *RedisCache) FindshortUrl(longurl string) (sUrl string, err error) {

	val, err := this.Client.Get(longurl).Result()

	if err != nil {
		log.Println(err)
		return val, ErrConnection
	}

	return val, err

}

func (this *RedisCache) FindlongUrl(shortUrl string) (lUrl string, err error) {

	val, err := this.Client.Get(shortUrl).Result()

	if err != nil {
		log.Println(err)
		return val, ErrConnection
	}

	return val, err

}

func (this *RedisCache) AddUrls(longUrl string, shortUrl string) (err error) {

	val, err := this.Client.Exists(shortUrl).Result()

	if err == redis.Nil {

		return ErrNotFound

	} else if err != nil {

		return ErrConnection

	} else if val == false {

		errSet := this.Client.Set(shortUrl, longUrl, 0).Err()

		if errSet != nil {
			log.Println(errSet)
		}

		log.Printf("Saving short:%s long:%s", shortUrl, longUrl)

	} else {

		return ErrDuplicate

	}

	return

}
