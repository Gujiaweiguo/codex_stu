# Codex Open Spec 开发指南：Go RBAC + Vue Web/H5

## 读完可以
- 用 Open Spec 模式驱动 Codex：写 Spec、澄清、计划、小步实现、验证、交付。
- 完成一套 Go RBAC 后端（Gin/Gorm/JWT/SQLite）+ Vue3 Web/H5 前端（Vite + Vue Router + Pinia + Axios），含验证步骤。
- 知道对 Codex 下指令的脚本式话术，避免“一口气全做完”。

## 1. Open Spec 工作法（6 步）
1) 写 Spec：目标、范围/非范围、需求细项、验收标准、测试计划、约束、交付物。
2) 让 Codex 阅读并提问：先读 Spec/代码，明确疑问再动手。
3) 要计划：3~6 步，覆盖实现与测试，确认后执行。
4) 小步实现：一次只做一小块；限制修改范围；改完解释文件与理由。
5) 验证：跑测试/构建/手动步骤，对照验收标准；失败先给修复方案。
6) 收尾：总结改动、影响面、测试状态、剩余风险/待办。

### 推荐 Spec 模版
```markdown
# 项目/需求名称
## 背景 & 目标
## 范围 / 非范围
- 范围：
- 非范围：
## 需求细项
- 后端：…
- 前端：…
## 验收标准
- 接口返回/状态码/权限；前端页面/交互/适配；错误场景。
## 测试计划
- 自动化命令、需要补的用例；手动验证步骤。
## 约束 / 依赖
- 技术栈、性能、兼容性、API 契约、端口等。
## 交付物
- 代码文件、配置、文档、截图/录屏（可选）。
```

## 2. 环境准备
- 后端：Go >= 1.20，SQLite（默认文件型），`rg`/`curl`，可运行 `go test ./...`、`go run .`。
- 前端：Node >= 18，包管理器（pnpm/npm/yarn），浏览器或 `vite preview`；移动端用浏览器设备模式或真机。
- Codex CLI：`codex --version` 确认可用；对话尽量在仓库根目录。

## 3. 对话套路（可直接复制改词）
- 读上下文：`请先看 go.mod、main.go（或 cmd/api/main.go，以你的项目结构为准）、frontend/package.json，总结依赖与入口。`
- 提计划：`基于 Spec 给 4~6 步计划，含实现与测试；再问我不清楚的。`
- 控范围：`现在执行计划第 1 步，只改 internal/auth 下文件，改完解释改动与风险。`
- 前端约束：`只改 frontend/src 下，保持 Vite 配置不变；给出修改文件列表。`
- 测试/构建：`运行 go test ./... 和 pnpm test（或 pnpm build），贴失败摘要。`
- 收尾：`总结改动、测试结果、剩余风险/待办，列出改动的文件路径。`

## 4. Open Spec 示例（Go RBAC + Vue Web/H5）
### 4.1 目标与栈
- 后端：用户注册/登录（JWT），角色/权限管理，受保护接口示例，SQLite 持久化。
- 前端：Vue3 + Vite + TypeScript + Pinia + Vue Router + Axios，Web/H5 响应式布局；提供登录、权限受控的 Demo 页面与角色/权限管理界面。

### 4.2 初始化命令（你在终端执行；我在对话框里给命令/解释）
建议流程：
1) 你在对话框里让 Codex 先“只出计划不写代码”（见 6.0 的模板 1）。
2) 你在终端执行初始化命令；把关键输出（成功/报错）粘贴回对话框。

