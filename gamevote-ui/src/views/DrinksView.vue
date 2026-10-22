<template>
  <main class="max-w-6xl mx-auto px-4 sm:px-6 py-8">
    <!-- Toast Notification -->
    <Transition name="toast">
      <div
        v-if="toastMsg"
        class="fixed bottom-6 right-6 z-50 flex items-center gap-3 px-5 py-3 rounded-2xl shadow-2xl text-white font-medium text-sm backdrop-blur-md"
        style="background: linear-gradient(135deg, rgba(99,102,241,0.95), rgba(139,92,246,0.95)); border: 1px solid rgba(255,255,255,0.2);"
      >
        <span>{{ toastMsg }}</span>
      </div>
    </Transition>

    <!-- Header & User Info -->
    <header class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-8">
      <div>
        <div class="flex items-center gap-3 flex-wrap">
          <h1 class="text-3xl font-extrabold tracking-tight">🍺 Drinks Tracker</h1>
          <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            Live Sync
          </span>
          <span class="px-2.5 py-1 rounded-full text-xs font-mono font-medium" style="background:var(--c-surface);border:1px solid var(--c-border);color:var(--c-muted)">
            Party {{ code }}
          </span>
        </div>
        <p class="text-sm mt-1" style="color:var(--c-muted)">
          Real-time group drinks & beer-equivalent counter
        </p>
      </div>

      <div class="flex items-center gap-3 w-full md:w-auto justify-between md:justify-end flex-wrap">
        <!-- Current Logged In User Pill -->
        <div
          v-if="currentUser"
          class="flex items-center gap-2.5 px-3.5 py-2 rounded-xl"
          style="background:var(--c-surface);border:1px solid var(--c-border)"
        >
          <div class="w-7 h-7 rounded-lg flex items-center justify-center font-bold text-xs" style="background:linear-gradient(135deg, var(--c-accent), var(--c-accent2));color:#fff">
            {{ currentUser.charAt(0).toUpperCase() }}
          </div>
          <div class="text-xs">
            <div class="text-xs text-muted leading-none">Logged in as</div>
            <div class="font-bold text-sm leading-tight text-white">{{ currentUser }}</div>
          </div>
          <span class="ml-1 text-xs px-2 py-0.5 rounded-full font-semibold bg-indigo-500/15 text-indigo-300 border border-indigo-500/30">
            {{ myBeersTotal.toFixed(1) }} 🍺
          </span>
        </div>

        <RouterLink :to="`/parties/${code}`" class="btn btn-ghost text-sm py-2">
          ← Back to Party
        </RouterLink>
      </div>
    </header>

    <!-- Stats Ribbon -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="stat-card">
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span>Party Beer Total</span>
          <span class="text-lg">🍻</span>
        </div>
        <div class="text-2xl sm:text-3xl font-extrabold text-white">
          {{ totalBeerEquivalents.toFixed(1) }}
          <span class="text-xs sm:text-sm font-normal text-muted">beers</span>
        </div>
        <div class="text-xs text-muted mt-1">Sum of all alcohol units</div>
      </div>

      <div class="stat-card">
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span>Total Drinks</span>
          <span class="text-lg">🥂</span>
        </div>
        <div class="text-2xl sm:text-3xl font-extrabold text-white">
          {{ rawBeerLog.length }}
          <span class="text-xs sm:text-sm font-normal text-muted">drinks</span>
        </div>
        <div class="text-xs text-muted mt-1">Logged across all attendees</div>
      </div>

      <div class="stat-card">
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span>Party MVP</span>
          <span class="text-lg">👑</span>
        </div>
        <div class="text-xl sm:text-2xl font-extrabold truncate text-amber-300">
          {{ topDrinker ? topDrinker.name : '—' }}
        </div>
        <div class="text-xs text-muted mt-1 truncate">
          {{ topDrinker ? `${topDrinker.beers.toFixed(1)} beers logged` : 'No drinks yet' }}
        </div>
      </div>

      <div class="stat-card">
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span>Your Drinks</span>
          <span class="text-lg">🎯</span>
        </div>
        <div class="text-2xl sm:text-3xl font-extrabold text-indigo-400">
          {{ myDrinksCount }}
          <span class="text-xs sm:text-sm font-normal text-muted">({{ myBeersTotal.toFixed(1) }} 🍺)</span>
        </div>
        <div class="text-xs text-muted mt-1">
          {{ myBeerSharePercent }}% of party total
        </div>
      </div>
    </div>

    <!-- Quick Drinks Shelf -->
    <section class="mb-10">
      <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div>
          <h2 class="text-lg font-bold flex items-center gap-2">
            <span>⚡ Quick Log Drink</span>
            <span class="text-xs font-normal text-muted">Tap to log 1 drink for <strong>{{ currentUser || 'yourself' }}</strong></span>
          </h2>
        </div>
        <button
          type="button"
          class="btn btn-ghost btn-sm text-xs font-semibold flex items-center gap-1.5"
          :class="{ 'text-indigo-400 border-indigo-500/40': showCustomBuilder }"
          @click="showCustomBuilder = !showCustomBuilder"
        >
          <span>{{ showCustomBuilder ? '✕ Close Builder' : '✨ + Add Custom Drink' }}</span>
        </button>
      </div>

      <!-- Presets Grid -->
      <div v-if="loadingPresets" class="flex justify-center py-10">
        <div class="spinner" style="width:32px;height:32px;border-width:3px"></div>
      </div>

      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
        <button
          v-for="dt in presets"
          :key="dt.id"
          type="button"
          class="drink-btn group"
          :disabled="loggingDrink"
          @click="logDrink(dt)"
        >
          <div class="text-3xl sm:text-4xl mb-1.5 transition-transform duration-200 group-hover:scale-110 group-active:scale-95">
            {{ drinkEmoji(dt.name || '') }}
          </div>
          <div class="font-bold text-sm text-white truncate w-full text-center">
            {{ dt.name }}
          </div>
          <div class="text-xs text-muted truncate w-full text-center mt-0.5">
            {{ dt.unitName || `${dt.volumeMl}ml` }}
          </div>
          <div class="mt-2 text-xs font-semibold px-2 py-0.5 rounded-full bg-white/5 border border-white/10 group-hover:bg-indigo-500/20 group-hover:border-indigo-500/40 group-hover:text-indigo-300 transition-colors">
            ≈ {{ dt.beerEquivalent }} 🍺
          </div>
        </button>

        <!-- Custom drink trigger tile -->
        <button
          type="button"
          class="drink-btn border-dashed border-white/20 hover:border-indigo-500/60 bg-transparent hover:bg-indigo-500/5 group text-center"
          @click="showCustomBuilder = true"
        >
          <div class="text-3xl mb-1.5 text-indigo-400 group-hover:scale-110 transition-transform">
            🪄
          </div>
          <div class="font-bold text-sm text-indigo-300">
            Custom
          </div>
          <div class="text-xs text-muted mt-0.5">
            New drink type
          </div>
          <div class="mt-2 text-xs font-semibold text-indigo-400">
            + Create
          </div>
        </button>
      </div>
    </section>

    <!-- Intuitive Custom Drink Builder (Expandable / Inline) -->
    <Transition name="expand">
      <section v-if="showCustomBuilder" class="card mb-10 overflow-hidden border-indigo-500/40 shadow-xl" style="background: linear-gradient(180deg, var(--c-surface), #121622);">
        <div class="flex items-center justify-between mb-5 border-b border-white/5 pb-4">
          <div>
            <h3 class="text-lg font-bold text-white flex items-center gap-2">
              <span>🪄 Add Custom Drink</span>
              <span class="text-xs px-2 py-0.5 rounded-full bg-indigo-500/20 text-indigo-300 border border-indigo-500/30">
                Calculator & Presets
              </span>
            </h3>
            <p class="text-xs text-muted mt-0.5">
              Pick a quick drink category or enter custom volume and ABV. We calculate the beer-equivalent automatically!
            </p>
          </div>
          <button
            type="button"
            class="text-muted hover:text-white text-lg p-1"
            @click="showCustomBuilder = false"
          >
            ✕
          </button>
        </div>

        <!-- 1. Quick Category Chips -->
        <div class="mb-5">
          <label class="block text-xs font-semibold uppercase tracking-wider text-muted mb-2">
            1. Quick Category Preset
          </label>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="cat in categoryPresets"
              :key="cat.label"
              type="button"
              class="px-3 py-1.5 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 border"
              :class="activeCategory === cat.label
                ? 'bg-indigo-500 text-white border-indigo-400 shadow-md'
                : 'bg-white/5 text-slate-300 border-white/10 hover:border-white/20 hover:bg-white/10'"
              @click="applyCategory(cat)"
            >
              <span>{{ cat.emoji }}</span>
              <span>{{ cat.label }}</span>
              <span class="opacity-60 text-[10px]">({{ cat.volumeMl }}ml, {{ cat.alcoholPercent }}%)</span>
            </button>
          </div>
        </div>

        <!-- 2. Inputs & Live Preview -->
        <div class="grid md:grid-cols-3 gap-6 items-start">
          <div class="md:col-span-2 flex flex-col gap-4">
            <div>
              <label class="block text-xs font-semibold text-muted mb-1.5">
                Drink Name
              </label>
              <div class="relative">
                <input
                  v-model="custom.name"
                  type="text"
                  class="input pl-10"
                  placeholder="e.g. Gin Tonic, Tequila Shot, IPA..."
                />
                <span class="absolute left-3 top-2.5 text-lg select-none">
                  {{ drinkEmoji(custom.name || activeCategoryEmoji) }}
                </span>
              </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <label class="block text-xs font-semibold text-muted mb-1.5">
                  Volume (ml)
                </label>
                <input
                  v-model.number="custom.volumeMl"
                  type="number"
                  class="input"
                  min="1"
                  max="5000"
                  placeholder="500"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-muted mb-1.5">
                  Alcohol (% vol)
                </label>
                <input
                  v-model.number="custom.alcoholPercent"
                  type="number"
                  step="0.1"
                  min="0"
                  max="100"
                  class="input"
                  placeholder="5.0"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-muted mb-1.5">
                  Unit Name
                </label>
                <input
                  v-model="custom.unitName"
                  type="text"
                  class="input"
                  placeholder="e.g. Shot, Glass, Pint"
                />
              </div>
            </div>
          </div>

          <!-- Live Calculation Result Card -->
          <div
            class="p-4 rounded-2xl border flex flex-col justify-between"
            style="background:var(--c-bg);border-color:var(--c-border)"
          >
            <div>
              <div class="text-xs text-muted font-semibold uppercase tracking-wider mb-1">
                Beer Equivalent Meter
              </div>
              <div class="flex items-baseline gap-2 mb-2">
                <span class="text-3xl font-extrabold text-amber-400">
                  {{ calculatedBeerEq }}
                </span>
                <span class="text-sm font-semibold text-white">🍺 beers</span>
              </div>

              <!-- Pure Alcohol calculation breakdown -->
              <div class="text-xs text-muted space-y-1 bg-white/5 p-2.5 rounded-xl border border-white/5">
                <div class="flex justify-between">
                  <span>Pure Alcohol:</span>
                  <span class="font-mono text-white">{{ calculatedPureAlcohol }} ml</span>
                </div>
                <div class="flex justify-between">
                  <span>Standard 0.5L Beer:</span>
                  <span class="font-mono text-white">25 ml alcohol</span>
                </div>
                <div class="border-t border-white/10 pt-1 flex justify-between text-[11px] text-indigo-300">
                  <span>Drink Strength:</span>
                  <span class="font-bold">{{ drinkStrengthLabel }}</span>
                </div>
              </div>
            </div>

            <div class="mt-4 flex gap-2">
              <button
                type="button"
                class="btn btn-primary flex-1 text-sm justify-center py-2.5"
                :disabled="!custom.name.trim() || savingCustom"
                @click="addCustomDrink"
              >
                <span v-if="savingCustom">Saving...</span>
                <span v-else>Save & Drink Now 🍻</span>
              </button>
              <button
                type="button"
                class="btn btn-ghost text-sm py-2.5 px-3"
                @click="showCustomBuilder = false"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </section>
    </Transition>

    <!-- Analytics & Charts Grid -->
    <section class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-12">
      <!-- Total Beers per Person (Leaderboard Bar Chart) -->
      <div class="card flex flex-col">
        <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
          <div>
            <h2 class="text-base font-bold text-white flex items-center gap-2">
              <span>🏆 Total Beer-Equivalents</span>
            </h2>
            <p class="text-xs text-muted">Leaderboard of all attendees</p>
          </div>
          <span class="text-xs px-2.5 py-1 rounded-full font-semibold bg-white/5 text-muted border border-white/10">
            {{ attendeeLeaderboard.length }} participants
          </span>
        </div>

        <div v-if="!rawBeerLog.length" class="flex-1 flex flex-col items-center justify-center py-12 text-center text-muted">
          <div class="text-4xl mb-2">🍺</div>
          <p class="text-sm font-semibold text-white">No drinks recorded yet</p>
          <p class="text-xs text-muted mt-1">Tap any drink above to record the first one!</p>
        </div>

        <div v-else class="flex-1 flex flex-col justify-center min-h-[240px]">
          <canvas ref="barChartEl" style="max-height:240px;width:100%"></canvas>
        </div>
      </div>

      <!-- Consumption Over Time (Line Chart with Real Timestamps) -->
      <div class="card flex flex-col">
        <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
          <div>
            <h2 class="text-base font-bold text-white flex items-center gap-2">
              <span>📈 Consumption Over Time</span>
            </h2>
            <p class="text-xs text-muted">Pacing & drinks timeline with real timestamps</p>
          </div>
          <span class="text-xs px-2.5 py-1 rounded-full font-semibold bg-white/5 text-muted border border-white/10">
            Cumulative Beers
          </span>
        </div>

        <div v-if="!rawBeerLog.length" class="flex-1 flex flex-col items-center justify-center py-12 text-center text-muted">
          <div class="text-4xl mb-2">⏱️</div>
          <p class="text-sm font-semibold text-white">Timeline is empty</p>
          <p class="text-xs text-muted mt-1">Timestamps and pacing will appear once drinks are logged.</p>
        </div>

        <div v-else class="flex-1 flex flex-col justify-center min-h-[240px]">
          <canvas ref="lineChartEl" style="max-height:240px;width:100%"></canvas>
        </div>
      </div>
    </section>

    <!-- Drink Log (Placed at the bottom as requested) -->
    <section class="card mb-8">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-6 border-b border-white/5 pb-4">
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-lg font-bold text-white">📋 Live Activity Log</h2>
            <span class="px-2 py-0.5 rounded-full text-xs font-semibold bg-white/10 text-white">
              {{ filteredLog.length }}
            </span>
          </div>
          <p class="text-xs text-muted mt-0.5">
            Every drink logged in this party with persistent timestamps
          </p>
        </div>

        <!-- Filter tabs: All vs Mine -->
        <div class="flex items-center gap-1.5 p-1 rounded-xl bg-black/20 border border-white/5 text-xs">
          <button
            type="button"
            class="px-3 py-1.5 rounded-lg font-semibold transition-all"
            :class="logFilter === 'all' ? 'bg-indigo-600 text-white shadow' : 'text-muted hover:text-white'"
            @click="logFilter = 'all'"
          >
            All Attendees
          </button>
          <button
            v-if="currentUser"
            type="button"
            class="px-3 py-1.5 rounded-lg font-semibold transition-all"
            :class="logFilter === 'me' ? 'bg-indigo-600 text-white shadow' : 'text-muted hover:text-white'"
            @click="logFilter = 'me'"
          >
            My Drinks Only
          </button>
        </div>
      </div>

      <div v-if="!filteredLog.length" class="text-center py-12 text-muted">
        <div class="text-4xl mb-2">🧊</div>
        <p class="text-sm font-semibold text-white">No drinks found</p>
        <p class="text-xs text-muted mt-1">
          {{ logFilter === 'me' ? "You haven't logged any drinks yet." : 'No drinks recorded in this party yet.' }}
        </p>
      </div>

      <div v-else class="flex flex-col gap-2.5 max-h-[380px] overflow-y-auto pr-1">
        <TransitionGroup name="list" tag="div" class="flex flex-col gap-2">
          <div
            v-for="entry in filteredLog"
            :key="entry.id"
            class="flex items-center justify-between px-4 py-3 rounded-xl transition-all hover:translate-x-0.5"
            :class="entry.attendee === currentUser ? 'border-indigo-500/30' : 'border-white/5'"
            style="background:var(--c-bg);border-width:1px;"
          >
            <!-- Left: Emoji + Attendee + Drink Name -->
            <div class="flex items-center gap-3 min-w-0">
              <span class="text-2xl flex-shrink-0 select-none">
                {{ drinkEmoji(entry.drinkName) }}
              </span>
              <div class="min-w-0">
                <div class="flex items-center gap-2 flex-wrap">
                  <span class="font-bold text-sm text-white truncate">
                    {{ entry.attendee }}
                  </span>
                  <span
                    v-if="entry.attendee === currentUser"
                    class="text-[10px] px-1.5 py-0.2 rounded font-semibold bg-indigo-500/20 text-indigo-300 border border-indigo-500/30"
                  >
                    You
                  </span>
                  <span class="text-xs text-muted truncate">
                    had a <strong class="text-slate-200">{{ entry.drinkName }}</strong>
                  </span>
                </div>
                <div class="text-xs text-muted flex items-center gap-2 mt-0.5">
                  <span>{{ formatFullTime(entry.createdAt) }}</span>
                  <span>•</span>
                  <span>{{ formatRelativeTime(entry.createdAt) }}</span>
                </div>
              </div>
            </div>

            <!-- Right: Beer Equivalent Badge -->
            <div class="flex-shrink-0 flex items-center gap-2">
              <span class="text-xs font-bold px-2.5 py-1 rounded-full bg-emerald-500/15 text-emerald-400 border border-emerald-500/20">
                +{{ Number(entry.beerEquivalent).toFixed(1) }} 🍺
              </span>
            </div>
          </div>
        </TransitionGroup>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getParty, listDrinks, createDrink, addBeer, client } from '@/api-client'
