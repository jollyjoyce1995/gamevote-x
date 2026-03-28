<template>
  <section>
    <div class="card mb-4">
      <h2 class="text-lg font-bold mb-4">🎮 Nominate Games</h2>

      <!-- Steam search -->
      <div class="flex gap-4 mb-4">
        <input
          v-model="gameSearch"
          name="gameSearch"
          class="input py-0"
          placeholder="Search Steam games..."
          @input="onSearchInput"
        />
        <button class="btn btn-primary" @click="submitCustomGame">Add Custom</button>
      </div>

      <!-- Steam suggestions -->
      <div v-if="suggestions.length" class="rounded-xl overflow-hidden flex flex-col gap-2">
        <GameItem
          v-for="g in suggestions"
          :key="g.appId"
          mode="suggestion"
          :game="g"
          @nominate="pickGame(g)"
        />
      </div>

      <!-- Nominated list -->
      <div class="mt-2">
        <p class="text-sm font-semibold mb-3" style="color:var(--c-muted)">Nominated Games ({{ party?.options?.length || 0 }})</p>
        <TransitionGroup name="list" tag="div" class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <GameItem
            v-for="opt in party?.options"
            :key="opt.name"
            mode="nominated"
            :game="opt"
            @remove="removeOption"
          />
        </TransitionGroup>
        <p v-if="!party?.options?.length" class="text-sm text-center py-8" style="color:var(--c-muted)">No games nominated yet. Use the search bar above to add some!</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { usePartyStore } from '@/stores/party'
import { searchGames, postOption, deleteOption } from '@/api-client'
import GameItem from '@/components/party/GameItem.vue'
import type {ModelsGame} from "@/generated-api";

const route = useRoute()
const partyStore = usePartyStore()

const code = computed(() => route.params.code as string)
const party = computed(() => partyStore.party)

const gameSearch = ref('')
const suggestions = ref<ModelsGame[]>([])
let searchTimeout: ReturnType<typeof setTimeout> | null = null

function onSearchInput() {
  if (searchTimeout) clearTimeout(searchTimeout)
  if (!gameSearch.value.trim()) { suggestions.value = []; return }
  searchTimeout = setTimeout(async () => {
    const res = await searchGames({ query: { q: gameSearch.value } }).catch(() => null)
    suggestions.value = res?.data || []
  }, 300)
}

async function pickGame(item: ModelsGame | string) {
  const name = typeof item === 'string' ? item.trim() : item.name
  if (!name) return

  const exists = partyStore.party!.options?.some((opt: any) => opt.name?.toLowerCase() === name.toLowerCase())
  if (exists) {
    alert('This game is already nominated!')
    return
  }

  const option = typeof item === 'string' 
    ? { name } 
    : { name: item.name, appId: item.appId, imageUrl: item.imageUrl }

  gameSearch.value = ''
  suggestions.value = []
  
  try {
    await postOption({
      path: { code: code.value },
      body: option
    })
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
  if (!name) return

  await deleteOption({
      path: { code: code.value, gameName: name }
  }).catch((e: any) => {
    console.error('Failed to remove option:', e)
  })
}
</script>
