<template>
  <main class="min-h-screen flex items-center justify-center p-4">
    <div class="card w-full max-w-sm text-center">
      <div class="mb-8">
        <div class="text-5xl mb-3">🎲</div>
        <h1 class="text-3xl font-extrabold grad-text mb-1">GameVote</h1>
        <p class="text-sm" style="color:var(--c-muted)">Vote on what to play tonight</p>
      </div>

      <form @submit.prevent="handleLogin" class="flex flex-col gap-4">
        <div class="text-left">
          <label class="block text-sm font-medium mb-1.5">Your Name</label>
          <input
            id="username-input"
            v-model="username"
            class="input"
            placeholder="Enter your username..."
            autofocus
            required
          />
        </div>
        <button id="login-btn" type="submit" class="btn btn-primary" :disabled="loading || !username.trim()">
          <span v-if="loading">Joining...</span>
          <span v-else>Join 🚀</span>
        </button>
        <p v-if="error" class="text-sm" style="color:var(--c-danger)">{{ error }}</p>
      </form>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const username = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  if (!username.value.trim()) return
  loading.value = true
  error.value = ''
  try {
    await authStore.login(username.value.trim())
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>
