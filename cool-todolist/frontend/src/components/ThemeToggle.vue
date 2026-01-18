<script setup lang="ts">
import type { Theme } from '@/types'

const props = defineProps<{
  theme: Theme
  isDark: boolean
}>()

const emit = defineEmits<{
  setTheme: [theme: Theme]
}>()

const themes: { value: Theme; label: string; icon: string }[] = [
  {
    value: 'light',
    label: '浅色',
    icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>'
  },
  {
    value: 'dark',
    label: '深色',
    icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path></svg>'
  },
  {
    value: 'auto',
    label: '跟随系统',
    icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>'
  },
]

function setTheme(theme: Theme) {
  emit('setTheme', theme)
}
</script>

<template>
  <div class="glass-card rounded-xl p-2 inline-flex gap-1">
    <button
      v-for="t in themes"
      :key="t.value"
      @click="setTheme(t.value)"
      :class="[
        'p-2 rounded-lg transition-all duration-200',
        theme === t.value
          ? 'bg-primary-500 text-white shadow-md'
          : 'text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700'
      ]"
      :title="t.label"
      v-html="t.icon"
    />
  </div>
</template>
