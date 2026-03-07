#!/bin/bash
# Install OpenAPI Generator CLI if not available
# npm install @openapitools/openapi-generator-cli -g

npx @openapitools/openapi-generator-cli generate \
    -i docs/swagger.json \
    -g typescript-fetch \
    -o client-ts \
    --additional-properties=typescriptThreePlus=true
