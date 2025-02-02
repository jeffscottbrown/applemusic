FROM golang:1.23-alpine AS appbuilder

ARG GIT_COMMIT

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-X github.com/jeffscottbrown/applemusic/commit.Hash=$GIT_COMMIT -X github.com/jeffscottbrown/applemusic/commit.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -o musicsearch .

FROM gcr.io/distroless/static-debian12

ARG PROJECT_ID
ENV PROJECT_ID=$PROJECT_ID

WORKDIR /app
COPY --from=appbuilder /build/musicsearch .
COPY --from=appbuilder /build/web/templates/ ./web/templates/
COPY --from=appbuilder /build/web/assets/ ./web/assets/
CMD ["./musicsearch"]
