import {defineStore} from 'pinia'
import {ref} from 'vue'
import Cookies from 'js-cookie'
import {validateUser, login as apiLogin} from '@/api-client'

export const useAuthStore = defineStore('auth', () => {
    const username = ref<string | null>(localStorage.getItem('username'))
    const isValidating = ref(false)

    async function login(name: string) {
        username.value = name
        Cookies.set('username', name)
        try {
            await apiLogin({body: {username: name}})
        } catch (e) {
            console.error(e)
        }
    }

    async function validateSession(): Promise<boolean> {
        const cookieUsername = Cookies.get('username')
        if (!cookieUsername) {
            username.value = null
            return false
        }

        // If we already have the username and it matches cookie, assume it's valid
        // This prevents unnecessary API calls on every route change
        if (username.value === cookieUsername) {
            return true
        }

        isValidating.value = true
        try {
            await validateUser({
                path: {
                    username: cookieUsername
                }
            })
            username.value = cookieUsername
            return true
        } catch (error) {
            // User doesn't exist in database, clear invalid cookie
            console.warn('Session invalid: user not found in database, clearing cookie')
            username.value = null
            Cookies.remove('username')
            return false
        } finally {
            isValidating.value = false
        }
    }

    function logout() {
        username.value = null
        Cookies.remove('username')
    }

    return {username, isValidating, login, logout, validateSession}
})
