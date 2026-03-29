import {defineStore} from 'pinia'
import {ref} from 'vue'
import Cookies from 'js-cookie'
import {validateUser, login as apiLogin} from '@/api-client'
import {resolve} from "chart.js/helpers";

export const useAuthStore = defineStore('auth', () => {
    const username = ref<string | undefined>(Cookies.get('username'))
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
            username.value = undefined
            return false
        }

        isValidating.value = true

        const response = await validateUser({
            path: {
                username: cookieUsername
            }
        })

        if(response.error) {
            // User doesn't exist in database, clear invalid cookie
            console.warn('Session invalid: user not found in database, clearing cookie')
            username.value = undefined
            Cookies.remove('username')
            console.log("asdf4")
            isValidating.value = false
            return false
        }else{
            username.value = cookieUsername
            isValidating.value = false
            return true
        }
    }

    function logout() {
        username.value = undefined
        Cookies.remove('username')
    }

    return {username, isValidating, login, logout, validateSession}
})
