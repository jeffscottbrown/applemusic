## Running The App

The project requires 3 OS environment variables to be set.

SESSION_SECRET - used internally by Gothic to manage session data. The value
should be random and at least 32 bytes.

GOOGLE_ID - Client id from Google OAuth2 API credentials.

GOOGLE_SECRET - Client secret from Google OAuth2 API credentials.

With those environment variables set, the application should be ready to run.
```bash
go run .
```

That will serve the browswer interface at http://localhost:8080 and
will serve a JSON API at `http://localhost:8080/search/[band name goes here]`,
for example http://localhost:8080/search/Phish.




