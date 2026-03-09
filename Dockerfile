# BACKEND BUILD 
FROM golang:1.26 AS backend-build
WORKDIR /app
COPY apps/backend ./apps/backend
WORKDIR /app/apps/backend
RUN go mod tidy && go build -o /app/bin/ai-wardrobe ./cmd/ai-wardrobe

FROM debian:bookworm-slim AS backend
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=backend-build /app/bin/ai-wardrobe ./ai-wardrobe
COPY apps/backend/config ./config
COPY .env .env
EXPOSE 8002
CMD ["./ai-wardrobe"]

# FRONTEND BUILD
FROM node:20 AS frontend-build
WORKDIR /app
COPY apps/frontend ./apps/frontend
WORKDIR /app/apps/frontend
RUN npm install
RUN npm run build

# FRONTEND SERVE
FROM nginx:stable-alpine AS frontend
COPY --from=frontend-build /app/apps/frontend/dist /usr/share/nginx/html
EXPOSE 3002
CMD ["nginx", "-g", "daemon off;"]