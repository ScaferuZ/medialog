<template>
  <div class="space-y-6">
    <section class="rounded-2xl border border-white/10 bg-white/5 p-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-sm uppercase tracking-[0.2em] text-cyan-300">Media</p>
          <h1 class="mt-2 text-3xl font-bold text-white">Browse media</h1>
          <p class="mt-2 text-sm text-slate-400">Search and filter your collection.</p>
        </div>

        <div class="grid gap-3 sm:grid-cols-2 lg:w-[32rem]">
          <label class="block">
            <span class="mb-1 block text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Search</span>
            <input v-model="query" type="search" placeholder="Search title or description"
              class="w-full rounded-xl border border-white/10 bg-slate-950/60 px-4 py-3 text-white outline-none transition placeholder:text-slate-500 focus:border-cyan-400/60" />
          </label>

          <label class="block">
            <span class="mb-1 block text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Type</span>
            <select v-model="typeFilter"
              class="w-full rounded-xl border border-white/10 bg-slate-950/60 px-4 py-3 text-white outline-none transition focus:border-cyan-400/60">
              <option value="all">All</option>
              <option v-for="type in typeOptions" :key="type" :value="type">{{ type }}</option>
            </select>
          </label>
        </div>
      </div>
    </section>

    <div v-if="pending" class="rounded-2xl border border-white/10 bg-white/5 p-8 text-slate-400">Loading media...</div>
    <div v-else-if="error" class="rounded-2xl border border-red-500/20 bg-red-500/10 p-8 text-red-200">Failed to load media.</div>

    <section v-else>
      <div v-if="filteredMedia.length" class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <NuxtLink v-for="item in filteredMedia" :key="item.id" :to="`/media/${item.id}`"
          class="group overflow-hidden rounded-2xl border border-white/10 bg-white/5 transition hover:-translate-y-1 hover:border-cyan-400/40 hover:bg-white/10">
          <div class="aspect-[2/3] overflow-hidden bg-slate-900">
            <img v-if="item.image_url || item.poster || item.poster_url" :src="item.image_url || item.poster || item.poster_url"
              :alt="item.title" class="h-full w-full object-cover transition duration-500 group-hover:scale-105" />
            <div v-else class="flex h-full items-center justify-center text-sm text-slate-500">No image</div>
          </div>

          <div class="space-y-2 p-4">
            <div class="flex items-start justify-between gap-3">
              <h2 class="line-clamp-1 text-base font-semibold text-white">{{ item.title }}</h2>
              <span class="shrink-0 rounded-full border border-white/10 px-2 py-1 text-[11px] uppercase tracking-[0.2em] text-slate-400">
                {{ item.type || 'media' }}
              </span>
            </div>
            <p class="line-clamp-2 text-sm text-slate-400">{{ item.description || 'No description available.' }}</p>
            <div class="flex items-center justify-between text-sm text-slate-300">
              <span>Rating</span>
              <span class="text-cyan-300">{{ item.rating ?? '—' }}</span>
            </div>
          </div>
        </NuxtLink>
      </div>

      <div v-else class="rounded-2xl border border-white/10 bg-white/5 p-8 text-slate-400">No media found.</div>
    </section>
  </div>
</template>

<script setup lang="ts">
type MediaItem = {
  id: string
  title: string
  description?: string
  type?: string
  rating?: number | string
  image_url?: string
  poster?: string
  poster_url?: string
}

const { listMedia } = useMedia()
const query = ref('')
const typeFilter = ref('all')

const { data, pending, error } = await useAsyncData<MediaItem[]>('media-list', async () => {
  const response = await listMedia()
  return Array.isArray(response) ? response : (response as { items?: MediaItem[] })?.items || []
})

const media = computed(() => data.value ?? [])
const typeOptions = computed(() => [...new Set(media.value.map(item => item.type).filter(Boolean))] as string[])

const filteredMedia = computed(() => {
  const q = query.value.trim().toLowerCase()
  return media.value.filter((item) => {
    const matchesQuery = !q || `${item.title} ${item.description || ''}`.toLowerCase().includes(q)
    const matchesType = typeFilter.value === 'all' || item.type === typeFilter.value
    return matchesQuery && matchesType
  })
})
</script>
