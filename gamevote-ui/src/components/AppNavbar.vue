<template>
  <nav class="glass sticky top-0 z-50 px-6 py-5 flex items-center justify-between gap-4">
    <!-- Left: Logo + Party code -->
    <div class="flex items-center gap-4">
      <RouterLink to="/" class="text-xl font-extrabold grad-text tracking-tight">🎲 GameVote</RouterLink>
      <template v-if="partyCode">
        <span class="text-muted text-sm">›</span>
        <span class="text-sm font-semibold text-indigo-300">{{ partyCode }}</span>
        <RouterLink :to="`/parties/${partyCode}/drinks`" class="btn btn-ghost btn-sm ml-2">🍺 Drinks</RouterLink>
      </template>
    </div>

    <!-- Center: Online users in party -->
    <div v-if="partyCode && onlineUsers.length" class="flex items-center gap-2">
      <span class="text-xs text-muted">Online:</span>
      <div class="flex gap-1">
        <span
          v-for="u in onlineUsers"
          :key="u"
          class="badge badge-nomination"
        >{{ u }}</span>
      </div>
    </div>

    <!-- Right: Username + logout -->
    <div class="flex items-center gap-3">
      <span class="text-sm text-muted">👤 <strong class="text-text">{{ authStore.username }}</strong></span>
      <button class="btn btn-ghost btn-sm" @click="handleLogout">Logout</button>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePartyStore } from '@/stores/party'

const authStore = useAuthStore()
const partyStore = usePartyStore()
const route = useRoute()
const router = useRouter()

const partyCode = computed(() => route.params.code as string | undefined)
const onlineUsers = computed(() => partyStore.onlineUsers)

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