import type { DrinkType, Party, BeerEntry } from '@/generated-api'
import {
  Chart,
  BarElement,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Tooltip,
  Legend,
  BarController,
  LineController
} from 'chart.js'

Chart.register(
  BarElement,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Tooltip,
  Legend,
  BarController,
  LineController
)

const route = useRoute()
const authStore = useAuthStore()
const code = computed(() => route.params.code as string)

// Currently logged-in user: FIXED, no selection in frontend
const currentUser = computed(() => authStore.username || '')

const party = ref<Party | null>(null)
const presets = ref<DrinkType[]>([])
const loadingPresets = ref(true)
const loggingDrink = ref(false)
const savingCustom = ref(false)
const showCustomBuilder = ref(false)
const toastMsg = ref<string | null>(null)
let toastTimer: any = null

function showToast(msg: string) {
  toastMsg.value = msg
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastMsg.value = null
  }, 2500)
}

// ── Persistent Beer Log from Party API ──────────────────────────────────────
const rawBeerLog = computed<BeerEntry[]>(() => {
  return party.value?.beerLog ?? []
})

// Log filter
const logFilter = ref<'all' | 'me'>('all')

const filteredLog = computed(() => {
  const list = [...rawBeerLog.value].reverse() // latest first
  if (logFilter.value === 'me' && currentUser.value) {
    return list.filter(e => e.attendee === currentUser.value)
  }
  return list
})

