<template>
  <main class="min-h-screen flex items-center justify-center px-6 py-10 bg-slate-950/20">
    <div class="card w-full max-w-xl shadow-2xl border border-indigo-500/20">
      <h1 class="text-3xl font-extrabold mb-1 grad-text">Create a Party</h1>
      <p class="text-sm mb-8" style="color:var(--c-muted)">Select the attendees for tonight's session</p>

      <!-- Attendees -->
      <div class="mb-8">
        <div class="flex items-center justify-between mb-4">
          <label class="text-sm font-bold uppercase tracking-wider text-indigo-400">Available Users</label>
          <span class="text-xs text-muted">{{ attendees.length }} selected</span>
        </div>
        
        <div class="grid grid-cols-1 gap-2 max-h-64 overflow-y-auto pr-2 custom-scrollbar">
          <div
            v-for="u in allUsers"
            :key="u.username"
            @click="toggleAttendee(u.username || '')"
            class="flex items-center justify-between px-4 py-3 rounded-xl cursor-pointer transition-all duration-200 group"
            :class="attendees.includes(u.username || '') 
              ? 'bg-indigo-500/10 border-indigo-500/40 border' 
              : 'bg-white/5 border-transparent border hover:bg-white/10 hover:border-white/10'"
          >
            <div class="flex items-center gap-3">
              <div 
                class="w-10 h-10 rounded-full flex items-center justify-center text-lg shadow-inner"
                :class="attendees.includes(u.username || '') ? 'bg-indigo-500 text-white' : 'bg-slate-800 text-slate-400'"
              >
                {{ u.username?.charAt(0).toUpperCase() }}
              </div>
              <span class="font-medium" :class="attendees.includes(u.username || '') ? 'text-indigo-200' : 'text-slate-300'">
                {{ u.username }}
              </span>
            </div>
            <div 
              class="w-6 h-6 rounded-full border-2 flex items-center justify-center transition-colors"
              :class="attendees.includes(u.username || '') 
                ? 'bg-indigo-500 border-indigo-500' 
                : 'border-slate-700 group-hover:border-slate-500'"
            >
              <span v-if="attendees.includes(u.username || '')" class="text-white text-xs">✓</span>
            </div>
          </div>
        </div>
      </div>

      <p v-if="error" class="text-sm mb-6 p-3 rounded-lg bg-red-500/10 border border-red-500/20" style="color:var(--c-danger)">
        ⚠️ {{ error }}
      </p>

      <div class="flex gap-4">
        <RouterLink to="/" class="btn btn-ghost flex-1 justify-center py-3">Cancel</RouterLink>
        <button
          class="btn btn-primary flex-1 py-3 text-lg font-bold shadow-lg shadow-indigo-500/20"
          :disabled="loading || attendees.length === 0"
          @click="handleCreate"
        >
          <span v-if="loading">Creating...</span>
          <span v-else>🚀 Start Party</span>
        </button>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { createParty, postAttendee, getUsers } from '@/api-client'
import type { ServiceUserDto as UserDTO } from '@/generated-api'

const router = useRouter()
const allUsers = ref<UserDTO[]>([])
const attendees = ref<string[]>([])
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    const usersResp = await getUsers()
    allUsers.value = usersResp.data || []
    // Preselect everyone
    attendees.value = allUsers.value.map(u => u.username || '')
  } catch (e) {
    console.error('Failed to fetch users:', e)
    error.value = 'Could not load users. Please check if the backend is running.'
  }
})

function toggleAttendee(username: string) {
  const idx = attendees.value.indexOf(username)
  if (idx > -1) {
    attendees.value.splice(idx, 1)
  } else {
    attendees.value.push(username)
  }
}

async function handleCreate() {
  if (!attendees.value.length) return
  loading.value = true
  error.value = ''
  try {
    const partyResp = await createParty({ body: {} })
    const partyCode = partyResp.data?.code!
    // Add each attendee via the API
    for (const a of attendees.value) {
      await postAttendee({ path: { code: partyCode }, body: { value: a } }).catch(() => {})
    }
    router.push(`/parties/${partyCode}`)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to create party'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(99, 102, 241, 0.2);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(99, 102, 241, 0.4);
}
</style>
