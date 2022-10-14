# twoStage
Simple Go proejct for two-stage operations(2FA, SMS codes etc.)

## How to start
1. `docker-compose up -d` for tests and checking what is going on
2. `make`/`make test` depending on what you want
3. Start the `./main` binary

## Contract

* GET `/init` - returns you operation ID and saved code in JSON. Always returns `200 OK`
* POST `/commit` - with JSON `{"id":"...", "code":"..."}` checks if operation with such ID and code still exists in the system.
`417` on code mismatch, `404` on wrong/non-existing ID, `400` on malformed body
* `403` on every other route