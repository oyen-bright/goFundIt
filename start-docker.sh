#!/bin/bash

# Ensure the script stops on errors
set -e

# Run docker-compose with the specified env file
echo "Starting Docker Compose..."
docker-compose --env-file config/env/.env.dev up
