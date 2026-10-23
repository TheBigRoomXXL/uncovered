
go test -coverpkg=./... -coverprofile=coverage.txt ./...
gocover-cobertura < coverage.txt > coverage.xml
go run . coverage.xml
