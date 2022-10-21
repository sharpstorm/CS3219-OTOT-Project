#!make

include .env
export

run:
	go run .

build:
	go build -o ../dist/backend

generate:
	go generate ./...

test:
	go test backend.cs3219.comp.nus.edu.sg/auth backend.cs3219.comp.nus.edu.sg/controller  backend.cs3219.comp.nus.edu.sg/database

test-e2e:
	go test -p 1 backend.cs3219.comp.nus.edu.sg/e2e
