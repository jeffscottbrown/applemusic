FROM golang:1.23

ARG GIT_COMMIT
ARG PROJECT_ID
ENV SESSION_SECRET=`uuidgen`

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -ldflags "-X github.com/jeffscottbrown/applemusic/secrets.projectId=$PROJECT_ID -X github.com/jeffscottbrown/applemusic/commit.Hash=$GIT_COMMIT -X github.com/jeffscottbrown/applemusic/commit.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -o app .

CMD ["./app"]
