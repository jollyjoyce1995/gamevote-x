<template>
  <div
    class="flex items-center gap-3 rounded-2xl relative group overflow-hidden transition-all duration-200"
    :class="containerClass"
    :style="containerStyle"
    @click="handleClick"
  >
    <template v-if="mode === 'nominated'">
      <div class="w-16 h-8 flex items-center justify-center rounded bg-gray-800 text-[10px] relative z-10" style="color:var(--c-muted)">GAME</div>
      <span class="text-sm font-bold flex-1 truncate relative z-10">{{ gameName }}</span>
      <button 
        class="btn btn-ghost btn-sm text-error opacity-0 group-hover:opacity-100 transition-opacity relative z-10 p-2"
        title="Remove nomination"
        @click.stop="$emit('remove', gameName)"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
      </button>
    </template>
    
    <template v-else-if="mode === 'suggestion'">
      <img v-if="imageUrl" :src="imageUrl" class="w-16 h-8 object-cover rounded shadow-sm" alt="" />
      <div class="flex-1">
        <p class="text-sm font-semibold">{{ gameName }}</p>
        <p v-if="appId" class="text-xs" style="color:var(--c-muted)">Steam AppID {{ appId }}</p>
      </div>
      <button class="btn btn-primary btn-sm" @click.stop="$emit('nominate')">Nominate</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ModelsPartyOption } from '@/generated-api'

const props = defineProps<{
  mode: 'suggestion' | 'nominated'
  game: ModelsPartyOption | { name: string, appId?: number | string, imageUrl?: string }
}>()

const emit = defineEmits<{
  (e: 'nominate'): void
  (e: 'remove', name: string): void
}>()

const gameName = computed(() => props.game.name || '')
const appId = computed(() => 'appId' in props.game ? props.game.appId : undefined)
const imageUrl = computed(() => 'imageUrl' in props.game ? props.game.imageUrl : undefined)

const containerClass = computed(() => {
  if (props.mode === 'suggestion') return 'px-4 py-3 cursor-pointer hover:bg-[var(--c-border)]'
  return 'px-3 py-3'
})

const containerStyle = computed(() => {
  if (props.mode === 'nominated') return { background: 'var(--c-bg-card)', border: '1px solid var(--c-border)' }
  return {}
})

function handleClick() {
  if (props.mode === 'suggestion') {
    emit('nominate')
  }
}
</script>
