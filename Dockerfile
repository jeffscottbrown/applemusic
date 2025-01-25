FROM golang:1.23

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -ldflags "-X github.com/jeffscottbrown/applemusic/commit.Hash=$(git rev-parse --short HEAD) -X github.com/jeffscottbrown/applemusic/commit.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -o app .

CMD ["./app"]
