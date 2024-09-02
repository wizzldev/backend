FROM golang:1.21-alpine

WORKDIR /app

# Install WebP image dependencies for go-webp
# https://github.com/kolesa-team/go-webp
RUN apk add --no-cache gcc musl-dev linux-headers libwebp-dev
RUN rm -rf var/cache/*

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .

ENV CGO_ENABLED = 1
RUN go build -o wizzl
CMD ["./wizzl", "-"]