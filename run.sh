#!/bin/bash

function check(){
    if test $( pgrep -f $1 | wc -l ) -neq 0
    then
        nohup ./$1 -c ~/free-im/free.yaml > http_app.out 2>&1 &
    fi
}

go build ./cmd/http_app
go build ./cmd/logic
go build ./cmd/tcp_conn