下面是一个可复制的初始化脚本（示例目录叫 `rbac-demo`，你可自行改名）：
```bash
# backend
mkdir rbac-demo && cd rbac-demo
go mod init rbac-demo
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/sqlite github.com/golang-jwt/jwt/v4 golang.org/x/crypto/bcrypt

# frontend
mkdir frontend && cd frontend
npm create vite@latest . -- --template vue-ts
npm install
npm install axios pinia vue-router
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

### 4.3 后端开发（按“计划 → 确认 → 分步执行”）
你的动作始终是三件事：在对话框下指令 → 在终端跑验证 → 在 Diff 审阅是否越界。

**A) 先要后端计划（对话框）**
```text
先不要写代码。请阅读项目根目录结构和 go.mod（我会打开），然后给出 4–6 步后端计划。
每一步写清：要改哪些文件、要提供哪些接口、怎么验证（终端命令）、预期结果。
```

**B) 执行计划第 N 步（对话框）**
每一步都用“强约束”避免跑偏（用 6.0 模板 2 也行）。例如：
```text
开始后端计划第 1 步。这一步只允许改 main.go。
目标：先跑起来一个最小服务：SQLite 初始化 + 自动迁移 + /health 200 + 基础路由结构。
改完给：改动摘要 + Diff 审阅点 + 我在终端要跑的验证命令（含预期）。
```

**C) 你在终端验证，把输出贴回对话框**
示例命令（按 Codex 给你的为准）：
```bash
go test ./...
go run .
curl -i http://127.0.0.1:8080/health
```

后端内容通常按下面顺序小步完成（用于你检查计划是否覆盖到位）：
1) 骨架：`main.go`（或你项目入口）包含 SQLite 初始化、模型 User/Role/Permission、路由注册、`/health`。
2) 认证：`/signup` + bcrypt，`/login` + JWT；中间件解析 Bearer token。
3) 权限校验：`requirePerm(perms...)`（或等价机制），权限命名 `perm:*`。
4) RBAC 管理：`/roles`、`/permissions` 创建；分配关系（角色-权限、用户-角色）。
5) 受保护接口：`/demo` 需 `perm:demo:read`；`/me` 返回用户/角色/权限。
6) 种子数据与测试：初始化 admin 用户/默认权限；`go test ./...`（必要时补 1–2 个 httptest 覆盖 `/login`、`/demo`）。

### 4.4 前端开发（Web/H5，同一代码基；同样按分步逻辑）
**A) 先要前端计划（对话框）**
```text
先不要写代码。请阅读 frontend/package.json 和 src 目录结构（我会打开），给出 4–6 步前端计划：
路由/状态/API/页面/适配/构建（可选测试）。
每一步写清：要改哪些文件、怎么验证（终端命令）、预期结果。
```

**B) 执行计划第 M 步（对话框）**
例如（限制只改 `frontend/src`）：
```text
开始前端计划第 1 步。这一步只允许改 frontend/src 下文件（列出你要改的路径）。
目标：接入路由 + 登录页骨架 + token 存储（先不追求样式）。
改完给：改动摘要 + Diff 审阅点 + 验证命令（含预期）。
```

**C) 你在终端验证**
```bash
npm run dev
npm run build
npm run preview
```

前端内容通常按下面顺序小步完成（用于你检查计划是否覆盖到位）：
1) 路由：`/login`、`/demo`、`/roles`、`/permissions`（`/users` 可选）。
2) 状态与 API：Axios 自动带 Bearer；Pinia 存 token/user；必要时用 localStorage 持久化。
3) 路由守卫：无 token → 跳 `/login`；可用 meta 声明 `permission` 并在守卫中校验（调用 `/me` 或缓存）。
4) 页面：登录页、Demo 页、角色/权限管理页（先能用，再做美化）。
5) H5 适配：窄屏下表单/按钮堆叠不溢出；检查 viewport 与点击区域。
6) 构建：`npm run build && npm run preview` 成功。

#### 前端关键片段示例
`src/api/http.ts`
```ts
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const http = axios.create({ baseURL: 'http://localhost:8080' })

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) config.headers.Authorization = `Bearer ${auth.token}`
  return config
})

export default http
```

`src/stores/auth.ts`
```ts
import { defineStore } from 'pinia'
import http from '../api/http'

type User = { id: number; username: string; roles: string[]; permissions: string[] }

export const useAuthStore = defineStore('auth', {
  state: () => ({ token: '', user: null as User | null }),
  actions: {
    async login(username: string, password: string) {
      const { data } = await http.post('/login', { username, password })
      this.token = data.token
      await this.fetchMe()
    },
    async fetchMe() {
      const { data } = await http.get('/me')
      this.user = data
    },
    logout() {
      this.token = ''
      this.user = null
    },
  },
})
```

`src/router/index.ts`
```ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Login from '../views/Login.vue'
import Demo from '../views/Demo.vue'

const routes = [
  { path: '/login', component: Login },
  { path: '/demo', component: Demo, meta: { requiresAuth: true } },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.token) return '/login'
  return true
})

