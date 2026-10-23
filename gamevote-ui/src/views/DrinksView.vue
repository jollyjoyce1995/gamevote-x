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
          Real-time group drinks, blood alcohol concentration (BAC) & consumption timelines
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
    <div class="grid grid-cols-2 lg:grid-cols-5 gap-3.5 mb-8">
      <!-- 1. Party Beers -->
      <div class="stat-card">
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span>Party Total</span>
          <span class="text-lg">🍻</span>
        </div>
        <div class="text-2xl sm:text-3xl font-extrabold text-white">
          {{ totalBeerEquivalents.toFixed(1) }}
          <span class="text-xs sm:text-sm font-normal text-muted">beers</span>
        </div>
        <div class="text-xs text-muted mt-1">{{ rawBeerLog.length }} total drinks</div>
      </div>

      <!-- 2. Your Drinks -->
      <div class="stat-card">
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span>Your Drinks</span>
          <span class="text-lg">🎯</span>
        </div>
        <div class="text-2xl sm:text-3xl font-extrabold text-indigo-400">
          {{ myDrinksCount }}
          <span class="text-xs sm:text-sm font-normal text-muted">({{ myBeersTotal.toFixed(1) }} 🍺)</span>
        </div>
        <div class="text-xs text-muted mt-1">{{ myBeerSharePercent }}% of party total</div>
      </div>

      <!-- 3. Current Estimated Blood Alcohol (BAC) -->
      <div
        class="stat-card border-rose-500/25 cursor-pointer select-none group"
        title="Click to toggle between Promille (‰) and Percent (BAC)"
        @click="toggleBacUnit"
      >
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span class="flex items-center gap-1.5">
            <span>Your Blood Alcohol</span>
            <span class="text-[10px] font-mono px-1.5 py-0.2 rounded bg-white/10 text-indigo-300 group-hover:bg-indigo-500/30 transition-colors">
              {{ bacUnit === 'promille' ? '‰' : '%' }}
            </span>
          </span>
          <span class="text-lg transition-transform group-hover:scale-110">🩸</span>
        </div>
        <div class="flex items-baseline gap-1.5 flex-wrap">
          <span class="text-2xl sm:text-3xl font-extrabold" :class="bacColorClass(userBACSummary.currentBAC)">
            {{ bacUnit === 'promille' ? formatPromille(userBACSummary.currentBAC) : formatBAC(userBACSummary.currentBAC) }}
          </span>
          <span class="text-xs font-semibold text-muted">
            ({{ bacUnit === 'promille' ? formatBAC(userBACSummary.currentBAC) : formatPromille(userBACSummary.currentBAC) }})
          </span>
        </div>
        <div class="flex items-center gap-1.5 mt-1 text-xs">
          <span class="w-2 h-2 rounded-full" :class="bacDotClass(userBACSummary.currentBAC)"></span>
          <span class="font-medium text-slate-300 truncate">{{ getSobrietyStage(userBACSummary.currentBAC).label }}</span>
        </div>
      </div>

      <!-- 4. Sober Countdown -->
      <div class="stat-card">
        <div class="flex items-center justify-between text-muted text-xs font-semibold uppercase tracking-wider mb-1">
          <span>Sober Estimate</span>
          <span class="text-lg">⏱️</span>
        </div>
        <div class="text-xl sm:text-2xl font-extrabold text-white truncate">
          {{ formatSoberEta(userBACSummary.soberTime) }}
        </div>
        <div class="text-xs text-muted mt-1 truncate">
          {{ userBACSummary.currentBAC <= 0.0001 ? 'Currently at 0.00%' : `Sober at ~${formatTimeOnly(userBACSummary.soberTime)}` }}
        </div>
      </div>

      <!-- 5. Party MVP -->
      <div class="stat-card col-span-2 lg:col-span-1">
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
        <div class="flex items-center gap-2 flex-wrap">
          <!-- Hidable Beer Calculation Explanation Toggle -->
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-semibold transition-all border flex items-center gap-1.5 cursor-pointer"
            :class="showFormulaGuide 
              ? 'bg-amber-500/20 text-amber-200 border-amber-500/50 shadow-sm' 
              : 'bg-white/5 text-slate-300 border-white/10 hover:border-white/20 hover:text-white'"
            @click="showFormulaGuide = !showFormulaGuide"
          >
            <span>📐 How are beers calculated?</span>
            <span class="text-[10px]">{{ showFormulaGuide ? '▲ Hide' : '▼ Show' }}</span>
          </button>

          <!-- Optional Toggle for Party Custom Drinks -->
          <button
            v-if="partyCustomDrinks.length > 0"
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-semibold transition-all border flex items-center gap-1.5 cursor-pointer"
            :class="showCustomPartyDrinks 
              ? 'bg-purple-600/25 text-purple-200 border-purple-500/50 shadow-sm' 
              : 'bg-white/5 text-slate-300 border-white/10 hover:border-white/20 hover:text-white'"
            @click="showCustomPartyDrinks = !showCustomPartyDrinks"
          >
            <span>🪄 Party Custom Drinks ({{ partyCustomDrinks.length }})</span>
            <span class="text-[10px]">{{ showCustomPartyDrinks ? '▲ Hide' : '▼ Show' }}</span>
          </button>

          <!-- Add Custom Drink Button -->
          <button
            type="button"
            class="btn btn-ghost btn-sm text-xs font-semibold flex items-center gap-1.5 cursor-pointer"
            :class="{ 'text-indigo-400 border-indigo-500/40': showCustomBuilder }"
            @click="showCustomBuilder = !showCustomBuilder"
          >
            <span>{{ showCustomBuilder ? '✕ Close Builder' : '✨ + Add Custom Drink' }}</span>
          </button>
        </div>
      </div>

      <!-- Hidable Beer Calculation Guide & Live Converter -->
      <Transition name="guide-expand">
        <div
          v-if="showFormulaGuide"
          class="card mb-6 p-5 border-amber-500/30 shadow-xl overflow-hidden"
          style="background: linear-gradient(180deg, rgba(28, 25, 23, 0.95), var(--c-surface));"
        >
          <!-- Top Row: Title, Badge & Hide Button -->
          <div class="flex items-start justify-between pb-3.5 border-b border-white/10 gap-4 flex-wrap">
            <div>
              <div class="flex items-center gap-2 flex-wrap">
                <span class="text-xl">📐</span>
                <h3 class="text-base sm:text-lg font-bold text-white">
                  Beer-Equivalent Calculation Logic
                </h3>
                <span class="text-xs px-2.5 py-0.5 rounded-full bg-amber-500/20 text-amber-300 border border-amber-500/30 font-medium">
                  Dynamic Physical Standard
                </span>
              </div>
              <p class="text-xs text-muted mt-1">
                Every beverage in this tracker is normalized into standard beer units based on exact pure ethanol volume.
              </p>
            </div>
            <button
              type="button"
              class="text-xs px-2.5 py-1 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300 hover:text-white border border-white/10 transition-colors cursor-pointer flex items-center gap-1"
              @click="showFormulaGuide = false"
            >
              <span>✕</span>
              <span>Hide Guide</span>
            </button>
          </div>

          <!-- 3 Educational Cards -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3.5 my-4">
            <!-- 1. The Benchmark -->
            <div class="p-3.5 rounded-xl bg-white/5 border border-white/10 flex flex-col justify-between">
              <div>
                <div class="text-[11px] font-bold uppercase tracking-wider text-amber-400 mb-1 flex items-center gap-1.5">
                  <span>🍺</span>
                  <span>1. The Benchmark</span>
                </div>
                <div class="text-white text-xs sm:text-sm font-semibold mb-1">Standard Halbe (500 ml @ 5.0%)</div>
                <p class="text-xs text-muted leading-relaxed">
                  1.0 beer equivalent equals exactly <strong>25 ml</strong> of pure ethanol:
                </p>
              </div>
              <div class="mt-3 p-2 rounded-lg bg-black/40 font-mono text-xs text-amber-300 text-center border border-amber-500/20">
                500 ml × 0.05 = 25 ml alcohol
              </div>
            </div>

            <!-- 2. The Formula -->
            <div class="p-3.5 rounded-xl bg-white/5 border border-white/10 flex flex-col justify-between">
              <div>
                <div class="text-[11px] font-bold uppercase tracking-wider text-indigo-400 mb-1 flex items-center gap-1.5">
                  <span>⚡</span>
                  <span>2. Dynamic Formula</span>
                </div>
                <div class="text-white text-xs sm:text-sm font-semibold mb-1">Direct Proportional Scaling</div>
                <p class="text-xs text-muted leading-relaxed">
                  Divide the drink's total ethanol by the 25 ml reference baseline:
                </p>
              </div>
              <div class="mt-3 p-2 rounded-lg bg-black/40 font-mono text-xs text-indigo-300 text-center border border-indigo-500/20">
                Beer Eq = (Volume × ABV %) ÷ 2500
              </div>
            </div>

            <!-- 3. BAC Physical Link -->
            <div class="p-3.5 rounded-xl bg-white/5 border border-white/10 flex flex-col justify-between">
              <div>
                <div class="text-[11px] font-bold uppercase tracking-wider text-emerald-400 mb-1 flex items-center gap-1.5">
                  <span>🩸</span>
                  <span>3. BAC / Promille Link</span>
                </div>
                <div class="text-white text-xs sm:text-sm font-semibold mb-1">Widmark Ethanol Mass</div>
                <p class="text-xs text-muted leading-relaxed">
                  With alcohol density ρ = 0.8 g/ml, each 1.0 beer injects exactly <strong>20.0 g</strong> pure alcohol into your BAC curve.
                </p>
              </div>
              <div class="mt-3 p-2 rounded-lg bg-black/40 font-mono text-xs text-emerald-300 text-center border border-emerald-500/20">
                1.0 🍺 = 25 ml × 0.8 = 20.0 g ethanol
              </div>
            </div>
          </div>

          <!-- Interactive Live Converter Playground -->
          <div class="p-4 rounded-xl bg-black/40 border border-white/10 mb-4">
            <div class="flex items-center justify-between mb-3 flex-wrap gap-2">
              <div class="text-xs font-bold uppercase tracking-wider text-white flex items-center gap-1.5">
                <span>🧮 Interactive Live Converter</span>
                <span class="text-muted font-normal text-[11px]">— Test any volume & ABV live</span>
              </div>
              <div class="flex items-center gap-1.5 flex-wrap">
                <button
                  v-for="qp in guideQuickPillPresets"
                  :key="qp.label"
                  type="button"
                  class="text-[11px] px-2 py-0.5 rounded-md bg-white/5 hover:bg-white/15 border border-white/10 text-slate-300 hover:text-white transition-colors cursor-pointer"
                  @click="setGuidePreset(qp.volume, qp.abv)"
                >
                  {{ qp.label }}
                </button>
              </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-12 gap-3 items-center">
              <!-- Volume input -->
              <div class="sm:col-span-4">
                <label class="block text-[11px] font-semibold text-muted mb-1">Volume (ml)</label>
                <div class="relative">
                  <input
                    v-model.number="guideCalcVolume"
                    type="number"
                    min="1"
                    max="5000"
                    class="input text-sm py-1.5 pr-8"
                    placeholder="40"
                  />
                  <span class="absolute right-3 top-2 text-xs text-muted">ml</span>
                </div>
              </div>

              <!-- ABV input -->
              <div class="sm:col-span-4">
                <label class="block text-[11px] font-semibold text-muted mb-1">Alcohol (% ABV)</label>
                <div class="relative">
                  <input
                    v-model.number="guideCalcAbv"
                    type="number"
                    step="0.1"
                    min="0"
                    max="100"
                    class="input text-sm py-1.5 pr-8"
                    placeholder="40.0"
                  />
                  <span class="absolute right-3 top-2 text-xs text-muted">%</span>
                </div>
              </div>

              <!-- Live Output Card -->
              <div class="sm:col-span-4 p-3 rounded-xl bg-white/5 border border-white/10 flex items-center justify-between">
                <div>
                  <div class="text-[10px] uppercase font-bold text-muted">Result</div>
                  <div class="text-2xl font-extrabold text-amber-300">
                    = {{ guideCalcBeerEq }} 🍺
                  </div>
                </div>
                <div class="text-right text-xs">
                  <div class="text-white font-mono font-semibold">{{ guideCalcPureMl }} ml alcohol</div>
                  <div class="text-muted text-[11px]">{{ guideCalcGrams }} g ethanol</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Everyday Drinks Reference Breakdown -->
          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-muted mb-2">
              📋 Everyday Drinks Reference Guide
            </div>
            <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-8 gap-2">
              <div
                v-for="item in referenceExamples"
                :key="item.name"
                class="p-2 rounded-xl bg-white/5 border border-white/5 flex flex-col items-center text-center hover:bg-white/10 transition-colors"
              >
                <span class="text-xl mb-0.5">{{ item.emoji }}</span>
                <span class="text-xs font-bold text-white truncate w-full">{{ item.name }}</span>
                <span class="text-[10px] text-muted">{{ item.vol }}ml @ {{ item.abv }}%</span>
                <span class="mt-1.5 text-xs font-extrabold text-amber-300 px-1.5 py-0.5 rounded bg-amber-500/10 border border-amber-500/20">
                  {{ item.beerEq }} 🍺
                </span>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Presets Grid -->
      <div v-if="loadingPresets" class="flex justify-center py-10">
        <div class="spinner" style="width:32px;height:32px;border-width:3px"></div>
      </div>

      <div v-else class="space-y-4">
        <!-- 1. Default Standard Presets (Shown by default) -->
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
          <button
            v-for="dt in standardPresets"
            :key="dt.id"
            type="button"
            class="drink-btn group cursor-pointer"
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
              ≈ {{ getDrinkBeerEq(dt) }} 🍺
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

        <!-- 2. Party Custom Drinks (Hidden by default, optionally shown) -->
        <Transition name="expand">
          <div v-if="showCustomPartyDrinks && partyCustomDrinks.length > 0" class="pt-3 border-t border-white/5">
            <div class="text-xs font-semibold text-purple-300 uppercase tracking-wider mb-2.5 flex items-center justify-between">
              <span class="flex items-center gap-1.5">
                <span>🪄 Custom Drinks for Party {{ code }}</span>
                <span class="text-muted font-normal text-[11px]">(Only drinks from this party)</span>
              </span>
              <button
                type="button"
                class="text-muted hover:text-white text-xs underline cursor-pointer"
                @click="showCustomPartyDrinks = false"
              >
                Hide
              </button>
            </div>
            <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
              <button
                v-for="dt in partyCustomDrinks"
                :key="dt.name"
                type="button"
                class="drink-btn group border-purple-500/30 hover:border-purple-400 bg-purple-950/20"
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
                  {{ dt.unitName || 'Custom' }}
                </div>
                <div class="mt-2 text-xs font-semibold px-2 py-0.5 rounded-full bg-purple-500/20 border border-purple-500/40 text-purple-300">
                  ≈ {{ getDrinkBeerEq(dt) }} 🍺
                </div>
              </button>
            </div>
          </div>
        </Transition>
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

    <!-- Analytics & Big Charts Section -->
    <section class="card mb-12 overflow-hidden">
      <!-- Chart Navigation Bar (Large Tabs) -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-5 mb-6 border-b border-white/5">
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-xl font-bold text-white">📊 Analytics & Timeline</h2>
            <span class="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-indigo-500/15 text-indigo-300 border border-indigo-500/30">
              Interactive View
            </span>
          </div>
          <p class="text-xs text-muted mt-1">
            Track blood alcohol progression (Verlauf), cumulative consumption, and attendee ranking
          </p>
        </div>

        <!-- View Switcher Tabs -->
        <div class="flex items-center gap-1.5 p-1 rounded-2xl bg-black/30 border border-white/5 overflow-x-auto text-xs">
          <button
            type="button"
            class="chart-tab-btn"
            :class="{ 'chart-tab-btn-active': activeChartTab === 'bac' }"
            @click="switchChartTab('bac')"
          >
            <span>🩸 Blood Alcohol Verlauf</span>
          </button>
          <button
            type="button"
            class="chart-tab-btn"
            :class="{ 'chart-tab-btn-active': activeChartTab === 'timeline' }"
            @click="switchChartTab('timeline')"
          >
            <span>📈 Consumption Timeline</span>
          </button>
          <button
            type="button"
            class="chart-tab-btn"
            :class="{ 'chart-tab-btn-active': activeChartTab === 'leaderboard' }"
            @click="switchChartTab('leaderboard')"
          >
            <span>🏆 Leaderboard</span>
          </button>
          <button
            type="button"
            class="chart-tab-btn"
            :class="{ 'chart-tab-btn-active': activeChartTab === 'all' }"
            @click="switchChartTab('all')"
          >
            <span>🔲 Show All</span>
          </button>
        </div>
      </div>

      <!-- Tab 1: Blood Alcohol Concentration (BAC) Verlauf -->
      <div v-show="activeChartTab === 'bac' || activeChartTab === 'all'" class="mb-8">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4">
          <div class="flex items-center gap-3 flex-wrap">
            <h3 class="text-base font-bold text-white flex items-center gap-2">
              <span>🩸 Blood Alcohol (BAC) Verlauf</span>
              <span class="text-xs font-normal text-muted">Absorption & elimination curve over time</span>
            </h3>

            <!-- Attendee selector for BAC -->
            <div class="flex items-center gap-1.5 bg-white/5 border border-white/10 rounded-xl px-2.5 py-1 text-xs">
              <span class="text-muted">Attendee:</span>
              <select v-model="bacTargetAttendee" class="bg-transparent text-white font-semibold outline-none cursor-pointer">
                <option v-for="a in attendeeList" :key="a" :value="a" class="bg-slate-900 text-white">
                  {{ a }}{{ a === currentUser ? ' (You)' : '' }}
                </option>
              </select>
            </div>
          </div>

          <!-- Profile / Settings Controls -->
          <div class="flex items-center gap-2 flex-wrap">
            <button
              type="button"
              class="px-2.5 py-1 rounded-xl text-xs font-semibold bg-white/5 border border-white/10 hover:border-white/20 text-slate-300 transition-colors flex items-center gap-1.5"
              @click="showBacSettings = !showBacSettings"
            >
              <span>⚙️ Body: {{ bacWeightKg }}kg • {{ bacGender === 'male' ? '♂ Male' : '♀ Female' }}</span>
              <span class="text-[10px] text-muted">{{ showBacSettings ? '▲' : '▼' }}</span>
            </button>

            <!-- Unit toggle: % BAC vs ‰ Promille -->
            <button
              type="button"
              class="px-2.5 py-1 rounded-xl text-xs font-semibold transition-colors border"
              :class="bacUnit === 'promille' ? 'bg-indigo-600/30 text-indigo-300 border-indigo-500/40' : 'bg-white/5 text-slate-300 border-white/10'"
              @click="toggleBacUnit"
            >
              Unit: {{ bacUnit === 'promille' ? '‰ Promille' : '% BAC' }}
            </button>
          </div>
        </div>

        <!-- Collapsible BAC Body Profile Settings -->
        <Transition name="expand">
          <div
            v-if="showBacSettings"
            class="p-4 rounded-2xl mb-5 border border-white/10 flex flex-wrap items-center gap-6"
            style="background:var(--c-bg)"
          >
            <div>
              <div class="text-xs font-semibold text-muted mb-1.5">Biological Sex (Widmark ratio)</div>
              <div class="flex gap-2">
                <button
                  type="button"
                  class="px-3 py-1 rounded-xl text-xs font-semibold border transition-all"
                  :class="bacGender === 'male' ? 'bg-indigo-600 text-white border-indigo-400' : 'bg-white/5 text-slate-300 border-white/10'"
                  @click="bacGender = 'male'"
                >
                  ♂ Male (r = 0.68)
                </button>
                <button
                  type="button"
                  class="px-3 py-1 rounded-xl text-xs font-semibold border transition-all"
                  :class="bacGender === 'female' ? 'bg-indigo-600 text-white border-indigo-400' : 'bg-white/5 text-slate-300 border-white/10'"
                  @click="bacGender = 'female'"
                >
                  ♀ Female (r = 0.55)
                </button>
              </div>
            </div>

            <div>
              <div class="text-xs font-semibold text-muted mb-1.5">Body Weight (kg)</div>
              <div class="flex items-center gap-2">
                <input
                  v-model.number="bacWeightKg"
                  type="number"
                  min="40"
                  max="200"
                  class="input py-1 px-3 w-24 text-center font-bold"
                />
                <span class="text-xs text-muted">kg</span>
              </div>
            </div>

            <div class="text-xs text-muted max-w-sm">
              💡 The Widmark equation calculates peak absorption (30m window) and linear metabolic elimination (~0.15‰ / 0.015% per hour).
            </div>
          </div>
        </Transition>

        <!-- Current Target BAC Summary Bar -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 p-3.5 rounded-2xl mb-4 border border-white/5 bg-black/20">
          <div>
            <div class="text-[11px] text-muted font-semibold uppercase">Current BAC</div>
            <div class="text-lg sm:text-xl font-bold" :class="bacColorClass(selectedTargetBAC.currentBAC)">
              {{ bacUnit === 'promille' ? formatPromille(selectedTargetBAC.currentBAC) : formatBAC(selectedTargetBAC.currentBAC) }}
              <span class="text-xs font-normal text-muted">({{ bacUnit === 'promille' ? formatBAC(selectedTargetBAC.currentBAC) : formatPromille(selectedTargetBAC.currentBAC) }})</span>
            </div>
          </div>

          <div>
            <div class="text-[11px] text-muted font-semibold uppercase">Peak Reached</div>
            <div class="text-lg sm:text-xl font-bold text-amber-300">
              {{ bacUnit === 'promille' ? formatPromille(selectedTargetBAC.peakBAC) : formatBAC(selectedTargetBAC.peakBAC) }}
              <span v-if="selectedTargetBAC.peakTime" class="text-xs font-normal text-muted">
                at {{ formatTimeOnly(selectedTargetBAC.peakTime) }}
              </span>
            </div>
          </div>

          <div>
            <div class="text-[11px] text-muted font-semibold uppercase">Status</div>
            <div class="text-sm font-semibold truncate text-slate-200 mt-0.5">
              {{ getSobrietyStage(selectedTargetBAC.currentBAC).label }}
            </div>
          </div>

          <div>
            <div class="text-[11px] text-muted font-semibold uppercase">Est. Sober At</div>
            <div class="text-sm font-semibold truncate text-emerald-400 mt-0.5">
              {{ selectedTargetBAC.soberTime ? `~${formatTimeOnly(selectedTargetBAC.soberTime)} (${formatSoberEta(selectedTargetBAC.soberTime)})` : 'Sober now' }}
            </div>
          </div>
        </div>

        <!-- Big BAC Verlauf Canvas -->
        <div v-if="!targetAttendeeDrinks.length" class="flex flex-col items-center justify-center py-16 text-center text-muted">
          <div class="text-4xl mb-2">🩸</div>
          <p class="text-sm font-semibold text-white">No drinks logged for {{ bacTargetAttendee }}</p>
          <p class="text-xs text-muted mt-1">Log a drink above to calculate and display blood alcohol progression.</p>
        </div>

        <div v-else class="relative w-full chart-container">
          <canvas ref="bacChartEl" style="height:380px;max-height:420px;width:100%"></canvas>
        </div>
      </div>

      <!-- Tab 2: Consumption Timeline (Beers Over Time) -->
      <div v-show="activeChartTab === 'timeline' || activeChartTab === 'all'" class="mb-8" :class="{ 'pt-6 border-t border-white/5': activeChartTab === 'all' }">
        <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
          <div>
            <h3 class="text-base font-bold text-white flex items-center gap-2">
              <span>📈 Consumption Over Time</span>
              <span class="text-xs font-normal text-muted">Real-world pacing & cumulative beers</span>
            </h3>
          </div>
          <span class="text-xs px-2.5 py-1 rounded-full font-semibold bg-white/5 text-muted border border-white/10">
            All Attendees
          </span>
        </div>

        <div v-if="!rawBeerLog.length" class="flex flex-col items-center justify-center py-16 text-center text-muted">
          <div class="text-4xl mb-2">⏱️</div>
          <p class="text-sm font-semibold text-white">Timeline is empty</p>
          <p class="text-xs text-muted mt-1">Timestamps and pacing will appear once drinks are logged.</p>
        </div>

        <div v-else class="relative w-full chart-container">
          <canvas ref="lineChartEl" style="height:380px;max-height:420px;width:100%"></canvas>
        </div>
      </div>

      <!-- Tab 3: Attendee Leaderboard Bar Chart -->
      <div v-show="activeChartTab === 'leaderboard' || activeChartTab === 'all'" :class="{ 'pt-6 border-t border-white/5': activeChartTab === 'all' }">
        <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
          <div>
            <h3 class="text-base font-bold text-white flex items-center gap-2">
              <span>🏆 Total Beer-Equivalents Leaderboard</span>
              <span class="text-xs font-normal text-muted">Ranked by total alcohol units</span>
            </h3>
          </div>
          <span class="text-xs px-2.5 py-1 rounded-full font-semibold bg-white/5 text-muted border border-white/10">
            {{ attendeeLeaderboard.length }} participants
          </span>
        </div>

        <div v-if="!rawBeerLog.length" class="flex flex-col items-center justify-center py-16 text-center text-muted">
          <div class="text-4xl mb-2">🍺</div>
          <p class="text-sm font-semibold text-white">No drinks recorded yet</p>
          <p class="text-xs text-muted mt-1">Tap any drink above to record the first one!</p>
        </div>

        <div v-else class="relative w-full chart-container">
          <canvas ref="barChartEl" style="height:380px;max-height:420px;width:100%"></canvas>
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

      <div v-else class="flex flex-col gap-2.5 max-h-[420px] overflow-y-auto pr-1">
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
import { getParty, listDrinks, createDrink, addBeer, joinParty, client } from '@/api-client'
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
  LineController,
  Filler
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
  LineController,
  Filler
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
const showCustomPartyDrinks = ref(false)
const showFormulaGuide = ref(false)
const localPartyCustomDrinkNames = ref<string[]>([])
const toastMsg = ref<string | null>(null)
let toastTimer: any = null

