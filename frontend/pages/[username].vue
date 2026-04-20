<template>
  <div class="min-h-screen bg-slate-950 px-4 py-10 text-slate-100 sm:px-6 lg:px-8">
    <div class="mx-auto max-w-3xl">
      <div class="overflow-hidden rounded-3xl border border-white/10 bg-white/5 shadow-2xl shadow-black/30 backdrop-blur">
        <div class="h-28 bg-gradient-to-r from-emerald-500 via-cyan-500 to-blue-500"></div>

        <div class="px-6 pb-6 sm:px-8">
          <div class="-mt-12 flex items-end justify-between gap-4">
            <div class="flex items-end gap-4">
              <div class="flex h-24 w-24 items-center justify-center rounded-2xl border-4 border-slate-950 bg-slate-800 text-3xl font-bold text-white shadow-lg">
                <img
                  v-if="profile?.avatar_url"
                  :src="profile.avatar_url"
                  :alt="profile.username"
                  class="h-full w-full rounded-2xl object-cover"
                />
                <span v-else>{{ initials }}</span>
              </div>

              <div class="pb-2">
                <h1 class="text-2xl font-semibold tracking-tight">{{ profile?.username || username }}</h1>
                <p v-if="profile?.display_name" class="mt-1 text-sm font-medium text-slate-300">{{ profile.display_name }}</p>
                <p class="mt-1 text-sm text-slate-400">{{ profile?.bio || 'No bio yet.' }}</p>
              </div>
            </div>

            <div v-if="!hasMounted" class="h-10 w-24 rounded-full border border-white/10 bg-white/5"></div>

            <div v-else-if="isOwnProfile" class="rounded-full border border-cyan-400/30 bg-cyan-400/10 px-4 py-2 text-sm font-medium text-cyan-200">
              Your profile
            </div>

            <button
              v-else
              class="rounded-full bg-white px-5 py-2 text-sm font-semibold text-slate-950 transition hover:scale-[1.02] hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="following"
              @click="onFollow"
            >
              {{ following ? 'Following' : 'Follow' }}
            </button>
          </div>

          <div class="mt-8 grid grid-cols-2 gap-3 sm:grid-cols-4">
            <div v-for="item in statsCards" :key="item.label" class="rounded-2xl border border-white/10 bg-slate-900/60 p-4">
              <div class="text-2xl font-semibold">{{ item.value }}</div>
              <div class="mt-1 text-sm text-slate-400">{{ item.label }}</div>
            </div>
          </div>

          <p v-if="error" class="mt-6 rounded-2xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200">
            {{ error }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const username = computed(() => String(route.params.username || ''))
const { getUserProfile, getUserStats, getFollowStatus, followUser } = useUsers()
const { user, isAuthenticated } = useAuth()

interface UserProfile {
  id: string
  username: string
  display_name?: string
  bio?: string
  avatar_url?: string
  is_public: boolean
  created_at: string
}

interface UserStats {
  completed_count: number
  in_progress_count: number
  followers_count: number
  following_count: number
}

const profile = ref<UserProfile | null>(null)
const stats = ref({ completed: 0, inProgress: 0, followers: 0, following: 0 })
const error = ref('')
const following = ref(false)
const hasMounted = ref(false)
const isOwnProfile = computed(() => isAuthenticated.value && user.value?.username === username.value)

const initials = computed(() => (profile.value?.username || username.value || '?').slice(0, 1).toUpperCase())
const statsCards = computed(() => [
  { label: 'Completed', value: stats.value.completed },
  { label: 'In Progress', value: stats.value.inProgress },
  { label: 'Followers', value: stats.value.followers },
  { label: 'Following', value: stats.value.following }
])

const syncFollowState = async () => {
  if (!isAuthenticated.value || isOwnProfile.value) {
    following.value = false
    return
  }

  try {
    const response = await getFollowStatus(username.value) as { following?: boolean }
    following.value = !!response.following
  } catch {
    following.value = false
  }
}

await useAsyncData(`user-profile-${username.value}`, async () => {
  try {
    const [userData, statsData] = await Promise.all([
      getUserProfile(username.value),
      getUserStats(username.value)
    ]) as [UserProfile, UserStats]

    profile.value = userData
    stats.value = {
      completed: statsData?.completed_count ?? 0,
      inProgress: statsData?.in_progress_count ?? 0,
      followers: statsData?.followers_count ?? 0,
      following: statsData?.following_count ?? 0
    }

    await syncFollowState()
  } catch (e: any) {
    error.value = e?.message || 'Failed to load profile.'
  }
})

watch([isAuthenticated, isOwnProfile, username], () => {
  void syncFollowState()
})

onMounted(() => {
  hasMounted.value = true
})

const onFollow = async () => {
  if (isOwnProfile.value) {
    return
  }

  try {
    await followUser(username.value)
    following.value = true
    stats.value.followers += 1
  } catch (e: any) {
    error.value = e?.message || 'Failed to follow user.'
  }
}
</script>