export default router
```

> 让 Codex 生成完整页面和样式时，要求：文件路径、依赖、移动端适配说明、如何运行与预览。

### 4.5 验收与验证
- 后端：
  - `go run .`；`/health` 200。
  - `/login` admin -> 得 token；受保护接口 `/demo` 返回 200；无权限用户返回 403。
  - `/roles`/`/permissions`/`/users/:id/roles` 能创建/分配成功；`/me` 返回角色/权限。
- 前端：
  - `npm run dev` 可访问登录页；窄屏下表单/按钮正常排列。
  - 登录 admin 后能访问 Demo 页面并看到受保护数据。
  - 角色/权限页面能创建并分配，操作失败有提示。
  - 构建：`npm run build && npm run preview` 成功；H5 视口正常。
- 测试：可要求 Codex 补充 1-2 个后端 httptest、前端组件/路由守卫测试（如 Vitest）。

## 5. OpenSpec CLI（可选，Fission-AI/OpenSpec）
- 安装与初始化（需要联网）：`npm install -g @fission-ai/openspec@latest` → `openspec --version` → 在项目根执行 `openspec init`。会生成 `openspec/` 目录与项目根的 `AGENTS.md`，并为支持的助手写入 `/openspec-*` 命令（Codex 原生支持）。
- 核心命令（精简版）：`openspec list` 查看变更；`openspec show <change>` 查看 proposal/tasks/spec delta；`openspec validate <change>` 校验结构；`openspec archive <change> --yes` 归档并把 delta 合并到 specs。
- 典型流程：`/openspec-proposal <name>`（或让 Codex/你执行同名命令）生成变更；确认后 `/openspec-apply <name>` 逐项落实 tasks；完成后 `/openspec-archive <name>` 归档。非原生助手用自然语言描述同样动作即可。
- 用法提示：保持 Spec/changes 都在 `openspec/` 下，Codex 对话时让它读取对应文件；CLI 报告的结果（list/validate/show）贴给 Codex，对齐后再动手。
- 原则：Spec 是单一事实来源，CLI 输出只是辅助；以 Spec + 测试结果为最终准绳。

## 6. 让 Codex 带你一步步产出代码（分步提示词）
### 6.0 VS Code 插件里：你在哪做什么
把它记成三件事：**对话框指挥**、**终端验证**、**Diff 审阅**。

- 你在 **Codex 对话框** 里做：说清目标/约束、指定“这一步只改哪些文件”、贴终端输出/报错、确认继续下一步。
- 你在 **VS Code 集成终端** 里做：实际运行/测试/构建/安装依赖，把输出原样粘贴回对话框（成功也贴一行关键输出）。
- 你在 **Source Control / Diff** 里做：审阅改动是否越界（是否改了不该改的文件/引入了不需要的依赖/动了目录结构），不满意就要求回滚或按规则重改。

下面给一组“可直接复制到对话框”的模板：

1) 先只要计划（不写代码）：
```text
先不要写代码。请阅读我工作区里与本任务相关的文件（Spec/README/go.mod/目录结构），然后给出 4–6 步计划。
要求：每一步写清（要改哪些文件、怎么验证=终端命令、预期结果）。
我回复“开始第 1 步”后你再实现。
```

2) 开始某一步（强约束，避免跑偏）：
```text
开始第 N 步。这一步只允许改这些文件：<填文件列表>。
不要引入新依赖；不要移动/重命名文件；不要做无关重构。
改完只交付：改动摘要 + 我看 diff 的 3 个点 + 终端验证命令（含预期成功标准）。
```

3) 我跑完验证，把结果贴回去：
```text
我已运行你给的验证命令，输出如下（原样粘贴）：
<这里粘贴终端输出/报错>
请判断是否通过；如果没通过，给最小修复方案，并继续按“只改指定文件”的约束修复。
```

0) 初始化（可选）  
`已安装 openspec CLI，执行 openspec init 生成 openspec/ 与 AGENTS.md。`

1) 同步目标与上下文  
`请阅读项目内的 Spec（如 openspec/project.md、openspec/specs/... 或 docs/spec.md），确认目标（本例是 Go RBAC + Vue Web/H5）和验收标准；先提需要澄清的问题。若没有现成 Spec，请先和我一起写。`
写 Spec 的详细步骤（可让 Codex 执行并输出文件草稿）：
- 给背景与目标：业务场景、要解决的问题、成功标准。
- 写业务目标/成功指标：明确业务侧想要的可衡量结果（如活跃、转化、工单下降等），并与后续验收标准对应。
- 定范围/非范围：明确不做什么，避免膨胀。
- 列需求细项：后端接口/鉴权/数据模型，前端路由/状态/页面/适配；必要时给输入输出示例。
- 写验收标准：正常/错误/边界场景，对应的返回码、文案、UI 反馈；前端交互和响应式要求。
- 写测试计划：自动化命令（如 `go test ./...`、`npm run test`）、手动步骤（curl 或页面操作+预期），特殊环境说明。
- 写约束/依赖：技术栈、端口、数据源、鉴权方式、性能/安全要求。
- 写交付物：代码文件、配置、文档/截图等。
- 让 Codex 输出疑问清单，补充后再让它精炼成最终 Spec 文件。
示例提示词：`请按上面的步骤，基于“Go RBAC + Vue Web/H5”起草 Spec，包含范围/非范围、需求细项、验收标准（含错误/边界）、测试计划、约束/交付物，并列出不确定点。`

从 0 起草的 Spec 示例（压缩版，可让 Codex 生成更长正式稿）：
```markdown
# RBAC 后台 + Vue Web/H5
## 背景 & 目标
- 为内部后台提供账号登录与权限控制，支持 Web/H5 自适应。
- 业务目标：上线后 2 周内替代旧后台，减少权限相关工单 50%。

