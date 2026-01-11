# Codex 编程入门（VS Code 插件场景）：配置、用法与 Go 权限后台示例（MySQL 8）

## 读完可以
- 在 VS Code 里用 Codex 插件搭好开发环境，知道如何把上下文喂给 Codex。
- 掌握插件场景的对话套路：读代码、做计划、改代码、跑测试、在 VS Code 看 diff、总结交付。
- 从新建一个空目录开始，跟着脚本完成一个 Go RBAC 后台（Gin + Gorm + JWT + MySQL 8），含验证步骤。

## 1. 环境与配置
- 安装前置：VS Code、Codex 插件、Go >= 1.20、MySQL 8（本地已装）、Git、curl/httpie、rg（推荐用于搜索）。
- VS Code 建议：打开“工作区/文件夹”作为项目根；在集成终端跑命令；用 Source Control 面板看 diff/回滚。
- 验证脚本：优先用 `python3` 从 JSON 提取 token/id（也可用 `jq` 或手动复制）。
- 推荐习惯：新建目录后先 `git init`；让 Codex 先读 `AGENTS.md`（若有）、`README.md`、`go.mod`、目录结构；保持改动小步可审阅。
- 推荐加一份项目指引：在项目根写 `AGENTS.md`，把“默认中文沟通/文档中文、修改范围、验证命令”等规则固化下来，减少反复强调。
- 日常命令：`rg --files` 查找文件，`rg keyword` 搜索，`go list` 查看模块，`ls`/`tree` 给 Codex 看结构。

## 2. VS Code 插件的基本用法（对话示例）
> 插件场景的关键是“上下文”：尽量让 Codex 读到你当前工作区的关键文件；改动后用 VS Code 的 diff 审阅再继续下一步。

- 读上下文（新建目录）：`这是一个新建工作区。请先看目录结构（ls/tree）+ README/go.mod（若有），总结现状并给出下一步建议。`
- 读当前文件/选区：`请只基于我当前打开文件（以及选中的片段）分析问题并给修改建议。`
- 计划：`基于需求给 3~6 步计划（含实现与验证），每一步写清要改哪些文件、要跑什么命令。`
- 控制修改范围：`这一步只允许改 main.go（或指定文件列表），不要动其它文件。`
- 跑验证：`请告诉我在 VS Code 终端要跑哪些命令验证（例如 go test ./...），以及预期输出。`
- 迭代：每完成一小块，就让 Codex 总结“改了什么/为什么/如何验证/下一步”。

### 2.1 推荐：计划 → 确认 → 分步执行（更像“计划模式”）
把下面模板直接粘到 VS Code 插件对话里用（适合你这个 Go RBAC 案例，能显著减少跑偏和大改动）。

1) 先只产出计划（不写代码）：
```text
先不要写代码。请先阅读当前工作区关键文件（AGENTS.md/README/go.mod/目录结构），然后给出 3~6 步计划。
要求：每一步写清（要改哪些文件、要新增哪些接口/结构、要跑什么命令验证、预期结果）。
我回复“开始第 1 步”后你再实现。
```

2) 开始某一步（带范围约束）：
```text
开始第 1 步。这一步只允许改这些文件：<填文件列表>。
实现完成后只交付：改动摘要 + VS Code diff 审阅要点 + 最小验证命令（含预期输出）。
如果发现计划需要调整，先停下来说明原因并给出更新后的计划，再等我确认。
```

3) 每步收尾复盘（强制可复现）：
```text
请用 5 行内总结：这一步改了什么/为什么/如何验证/可能风险/下一步做什么。
```

### 2.2 一组高性价比“提示词模板”（不靠运气，减少返工）
下面这些不需要每次全用：按场景挑一条粘贴就行。

**只诊断不动手（先定位再修改）**
```text
先不要改代码。请只基于我当前打开文件/我粘贴的日志，定位问题可能原因，给 3 个最可能的根因与各自的验证方法。
确认根因后我再让你改。
```

**先复述需求与约束（防止跑偏）**
```text
在开始前请先复述：需求目标、明确约束（不能改哪些文件/不能引入哪些依赖/是否允许联网）、以及你准备修改的文件列表。
如果有任何不确定点，先提 1–3 个问题再继续。
```

