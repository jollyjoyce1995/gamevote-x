import {defineStore} from 'pinia'
import {ref} from 'vue'
import Cookies from 'js-cookie'
import {login as apiLogin} from '@/api-client'

export const useAuthStore = defineStore('auth', () => {
    const username = ref<string | undefined>(Cookies.get('username'))
    const isValidating = ref(false)

    async function login(name: string) {
        const response = await apiLogin({body: {username: name}})
        if (response.error) {
            throw new Error(response.error.error)
        }
        username.value = response.data.username
        Cookies.set('username', response.data.username, {expires: 365})
    }

    async function validateSession(): Promise<boolean> {
        const cookieUsername = Cookies.get('username')
        if (!cookieUsername) {
            username.value = undefined
            return false
        }

        isValidating.value = true

        // Logging in is idempotent: it re-registers the name if the database
        // was reset and refreshes the last-login time.
        const response = await apiLogin({body: {username: cookieUsername}})
        isValidating.value = false

        if (response.error) {
            console.warn('Session invalid, clearing cookie')
            username.value = undefined
            Cookies.remove('username')
            return false
        }
        username.value = response.data.username
        return true
    }

    function logout() {
        username.value = undefined
        Cookies.remove('username')
    }

    return {username, isValidating, login, logout, validateSession}
})
