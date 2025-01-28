# Golang Apple Music Search

## Overview

This is a small playground app written in Golang that uses the Apple Music API.

An instance of the application is generally running at 
https://godemo.jeffandbetsy.net right now.

## Running The App Locally

The project requires 3 OS environment variables to be set.

SESSION_SECRET - used internally by Gothic to manage session data. The value
should be random and at least 32 bytes.

GOOGLE_ID - Client id from Google OAuth2 API credentials.

GOOGLE_SECRET - Client secret from Google OAuth2 API credentials.

Optionally you can set GOOGLE_CALLBACK_URL which will be used to tell Google
OAuth2 where to call back into your application after authentication.  If you
do not set this value, a default value of
`http://localhost:8080/auth/google/callback` will be used.  The value must
end with `/auth/google/callback` and you may change the host name and/or the
port number. Note that this URL must be configured as authorized redirect URIs 
when configuring the credentials in the Google API console.

With those environment variables set, the application should be ready to run.
```bash
go run .
```

That will serve the browswer interface at http://localhost:8080 and
will serve a JSON API at `http://localhost:8080/search/[band name goes here]`,
for example http://localhost:8080/search/Phish.

Alternatively the application may be run in a Docker container.  

docker run -p 8080:8080 \
    -e GOOGLE_SECRET=$GOOGLE_SECRET \
    -e GOOGLE_ID=$GOOGLE_ID \
    -e GOOGLE_CALLBACK_URL=$GOOGLE_CALLBACK_URL \
    -p 8080:8080 docker.io/jeffscottbrown/applemusic:latest
