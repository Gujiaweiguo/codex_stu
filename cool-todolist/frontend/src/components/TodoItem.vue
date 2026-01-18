<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Todo } from '@/types'
import { useMotionValue, useMotionTemplate } from '@vueuse/motion'

const props = defineProps<{
  todo: Todo
  dragging?: boolean
}>()

const emit = defineEmits<{
  toggle: [id: number]
  delete: [id: number]
  edit: [id: number, title: string]
}>()

const isEditing = ref(false)
const editInput = ref<HTMLInputElement>()
const editTitle = ref(props.todo.title)

const titleClass = computed(() =>
  props.todo.completed
    ? 'line-through text-gray-400 dark:text-gray-500'
    : 'text-gray-800 dark:text-gray-100'
)

function startEdit() {
  isEditing.value = true
  editTitle.value = props.todo.title
  setTimeout(() => editInput.value?.focus(), 0)
}

function cancelEdit() {
  isEditing.value = false
  editTitle.value = props.todo.title
}

function saveEdit() {
  if (editTitle.value.trim()) {
    emit('edit', props.todo.id, editTitle.value.trim())
  }
  isEditing.value = false
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') saveEdit()
  if (e.key === 'Escape') cancelEdit()
}
</script>

<template>
  <div
    class="todo-item group flex items-center gap-4 p-4 mb-3 rounded-xl glass-card
           hover:shadow-lg hover:scale-[1.02] active:scale-[0.98]"
    :class="{ 'opacity-50': dragging }"
  >
    <!-- 复选框 -->
    <input
      type="checkbox"
      :checked="todo.completed"
      @change="emit('toggle', todo.id)"
      class="checkbox-custom flex-shrink-0"
    />

    <!-- 任务标题 -->
    <div class="flex-1 min-w-0">
      <input
        v-if="isEditing"
        ref="editInput"
        v-model="editTitle"
        @blur="saveEdit"
        @keydown="handleKeydown"
        class="w-full px-2 py-1 rounded bg-white/50 dark:bg-gray-700/50
               focus:outline-none focus:ring-2 focus:ring-primary-500"
      />
      <span
        v-else
        @dblclick="startEdit"
        :class="[
          'block truncate text-lg cursor-pointer transition-all duration-200',
          titleClass
        ]"
      >
        {{ todo.title }}
      </span>
    </div>

    <!-- 操作按钮 -->
    <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
      <button
        @click="startEdit"
        class="p-2 rounded-lg text-gray-400 hover:text-primary-500 hover:bg-primary-50
               dark:hover:bg-primary-900/30 transition-all"
        title="编辑"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
        </svg>
      </button>
      <button
        @click="emit('delete', todo.id)"
        class="p-2 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50
               dark:hover:bg-red-900/30 transition-all"
        title="删除"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>

    <!-- 拖拽手柄 -->
    <div class="cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-500">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 8h16M4 16h16" />
      </svg>
    </div>
  </div>
</template>
