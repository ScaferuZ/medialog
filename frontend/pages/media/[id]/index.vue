<template>
  <div v-if="pending" class="rounded-2xl border border-white/10 bg-white/5 p-8 text-slate-400">Loading media...</div>
  <div v-else-if="error || !media" class="rounded-2xl border border-red-500/20 bg-red-500/10 p-8 text-red-200">Media not found.</div>

  <div v-else class="space-y-6">
    <div v-if="showSavedLogBanner" class="flex items-start justify-between gap-4 rounded-2xl border border-emerald-400/30 bg-emerald-400/10 px-5 py-4 text-emerald-100">
      <div>
        <p class="text-sm font-semibold uppercase tracking-[0.22em] text-emerald-300">Saved</p>
        <p class="mt-1 text-sm text-emerald-100/90">Your log was saved and is now attached to this title.</p>
      </div>
      <button type="button" class="text-sm font-medium text-emerald-200 transition hover:text-white" @click="showSavedLogBanner = false">Dismiss</button>
    </div>

    <article class="grid gap-8 lg:grid-cols-[320px_1fr]">
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
          <div v-if="userLog" class="rounded-full border border-emerald-400/20 bg-emerald-400/10 px-4 py-2 text-emerald-100">Logged: <span class="font-semibold">{{ formatStatus(userLog.status) }}</span></div>
        </div>

        <button
          class="rounded-xl bg-cyan-400 px-5 py-3 font-semibold text-slate-950 transition hover:bg-cyan-300"
          @click="handleLogThis"
        >
          {{ userLog ? 'Edit Log' : 'Log This' }}
        </button>
      </section>
    </article>

    <section v-if="userLog || communityReviews.length" class="rounded-3xl border border-white/10 bg-white/5 p-6">
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p class="text-sm uppercase tracking-[0.2em] text-cyan-300">Reviews & activity</p>
          <h2 class="mt-2 text-2xl font-semibold text-white">What people have said about this title</h2>
        </div>

        <div class="flex flex-wrap gap-3 text-sm text-slate-300">
          <div v-if="userLog" class="rounded-full border border-emerald-400/20 bg-emerald-400/10 px-4 py-2 text-emerald-100">You logged this</div>
          <div class="rounded-full border border-white/10 px-4 py-2">Reviews: <span class="text-white">{{ reviewCount }}</span></div>
        </div>
      </div>

      <div class="mt-6 grid gap-4 lg:grid-cols-2">
        <article v-if="userLog" class="rounded-2xl border border-emerald-400/20 bg-emerald-400/8 p-5">
          <div class="flex flex-wrap items-start justify-between gap-3">
            <div>
              <p class="text-sm font-semibold uppercase tracking-[0.22em] text-emerald-300">Your review</p>
              <h3 class="mt-2 text-xl font-semibold text-white">{{ formatStatus(userLog.status) }}</h3>
            </div>

            <div class="flex flex-wrap gap-3 text-sm">
              <div class="rounded-full border border-white/10 px-4 py-2 text-slate-300">Rating: <span class="text-cyan-300">{{ userLog.rating ?? '—' }}</span></div>
              <div class="rounded-full border border-white/10 px-4 py-2 text-slate-300">Updated: <span class="text-white">{{ formatDate(userLog.updated_at || userLog.created_at) }}</span></div>
            </div>
          </div>

          <p class="mt-4 whitespace-pre-wrap text-slate-200">{{ userLog.note || 'You logged this title, but haven’t written a review yet.' }}</p>
        </article>

        <article v-for="review in communityReviews" :key="review.id" class="rounded-2xl border border-white/10 bg-slate-950/60 p-5">
          <div class="flex flex-wrap items-start justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-white">{{ review.displayName || review.username }}</p>
              <p class="mt-1 text-xs uppercase tracking-[0.18em] text-slate-500">{{ review.username }}</p>
            </div>

            <div class="flex flex-wrap gap-3 text-sm">
              <div class="rounded-full border border-white/10 px-4 py-2 text-slate-300">Rating: <span class="text-cyan-300">{{ review.rating ?? '—' }}</span></div>
              <div class="rounded-full border border-white/10 px-4 py-2 text-slate-300">Updated: <span class="text-white">{{ formatDate(review.updatedAt || review.createdAt) }}</span></div>
            </div>
          </div>

          <p class="mt-4 whitespace-pre-wrap text-slate-300">{{ review.content || 'No review text provided.' }}</p>
        </article>
      </div>
    </section>
  </div>
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
  Metadata?: string
  TmdbID?: string
}