**最小改动 + 可审阅（避免大重构）**
```text
请做“最小且合理”的改动：优先补缺/修错，不做无关重构、不改命名、不移动文件。
改完请给：改动文件列表 + 每个文件一句话说明 + 需要我在 VS Code diff 重点看的 3 个点。
```

**强制范围约束（只改指定文件）**
```text
这一步只允许改这些文件：<填文件列表>。如果必须新增文件，请先说明原因并等我确认。
```

**交付必须带可复现验证**
```text
改完请给出在 VS Code 集成终端可直接复制的验证命令（包含预期输出/成功标准）。
如果需要我先准备环境（MySQL/环境变量），请列清单。
```

**提供回滚/撤销（方便你大胆试）**
```text
如果这一步改坏了，请给出回滚方式（例如用 git restore/checkout 恢复哪些文件），并说明回滚会影响哪些改动。
```

**出现报错时如何喂上下文（提高一次命中率）**
```text
我将粘贴：报错日志 + 复现步骤 + 相关文件片段。你先给出最小复现与最小修复方案，再给出具体修改。
```

**需要查资料/用外部工具时先征求确认**
```text
如果你需要查官方文档/仓库/Issue，先说明你准备查哪里、要解决什么问题，再等我确认是否允许（或是否启用 MCP）。
```

## 3. 实战：用 Codex 写 Go 权限管理后台（RBAC）
### 3.1 目标与技术栈
- 提供用户注册、登录（JWT），角色与权限管理，受保护接口示例。
- 技术：Go + Gin（HTTP）+ Gorm + MySQL 8 + bcrypt（密码）+ JWT。

### 3.2 初始化项目（本地命令）
在 VS Code 中：新建一个空文件夹并“打开文件夹”作为工作区，然后在集成终端执行：
```bash
mkdir rbac-demo && cd rbac-demo
git init

# 建议：用环境变量存 MySQL DSN，把 .env 加进 .gitignore
printf ".env\n" > .gitignore

go mod init rbac-demo
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/mysql github.com/golang-jwt/jwt/v4 golang.org/x/crypto/bcrypt

# 可选：先提交一个“空项目”基线，方便后续让 Codex 小步改动
git add .
git commit -m "chore: init project"
```

### 3.2.1 准备 MySQL（本地已安装）
先确保 MySQL 8 已启动，并创建数据库（示例叫 `rbac_demo`）：
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS rbac_demo DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
```
然后在项目根目录创建 `.env`（示例 DSN；按你本机账号密码修改）：
```bash
cat > .env <<'EOF'
MYSQL_DSN=root:your_password@tcp(127.0.0.1:3306)/rbac_demo?charset=utf8mb4&parseTime=True&loc=Local
EOF
```

### 3.2.2 （推荐）生成 `AGENTS.md`：默认中文沟通 + 中文文档
很多 AI 编码助手会自动读取项目根的 `AGENTS.md`，把它当作“项目内的固定规则”。建议你在新建目录后就生成一份（按需再改）：
```bash
cat > AGENTS.md <<'EOF'
<INSTRUCTIONS>
# Project Collaboration Guide (Flexible)

## Language
- Default to communicating and writing docs in Chinese, unless I explicitly request English/bilingual.

## Scope
- Unless explicitly requested, prefer modifying only files under the current project directory; avoid large rewrites and keep changes in small steps.
- Unless explicitly requested, do not automatically `git commit`/`git push`.

## Safety
- Do not write sensitive information into the repo (passwords, tokens, private keys, MySQL DSNs, etc.); prefer injecting via environment variables/`.env`, and make sure `.env` is included in `.gitignore`.

