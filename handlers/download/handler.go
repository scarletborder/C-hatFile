package download

import (
	chats3 "chatFileBackend/utils/publish/s3"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DownloadHandler(c *gin.Context) {
	// get
	obj_name, ok := c.GetQuery("name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no obj name"})
		return
	}
	reader, leng, err := chats3.GetDownlodReaderByObjectName(obj_name)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}
	// 推断文件的MIME类型
	ext := filepath.Ext(obj_name)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream" // 如果无法推断类型，使用通用类型
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+obj_name)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprint(leng))

	// 将io.Reader写入响应
	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		c.String(500, "Error while sending file: %s", err.Error())
	}
}
