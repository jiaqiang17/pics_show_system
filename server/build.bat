# 交叉编译到Linux
set GOOS=linux
set GOARCH=amd64
go build -o discover-linux
