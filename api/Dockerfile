FROM golang:1.23.6-alpine AS build

WORKDIR /app
COPY . .
RUN cd api &&\
    CGO_ENABLED=0 go build -o bin/habapi ./cmd/habapi/main.go

FROM golang:1.23.6-alpine

WORKDIR /app
USER 1000:1000
COPY --from=build --chown=1000:1000 /app/api/bin/habapi ./habapi

CMD ["/app/habapi"]