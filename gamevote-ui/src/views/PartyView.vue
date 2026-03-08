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
            <span :class="statusBadge(party.status)" class="badge">{{ party.status }}</span>
          </div>
          <div class="flex gap-4 text-sm" style="color:var(--c-muted)">
            <span>👥 {{ party.attendees.join(', ') }}</span>
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

      <!-- ─── NOMINATION PHASE ──────────────────── -->
      <section v-if="party.status === 'NOMINATION'">
        <div class="card mb-4">
          <h2 class="text-lg font-bold mb-4">🎮 Nominate Games</h2>

          <!-- Steam search -->
          <div class="flex gap-2 mb-3">
            <input
              v-model="gameSearch"
              class="input"
              placeholder="Search Steam games..."
              @input="onSearchInput"
            />
            <button class="btn btn-primary" @click="submitCustomGame">Add Custom</button>
          </div>

          <!-- Steam suggestions -->
          <div v-if="suggestions.length" class="rounded-xl overflow-hidden mb-4" style="border:1px solid var(--c-border); background: var(--c-bg-card)">
            <div
              v-for="g in suggestions"
              :key="g.appId"
              class="flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors hover:bg-[var(--c-border)]"
              @click="pickGame(g)"
            >
              <img v-if="g.imageUrl" :src="g.imageUrl" class="w-16 h-8 object-cover rounded shadow-sm" alt="" />
              <div class="flex-1">
                <p class="text-sm font-semibold">{{ g.name }}</p>
                <p class="text-xs" style="color:var(--c-muted)">Steam AppID {{ g.appId }}</p>
              </div>
              <button class="btn btn-primary btn-sm">Nominate</button>
            </div>
          </div>

          <!-- Nominated list -->
          <div class="mt-2">
            <p class="text-sm font-semibold mb-3" style="color:var(--c-muted)">Nominated Games ({{ party.options?.length || 0 }})</p>
            <TransitionGroup name="list" tag="div" class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div
                v-for="optName in party.options"
                :key="optName"
                class="flex items-center gap-3 px-3 py-3 rounded-2xl relative group overflow-hidden"
                style="background:var(--c-bg-card);border:1px solid var(--c-border)"
              >
                <div class="w-16 h-8 flex items-center justify-center rounded bg-gray-800 text-[10px] relative z-10" style="color:var(--c-muted)">GAME</div>
                
                <span class="text-sm font-bold flex-1 truncate relative z-10">{{ optName }}</span>
                
                <button 
                  class="btn btn-ghost btn-sm text-error opacity-0 group-hover:opacity-100 transition-opacity relative z-10 p-2"
                  title="Remove nomination"
                  @click="removeOption(optName)"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </TransitionGroup>
            <p v-if="!party.options?.length" class="text-sm text-center py-8" style="color:var(--c-muted)">No games nominated yet. Use the search bar above to add some!</p>
          </div>
        </div>
      </section>

      <!-- ─── VOTING PHASE ──────────────────────── -->
      <section v-else-if="party.status === 'VOTING'">
        <div v-if="alreadyVoted" class="card text-center py-8">
          <p class="text-lg mb-2">✅ You've voted!</p>
          <p class="text-sm" style="color:var(--c-muted)">Waiting for others: {{ outstanding.join(', ') || 'All done!' }}</p>
        </div>
        <div v-else class="card">
          <h2 class="text-lg font-bold mb-6">👍👎 Vote on Games</h2>
          <div class="flex flex-col gap-4">
            <div
              v-for="optName in party.options"
              :key="optName"
              class="flex items-center gap-3 p-3 rounded-xl"
              style="background:var(--c-bg);border:1px solid var(--c-border)"
            >
              <span class="font-medium flex-1">{{ optName }}</span>
              <div class="flex gap-2">
                <button
                  class="vote-btn like"
                  :class="{ active: currentVotes[optName] === 1 }"
                  @click="setVote(optName, 1)"
                >👍</button>
                <button
                  class="vote-btn dislike"
                  :class="{ active: currentVotes[optName] === -1 }"
                  @click="setVote(optName, -1)"
                >👎</button>
              </div>
            </div>
          </div>
          <button class="btn btn-primary w-full mt-6" :disabled="submittingVote" @click="submitVotes">
            {{ submittingVote ? 'Submitting...' : '✅ Submit Votes' }}
          </button>
        </div>
        <!-- Outstanding voters -->
        <div v-if="outstanding.length" class="card mt-4">
          <p class="text-sm font-semibold mb-2" style="color:var(--c-muted)">Still voting:</p>
          <div class="flex gap-2 flex-wrap">
            <span v-for="u in outstanding" :key="u" class="badge badge-voting">{{ u }}</span>
          </div>
        </div>
      </section>

      <!-- ─── RESULTS PHASE ─────────────────────── -->
      <section v-else-if="party.status === 'RESULTS'">
        <div class="card">
          <h2 class="text-lg font-bold mb-6">🏆 Results</h2>
          <div class="flex flex-col gap-3 mb-8">
            <div
              v-for="(score, gameName) in sortedResults"
              :key="gameName"
              class="flex items-center gap-4"
            >
              <span class="font-medium flex-1 text-sm">{{ gameName }}</span>
              <div class="flex-1 h-3 rounded-full overflow-hidden" style="background:var(--c-border)">
                <div
                  class="h-full rounded-full transition-all"
                  :style="{ width: barWidth(score) + '%', background: score > 0 ? 'var(--c-success)' : score < 0 ? 'var(--c-danger)' : 'var(--c-muted)' }"
                ></div>
              </div>
              <span class="text-sm font-bold w-8 text-right" :style="{ color: score > 0 ? 'var(--c-success)' : score < 0 ? 'var(--c-danger)' : 'var(--c-muted)' }">{{ score > 0 ? '+' : '' }}{{ score }}</span>
            </div>
          </div>
          <!-- Pie Chart -->
          <canvas ref="pieChartEl" style="max-height:260px"></canvas>
        </div>
      </section>
    </template>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePartyStore } from '@/stores/party'
