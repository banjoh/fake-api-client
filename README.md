# Fake API client

A client library for our new and fresh Fake API service.

## Design choices
* The library is structured to have each resource type as a subpackage. At the moment, we only have the accounts resource
* Deprecated fields will not be implemented
* The library constructs it's own default HTTP client which has sane defaults for a production environment, but also allows users to inject their own HTTP client instance
* Error handling is implemented by capturing context of API specific errors in `APIError` error and wrapping other errors with error messages indicating the reason for the error.
* All APIs have a `context.Context` object that users can use to manage the lifecycle of the request. They could for example have the request timeot after some duration.

## Example library usage

Creating a new client instance
```go
client, err := accounts.New()

client, err := accounts.NewWithClient(&http.Client{}, &client.DefaultRetrySleeper{})
```

Making a request to the backend
```go
ctx := context.Background()
acc, err := client.Create(ctx, &accCreate)
```

Making a request that can timeout to the backend
```go
ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
defer cancel()
acc, err := client.Create(ctx, &accCreate)
```

## Development requirements

To make changes to this project you will need to have the following tools in your environment.

* golang
* pre-commit (developer gating checks)
* An IDE

## Testing requirements
* Install `docker` & `docker-compose` if you prefer to run tests using using containers
* Execute tests using `./scripts/tests` or `docker-compose up --build`

## Improvements considerations
* Versioning the client library in conjunction with the platform APIs will be important to ensure compatibility.
* Since this client library essentially exposes a set of platform APIs, testing it using the contract testing approach will be very benefitial. A solution like PACT would be a good candidate.
* A ResourceAPI interface would be a good addition to enforce a contract for all resource types
* The implementation is missing rate limiting implementation which the platform APIs enforce. A back-off and retry logic for this will be required.
* Accepted technical debt is commented in the code base using a DEBT tag

## TODO
* Review http.Client connection parameters
