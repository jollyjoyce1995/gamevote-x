<template>
  <div
      class="flex items-center gap-4 px-2 rounded-2xl relative group overflow-hidden transition-all duration-200 cursor-pointer hover:bg-gray-700"
      @click="handleClick"
  >
    <img v-if="imageUrl" :src="imageUrl" class="w-1/4 object-cover rounded shadow-sm" alt=""/>
    <div class="flex-1">
      <p class="text-sm font-semibold">{{ gameName }}</p>
      <p v-if="appId" class="text-xs" style="color:var(--c-muted)">Steam AppID {{ appId }}</p>
    </div>
    <button v-if="mode == 'suggestion'" class="btn btn-primary btn-sm" @click.stop="$emit('nominate')">Nominate</button>
    <button v-if="mode == 'nominated'" class="btn btn-primary btn-sm" @click.stop="$emit('remove', gameName)">Remove
    </button>
    <div v-if="mode == 'voting'" class="flex gap-2">
      <button
          class="vote-btn like"
          @click="setVote(game, true)"
      >👍</button>
      <button
          class="vote-btn dislike"
          @click="setVote(game, false)"
      >👎</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import type {ModelsPartyOption} from '@/generated-api'

const props = defineProps<{
  mode: 'suggestion' | 'nominated' | 'voting',
  game: ModelsPartyOption
}>()

const emit = defineEmits<{
  (e: 'nominate'): void
  (e: 'remove', name: string): void,
  (e: 'like', name: string, like: boolean): void,
}>()

const mode = computed(() => props.mode)
const gameName = computed(() => props.game.name || '')
const appId = computed(() => 'appId' in props.game ? props.game.appId : undefined)
const imageUrl = computed(() => 'imageUrl' in props.game ? props.game.imageUrl : undefined)

function handleClick() {
  if (props.mode === 'suggestion') {
    emit('nominate')
  }
}

function setVote(game: ModelsPartyOption, like: boolean) {
  emit('like', game.name!, like)
}
</script>
