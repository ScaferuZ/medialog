<template>
  <div class="flex min-h-[60vh] items-center justify-center">
    <div class="w-full max-w-md rounded-2xl border border-white/10 bg-white/5 p-8 shadow-xl">
      <h1 class="mb-6 text-center text-2xl font-bold text-white">Create Account</h1>
      
      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label for="username" class="mb-1 block text-sm font-medium text-slate-300">Username</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            minlength="3"
            class="w-full rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-white placeholder-slate-400 focus:border-cyan-400 focus:outline-none focus:ring-1 focus:ring-cyan-400"
            placeholder="Choose a username"
          />
        </div>
        
        <div>
          <label for="email" class="mb-1 block text-sm font-medium text-slate-300">Email</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            class="w-full rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-white placeholder-slate-400 focus:border-cyan-400 focus:outline-none focus:ring-1 focus:ring-cyan-400"
            placeholder="Enter your email"
          />
        </div>
        
        <div>
          <label for="password" class="mb-1 block text-sm font-medium text-slate-300">Password</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            minlength="8"
            class="w-full rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-white placeholder-slate-400 focus:border-cyan-400 focus:outline-none focus:ring-1 focus:ring-cyan-400"
            placeholder="Create a password (min 8 characters)"
          />
        </div>
        
        <div>
          <label for="confirmPassword" class="mb-1 block text-sm font-medium text-slate-300">Confirm Password</label>
          <input
            id="confirmPassword"
            v-model="form.confirmPassword"
            type="password"
            required
            class="w-full rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-white placeholder-slate-400 focus:border-cyan-400 focus:outline-none focus:ring-1 focus:ring-cyan-400"
            placeholder="Confirm your password"
          />
        </div>
        
        <div v-if="validationError" class="rounded-lg bg-red-500/10 px-4 py-2 text-sm text-red-400">
          {{ validationError }}
        </div>
        
        <div v-if="error" class="rounded-lg bg-red-500/10 px-4 py-2 text-sm text-red-400">
          {{ error }}
        </div>
        
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full rounded-lg bg-cyan-500 px-4 py-2 font-medium text-slate-900 transition-colors hover:bg-cyan-400 disabled:cursor-not-allowed disabled:opacity-50"
        >
          <span v-if="isLoading">Creating account...</span>
          <span v-else>Create Account</span>
        </button>
      </form>
      
      <p class="mt-6 text-center text-sm text-slate-400">
        Already have an account?
        <NuxtLink to="/login" class="text-cyan-400 hover:text-cyan-300">Sign in</NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
const { register, isLoading, error } = useAuth()
const router = useRouter()

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validationError = ref('')

const handleRegister = async () => {
  validationError.value = ''
  
  if (form.password !== form.confirmPassword) {
    validationError.value = 'Passwords do not match'
    return
  }
  
  if (form.password.length < 8) {
    validationError.value = 'Password must be at least 8 characters'
    return
  }
  
  try {
    await register({
      username: form.username,
      email: form.email,
      password: form.password
    })
    router.push('/dashboard')
  } catch {
    // Error is already handled in the composable
  }
}
</script>
