# export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I proto/ --go_out=. proto/*.proto