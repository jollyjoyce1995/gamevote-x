<template>
  <section>
    <section>
      <div class="card">
        <div class="flex justify-between items-center mb-6 flex-wrap gap-2">
          <h2 class="text-lg font-bold">🏆 Results</h2>
          <div v-if="availableRounds.length > 0" class="flex gap-2 bg-black/20 p-1 rounded-full overflow-x-auto">
            <button
                v-for="r in availableRounds"
                :key="r"
                class="px-4 py-1.5 text-sm font-semibold rounded-full transition-all whitespace-nowrap"
                :class="selectedRound === r ? 'shadow-lg bg-surface-1 text-primary' : 'text-muted hover:text-white'"
                :style="selectedRound === r ? 'background:var(--c-surface-1);color:var(--c-primary)' : ''"
                @click="selectedRound = r"
            >
              Round {{ r }}
            </button>
          </div>
        </div>
        <div class="flex flex-col gap-4 mb-8">
          <div
              v-for="entry in sortedResults"
              :key="entry.name"
              class="flex items-center gap-4 p-2 rounded-xl"
              style="background: var(--c-surface, rgba(255,255,255,0.03))"
          >
            <!-- Game thumbnail -->
            <img
                v-if="entry.imageUrl"
                :src="entry.imageUrl"
                :alt="entry.name"
                class="w-16 h-10 object-cover rounded-lg flex-shrink-0 shadow"
            />
            <div v-else class="w-16 h-10 rounded-lg flex-shrink-0 flex items-center justify-center text-2xl"
                 style="background:var(--c-border)">
              🎮
            </div>

            <!-- Name + AppID -->
            <div class="min-w-0 w-32 flex-shrink-0">
              <p class="text-sm font-semibold truncate">{{ entry.name }}</p>
              <p v-if="entry.appId" class="text-xs" style="color:var(--c-muted)">AppID {{ entry.appId }}</p>
            </div>

            <!-- Score bar -->
            <div class="flex-1 h-3 rounded-full overflow-hidden" style="background:var(--c-border)">
              <div
                  class="h-full rounded-full transition-all"
                  :style="{ width: barWidth(entry.score) + '%', background: entry.score > 0 ? 'var(--c-success)' : entry.score < 0 ? 'var(--c-danger)' : 'var(--c-muted)' }"
              ></div>
            </div>

            <!-- Score label -->
            <span
                class="text-sm font-bold w-8 text-right flex-shrink-0"
                :style="{ color: entry.score > 0 ? 'var(--c-success)' : entry.score < 0 ? 'var(--c-danger)' : 'var(--c-muted)' }"
            >{{ entry.score > 0 ? '+' : '' }}{{ entry.score }}</span>
          </div>
        </div>
        <!-- Pie Chart -->
        <canvas ref="pieChartEl" style="max-height:260px"></canvas>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import {ref, computed, onMounted, onUnmounted, watch, nextTick} from 'vue'
import {usePartyStore} from '@/stores/party'
import {Chart, ArcElement, Tooltip, Legend, PieController} from 'chart.js'

Chart.register(ArcElement, Tooltip, Legend, PieController)

const partyStore = usePartyStore()
const party = computed(() => partyStore.party)

const pieChartEl = ref<HTMLCanvasElement | null>(null)
let pieChart: Chart | null = null

const selectedRound = ref<number | null>(null)

const availableRounds = computed(() => {
  if (!party.value?.roundResults) return party.value?.currentRound ? [party.value.currentRound] : []
  const keys = Object.keys(party.value.roundResults).map(Number).sort((a, b) => a - b)
  return keys.length > 0 ? keys : party.value?.currentRound ? [party.value.currentRound] : []
})

watch(() => party.value?.currentRound, (newVal) => {
  if (selectedRound.value === null && newVal) {
    selectedRound.value = newVal
  }
}, {immediate: true})

const currentRoundResults = computed(() => {
  if (!party.value) return {}
  const round = selectedRound.value || party.value.currentRound || 1
  if (party.value.roundResults && party.value.roundResults[round]) {
    return party.value.roundResults[round]
  }
  return party.value.results || {}
})

const sortedResults = computed(() => {
  if (!party.value?.options) return []

  const resultsObj = currentRoundResults.value

  // Create an array of results with full game details
  const results = party.value.options.map(opt => {
    return {
      name: opt.name || '',
      appId: opt.appId,
      imageUrl: opt.imageUrl,
      score: resultsObj[opt.name || ''] ?? 0
    }
  })

  // Optionally handle results for games not in options (fallback)
  const inOptions = new Set(party.value.options.map(o => o.name))
  Object.entries(resultsObj).forEach(([name, score]) => {
    if (!inOptions.has(name)) {
      results.push({name, score: score as number, appId: undefined, imageUrl: undefined})
    }
  })

  return results
})

const maxScore = computed(() =>
    Math.max(1, ...Object.values(currentRoundResults.value as Record<string, number>), 1)
)

function barWidth(score: number) {
  return Math.max(2, (Math.abs(score) / maxScore.value) * 100)
}

function renderPieChart() {
  if (!pieChartEl.value || !currentRoundResults.value) return
  const resultsObj = currentRoundResults.value
  const labels = Object.keys(resultsObj)
  const data = Object.values(resultsObj).map(v => Math.max(0, v as number))

  if (pieChart) pieChart.destroy()
  pieChart = new Chart(pieChartEl.value, {
    type: 'pie',
    data: {
      labels,
      datasets: [{
        data,
        backgroundColor: ['#6366f1', '#8b5cf6', '#22c55e', '#f59e0b', '#ef4444', '#06b6d4', '#ec4899'],
        borderColor: '#161b26',
        borderWidth: 3,
      }]
    },
    options: {
      plugins: {
        legend: {labels: {color: '#e2e8f0', font: {size: 12}}},
      }
    }
  })
}

onMounted(() => {
  if (party.value?.status === 'RESULTS') {
    nextTick(() => renderPieChart())
  }
})

onUnmounted(() => {
  if (pieChart) pieChart.destroy()
})

watch(() => party.value?.status, (newStatus) => {
  if (newStatus === 'RESULTS') nextTick(() => renderPieChart())
})

// Re-render when selected round changes
watch(selectedRound, () => {
  if (party.value?.status === 'RESULTS') nextTick(() => renderPieChart())
})

// Also re-render if results update while already in RESULTS phase
watch(() => party.value?.roundResults, () => {
  if (party.value?.status === 'RESULTS') nextTick(() => renderPieChart())
}, {deep: true})
</script>
