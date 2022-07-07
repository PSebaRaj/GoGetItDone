install_swag:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: install_swag
	swagger generate spec -o swagger.yaml --scan-models
