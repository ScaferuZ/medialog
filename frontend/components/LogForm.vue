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

const status = ref(props.status ?? 'planned')
const rating = ref<number | null>(props.rating ?? null)
const note = ref(props.note ?? '')

const statusOptions = [
  { value: 'planned', label: 'Planned', description: 'On your radar' },
  { value: 'in_progress', label: 'In progress', description: 'Currently watching' },
  { value: 'completed', label: 'Completed', description: 'Finished and rated' },
  { value: 'dropped', label: 'Dropped', description: 'Stopped for now' }
] as const

const setRating = (value: number) => {
  rating.value = value
}

const clearRating = () => {
  rating.value = null
}

const filledWidthForStar = (starIndex: number) => {
  if (rating.value == null) {
    return '0%'
  }

  const fullValue = starIndex * 2
  const halfValue = fullValue - 1

  if (rating.value >= fullValue) {
    return '100%'
  }

  if (rating.value === halfValue) {
    return '50%'
  }

  return '0%'
}

watch(() => props.status, value => { if (value !== undefined) status.value = value })
watch(() => props.rating, value => { if (value !== undefined) rating.value = value ?? null })
watch(() => props.note, value => { if (value !== undefined) note.value = value ?? '' })

const handleSubmit = () => emit('submit', { status: status.value, rating: rating.value, note: note.value })
</script>

<template>
  <form class="space-y-8" @submit.prevent="handleSubmit">
    <div class="space-y-3">
      <div>
        <p class="text-sm font-semibold uppercase tracking-[0.24em] text-cyan-300">Status</p>
        <p class="mt-1 text-sm text-slate-400">Set where you are in this title right now.</p>
      </div>

      <div class="grid gap-3 sm:grid-cols-2">
        <button
          v-for="option in statusOptions"
          :key="option.value"
          type="button"
          class="rounded-2xl border px-4 py-4 text-left transition"
          :class="status === option.value
            ? 'border-cyan-400/60 bg-cyan-400/12 shadow-[0_0_0_1px_rgba(34,211,238,0.2)]'
            : 'border-white/10 bg-slate-900/70 hover:border-white/20 hover:bg-slate-900'"
          @click="status = option.value"
        >
          <div class="font-semibold text-white">{{ option.label }}</div>
          <div class="mt-1 text-sm text-slate-400">{{ option.description }}</div>
        </button>
      </div>
    </div>

    <div class="space-y-3">
      <div class="flex items-center justify-between gap-4">
        <div>
          <p class="text-sm font-semibold uppercase tracking-[0.24em] text-cyan-300">Rating</p>
          <p class="mt-1 text-sm text-slate-400">Use half-stars for a 10-point score.</p>
        </div>
        <div class="text-right">
          <div class="text-2xl font-bold text-white">{{ rating ?? '—' }}<span class="text-sm font-medium text-slate-400"> / 10</span></div>
          <button type="button" class="mt-1 text-xs font-medium text-slate-400 hover:text-white" @click="clearRating">Clear</button>
        </div>
      </div>

      <div class="flex flex-wrap gap-2">
        <div v-for="starIndex in 5" :key="starIndex" class="relative h-12 w-12">
          <button type="button" class="absolute inset-y-0 left-0 z-10 w-1/2" :aria-label="`Rate ${starIndex * 2 - 1} out of 10`" @click="setRating(starIndex * 2 - 1)"></button>
          <button type="button" class="absolute inset-y-0 right-0 z-10 w-1/2" :aria-label="`Rate ${starIndex * 2} out of 10`" @click="setRating(starIndex * 2)"></button>
          <span class="absolute inset-0 text-center text-4xl leading-[3rem] text-slate-700">★</span>
          <span
            class="absolute inset-y-0 left-0 overflow-hidden text-center text-4xl leading-[3rem] text-cyan-300"
            :style="{ width: filledWidthForStar(starIndex) }"
          >★</span>
        </div>
      </div>
    </div>

    <div class="space-y-3">
      <div>
        <p class="text-sm font-semibold uppercase tracking-[0.24em] text-cyan-300">Review</p>
        <p class="mt-1 text-sm text-slate-400">Capture your quick thoughts, reaction, or mini review.</p>
      </div>
      <textarea
        v-model="note"
        rows="6"
        class="w-full rounded-2xl border border-white/10 bg-slate-900/80 px-4 py-3 text-base text-white placeholder:text-slate-500 focus:border-cyan-400/60 focus:outline-none focus:ring-2 focus:ring-cyan-400/20"
        placeholder="What stood out to you?"
      ></textarea>
    </div>

    <button
      type="submit"
      :disabled="submitting"
      class="inline-flex items-center rounded-xl bg-cyan-400 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-300 disabled:cursor-not-allowed disabled:opacity-60"
    >
      {{ submitting ? 'Saving...' : 'Save log' }}
    </button>
  </form>
</template>