// ── Metrics ─────────────────────────────────────────────────────────────────
const totalBeerEquivalents = computed(() => {
  return rawBeerLog.value.reduce((sum, e) => sum + (e.beerEquivalent || 0), 0)
})

const myDrinksCount = computed(() => {
  if (!currentUser.value) return 0
  return rawBeerLog.value.filter(e => e.attendee === currentUser.value).length
})

const myBeersTotal = computed(() => {
  if (!currentUser.value) return 0
  return rawBeerLog.value
    .filter(e => e.attendee === currentUser.value)
    .reduce((sum, e) => sum + (e.beerEquivalent || 0), 0)
})

const myBeerSharePercent = computed(() => {
  if (totalBeerEquivalents.value === 0) return 0
  return Math.round((myBeersTotal.value / totalBeerEquivalents.value) * 100)
})

const attendeeLeaderboard = computed(() => {
  const allAttendees = party.value?.attendees?.length
    ? party.value.attendees
    : [...new Set(rawBeerLog.value.map(e => e.attendee))]

  const map: Record<string, number> = {}
  for (const a of allAttendees) map[a] = 0
  for (const e of rawBeerLog.value) {
    map[e.attendee] = (map[e.attendee] ?? 0) + (e.beerEquivalent || 0)
  }

  return Object.entries(map)
    .map(([name, beers]) => ({ name, beers }))
    .sort((a, b) => b.beers - a.beers)
})

