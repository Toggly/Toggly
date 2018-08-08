#!/usr/bin/env bash

go test ./... -coverprofile=sonar-coverage.out -json > sonar-report.json
