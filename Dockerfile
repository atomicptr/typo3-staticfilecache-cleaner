FROM alpine:3.11 as base

RUN apk add --no-cache build-base go git

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"'

FROM scratch

ENV CLEAN_PATH "/data"

COPY --from=base /app /app
WORKDIR /app

ENTRYPOINT ["/app/typo3-staticfilecache-cleaner"]