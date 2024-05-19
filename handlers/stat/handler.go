package stat

import "github.com/gin-gonic/gin"

func StatRankHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"stamp":   update_stamp,
		"results": rank_result})
}
