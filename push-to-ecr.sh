#!/bin/bash

# AWS ECR repository info
AWS_ACCOUNT_ID=034362052544
REGION=us-east-1
REPOSITORY_NAME=auth-lambda

# Docker image tag
TAG=latest
ECR_URI="$AWS_ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/$REPOSITORY_NAME:$TAG"

# Log in to ECR
echo "Logging in to ECR..."
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com

# Build the Docker image using docker-compose
echo "Building the Docker image with docker-compose..."
docker-compose up --build -d

# Tag the Docker image
echo "Tagging the image..."
docker tag auth-lambda:latest $ECR_URI

# Push the image to ECR
echo "Pushing the image to ECR..."
docker push $ECR_URI

echo "Docker image pushed to ECR: $ECR_URI"
