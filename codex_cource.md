# Codex 编程入门：配置、用法与 Go 权限后台示例

## 读完可以
- 搭好 Codex 使用环境，知道如何把上下文喂给 Codex。
- 掌握常用对话套路：读代码、做计划、改代码、跑测试、总结。
- 跟着脚本完成一个 Go RBAC 后台（Gin + Gorm + JWT + SQLite），含验证步骤。

## 1. 环境与配置
- 安装前置：Go >= 1.20、Git、curl/httpie、rg（推荐用于搜索）；验证脚本用 `python3` 从 JSON 提取 token/id（也可用 `jq` 或手动复制）。
- 安装 Codex CLI 并确认可用：`codex --version`（具体安装方式按官方指引）。
- 推荐习惯：在项目根目录对话；让 Codex 先看 `go.mod`、`README`；保持分支可写或明确脏改。
- 运行能力：能执行 `go test ./...` 和 `go run .`；如 sandbox 限制执行，手动跑并把结果贴给 Codex。
- 日常命令：`rg --files` 查找文件，`rg keyword` 搜索，`go list` 查看模块，`ls`/`tree` 给 Codex 看结构。

## 2. 基本使用姿势（对话示例）
- 读上下文：`请先看 go.mod 和 main.go（或 cmd/**/main.go），总结依赖和入口。`
- 计划：`基于需求，给 3~6 步计划，包含实现与测试。`
- 改代码：`只改 internal/auth/ 下的文件；改完解释路径和理由。`
- 跑测试：`运行 go test ./...，贴关键失败信息。`（受限时请你自行运行并反馈）
- 迭代：每完成一小块，就让 Codex 总结改动、风险、下一步。
- 需要具体代码：`输出完整 main.go 代码块，并说明放在项目根目录。`

## 3. 实战：用 Codex 写 Go 权限管理后台（RBAC）
### 3.1 目标与技术栈
- 提供用户注册、登录（JWT），角色与权限管理，受保护接口示例。
- 技术：Go + Gin（HTTP）+ Gorm + SQLite（默认方便本地）+ bcrypt（密码）+ JWT。

### 3.2 初始化项目（本地命令）
```bash
mkdir rbac-demo && cd rbac-demo
go mod init rbac-demo
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/sqlite github.com/golang-jwt/jwt/v4 golang.org/x/crypto/bcrypt
```

### 3.3 跟 Codex 的对话脚本
1) 读取上下文  
`这是 Go RBAC 小项目，依赖 Gin/Gorm/SQLite/JWT。先看 go.mod，确认依赖是否齐全，再给出骨架设计。`
2) 生成骨架  
`创建 main.go，包含：DB 初始化（SQLite）、模型 User/Role/Permission、路由注册、基础中间件。`
3) 实现认证与授权  
`补充 /signup /login，使用 bcrypt + JWT；写 authMiddleware 解析 Bearer Token；写 requirePerm(perms...) 校验权限。`
4) 角色与权限管理  
`新增 /roles /permissions 创建接口，/roles/:id/permissions 分配权限，/users/:id/roles 分配角色，权限保护按命名 perm:*。`
5) 受保护示例接口  
`新增 /demo 需要 perm:demo:read；/me 返回当前用户、角色、权限列表。`
6) 测试与运行  
`运行 go run . 后提供 curl 验证脚本；如需测试用例，帮忙加一两个 handler 的单测。`

### 3.4 参考代码（最小可跑示例）
将以下内容保存为 `main.go`（可让 Codex 直接生成文件）：
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var jwtSecret = []byte("change-me")

type User struct {
	ID           uint   `json:"id"`
	Username     string `gorm:"uniqueIndex" json:"username"`
	PasswordHash string `json:"-"`
	Roles        []Role `gorm:"many2many:user_roles;" json:"roles"`
}

