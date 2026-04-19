<template>
  <div class="space-y-6">
    <section class="rounded-2xl border border-white/10 bg-white/5 p-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-sm uppercase tracking-[0.2em] text-cyan-300">Media</p>
          <h1 class="mt-2 text-3xl font-bold text-white">Browse media</h1>
          <p class="mt-2 text-sm text-slate-400">Discover titles from the local catalog and TMDB, then sync them into your library when you open them.</p>
        </div>

        <div class="grid gap-3 sm:grid-cols-2 lg:w-[32rem]">
          <label class="block">
            <span class="mb-1 block text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Search</span>
            <input
              v-model="query"
              type="search"
              placeholder="Search title or description"
              class="w-full rounded-xl border border-white/10 bg-slate-950/60 px-4 py-3 text-white outline-none transition placeholder:text-slate-500 focus:border-cyan-400/60"
            />
          </label>

          <label class="block">
            <span class="mb-1 block text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Type</span>
            <select
              v-model="typeFilter"
              class="w-full rounded-xl border border-white/10 bg-slate-950/60 px-4 py-3 text-white outline-none transition focus:border-cyan-400/60"
            >
              <option value="all">All</option>
              <option v-for="type in typeOptions" :key="type" :value="type">{{ type }}</option>
            </select>
          </label>
        </div>
      </div>
    </section>

    <div v-if="pending" class="rounded-2xl border border-white/10 bg-white/5 p-8 text-slate-400">Loading media...</div>
    <div v-else-if="error" class="rounded-2xl border border-red-500/20 bg-red-500/10 p-8 text-red-200">{{ error }}</div>

    <section v-else>
      <div class="mb-4 flex items-center justify-between gap-3 text-sm text-slate-400">
        <span>{{ sourceLabel }}</span>
        <span v-if="isTMDBSource" class="rounded-full border border-cyan-400/20 bg-cyan-400/10 px-3 py-1 text-cyan-200">
          Titles sync locally when you open them
        </span>
      </div>

      <div v-if="media.length" class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <button
          v-for="item in media"
          :key="itemKey(item)"
          type="button"
          class="group overflow-hidden rounded-2xl border border-white/10 bg-white/5 text-left transition hover:-translate-y-1 hover:border-cyan-400/40 hover:bg-white/10 disabled:cursor-wait disabled:opacity-70"
          :disabled="syncingKey === itemKey(item)"
          @click="openMedia(item)"
        >
          <div class="aspect-[2/3] overflow-hidden bg-slate-900">
            <img
              v-if="item.image_url || item.poster || item.poster_url"
              :src="item.image_url || item.poster || item.poster_url"
              :alt="item.title"
              class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
            />
            <div v-else class="flex h-full items-center justify-center text-sm text-slate-500">No image</div>
          </div>

          <div class="space-y-2 p-4">
            <div class="flex items-start justify-between gap-3">
              <h2 class="line-clamp-1 text-base font-semibold text-white">{{ item.title }}</h2>
              <span class="shrink-0 rounded-full border border-white/10 px-2 py-1 text-[11px] uppercase tracking-[0.2em] text-slate-400">
                {{ item.type || 'film' }}
              </span>
            </div>
            <p class="line-clamp-2 text-sm text-slate-400">{{ item.description || 'No description available.' }}</p>
            <div class="flex items-center justify-between text-sm text-slate-300">
              <span>{{ syncingKey === itemKey(item) ? 'Syncing…' : 'Rating' }}</span>
              <span class="text-cyan-300">{{ item.rating ?? '—' }}</span>
            </div>
          </div>
        </button>
      </div>

      <div v-else class="rounded-2xl border border-white/10 bg-white/5 p-8 text-slate-400">No media found.</div>
    </section>
  </div>
</template>

