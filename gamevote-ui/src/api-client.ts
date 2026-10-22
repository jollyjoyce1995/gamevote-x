import { client } from './generated-api/client.gen'

// The Go server serves the UI and the API from the same origin
// (in dev, vite proxies /api to it).
client.setConfig({ baseUrl: '/api' })

export * from './generated-api'
export { client }
