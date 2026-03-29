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
      <PartyCard
        v-for="p in parties"
        :key="p.code"
        :party="p"
      />
    </TransitionGroup>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { getParties } from '@/api-client'
import type { ServicePartyDto as PartyDTO } from '@/generated-api'
import PartyCard from '@/components/PartyCard.vue'

const parties = ref<PartyDTO[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await getParties()
    const data = (res.data as PartyDTO[]) || []
    // Sort newest first by createdAt
    parties.value = data.sort((a, b) => {
      const da = (a as any).createdAt ? new Date((a as any).createdAt).getTime() : 0
      const db = (b as any).createdAt ? new Date((b as any).createdAt).getTime() : 0
      return db - da
    })
  } finally {
    loading.value = false
  }
})
</script>
