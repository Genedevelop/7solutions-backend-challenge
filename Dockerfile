# Stage 1: build a static binary
FROM golang:1.27-alpine AS build
WORKDIR /src
# Step 1: download modules first so this layer stays cached until go.mod or go.sum changes
COPY go.mod go.sum ./
RUN go mod download
# Step 2: copy the source and build without cgo so the binary runs on a distroless image
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/http

# Stage 2: run on distroless as a non-root user, no shell and no package manager in the final image
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api"]
