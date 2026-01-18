# 酷炫 TodoList

一个功能完善、界面酷炫的待办事项应用，使用 Vue 3 + Go 构建。

## 特性

- ✅ **任务管理**：添加、编辑、删除、完成任务
- ✨ **拖拽排序**：通过拖拽调整任务顺序
- 🎨 **主题切换**：支持浅色/深色/跟随系统主题
- 📊 **数据统计**：实时显示任务完成度和统计信息
- 🎭 **动画效果**：流畅的过渡动画和交互反馈
- 📱 **响应式设计**：适配桌面和移动端
- 💾 **数据持久化**：SQLite 本地存储

## 技术栈

### 后端
- Go 1.21+
- Gin (Web 框架)
- GORM (ORM)
- SQLite (数据库)

### 前端
- Vue 3 (Composition API)
- TypeScript
- Vite (构建工具)
- TailwindCSS (样式)
- Pinia (状态管理)
- VueUse (组合式工具库)

## 项目结构

```
cool-todolist/
├── backend/                # Go 后端
│   ├── main.go            # 入口文件
│   ├── models/            # 数据模型
│   ├── handlers/          # API 处理器
│   ├── database/          # 数据库配置
│   └── data/              # SQLite 数据库文件
└── frontend/              # Vue 前端
    ├── src/
    │   ├── components/    # Vue 组件
    │   ├── composables/   # 组合式函数
    │   ├── api/          # API 调用
    │   └── types/        # TypeScript 类型
    └── package.json
```

## 快速开始

### 后端

```bash
cd backend

# 安装依赖
go mod download

# 启动服务
go run .
```

后端服务将在 `http://localhost:8080` 启动。

### 前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

## API 接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/todos` | 获取所有任务 |
| POST | `/api/todos` | 创建任务 |
| PUT | `/api/todos/:id` | 更新任务 |
| PATCH | `/api/todos/:id/toggle` | 切换完成状态 |
| DELETE | `/api/todos/:id` | 删除任务 |
| POST | `/api/todos/reorder` | 重新排序 |

## 使用说明

1. **添加任务**：在输入框中输入任务名称，按回车或点击 + 按钮
2. **完成任务**：点击任务左侧的复选框
3. **编辑任务**：双击任务标题或点击编辑按钮
4. **删除任务**：点击删除按钮（需要确认）
5. **排序任务**：拖拽任务右侧的手柄进行调整
6. **切换主题**：点击右上角的主题切换按钮

## 构建部署

### 后端构建

```bash
cd backend
go build -o todolist-server .
./todolist-server
```

### 前端构建

```bash
cd frontend
npm run build
```

构建产物在 `frontend/dist` 目录。

## 开发规范

本项目遵循 [OpenSpec 编程规范](../docs/claude_openspec_coding_standards.md)：

1. 先写 Spec，再实现
2. 小步交付，每步验证
3. 对话驱动开发

## 许可证

MIT
