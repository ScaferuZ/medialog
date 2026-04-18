<template>
  <div class="flex min-h-[60vh] items-center justify-center">
    <div class="w-full max-w-md rounded-2xl border border-white/10 bg-white/5 p-8 shadow-xl">
      <h1 class="mb-6 text-center text-2xl font-bold text-white">Sign In</h1>
      
      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label for="username" class="mb-1 block text-sm font-medium text-slate-300">Username</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            class="w-full rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-white placeholder-slate-400 focus:border-cyan-400 focus:outline-none focus:ring-1 focus:ring-cyan-400"
            placeholder="Enter your username"
          />
        </div>
        
        <div>
          <label for="password" class="mb-1 block text-sm font-medium text-slate-300">Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            class="w-full rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-white placeholder-slate-400 focus:border-cyan-400 focus:outline-none focus:ring-1 focus:ring-cyan-400"
            placeholder="Enter your password"
          />
        </div>
        
        <div v-if="error" class="rounded-lg bg-red-500/10 px-4 py-2 text-sm text-red-400">
          {{ error }}
        </div>
        
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full rounded-lg bg-cyan-500 px-4 py-2 font-medium text-slate-900 transition-colors hover:bg-cyan-400 disabled:cursor-not-allowed disabled:opacity-50"
        >
          <span v-if="isLoading">Signing in...</span>
          <span v-else>Sign In</span>
        </button>
      </form>
      
      <p class="mt-6 text-center text-sm text-slate-400">
        Don't have an account?
        <NuxtLink to="/register" class="text-cyan-400 hover:text-cyan-300">Sign up</NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
const { login, isLoading, error } = useAuth()
const router = useRouter()

const form = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  try {
    await login({
      username: form.username,
      password: form.password
    })
    router.push('/dashboard')
  } catch {
    // Error is already handled in the composable
  }
}
</script>
