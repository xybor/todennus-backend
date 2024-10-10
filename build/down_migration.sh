#!/bin/bash

set -a # automatically export all variables
. ./config/.env
set +a

if [ $1 = "postgres" ];
then
    echo "Down postgres"
    migrate -source file://./infras/database/postgres/migration -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:5432/$POSTGRES_DB down $2
fi
