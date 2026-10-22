import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Party } from '@/generated-api'

export const usePartyStore = defineStore('party', () => {
    const party = ref<Party | null>(null)
    const onlineUsers = ref<string[]>([])
    let eventSource: EventSource | null = null

    function setParty(p: Party) {
        party.value = p
    }

    function setOnlineUsers(users: string[]) {
        onlineUsers.value = users
    }

    function setEventSource(es: EventSource | null) {
        eventSource = es
    }

    function closeStream() {
        if (eventSource) {
            eventSource.close()
            eventSource = null
        }
        party.value = null
        onlineUsers.value = []
    }

    return { party, onlineUsers, setParty, setOnlineUsers, setEventSource, closeStream }
})
