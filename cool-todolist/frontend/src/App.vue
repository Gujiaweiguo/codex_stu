<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useTodos } from './composables/useTodos'
import { useTheme } from './composables/useTheme'
import TodoForm from './components/TodoForm.vue'
import TodoItem from './components/TodoItem.vue'
import TodoStats from './components/TodoStats.vue'
import ThemeToggle from './components/ThemeToggle.vue'
import EmptyState from './components/EmptyState.vue'

const { todos, loading, error, stats, loadTodos, addTodo, removeTodo, toggleTodo, updateTodo } = useTodos()
const { theme, isDark, setTheme } = useTheme()

const draggedItem = ref<number | null>(null)

// 初始化加载
onMounted(() => {
  loadTodos()
})

// 添加任务
async function handleAdd(title: string) {
  await addTodo(title)
}

// 删除任务（带确认）
async function handleDelete(id: number) {
  if (confirm('确定要删除这个任务吗？')) {
    await removeTodo(id)
  }
}

// 编辑任务
async function handleEdit(id: number, title: string) {
  await updateTodo(id, { title })
}

// 拖拽开始
function handleDragStart(id: number) {
  draggedItem.value = id
}

// 拖拽结束
function handleDragEnd() {
  draggedItem.value = null
}

// 拖拽悬停
function handleDragOver(id: number) {
  if (draggedItem.value === null || draggedItem.value === id) return

  const draggedIndex = todos.value.findIndex(t => t.id === draggedItem.value)
  const targetIndex = todos.value.findIndex(t => t.id === id)

  if (draggedIndex !== -1 && targetIndex !== -1) {
    const newTodos = [...todos.value]
    const [removed] = newTodos.splice(draggedIndex, 1)
    newTodos.splice(targetIndex, 0, removed)
    todos.value = newTodos
  }
}

// 拖拽释放
async function handleDrop() {
  if (draggedItem.value !== null) {
    const newOrder = todos.value.map(t => t.id)
    await updateTodo(draggedItem.value, {})
    // 这里应该调用 reorderTodos API
    draggedItem.value = null
  }
}
</script>

<template>
  <div class="min-h-screen transition-colors duration-300">
    <!-- 背景装饰 -->
    <div class="fixed inset-0 -z-10 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-purple-300/30 dark:bg-purple-700/20
                    rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-300/30 dark:bg-blue-700/20
                    rounded-full blur-3xl"></div>
    </div>

    <!-- 主容器 -->
    <div class="max-w-2xl mx-auto px-4 py-8">
      <!-- 头部 -->
      <header class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-4xl font-bold bg-gradient-to-r from-primary-500 to-purple-500
                       bg-clip-text text-transparent">
            酷炫 TodoList
          </h1>
          <p class="text-gray-500 dark:text-gray-400 mt-1">
            管理你的日常任务
          </p>
        </div>
        <ThemeToggle :theme="theme" :is-dark="isDark" @set-theme="setTheme" />
      </header>

      <!-- 错误提示 -->
      <div
        v-if="error"
        class="mb-6 p-4 rounded-xl bg-red-50 dark:bg-red-900/30 text-red-600
               dark:text-red-400 flex items-center gap-3"
      >
        <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ error }}
      </div>

      <!-- 统计信息 -->
      <TodoStats :stats="stats" />

      <!-- 添加任务表单 -->
      <TodoForm :loading="loading" @add="handleAdd" />

      <!-- 任务列表 -->
      <TransitionGroup
        tag="div"
        name="todo-item"
        class="mt-6"
      >
        <TodoItem
          v-for="todo in todos"
          :key="todo.id"
          :todo="todo"
          :dragging="draggedItem === todo.id"
          draggable="true"
          @toggle="toggleTodo"
          @delete="handleDelete"
          @edit="handleEdit"
          @dragstart="handleDragStart(todo.id)"
          @dragend="handleDragEnd"
          @dragover.prevent="handleDragOver(todo.id)"
          @drop="handleDrop"
        />
      </TransitionGroup>

      <!-- 空状态 -->
      <EmptyState v-if="!loading && todos.length === 0" />

      <!-- 加载中 -->
      <div v-if="loading" class="flex justify-center py-8">
        <div class="w-8 h-8 border-4 border-primary-500 border-t-transparent
                      rounded-full animate-spin"></div>
      </div>
    </div>
  </div>
</template>

<style>
.todo-item-enter-active,
.todo-item-leave-active {
  transition: all 0.3s ease;
}

.todo-item-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.todo-item-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.todo-item-move {
  transition: transform 0.3s ease;
}
</style>
