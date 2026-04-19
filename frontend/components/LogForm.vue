<script setup lang="ts">
const props = defineProps<{
  status?: string
  rating?: number | null
  note?: string
  submitting?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: { status: string; rating: number | null; note: string }]
}>()

const status = ref(props.status ?? 'plan')
const rating = ref<number | null>(props.rating ?? null)
const note = ref(props.note ?? '')

watch(() => props.status, value => { if (value !== undefined) status.value = value })
watch(() => props.rating, value => { if (value !== undefined) rating.value = value ?? null })
watch(() => props.note, value => { if (value !== undefined) note.value = value ?? '' })

const handleSubmit = () => emit('submit', { status: status.value, rating: rating.value, note: note.value })
</script>

<template>
  <form class="space-y-5 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm" @submit.prevent="handleSubmit">
    <div>
      <label class="mb-2 block text-sm font-medium text-slate-700">Status</label>
      <select v-model="status" class="w-full rounded-lg border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-200">
        <option value="plan">Plan</option>
        <option value="watching">Watching</option>
        <option value="completed">Completed</option>
        <option value="dropped">Dropped</option>
      </select>
    </div>

    <div>
      <label class="mb-2 block text-sm font-medium text-slate-700">Rating</label>
      <input v-model.number="rating" type="number" min="0" max="10" step="1" class="w-full rounded-lg border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-200" placeholder="Optional" />
    </div>

    <div>
      <label class="mb-2 block text-sm font-medium text-slate-700">Note</label>
      <textarea v-model="note" rows="5" class="w-full rounded-lg border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-200" placeholder="Add your thoughts..."></textarea>
    </div>

    <button type="submit" :disabled="submitting" class="inline-flex items-center rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60">
      {{ submitting ? 'Saving...' : 'Save log' }}
    </button>
  </form>
</template>
