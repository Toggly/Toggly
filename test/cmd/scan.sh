#!/usr/bin/env bash

sonar-scanner \
  -Dsonar.projectKey=TooglyCore \
  -Dsonar.sources=app \
  -Dsonar.go.coverage.reportPaths=sonar-coverage.out \
  -Dsonar.go.tests.reportPaths=sonar-report.json \
  -Dsonar.host.url=http://localhost:9000
  