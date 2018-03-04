# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=midimaggot

all: build
build:
	$(GOBUILD) -v
	$(GOBUILD) -v ./cmd/midimaggot
	cd ../..
install:	
	$(GOINSTALL) -v
	$(GOINSTALL) -v ./cmd/midimaggot
	cd ../..

