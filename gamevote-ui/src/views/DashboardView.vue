<template>
  <main class="min-h-[calc(100vh-120px)] flex flex-col items-center justify-center px-6 py-10">
    <div class="w-full max-w-5xl">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-extrabold">
          <span class="grad-text">Parties</span>
        </h1>
        <p class="text-sm mt-1" style="color:var(--c-muted)">Pick a party to join or create a new one</p>
      </div>
      <RouterLink to="/parties/new" class="btn btn-primary">
        ✨ New Party
      </RouterLink>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-20">
      <div class="spinner"></div>
    </div>

    <!-- Empty -->
    <div v-else-if="parties.length === 0" class="card text-center py-16">
      <div class="text-5xl mb-4">🎉</div>
      <p class="text-lg font-semibold mb-2">No parties yet!</p>
      <p class="text-sm mb-6" style="color:var(--c-muted)">Be the first to create one</p>
      <RouterLink to="/parties/new" class="btn btn-primary">Create Party</RouterLink>
    </div>

    <!-- Party list -->
    <TransitionGroup name="list" tag="div" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <RouterLink
        v-for="p in parties"
        :key="p.code"
        :to="`/parties/${p.code}`"
        class="card group cursor-pointer no-underline"
        style="text-decoration:none"
      >
        <div class="flex items-start justify-between mb-4">
          <div>
            <div class="text-xs font-mono mb-1" style="color:var(--c-muted)">{{ p.code }}</div>
            <div class="flex items-center gap-2">
              <span :class="statusBadge(p.status)" class="badge">{{ p.status }}</span>
            </div>
          </div>
          <span class="text-2xl">🎮</span>
        </div>
        <div class="flex gap-4 text-sm" style="color:var(--c-muted)">
          <span>👥 {{ p.attendees.length }}</span>
          <span>📋 {{ p.options.length }} games</span>
          <span>🍺 {{ p.beerCount }}</span>
        </div>
      </RouterLink>
    </TransitionGroup>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { getParties, type PartyDTO } from '@/api'

const parties = ref<PartyDTO[]>([])
const loading = ref(true)

function statusBadge(status: string) {
  return {
    'badge-nomination': status === 'NOMINATION',
    'badge-voting': status === 'VOTING',
    'badge-results': status === 'RESULTS',
  }
}

onMounted(async () => {
  try {
    parties.value = await getParties()
  } finally {
    loading.value = false
  }
})
</script>
