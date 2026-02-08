#!/bin/bash

trap 'exit 1' ERR

cd "$(dirname "$0")"
git stash
git pull

cd backend
docker build -t dropmail-backend:latest .
cd ../frontend
docker build -t dropmail-frontend:latest .
cd ..

docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d

docker system prune -af

trap - ERR