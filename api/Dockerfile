FROM golang:1.23.6-alpine AS build

RUN apk add make git

WORKDIR /app
COPY . .
RUN cd api &&\
    make bin/habapi

FROM golang:1.23.6-alpine

WORKDIR /app
USER 1000:1000
COPY --from=build --chown=1000:1000 /app/api/bin/habapi ./habapi

CMD ["/app/habapi"]