#!/bin/bash

export CI_NAME="heroku"
export GIT_COMMITTED_AT="$(date +%s)"
export TZ="UTC"

go test -coverprofile=c.out
RETURN_VALUE=$?
./cc-test-reporter after-build -t gocov
exit $RETURN_VALUE
