package config

type Config struct {
	CacheFile       *string    `json:"cache_file"`
	Regions         *[]*string `json:"regions"`
	CacheTTLSeconds *int64     `json:"cache_ttl_seconds"`
}
