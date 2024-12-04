#!/bin/bash

docker tag backend-syn-auth-service:latest sakku116/backend-syn-auth-service:latest
docker push sakku116/backend-syn-auth-service:latest

docker tag backend-syn-author-service:latest sakku116/backend-syn-author-service:latest
docker push sakku116/backend-syn-author-service:latest