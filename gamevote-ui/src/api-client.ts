import { Configuration, PartiesApi, PollsApi } from './generated-api'

const config = new Configuration({
    basePath: 'http://localhost:8080',
    // You can add middleware here for auth if needed
})

export const partiesApi = new PartiesApi(config)
export const pollsApi = new PollsApi(config)
