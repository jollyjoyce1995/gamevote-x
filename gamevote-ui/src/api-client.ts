import { Configuration, PartiesApi } from './generated-api'

const config = new Configuration({
    basePath: 'http://localhost:8080',
    // You can add middleware here for auth if needed
})

export const partiesApi = new PartiesApi(config)
