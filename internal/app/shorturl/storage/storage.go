package storage

import "github.com/SnDragon/go-treasure-chest/internal/app/shorturl/entity"

// Storage 短链服务抽象接口
type Storage interface {
	Shorten(url string, expSecond int64) (string, error)     // 将长链转成短链,并设置过期时间
	ShortLinkInfo(sid string) (*entity.UrlDetailInfo, error) // 根据短链id获取详情
	UnShorten(sid string) (string, error)                    // 根据短链id转成原始长链
}
