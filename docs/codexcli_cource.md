# Codex CLI 编程入门：配置、用法与 Go 权限后台示例（MySQL 8）

> 本文面向 **Codex CLI**（终端场景）。如果你要看 VS Code 插件场景，请看 `docs/codex_cource.md`。

## 读完可以
- 在终端用 Codex CLI 跑通：读上下文 → 出计划 → 分步改动 → 终端验证 → `git diff` 审阅。
- 理解 CLI 里最关键的 3 个概念：**交互模式**、**非交互 exec**、**沙箱与审批**。
- 跟着脚本从零做一个 Go RBAC 小后台（Gin + Gorm + JWT + MySQL 8），含 curl 验证与最小回滚方式。
- 知道怎么接入一个“真的有用”的 MCP（文件系统 / GitHub），并在提示词里让它“明确使用”。

---

## 1. 环境与检查

### 1.1 必备工具
- `git`：用于审阅/回滚（强烈建议）。
- `rg`（ripgrep，可选但强烈建议）：用于快速搜索。
- `python3`：用于从 JSON 输出里提取字段（避免依赖 `jq`）。
- Go 示例需要：Go >= 1.20、MySQL 8、`curl`。

先跑一遍自检（有啥缺啥补）：
```bash
codex --version
git --version
python3 --version
rg --version || true
go version || true
mysql --version || true
curl --version
```

### 1.2 安全与凭证
- 不要把任何密码、token、私钥、DSN 写进仓库。
- 需要凭证时用环境变量或 `.env`（并确保 `.env` 已加入 `.gitignore`）。

---

## 2. CLI 里的 3 个核心概念

### 2.1 交互模式 vs 非交互模式
- **交互模式**：`codex`  
  适合探索、反复迭代、边看边改。
- **非交互模式**：`codex exec "..."`  
  适合把一次任务“固定成一条可复现命令”，例如：让它做代码审阅、生成计划、输出一段可复制的命令序列等。

示例：
```bash
codex
```
```bash
codex exec "请读取当前目录结构并总结项目入口与下一步建议。"
```

### 2.2 沙箱（sandbox）与权限边界
你可能会看到/使用类似参数（不同环境配置略有差异）：
- `-s workspace-write`：只允许读写工作区（最常用、也最安全）。
- 有些环境还支持“只读”沙箱（例如 `read-only`）：只允许读取，禁止改文件（适合“先分析/先出计划”）。

建议默认用 `workspace-write`，但在“先看不改”的阶段强制只读更稳：
```bash
codex exec -s workspace-write "先不要改文件。请读目录结构与 README，总结现状并给 5 步计划。"
```

### 2.3 审批（approval）与网络访问
很多环境会对“可能写入/可能联网/可能访问工作区外路径”的操作要求审批。实战原则：
- 需要联网（下载依赖、拉取 MCP Server、访问第三方 API）之前，先确认环境是否允许。
- 不确定就让 Codex 先说明“要做什么、影响什么”，再决定是否放行。

---

## 3. 推荐工作流：计划 → 确认 → 分步执行（CLI 版）

### 3.1 第一句话：让它先读规约
如果你的项目根目录有 `AGENTS.md`，建议在每个新会话第一句提醒：
```text
请先阅读工作区根目录的 AGENTS.md，并按其中工作流执行。
```

### 3.2 先出计划（不改文件）
```bash
codex exec -s workspace-write "先不要改文件。请读取目录结构（ls/tree）+ README/go.mod（如有）+ AGENTS.md（如有），给 3–6 步计划：每步写清要改哪些文件、怎么验证、预期结果。"
```

### 3.3 执行某一步（强制范围约束 + 可审阅）
把“范围”写死，避免大改：
```bash
codex exec -s workspace-write "开始第 1 步。这一步只允许改：main.go 和 internal/auth/*。改完请只交付：改动摘要 + git diff 审阅要点 + 最小验证命令（含预期输出）。如果计划需要调整，先停下说明原因并等我确认。"
```

