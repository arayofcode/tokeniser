FROM golang:1.22-alpine as build

WORKDIR /build
COPY . .
RUN apk add --no-cache make=4.4
RUN make dep
RUN make build

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=build /build/bin/tokeniser .
COPY --from=build /build/router/templates router/templates
COPY --from=build /build/router/static router/static

CMD [ "/app/tokeniser" ]