# build stage
FROM golang:1.13.1-stretch AS build-env
ENV APP_DIR=$GOPATH/src/github.com/mono83/tame
ADD . $APP_DIR
RUN cd $APP_DIR && make build && mv $APP_DIR/release/tame /goapp

# final stage
FROM alpine:3.10
WORKDIR /app
COPY --from=build-env /goapp /app/tame
ENTRYPOINT ["./tame"]
CMD []