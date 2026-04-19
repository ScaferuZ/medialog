<template>
  <div class="min-h-screen bg-slate-950 px-4 py-10 text-slate-100 sm:px-6 lg:px-8">
    <div class="mx-auto max-w-3xl">
      <h1 class="mb-8 text-3xl font-bold">Timeline</h1>
      
      <div v-if="pending" class="text-center py-12">
        <p class="text-slate-400">Loading...</p>
      </div>
      
      <div v-else-if="error" class="rounded-2xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200">
        {{ error }}
      </div>
      
      <div v-else-if="logs.length === 0" class="text-center py-12 text-slate-400">
        <p>No activity yet. Follow some users to see their logs!</p>
      </div>
      
      <div v-else class="space-y-4">
        <div
          v-for="log in logs"
          :key="log.id"
          class="rounded-2xl border border-white/10 bg-white/5 p-6"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-full bg-gradient-to-br from-emerald-500 to-blue-500 flex items-center justify-center text-sm font-bold">
                {{ log.username?.[0]?.toUpperCase() || '?' }}
              </div>
              <div>
                <p class="font-medium">{{ log.username || 'Unknown' }}</p>
                <p class="text-sm text-slate-400">{{ new Date(log.created_at).toLocaleDateString() }}</p>
              </div>
            </div>
            <span
              :class="{
                'rounded-full px-3 py-1 text-xs font-medium': true,
                'bg-emerald-500/20 text-emerald-300': log.status === 'completed',
                'bg-blue-500/20 text-blue-300': log.status === 'in_progress',
                'bg-yellow-500/20 text-yellow-300': log.status === 'planned',
                'bg-red-500/20 text-red-300': log.status === 'dropped'
              }"
            >
              {{ log.status?.replace('_', ' ') }}
            </span>
          </div>
          
          <div v-if="log.rating" class="mt-3 flex items-center gap-1 text-yellow-400">
            <span>★</span>
            <span>{{ log.rating }}/10</span>
          </div>
          
          <p v-if="log.note" class="mt-3 text-slate-300">{{ log.note }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { getTimeline } = useLogs()

const { data: logs, pending, error } = await useAsyncData('timeline', async () => {
  try {
    const response = await getTimeline()
    return response.logs || []
  } catch (e: any) {
    throw new Error(e?.message || 'Failed to load timeline')
  }
}, {
  server: false
})
</script>