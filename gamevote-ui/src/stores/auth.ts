import { defineStore } from 'pinia'
import { ref } from 'vue'
import { loginUser } from '@/api'
import Cookies from 'js-cookie'

export const useAuthStore = defineStore('auth', () => {
    const username = ref<string | null>(Cookies.get('username') || null)

    async function login(name: string) {
        await loginUser(name)
        username.value = name
        Cookies.set('username', name, { expires: 30 })
    }

    function logout() {
        username.value = null
        Cookies.remove('username')
    }

    return { username, login, logout }
})
