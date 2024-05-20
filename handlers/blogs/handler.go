package blogs

import (
	"net/http"
	"strconv"

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
	id := c.Param("id")
	fileID, err := strconv.Atoi(id)
	if err != nil || fileID < 0 || fileID >= len(files) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file ID"})
		return
	}
	content, err := GetBlogContent(fileID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read file"})
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}