const DEFAULT_PRESET_NAMES = new Set([
  'beer',
  'bier',
  'vodka',
  'wodka',
  'wein',
  'wine',
  'spritzer',
  'gspritzter',
  'gin',
  'jägermeister',
  'jaegermeister',
  'jager',
  'jägermeister orange',
  'jaegermeister orange',
  'jager orange',
  'jägermeister winter',
  'jaegermeister winter',
  'jager winter',
  'berliner luft',
  'klopfer',
  'nussschnaps',
  'nuss schnaps',
  'tatra tea black',
  'tatra tea',
  'tatra',
  'weißer rum',
  'weisser rum',
  'white rum',
  'rum',
  'xuxu',
  'radler',
  'cocktail',
  'shot',
  'tequila',
  'whiskey',
  'whisky'
])

function isDefaultPreset(drink: DrinkType | string): boolean {
  if (typeof drink === 'object' && drink !== null) {
    if (typeof drink.id === 'number' && drink.id > 0 && drink.id <= 14) {
      return true
    }
    return isDefaultPreset(drink.name || '')
  }
  const clean = (drink || '').toLowerCase().trim()
  return DEFAULT_PRESET_NAMES.has(clean)
}

// 1. Standard Presets (Shown by default)
const standardPresets = computed(() => {
  return presets.value.filter(dt => isDefaultPreset(dt))
})