## Delivery and Verification
- If you touch commands/code snippets, provide a minimal verification path (e.g., `go test ./...`, `go run .` + a `curl` check).
- In output, prioritize executable steps/conclusions first, then only the necessary explanation; when referencing files, include paths (and line numbers when needed).
</INSTRUCTIONS>
EOF
```

### 3.2.3 （可选）配置 MCP：让插件接入外部工具/数据源
MCP（Model Context Protocol）可以把“外部系统的能力”以标准方式提供给 Codex，例如：文件系统、GitHub、浏览器自动化等（取决于你接入的 MCP Server，以及你使用的 Codex 插件/运行方式是否支持 MCP）。

如果你用的是 Codex CLI（而不是/不只是 VS Code 插件），可以用 `codex mcp ...` 管理 MCP Server（实验性）：`add/list/get/remove`。

先用一句话理解 MCP：
- **MCP Server**：一个“工具提供者”进程（本地起一个命令，或远程 HTTP 服务）。
- **Codex（客户端）**：把 Server 提供的能力注册成可调用的工具（tools），并在对话中按需调用。

示例：接入一个“文件系统 MCP Server”（Node 版），把某个目录暴露给 MCP（会写入 `~/.codex/config.toml`）：
```bash
codex mcp add filesystem -- npx -y @modelcontextprotocol/server-filesystem "$(pwd)"
codex mcp list
codex mcp get filesystem
codex mcp remove filesystem
```

#### 一个真正“跑得出来”的 MCP 最小例子（推荐按这个做一遍）
上面用 `$(pwd)` 作为目录时，你可能感觉“没变化”（因为 Codex 本来就能读你工作区里的文件）。所以这里故意用一个**工作区外**的目录做演示：在 `workspace-write` 沙箱下，Codex 默认不该直接读取 `$HOME` 下的工作区外目录（或其它工作区外路径），只有通过你接入的 MCP 工具才有机会读到。

前置：本机有 `node`/`npx`（Node 20+ 更稳）：
```bash
node -v
npx -v
```

1) 准备一个工作区外的目录与文件：
```bash
mkdir -p "$HOME/mcp-demo"
cat > "$HOME/mcp-demo/hello.txt" <<'EOF'
hello from mcp
EOF
```
如果你只是临时演示，用 `/tmp/mcp-demo` 也可以（可能会被系统清理）。

2) 注册一个 MCP Server（stdio 方式）：
```bash
codex mcp add demo-fs -- npx -y @modelcontextprotocol/server-filesystem "$HOME/mcp-demo"
codex mcp get demo-fs
```

3) 让 Codex **明确使用**这个 MCP 工具来读文件（示例用非交互 `codex exec`；你也可以直接开交互 `codex` 再粘贴同样的提示词）：
```bash
codex exec -s workspace-write "请使用 MCP 工具 demo-fs：列出根目录有哪些文件，然后读取 hello.txt，把文件内容原样打印出来。最后说明你调用了哪个工具完成的。"
```

4) 用完就移除（避免长期暴露目录）：
```bash
codex mcp remove demo-fs
```

> 你也可以不通过 `codex mcp add`，直接手动在 `~/.codex/config.toml` 配置（格式大致如下）：
```toml
[mcp_servers.demo-fs]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-filesystem", "/home/your_name/mcp-demo"]
```
把 `/home/your_name` 替换成你机器上的实际 HOME 路径（可用 `echo "$HOME"` 查看）。

注意：MCP Server 可能拥有较高权限（读写文件、访问网络/第三方 API 等）。只添加你信任的 Server，并优先用环境变量/`.env` 传递凭证，避免写进仓库。

### 3.3 跟 Codex 的对话脚本
1) 读取上下文  
`这是一个从零开始的新目录（已 git init + go mod init）。目标是做 Go RBAC 小项目（Gin/Gorm/MySQL8/JWT）。请先看目录结构、go.mod、README（若有），确认依赖与入口建议，再给出骨架设计与下一步计划。`
2) 生成骨架  
`创建 main.go，包含：DB 初始化（MySQL）、模型 User/Role/Permission、路由注册、基础中间件。`
3) 实现认证与授权  
`补充 /signup /login，使用 bcrypt + JWT；写 authMiddleware 解析 Bearer Token；写 requirePerm(perms...) 校验权限。`
4) 角色与权限管理  
`新增 /roles /permissions 创建接口，/roles/:id/permissions 分配权限，/users/:id/roles 分配角色，权限保护按命名 perm:*。`
5) 受保护示例接口  
`新增 /demo 需要 perm:demo:read；/me 返回当前用户、角色、权限列表。`
6) 测试与运行  
`运行 go run . 后提供 curl 验证脚本；如需测试用例，帮忙加一两个 handler 的单测。`

### 3.4 参考代码（最小可跑示例）
将以下内容保存为 `main.go`（更推荐让 Codex 按 3.3 小步生成，下面可作为对照）：
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
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
	dsn := getenv("MYSQL_DSN", "root:password@tcp(127.0.0.1:3306)/rbac_demo?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
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
- MySQL 连接失败：确认 MySQL 服务已启动、账号有权限、DSN 正确（尤其是 `parseTime=True`），以及数据库（如 `rbac_demo`）已创建。
