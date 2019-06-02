#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
[ -f ~/.washhubrc ] && source ~/.washhubrc
# echo "$@"
# development mode
# cd ${DIR} && go run main.go "$@"
# production mode
${DIR}/washhub "$@"
