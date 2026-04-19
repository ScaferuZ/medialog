<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const { useLogs } = await import('~/composables/useLogs')

const { getMyLogs, createLog, updateLog } = useLogs()

const mediaId = computed(() => String(route.params.id))
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const existingLog = ref<any>(null)

const form = reactive({ status: 'plan', rating: null as number | null, note: '' })

const load = async () => {
  loading.value = true
  error.value = ''
  try {
    const logs = await getMyLogs()
    existingLog.value = Array.isArray(logs) ? logs.find((log: any) => String(log.media_id) === mediaId.value) : null
    if (existingLog.value) {
      form.status = existingLog.value.status ?? 'plan'
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
</script>

<template>
  <div class="mx-auto max-w-2xl px-4 py-10">
    <div class="mb-6">
      <p class="text-sm uppercase tracking-[0.25em] text-slate-500">Media log</p>
      <h1 class="mt-2 text-3xl font-bold text-slate-900">{{ existingLog ? 'Edit log' : 'Create log' }}</h1>
    </div>

    <div v-if="loading" class="rounded-2xl border border-slate-200 bg-white p-6 text-slate-600">Loading…</div>

    <div v-else class="space-y-4">
      <p v-if="error" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{{ error }}</p>
      <LogForm :status="form.status" :rating="form.rating" :note="form.note" :submitting="saving" @submit="handleSubmit" />
    </div>
  </div>
</template>
