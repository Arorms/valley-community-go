#!/usr/bin/env bash

curl -X POST http://localhost:3000/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Test content"}'