# frontend-service

Simple service for delivering the frontend to users.

## Installing

Clone repository

```bash
git clone https://github.com/polyse/frontend-service.git
cd frontend-service
```
Quick start:

```bash
make
```

Quick start in Docker container:

```bash
docker build -t frontend-service .
docker run --rm -p 9900:9900 frontend-service
```