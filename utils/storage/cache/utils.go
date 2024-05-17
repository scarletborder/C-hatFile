package cached

import (
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

/*
CacheGet

通过key+cacheditem的一个独一无二的作为键去访问cache

必须由CacheSet设置

@param key 领域

@param obj 指针
*/
func CacheGet(key string, obj interface{}) (bool, error) {
	c := GetOrCreateCache(key)
	wrappedObj := obj.(cachedItem)

	key = key + id2Str(wrappedObj.GetID()) // cache的key
	if val, found := c.Get(key); found {
		c.Set(key, val, 5*time.Minute)
		dstVal := reflect.ValueOf(obj)
		dstVal.Elem().Set(reflect.ValueOf(val))
		return true, nil
	}
	// 缓存没有命中
	return false, nil
}

/*
CacheGetByStr

通过key+cacheditem(的特征值) 的一个独一无二的作为键去访问cache

必须由CacheSetByStr设置

@param key 领域

@param obj 指针
*/
func CacheGetByStr(key string, obj interface{}) (bool, error) {
	c := GetOrCreateCache(key)
	wrappedObj := obj.(cachedItem)
	key = key + wrappedObj.GetFeature()
	if val, found := c.Get(key); found {
		c.Set(key, val, 6*time.Minute)
		dstVal := reflect.ValueOf(obj)
		dstVal.Elem().Set(reflect.ValueOf(val))
		return true, nil
	}
	// 缓存没有命中
	return false, nil
}

// CacheSet
//
// 一般Fallback到db后设置某cacheditem
//
// @param obj 值传递
func CacheSet(key string, obj interface{}) {
	c := GetOrCreateCache(key)
	wrappedObj := obj.(cachedItem)

	key = key + id2Str(wrappedObj.GetID()) // 真key
	if _, find := c.Get(key); find {
		wrappedObj.Dirty()
	}
	c.Set(key, wrappedObj, 6*time.Minute)
}

// CacheSetByStr
//
// 一般Fallback到db后设置某cacheditem
//
// @param obj 值传递
func CacheSetByStr(key string, obj interface{}) {
	c := GetOrCreateCache(key)
	wrappedObj := obj.(cachedItem)

	key = key + wrappedObj.GetFeature() // 真key
	if _, find := c.Get(key); find {
		wrappedObj.Dirty()
	}
	c.Set(key, wrappedObj, 6*time.Minute)
}

// StartDBSync
//
// 开始和db同步，定时写入
//
// sync_fn 接受CacheGet的数据，并通过某种方式同步db
func StartDBSync(key string, sync_handler func(chunk []interface{}) error, t time.Duration) {
	go func() {
		daemon := time.NewTimer(t)
		for range daemon.C { // 死循环直至程序终止
			// 执行传递的函数
			c := GetOrCreateCache(key)
			items := c.Items()
			item_slice := make([]interface{}, 0, len(items))

			for _, v := range items {
				if cv, ok := v.Object.(cachedItem); ok && cv.IsDirty() {
					item_slice = append(item_slice, v.Object)
				}
			}
			c.DeleteExpired()
			go func() {
				err := sync_handler(item_slice)
				if err != nil {
					logrus.Errorln("sync error in ", key)
				}
			}()

			// err := sync_fn
			// if err != nil {
			// 	// fmt.Println("Error executing action:", err)
			// }
		}
	}()
}
