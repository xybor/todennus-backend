#!/bin/bash

migrate create -ext=sql -dir=infras/database/$1/migration -seq $2