<script setup lang="ts">
type LocalMediaItem = {
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

type TMDBMovieItem = {
  id: number
  title: string
  overview?: string
  poster_url?: string
  vote_average?: number
  release_date?: string
}

type DisplayMediaItem = LocalMediaItem & {
  source: 'local' | 'tmdb'
  tmdb_id?: number
}

type SyncMovieResponse = {
  media?: {
    id: string
    ID?: string
  }
}

const router = useRouter()
const config = useRuntimeConfig()
const { listMedia, searchTMDB, getPopularTMDB, syncTMDBMovie } = useMedia()

const query = ref('')
const typeFilter = ref('all')
const media = ref<DisplayMediaItem[]>([])
const pending = ref(true)
const error = ref('')
const syncingKey = ref<string | null>(null)
const sourceLabel = ref('')
const tmdbEnabled = computed(() => config.public.tmdbEnabled)

const isTMDBSource = computed(() => media.value.some(item => item.source === 'tmdb'))
const typeOptions = computed(() => ['film'])

const itemKey = (item: DisplayMediaItem) => `${item.source}:${item.id}`

const normalizeLocalMediaItem = (item: LocalMediaItem): DisplayMediaItem => ({
  id: item.id || item.ID || '',
  title: item.title || item.Title || 'Untitled media',
  description: item.description || item.Description,
  type: item.type || item.Type,
  rating: item.rating,
  image_url: item.image_url,
  poster: item.poster,
  poster_url: item.poster_url || item.CoverImage,
  source: 'local'
})

const toTMDBDisplayItem = (item: TMDBMovieItem): DisplayMediaItem => ({
  id: `tmdb-${item.id}`,
  tmdb_id: item.id,
  title: item.title,
  description: item.overview,
  type: 'film',
  rating: item.vote_average,
  poster_url: item.poster_url,
  source: 'tmdb'
})

const loadLocalCatalog = async () => {
  const response = await listMedia() as { media?: LocalMediaItem[] }
  return (response.media ?? []).map(normalizeLocalMediaItem)
}

const filterLocalMedia = (items: DisplayMediaItem[], trimmedQuery: string) => {
  if (!trimmedQuery) {
    return items
  }

  const loweredQuery = trimmedQuery.toLowerCase()
  return items.filter((item) => `${item.title} ${item.description || ''}`.toLowerCase().includes(loweredQuery))
}

const loadMedia = async () => {
  pending.value = true
  error.value = ''

  try {
    const trimmedQuery = query.value.trim()
    const wantsFilmBrowse = typeFilter.value === 'all' || typeFilter.value === 'film'
    const localMedia = await loadLocalCatalog()

    if (trimmedQuery && wantsFilmBrowse && tmdbEnabled.value) {
      try {
        const response = await searchTMDB(trimmedQuery) as { movies?: TMDBMovieItem[] }
        media.value = (response.movies ?? []).map(toTMDBDisplayItem)
        sourceLabel.value = 'Showing TMDB search results'
      } catch {
        media.value = filterLocalMedia(localMedia, trimmedQuery)
        sourceLabel.value = 'TMDB search is unavailable, so these results are from your local catalog'
      }
      return
    }

    if (trimmedQuery) {
      media.value = filterLocalMedia(localMedia, trimmedQuery)
      sourceLabel.value = 'Showing results from your local catalog'
      return
    }

    if (wantsFilmBrowse && tmdbEnabled.value) {
      try {
        const tmdbResponse = await getPopularTMDB() as { movies?: TMDBMovieItem[] }
        media.value = (tmdbResponse.movies ?? []).map(toTMDBDisplayItem)
        sourceLabel.value = 'Showing popular films from TMDB'
        return
      } catch {
        if (localMedia.length > 0) {
          media.value = localMedia
          sourceLabel.value = 'TMDB browse is unavailable, so these results are from your local catalog'
          return
        }

        media.value = []
        sourceLabel.value = 'TMDB browse is unavailable right now. Add TMDB_API_KEY to enable external discovery.'
        return
      }
    }

    if (localMedia.length > 0 || !wantsFilmBrowse) {
      media.value = localMedia
      sourceLabel.value = localMedia.length > 0
        ? 'Showing your local media catalog'
        : 'No local media found for this filter'
      return
    }

    media.value = []
    sourceLabel.value = 'TMDB browse is disabled. Set NUXT_PUBLIC_TMDB_ENABLED=true and TMDB_API_KEY to enable external discovery.'
  } catch (e: any) {
    error.value = e?.message || 'Failed to load media.'
    media.value = []
  } finally {
    pending.value = false
  }
}

const openMedia = async (item: DisplayMediaItem) => {
  if (item.source === 'local') {
    await router.push(`/media/${item.id}`)
    return
  }

  if (!item.tmdb_id) {
    error.value = 'Missing TMDB id for this title.'
    return
  }

  syncingKey.value = itemKey(item)
  error.value = ''

  try {
    const response = await syncTMDBMovie(item.tmdb_id) as SyncMovieResponse
    const localMediaID = response.media?.id || response.media?.ID

    if (!localMediaID) {
      throw new Error('Failed to create a local media record.')
    }

    await router.push(`/media/${localMediaID}`)
  } catch (e: any) {
    error.value = e?.message || 'Failed to sync media from TMDB.'
  } finally {
    syncingKey.value = null
  }
}

await loadMedia()
watch([query, typeFilter], loadMedia)
</script>
