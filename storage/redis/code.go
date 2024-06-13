package redis_service

import (
	"app/pkg/logs"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CodeCache struct {
	client *redis.Client
	log    logs.LoggerInterface
}

type VerificationCode map[string]interface{}

func NewCodeCache(r *redis.Client, log logs.LoggerInterface) *CodeCache {
	return &CodeCache{client: r, log: log}
}

func (cc *CodeCache) validateCodeObjectKey(s string) string {
	return fmt.Sprintf("code-object-e-commerce-market-v2:%s", s)
}

func (cc *CodeCache) SetCode(ctx context.Context, email, code string, expire time.Time) error {
	data := VerificationCode{
		"code":   code,
		"expire": expire.Format(time.RFC3339Nano),
	}
	jsonData, err := json.Marshal(&data)
	if err != nil {
		cc.log.Error("could not marshal redis map object", logs.Error(err))
		return err
	}

	err = cc.client.Set(
		ctx,
		cc.validateCodeObjectKey(email),
		jsonData,
		time.Hour*1,
	).Err()

	if err != nil {
		cc.log.Error("could not set code object to redis", logs.Error(err))
		return err
	}
	return nil
}

func (cc *CodeCache) GetCode(ctx context.Context, email string) (string, time.Time, error) {
	val, err := cc.client.Get(ctx, cc.validateCodeObjectKey(email)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", time.Time{}, ErrKeyNotFound
		}
		return "", time.Time{}, err
	}

	var verCode VerificationCode
	err = json.Unmarshal([]byte(val), &verCode)
	if err != nil {
		cc.log.Error("could not unmarshal code object from redis", logs.Error(err))
		return "", time.Time{}, err
	}

	code, ok := verCode["code"].(string)
	if !ok {
		cc.log.Error("something not ok with code from redis object", logs.String("code", verCode["code"].(string)))
		return "", time.Time{}, err
	}

	expireStr, ok := verCode["expire"].(string)
	if !ok {
		cc.log.Error("something not ok with expire from redis object", logs.String("expire", verCode["expire"].(string)))
	}

	expire, err := time.Parse(time.RFC3339Nano, expireStr)
	if err != nil {
		cc.log.Error("could not parse expire object from RFC3339Nano to Time", logs.Error(err))
		return "", time.Time{}, err
	}

	return code, expire, nil
}

func (cc *CodeCache) DeleteCode(ctx context.Context, key string) error {
	_, err := cc.client.Del(ctx, cc.validateCodeObjectKey(key)).Result()
	if err != nil {
		cc.log.Error("could not delete code-object from redis", logs.Error(err), logs.String("redis-key", key))
		return err
	}
	return nil
}
