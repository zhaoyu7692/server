# !/usr/bash
PROJECT_PATH=/home/zhaoyu/桌面/OpenJudge/Server
go env -w GOROOT=/usr/local/go #gosetup
go env -w GOPATH=$PROJECT_PATH #gosetup
/usr/local/go/bin/go build -o $PROJECT_PATH/server $PROJECT_PATH/src/main/main.go #gosetup