import {
  searchGames,
  openPartyStream,
  type Game
} from '@/api'
import { partiesApi, pollsApi } from '@/api-client'
import type { ServicePartyDTO as PartyDTO } from '@/generated-api'
import { Chart, ArcElement, Tooltip, Legend, PieController } from 'chart.js'
Chart.register(ArcElement, Tooltip, Legend, PieController)

const route = useRoute()
const authStore = useAuthStore()
const partyStore = usePartyStore()

const code = computed(() => route.params.code as string)
const party = computed(() => partyStore.party)

const loading = ref(true)
const error = ref('')
const advancing = ref(false)
const gameSearch = ref('')
const suggestions = ref<Game[]>([])
const currentVotes = ref<Record<string, number>>({})
const submittingVote = ref(false)
const alreadyVoted = ref(false)
const outstanding = ref<string[]>([])
const pieChartEl = ref<HTMLCanvasElement | null>(null)
let pieChart: Chart | null = null
let searchTimeout: ReturnType<typeof setTimeout> | null = null

// ── Computed ────────────────────────────────────────────────────
const sortedResults = computed(() => {
  if (!party.value?.results) return {}
  return Object.fromEntries(
    Object.entries(party.value.results).sort(([, a], [, b]) => b - a)
  )
})

const maxScore = computed(() =>
  Math.max(1, ...Object.values(party.value?.results ?? {}), 1)
)

