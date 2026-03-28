<template>
  <main class="max-w-5xl mx-auto px-6 py-8">
    <div class="flex items-center justify-between mb-8 flex-wrap gap-3">
      <div>
        <h1 class="text-2xl font-extrabold">🍺 Drinks Tracker</h1>
        <p class="text-sm mt-1" style="color:var(--c-muted)">Party <span class="font-mono text-indigo-400">{{ code }}</span></p>
      </div>
      <RouterLink :to="`/parties/${code}`" class="btn btn-ghost">← Back to Party</RouterLink>
    </div>

    <div class="grid lg:grid-cols-3 gap-6">
      <!-- Left: Drink buttons + form -->
      <div class="lg:col-span-1 flex flex-col gap-4">
        <!-- Who am I drinking for? -->
        <div class="card">
          <label class="block text-sm font-semibold mb-2">Drinking for</label>
          <select v-model="selectedAttendee" class="input">
            <option v-for="a in attendees" :key="a" :value="a">{{ a }}</option>
          </select>
        </div>

        <!-- Preset buttons -->
        <div class="card">
          <p class="text-sm font-semibold mb-3">Common Drinks</p>
          <div v-if="loadingPresets" class="flex justify-center py-4"><div class="spinner" style="width:24px;height:24px;border-width:2px"></div></div>
          <div v-else class="flex flex-wrap gap-2">
            <button
              v-for="dt in presets"
              :key="dt.id"
              class="btn btn-ghost btn-sm"
              :title="`≈ ${dt.beerEquivalent} beer`"
              @click="logDrink(dt)"
            >
              {{ drinkEmoji(dt.name || '') }} {{ dt.name }}<br />
              <span class="text-xs font-normal opacity-60">{{ dt.unitName }}</span>
            </button>
          </div>
        </div>

        <!-- Custom drink form -->
        <details class="card">
          <summary class="cursor-pointer text-sm font-semibold text-indigo-400">+ Add Custom Drink</summary>
          <div class="mt-4 flex flex-col gap-3">
            <input v-model="custom.name" class="input" placeholder="Drink name" />
            <div class="flex gap-2">
              <input v-model.number="custom.volumeMl" class="input" type="number" placeholder="Volume (ml)" min="1" />
              <input v-model.number="custom.alcoholPercent" class="input" type="number" placeholder="Alcohol %" min="0" max="100" step="0.1" />
            </div>
            <input v-model="custom.unitName" class="input" placeholder="Unit name (e.g. Shot)" />
            <button class="btn btn-primary" :disabled="!custom.name" @click="addCustomDrink">Save & Log Drink</button>
          </div>
        </details>
      </div>

      <!-- Right: Log + Charts -->
      <div class="lg:col-span-2 flex flex-col gap-6">
        <!-- Drink Log -->
        <div class="card">
          <h2 class="text-lg font-bold mb-4">📋 Log</h2>
          <div v-if="!log.length" class="text-center py-6" style="color:var(--c-muted)">No drinks logged yet</div>
          <div class="flex flex-col gap-2 max-h-64 overflow-y-auto">
            <TransitionGroup name="list" tag="div" class="flex flex-col gap-2">
              <div
                v-for="(entry, i) in logReversed"
                :key="i"
                class="flex items-center justify-between px-3 py-2 rounded-xl"
                style="background:var(--c-bg);border:1px solid var(--c-border)"
              >
                <span class="text-sm">{{ drinkEmoji(entry.name) }} <strong>{{ entry.attendee }}</strong> – {{ entry.name }}</span>
                <span class="text-xs" style="color:var(--c-muted)">{{ entry.time }}</span>
              </div>
            </TransitionGroup>
          </div>
        </div>

        <!-- Bar chart: total beers per person -->
        <div class="card">
          <h2 class="text-base font-bold mb-4">🏆 Total Beer-Equivalents</h2>
          <canvas ref="barChartEl" style="max-height:220px"></canvas>
        </div>

        <!-- Line chart: consumption over time -->
        <div class="card">
          <h2 class="text-base font-bold mb-4">📈 Consumption Over Time</h2>
          <canvas ref="lineChartEl" style="max-height:220px"></canvas>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getParty, getDrinkTypes, postDrinkType, postBeer } from '@/api-client'
import type { ModelsDrinkType as DrinkType, ServicePartyDto as PartyDTO } from '@/generated-api'
import { Chart, BarElement, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, BarController, LineController } from 'chart.js'
Chart.register(BarElement, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend, BarController, LineController)

const route = useRoute()
const authStore = useAuthStore()
const code = computed(() => route.params.code as string)

const party = ref<PartyDTO | null>(null)
const attendees = computed(() => party.value?.attendees ?? [])
const selectedAttendee = ref(authStore.username ?? '')
const presets = ref<DrinkType[]>([])
const loadingPresets = ref(true)

interface LogEntry { attendee: string; name: string; beers: number; time: string; ts: number }
const log = ref<LogEntry[]>([])
const logReversed = computed(() => [...log.value].reverse())

const custom = ref({ name: '', volumeMl: 500, alcoholPercent: 5, unitName: 'custom' })

const barChartEl = ref<HTMLCanvasElement | null>(null)
const lineChartEl = ref<HTMLCanvasElement | null>(null)
let barChart: Chart<'bar'> | null = null
let lineChart: Chart<'line'> | null = null