// 2. Custom drinks belonging ONLY to THIS party
const partyCustomDrinks = computed(() => {
  const partyNames = new Set<string>()

  // Names logged in this party's beerLog
  for (const entry of rawBeerLog.value) {
    const n = (entry.drinkName || '').trim()
    if (n && !isDefaultPreset(n)) {
      partyNames.add(n.toLowerCase())
    }
  }

  // Names created locally in this party
  for (const n of localPartyCustomDrinkNames.value) {
    if (n && !isDefaultPreset(n)) {
      partyNames.add(n.toLowerCase().trim())
    }
  }

  const result: (DrinkType | { id?: number; name: string; beerEquivalent: number; unitName?: string; volumeMl?: number; alcoholPercent?: number })[] = []
  const added = new Set<string>()

  // Look for matching drink types from presets.value
  for (const dt of presets.value) {
    const lower = (dt.name || '').toLowerCase().trim()
    if (partyNames.has(lower) && !added.has(lower)) {
      result.push(dt)
      added.add(lower)
    }
  }

  // Fallback: If logged in beerLog but not returned in presets list
  for (const entry of rawBeerLog.value) {
    const lower = (entry.drinkName || '').toLowerCase().trim()
    if (partyNames.has(lower) && !added.has(lower)) {
      result.push({
        id: -entry.id,
        name: entry.drinkName,
        beerEquivalent: entry.beerEquivalent,
        unitName: 'Custom'
      })
      added.add(lower)
    }
  }

  return result
})

