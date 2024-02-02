#!/bin/bash

# показывает процент покрытия тестами
# игнорирует файлы и папки с именем mock и generated и ...
# указать вот тут:
FILTER="mock|generated|easyjson|migrations|tool"

go test -v -coverpkg=./... -coverprofile=profile.cov.unfiltered ./... 
cat profile.cov.unfiltered | grep -Ev $FILTER > profile.cov
go tool cover -func profile.cov
unlink profile.cov
unlink profile.cov.unfiltered