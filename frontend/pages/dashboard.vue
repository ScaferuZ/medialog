<template>
  <div class="space-y-8">
    <section class="rounded-2xl border border-white/10 bg-white/5 p-8">
      <h1 class="text-3xl font-bold text-white">Welcome, {{ user?.username || 'User' }}</h1>
      <p class="mt-2 text-slate-400">Track and organize your media journey.</p>
    </section>
    
    <section class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <div class="rounded-xl border border-white/10 bg-white/5 p-6">
        <h3 class="text-lg font-semibold text-white">Your Logs</h3>
        <p class="mt-1 text-sm text-slate-400">Start logging your media</p>
        <NuxtLink to="/media" class="mt-4 inline-block text-cyan-400 hover:text-cyan-300">
          Browse Media →
        </NuxtLink>
      </div>
      
      <div class="rounded-xl border border-white/10 bg-white/5 p-6">
        <h3 class="text-lg font-semibold text-white">Timeline</h3>
        <p class="mt-1 text-sm text-slate-400">See what friends are watching</p>
        <NuxtLink to="/timeline" class="mt-4 inline-block text-cyan-400 hover:text-cyan-300">
          View Timeline →
        </NuxtLink>
      </div>
      
      <div class="rounded-xl border border-white/10 bg-white/5 p-6">
        <h3 class="text-lg font-semibold text-white">Profile</h3>
        <p class="mt-1 text-sm text-slate-400">Manage your account</p>
        <NuxtLink :to="`/${user?.username}`" class="mt-4 inline-block text-cyan-400 hover:text-cyan-300">
          View Profile →
        </NuxtLink>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { user, isAuthenticated } = useAuth()

if (!isAuthenticated.value) {
  navigateTo('/login')
}
</script>
