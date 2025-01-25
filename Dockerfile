FROM golang:1.23

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -ldflags "-X github.com/jeffscottbrown/applemusic/commit.Hash=$GIT_COMMIT -X github.com/jeffscottbrown/applemusic/commit.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -o app .

CMD ["./app"]