type Role struct {
	ID          uint         `json:"id"`
	Name        string       `gorm:"uniqueIndex" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

type Permission struct {
	ID   uint   `json:"id"`
	Name string `gorm:"uniqueIndex" json:"name"`
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type named struct {
	Name string `json:"name"`
}

type assignPermissionReq struct {
	Permission string `json:"permission"`
}

type assignRoleReq struct {
	Role string `json:"role"`
}

func main() {
	db := initDB()
	seedBaseData(db)
	r := setupRouter(db)
	log.Println("server on :8080")
	log.Fatal(r.Run(":8080"))
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("rbac.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&User{}, &Role{}, &Permission{}); err != nil {
		log.Fatal(err)
	}
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.POST("/signup", signUpHandler(db))
	r.POST("/login", loginHandler(db))

	auth := r.Group("/")
	auth.Use(authMiddleware())
	auth.GET("/me", meHandler(db))
	auth.GET("/demo", requirePerm(db, "perm:demo:read"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "you can read demo"})
	})
	auth.POST("/roles", requirePerm(db, "perm:roles:create"), createRoleHandler(db))
	auth.POST("/permissions", requirePerm(db, "perm:permissions:create"), createPermissionHandler(db))
	auth.POST("/roles/:roleID/permissions", requirePerm(db, "perm:roles:assign"), assignPermissionHandler(db))
	auth.POST("/users/:userID/roles", requirePerm(db, "perm:users:assign"), assignRoleToUserHandler(db))

	return r
}

func signUpHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req credentials
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		username := strings.TrimSpace(req.Username)
		if len(req.Password) < 6 || username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		var exists int64
		db.Model(&User{}).Where("username = ?", username).Count(&exists)
		if exists > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
			return
		}
		hash, err := hashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "hash password failed"})
			return
		}
		user := User{Username: username, PasswordHash: hash}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create user failed"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
	}
}

func loginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req credentials
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		username := strings.TrimSpace(req.Username)
		if username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		var user User
		db.Where("username = ?", username).First(&user)
		if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})
		signed, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": signed})
	}
}

func createRoleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req named
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		}
		role := Role{Name: strings.TrimSpace(req.Name)}
		if err := db.FirstOrCreate(&role, "name = ?", role.Name).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
			return
		}
		c.JSON(http.StatusCreated, role)
	}
}

func createPermissionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req named
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		}
		perm := Permission{Name: strings.TrimSpace(req.Name)}
		if err := db.FirstOrCreate(&perm, "name = ?", perm.Name).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
			return
		}
		c.JSON(http.StatusCreated, perm)
	}
}

func assignPermissionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, err := strconv.ParseUint(c.Param("roleID"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
			return
		}
		var req assignPermissionReq
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Permission) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission"})
			return
		}
		var role Role
		if err := db.First(&role, roleID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
		var perm Permission
		if err := db.FirstOrCreate(&perm, "name = ?", strings.TrimSpace(req.Permission)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "permission save failed"})
			return
		}
		if err := db.Model(&role).Association("Permissions").Append(&perm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "assign failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"role": role.Name, "permission": perm.Name})
	}
}

func assignRoleToUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var req assignRoleReq
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Role) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
		var user User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		var role Role
		if err := db.FirstOrCreate(&role, "name = ?", strings.TrimSpace(req.Role)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "role save failed"})
			return
		}
		if err := db.Model(&user).Association("Roles").Append(&role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "assign failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user.Username, "role": role.Name})
	}
}

func meHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := currentUserID(c)
		var user User
		if err := db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"roles":       namesFromRoles(user.Roles),
			"permissions": permissionNames(user.Roles),
		})
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		userID, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		c.Set("userID", uint(userID))
		c.Next()
	}
}

func requirePerm(db *gorm.DB, perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := currentUserID(c)
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing permission"})
			return
		}
		if len(perms) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission not specified"})
			return
		}
		for _, perm := range perms {
			if userHasPermission(db, userID, perm) {
				c.Next()
				return
			}
		}
		var need any = perms
		if len(perms) == 1 {
			need = perms[0]
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing permission", "need": need})
		return
	}
}

func currentUserID(c *gin.Context) uint {
	if v, ok := c.Get("userID"); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

func userHasPermission(db *gorm.DB, userID uint, perm string) bool {
	var count int64
	db.Table("permissions").
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = rp.role_id").
		Joins("JOIN user_roles ur ON ur.role_id = roles.id").
		Where("ur.user_id = ? AND permissions.name = ?", userID, perm).
		Count(&count)
	return count > 0
}

func namesFromRoles(roles []Role) []string {
	out := make([]string, 0, len(roles))
	for _, r := range roles {
		out = append(out, r.Name)
	}
	return out
}

func permissionNames(roles []Role) []string {
	seen := make(map[string]struct{})
	out := []string{}
	for _, r := range roles {
		for _, p := range r.Permissions {
			if _, ok := seen[p.Name]; !ok {
				seen[p.Name] = struct{}{}
				out = append(out, p.Name)
			}
		}
	}
	return out
}

func seedBaseData(db *gorm.DB) {
	perms := []string{
		"perm:roles:create",
		"perm:permissions:create",
		"perm:roles:assign",
		"perm:users:assign",
		"perm:demo:read",
	}
	created := make([]Permission, 0, len(perms))
	for _, name := range perms {
		p := Permission{Name: name}
		db.FirstOrCreate(&p, "name = ?", name)
		created = append(created, p)
	}

	var admin Role
	db.FirstOrCreate(&admin, "name = ?", "admin")
	db.Model(&admin).Association("Permissions").Replace(created)

	var adminUser User
	db.Where("username = ?", "admin").First(&adminUser)
	if adminUser.ID == 0 {
		hash, err := hashPassword("admin123")
		if err != nil {
			log.Fatal(err)
		}
		adminUser = User{Username: "admin", PasswordHash: hash}
		db.Create(&adminUser)
	}
	db.Model(&adminUser).Association("Roles").Replace([]Role{admin})
	log.Println("seeded admin user=admin password=admin123 (change in production)")
}

func hashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
```

### 3.5 运行与验证（手动或让 Codex 协助）
- 启动：`go run .`
- 管理员登录拿 token（用 `python3` 提取；也可用 `jq` 或手动复制 JSON 里的 `token`）：
```bash
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])')
```
- 新建用户并记录 `id`：
```bash
ALICE_ID=$(curl -s -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}' | python3 -c 'import sys, json; print(json.load(sys.stdin)["id"])')
```
- 创建角色并记录 `id`，分配权限，再把角色分配给用户：
```bash
READER_ROLE_ID=$(curl -s -X POST http://localhost:8080/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"reader"}' | python3 -c 'import sys, json; print(json.load(sys.stdin)["id"])')

curl -s -X POST "http://localhost:8080/roles/$READER_ROLE_ID/permissions" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"permission":"perm:demo:read"}'

curl -s -X POST "http://localhost:8080/users/$ALICE_ID/roles" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"role":"reader"}'
```
- 用户登录并访问受保护接口：
```bash
USER_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}' | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])')

curl -s -H "Authorization: Bearer $USER_TOKEN" http://localhost:8080/demo
curl -s -H "Authorization: Bearer $USER_TOKEN" http://localhost:8080/me
```
- 若想新增权限/角色，只需用管理员 token 调用对应接口；生产环境请更换 `jwtSecret`、admin 密码、数据库。

### 3.6 可扩展练习
- 换数据库驱动（Postgres/MySQL），抽离配置文件或环境变量。
- 增加刷新 token、密码找回、操作审计日志。
- 添加 Gin 中间件做请求日志、请求限流或跨域。
- 补充集成测试：用 httptest 对 `/login`、`/demo` 做权限回归。

## 4. 常见问题
- 依赖缺失：运行 `go mod tidy`；若失败，把错误贴给 Codex。
- 端口占用：修改 `r.Run(":8080")` 的端口或释放占用进程。
- 权限不足：确认使用管理员 token；或在 seed 阶段新增自己的初始用户。
- SQLite 文件权限：确保运行目录可写，或改用内存 `sqlite.Open(":memory:")`。
