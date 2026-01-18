<script setup lang="ts">
import { useRouter } from 'vue-router'
</script>

<template>
  <div class="flex flex-col items-center justify-center py-16 px-4">
    <!-- 动画插图 -->
    <div class="relative mb-8">
      <div class="w-32 h-32 rounded-full bg-gradient-to-br from-primary-400 to-purple-500
                    animate-pulse opacity-20"></div>
      <div class="absolute inset-0 flex items-center justify-center">
        <svg class="w-16 h-16 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
        </svg>
      </div>
    </div>

    <!-- 提示文字 -->
    <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-300 mb-2">
      还没有任务
    </h3>
    <p class="text-gray-500 dark:text-gray-400 text-center mb-6">
      添加一个新任务开始吧！
    </p>

    <!-- 示例任务建议 -->
    <div class="flex flex-wrap gap-2 justify-center">
      <span class="px-3 py-1 text-sm bg-white/50 dark:bg-gray-700/50
                   rounded-full text-gray-600 dark:text-gray-300">
        📝 阅读一本书
      </span>
      <span class="px-3 py-1 text-sm bg-white/50 dark:bg-gray-700/50
                   rounded-full text-gray-600 dark:text-gray-300">
        🏃 跑步 30 分钟
      </span>
      <span class="px-3 py-1 text-sm bg-white/50 dark:bg-gray-700/50
                   rounded-full text-gray-600 dark:text-gray-300">
        💻 学习 Vue 3
      </span>
    </div>
  </div>
</template>
