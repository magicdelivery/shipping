# Magic delivery: the 'Shipping' microservice

## Intro

This is part of a small project for learning microservices development.

Currently, the 'Shipping' application manages client addresses and finds the geographically closest Parcel Locker addresses.

## Install and run

Before setting up the project, you need to have a Docker infrastructure installed to run the containers, as well as the `make` utility.
First, start the ['Parcel Locker'](https://github.com/magicdelivery/parcel_locker) microservice, followed by the 'Shipping' microservice, as the latter interacts with the former through its network.

```sh
mkdir magicdelivery && cd magicdelivery
git clone https://github.com/magicdelivery/parcel_locker.git
cd parcel_locker
make up
cd ..
git clone https://github.com/magicdelivery/shipping.git
cd shiping
make up
```

## Details

The application of this microservice is written in Go. The architecture is designed following Domain Driven Design (DDD) principles. The code structure is divided into the following parts:

* `internal/app` - application logic,
* `internal/infra` - infrastructure logic,
* `internal/domain` - domain business logic.

It is expected that higher levels can call lower levels, but not vice versa.

The application layer uses the [Gin](https://gin-gonic.com/) framework to serve the RESTful API.

The application uses Redis for data storage and caching the results of requests to the 'Parcel Locker' microservice. Caching is implemented using the [gocache](https://github.com/eko/gocache) library.

Communication with the 'Parcel Locker' microservice is carried out via synchronous request-response interactions over the HTTP protocol. To enhance reliability, request retries are implemented using an [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff) algorithm.

The [testify](https://github.com/stretchr/testify) library is used for unit testing the logic, allowing for the generation of mocks for out-of-process dependencies.

E2E testing of the API is implemented using the [hurl](https://hurl.dev/) command line tool. You can run the tests with the command `make testapi`, which will build and start the necessary containers and then execute the hurl test requests.

## Contribution

Since this is a learning project focused on microservices with Golang, I would be incredibly grateful for any advice or ideas for improvement!

## Licence

MIT