type UserLog = {
  id: string
  ID?: string
  status: string
  Status?: string
  rating?: number | null
  Rating?: number | null
  note?: string
  Note?: string
  created_at?: string
  CreatedAt?: string
  updated_at?: string
  UpdatedAt?: string
}

type MediaReview = {
  id: string
  ID?: string
  logId?: string
  LogID?: string
  username: string
  Username?: string
  displayName?: string
  DisplayName?: string
  content?: string
  Content?: string
  rating?: number | string | null
  Rating?: number | string | null
  createdAt?: string
  CreatedAt?: string
  updatedAt?: string
  UpdatedAt?: string
}

const decodeBase64JSON = (value?: string) => {
  if (!value) {
    return null
  }

  try {
    const decoded = process.server
      ? Buffer.from(value, 'base64').toString('utf-8')
      : window.atob(value)

    return JSON.parse(decoded) as Record<string, unknown>
  } catch {
    return null
  }
}

const normalizeMediaItem = (item: MediaItem | null | undefined): MediaItem | null => {
  if (!item) {
    return null
  }

  const metadata = decodeBase64JSON(item.Metadata)
  const metadataRating = typeof metadata?.vote_average === 'number' ? metadata.vote_average : undefined

  return {
    id: item.id || item.ID || '',
    title: item.title || item.Title || 'Untitled media',
    description: item.description || item.Description,
    type: item.type || item.Type,
    rating: item.rating ?? metadataRating,
    image_url: item.image_url,
    poster: item.poster,
    poster_url: item.poster_url || item.CoverImage,
    TmdbID: item.TmdbID
  }
}

const route = useRoute()
const router = useRouter()
const { getMedia, getMediaReviews } = useMedia()
const { isAuthenticated, token, user } = useAuth()
const { getLogForMedia } = useLogs()

const showSavedLogBanner = ref(route.query.savedLog === '1')

const normalizeUserLog = (item: UserLog | null | undefined) => {
  if (!item) {
    return null
  }

  return {
    id: item.id || item.ID || '',
    status: item.status || item.Status || 'planned',
    rating: item.rating ?? item.Rating ?? null,
    note: item.note || item.Note || '',
    created_at: item.created_at || item.CreatedAt || '',
    updated_at: item.updated_at || item.UpdatedAt || ''
  }
}

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

const normalizeReview = (item: MediaReview | null | undefined) => {
  if (!item) {
    return null
  }

  return {
    id: item.id || item.ID || '',
    logId: item.logId || item.LogID || '',
    username: item.username || item.Username || 'unknown',
    displayName: item.displayName || item.DisplayName || '',
    content: item.content || item.Content || '',
    rating: item.rating ?? item.Rating ?? null,
    createdAt: item.createdAt || item.CreatedAt || '',
    updatedAt: item.updatedAt || item.UpdatedAt || ''
  }
}

const reviewCount = computed(() => communityReviews.value.length + (userLog.value ? 1 : 0))

const { data, pending, error } = await useAsyncData<MediaItem>(`media-${route.params.id}`, async () => {
  const response = await getMedia(String(route.params.id)) as { media?: MediaItem }
  return normalizeMediaItem(response.media)
})

const media = computed(() => data.value)
const { data: userLogData } = await useAsyncData(`media-log-${route.params.id}`, async () => {
  if (!token.value) {
    return null
  }

  try {
    const response = await getLogForMedia(String(route.params.id)) as UserLog
    return normalizeUserLog(response)
  } catch {
    return null
  }
})

const userLog = computed(() => userLogData.value)
const { data: reviewsData } = await useAsyncData(`media-reviews-${route.params.id}`, async () => {
  try {
    const response = await getMediaReviews(String(route.params.id)) as { reviews?: MediaReview[] }
    return (response.reviews ?? [])
      .map(normalizeReview)
      .filter((review): review is NonNullable<ReturnType<typeof normalizeReview>> => !!review)
  } catch {
    return []
  }
})

const reviews = computed(() => reviewsData.value ?? [])

const communityReviews = computed(() => {
  if (!userLog.value) {
    return reviews.value
  }

  return reviews.value.filter((review) => review.logId !== userLog.value?.id && review.username !== user.value?.username)
})

const handleLogThis = async () => {
  if (!isAuthenticated.value) {
    await router.push('/login')
    return
  }

  await router.push(`/media/${route.params.id}/log`)
}

onMounted(() => {
  if (route.query.savedLog === '1') {
    const nextQuery = { ...route.query }
    delete nextQuery.savedLog
    router.replace({ query: nextQuery })
  }
})
</script>
