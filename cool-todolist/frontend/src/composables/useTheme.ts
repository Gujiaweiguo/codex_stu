import { ref, watch, onMounted } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import type { Theme } from '@/types'

const THEME_KEY = 'todolist-theme'

export function useTheme() {
  const isDark = useDark()
  const toggleDark = useToggle(isDark)

  const theme = ref<Theme>('auto')

  // 从 localStorage 加载主题设置
  onMounted(() => {
    const saved = localStorage.getItem(THEME_KEY) as Theme
    if (saved) {
      theme.value = saved
    }
  })

  // 监听主题变化
  watch(theme, (newTheme) => {
    localStorage.setItem(THEME_KEY, newTheme)

    // 根据主题设置 dark mode
    if (newTheme === 'auto') {
      // 跟随系统
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      if (prefersDark !== isDark.value) {
        toggleDark()
      }
    } else if (newTheme === 'dark' && !isDark.value) {
      isDark.value = true
    } else if (newTheme === 'light' && isDark.value) {
      isDark.value = false
    }
  })

  function setTheme(newTheme: Theme) {
    theme.value = newTheme
  }

  return {
    isDark,
    theme,
    setTheme,
    toggleDark,
  }
}
