<img src="https://longerwu-1252728875.cos.ap-guangzhou.myqcloud.com/blogs/shorturl_mind.png" width = "500" height = "500" alt="思维导图" />

## 简介

### 什么是短地址(短链)服务?

> 可以将原本的一大串很长的地址缩短到一个很短的地址，用户访问这个短地址可以重定向到原本的长地址

比较著名的短链域名有腾讯的`url.cn`,微博的`t.cn`等，也有很多公司提供了免费/收费的短链接开放API，感兴趣的可以自行网上搜索。

### 使用场景
* 提升用户体验, 例如淘宝商品详情url不可避免地会跟着很多参数，相较繁长的字符串, 使用简短的链接分享对用户来说观感更好
* 体现品牌专业度,就像QQ号@qq.com会显得不够专业一样，使用短地址而非长地址进行营销推广在一定程度上能提升用户对品牌的好感
* 避免url被截断，例如当调用第三方平台接口发送消息时，如果消息内容中包含url，里面的#、?等特殊字符可能在客户端被截断,导致收到消息的用户打不开该链接

### 弊端
使用短链服务可能会有哪些弊端呢?
1. 成本，不使用第三方服务的话,需要自行搭建一套服务来实现长链转短链,以及支持短链跳转,使用第三方服务则可能需要收费
2. 安全,如果短链服务是对外开放的，可能会被黑产或不法分子利用
3. 时效,例如短链设置的过期时间较长，可能原本的长链已经失效了，这时短链就会打不开
4. 速度, 因为有个短链重定向到长链的过程

抛开以上弊端不讲,本着学习的心态，本文将介绍如何使用Go语言开发一
个简单的短链服务,简述背后的原理, 末尾会有项目代码。

## 实现原理
<!-- more -->
### 流程
实现短链服务并不复杂,我们先看看使用短链的流程，如图所示:

![流程图](https://longerwu-1252728875.cos.ap-guangzhou.myqcloud.com/blogs/short_url_flow_chart.png)
1. 调用短链服务转换接口，输入长链
2. 短链服务生成对应的短链,并存储映射关系
3. 客户端(一般是浏览器)访问短链
4. 短链服务查找映射关系
5. 如果能找到对应的长链且在有效期内，则重定向到长链，否则返回404
6. 客户端访问长链


可见，这里的短链服务核心是如何长链映射成短链并存储，以及根据短链查找对应的长链并重定向给客户端，即短链生成与查找

### 短链生成
需要保证的点:
1. 全局唯一
2. 尽可能短
3. 利于查找

常用的短链接算法主要有两种:

| 算法 | 简述 | 优点 |缺点 |
| --- | --- | --- |--- |
| 哈希 | 原地址通过哈希函数生成哈希值 | 本地计算,无需依赖第三方组件 | 可能存在哈希碰撞，当哈希冲突时需要rehash或其他处理，且哈希值一般不会很短
| 分布式ID | 借助一定算法或外部组件生成全局唯一ID | 可保证全局唯一,其中自增ID的形式较短 | 通常需要依赖外部组件如MySQL,Redis,Zookeeper等

### 短链映射存储
选择较多,基本要求:
* 快速查找
* 可以设置过期时间
* 持久化


## 开发思路
由于Redis具有丰富的数据结构，支持设置过期时间，持久化，高性能等特点,这里我们使用Redis来生成自增id,并存储映射关系,但在设计会考虑可扩展为其他存储的能力。

额外说明的一点: 这里自增id会转成62进制的字符串(A-Za-z0-0)，作为短链id(下文都称sid), 可以让整体字符串更短些
*不使用base64是因为base64包含+、/两个特殊符号, 对URL不友好*

### 项目用到的库:
* [gorilla/mux](https://github.com/gorilla/mux): http路由
* [go-redis](https://github.com/go-redis/redis): redis客户端
* [base62](https://github.com/mattheath/base62): 可以将数字转成62进制的字符串
* [validator](https://gopkg.in/go-playground/validator.v9): 请求参数验证

### 实现接口
为方便后续扩展，我们可以定义一个抽象的Storage接口，并提供一个Redis的实现，后续也可替换成其他实现
```go
// Storage 短链服务抽象接口
type Storage interface {
	Shorten(url string, expSecond int64) (string, error)     // 将长链转成短链,并设置过期时间
	ShortLinkInfo(sid string) (*entity.UrlDetailInfo, error) // 根据短链id获取详情
	UnShorten(sid string) (string, error)                    // 根据短链id转成原始长链
}
```


```go
// RedisStorage Redis实现短链服务
type RedisStorage struct {
	redisCli *redis.Client
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
```

## 效果测试
### 配置
修改configs/shorturl/app.yaml的配置
```yaml
base_url: http://myurl.cn/
redis_config:
  db_host: 127.0.0.1
  db_port: 6379
  db_passwd:
  db: 0
```

假设我们的短链域名为myurl.cn,配置host
```
127.0.0.1 myurl.cn
```
### 启动服务
进入项目根mulu,启动服务:
```bash
$ make run_short_url
go build -o build/shorturl ./cmd/shorturl  && ./build/shorturl
2022/12/14 20:58:47 resource.go:35: confFile: configs/shorturl/app.yaml
2022/12/14 20:58:47 resource.go:44: app config init succeed, conf: &{BaseHost:127.0.0.1 BasePort:80 RedisConfig:{DBHost:127.0.0.1 DBPort:6379 DBPasswd: DB:0}}
2022/12/14 20:58:47 app.go:40: App run in :80...
```

### 长链转短链接口
```bash
curl -X POST \
  http://myurl.cn/api/shorten \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{"url":"https://www.baidu.com?name=SnDragon","expire_seconds":100}'
```
返回:
```
{
    "code": 0,
    "msg": "ok",
    "short_url": "http://myurl.cn/4C99"
}
```

### 重定向
浏览器打开`http://myurl.cn/4C99`

<img src="https://longerwu-1252728875.cos.ap-guangzhou.myqcloud.com/blogs/short_url_redirect.png" alt="重定向截图" width = "500" height = "500" />

### 查看详情
```
curl -X GET \
  'http://myurl.cn/api/info?sid=4C99' \
  -H 'cache-control: no-cache' \
  -H 'postman-token: 82323d30-14b1-bac8-8198-2632fbe008e1'
```

返回:
```
{
    "code": 0,
    "msg": "ok",
    "info": {
        "origin_url": "https://www.baidu.com?name=SnDragon",
        "created_at": 1671023283,
        "expired_at": 1671023383,
        "counter": 1
    }
}
```

### 仓库地址
* https://github.com/SnDragon/go-treasure-chest