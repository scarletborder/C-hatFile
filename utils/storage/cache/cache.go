package cached

import (
	"strconv"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var cacheMaps = make(map[string]*cache.Cache)

var m = new(sync.Mutex)

type cachedItem interface {
	GetID() uint64      // 获得唯一后缀字符串
	GetFeature() string // 特征字符串
	IsDirty() bool      // 是否被修改
	Dirty()             // 标记进行修改
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
			cc = cache.New(5*time.Minute, 10*time.Minute)
			cacheMaps[name] = cc
			return cc
		}
		return cc
	}
	return cc
}
