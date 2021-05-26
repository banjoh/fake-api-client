# Fake API client

A client library for our new and fresh Fake API service.

## Design choices
* The library is structured to have each resource type as a subpackage. At the moment, we only have the accounts resource
* Deprecated fields will not be implemented
* The library constructs it's own default HTTP client which has sane defaults for a production environment, but also allows users to inject their own HTTP client instance
* Error handling is implemented by capturing API specific errors in `APIError` error and wrapping other errors in an `error` object containing a description of the reason the error occured.
* All APIs have a `context.Context` object that users can use to manage the lifecycle of the request. They could for example have the request timeout after some duration.

## Example library usage

Creating a new client instance
```go
// Default constructor
client, err := accounts.New()

// Using dependency injections
client, err := accounts.NewWithClient(&http.Client{}, &client.DefaultRetrySleeper{})
```

Making a request to the backend
```go
accCreate := accounts.AccountCreate{}
ctx := context.Background()
acc, err := client.Create(ctx, &accCreate)
```

Making a request that we need to timeout if not complete by our given duration. This approach is also ideal if the library is used within a web application context where one would like to have the API calls cancellable.
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
* Install `docker` & `docker-compose` if you prefer to run tests in a containerized environment without having to install libraries and golang.
* Execute tests using `./scripts/tests` (required golang environment) or `docker-compose up --build`

## Improvement considerations
* Versioning the client library in conjunction with the platform APIs will be important to ensure compatibility.
* Since this client library essentially exposes a set of platform APIs, testing it using the contract testing approach will be very benefitial. A solution like [pact.io](https://docs.pact.io/) would be a good candidate.
* A ResourceAPI interface to have resource structs implement would be a good addition to enforce a contract for all resource types
* The implementation is missing rate limiting implementation which the platform APIs enforce. A back-off and retry logic for this will be required.
* Accepted technical debt is commented in the code base using a DEBT tag
