#!/bin/bash

function check(){
    if test $( pgrep -f $1 | wc -l ) -ne 0
    then
        echo "结束进程: $1"
        killall $1
    fi
    nohup ./$1 -c ~/free-im/free.yaml > http_app.out 2>&1 &
}

go build ./cmd/http_app
go build ./cmd/logic
go build ./cmd/tcp_conn

check http_app
check logic
check tcp_conn