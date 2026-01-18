import { ref, computed } from 'vue'
import type { Todo, TodoStats } from '@/types'
import { todoApi } from '@/api/todos'

export function useTodos() {
  const todos = ref<Todo[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 统计数据
  const stats = computed<TodoStats>(() => {
    const total = todos.value.length
    const completed = todos.value.filter(t => t.completed).length
    const pending = total - completed
    const completionRate = total > 0 ? Math.round((completed / total) * 100) : 0
    return { total, completed, pending, completionRate }
  })

  // 加载所有任务
  async function loadTodos() {
    loading.value = true
    error.value = null
    try {
      todos.value = await todoApi.getAll()
    } catch (e) {
      error.value = '加载任务失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  // 添加任务
  async function addTodo(title: string) {
    if (!title.trim()) return
    const maxOrder = Math.max(0, ...todos.value.map(t => t.order))
    const newTodo = await todoApi.create(title, maxOrder + 1)
    todos.value.push(newTodo)
  }

  // 删除任务
  async function removeTodo(id: number) {
    await todoApi.delete(id)
    todos.value = todos.value.filter(t => t.id !== id)
  }

  // 切换完成状态
  async function toggleTodo(id: number) {
    const updated = await todoApi.toggle(id)
    const index = todos.value.findIndex(t => t.id === id)
    if (index !== -1) {
      todos.value[index] = updated
    }
  }

  // 更新任务
  async function updateTodo(id: number, updates: Partial<Todo>) {
    const updated = await todoApi.update(id, updates)
    const index = todos.value.findIndex(t => t.id === id)
    if (index !== -1) {
      todos.value[index] = updated
    }
  }

  // 重新排序
  async function reorderTodos(newOrder: number[]) {
    await todoApi.reorder(newOrder)
    // 更新本地顺序
    const newTodos: Todo[] = []
    newOrder.forEach((id, index) => {
      const todo = todos.value.find(t => t.id === id)
      if (todo) {
        newTodos.push({ ...todo, order: index })
      }
    })
    todos.value = newTodos
  }

  return {
    todos,
    loading,
    error,
    stats,
    loadTodos,
    addTodo,
    removeTodo,
    toggleTodo,
    updateTodo,
    reorderTodos,
  }
}
