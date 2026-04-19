<template>
  <div v-if="pending" class="rounded-2xl border border-white/10 bg-white/5 p-8 text-slate-400">Loading media...</div>
  <div v-else-if="error || !media" class="rounded-2xl border border-red-500/20 bg-red-500/10 p-8 text-red-200">Media not found.</div>

  <article v-else class="grid gap-8 lg:grid-cols-[320px_1fr]">
    <div class="overflow-hidden rounded-3xl border border-white/10 bg-white/5">
      <img v-if="media.image_url || media.poster || media.poster_url" :src="media.image_url || media.poster || media.poster_url"
        :alt="media.title" class="aspect-[2/3] w-full object-cover" />
      <div v-else class="flex aspect-[2/3] items-center justify-center text-slate-500">No image</div>
    </div>

    <section class="space-y-5 rounded-3xl border border-white/10 bg-white/5 p-6">
      <div>
        <p class="text-sm uppercase tracking-[0.2em] text-cyan-300">{{ media.type || 'Media' }}</p>
        <h1 class="mt-2 text-4xl font-bold text-white">{{ media.title }}</h1>
      </div>

      <p class="text-slate-300">{{ media.description || 'No description available.' }}</p>

      <div class="flex flex-wrap gap-3 text-sm">
        <div class="rounded-full border border-white/10 px-4 py-2 text-slate-300">Rating: <span class="text-cyan-300">{{ media.rating ?? '—' }}</span></div>
      </div>

      <button class="rounded-xl bg-cyan-400 px-5 py-3 font-semibold text-slate-950 transition hover:bg-cyan-300">
        Log This
      </button>
    </section>
  </article>
</template>

<script setup lang="ts">
type MediaItem = {
  id: string
  ID?: string
  title: string
  Title?: string
  description?: string
  Description?: string
  type?: string
  Type?: string
  rating?: number | string
  image_url?: string
  poster?: string
  poster_url?: string
  CoverImage?: string
}

const normalizeMediaItem = (item: MediaItem | null | undefined): MediaItem | null => {
  if (!item) {
    return null
  }

  return {
    id: item.id || item.ID || '',
    title: item.title || item.Title || 'Untitled media',
    description: item.description || item.Description,
    type: item.type || item.Type,
    rating: item.rating,
    image_url: item.image_url,
    poster: item.poster,
    poster_url: item.poster_url || item.CoverImage
  }
}

const route = useRoute()
const { getMedia } = useMedia()

const { data, pending, error } = await useAsyncData<MediaItem>(`media-${route.params.id}`, async () => {
  const response = await getMedia(String(route.params.id)) as { media?: MediaItem }
  return normalizeMediaItem(response.media)
})

const media = computed(() => data.value)
</script>
