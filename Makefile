GO ?= go
GOPATH ?= $(HOME)/go

$(GOPATH)/bin/reflex: go.mod
	$(GO) install github.com/cespare/reflex@v0.3.1

run: $(GOPATH)/bin/reflex
	reflex -c reflex.conf
