<template>
  <RouterLink
    :to="`/parties/${party.code}`"
    class="card group cursor-pointer no-underline"
    style="text-decoration:none"
  >
    <div class="flex items-start justify-between mb-4">
      <div>
        <div class="text-xs font-mono mb-1" style="color:var(--c-muted)">{{ party.code }}</div>
        <div class="flex items-center gap-2">
          <span :class="statusBadge" class="badge">{{ party.status }}</span>
        </div>
      </div>
      <span class="text-2xl">🎮</span>
    </div>
    <div class="flex gap-4 text-sm" style="color:var(--c-muted)">
      <span>👥 {{ party.attendees?.length || 0 }}</span>
      <span>📋 {{ party.options?.length || 0 }} games</span>
      <span>🍺 {{ party.beerCount || 0 }}</span>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import type { PartyDTO } from '@/api'

const props = defineProps<{
  party: PartyDTO
}>()

const statusBadge = computed(() => ({
  'badge-nomination': props.party.status === 'NOMINATION',
  'badge-voting': props.party.status === 'VOTING',
  'badge-results': props.party.status === 'RESULTS',
}))
</script>
