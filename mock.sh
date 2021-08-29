export GOPATH=/Users/klint/go
export PATH=$PATH:$(go env GOPATH)/bin

go get -u github.com/golang/mock/gomock
go get -u github.com/golang/mock/mockgen

mockgen --source server.go -package kim -destination server_mock.go
mockgen --source storage.go -package kim -destination storage_mock.go
mockgen --source dispatcher.go -package kim -destination dispatcher_mock.go
