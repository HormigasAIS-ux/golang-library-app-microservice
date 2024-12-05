#!/bin/bash

docker build -t backend-syn-auth-service:latest ./auth_service
docker build -t backend-syn-author-service:latest ./author_service
docker build -t backend-syn-book-service:latest ./book_service