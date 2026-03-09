<template>
  <section>
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

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { usePartyStore } from '@/stores/party'
import { Chart, ArcElement, Tooltip, Legend, PieController } from 'chart.js'
Chart.register(ArcElement, Tooltip, Legend, PieController)

const partyStore = usePartyStore()
const party = computed(() => partyStore.party)

const pieChartEl = ref<HTMLCanvasElement | null>(null)
let pieChart: Chart | null = null

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

// Also re-render if results update while already in RESULTS phase
watch(() => party.value?.results, () => {
  if (party.value?.status === 'RESULTS') nextTick(() => renderPieChart())
}, { deep: true })
</script>
