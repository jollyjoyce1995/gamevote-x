<template>
  <div id="app">
    <!-- Google Font -->
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet" />

    <!-- Global Nav (only when logged in and in a party context) -->
    <AppNavbar v-if="authStore.username" />

    <!-- Page transitions -->
    <RouterView v-slot="{ Component }">
      <Transition name="page" mode="out-in">
        <component :is="Component" :key="$route.fullPath" />
      </Transition>
    </RouterView>
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppNavbar from '@/components/AppNavbar.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<style>
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
</style>