const topDrinker = computed(() => {
  const first = attendeeLeaderboard.value[0]
  if (!first || first.beers === 0) {
    return null
  }
  return first
})

// ── Intuitive Custom Drink Builder ──────────────────────────────────────────
interface CategoryPreset {
  label: string
  emoji: string
  volumeMl: number
  alcoholPercent: number
  unitName: string
}

const categoryPresets: CategoryPreset[] = [
  { label: 'Shot', emoji: '🥃', volumeMl: 40, alcoholPercent: 40.0, unitName: 'Shot' },
  { label: 'Pint', emoji: '🍺', volumeMl: 500, alcoholPercent: 5.0, unitName: 'Pint' },
  { label: 'Bottle / Can', emoji: '🥫', volumeMl: 330, alcoholPercent: 5.0, unitName: 'Can' },
  { label: 'Wine Glass', emoji: '🍷', volumeMl: 150, alcoholPercent: 12.0, unitName: 'Glass' },
  { label: 'Cocktail', emoji: '🍹', volumeMl: 250, alcoholPercent: 15.0, unitName: 'Cocktail' },
  { label: 'Spritzer', emoji: '🍾', volumeMl: 250, alcoholPercent: 6.5, unitName: 'Spritzer' },
]

const activeCategory = ref<string>('Shot')
const activeCategoryEmoji = ref<string>('🥃')

