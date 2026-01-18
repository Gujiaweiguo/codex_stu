<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  loading?: boolean
}>()

const emit = defineEmits<{
  add: [title: string]
}>()

const input = ref<HTMLInputElement>()
const title = ref('')

async function handleSubmit() {
  if (!title.value.trim() || props.loading) return
  emit('add', title.value)
  title.value = ''
  input.value?.focus()
}

defineExpose({
  focus: () => input.value?.focus(),
})
</script>

<template>
  <form @submit.prevent="handleSubmit" class="relative">
    <input
      ref="input"
      v-model="title"
      type="text"
      placeholder="添加新任务..."
      :disabled="loading"
      class="w-full px-6 py-4 pr-14 text-lg rounded-2xl glass-card
             focus:outline-none focus:ring-2 focus:ring-primary-500
             placeholder-gray-400 dark:placeholder-gray-500
             disabled:opacity-50 disabled:cursor-not-allowed
             transition-all duration-200"
    />
    <button
      type="submit"
      :disabled="loading || !title.trim()"
      class="absolute right-2 top-1/2 -translate-y-1/2
             w-10 h-10 rounded-xl bg-primary-500 hover:bg-primary-600
             text-white font-bold text-lg
             disabled:opacity-50 disabled:cursor-not-allowed
             transition-all duration-200
             hover:scale-105 active:scale-95"
    >
      +
    </button>
  </form>
</template>
