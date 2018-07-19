#!/usr/bin/env bash

go test ./app/... -coverprofile=sonar-coverage.out -json > sonar-report.json

sonar-scanner \
  -Dsonar.projectKey=TooglyCore \
  -Dsonar.sources=app \
  -Dsonar.go.coverage.reportPaths=sonar-coverage.out \
  -Dsonar.go.tests.reportPaths=sonar-report.json \
  -Dsonar.host.url=http://localhost:9000
  