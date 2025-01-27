FROM golang:1.23

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

ARG GOOGLE_ID=$GOOGLE_ID

ARG GOOGLE_SECRET=$GOOGLE_SECRET
ENV SESSION_SECRET=`uuidgen`

ARG GOOGLE_CALLBACK_URL=$GOOGLE_CALLBACK_URL


WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -ldflags "-X github.com/jeffscottbrown/applemusic/auth.googleSecret=$GOOGLE_SECRET -X github.com/jeffscottbrown/applemusic/auth.googleId=GOOGLE_ID -X github.com/jeffscottbrown/applemusic/auth.callbackUrl=$GOOGLE_CALLBACK_URL -X github.com/jeffscottbrown/applemusic/commit.Hash=$GIT_COMMIT -X github.com/jeffscottbrown/applemusic/commit.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" -o app .

CMD ["./app"]
