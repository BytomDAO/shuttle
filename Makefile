install:
	@echo "Installing swap to $(GOPATH)/bin"
	@go install ./cmd/swap

clean:
	@rm -rf $(GOPATH)/bin/swap