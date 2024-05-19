package search

import (
	publish_utils "chatFileBackend/utils/publish/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {
	// 模拟搜索结果
	// for (let i = 0; i < 123; i++) {
	//     results.push({
	//         title: title + " Result " + (i + 1),
	//         url: "http://example.com",
	//         tags: [tag]
	//     });
	// }
	var noarg bool = true // 没有限制任何搜索

	title_str, ok := c.GetQuery("title")
	if ok {
		noarg = false
	}

	tags_str, ok := c.GetQuery("tags")
	if ok {
		noarg = false
	}

	if noarg {
		c.JSON(500, gin.H{"message": "Need at least one argument"})
		return
	}

	// 模拟搜索
	tags := publish_utils.Str2Tags(tags_str)

	var results []SearchResult = search(title_str, tags)

	// for i := 0; i < 100; i += 1 {
	// 	restmp := SearchResult{Title: title_str + fmt.Sprint(i),
	// 		URL:  "123.com",
	// 		Tags: tags}

	// 	results = append(results, restmp)
	// }
	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"message": "success"})
}
