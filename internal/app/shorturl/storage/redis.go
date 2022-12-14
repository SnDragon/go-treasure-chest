package storage

import (
	"fmt"
	"github.com/SnDragon/go-treasure-chest/encoding"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/config"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/entity"
	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"log"
	"net/http"
	"time"

	serrors "github.com/SnDragon/go-treasure-chest/internal/app/shorturl/errors"
)

const (
	Offset = 1000000
	// RedisKeyUrlGlobalId 全局Id
	RedisKeyUrlGlobalId = "url:global:id"
	RedisKeyShortUrl    = "url:short:%s"
	RedisKeyUrlDetail   = "url:detail:%s"
	RedisKeyUrlCounter  = "url:counter:%s"
)

// RedisStorage Redis实现短链服务
type RedisStorage struct {
	redisCli *redis.Client
}

func NewRedisStorage() (*RedisStorage, error) {
	redisConf := config.AppConfig.RedisConfig
	redisStorage := &RedisStorage{
		redisCli: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", redisConf.DBHost, redisConf.DBPort),
			Password: redisConf.DBPasswd,
			DB:       redisConf.DB,
		}),
	}
	// 连接redis
	if _, err := redisStorage.redisCli.Ping().Result(); err != nil {
		return nil, err
	}
	return redisStorage, nil
}

func (r *RedisStorage) Shorten(url string, expSecond int64) (string, error) {
	// 1. 获取自增id
	id, err := r.redisCli.Incr(RedisKeyUrlGlobalId).Result()
	if err != nil {
		return "", errors.Wrap(err, "[Shorten] incr global id err")
	}
	// 2. 转成base62(base64包含`+`、`/`字符,对URL不友好)
	sid := base62.EncodeInt64(Offset + id)
	// 3. 设置短url对应的原始url
	if err := r.redisCli.Set(fmt.Sprintf(RedisKeyShortUrl, sid), url,
		time.Second*time.Duration(expSecond)).Err(); err != nil {
		return "", errors.Wrap(err, "[Shorten] set RedisKeyShortUrl err")
	}
	// 4. 设置详情
	urlDetail := &entity.UrlDetailInfo{
		OriginUrl: url,
		CreatedAt: time.Now().Unix(),
		ExpiredAt: time.Now().Unix() + expSecond,
	}
	if err := r.redisCli.Set(fmt.Sprintf(RedisKeyUrlDetail, sid),
		encoding.JsonMarshalString(urlDetail), 0).Err(); err != nil {
		return "", errors.Wrap(err, "[Shorten] set RedisKeyUrlDetail err")
	}
	return config.AppConfig.BaseUrl + sid, nil
}

func (r *RedisStorage) ShortLinkInfo(sid string) (*entity.UrlDetailInfo, error) {
	// 1. 获取详情
	data, err := r.redisCli.Get(fmt.Sprintf(RedisKeyUrlDetail, sid)).Result()
	if err != nil {
		return nil, errors.Wrap(err, "[ShortLinkInfo] get url detail err")
	}
	// 2. 反序列化
	info := &entity.UrlDetailInfo{}
	if err := encoding.JsonUnMarshalString(data, info); err != nil {
		return nil, errors.Wrapf(err, "[ShortLinkInfo] JsonUnMarshalString err: %v", data)
	}
	// 3. 获取计数器
	countRet, err := r.redisCli.Get(fmt.Sprintf(RedisKeyUrlCounter, sid)).Result()
	if err == redis.Nil {
		countRet = "0"
	} else if err != nil {
		return nil, errors.Wrapf(err, "[ShortLinkInfo] get RedisKeyUrlCounter err, sid: %v", sid)
	}
	info.Counter = cast.ToInt64(countRet)
	return info, nil
}

func (r *RedisStorage) UnShorten(sid string) (string, error) {
	// 1. 获取对应长链
	val, err := r.redisCli.Get(fmt.Sprintf(RedisKeyShortUrl, sid)).Result()
	if err == redis.Nil {
		return "", &serrors.StatusError{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("unknown url: %v", sid),
		}
	} else if err != nil {
		return "", errors.Wrap(err, "get RedisKeyShortUrl err")
	}
	// 2. 访问计数器+1
	if err := r.redisCli.Incr(fmt.Sprintf(RedisKeyUrlCounter, sid)).Err(); err != nil {
		// 只影响统计，不影响主流程，打印错误日志即可
		log.Printf("[UnShorten]Incr RedisKeyUrlCounter err, sid: %v\n", sid)
	}
	return val, nil
}