## 范围 / 非范围
- 范围：用户注册/登录、JWT 鉴权、角色/权限管理、受保护接口示例、Vue 前端登录/演示页/角色权限管理。
- 非范围：第三方登录、支付、消息推送。

## 需求细项
- 后端：Gin + Gorm + SQLite；接口 /signup /login /me /demo；角色/权限 CRUD 与分配；权限命名 perm:*。
- 前端：Vue3+Vite+TS，路由 /login /demo /roles /permissions；Pinia 管理 token/user；Axios 带 Bearer；窄屏适配。

## 验收标准
- 正常：/login 返回 200+token；携 token 访问 /demo 得 200，未授权 403；/roles /permissions 创建成功 201。
- 错误/边界：重复用户名 409，密码<6 报 400；无 token 401；token 过期/伪造 401。
- 前端：无 token 访问受保护页自动跳转 /login；登录后能看到 Demo 文案；窄屏表单不溢出。

## 测试计划
- 自动化：`go test ./...`；必要时补 httptest 覆盖 /login、/demo；前端可补 Vitest 路由守卫测试。
- 手动：`go run .`；curl 登录、访问 /demo、创建角色/权限、分配角色，再用新用户访问；`npm run dev` 验证登录跳转与 Demo/管理页；`npm run build && npm run preview`。

## 约束 / 依赖
- 技术：Go >=1.20，Node >=18；默认端口 8080/5173；JWT HS256，secret 需可配置；SQLite 文件需可写。

## 交付物
- main.go（或分模块）后端代码；frontend/src 下页面/路由/状态；README 或 docs 更新；若有截图可附。
```

2) 后端计划  
`基于文档和 go.mod，给出 4~6 步后端计划（含实现与测试），只列步骤，不执行。`

3) 后端实现（逐步执行，N 从 1 开始）  
`现在执行后端计划第 N 步，只改 main.go（或你要指定的路径），改完说明改动与风险。`  
完成全部后端步骤后：`运行 go test ./...（如受限，请给我要运行的命令和预期输出），总结后端状态。`

4) 前端计划  
`阅读 frontend/package.json 和 src 结构，基于文档给出 4~6 步前端计划（路由/状态/API/页面/适配/测试），只列步骤，不执行。`

5) 前端实现（逐步执行，M 从 1 开始）  
`执行前端计划第 M 步，只改 frontend/src 下相关文件（列出路径），改完说明改动与风险。`  
完成全部前端步骤后：`运行 npm run build（或 pnpm/yarn 同步），如受限请给命令和预期输出；必要时运行单测。`

6) 联调与验收  
`给出后端 curl 验证脚本与前端 dev/preview 运行方式，检查验收标准是否全部满足；列出任何失败项及修复建议。`

7) 收尾交付  
`总结改动、测试结果、剩余风险/待办，列出修改的文件路径。`

## 7. 常见问题
- JWT/密钥：生产请更换 `jwtSecret`、admin 初始密码；敏感配置用 env。
- 跨域：前端若跨域，后端加 `gin` CORS 中间件或本地代理（Vite `server.proxy`）。
- SQLite 权限：确保运行目录可写；或改内存 `:memory:`；生产改用 MySQL/Postgres。
- H5 适配：检查 viewport、交互区域尺寸（按钮高度 >= 40px）、触摸反馈；必要时用 CSS 媒体查询提升点击区域。

以上即 Open Spec 版本的 Codex 指南：先写 Spec，再让 Codex 按计划小步交付，覆盖 Go 后端和 Vue Web/H5 前端的完整闭环。
