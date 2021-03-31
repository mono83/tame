# build stage
FROM golang:1.13.1-stretch AS build-env
WORKDIR /tame
ADD . /tame
RUN mkdir -p release && go test ./... && CGO_ENABLED=0 go build -o release/tame app/tame.go

# final stage
FROM alpine:3.10
WORKDIR /app
COPY --from=build-env /tame/release/tame /app/tame
ENTRYPOINT ["./tame"]
