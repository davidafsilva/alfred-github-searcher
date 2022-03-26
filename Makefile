all: build

build:
	alfred build
release:
	alfred pack
clean:
	alfred clean