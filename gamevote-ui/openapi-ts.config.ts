import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../gamevote-api-go/api/openapi.yaml',
  output: 'src/generated-api',
  plugins: ['@hey-api/client-fetch'],
});
