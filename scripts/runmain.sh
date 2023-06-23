#!/usr/bin/bash
export MONGO_URI="mongodb://admin:eumesmo@localhost:27017/test?authSource=admin"
export MONGO_DATABASE="pert"
export JWT_SECRET="sapeca"
go run main.go
