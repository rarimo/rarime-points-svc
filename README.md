# rarime-points-svc

## Description

Core service of Rarimo Points System

## Install

```
git clone github.com/rarimo/rarime-points-svc
cd rarime-points-svc
go build main.go
export KV_VIPER_FILE=./config.yaml
./main migrate up
./main run service
```

## API documentation

[Online docs](https://rarimo.github.io/rarime-points-svc/) are available.

All endpoints from docs MUST be publicly accessible.

### Private endpoints

Private endpoints are not documented and MUST only be accessible within the
internal network. They do not require authorization in order to simplify back-end
interactions with Points service. Package [connector](./pkg/connector) provides
functionality to interact with these endpoints.

The path for internal endpoints is `/integrations/rarime-points-svc/v1/private/*`.

### Add referrals

Private endpoint to set usage count for genesis referral code or create a new
_System user_ with genesis referral code. _System user_ is unable to claim events or
withdraw, it has `is_disabled` attribute set to `true`, so the client app should
not allow it interactions with the system, although it is technically possible
to do other actions.

Path: `/integrations/rarime-points-svc/v1/private/referrals`
Body:
```json
{
    "nullifier": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "count": 2,
}
```
Response:
```json
{
  "referral": "kPRQYQUcWzW",
  "usage_left": 2
}
```

Parameters:
- `nullifier` - nullifier to create or edit referrals for
- `count` - number of referral usage

### Local build

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it
```
  cd docs
  npm install
  npm run start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.

Note: if you are using Gitlab for building project `docs/spec/paths` folder must not be
empty, otherwise only `Build and Publish` job will be passed.  

## Running from Source

* Run dependencies, based on config example
* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `migrate up` command to create database schema
* Launch the service with `run service` command

### Database
For services, we do use ***PostgresSQL*** database. 
You can [install it locally](https://www.postgresql.org/download/) or use [docker image](https://hub.docker.com/_/postgres/).

## Testing
To run tests need to mock some logic in handlers/middleware:
```go
func AuthMiddleware(auth *auth.Client, log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// claims, err := auth.ValidateJWT(r)
			// if err != nil {
			// 	log.WithError(err).Info("Got invalid auth or validation error")
			// 	ape.RenderErr(w, problems.Unauthorized())
			// 	return
			// }

			// if len(claims) == 0 {
			// 	ape.RenderErr(w, problems.Unauthorized())
			// 	return
			// }

			ctx := CtxUserClaims([]resources.Claim{{Nullifier: r.Header.Get("nullifier")}})(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
```
and in handlers/verify_passport:
```go
	// never panics because of request validation
	// proof.PubSignals[zk.Nullifier] = mustHexToInt(nullifier)
	// err = Verifier(r).VerifyProof(*proof)
	// if err != nil {
	// 	return nil, problems.BadRequest(err)
	// }
```

Run service with standart config (you need to configure db url only) and run tests.