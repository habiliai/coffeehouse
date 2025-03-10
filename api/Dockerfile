FROM golang:1.24-alpine3.21 AS build

RUN apk add make

WORKDIR /app
COPY . .
RUN cd api &&\
    make bin/alice

FROM alpine:3.21

WORKDIR /app
USER 1000:1000
COPY --from=build --chown=1000:1000 /app/api/bin/alice ./alice

CMD ["/app/alice"]