function barWidth(score: number) {
  return Math.max(2, (Math.abs(score) / maxScore.value) * 100)
}

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
    if (party.value?.status === 'VOTING') loadOutstanding()
    if (party.value?.status === 'RESULTS') nextTick(() => renderPieChart())
  })

  es.addEventListener('online_users', (e: MessageEvent) => {
    partyStore.setOnlineUsers(JSON.parse(e.data) as string[])
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

// ── Nomination ───────────────────────────────────────────────────
function onSearchInput() {
  if (searchTimeout) clearTimeout(searchTimeout)
  if (!gameSearch.value.trim()) { suggestions.value = []; return }
  searchTimeout = setTimeout(async () => {
    suggestions.value = await searchGames(gameSearch.value).catch(() => [])
  }, 300)
}

async function pickGame(item: Game | string) {
  const name = typeof item === 'string' ? item.trim() : item.name
  if (!name) return

  // Check for duplicates
  if (party.value?.options?.some(optName => optName.toLowerCase() === name.toLowerCase())) {
    alert('This game is already nominated!')
    return
  }

  const option = typeof item === 'string' 
    ? { name } 
    : { name: item.name, appId: item.appId, imageUrl: item.imageUrl }

  gameSearch.value = ''
  suggestions.value = []
  
  try {
    const res = await partiesApi.partiesCodeOptionsPost({
      code: code.value,
      value: { value: name }
    })
    // Note: The UI expects the option object to be returned or reactive update via SSE
  } catch (e) {
    console.error('Failed to nominate:', e)
  }
}

async function submitCustomGame() {
  const name = gameSearch.value.trim()
  if (!name) return
  await pickGame(name)
}

async function removeOption(name: string) {
  // The generated API uses index for delete, but the old one used name.
  // Checking the swagger... it says "Delete an option from a party by index"
  // but the path is /parties/{code}/options/{optionId}.
  // In the manual api.ts it was encodeURIComponent(name).
  // Let's check the swagger index type... it says integer.
  // This might be a breaking change if the backend actually expects index now.
  const index = party.value?.options?.findIndex(o => o === name)
  if (index === undefined || index === -1) return

  await partiesApi.partiesCodeOptionsOptionIdDelete({
    code: code.value,
    optionId: index
  }).catch((e) => {
    console.error('Failed to remove option:', e)
  })
}

// ── Voting ───────────────────────────────────────────────────────
function setVote(game: string, val: number) {
  currentVotes.value[game] = currentVotes.value[game] === val ? 0 : val
}

async function submitVotes() {
  if (!authStore.username) return
  submittingVote.value = true
  try {
    const pollId = party.value?.links?.['poll']?.href?.split('/').pop()
    if (!pollId) return
    await pollsApi.pollsIdVotesAttendeePut({
      id: pollId,
      attendee: authStore.username,
      choices: currentVotes.value
    })
    alreadyVoted.value = true
    await loadOutstanding()
  } finally {
    submittingVote.value = false
  }
}

async function loadOutstanding() {
  const pollId = party.value?.links?.['poll']?.href?.split('/').pop()
  if (!pollId) return
  outstanding.value = await pollsApi.pollsIdOutstandingGet({ id: pollId }).catch(() => [])
  if (!outstanding.value.includes(authStore.username!)) {
    alreadyVoted.value = true
  }
}

// ── Results Pie Chart ────────────────────────────────────────────
function renderPieChart() {
  if (!pieChartEl.value || !party.value?.results) return
  const labels = Object.keys(party.value.results)
  const data = Object.values(party.value.results).map(v => Math.max(0, v))

  if (pieChart) pieChart.destroy()
  pieChart = new Chart(pieChartEl.value, {
    type: 'pie',
    data: {
      labels,
      datasets: [{
        data,
        backgroundColor: ['#6366f1','#8b5cf6','#22c55e','#f59e0b','#ef4444','#06b6d4','#ec4899'],
        borderColor: '#161b26',
        borderWidth: 3,
      }]
    },
    options: {
      plugins: {
        legend: { labels: { color: '#e2e8f0', font: { size: 12 } } },
      }
    }
  })
}

// ── Lifecycle ────────────────────────────────────────────────────
onMounted(async () => {
  try {
    const p = await partiesApi.partiesCodeGet({ code: code.value })
    partyStore.setParty(p)
    connectSSE()
    if (p.status === 'VOTING') await loadOutstanding()
    if (p.status === 'RESULTS') nextTick(() => renderPieChart())
  } catch (e: any) {
    console.error('Party load error:', e)
    error.value = e.message || 'Party not found'
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  partyStore.closeStream()
  if (pieChart) pieChart.destroy()
})

watch(() => party.value?.status, (newStatus) => {
  if (newStatus === 'RESULTS') nextTick(() => renderPieChart())
  if (newStatus === 'VOTING') {
    alreadyVoted.value = false
    currentVotes.value = {}
    loadOutstanding()
  }
})
</script>
