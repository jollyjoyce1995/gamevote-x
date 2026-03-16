<template>
  <section>
    <div v-if="alreadyVoted" class="card text-center py-8">
      <p class="text-lg mb-2">✅ You've voted!</p>
      <p class="text-sm" style="color:var(--c-muted)">Waiting for others: {{
          outstanding.join(', ') || 'All done!'
        }}</p>
    </div>
    <div v-else class="card">
      <h2 class="text-lg font-bold mb-6">👍👎 Vote on Games</h2>
      <div class="flex flex-col gap-4">
        <GameItem
            v-for="opt in party?.options"
            mode="voting"
            :key="opt.name"
            :game="opt"
            :like="currentVotes[opt.name!] ?? 0"
            @like="vote"
            class="flex items-center gap-3 p-3 rounded-xl"
        />
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
</template>

<script setup lang="ts">
import {ref, computed, onMounted, watch} from 'vue'
import {useAuthStore} from '@/stores/auth'
import {usePartyStore} from '@/stores/party'
import {partiesApi} from '@/api-client'
import GameItem from "@/components/party/GameItem.vue";

const authStore = useAuthStore()
const partyStore = usePartyStore()

const party = computed(() => partyStore.party)

const currentVotes = ref<Record<string, number>>({})
const submittingVote = ref(false)

const outstanding = computed(() => {
  return party.value?.outstandingVoters?.filter(x => x != authStore.username) || []
})
const alreadyVoted = computed(() => {
  if (party.value?.status !== 'VOTING') return false
  if (!authStore.username) return false
  return !party.value?.outstandingVoters?.includes(authStore.username)
})

function vote(name: string, like: number) {
  currentVotes.value[name] = like;
}

async function submitVotes() {
  if (!authStore.username) return
  submittingVote.value = true
  if(!party.value?.code){
    return
  }
  try {
    await partiesApi.postVote({
      code: party.value.code,
      attendee: authStore.username,
      choices: currentVotes.value,
    })
    // No need to loadOutstanding manually. The SSE event `party_updated` 
    // will push the updated outstandingVoters list instantly.
  } finally {
    submittingVote.value = false
  }
}

watch(() => party.value?.status, (newStatus) => {
  if (newStatus === 'VOTING') {
    currentVotes.value = {}
  }
})
</script>
