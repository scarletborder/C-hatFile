# chatFile Backend
从zb_sccrystal中的同名插件中单独拆分出后端逻辑整合成的后端

将为了适合通用后端将模糊bot环境中的概念



由于您使用的是Gin框架，我将使用Gin的特性来改写之前的代码示例。这将包括创建登录端点、生成JWT Token以及中间件来验证请求中的Token。

### 1. 生成JWT Token的函数

```go
import (
    "github.com/dgrijalva/jwt-go"
    "time"
)

// 生成JWT Token
func generateToken(username string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username, // 将用户名编入Token
        "exp": time.Now().Add(time.Hour * 72).Unix(), // Token有效期
    })

    // 用一个安全的密钥签名Token
    tokenString, err := token.SignedString([]byte("your-secret-key"))
    return tokenString, err
}
```

### 2. 登录路由处理器

```go
import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func loginHandler(c *gin.Context) {
    var loginInfo struct {
        Username string `json:"username"`
        Password string `json:"password"` // 假设前端发送的是加密后的密码
    }

    if err := c.ShouldBindJSON(&loginInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // 模拟用户验证过程
    if loginInfo.Username != "expectedUsername" || loginInfo.Password != "expectedPassword" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
        return
    }

    tokenString, err := generateToken(loginInfo.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
```

### 3. Token验证中间件

```go
import (
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

func TokenAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, gin.ErrorResponse{Err: "Unexpected signing method"}
            }
            return []byte("your-secret-key"), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Set("username", claims["username"])
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 4. 使用中间件和处理资源请求的路由

```go
func main() {
    router := gin.Default()

    router.POST("/login", loginHandler)
    router.GET("/resource", TokenAuthMiddleware(), func(c *gin.Context) {
        username := c.MustGet("username").(string)
        c.JSON(http.StatusOK, gin.H{"user": username, "data": "Here is your protected resource"})
    })

    router.Run(":8080")
}
```

在这个代码中，我设置了一个Gin路由器，一个用于登录的POST端点，和一个使用Token验证中间件保护的GET端点。这种结构可以很好地适用于基于Token的用户身份验证和资源请求处理。



# 原来的插件描述
为群聊环境设计的一款允许文件上传、下载、查询的zerobot(fork version)插件

## 前言

笔者先前已经开始着手编写该插件，但是原先的插件功能较少仅有上传、查询和下载功能，并且没有很好的权限控制，代码结构也较为混乱。由此笔者借助此次实验任务机会重构代码。

## 需求分析

### 特性
- `admin`, `管理员`, `匿名用户`三种用户权限能力分隔
- 支持标签上传，标题标签查询文件
- 方便的获取文件预签名下载链接
- 匿名用户查询个人上传的所有文件
- 匿名用户删除个人上传的文件
- 管理员删除文件
- admin授权管理员
- 一些统计信息的查看

### 数据字典

1. **用户类型**
   - **类型**：枚举
   - **可能的值**：匿名用户、管理员、admin
   - **描述**：定义用户的权限级别。

2. **文件**
   - **类型**：复合数据类型
   - **属性**：
     - 文件ID：字符串，唯一标识符
     - 文件标题：字符串
     - 文件大小：int64
     - 文件标签：标签数组
     - 上传者用户ID：视bot环境而适配，如qq是qq号
     - 上传时间：golang time.time类型
   - **描述**：在系统中上传和查询的文件详情。

3. **管理员密钥**
   - **类型**：字符串
   - **描述**：允许admin生成的特定密钥，用于给管理员授权高级操作权限。

4. **查询类型**
   - **类型**：复合数据类型
   - **属性**：
     - 文件标签：标签数组
     - 文件标题：字符串
   - **描述**：用户用于查询文件的方法，至少一项不为空。

5. **查询结果**
   - **类型**：字符串数组
   - **描述**：根据查询类型返回的文件ID列表。

6. **下载链接请求**
    - **类型**：字符串
    - **描述**：匿名用户提交的申请下载的文件ID。

7. **删除请求**
   - **类型**：字符串
   - **描述**：管理员提交的用于删除特定文件ID的请求。

### 数据流

1. **文件上传数据流**
   - **包含**：文件标题、文件标签
   - **描述**：用户上传文件时提供的数据。

2. **文件查询请求**
   - **包含**：查询类型、查询关键词（标题或标签）
   - **描述**：用户用于搜索文件的请求数据。

3. **管理员授权请求**
   - **类型**：操作类型（授权/删除授权）、管理员ID
   - **描述**：admin处理管理员权限的请求。

## 实验内容

### 系统设计与数据库设计

#### 系统设计

- **View/API**：`lagrange.onebot`实现qq-nt协议通信，`zerobot`提供机器人框架,`zb_sccrystal/Global`提供更多有关数据库和AWS S3交互以及拓展qq通信的API
- **Application**: `zb_sccrystal/Plugin/chatFile`包
  - 提供本次实验用到的所有功能入口，和对应逻辑。
- **Domain**： 
  - 上下文用到的数据结构定义在`Global/Storage/model.go`中，对`zb_sccrystal`统一设计的一种文件对象。
  - `Global/Storage/utils.go`提供了一种通过FNV算法生成minio对象存储的唯一的文件字符串索引
- Infrastructure:
  - Storage:
    - 数据库：Mysql
    - 文件管理：Minio



#### 数据库设计



#### 更多细节数据库设计

##### 权限控制

在用户权限管理方面，主要是以用户层的判断为主。为了体现用户口令哈希存储的实验要求，在应用层保证user名只和对应即时通讯软件的用户id(如qq号)相关下，设置用户权限访问控制。

例如，admin(实验环境是继承于`Global/Config/config.toml`的super_ids配置，即qq号in [1581822568])可以享有所有的操作权限。