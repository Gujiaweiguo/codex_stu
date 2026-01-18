<script setup lang="ts">
import type { TodoStats } from '@/types'

defineProps<{
  stats: TodoStats
}>()
</script>

<template>
  <div class="glass-card rounded-2xl p-6 mb-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-700 dark:text-gray-200">
        任务统计
      </h3>
      <span class="text-3xl font-bold text-primary-500">
        {{ stats.completionRate }}%
      </span>
    </div>

    <!-- 进度条 -->
    <div class="relative h-4 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden mb-4">
      <div
        class="absolute inset-y-0 left-0 bg-gradient-to-r from-primary-400 to-primary-600
               rounded-full transition-all duration-500 ease-out"
        :style="{ width: `${stats.completionRate}%` }"
      >
        <div class="absolute inset-0 bg-white/30 animate-pulse"></div>
      </div>
    </div>

    <!-- 详细数据 -->
    <div class="grid grid-cols-3 gap-4 text-center">
      <div class="p-3 rounded-xl bg-blue-50 dark:bg-blue-900/30">
        <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">
          {{ stats.total }}
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400">总计</div>
      </div>
      <div class="p-3 rounded-xl bg-green-50 dark:bg-green-900/30">
        <div class="text-2xl font-bold text-green-600 dark:text-green-400">
          {{ stats.completed }}
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400">已完成</div>
      </div>
      <div class="p-3 rounded-xl bg-orange-50 dark:bg-orange-900/30">
        <div class="text-2xl font-bold text-orange-600 dark:text-orange-400">
          {{ stats.pending }}
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400">待完成</div>
      </div>
    </div>
  </div>
</template>
