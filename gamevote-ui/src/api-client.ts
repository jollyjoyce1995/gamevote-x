import { client } from './generated-api/client.gen'

client.setConfig({
    baseUrl: 'http://localhost:8080',
})

export * from './generated-api'
export { client }
