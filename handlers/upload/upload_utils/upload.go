package upload_utils

import (
	"chatFileBackend/models"
	"chatFileBackend/utils/publish"
	"io"
)

//  根据文件流和元数据创建db相关line和往对象存储中塞入数据

func UploadFile(file io.Reader, meta *models.MetaData) (msg string, err error) {
	return publish.UploadDocument(file, meta)
}

/*
你可以使用Gorm进行数据库操作，并实现所需的插入、上传和回滚操作。以下是一个示例代码，展示如何实现这个流程：

1. **插入文件元数据并获取自增ID**
2. **尝试上传文件**
3. **如果上传失败，则删除插入的元数据**

假设你有一个 `FileMetadata` 结构体来表示文件元数据，并且已经配置好了Gorm连接。

```go
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// FileMetadata 表示文件元数据
type FileMetadata struct {
	ID       uint   `gorm:"primaryKey"`
	FileName string `gorm:"not null"`
	FileSize int64  `gorm:"not null"`
}

func main() {
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// 假设元数据已经准备好
	fileMeta := FileMetadata{
		FileName: "example.txt",
		FileSize: 1024,
	}

	// 插入元数据
	result := db.Create(&fileMeta)
	if result.Error != nil {
		log.Fatal("failed to insert file metadata: ", result.Error)
	}

	// 获取插入的ID
	insertedID := fileMeta.ID
	fmt.Printf("Inserted metadata with ID: %d\n", insertedID)

	// 假设上传文件的函数
	err = uploadFile(fileMeta)
	if err != nil {
		// 如果上传失败，则删除之前插入的元数据
		db.Delete(&fileMeta, insertedID)
		log.Fatal("failed to upload file, rolled back metadata insertion: ", err)
	}

	fmt.Println("File uploaded successfully")
}

// uploadFile 是一个模拟文件上传的函数
func uploadFile(meta FileMetadata) error {
	// 这里可以实现实际的文件上传逻辑
	// 如果失败，返回错误
	return fmt.Errorf("mock upload error")
}
```

### 解释：

1. **定义 `FileMetadata` 结构体**：

   ```go
   type FileMetadata struct {
       ID       uint   `gorm:"primaryKey"`
       FileName string `gorm:"not null"`
       FileSize int64  `gorm:"not null"`
   }
   ```

2. **插入文件元数据**：

   ```go
   result := db.Create(&fileMeta)
   if result.Error != nil {
       log.Fatal("failed to insert file metadata: ", result.Error)
   }
   ```

3. **获取插入的ID**：

   ```go
   insertedID := fileMeta.ID
   fmt.Printf("Inserted metadata with ID: %d\n", insertedID)
   ```

4. **上传文件**：

   ```go
   err = uploadFile(fileMeta)
   if err != nil {
       db.Delete(&fileMeta, insertedID)
       log.Fatal("failed to upload file, rolled back metadata insertion: ", err)
   }
   ```

5. **上传文件函数（模拟）**：

   ```go
   func uploadFile(meta FileMetadata) error {
       // 这里可以实现实际的文件上传逻辑
       // 如果失败，返回错误
       return fmt.Errorf("mock upload error")
   }
   ```

这个代码展示了如何使用Gorm处理自增主键的插入和错误处理，并在文件上传失败时进行回滚操作。根据你的实际需求，你可以替换 `uploadFile` 函数的实现。
*/
