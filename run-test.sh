#!/bin/bash

COMPOSE_PROJECT_NAME=otot-b-test docker-compose -f test.docker-compose.yaml up --abort-on-container-exit
