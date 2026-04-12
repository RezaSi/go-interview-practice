#!/bin/bash

# Test script for Railway deployment
# This script tests the Docker build and basic functionality

set -e

echo "🚀 Testing Railway Deployment Configuration..."

# Test 1: Docker build
echo "📦 Testing Docker build..."
if docker build -f web-ui/Dockerfile -t go-interview-practice-test . ; then
    echo "✅ Docker build successful"
else
    echo "❌ Docker build failed"
    exit 1
fi

# Test 2: Container startup
echo "🔄 Testing container startup..."
CONTAINER_ID=$(docker run -d -p 8080:8080 go-interview-practice-test)

# Wait for container to start
echo "⏳ Waiting for container to start..."
sleep 10

# Test 3: Health check
echo "🏥 Testing health endpoint..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed"
    docker logs $CONTAINER_ID
    docker stop $CONTAINER_ID
    exit 1
fi

# Test 4: Main page
echo "🏠 Testing main page..."
if curl -f http://localhost:8080/ > /dev/null 2>&1; then
    echo "✅ Main page accessible"
else
    echo "❌ Main page not accessible"
    docker logs $CONTAINER_ID
    docker stop $CONTAINER_ID
    exit 1
fi

# Test 5: API endpoints
echo "🔌 Testing API endpoints..."
if curl -f http://localhost:8080/api/challenges > /dev/null 2>&1; then
    echo "✅ API endpoints working"
else
    echo "❌ API endpoints not working"
    docker logs $CONTAINER_ID
    docker stop $CONTAINER_ID
    exit 1
fi

# Cleanup
echo "🧹 Cleaning up..."
docker stop $CONTAINER_ID
docker rm $CONTAINER_ID
docker rmi go-interview-practice-test

echo "🎉 All tests passed! Railway deployment configuration is ready."
echo ""
echo "📋 Deployment Summary:"
echo "  ✅ Docker build works"
echo "  ✅ Container starts successfully"
echo "  ✅ Health check endpoint responds"
echo "  ✅ Main application is accessible"
echo "  ✅ API endpoints are functional"
echo ""
echo "🚀 Ready to deploy to Railway!"
echo "   Use the railway-template.json or railway-template-config.json"
echo "   Deploy URL: https://railway.app/template/go-interview-practice"
