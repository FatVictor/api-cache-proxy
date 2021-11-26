package services

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewCaching(config ...fiber.Config) fiber.Handler {
	options, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	rdb := redis.NewClient(options)
	_ttl := os.Getenv("CACHE_TTL")
	ttl, err := strconv.Atoi(_ttl)
	if err != nil || ttl < 0 {
		ttl = 10 * 60
	}
	return func(c *fiber.Ctx) error {
		key := hash(c.Request())
		val, _err := rdb.Get(context.Background(), key).Result()

		if _err != nil {
			if _err == redis.Nil {
				log.Println("There is no value in cache, requesting to original site")
			} else {
				log.Println("Not able to connect to redis, please help to very connection", _err.Error())
			}
			if e := c.Next(); e != nil {
				return e
			}
			c.Response().Header.Del(fiber.HeaderServer)
			if c.Response().StatusCode() != fiber.StatusOK {
				return nil
			}
			headers := map[string]string{}
			cookies := map[string]string{}
			c.Response().Header.VisitAll(func(key, value []byte) {
				headers[string(key)] = string(value)
			})
			c.Response().Header.VisitAllCookie(func(key, value []byte) {
				cookies[string(key)] = string(value)
			})
			_json, _ := json.Marshal(ProxyResponse{
				Headers: headers,
				Cookies: cookies,
				Content: base64.StdEncoding.EncodeToString(c.Response().Body()),
			})
			rdb.Set(context.Background(), key, string(_json), time.Second*time.Duration(ttl))
		} else {
			var r ProxyResponse
			json.Unmarshal([]byte(val), &r)
			data, _ := base64.StdEncoding.DecodeString(r.Content)
			c.Send(data)
			for k, v := range r.Headers {
				c.Response().Header.Set(k, v)
			}
		}
		return nil
	}
}
func hash(req *fiber.Request) string {
	h := sha1.New()
	h.Write(req.RequestURI())
	h.Write(req.Body())
	req.Header.VisitAll(func(key, value []byte) {
		if strings.EqualFold(fiber.HeaderAccept, string(key)) ||
			strings.EqualFold(fiber.HeaderAcceptEncoding, string(key)) ||
			strings.EqualFold(fiber.HeaderAcceptLanguage, string(key)) {
			h.Write(value)
		}
	})
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

type ProxyResponse struct {
	Headers map[string]string `json:"headers"`
	Cookies map[string]string `json:"cookies"`
	Content string            `json:"content"`
}
