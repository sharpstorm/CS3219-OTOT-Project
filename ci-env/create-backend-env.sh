#!/bin/bash

echo -e "DATABASE_URL=postgres:5432\n"\
"DATABASE_USERNAME=ciuser\n"\
"DATABASE_PASSWORD=cipassword\n"\
"DATABASE_NAME=cidb" > backend/.env