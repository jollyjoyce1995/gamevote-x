#!/usr/bin/env bash

# bun install -g @hey-api/openapi-ts
swag init -g .\cmd\api\main.go
cd ../gamevote-ui && npx openapi-ts
