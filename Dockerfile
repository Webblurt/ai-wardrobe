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