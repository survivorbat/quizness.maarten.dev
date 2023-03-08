# QQ

QQ is a quiz 

## Prerequisites

In order to run the database you need [Docker](https://docs.docker.com/get-docker/)
installed. For development you need [Go](https://golang.org/dl/).
If you want to use the Makefile for easy commands you also need to have [Make](https://www.gnu.org/software/make/) on
your system.

## Getting started

1. Set the `AUTH_CLIENT_SECRET` in your environment
2. Start the database using `make dr`
3. Start the server using `make server`
4. Start the frontend using `make ui`
5. Visit [localhost:3000](http://localhost:3000)