const drinkEmojis: Record<string, string> = {
  beer: '🍺', vodka: '🥃', wein: '🍷', spritzer: '🍾', gin: '🟢',
  jägermeister: '🦌', rum: '🍹', 'berliner luft': '💨', klopfer: '💥',
}
function drinkEmoji(name: string) {
  return drinkEmojis[name.toLowerCase()] ?? '🥂'
}

// ── Log a drink ────────────────────────────────────────────────────────
async function logDrink(dt: DrinkType) {
  const attendee = selectedAttendee.value
  if (!attendee) return
  await postBeer({ path: { code: code.value }, body: { attendee } }).catch(() => {})
  log.value.push({
    attendee,
    name: dt.name || 'Unknown',
    beers: dt.beerEquivalent || 0,
    time: new Date().toLocaleTimeString(),
    ts: Date.now(),
  })
  nextTick(() => renderCharts())
}

// ── Add custom drink type + log it ────────────────────────────────────
async function addCustomDrink() {
  if (!custom.value.name) return
  const pureAlcohol = custom.value.volumeMl * (custom.value.alcoholPercent / 100)
  const beerEq = parseFloat((pureAlcohol / 25).toFixed(2))
  const newDt: Omit<DrinkType, 'id'> = {
    name: custom.value.name,
    volumeMl: custom.value.volumeMl,
    alcoholPercent: custom.value.alcoholPercent,
    beerEquivalent: beerEq,
    unitName: custom.value.unitName,
  }
  const savedResp = await postDrinkType({ body: newDt as any }).catch(() => ({ data: { ...newDt, id: Date.now().toString() } }))
  const saved = savedResp.data as DrinkType
  presets.value.push(saved)
  await logDrink(saved)
  custom.value = { name: '', volumeMl: 500, alcoholPercent: 5, unitName: 'custom' }
}

// ── Charts ─────────────────────────────────────────────────────────────
function renderCharts() {
  const allAttendees = attendees.value.length ? attendees.value : [...new Set(log.value.map(e => e.attendee))]
  const totals: Record<string, number> = {}
  for (const a of allAttendees) totals[a] = 0
  for (const e of log.value) totals[e.attendee] = (totals[e.attendee] ?? 0) + e.beers

  const colors = ['#6366f1', '#8b5cf6', '#22c55e', '#f59e0b', '#ef4444', '#06b6d4', '#ec4899']

  // Bar chart
  if (barChartEl.value) {
    if (barChart) barChart.destroy()
    barChart = new Chart(barChartEl.value, {
      type: 'bar',
      data: {
        labels: Object.keys(totals),
        datasets: [{
          label: 'Beer equivalents',
          data: Object.values(totals),
          backgroundColor: colors,
          borderRadius: 8,
        }]
      },
      options: {
        plugins: { legend: { display: false } },
        scales: {
          x: { ticks: { color: '#94a3b8' }, grid: { color: '#1e2739' } },
          y: { ticks: { color: '#94a3b8' }, grid: { color: '#1e2739' }, beginAtZero: true },
        }
      }
    })
  }

  // Line chart (cumulative over time)
  if (lineChartEl.value) {
    if (lineChart) lineChart.destroy()

    // Build cumulative totals over time
    const sorted = [...log.value].sort((a, b) => a.ts - b.ts)
    const cumulative: Record<string, number[]> = {}
    const timeLabels: string[] = []
    const running: Record<string, number> = {}

    for (const a of allAttendees) { cumulative[a] = []; running[a] = 0 }

    for (const entry of sorted) {
      timeLabels.push(entry.time)
      running[entry.attendee] = (running[entry.attendee] ?? 0) + entry.beers
      for (const a of allAttendees) {
        ;(cumulative[a] ??= []).push(running[a] ?? 0)
      }
    }

    lineChart = new Chart(lineChartEl.value, {
      type: 'line',
      data: {
        labels: timeLabels,
        datasets: allAttendees.map((a, i) => ({
          label: a,
          data: cumulative[a] as number[],
          borderColor: colors[i % colors.length],
          backgroundColor: colors[i % colors.length] + '22',
          fill: true,
          tension: 0.4,
          pointRadius: 4,
        }))
      },
      options: {
        plugins: { legend: { labels: { color: '#e2e8f0' } } },
        scales: {
          x: { ticks: { color: '#94a3b8' }, grid: { color: '#1e2739' } },
          y: { ticks: { color: '#94a3b8' }, grid: { color: '#1e2739' }, beginAtZero: true },
        }
      }
    })
  }
}

onMounted(async () => {
  const [pResp, dtsResp] = await Promise.all([
    getParty({ path: { code: code.value } }).catch(() => null),
    getDrinkTypes().catch(() => ({ data: [] }))
  ])
  const p = pResp?.data
  if (p) {
    party.value = p as PartyDTO
    if (!selectedAttendee.value && p.attendees?.length) {
      selectedAttendee.value = p.attendees[0] as string
    }
    // Prefill log from beer counts (without time)
    for (const [attendee, count] of Object.entries(p.beerPerAttendee || {})) {
      for (let i = 0; i < count; i++) {
        log.value.push({ attendee, name: 'Beer', beers: 1, time: '?', ts: 0 })
      }
    }
  }
  presets.value = dtsResp?.data as DrinkType[] || []
  loadingPresets.value = false
  loadingPresets.value = false
  await nextTick()
  renderCharts()
})
</script>