// ── Dynamic Beer-Equivalent Calculation & Guide State ───────────────────────
function getDrinkBeerEq(dt: { volumeMl?: number; alcoholPercent?: number; beerEquivalent?: number }): number {
  if (dt.volumeMl && dt.volumeMl > 0 && dt.alcoholPercent && dt.alcoholPercent > 0) {
    return parseFloat(((dt.volumeMl * dt.alcoholPercent) / 2500).toFixed(2))
  }
  return Number(dt.beerEquivalent ?? 1.0)
}

// Live interactive converter inputs
const guideCalcVolume = ref<number>(40)
const guideCalcAbv = ref<number>(40.0)

const guideCalcPureMl = computed(() => {
  const vol = guideCalcVolume.value || 0
  const abv = guideCalcAbv.value || 0
  return parseFloat(((vol * abv) / 100).toFixed(1))
})

const guideCalcGrams = computed(() => {
  return parseFloat((guideCalcPureMl.value * 0.8).toFixed(1))
})

const guideCalcBeerEq = computed(() => {
  const vol = guideCalcVolume.value || 0
  const abv = guideCalcAbv.value || 0
  if (vol <= 0 || abv <= 0) return 0
  return parseFloat(((vol * abv) / 2500).toFixed(2))
})

interface GuideQuickPreset {
  label: string
  volume: number
  abv: number
}

