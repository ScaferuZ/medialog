<template>
  <div class="min-h-screen bg-gradient-to-b from-slate-950 to-slate-900 text-slate-100">
    <AppHeader />
    <main class="mx-auto w-full max-w-6xl px-4 py-10 sm:px-6 lg:px-8">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>
