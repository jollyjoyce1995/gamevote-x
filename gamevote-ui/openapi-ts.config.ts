import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../gamevote-api-go/docs/swagger.json',
  output: 'src/generated-api',
  plugins: ['@hey-api/client-fetch'],
});
