# This being a container for running tests, not much
# optimization for size reduction is required
FROM golang:buster

WORKDIR /accounts-client
COPY . .
RUN ls -al
