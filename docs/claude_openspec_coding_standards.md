# Claude + OpenSpec 编程规范

> 本规范定义了使用 Claude Code 结合 OpenSpec 进行规范驱动开发的标准流程、话术模板和交付标准。

---

## 目录

1. [核心原则](#1-核心原则)
2. [OpenSpec 安装与配置](#2-openspec-安装与配置)
3. [OpenSpec 使用详解](#3-openspec-使用详解)
4. [Spec 编写规范](#4-spec-编写规范)
5. [Claude 协作话术](#5-claude-协作话术)
6. [代码实现规范](#6-代码实现规范)
7. [验收与验证规范](#7-验收与验证规范)
8. [常见反模式](#8-常见反模式)

---

## 1. 核心原则

### 1.1 规范驱动 (Spec-Driven)

- **先 Spec，后代码**：任何功能开发必须先有明确的 Spec
- **Spec 是单一事实来源**：代码实现必须与 Spec 保持一致
- **可变更**：Spec 可以演进，但必须通过显式的 Change 流程

### 1.2 小步交付 (Incremental)

- **每步可验证**：每个实现步骤必须能独立验证
- **限制变更范围**：一次只改相关文件，避免大规模重构
- **可回滚**：每步改动都是可审阅、可回滚的

### 1.3 对话驱动 (Conversation-Driven)

- **澄清优先**：动手前先提 1-3 个关键问题
- **计划确认**：给出计划后等待用户确认再执行
- **透明交付**：每次交付包含改动摘要、审阅点、验证命令

---

## 2. OpenSpec 安装与配置

### 2.1 是否需要安装 OpenSpec CLI？

| 使用场景 | 是否需要 OpenSpec CLI | 说明 |
|---------|----------------------|------|
| **仅使用 Claude Code** | ❌ 不需要 | Claude Code 原生支持 OpenSpec，可用 `/openspec-*` 命令 |
| **需要独立管理 Spec** | ✅ 需要 | 在终端直接运行 `openspec` 命令 |
| **团队协作/CI 集成** | ✅ 需要 | 需要独立的 Spec 校验和归档 |

**推荐**：如果是个人开发且主要使用 Claude Code，可以不安装 OpenSpec CLI。

### 2.2 安装 OpenSpec CLI

#### 前置条件

检查是否已安装 Node.js（需要 >= 18.x）：

```bash
# 检查 Node.js 版本
node --version
```

如果未安装或版本过低，请先安装 Node.js：
- **macOS**: `brew install node`
- **Ubuntu/Debian**: `sudo apt update && sudo apt install nodejs npm`
- **Windows**: 从 [nodejs.org](https://nodejs.org/) 下载安装

#### 安装步骤

```bash
# 1. 全局安装 OpenSpec CLI
npm install -g @fission-ai/openspec@latest

# 2. 验证安装
openspec --version

# 3. 查看帮助信息
openspec --help
```

#### 验证安装成功的输出

```bash
$ openspec --version
openspec/0.x.x

$ openspec --help
Usage: openspec [options] [command]

Options:
  -v, --version          output the version number
  -h, --help             display help for command

Commands:
  init                   初始化 OpenSpec 项目
  list                   列出所有变更
  show <change>          查看变更详情
  validate <change>      校验变更结构
  archive <change>       归档变更
  ...
```

### 2.3 初始化 OpenSpec 项目

#### 在现有项目中初始化

```bash
# 1. 进入项目根目录
cd /path/to/your/project

# 2. 初始化 OpenSpec
openspec init
```

初始化后会生成以下结构：

```
project/
├── openspec/
│   ├── specs/          # 规格说明书目录
│   ├── changes/        # 变更请求目录
│   └── config.json     # OpenSpec 配置文件
├── AGENTS.md           # AI 协作指南（自动生成）
└── .openspecrc         # OpenSpec 配置（可选）
```

#### 验证初始化

```bash
# 检查目录是否创建
ls -la openspec/

# 查看配置
cat openspec/config.json
```

### 2.4 创建 AGENTS.md（如未自动生成）

如果 `openspec init` 没有自动生成 `AGENTS.md`，手动创建：

```bash
# 在项目根目录创建
cat > AGENTS.md << 'EOF'
<INSTRUCTIONS>
# 项目协作约定

## 仓库定位
- 本仓库使用 OpenSpec 进行规范驱动开发
- 所有需求变更必须先写 Spec，再实现

## 默认工作流
1. 阅读相关 Spec 文件（openspec/specs/）
2. 提出 1-3 个澄清问题
3. 给出 4-6 步实现计划
4. 小步实现，每步可验证

## 范围与安全
- 默认只修改 Spec 中指定的文件
- 不引入新依赖，除非明确批准
- 不自动 git commit/git push
- 敏感信息使用环境变量或占位符

## 验证与交付
- 每步交付包含：改动摘要、审阅点、验证命令
- 验收标准必须与 Spec 一致
</INSTRUCTIONS>
EOF
```

---

## 3. OpenSpec 使用详解

### 3.1 OpenSpec 工作流全景图

```
┌────────────────────────────────────────────────────────────────────────┐
│                         OpenSpec 完整工作流                             │
├────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐             │
│  │  1.写   │───▶│  2.创建  │───▶│  3.澄清  │───▶│  4.计划  │             │
│  │  Spec   │    │  Change │    │  疑问   │    │         │             │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘             │
│       │              │              │              │                   │
│       ▼              ▼              ▼              ▼                   │
│  openspec/    openspec/    与 Claude   与 Claude                     │
│  specs/       changes/     对话确认     对话确认                       │
│                                                                         │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐                             │
│  │  5.小步  │◀───│  6.验证  │◀───│  7.实现  │                             │
│  │  交付   │    │  测试   │    │  Tasks  │                             │
│  └─────────┘    └─────────┘    └─────────┘                             │
│       │              │              │                                  │
│       ▼              ▼              ▼                                  │
│  文档总结      终端验证      限制文件范围                             │
│                                                                         │
│  ┌─────────┐                                                             │
│  │  8.归档  │                                                             │
│  └─────────┘                                                             │
│       │                                                                  │
│       ▼                                                                  │
│  openspec archive                                                        │
│  合并到 specs/                                                           │
│                                                                         │
└────────────────────────────────────────────────────────────────────────┘
```

### 3.2 核心 CLI 命令详解

#### 命令一览表

| 命令 | 简写 | 用途 | 输出 |
|------|------|------|------|
| `openspec init` | - | 初始化项目 | 创建 `openspec/` 目录 |
| `openspec list` | `openspec ls` | 列出所有变更 | 变更名称、状态、创建时间 |
| `openspec show <change>` | `openspec cat` | 查看变更详情 | proposal、tasks、spec delta |
| `openspec validate <change>` | `openspec check` | 校验变更结构 | 通过/失败信息 |
| `openspec archive <change>` | `openspec done` | 归档变更 | 合并到 specs/ |
| `openspec proposal <name>` | `openspec new` | 创建新变更 | 生成变更模板 |

#### 3.2.1 初始化项目

```bash
# 在项目根目录执行
openspec init

# 输出示例：
# ✓ Created openspec/specs/
# ✓ Created openspec/changes/
# ✓ Created openspec/config.json
# ✓ Created AGENTS.md
#
# OpenSpec project initialized successfully!
```

#### 3.2.2 列出所有变更

```bash
# 查看所有变更
openspec list

# 输出示例：
# ┌─────────────────────┬──────────┬─────────────────────┐
# │ CHANGE              │ STATUS   │ CREATED             │
# ├─────────────────────┼──────────┼─────────────────────┤
# │ add-user-auth       │ pending  │ 2026-01-18 10:30    │
# │ fix-login-bug       │ approved │ 2026-01-17 15:20    │
# │ refactor-api        │ archived │ 2026-01-15 09:00    │
# └─────────────────────┴──────────┴─────────────────────┘
```

#### 3.2.3 查看变更详情

```bash
# 查看指定变更
openspec show add-user-auth

# 或使用简写
openspec cat add-user-auth

# 输出示例：
# ┌─────────────────────────────────────────────────────────────┐
# │ Change: add-user-auth                                        │
# ├─────────────────────────────────────────────────────────────┤
# │                                                              │
# │ ## Proposal                                                  │
# │                                                             │
# │ ### 标题                                                     │
# │ 添加用户认证功能                                             │
# │                                                             │
# │ ### 目标                                                     │
# │ - 实现用户注册/登录                                          │
# │ - JWT Token 认证                                            │
# │                                                             │
# │ ### Tasks                                                    │
# │ 1. 创建用户模型                                              │
# │ 2. 实现注册接口                                              │
# │ 3. 实现登录接口                                              │
# │ 4. 添加认证中间件                                            │
# │                                                             │
# └─────────────────────────────────────────────────────────────┘
```

#### 3.2.4 校验变更

```bash
# 校验变更结构
openspec validate add-user-auth

# 输出示例（通过）：
# ✓ Change structure is valid
# ✓ Proposal has required sections
# ✓ Tasks are properly defined
#
# Validation passed!

# 输出示例（失败）：
# ✗ Change structure is invalid
# ✗ Missing required field: "背景" in proposal
# ✗ Task 2 has no verification command
#
# Validation failed. Please fix the issues above.
```

#### 3.2.5 归档变更

```bash
# 归档已完成的变更
openspec archive add-user-auth --yes

# 输出示例：
# ✓ Merged tasks to specs/
# ✓ Updated baseline
# ✓ Archived change: add-user-auth
#
# Change archived successfully!
```

### 3.3 使用 Claude Code 原生命令

如果你的 Claude Code 支持 OpenSpec，可以直接使用斜杠命令：

```text
# 在 Claude Code 对话中使用

/openspec-proposal add-user-auth
# → 创建新的变更提案

/openspec-apply add-user-auth
# → 开始实现变更任务

/openspec-archive add-user-auth
# → 归档变更
```

**注意**：如果斜杠命令不可用，使用自然语言描述同样动作即可。

### 3.4 完整使用示例

#### 场景：添加用户登录功能

**Step 1: 创建 Spec 文件**

```bash
# 创建 spec
cat > openspec/specs/user-login.md << 'EOF'
# 用户登录功能

## 背景 & 目标
- 为系统添加用户登录能力
- 支持邮箱/密码登录
- 返回 JWT Token

## 范围
- 登录接口 POST /login
- JWT Token 生成
- 基础错误处理

## 非范围
- 用户注册（另一个 change）
- 第三方登录
- 密码找回

## 验收标准
- 正确邮箱/密码返回 200 + token
- 错误凭证返回 401
- 缺少参数返回 400

## 测试计划
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```
EOF
```

**Step 2: 创建 Change**

```bash
# 使用 CLI 创建
openspec proposal add-user-login

# 或手动创建
mkdir -p openspec/changes/add-user-login
cat > openspec/changes/add-user-login/proposal.md << 'EOF'
# 变更：添加用户登录功能

## 背景
当前系统无用户认证，需要添加登录能力。

## 目标
- [ ] 实现 /login 接口
- [ ] 返回 JWT Token
- [ ] 错误处理完善

## 任务
1. 创建 User 模型
2. 实现 JWT 工具函数
3. 实现登录 handler
4. 添加测试
EOF
```

**Step 3: 与 Claude 协作实现**

在 Claude Code 中输入：

```text
请阅读 openspec/specs/user-login.md 和 openspec/changes/add-user-login/proposal.md
基于此给出 4-6 步实现计划，每步说明：
1. 要改哪些文件
2. 怎么验证
3. 预期结果
```

**Step 4: 验证与归档**

```bash
# 实现完成后，验证变更结构
openspec validate add-user-login

# 验证通过后归档
openspec archive add-user-login --yes
```

---

## 4. Spec 编写规范

### 4.1 Spec 模板

```markdown
# [项目/需求名称]

## 背景 & 目标
- 业务场景：...
- 要解决的问题：...
- 成功标准：...

## 范围 / 非范围
- **范围**：
  - 后端：...
  - 前端：...
- **非范围**：
  - 明确不做：...

## 需求细项
- **后端**：
  - 接口：...
  - 数据模型：...
  - 鉴权：...
- **前端**：
  - 路由：...
  - 页面：...
  - 状态管理：...
  - 适配：...

## 验收标准
- **正常场景**：
  - 接口返回：...
  - 状态码：...
  - UI 交互：...
- **错误/边界场景**：
  - 错误码：...
  - 错误提示：...
  - 降级策略：...

## 测试计划
- **自动化**：
  ```bash
  go test ./...
  npm run test
  ```
- **手动验证**：
  1. ...
  2. ...

## 约束 / 依赖
- 技术栈：...
- 端口：...
- 性能要求：...
- 安全要求：...

## 交付物
- 代码文件：...
- 配置：...
- 文档：...
```

### 4.2 Spec 编写检查清单

在开始实现前，确认 Spec 包含：

- [ ] 明确的业务目标和成功指标
- [ ] 清晰的范围/非范围边界
- [ ] 可执行的验收标准（含错误场景）
- [ ] 可复制的测试命令
- [ ] 技术约束说明
- [ ] 完整的交付物列表

---

## 5. Claude 协作话术

### 5.1 阶段一：阅读与澄清

**让 Claude 阅读上下文**
```text
请先阅读以下文件，总结依赖与入口：
- openspec/specs/xxx.md
- go.mod 或 package.json
- main.go 或 src/main.ts
```

**让 Claude 提出问题**
```text
基于 Spec，列出 1-3 个需要澄清的问题。
如果没有疑问，请明确说明"已理解，可以开始计划"。
```

### 5.2 阶段二：制定计划

**只出计划，不写代码**
```text
先不要写代码。
请基于 Spec 给出 4-6 步实现计划。
每一步必须说明：
1. 要改哪些文件
2. 怎么验证（终端命令）
3. 预期结果
```

### 5.3 阶段三：小步实现

**强约束执行（防止越界）**
```text
开始第 N 步。
约束：
- 只允许改以下文件：<文件列表>
- 不引入新依赖
- 不做无关重构
交付：
- 改动摘要
- Diff 审阅的 3 个关键点
- 终端验证命令（含预期输出）
```

### 5.4 阶段四：验证与反馈

**验证失败时的反馈**
```text
我已运行验证命令，输出如下：

<粘贴终端输出>

请判断是否通过。
如果失败，给出最小修复方案。
```

### 5.5 阶段五：收尾交付

**总结交付**
```text
总结本次开发：
1. 改动的文件列表
2. 测试结果
3. 剩余风险/待办
4. 与 Spec 的一致性说明
```

---

## 6. 代码实现规范

### 6.1 通用规范

| 规范 | 说明 | 示例 |
|------|------|------|
| 文件命名 | kebab-case | `user-service.ts`, `auth_middleware.go` |
| 目录结构 | 按功能分层 | `internal/`, `pkg/`, `src/` |
| 注释 | 解释"为什么"，非"是什么" | `// 使用 bcrypt 防止彩虹表攻击` |
| 错误处理 | 明确错误类型，不吞错误 | `return fmt.Errorf("...: %w", err)` |
| 日志 | 结构化日志，包含上下文 | `log.WithError(err).WithField("user", id)...` |

### 6.2 后端规范 (Go)

```go
// 1. 项目结构
project/
├── cmd/api/          # 主程序入口
├── internal/         # 私有代码
│   ├── auth/         # 认证逻辑
│   ├── handler/      # HTTP 处理器
│   ├── model/        # 数据模型
│   └── repo/         # 数据访问
├── pkg/              # 可导出代码
└── go.mod

// 2. 错误处理示例
func (s *Service) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, fmt.Errorf("user not found: %d", id)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}

// 3. API 响应格式
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 6.3 前端规范 (Vue/TS)

```typescript
// 1. 目录结构
src/
├── api/          # API 调用
├── assets/       # 静态资源
├── components/   # 公共组件
├── views/        # 页面组件
├── router/       # 路由配置
├── stores/       # 状态管理
└── types/        # 类型定义

// 2. API 封装示例
// api/http.ts
import axios from 'axios'

const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default http

// 3. 组件命名
// 使用 PascalCase，与功能相关
<UserProfile />
<LoginForm />
```

### 6.4 安全规范

| 规范 | 说明 |
|------|------|
| 敏感信息 | 不写进代码，使用环境变量 |
| Token | 使用安全存储（localStorage 限开发，生产用 httpOnly cookie） |
| 输入验证 | 前后端都要验证，不信任客户端 |
| SQL 注入 | 使用参数化查询/ORM |
| XSS | 输出时转义，避免 `v-html` 除非必要 |

---

## 7. 验收与验证规范

### 7.1 验收标准模板

```markdown
### 后端验收
- [ ] 健康检查返回 200
- [ ] 接口返回正确的状态码（200/201/400/401/403/404/500）
- [ ] 错误场景有明确提示
- [ ] `go test ./...` 全部通过
- [ ] 无安全漏洞（sql注入、xss等）

### 前端验收
- [ ] 页面正常渲染，无控制台错误
- [ ] 交互符合设计预期
- [ ] 移动端适配正常
- [ ] `npm run build` 成功
- [ ] 关键用户流程可完成
```

### 7.2 验证命令清单

```bash
# 后端
go mod tidy                    # 整理依赖
go test ./... -v              # 运行测试
go run .                      # 启动服务
curl -i http://localhost:8080/health  # 健康检查

# 前端
npm install                   # 安装依赖
npm run dev                   # 开发模式
npm run build                 # 构建
npm run preview               # 预览构建结果
```

---

## 8. 常见反模式

### 8.1 避免这些做法

| 反模式 | 问题 | 正确做法 |
|--------|------|----------|
| "一口气全做完" | 无法审阅，风险高 | 分步执行，每步验证 |
| "你看着改" | 范围失控 | 明确文件列表 |
| 跳过 Spec 直接写代码 | 需求不明确，返工率高 | 先 Spec，再计划，后实现 |
| 只关注正常流程 | 边界场景出问题 | 验收标准包含错误场景 |
| 没有验证命令 | 无法确认是否成功 | 每步提供可复制的验证命令 |

### 8.2 Claude 对话红线

```text
# 不要说
- "你帮我改一下整个项目"
- "把代码重构一下"
- "随便你怎么改"

# 要说
- "改 internal/auth 下文件，实现登录功能"
- "第 1 步只改 main.go，添加健康检查接口"
- "改完后给验证命令和预期输出"
```

---

## 附录：快速参考

### A.1 决策树：我需要安装 OpenSpec CLI 吗？

```
开始
  │
  ├─ 你主要使用 Claude Code 吗？
  │   ├─ 是 → 不需要安装 CLI，使用 /openspec-* 命令
  │   └─ 否 → 需要安装 CLI
  │
  ├─ 你需要独立管理 Spec（不用 Claude）吗？
  │   ├─ 是 → 需要安装 CLI
  │   └─ 否 → 可选安装
  │
  └─ 你需要 CI/CD 集成吗？
      ├─ 是 → 需要安装 CLI
      └─ 否 → 可选安装
```

### A.2 安装速查

```bash
# 1. 安装
npm install -g @fission-ai/openspec@latest

# 2. 验证
openspec --version

# 3. 初始化项目
cd /path/to/project
openspec init
```

### A.3 常用命令速查

```bash
openspec init              # 初始化
openspec list              # 列出变更
openspec show <name>       # 查看变更
openspec validate <name>   # 校验变更
openspec archive <name>    # 归档变更
```

### A.4 Spec 结构速查

```
背景目标 → 范围边界 → 需求细项 → 验收标准 → 测试计划 → 约束依赖 → 交付物
```

### A.5 对话流程速查

```
读 Spec → 提问题 → 定计划 → 确认 → 执行第 N 步 → 验证 → 下一步
```

### A.6 交付检查清单

每次交付必须包含：

- [ ] 改动的文件路径列表
- [ ] 改动摘要（做了什么）
- [ ] Diff 审阅要点（看哪里）
- [ ] 验证命令（怎么验证）
- [ ] 预期输出（应该是什么样）

---

*本规范遵循 OpenSpec 开源项目规范，详见 [Fission-AI/OpenSpec](https://github.com/fission-ai/openspec)*
