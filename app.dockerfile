FROM golang:1.23-alpine as app-builder
RUN apk update && apk add curl make git

WORKDIR /src
COPY go.mod .
COPY go.sum .
COPY api/ api/

RUN go env -w GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -gcflags="all=-N -l" -o app ./cmd/app

FROM alpine:latest
RUN apk update && apk add --no-cache curl
WORKDIR /src
COPY --from=app-builder /src/app .
COPY --from=app-builder /src/api .

CMD ["./app"]
