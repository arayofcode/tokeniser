FROM golang:1.22-alpine as base

WORKDIR /build

RUN apk add make golangci-lint

COPY go.* .
COPY Makefile .
COPY main.go .
COPY src/ ./src/
RUN make install

FROM base as build
RUN make build

FROM base as development
CMD [ "go", "run", "/build/main.go" ]

FROM gcr.io/distroless/static-debian12 as production

WORKDIR /app
COPY --from=build /build/bin/tokeniser-linux ./tokeniser
COPY --from=build /build/src/router/templates ./src/router/templates
COPY --from=build /build/src/router/static ./src/router/static

CMD [ "/app/tokeniser" ]