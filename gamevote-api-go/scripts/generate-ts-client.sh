#!/bin/bash
# Install OpenAPI Generator CLI if not available
# npm install @openapitools/openapi-generator-cli -g
swag init -g .\cmd\api\main.go
bun x @openapitools/openapi-generator-cli generate \
    -i docs/swagger.json \
    -g typescript-fetch \
    -o client-ts \
    --additional-properties=typescriptThreePlus=true