const custom = ref({
  name: '',
  volumeMl: 40,
  alcoholPercent: 40.0,
  unitName: 'Shot'
})

function applyCategory(cat: CategoryPreset) {
  activeCategory.value = cat.label
  activeCategoryEmoji.value = cat.emoji
  custom.value.volumeMl = cat.volumeMl
  custom.value.alcoholPercent = cat.alcoholPercent
  custom.value.unitName = cat.unitName
  if (!custom.value.name.trim()) {
    custom.value.name = cat.label
  }
}

const calculatedPureAlcohol = computed(() => {
  const ml = custom.value.volumeMl || 0
  const abv = custom.value.alcoholPercent || 0
  return parseFloat(((ml * abv) / 100).toFixed(1))
})

const calculatedBeerEq = computed(() => {
  return parseFloat((calculatedPureAlcohol.value / 25).toFixed(2))
})

const drinkStrengthLabel = computed(() => {
  const eq = calculatedBeerEq.value
  if (eq <= 0.4) return '🍃 Light & Sessionable'
  if (eq <= 1.2) return '🍺 Standard Drink'
  if (eq <= 2.2) return '⚡ Strong Drink'
  return '🔥 Heavy Hitter!'
})

// ── Emoji Mapping ───────────────────────────────────────────────────────────
const drinkEmojis: Record<string, string> = {
  beer: '🍺',
  bier: '🍺',
  vodka: '🥃',
  wodka: '🥃',
  wein: '🍷',
  wine: '🍷',
  spritzer: '🍾',
  sekt: '🍾',
  champagne: '🍾',
  gin: '🟢',
  'gin tonic': '🟢',
  jägermeister: '🦌',
  jager: '🦌',
  rum: '🍹',
  cocktail: '🍹',
  'berliner luft': '💨',
  klopfer: '💥',
  shot: '🥃',
  tequila: '🌵',
  whiskey: '🥃',
  whisky: '🥃',
  cider: '🍏',
}

