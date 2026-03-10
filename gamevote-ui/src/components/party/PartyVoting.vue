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
import {pollsApi} from '@/api-client'
import GameItem from "@/components/party/GameItem.vue";

const authStore = useAuthStore()
const partyStore = usePartyStore()

const party = computed(() => partyStore.party)

const currentVotes = ref<Record<string, number>>({})
const submittingVote = ref(false)
const alreadyVoted = ref(false)
const outstanding = ref<string[]>([])

function setVote(game: string, val: number) {
  currentVotes.value[game] = currentVotes.value[game] === val ? 0 : val
}

async function submitVotes() {
  if (!authStore.username) return
  submittingVote.value = true
  try {
    const pollId = party.value?.links?.['poll']?.href?.split('/').pop()
    if (!pollId) return
    await pollsApi.pollsIdVotesAttendeePut({
      id: pollId,
      attendee: authStore.username,
      choices: currentVotes.value
    })
    alreadyVoted.value = true
    await loadOutstanding()
  } finally {
    submittingVote.value = false
  }
}

async function loadOutstanding() {
  const pollId = party.value?.links?.['poll']?.href?.split('/').pop()
  if (!pollId) return
  outstanding.value = await pollsApi.pollsIdOutstandingGet({id: pollId}).catch(() => [])
  if (!outstanding.value.includes(authStore.username!)) {
    alreadyVoted.value = true
  }
}

function vote(name: string, like: boolean) {
  console.log('vote', name, like)
}

onMounted(() => {
  if (party.value?.status === 'VOTING') {
    loadOutstanding()
  }
})

watch(() => party.value?.status, (newStatus) => {
  if (newStatus === 'VOTING') {
    alreadyVoted.value = false
    currentVotes.value = {}
    loadOutstanding()
  }
})
</script>
