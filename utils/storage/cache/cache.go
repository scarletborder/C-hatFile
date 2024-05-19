package cached

import (
	"strconv"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	cacheMaps = make(map[string]*cache.Cache)
	m         = new(sync.Mutex)

	// time settings
	Expired_time  = 5 * time.Minute
	Sync_interval = 6 * time.Second
)

type cachedItem interface {
	SetDirty() // 标记进行修改

	FlushDirty()

	GetID() uint64 // 获得唯一后缀字符串

	GetFeature() string // 特征字符串

	IsDirty() bool // 是否被修改
}

func id2Str(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// 内置方法，不建议直接使用
func GetOrCreateCache(name string) *cache.Cache {
	cc, ok := cacheMaps[name]
	if !ok {
		// 注册对应领域的cache
		m.Lock()
		defer m.Unlock()
		cc, ok := cacheMaps[name]
		if !ok {
			cc = cache.New(Expired_time, 10*time.Minute)
			cacheMaps[name] = cc
			return cc
		}
		return cc
	}
	return cc
}