const guideQuickPillPresets: GuideQuickPreset[] = [
  { label: '💥 Klopfer (20ml @ 16%)', volume: 20, abv: 16.0 },
  { label: '🥃 Shot (40ml @ 40%)', volume: 40, abv: 40.0 },
  { label: '🍋 Radler (500ml @ 2.5%)', volume: 500, abv: 2.5 },
  { label: '🥫 Beer Can (330ml @ 5%)', volume: 330, abv: 5.0 },
  { label: '🍺 Halbe (500ml @ 5%)', volume: 500, abv: 5.0 },
  { label: '🍷 Wine Glass (200ml @ 12.5%)', volume: 200, abv: 12.5 },
  { label: '🍹 Cocktail (300ml @ 15%)', volume: 300, abv: 15.0 },
]

function setGuidePreset(vol: number, abv: number) {
  guideCalcVolume.value = vol
  guideCalcAbv.value = abv
}

const referenceExamples = [
  { name: 'Klopfer', emoji: '💥', vol: 20, abv: 16.0, beerEq: '0.13' },
  { name: 'Schnaps Shot', emoji: '🥃', vol: 20, abv: 35.0, beerEq: '0.28' },
  { name: 'Standard Shot', emoji: '🥃', vol: 40, abv: 40.0, beerEq: '0.64' },
  { name: 'Small Beer', emoji: '🥫', vol: 330, abv: 5.0, beerEq: '0.66' },
  { name: 'Radler', emoji: '🍋', vol: 500, abv: 2.5, beerEq: '0.50' },
  { name: 'Standard Halbe', emoji: '🍺', vol: 500, abv: 5.0, beerEq: '1.00' },
  { name: 'Glass of Wine', emoji: '🍷', vol: 200, abv: 12.5, beerEq: '1.00' },
  { name: 'Strong Cocktail', emoji: '🍹', vol: 300, abv: 15.0, beerEq: '1.80' },
]

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

const attendeeList = computed(() => {
  const attendees = party.value?.attendees ?? []
  if (attendees.length > 0) return attendees
  return [...new Set(rawBeerLog.value.map(e => e.attendee))]
})

// Log filter
function isSameUser(a?: string, b?: string): boolean {
  if (!a || !b) return false
  return a.trim().toLowerCase() === b.trim().toLowerCase()
}

// Log filter
const logFilter = ref<'all' | 'me'>('all')

const filteredLog = computed(() => {
  const list = [...rawBeerLog.value].reverse() // latest first
  if (logFilter.value === 'me' && currentUser.value) {
    return list.filter(e => isSameUser(e.attendee, currentUser.value))
  }
  return list
})

// ── Metrics ─────────────────────────────────────────────────────────────────
const totalBeerEquivalents = computed(() => {
  return rawBeerLog.value.reduce((sum, e) => sum + (e.beerEquivalent || 0), 0)
})

const myDrinksCount = computed(() => {
  if (!currentUser.value) return 0
  return rawBeerLog.value.filter(e => isSameUser(e.attendee, currentUser.value)).length
})

const myBeersTotal = computed(() => {
  if (!currentUser.value) return 0
  return rawBeerLog.value
    .filter(e => isSameUser(e.attendee, currentUser.value))
    .reduce((sum, e) => sum + (e.beerEquivalent || 0), 0)
})

const myBeerSharePercent = computed(() => {
  if (totalBeerEquivalents.value === 0) return 0
  return Math.round((myBeersTotal.value / totalBeerEquivalents.value) * 100)
})

