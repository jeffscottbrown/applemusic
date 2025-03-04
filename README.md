# Golang Apple Music Search

## Overview

This is a small playground app written in Golang that uses the Apple Music API.

An instance of the application is generally running at 
https://godemo.jeffandbetsy.net right now.

## Running The App Locally

The project requires OS environment variables to be set in order
to configure OAuth.  Both Google and GitHub are supported.

`GOOGLE_ID`, `GITHUB_ID` - Client id for OAuth2 API credentials.

`GOOGLE_SECRET`, `GITHUB_SECRET` - Client secret for OAuth2 API credentials.

Optionally you can set `GOOGLE_CALLBACK_URL` and `GITHUB_CALLBACK_URL` which will be used 
to tell the provider where to call back into your application after authentication.  If you
do not set this value, a default value of
`http://localhost:8080/auth/google/callback` will be used.  The value must
end with `/auth/google/callback` and you may change the host name and/or the
port number. Note that this URL must be configured as an authorized redirect URI 
when configuring the credentials.

When deploying to GCP the system will attempt to resolve all of those values
from the Google Secrets Manager.  In order for that to work you will need
to set an environment variable named `PROJECT_ID` which is the id of the 
GCP project which contains the secrets.

With those environment variables set, the application should be ready to run.

```bash
go tool templ generate
go run .
```

The application may be run using [Air](https://github.com/air-verse/air).  Note that 
when using air to launch the app, air is configured to generate the templates 
automatically so there is no need to do that explicitly.

```bash
go tool air
```

That will serve the browswer interface at http://localhost:8080.

Alternatively the application may be run in a Docker container.  

```bash
docker run -p 8080:8080 \
    -e GOOGLE_SECRET=$GOOGLE_SECRET \
    -e GOOGLE_ID=$GOOGLE_ID \
    -e GOOGLE_CALLBACK_URL=$GOOGLE_CALLBACK_URL \
    -e GITHUB_SECRET=$GITHUB_SECRET \
    -e GITHUB_ID=$GITHUB_ID \
    -e GITHUB_CALLBACK_URL=$GITHUB_CALLBACK_URL \
    docker.io/jeffscottbrown/applemusic:latest
```
