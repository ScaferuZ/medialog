<template>
  <div class="space-y-8">
    <section class="rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl shadow-black/20 sm:p-12">
      <p class="mb-3 text-sm font-medium uppercase tracking-[0.2em] text-cyan-300">Welcome</p>
      <h1 class="text-4xl font-bold tracking-tight text-white sm:text-5xl">
        Medialogg - Track your media
      </h1>
      <p class="mt-4 max-w-2xl text-lg leading-8 text-slate-300">
        Log, review, and organize everything you watch, read, and play in one place.
      </p>
    </section>

    <section class="rounded-3xl border border-white/10 bg-white/5 p-6 sm:p-8">
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p class="text-sm font-medium uppercase tracking-[0.2em] text-cyan-300">Latest public activity</p>
          <h2 class="mt-2 text-2xl font-semibold text-white">Recently logged titles</h2>
        </div>
      </div>

      <div v-if="pending" class="mt-6 text-slate-400">Loading latest activity...</div>
      <div v-else-if="error" class="mt-6 rounded-2xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-200">{{ error.message }}</div>
      <div v-else-if="activity.length === 0" class="mt-6 text-slate-400">No public activity yet.</div>

      <div v-else class="mt-6 grid gap-4 lg:grid-cols-2">
        <article
          v-for="item in activity"
          :key="item.id"
          class="overflow-hidden rounded-2xl border border-white/10 bg-slate-950/60"
        >
          <div class="grid grid-cols-[96px_1fr] gap-4 p-4">
            <div class="overflow-hidden rounded-2xl bg-white/5">
              <img v-if="item.coverImage" :src="item.coverImage" :alt="item.title" class="aspect-[2/3] w-full object-cover" />
              <div v-else class="flex aspect-[2/3] items-center justify-center text-xs text-slate-500">No image</div>
            </div>

            <div class="space-y-3">
              <div>
                <p class="text-xs uppercase tracking-[0.2em] text-cyan-300">{{ item.type }}</p>
                <NuxtLink :to="`/media/${item.mediaId}`" class="mt-1 block text-lg font-semibold text-white transition hover:text-cyan-300">{{ item.title }}</NuxtLink>
                <p class="mt-1 text-sm text-slate-400">Logged by <NuxtLink :to="`/${item.username}`" class="text-white transition hover:text-cyan-300">{{ item.displayName || item.username }}</NuxtLink> · {{ formatDate(item.createdAt) }}</p>
              </div>

              <div class="flex flex-wrap gap-3 text-sm">
                <div class="rounded-full border border-white/10 px-3 py-1 text-slate-300">Status: <span class="text-white">{{ formatStatus(item.status) }}</span></div>
                <div v-if="item.rating" class="rounded-full border border-white/10 px-3 py-1 text-slate-300">Rating: <span class="text-cyan-300">{{ item.rating }}</span></div>
              </div>

              <p v-if="item.note" class="line-clamp-3 text-sm text-slate-300">{{ item.note }}</p>
            </div>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
type ActivityItem = {
  id: string
  ID?: string
  mediaId: string
  MediaID?: string
  status: string
  Status?: string
  rating?: number | string | null
  Rating?: number | string | null
  note?: string
  Note?: string
  createdAt?: string
  CreatedAt?: string
  username: string
  Username?: string
  displayName?: string
  DisplayName?: string
  title: string
  Title?: string
  coverImage?: string
  CoverImage?: string
  type?: string
  Type?: string
}

const { getLatestActivity } = useActivity()

const normalizeActivityItem = (item: ActivityItem) => ({
  id: item.id || item.ID || '',
  mediaId: item.mediaId || item.MediaID || '',
  status: item.status || item.Status || 'completed',
  rating: item.rating ?? item.Rating ?? null,
  note: item.note || item.Note || '',
  createdAt: item.createdAt || item.CreatedAt || '',
  username: item.username || item.Username || 'unknown',
  displayName: item.displayName || item.DisplayName || '',
  title: item.title || item.Title || 'Untitled media',
  coverImage: item.coverImage || item.CoverImage || '',
  type: item.type || item.Type || 'media'
})

const formatStatus = (value: string) => {
  const labels: Record<string, string> = {
    planned: 'Planned',
    in_progress: 'In Progress',
    completed: 'Completed',
    dropped: 'Dropped'
  }

  return labels[value] || value
}

const formatDate = (value?: string) => {
  if (!value) {
    return 'Just now'
  }

  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) {
    return value
  }

  return parsed.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const { data: activity, pending, error } = await useAsyncData('latest-activity', async () => {
  const response = await getLatestActivity() as { activity?: ActivityItem[] }
  return (response.activity ?? []).map(normalizeActivityItem)
})
</script>
