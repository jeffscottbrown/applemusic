FROM golang:1.24-alpine AS appbuilder

ARG GITHUB_TAG

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/jeffscottbrown/applemusic/commit.Tag=$GITHUB_TAG -X github.com/jeffscottbrown/applemusic/commit.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -o musicsearch .

FROM gcr.io/distroless/static-debian12

ARG PROJECT_ID
ENV PROJECT_ID=$PROJECT_ID
ENV GIN_MODE=release

WORKDIR /app
COPY --from=appbuilder /build/musicsearch .
COPY --from=appbuilder /build/web/assets/ ./web/assets/
CMD ["./musicsearch"]
