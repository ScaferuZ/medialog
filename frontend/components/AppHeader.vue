<template>
  <header class="border-b border-white/10 bg-slate-950/80 backdrop-blur">
    <div class="mx-auto flex w-full max-w-6xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
      <NuxtLink to="/" class="block">
        <p class="text-lg font-semibold tracking-tight text-white">Medialogg</p>
        <p class="text-sm text-slate-400">Track your media</p>
      </NuxtLink>
      
      <nav class="flex items-center gap-4">
        <template v-if="isAuthenticated">
          <NuxtLink to="/dashboard" class="text-sm text-slate-300 hover:text-white">
            Dashboard
          </NuxtLink>
          <NuxtLink to="/media" class="text-sm text-slate-300 hover:text-white">
            Media
          </NuxtLink>
          <NuxtLink to="/timeline" class="text-sm text-slate-300 hover:text-white">
            Timeline
          </NuxtLink>
          <div class="flex items-center gap-2 pl-4 border-l border-white/10">
            <span class="text-sm text-slate-400">{{ user?.username }}</span>
            <button
              @click="logout"
              class="rounded-lg bg-white/10 px-3 py-1.5 text-sm text-white hover:bg-white/20"
            >
              Logout
            </button>
          </div>
        </template>
        
        <template v-else>
          <NuxtLink to="/login" class="text-sm text-slate-300 hover:text-white">
            Sign In
          </NuxtLink>
          <NuxtLink
            to="/register"
            class="rounded-lg bg-cyan-500 px-4 py-2 text-sm font-medium text-slate-900 hover:bg-cyan-400"
          >
            Sign Up
          </NuxtLink>
        </template>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
const { user, isAuthenticated, logout } = useAuth()
</script>
