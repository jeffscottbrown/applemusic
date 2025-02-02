FROM golang:1.23 AS appbuilder

ARG GIT_COMMIT
ARG PROJECT_ID

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-X github.com/jeffscottbrown/applemusic/secrets.projectId=$PROJECT_ID -X github.com/jeffscottbrown/applemusic/commit.Hash=$GIT_COMMIT -X github.com/jeffscottbrown/applemusic/commit.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -o musicsearch .

FROM scratch
WORKDIR /app
COPY --from=appbuilder /build/musicsearch .
CMD ["./musicsearch"]
