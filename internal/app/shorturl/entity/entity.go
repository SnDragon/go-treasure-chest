package entity

type UrlDetailInfo struct {
	OriginUrl string `json:"origin_url"` // 原始URL
	CreatedAt int64  `json:"created_at"` // 创建时间
	ExpiredAt int64  `json:"expired_at"` // 过期时间
	Counter   int64  `json:"counter"`    // 访问次数
}
