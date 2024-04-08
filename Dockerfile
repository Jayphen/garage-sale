FROM golang:1.22.2-alpine as builder

RUN apk add -v build-base
RUN apk add -v go 
RUN apk add -v ca-certificates
RUN apk add --no-cache \
    unzip \
    # this is needed only if you want to use scp to copy later your pb_data locally
    openssh



COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /go/bin/app .


FROM alpine:latest
WORKDIR /app

COPY --from=builder /go/bin/app .
COPY ./assets ./assets
EXPOSE 8080

# start PocketBase
CMD ["/app/app", "serve", "--http=0.0.0.0:8080"]
