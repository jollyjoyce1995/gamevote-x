<template>
  <main class="min-h-[calc(100vh-120px)] flex flex-col items-center justify-center px-6 py-8">
    <div class="w-full max-w-4xl">
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-20"><div class="spinner"></div></div>

    <!-- Error -->
    <div v-else-if="error" class="card text-center py-10">
      <p class="text-xl mb-4">❌ {{ error }}</p>
      <RouterLink to="/" class="btn btn-ghost">← Back</RouterLink>
    </div>

    <template v-else-if="party">
      <!-- Party Header -->
      <div class="flex items-center justify-between mb-6 flex-wrap gap-3">
        <div>
          <div class="flex items-center gap-3 mb-1">
            <h1 class="text-2xl font-extrabold">Party <span class="grad-text font-mono">{{ party.code }}</span></h1>
            <span :class="statusBadge(party.status!)" class="badge">{{ party.status }}</span>
          </div>
          <div class="flex gap-4 text-sm" style="color:var(--c-muted)">
            <span>👥 {{ (party.attendees || []).join(', ') }}</span>
          </div>
        </div>
        <button
          v-if="canAdvance"
          class="btn btn-primary"
          :disabled="advancing"
          @click="advance"
        >
          <span v-if="advancing">...</span>
          <span v-else>{{ advanceLabel }}</span>
        </button>
      </div>

      <!-- ─── PHASE COMPONENTS ──────────────────── -->
      <PartyNomination v-if="party.status === 'NOMINATION'" />
      <PartyVoting v-else-if="party.status === 'VOTING'" />
      <PartyResults v-else-if="party.status === 'RESULTS'" />
      
    </template>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePartyStore } from '@/stores/party'
import { openPartyStream } from '@/api'
import { partiesApi } from '@/api-client'
import type { ServicePartyDTO as PartyDTO } from '@/generated-api'

import PartyNomination from '@/components/party/PartyNomination.vue'
import PartyVoting from '@/components/party/PartyVoting.vue'
import PartyResults from '@/components/party/PartyResults.vue'

const route = useRoute()
const authStore = useAuthStore()
const partyStore = usePartyStore()

const code = computed(() => route.params.code as string)
const party = computed(() => partyStore.party)

const loading = ref(true)
const error = ref('')
const advancing = ref(false)

// ── Computed ────────────────────────────────────────────────────
const canAdvance = computed(() => {
  if (!party.value) return false
  if (party.value.status === 'NOMINATION') return (party.value.options?.length || 0) > 0
  if (party.value.status === 'VOTING') return true
  return false
})

const advanceLabel = computed(() => {
  if (party.value?.status === 'NOMINATION') return '▶ Start Voting'
  if (party.value?.status === 'VOTING') return '▶ Show Results'
  return ''
})

function statusBadge(status: string) {
  return {
    'badge-nomination': status === 'NOMINATION',
    'badge-voting': status === 'VOTING',
    'badge-results': status === 'RESULTS',
  }
}

// ── SSE ─────────────────────────────────────────────────────────
function connectSSE() {
  const es = openPartyStream(code.value, authStore.username!)
  partyStore.setEventSource(es)

  es.addEventListener('party_updated', (e: MessageEvent) => {
    partyStore.setParty(JSON.parse(e.data) as PartyDTO)
  })

  es.addEventListener('online_users', (e: MessageEvent) => {
    partyStore.setOnlineUsers(JSON.parse(e.data) as string[])
  })

  es.addEventListener('outstanding_voters_updated', (e: MessageEvent) => {
    partyStore.setOutstandingVoters(JSON.parse(e.data) as string[])
  })

  es.onerror = () => {
    // quietly ignore – EventSource will reconnect
  }
}

// ── Phase actions ────────────────────────────────────────────────
async function advance() {
  if (!party.value) return
  advancing.value = true
  try {
    const nextStatus = party.value.status === 'NOMINATION' ? 'VOTING' : 'RESULTS'
    const updated = await partiesApi.partiesCodePatch({
      code: code.value,
      patchReq: { status: nextStatus }
    })
    partyStore.setParty(updated)
  } finally {
    advancing.value = false
  }
}

// ── Lifecycle ────────────────────────────────────────────────────
onMounted(async () => {
  try {
    let p = await partiesApi.partiesCodeGet({ code: code.value })
    
    // Auto-join logic
    if (authStore.username && (!p.attendees || !p.attendees.includes(authStore.username))) {
      await partiesApi.partiesCodeAttendeesPost({
        code: code.value,
        value: { value: authStore.username }
      })
      // re-fetch party to get the updated attendees list immediately
      p = await partiesApi.partiesCodeGet({ code: code.value })
    }

    partyStore.setParty(p)
    connectSSE()
  } catch (e: any) {
    console.error('Party load error:', e)
    error.value = e.message || 'Party not found'
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  partyStore.closeStream()
})
</script>
