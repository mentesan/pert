#!/usr/bin/bash
#
#sudo docker run -d --name mongodb -v /opt/docker-data/mongo:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo

# Redis
sudo docker restart 61c00adf0791
# MongoDB
sudo docker restart 09cd29813dcc
