<template>
  <main class="min-h-[calc(100vh-120px)] flex flex-col items-center justify-center px-6 py-8">
    <div class="w-full max-w-4xl">
      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-20">
        <div class="spinner"></div>
      </div>

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
              <span v-if="party.currentRound" class="badge" style="background:var(--c-surface);color:var(--c-text)">Round {{ party.currentRound }}</span>
            </div>
            <div class="flex gap-4 text-sm" style="color:var(--c-muted)">
              <span>👥 {{ (party.attendees || []).join(', ') }}</span>
            </div>
          </div>
        </div>

        <!-- Tabs -->
        <div class="flex gap-2 mb-6 border-b" style="border-color: var(--c-border)">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            class="tab-btn"
            :class="{ 'tab-btn-active': activeTab === tab.key }"
            @click="activeTab = tab.key"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- Tab Panels -->
        <template v-if="activeTab === 'NOMINATION'">
          <PartyNomination/>
          <div class="mt-4 flex justify-end">
            <button
              class="btn btn-primary"
              :disabled="advancing || (party.options?.length ?? 0) === 0"
              @click="advance('VOTING')"
            >
              <span v-if="advancing">...</span>
              <span v-else>▶ Start Voting</span>
            </button>
          </div>
        </template>
        <template v-else-if="activeTab === 'VOTING'">
          <PartyVoting/>
          <div class="mt-4 flex justify-end">
            <button
              class="btn btn-primary"
              :disabled="advancing"
              @click="advance('RESULTS')"
            >
              <span v-if="advancing">...</span>
              <span v-else>▶ Show Results</span>
            </button>
          </div>
        </template>
        <template v-else-if="activeTab === 'RESULTS'">
          <PartyResults/>
          <div class="mt-4 flex justify-end">
            <button
              class="btn btn-primary"
              :disabled="restarting"
              @click="newRound"
            >
              <span v-if="restarting">...</span>
              <span v-else>🔄 New Round</span>
            </button>
          </div>
        </template>

      </template>
    </div>
  </main>
</template>

<script setup lang="ts">
import {ref, computed, watch, onMounted, onUnmounted} from 'vue'
import {RouterLink, useRoute} from 'vue-router'
import {useAuthStore} from '@/stores/auth'
import {usePartyStore} from '@/stores/party'
import { patchParty, getParty, postAttendee, client } from '@/api-client'
import type {ServicePartyDto as PartyDTO} from '@/generated-api'

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
const restarting = ref(false)

// ── Tab state ────────────────────────────────────────────────────
const tabs = [
  { key: 'NOMINATION' as const, label: '🎮 Nomination' },
  { key: 'VOTING'     as const, label: '👍 Voting' },
  { key: 'RESULTS'   as const, label: '🏆 Results' },
]
const activeTab = ref<'NOMINATION' | 'VOTING' | 'RESULTS'>('NOMINATION')

// Auto-switch tab when party status changes via SSE
watch(() => party.value?.status, (newStatus) => {
  if (newStatus) activeTab.value = newStatus as typeof activeTab.value
})

// ── Helpers ──────────────────────────────────────────────────────
function statusBadge(status: string) {
  return {
    'badge-nomination': status === 'NOMINATION',
    'badge-voting': status === 'VOTING',
    'badge-results': status === 'RESULTS',
  }
}

// ── SSE ─────────────────────────────────────────────────────────
function connectSSE() {
  const url = client.buildUrl({
    url: '/parties/{code}/stream',
    path: { code: code.value },
    query: { username: authStore.username! }
  })
  
  const es = new EventSource(url)
  partyStore.setEventSource(es as any)

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

// ── Phase advance ────────────────────────────────────────────────
async function advance(nextStatus: 'VOTING' | 'RESULTS') {
  if (!party.value) return
  advancing.value = true
  try {
    const updated = await patchParty({
      path: { code: code.value },
      body: { status: nextStatus }
    })
    partyStore.setParty(updated.data as PartyDTO)
    activeTab.value = nextStatus
  } finally {
    advancing.value = false
  }
}

// ── New Round ────────────────────────────────────────────────────
async function newRound() {
  if (!party.value) return
  restarting.value = true
  try {
    const updated = await patchParty({
      path: { code: code.value },
      body: { status: 'NOMINATION' }
    })
    partyStore.setParty(updated.data as PartyDTO)
    activeTab.value = 'NOMINATION'
  } finally {
    restarting.value = false
  }
}

// ── Lifecycle ────────────────────────────────────────────────────
onMounted(async () => {
  try {
    const pReq = await getParty({ path: { code: code.value } })
    let p = pReq.data as PartyDTO;

    // Auto-join logic
    if (authStore.username && (!p.attendees || !p.attendees.includes(authStore.username))) {
      await postAttendee({
        path: { code: code.value },
        body: { value: authStore.username }
      })
      // re-fetch party to get the updated attendees list immediately
      const p2 = await getParty({ path: { code: code.value } })
      p = p2.data as PartyDTO;
    }

    partyStore.setParty(p)
    // Set initial tab to match current party status
    if (p.status) activeTab.value = p.status as typeof activeTab.value
    await connectSSE()
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

<style scoped>
.tab-btn {
  padding: 0.5rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 600;
  border-bottom: 3px solid transparent;
  transition: color 0.2s, border-color 0.2s;
  color: var(--c-muted);
  background: transparent;
  cursor: pointer;
  margin-bottom: -1px;
}
.tab-btn:hover {
  color: var(--c-text);
}
.tab-btn-active {
  color: var(--c-primary);
  border-bottom-color: var(--c-primary);
}
</style>
