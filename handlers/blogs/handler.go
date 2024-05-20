package blogs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BlogPreviewHandler(c *gin.Context) {
	if errMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}

func BlogGetHandler(c *gin.Context) {
	title := c.Param("id")
	content, err := GetBlogContent(title)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read file"})
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}
