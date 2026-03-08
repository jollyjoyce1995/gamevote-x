const BASE_URL = 'http://localhost:8080'

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const res = await fetch(`${BASE_URL}${path}`, {
        method,
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: body ? JSON.stringify(body) : undefined,
    })
    if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }))
        throw new Error(err.error || `HTTP ${res.status}`)
    }
    if (res.status === 204 || res.headers.get('content-length') === '0') return undefined as T
    return res.json()
}

// ── Users ──────────────────────────────────────────────────────────────
export const loginUser = (username: string) =>
    request<{ id: string; username: string }>('POST', '/users', { username })
export const getUsers = () =>
    request<{ id: string; username: string }[]>('GET', '/users')

// ── Parties ────────────────────────────────────────────────────────────
export const getParties = () => request<PartyDTO[]>('GET', '/parties')
export const getParty = (code: string) => request<PartyDTO>('GET', `/parties/${code}`)
export const createParty = (attendees: string[]) =>
    request<PartyDTO>('POST', '/parties', { attendees, options: [] })
export const patchParty = (code: string, status: string) =>
    request<PartyDTO>('PATCH', `/parties/${code}`, { status })

// ── Options ────────────────────────────────────────────────────────────
export const addOption = (code: string, option: PartyOption) =>
    request<PartyOption>('POST', `/parties/${code}/options`, option)
export const deleteOption = (code: string, name: string) =>
    request<void>('DELETE', `/parties/${code}/options/${encodeURIComponent(name)}`)

// ── Attendees ──────────────────────────────────────────────────────────
export const addAttendee = (code: string, value: string) =>
    request<{ value: string }>('POST', `/parties/${code}/attendees`, { value })

// ── Beers / Drinks ─────────────────────────────────────────────────────
export const postBeer = (code: string, attendee: string) =>
    request<void>('POST', `/parties/${code}/beers`, { attendee })

// ── Drink Types ────────────────────────────────────────────────────────
export const getDrinkPresets = () => request<DrinkType[]>('GET', '/drinks/presets')
export const postDrinkPreset = (dt: Omit<DrinkType, 'id'>) =>
    request<DrinkType>('POST', '/drinks/presets', dt)

// ── Games ──────────────────────────────────────────────────────────────
export const searchGames = (q: string) => request<Game[]>('GET', `/games?q=${encodeURIComponent(q)}`)

// ── Polls ──────────────────────────────────────────────────────────────
export const getPoll = (id: string) => request<Poll>('GET', `/polls/${id}`)
export const getOutstanding = (id: string) => request<string[]>('GET', `/polls/${id}/outstanding`)
export const putVote = (pollId: string, attendee: string, choices: Record<string, number>) =>
    request<Record<string, number>>('PUT', `/polls/${pollId}/votes/${encodeURIComponent(attendee)}`, choices)
export const getResults = (pollId: string) => request<Record<string, number>>('GET', `/polls/${pollId}/results`)

// ── SSE ────────────────────────────────────────────────────────────────
export function openPartyStream(code: string, username: string): EventSource {
    return new EventSource(`${BASE_URL}/parties/${code}/stream?username=${encodeURIComponent(username)}`)
}

// ── Types ──────────────────────────────────────────────────────────────
export interface PartyOption {
    name: string
    appId?: number
    imageUrl?: string
}

export interface PartyDTO {
    id: string
    code: string
    attendees: string[]
    options: PartyOption[]
    status: 'NOMINATION' | 'VOTING' | 'RESULTS'
    results?: Record<string, number>
    beerCount: number
    beerPerAttendee: Record<string, number>
    _links?: Record<string, { href: string }>
}

export interface DrinkType {
    id?: string
    name: string
    volumeMl: number
    alcoholPercent: number
    beerEquivalent: number
    unitName: string
}

export interface Game {
    id?: string
    appId: number
    name: string
    imageUrl?: string
}

export interface Poll {
    id: string
    options: PartyOption[]
    attendees: string[]
    status: 'IN_PROGRESS' | 'COMPLETED'
}
