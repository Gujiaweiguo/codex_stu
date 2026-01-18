# 酷炫 TodoList 应用

## 背景 & 目标

### 业务场景
创建一个功能完善、界面酷炫的 TodoList 待办事项应用，帮助用户管理日常任务。

### 要解决的问题
- 现有 TodoList 应用界面简陋，缺乏吸引力
- 缺少数据统计功能，无法了解任务完成情况
- 没有主题切换，不适应不同使用场景

### 成功标准
- 用户可以添加、编辑、删除、完成任务
- 支持拖拽排序任务
- 支持深色/浅色主题切换
- 显示任务完成度统计
- 界面流畅，动画自然

## 范围 / 非范围

### 范围
**后端 (Go)**：
- RESTful API (`/api/todos`)
- SQLite 数据持久化
- CRUD 操作：创建、读取、更新、删除
- 任务状态切换（完成/未完成）
- 任务排序

**前端 (Vue 3)**：
- Vue 3 + Vite + TypeScript
- TailwindCSS 样式
- VueUse 组合式工具
- 任务列表页面
- 任务添加/编辑表单
- 主题切换组件
- 数据统计组件

### 非范围
- 用户注册/登录（单机版本）
- 多用户协作
- 任务分类/标签
- 任务提醒/通知
- 云同步

## 需求细项

### 后端 API

| 方法 | 路径 | 描述 | 请求体 | 响应 |
|------|------|------|--------|------|
| GET | `/api/todos` | 获取所有任务 | - | `[{id, title, completed, order, created_at}]` |
| POST | `/api/todos` | 创建任务 | `{title, order?}` | `{id, title, completed, order, created_at}` |
| PUT | `/api/todos/:id` | 更新任务 | `{title?, completed?, order?}` | 更新后的任务 |
| DELETE | `/api/todos/:id` | 删除任务 | - | `success: true` |
| PATCH | `/api/todos/:id/toggle` | 切换完成状态 | - | 更新后的任务 |

**数据模型**：
```go
type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    Order     int       `json:"order"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 前端功能

**核心功能**：
1. 添加任务（回车或点击按钮）
2. 编辑任务（双击或编辑按钮）
3. 删除任务（删除按钮）
4. 完成任务（点击复选框）
5. 拖拽排序（HTML5 Drag & Drop）

**酷炫特性**：
1. **动画效果**
   - 添加/删除任务时的过渡动画
   - 完成任务时的划线动画
   - 拖拽时的视觉反馈
   - 主题切换的平滑过渡

2. **主题切换**
   - 深色模式（Dark Mode）
   - 浅色模式（Light Mode）
   - 跟随系统（Auto）
   - LocalStorage 持久化

3. **数据统计**
   - 总任务数
   - 完成任务数
   - 完成率（进度条）
   - 今日完成任务数

**UI/UX**：
- 响应式设计（移动端适配）
- 空状态提示
- 操作成功/失败提示
- 确认删除对话框

## 验收标准

### 正常场景
- 启动后端服务，`GET /api/todos` 返回 200 和空数组 `[]`
- 添加任务后，列表中显示新任务
- 完成任务后，任务显示划线样式
- 删除任务后，列表中移除该任务
- 拖拽任务后，刷新页面顺序保持
- 切换主题后，界面立即变化
- 完成率统计实时更新

### 错误/边界场景
- 空标题提交时，显示错误提示
- 网络错误时，显示重试选项
- 删除任务时，弹出确认对话框
- 任务列表为空时，显示空状态插图

### 性能要求
- 页面首次加载 < 2 秒
- API 响应 < 100ms（本地）
- 动画帧率 >= 60fps

## 测试计划

### 后端测试
```bash
# 启动后端
cd backend && go run .

# 测试 API
curl http://localhost:8080/api/todos
curl -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{"title":"测试任务"}'
```

### 前端测试
```bash
# 启动前端
cd frontend && npm run dev

# 访问 http://localhost:5173
# 手动测试所有功能
```

### 构建测试
```bash
# 后端构建
cd backend && go build .

# 前端构建
cd frontend && npm run build
```

## 约束 / 依赖

### 技术栈
- **后端**：Go 1.20+, Gin, GORM, SQLite
- **前端**：Vue 3.3+, Vite 5+, TypeScript, TailwindCSS
- **动画**：VueUse Motion, CSS Transitions
- **拖拽**：@vueuse/core useDraggable

### 端口
- 后端：`8080`
- 前端开发服务器：`5173`

### 数据存储
- SQLite 文件：`backend/data/todos.db`

## 交付物

### 代码结构
```
cool-todolist/
├── backend/
│   ├── main.go
│   ├── models/todo.go
│   ├── handlers/todo.go
│   ├── database/db.go
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── TodoItem.vue
│   │   │   ├── TodoForm.vue
│   │   │   ├── ThemeToggle.vue
│   │   │   └── TodoStats.vue
│   │   ├── composables/
│   │   │   ├── useTodos.ts
│   │   │   └── useTheme.ts
│   │   ├── App.vue
│   │   └── main.ts
│   ├── index.html
│   └── package.json
└── README.md
```

### 文档
- README.md（项目说明、启动步骤）
- API.md（接口文档）