### 3.4 每步都留回滚手刹
你可以要求它在交付里带回滚指令。最通用的是：
```bash
git status
git diff
git restore -SW .
```

---

## 4. 一组高性价比提示词模板（CLI 版）

### 4.1 只诊断不动手（先定位再改）
```bash
codex exec -s workspace-write "先不要改文件。请只基于我粘贴的日志/当前文件内容，给 3 个最可能根因 + 各自的验证方法；等我确认根因后再改。"
```

### 4.2 先复述需求与约束（防跑偏）
```bash
codex exec -s workspace-write "开始前请先复述：需求目标、约束（不能改哪些文件/不能引入哪些依赖/是否允许联网）、以及你准备修改的文件列表。若不确定，先提 1–3 个问题。"
```

### 4.3 最小改动 + 可审阅（避免重构）
```bash
codex exec -s workspace-write "请做最小且合理的改动：优先补缺/修错，不做无关重构、不改命名、不移动文件。改完给：改动文件列表 + 每个文件一句话说明 + 需要我重点看 diff 的 3 个点。"
```

### 4.4 交付必须带可复现验证
```bash
codex exec -s workspace-write "改完请给出可直接复制的验证命令（含预期输出/成功标准）。若需要我先准备环境（MySQL/环境变量），请列清单。"
```

---

## 5. 实战：用 Codex CLI 做 Go RBAC 后台（Gin/Gorm/JWT/MySQL 8）

> 这一节的重点不是“抄代码”，而是把 Codex CLI 的工作流跑通：先计划，再分步生成/修改，最后用命令验证。

### 5.1 初始化项目（本地命令）
```bash
mkdir rbac-demo && cd rbac-demo
git init

printf ".env\n" > .gitignore
go mod init rbac-demo
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/mysql github.com/golang-jwt/jwt/v4 golang.org/x/crypto/bcrypt
```

### 5.2 准备 MySQL（本地已安装）
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS rbac_demo DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
```

创建 `.env`（示例 DSN；按你本机账号密码修改）：
```bash
cat > .env <<'EOF'
MYSQL_DSN=root:your_password@tcp(127.0.0.1:3306)/rbac_demo?charset=utf8mb4&parseTime=True&loc=Local
EOF
```

### 5.3 （推荐）写一份项目 `AGENTS.md`
```bash
cat > AGENTS.md <<'EOF'
<INSTRUCTIONS>
# Project Collaboration Guide (CLI Scenario)

## Language
- Default to communicating and writing docs in Chinese, unless I explicitly request English/bilingual.

## Workflow
- Default to "plan → confirmation → step-by-step execution": propose a 3–6 step plan first (each step: which files to change + how to verify + expected result), and only implement after I say "start step N".
- Each step’s delivery must include: change summary, review points (diff/risk), minimal verification commands (with expected output), and rollback suggestions.

## Scope
- Unless explicitly requested, make only the smallest reasonable change; no unrelated refactors, no renaming, no moving files.
- If I specify a file list, modify only those files; if you need to add new files, explain why and wait for my confirmation.

