#!/bin/bash

if [ "\$1" == "test" ]; then
  export $(cat test.env | xargs)
else
  export $(cat .env | xargs)
fi

goose -dir ./migrations postgres "$DATABASE_URL_MIGRATIONS" status
goose -dir ./migrations postgres "$DATABASE_URL_MIGRATIONS" up