function drinkEmoji(name: string) {
  const lower = name.toLowerCase().trim()
  for (const [key, emoji] of Object.entries(drinkEmojis)) {
    if (lower.includes(key)) return emoji
  }
  return '🥂'
}

// ── Time Formatting ─────────────────────────────────────────────────────────
function formatFullTime(isoStr: string) {
  try {
    const d = new Date(isoStr)
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } catch {
    return isoStr
  }
}

function formatRelativeTime(isoStr: string) {
  try {
    const diffMs = Date.now() - new Date(isoStr).getTime()
    const diffSec = Math.floor(diffMs / 1000)
    if (diffSec < 60) return 'Just now'
    const diffMin = Math.floor(diffSec / 60)
    if (diffMin < 60) return `${diffMin}m ago`
    const diffHours = Math.floor(diffMin / 60)
    if (diffHours < 24) return `${diffHours}h ago`
    return `${Math.floor(diffHours / 24)}d ago`
  } catch {
    return ''
  }
}

// ── Actions ─────────────────────────────────────────────────────────────────
async function logDrink(dt: DrinkType | { name: string; beerEquivalent: number; unitName?: string }) {
  const attendee = currentUser.value
  if (!attendee) {
    alert('Please log in with a username first!')
    return
  }

  const drinkName = dt.name || 'Drink'
  const beerEq = Number(dt.beerEquivalent ?? 1.0)
  loggingDrink.value = true

  try {
    await addBeer({
      path: { code: code.value },
      body: {
        attendee,
        drinkName,
        beerEquivalent: beerEq,
      }
    })

    showToast(`Logged ${drinkEmoji(drinkName)} ${drinkName} for ${attendee}! 🍻`)

    // Optimistically re-fetch party to guarantee instantaneous UI update
    const pResp = await getParty({ path: { code: code.value } })
    if (pResp?.data) {
      party.value = pResp.data
    }
  } catch (err) {
    console.error('Error logging drink:', err)
  } finally {
    loggingDrink.value = false
  }
}