## Safety
- Do not write sensitive information into the repo (passwords, tokens, private keys, MySQL DSNs, etc.); use environment variables/`.env`, and make sure `.env` is included in `.gitignore`.
</INSTRUCTIONS>
EOF
```

### 5.4 用 Codex CLI 分步实现（建议用交互模式）
启动交互会话：
```bash
codex
```

推荐按这个节奏对话（每步都让它给验证命令）：
1) 读取上下文并出计划  
```text
这是一个从零开始的新目录（已 git init + go mod init）。目标是做 Go RBAC 小项目（Gin/Gorm/MySQL8/JWT）。
请先读目录结构、go.mod、AGENTS.md，并给 5 步计划（每步：要改哪些文件 + 如何验证 + 预期结果）。
```

2) 生成骨架  
```text
开始第 1 步：创建 main.go（以及你认为最小必要的文件），包含：DB 初始化（MySQL）、模型 User/Role/Permission、路由注册、基础中间件。
只做最小可跑骨架；跑通后给 go run 与 /health 的 curl 验证。
```

3) 实现认证与授权  
```text
开始第 2 步：补充 /signup /login（bcrypt + JWT），写 authMiddleware 解析 Bearer Token，写 requirePerm(perms...) 校验权限。
```

4) 角色与权限管理  
```text
开始第 3 步：新增 /roles /permissions 创建接口，/roles/:id/permissions 分配权限，/users/:id/roles 分配角色，权限命名示例 perm:*。
```

5) 受保护示例接口 + 验证脚本  
```text
开始第 4 步：新增 /demo 需要 perm:demo:read；/me 返回当前用户、角色、权限列表。
给一组可复制的 curl 验证脚本，并用 python3 从 JSON 提取 token/id。
```

> 如果你更喜欢“一条命令做完一件事”，也可以把每一步改成 `codex exec -s workspace-write "..."` 的形式。

### 5.5 最小验证（curl + python3 提取字段）
管理员登录拿 token（也可用 `jq`）：
```bash
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -c 'import sys, json; print(json.load(sys.stdin)["token"])')
```

后续验证建议用同样方式提取 `id/token`，避免手动复制。

---

## 6. MCP：推荐一个“真有用”的接入方式（CLI 场景）

### 6.1 文件系统 MCP（最小可跑示例）
用途：让 Codex 通过 MCP 工具访问“工作区外”的目录（演示权限边界与工具调用）。

前置：本机有 `node`/`npx`（Node 20+ 更稳）：
```bash
node -v
npx -v
```

1) 准备一个工作区外目录：
```bash
mkdir -p "$HOME/mcp-demo"
cat > "$HOME/mcp-demo/hello.txt" <<'EOF'
hello from mcp
EOF
```

2) 注册 MCP Server：
```bash
codex mcp add demo-fs -- npx -y @modelcontextprotocol/server-filesystem "$HOME/mcp-demo"
codex mcp list
codex mcp get demo-fs
```

3) 让 Codex 明确调用这个工具：
```bash
codex exec -s workspace-write "请使用 MCP 工具 demo-fs：列出根目录有哪些文件，然后读取 hello.txt，把文件内容原样打印出来。最后说明你调用了哪个工具完成的。"
```

4) 用完移除：
```bash
codex mcp remove demo-fs
```

> 也可以手动在 `~/.codex/config.toml` 配置；路径里把 `/home/your_name` 换成 `echo \"$HOME\"` 的实际输出。

### 6.2 GitHub MCP（对 Go RBAC 很实用）
用途：查 `gin-gonic/gin`、`go-gorm/gorm`、`golang-jwt/jwt` 的权威示例/Issue/源码，减少“靠猜”。

建议 token 用环境变量提供（不要写进仓库）：
```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN"
```

添加（可能需要联网下载）：
```bash
codex mcp add github -- npx -y @modelcontextprotocol/server-github
```

使用示例：
```bash
codex exec -s workspace-write "请使用 MCP 工具 github：在 go-gorm/gorm 的资料中找到 many2many + Preload 的推荐写法与常见坑，并给出适用于 User-Role-Permission 的最小示例。"
```

---

## 7. 常见问题

### 7.1 Codex 不改文件/改了我不放心
- 让它输出“改动摘要 + 需要你重点看 diff 的点”，然后你用：
```bash
git status
git diff
```
- 不满意就回滚：
```bash
git restore -SW .
```

### 7.2 网络/依赖安装受限
- 先让 Codex 给出“离线可执行”的步骤或替代方案。
- 需要联网（比如安装 MCP Server、下载依赖）前先确认环境策略是否允许。

### 7.3 连接 MySQL 失败
- 检查 MySQL 是否启动、账号权限、DSN 是否正确（尤其 `parseTime=True`），数据库是否已创建。
