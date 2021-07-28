HOMEDIR := $(shell pwd)
OUTPUT := $(HOMEDIR)/output
SOURCE := $(shell ls *.go)
RM := rm
GO := $(GOROOT)/bin/go

all: clean build

build:
	for i in $(SOURCE); do \
		$(GO) build -o $(OUTPUT)/$(`basename $$i .go`) $$i; \
	done

clean:
	$(RM) -rf $(OUTPUT)

.PHONY: all build clean
