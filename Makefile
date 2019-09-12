install:
	@echo "Installing swap to $(GOPATH)/bin"
	@go install ./cmd/swap
	@echo "Install done."

clean:
	@echo "Cleaning $(GOPATH)/bin/swap"
	@rm -rf $(GOPATH)/bin/swap
	@echo "Clean done."