async function addCustomDrink() {
  const name = custom.value.name.trim()
  if (!name) return

  savingCustom.value = true
  const beerEq = calculatedBeerEq.value

  try {
    // Check if preset already exists
    const existing = presets.value.find(p => p.name?.toLowerCase() === name.toLowerCase())
    if (existing) {
      await logDrink(existing)
      showCustomBuilder.value = false
      return
    }

    const savedResp = await createDrink({
      body: {
        name,
        volumeMl: custom.value.volumeMl,
        alcoholPercent: custom.value.alcoholPercent,
        beerEquivalent: beerEq,
        unitName: custom.value.unitName || 'Custom',
      },
    })

    if (!savedResp.data) {
      // If conflict or error, try to log with custom values directly
      await logDrink({ name, beerEquivalent: beerEq, unitName: custom.value.unitName })
    } else {
      presets.value.push(savedResp.data)
      await logDrink(savedResp.data)
    }

    showCustomBuilder.value = false
    custom.value.name = ''
  } catch (err: any) {
    console.error('Failed to create custom drink:', err)
    // Fallback: log drink anyway
    await logDrink({ name, beerEquivalent: beerEq, unitName: custom.value.unitName })
    showCustomBuilder.value = false
  } finally {
    savingCustom.value = false
  }
}

// ── Chart Rendering ─────────────────────────────────────────────────────────
const barChartEl = ref<HTMLCanvasElement | null>(null)
const lineChartEl = ref<HTMLCanvasElement | null>(null)
let barChart: Chart | null = null
let lineChart: Chart | null = null

const palette = [
  '#6366f1', '#8b5cf6', '#22c55e', '#f59e0b', '#ec4899', '#06b6d4', '#10b981', '#3b82f6', '#f43f5e'
]

