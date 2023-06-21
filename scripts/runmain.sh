#!/usr/bin/bash
export MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin"
export MONGO_DATABASE="demo"
export JWT_SECRET="sapeca"
go run main.go
