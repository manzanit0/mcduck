FROM golang:1.22-alpine3.18 AS build

ENV GOWORK=off

ARG VERSION=dev
ARG RAILWAY_SERVICE_NAME

WORKDIR /workspace

COPY . ./
RUN --mount="type=cache,id=s/${RAILWAY_SERVICE_NAME}-/root/.cache/go-build,target=/root/.cache/go-build" go build -ldflags "-X main.version=${VERSION}" -o app ./cmd/${RAILWAY_SERVICE_NAME}

FROM alpine:3.18

COPY --from=build /workspace/app /app

CMD ["/app"]
