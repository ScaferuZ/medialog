<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const { useLogs } = await import('~/composables/useLogs')
const { getMedia } = useMedia()
const { isAuthenticated } = useAuth()

const { getMyLogs, createLog, updateLog } = useLogs()

const mediaId = computed(() => String(route.params.id))
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const existingLog = ref<any>(null)
const media = ref<any>(null)

const form = reactive({ status: 'planned', rating: null as number | null, note: '' })

const normalizeMedia = (item: any) => {
  if (!item) {
    return null
  }

  return {
    id: item.id || item.ID || '',
    title: item.title || item.Title || 'Untitled media',
    description: item.description || item.Description,
    poster: item.poster || item.poster_url || item.CoverImage || item.image_url,
    type: item.type || item.Type || 'film'
  }
}

const closeModal = async () => {
  await router.push(`/media/${mediaId.value}`)
}

const load = async () => {
  loading.value = true
  error.value = ''
  try {
    const [logResponse, mediaResponse] = await Promise.all([
      getMyLogs(),
      getMedia(mediaId.value)
    ])

    const logs = Array.isArray(logResponse)
      ? logResponse
      : Array.isArray(logResponse?.logs)
        ? logResponse.logs
        : []

    media.value = normalizeMedia(mediaResponse?.media)
    existingLog.value = logs.find((log: any) => String(log.media_id) === mediaId.value) ?? null
    if (existingLog.value) {
      form.status = existingLog.value.status ?? 'planned'
      form.rating = existingLog.value.rating ?? null
      form.note = existingLog.value.note ?? ''
    }
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to load log'
  } finally {
    loading.value = false
  }
}

const handleSubmit = async (payload: { status: string; rating: number | null; note: string }) => {
  saving.value = true
  error.value = ''
  try {
    const data = { media_id: mediaId.value, ...payload }
    if (existingLog.value?.id) await updateLog(existingLog.value.id, data)
    else await createLog(data)
    await router.push(`/media/${mediaId.value}`)
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to save log'
  } finally {
    saving.value = false
  }
}

onMounted(load)

if (!isAuthenticated.value) {
  await navigateTo('/login')
}
</script>

<template>
  <div class="fixed inset-0 z-50 overflow-y-auto bg-slate-950/85 backdrop-blur-sm">
    <div class="flex min-h-screen items-center justify-center p-4 sm:p-6" @click.self="closeModal">
      <div class="relative w-full max-w-5xl overflow-hidden rounded-[2rem] border border-cyan-400/20 bg-slate-950 shadow-[0_40px_120px_rgba(2,6,23,0.7)]">
        <button
          type="button"
          class="absolute right-4 top-4 z-10 rounded-full border border-white/10 bg-slate-900/80 px-3 py-2 text-sm font-medium text-slate-300 transition hover:border-white/20 hover:bg-slate-800 hover:text-white"
          @click="closeModal"
        >
          Close
        </button>

        <div class="grid lg:grid-cols-[280px_1fr]">
          <div class="border-b border-white/10 bg-gradient-to-br from-cyan-400/15 via-slate-950 to-slate-900 p-6 lg:border-b-0 lg:border-r">
            <div class="overflow-hidden rounded-[1.5rem] border border-white/10 bg-slate-900/80 shadow-2xl shadow-cyan-950/20">
              <img
                v-if="media?.poster"
                :src="media.poster"
                :alt="media.title"
                class="aspect-[2/3] w-full object-cover"
              />
              <div v-else class="flex aspect-[2/3] items-center justify-center text-slate-500">No image</div>
            </div>

            <div class="mt-6 space-y-3">
              <p class="text-xs font-semibold uppercase tracking-[0.28em] text-cyan-300">Media log</p>
              <h1 class="text-3xl font-bold tracking-tight text-white">{{ existingLog ? 'Update your log' : 'Create a new log' }}</h1>
              <p class="text-sm text-slate-400">Capture your progress, score it with half-stars, and leave a quick review without leaving the movie context.</p>
              <div class="rounded-2xl border border-white/10 bg-white/5 p-4">
                <p class="text-xs uppercase tracking-[0.24em] text-slate-500">Selected title</p>
                <p class="mt-2 text-lg font-semibold text-white">{{ media?.title || 'Loading title…' }}</p>
                <p class="mt-1 text-sm text-slate-400">{{ media?.type || 'film' }}</p>
              </div>
            </div>
          </div>

          <div class="p-6 sm:p-8">
            <div v-if="loading" class="rounded-2xl border border-white/10 bg-white/5 p-6 text-slate-300">Loading…</div>

            <div v-else class="space-y-4">
              <p v-if="error" class="rounded-2xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200">{{ error }}</p>
              <LogForm :status="form.status" :rating="form.rating" :note="form.note" :submitting="saving" @submit="handleSubmit" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