const attendeeLeaderboard = computed(() => {
  const allAttendees = attendeeList.value
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

// ── Blood Alcohol Concentration (BAC / Widmark Model) ───────────────────────
const bacGender = ref<'male' | 'female'>('male')
const bacWeightKg = ref<number>(75)
const bacUnit = ref<'promille' | 'percent'>('promille')
const showBacSettings = ref<boolean>(false)
const bacTargetAttendee = ref<string>('')

async function ensureJoinedParty() {
  const codeVal = code.value
  const user = currentUser.value
  if (!codeVal || !user) return

  if (party.value?.attendees && party.value.attendees.some(a => isSameUser(a, user))) {
    return
  }

  try {
    await joinParty({
      path: { code: codeVal, username: user }
    })
    if (party.value) {
      if (!party.value.attendees) party.value.attendees = []
      if (!party.value.attendees.some(a => isSameUser(a, user))) {
        party.value.attendees.push(user)
      }
    }
  } catch (err) {
    console.error('Failed to auto-join party:', err)
  }
}

// Target attendee defaults to currentUser whenever user is logged in
watch(currentUser, (user) => {
  if (user) {
    bacTargetAttendee.value = user
    ensureJoinedParty()
  }
}, { immediate: true })

watch(attendeeList, (list) => {
  if (!bacTargetAttendee.value && list.length > 0) {
    bacTargetAttendee.value = currentUser.value || list[0] || ''
  }
}, { immediate: true })

function toggleBacUnit() {
  bacUnit.value = bacUnit.value === 'promille' ? 'percent' : 'promille'
}

interface BACPoint {
  time: Date
  timeLabel: string
  bac: number // %
  promille: number // ‰
  isProjected: boolean
}

interface BACResult {
  points: BACPoint[]
  currentBAC: number
  peakBAC: number
  peakTime: Date | null
  soberTime: Date | null
}

function calculateBACVerlauf(
  drinks: BeerEntry[],
  weightKg: number,
  gender: 'male' | 'female'
): BACResult {
  if (!drinks || !drinks.length) {
    return { points: [], currentBAC: 0, peakBAC: 0, peakTime: null, soberTime: null }
  }

  const r = gender === 'female' ? 0.55 : 0.68
  const safeWeight = Math.max(30, weightKg || 75)
  const betaPerHour = 0.015 // ~0.015% per hour elimination rate

  // Sort drinks chronologically
  const sorted = [...drinks].sort(
    (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
  )

  const nowMs = Date.now()
  const firstDrinkMs = new Date(sorted[0]?.createdAt ?? Date.now()).getTime()

  // Track simulation state
  let runningBAC = 0
  let lastEventMs = firstDrinkMs
  let peakBAC = 0
  let peakTime: Date | null = null

  const rawPoints: { timeMs: number; bac: number; isProjected: boolean }[] = []

  function addPoint(tMs: number, val: number, isProj: boolean) {
    rawPoints.push({
      timeMs: tMs,
      bac: Math.max(0, val),
      isProjected: isProj
    })
  }

  // 1. Initial zero baseline point 5 min before first drink
  addPoint(firstDrinkMs - 5 * 60 * 1000, 0, false)

  // 2. Step through each drink in chronological order
  for (const drink of sorted) {
    const drinkMs = new Date(drink.createdAt).getTime()

    // If time elapsed between previous event and this drink, metabolize alcohol
    if (drinkMs > lastEventMs) {
      const elapsedHours = (drinkMs - lastEventMs) / 3600000
      if (drinkMs - lastEventMs > 20 * 60 * 1000) {
        const midMs = lastEventMs + (drinkMs - lastEventMs) / 2
        const midHours = (midMs - lastEventMs) / 3600000
        const midBAC = Math.max(0, runningBAC - betaPerHour * midHours)
        addPoint(midMs, midBAC, midMs > nowMs)
      }

      runningBAC = Math.max(0, runningBAC - betaPerHour * elapsedHours)
      addPoint(drinkMs - 1000, runningBAC, drinkMs > nowMs)
    }

    // Add this drink's alcohol content
    // 1.0 beer equivalent = 500ml @ 5% = 25ml ethanol * 0.8 g/ml = 20g pure alcohol
    const beerEq = Number(drink.beerEquivalent) > 0 ? Number(drink.beerEquivalent) : 1.0
    const alcoholGrams = beerEq * 20
    const bacBoost = alcoholGrams / (safeWeight * r * 10) // in % BAC

    runningBAC += bacBoost
    lastEventMs = drinkMs

    if (runningBAC > peakBAC) {
      peakBAC = runningBAC
      peakTime = new Date(drinkMs)
    }

    // Record point right after drinking
    addPoint(drinkMs, runningBAC, drinkMs > nowMs)
  }

  // 3. Bring forward to current time (nowMs)
  let currentBAC = runningBAC
  if (nowMs > lastEventMs) {
    const elapsedHours = (nowMs - lastEventMs) / 3600000
    currentBAC = Math.max(0, runningBAC - betaPerHour * elapsedHours)
    addPoint(nowMs, currentBAC, false)
  }

  // 4. Project forward into future until BAC reaches 0.00%
  let soberTime: Date | null = null

  if (currentBAC > 0.0001) {
    const hoursToZero = currentBAC / betaPerHour
    const soberMs = nowMs + hoursToZero * 3600000
    soberTime = new Date(soberMs)

    // Add intermediate projection points for smooth curve
    const steps = hoursToZero > 2 ? 4 : 2
    for (let i = 1; i < steps; i++) {
      const pMs = nowMs + (hoursToZero * i / steps) * 3600000
      const pElapsedHours = (pMs - nowMs) / 3600000
      const pBAC = Math.max(0, currentBAC - betaPerHour * pElapsedHours)
      addPoint(pMs, pBAC, true)
    }

    // Final sober zero point
    addPoint(soberMs, 0, true)
  } else {
    soberTime = new Date(nowMs)
  }

  rawPoints.sort((a, b) => a.timeMs - b.timeMs)

  const points: BACPoint[] = rawPoints.map(pt => {
    const d = new Date(pt.timeMs)
    return {
      time: d,
      timeLabel: d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      bac: parseFloat(pt.bac.toFixed(4)),
      promille: parseFloat((pt.bac * 10).toFixed(3)),
      isProjected: pt.isProjected
    }
  })

  return {
    points,
    currentBAC: parseFloat(currentBAC.toFixed(4)),
    peakBAC: parseFloat(peakBAC.toFixed(4)),
    peakTime: peakTime || new Date(firstDrinkMs),
    soberTime
  }
}

// User BAC summary for the stats ribbon
const userDrinks = computed(() => {
  if (!currentUser.value) return []
  return rawBeerLog.value.filter(e => isSameUser(e.attendee, currentUser.value))
})

const userBACSummary = computed(() => {
  return calculateBACVerlauf(userDrinks.value, bacWeightKg.value, bacGender.value)
})

// Target attendee drinks & calculation (for the chart)
const targetAttendeeDrinks = computed(() => {
  const target = bacTargetAttendee.value || currentUser.value
  if (!target) return []
  return rawBeerLog.value.filter(e => isSameUser(e.attendee, target))
})

const selectedTargetBAC = computed(() => {
  return calculateBACVerlauf(targetAttendeeDrinks.value, bacWeightKg.value, bacGender.value)
})

function formatBAC(bac: number) {
  if (bac <= 0.0001) return '0.00%'
  return bac < 0.1 ? `${bac.toFixed(3)}%` : `${bac.toFixed(2)}%`
}

function formatPromille(bac: number) {
  if (bac <= 0.0001) return '0.00‰'
  return `${(bac * 10).toFixed(2)}‰`
}

function getSobrietyStage(bac: number) {
  if (bac <= 0.005) return { label: 'Sober 🟢', color: '#22c55e' }
  if (bac <= 0.03) return { label: 'Mild Relaxation 🟡', color: '#eab308' }
  if (bac <= 0.05) return { label: 'Buzzed / Euphoric 🍺', color: '#f59e0b' }
  if (bac <= 0.08) return { label: 'At EU Driving Limit 🟠', color: '#f97316' }
  if (bac <= 0.15) return { label: 'Intoxicated 🔴', color: '#ef4444' }
  return { label: 'Heavy Intoxication ⛔', color: '#dc2626' }
}

function bacColorClass(bac: number) {
  if (bac <= 0.005) return 'text-emerald-400'
  if (bac <= 0.05) return 'text-amber-300'
  if (bac <= 0.08) return 'text-orange-400'
  return 'text-rose-400'
}

function bacDotClass(bac: number) {
  if (bac <= 0.005) return 'bg-emerald-400'
  if (bac <= 0.05) return 'bg-amber-300'
  if (bac <= 0.08) return 'bg-orange-400'
  return 'bg-rose-400'
}

function formatSoberEta(d: Date | null) {
  if (!d) return 'Sober'
  const diffMs = d.getTime() - Date.now()
  if (diffMs <= 0) return 'Sober'
  const totalMins = Math.round(diffMs / 60000)
  const h = Math.floor(totalMins / 60)
  const m = totalMins % 60
  if (h === 0) return `In ~${m}m`
  return `In ~${h}h ${m}m`
}

function formatTimeOnly(d: Date | null) {
  if (!d) return ''
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

// ── Chart Tab Navigation ────────────────────────────────────────────────────
type ChartTab = 'bac' | 'timeline' | 'leaderboard' | 'all'
const activeChartTab = ref<ChartTab>('bac')

function switchChartTab(tab: ChartTab) {
  activeChartTab.value = tab
  nextTick(() => renderCharts())
}

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
  'jägermeister orange': '🍊',
  'jaegermeister orange': '🍊',
  'jager orange': '🍊',
  'jägermeister winter': '❄️',
  'jaegermeister winter': '❄️',
  'jager winter': '❄️',
  jägermeister: '🦌',
  jaegermeister: '🦌',
  jager: '🦌',
  'berliner luft': '💨',
  'tatra tea': '🫖',
  tatra: '🫖',
  nussschnaps: '🌰',
  'nuss schnaps': '🌰',
  'weißer rum': '🥃',
  'weisser rum': '🥃',
  'white rum': '🥃',
  xuxu: '🍓',
  beer: '🍺',
  bier: '🍺',
  vodka: '🧊',
  wodka: '🧊',
  wein: '🍷',
  wine: '🍷',
  spritzer: '🍾',
  gspritzter: '🍾',
  sekt: '🍾',
  champagne: '🍾',
  'gin tonic': '🟢',
  gin: '🟢',
  rum: '🍹',
  cocktail: '🍹',
  klopfer: '💥',
  shot: '🥃',
  tequila: '🌵',
  whiskey: '🥃',
  whisky: '🥃',
  cider: '🍏',
  radler: '🍋',
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
async function logDrink(dt: DrinkType | { name: string; beerEquivalent?: number; unitName?: string; volumeMl?: number; alcoholPercent?: number }) {
  const attendee = currentUser.value
  if (!attendee) {
    alert('Please log in with a username first!')
    return
  }

  const drinkName = dt.name || 'Drink'
  const beerEq = getDrinkBeerEq(dt)
  loggingDrink.value = true

  // Ensure attendee is registered in the party
  if (party.value && !party.value.attendees.some(a => isSameUser(a, attendee))) {
    await ensureJoinedParty()
  }

  // Optimistically add to party.value.beerLog with a fresh array reference for instant UI feedback!
  if (party.value) {
    const newEntry: BeerEntry = {
      id: Date.now(),
      attendee,
      drinkName,
      beerEquivalent: beerEq,
      createdAt: new Date().toISOString()
    }
    const currentLog = party.value.beerLog ? [...party.value.beerLog] : []
    currentLog.push(newEntry)
    party.value.beerLog = currentLog

    party.value.beerCount = (party.value.beerCount || 0) + 1
    if (!party.value.beerPerAttendee) party.value.beerPerAttendee = {}
    party.value.beerPerAttendee[attendee] = (party.value.beerPerAttendee[attendee] || 0) + 1
  }

  // Immediately re-render charts so the graph and stats update instantaneously (<1ms)
  nextTick(() => renderCharts())

  try {
    const resp = await addBeer({
      path: { code: code.value },
      body: {
        attendee,
        drinkName,
        beerEquivalent: beerEq,
      }
    })

    if (resp.error) {
      console.error('Error logging drink on server:', resp.error)
      showToast(`⚠️ ${resp.error.error || 'Failed to log drink on server'}`)
      const pResp = await getParty({ path: { code: code.value } })
      if (pResp?.data) {
        party.value = pResp.data
        nextTick(() => renderCharts())
      }
      return
    }

    showToast(`Logged ${drinkEmoji(drinkName)} ${drinkName} for ${attendee}! 🍻`)

    // Re-fetch party from server to confirm persisted state
    const pResp = await getParty({ path: { code: code.value } })
    if (pResp?.data) {
      party.value = pResp.data
      nextTick(() => renderCharts())
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
    // Record as custom drink for this party
    if (!localPartyCustomDrinkNames.value.includes(name)) {
      localPartyCustomDrinkNames.value.push(name)
      try {
        localStorage.setItem(
          `gamevote_party_${code.value}_custom_drinks`,
          JSON.stringify(localPartyCustomDrinkNames.value)
        )
      } catch {}
    }
    showCustomPartyDrinks.value = true

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
const bacChartEl = ref<HTMLCanvasElement | null>(null)

let barChart: Chart | null = null
let lineChart: Chart | null = null
let bacChart: Chart | null = null

const palette = [
  '#6366f1', '#8b5cf6', '#22c55e', '#f59e0b', '#ec4899', '#06b6d4', '#10b981', '#3b82f6', '#f43f5e'
]

function renderCharts() {
  const logEntries = rawBeerLog.value
  const leaderboard = attendeeLeaderboard.value

  // 1. Bar Chart (Leaderboard)
  if (barChartEl.value && (activeChartTab.value === 'leaderboard' || activeChartTab.value === 'all')) {
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
            ticks: { color: '#94a3b8', font: { size: 11 } },
            grid: { color: '#1e2739' }
          }
        }
      }
    })
  }

  // 2. Line Chart (Consumption Over Time with REAL Timestamps)
  if (lineChartEl.value && (activeChartTab.value === 'timeline' || activeChartTab.value === 'all')) {
    if (lineChart) lineChart.destroy()

    if (logEntries.length > 0) {
      // Sort chronologically
      const sorted = [...logEntries].sort(
        (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      )

      const allAttendees = attendeeList.value
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
                boxWidth: 14,
                font: { size: 12, weight: 'bold' }
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
                font: { size: 11 }
              },
              grid: { color: '#1e2739' }
            },
            y: {
              beginAtZero: true,
              ticks: { color: '#94a3b8', font: { size: 11 } },
              grid: { color: '#1e2739' }
            }
          }
        }
      })
    }
  }

  // 3. BAC Verlauf Line Chart
  if (bacChartEl.value && (activeChartTab.value === 'bac' || activeChartTab.value === 'all')) {
    if (bacChart) bacChart.destroy()

    const bacData = selectedTargetBAC.value
    if (bacData.points.length > 0) {
      const labels = bacData.points.map(p => p.timeLabel)
      const values = bacData.points.map(p => {
        return bacUnit.value === 'promille' ? p.promille : p.bac
      })

      // Threshold line for driving limit: 0.5‰ or 0.05%
      const limitVal = bacUnit.value === 'promille' ? 0.5 : 0.05
      const limitLine = labels.map(() => limitVal)

      // Segment styling for past vs projected future
      const nowTime = Date.now()
      const segmentBorderDash = (ctx: any) => {
        const pt = bacData.points[ctx.p0DataIndex]
        return pt && pt.time.getTime() > nowTime ? [6, 6] : undefined
      }

      bacChart = new Chart(bacChartEl.value, {
        type: 'line',
        data: {
          labels,
          datasets: [
            {
              label: `${bacTargetAttendee.value || 'Attendee'} Blood Alcohol (${bacUnit.value === 'promille' ? '‰' : '%'})`,
              data: values,
              borderColor: '#f43f5e',
              backgroundColor: 'rgba(244, 63, 94, 0.12)',
              borderWidth: 3,
              fill: true,
              tension: 0.4,
              pointRadius: (ctx) => {
                const pt = bacData.points[ctx.dataIndex]
                return pt && pt.time.getTime() > nowTime ? 1 : 2
              },
              pointHoverRadius: 6,
              segment: {
                borderDash: segmentBorderDash,
                borderColor: (ctx: any) => {
                  const pt = bacData.points[ctx.p0DataIndex]
                  return pt && pt.time.getTime() > nowTime ? '#a855f7' : '#f43f5e'
                }
              }
            },
            {
              label: `Driving Limit (${limitVal}${bacUnit.value === 'promille' ? '‰' : '%'})`,
              data: limitLine,
              borderColor: 'rgba(239, 68, 68, 0.45)',
              borderWidth: 1.5,
              borderDash: [4, 4],
              pointRadius: 0,
              fill: false,
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          interaction: {
            mode: 'index',
            intersect: false
          },
          plugins: {
            legend: {
              position: 'top',
              labels: {
                color: '#e2e8f0',
                boxWidth: 14,
                font: { size: 12, weight: 'bold' }
              }
            },
            tooltip: {
              callbacks: {
                title: (items) => {
                  const idx = items[0]?.dataIndex ?? 0
                  const pt = bacData.points[idx]
                  const futureBadge = pt?.isProjected ? ' (Projected Recovery)' : ''
                  return `Time: ${pt?.timeLabel ?? ''}${futureBadge}`
                },
                label: (ctx) => {
                  if (ctx.datasetIndex === 1) {
                    return ` ${ctx.dataset.label}`
                  }
                  const idx = ctx.dataIndex
                  const pt = bacData.points[idx]
                  if (!pt) return ''
                  return ` BAC: ${pt.bac.toFixed(3)}% (${pt.promille.toFixed(2)}‰) - ${getSobrietyStage(pt.bac).label}`
                }
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
              ticks: {
                color: '#94a3b8',
                font: { size: 11 },
                callback: (val) => `${val}${bacUnit.value === 'promille' ? '‰' : '%'}`
              },
              grid: { color: '#1e2739' }
            }
          }
        }
      })
    }
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

// Watch changes to re-render charts
watch(
  [() => party.value?.beerLog, bacTargetAttendee, bacWeightKg, bacGender, bacUnit],
  () => {
    nextTick(() => renderCharts())
  },
  { deep: true }
)

onMounted(async () => {
  // Load saved profile if available
  try {
    const saved = localStorage.getItem('gamevote_bac_profile')
    if (saved) {
      const parsed = JSON.parse(saved)
      if (parsed.weightKg) bacWeightKg.value = parsed.weightKg
      if (parsed.gender) bacGender.value = parsed.gender
      if (parsed.unit) bacUnit.value = parsed.unit
    }
  } catch {}

  // Load custom drinks created for this specific party
  try {
    const savedPartyCustom = localStorage.getItem(`gamevote_party_${code.value}_custom_drinks`)
    if (savedPartyCustom) {
      localPartyCustomDrinkNames.value = JSON.parse(savedPartyCustom)
    }
  } catch {}

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

  // Auto-join party and ensure currentUser is targeted for charts
  if (currentUser.value) {
    bacTargetAttendee.value = currentUser.value
    await ensureJoinedParty()
  }

  await nextTick()
  renderCharts()
})

// Save profile settings
watch([bacWeightKg, bacGender, bacUnit], () => {
  try {
    localStorage.setItem(
      'gamevote_bac_profile',
      JSON.stringify({
        weightKg: bacWeightKg.value,
        gender: bacGender.value,
        unit: bacUnit.value
      })
    )
  } catch {}
})

onUnmounted(() => {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  if (barChart) barChart.destroy()
  if (lineChart) lineChart.destroy()
  if (bacChart) bacChart.destroy()
  if (toastTimer) clearTimeout(toastTimer)
})
</script>

<style scoped>
.stat-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: 1rem;
  padding: 1.15rem;
  transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
}
.stat-card:hover {
  transform: translateY(-2px);
  border-color: rgba(99,102,241,0.4);
  box-shadow: 0 4px 20px rgba(99,102,241,0.1);
}

.chart-tab-btn {
  padding: 0.5rem 1rem;
  border-radius: 0.85rem;
  font-weight: 600;
  color: var(--c-muted);
  white-space: nowrap;
  transition: all 0.2s;
  cursor: pointer;
}
.chart-tab-btn:hover {
  color: #fff;
}
.chart-tab-btn-active {
  background: var(--c-accent);
  color: #fff !important;
  box-shadow: 0 2px 10px rgba(99,102,241,0.35);
}

.chart-container {
  min-height: 380px;
  height: 380px;
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

.guide-expand-enter-active,
.guide-expand-leave-active {
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 1200px;
  overflow: hidden;
}
.guide-expand-enter-from,
.guide-expand-leave-to {
  opacity: 0;
  max-height: 0;
  transform: translateY(-8px);
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
