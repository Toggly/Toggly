#!/usr/bin/env bash

go test ./app/... -coverprofile=sonar-coverage.out -json > sonar-report.json
