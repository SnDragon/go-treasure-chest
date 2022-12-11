package storage

import "github.com/SnDragon/go-treasure-chest/internal/app/shorturl/entity"

type Storage interface {
	Shorten(url string, expSecond int64) (string, error)
	ShortLinkInfo(sid string) (*entity.UrlDetailInfo, error)
	UnShorten(sid string) (string, error)
}
