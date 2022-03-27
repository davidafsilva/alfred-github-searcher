all: build

build:
	alfred build
release: build
	alfred pack
clean:
	alfred clean