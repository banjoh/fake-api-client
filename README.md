# Fake API client

A client library for our new and fresh Fake API service.

### Design choices
* Have subpackage to allow future entities being added
* Not implementing deprecated attributes
* Repository design pattern to separate different entities (future improvements)

### Example library usage
TODO: Provide an example of using the client in an application

### Dev requirements

* golang
* pre-commit
* An IDE

### Testing requirements
* Install `docker` & `docker-compose`
* Execute tests using `go test ./...` or `docker-compose up --build`

### Possible improvements
* ResourceAPI interface for other types of resources to implement

## TODO
* Increase test coverage
    * Error cases for Create & Fetch
    * Network transport errors e.g timeouts, no connection, DNS??, broken connections (io.Reader)
    * Retry logic
* Race conditions
* Memory usage and leaks???
* Document API
    * Expected errors
    * Input parameters
    * Default http client
* Review http.Client connection parameters
