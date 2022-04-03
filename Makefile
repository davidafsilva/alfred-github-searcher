all: build

build:
	alfred build
link: build
	alfred link
unlink:
	alfred unlink
release: build
	alfred pack
clean:
	alfred clean