function renderCharts() {
  const logEntries = rawBeerLog.value
  const leaderboard = attendeeLeaderboard.value

  // 1. Bar Chart (Leaderboard)
  if (barChartEl.value) {
    if (barChart) barChart.destroy()

    const labels = leaderboard.map(l => l.name)
    const data = leaderboard.map(l => l.beers)
    const backgroundColors = labels.map((name, i) => {
      return name === currentUser.value ? '#6366f1' : palette[i % palette.length]
    })

    barChart = new Chart(barChartEl.value, {
      type: 'bar',
      data: {
        labels,
        datasets: [{
          label: 'Total Beer-Equivalents',
          data,
          backgroundColor: backgroundColors,
          borderRadius: 8,
          borderWidth: 0,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false },
          tooltip: {
            callbacks: {
              label: (ctx) => ` ${Number(ctx.raw).toFixed(2)} beers total`
            }
          }
        },
        scales: {
          x: {
            ticks: { color: '#94a3b8', font: { size: 12, weight: 'bold' } },
            grid: { display: false }
          },
          y: {
            beginAtZero: true,
            ticks: { color: '#94a3b8' },
            grid: { color: '#1e2739' }
          }
        }
      }
    })
  }

  // 2. Line Chart (Consumption Over Time with REAL Timestamps)
  if (lineChartEl.value) {
    if (lineChart) lineChart.destroy()

    if (!logEntries.length) return

    // Sort chronologically
    const sorted = [...logEntries].sort(
      (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
    )

    // Build timeline points
    const allAttendees = party.value?.attendees?.length
      ? party.value.attendees
      : [...new Set(logEntries.map(e => e.attendee))]

    const runningTotals: Record<string, number> = {}
    const attendeeSeries: Record<string, number[]> = {}
    for (const a of allAttendees) {
      runningTotals[a] = 0
      attendeeSeries[a] = []
    }

    const timeLabels: string[] = []

    for (const entry of sorted) {
      const timeStr = formatFullTime(entry.createdAt)
      timeLabels.push(timeStr)

      runningTotals[entry.attendee] = (runningTotals[entry.attendee] ?? 0) + (entry.beerEquivalent || 0)
      for (const a of allAttendees) {
        const list = attendeeSeries[a] || []
        list.push(parseFloat((runningTotals[a] ?? 0).toFixed(2)))
        attendeeSeries[a] = list
      }
    }

    const datasets = allAttendees.map((attendee, idx) => {
      const color = attendee === currentUser.value ? '#6366f1' : palette[idx % palette.length]
      return {
        label: attendee,
        data: attendeeSeries[attendee] || [],
        borderColor: color,
        backgroundColor: color + '22',
        borderWidth: attendee === currentUser.value ? 3 : 2,
        fill: false,
        tension: 0.35,
        pointRadius: 4,
        pointHoverRadius: 6,
      }
    })

    lineChart = new Chart(lineChartEl.value, {
      type: 'line',
      data: {
        labels: timeLabels,
        datasets,
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          mode: 'index',
          intersect: false,
        },
        plugins: {
          legend: {
            position: 'top',
            labels: {
              color: '#e2e8f0',
              boxWidth: 12,
              font: { size: 11, weight: 'bold' }
            }
          },
          tooltip: {
            callbacks: {
              title: (items) => `Time: ${items[0]?.label ?? ''}`,
              label: (ctx) => ` ${ctx.dataset.label}: ${ctx.raw} beers`
            }
          }
        },
        scales: {
          x: {
            ticks: {
              color: '#94a3b8',
              maxRotation: 45,
              minRotation: 0,
              font: { size: 10 }
            },
            grid: { color: '#1e2739' }
          },
          y: {
            beginAtZero: true,
            ticks: { color: '#94a3b8' },
            grid: { color: '#1e2739' }
          }
        }
      }
    })
  }
}

// ── SSE Live Updates ────────────────────────────────────────────────────────
let eventSource: EventSource | null = null

function setupSSE() {
  const streamUrl = client.buildUrl({
    url: '/parties/{code}/stream',
    path: { code: code.value },
    query: { username: currentUser.value || 'guest' }
  })

  eventSource = new EventSource(streamUrl)
  eventSource.addEventListener('party_updated', (e: MessageEvent) => {
    try {
      const updatedParty = JSON.parse(e.data) as Party
      party.value = updatedParty
    } catch (err) {
      console.error('Failed to parse SSE party_updated:', err)
    }
  })

  eventSource.onerror = () => {
    // quietly reconnect
  }
}

// Watch party changes to update charts
watch(() => party.value?.beerLog, () => {
  nextTick(() => renderCharts())
}, { deep: true })

onMounted(async () => {
  // 1. Setup SSE for live sync across all party attendees
  setupSSE()

  // 2. Fetch party & drink presets
  const [pResp, dtsResp] = await Promise.all([
    getParty({ path: { code: code.value } }).catch(() => null),
    listDrinks().catch(() => ({ data: [] }))
  ])

  if (pResp?.data) {
    party.value = pResp.data
  }
  presets.value = dtsResp?.data ?? []
  loadingPresets.value = false

  await nextTick()
  renderCharts()
})

onUnmounted(() => {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  if (barChart) barChart.destroy()
  if (lineChart) lineChart.destroy()
  if (toastTimer) clearTimeout(toastTimer)
})
</script>

<style scoped>
.stat-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: 1rem;
  padding: 1.25rem;
  transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
}
.stat-card:hover {
  transform: translateY(-2px);
  border-color: rgba(99,102,241,0.4);
  box-shadow: 0 4px 20px rgba(99,102,241,0.1);
}

.drink-btn {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: 1rem;
  padding: 1.25rem 0.75rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;
}
.drink-btn:hover:not(:disabled) {
  border-color: var(--c-accent);
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(99,102,241,0.22);
}
.drink-btn:active:not(:disabled) {
  transform: scale(0.96);
}
.drink-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Transitions */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 800px;
}
.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
  transform: scaleY(0.95);
}

.list-enter-active,
.list-leave-active {
  transition: all 0.25s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-15px);
}
</style>
