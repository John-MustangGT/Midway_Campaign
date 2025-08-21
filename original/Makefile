# Makefile for Midway_Campaign Go program

BINARY_NAME = midway_campaign
SOURCE_DIR = .

.PHONY: all build clean run

all: build

build:
	cd $(SOURCE_DIR) && go build -o ../$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

