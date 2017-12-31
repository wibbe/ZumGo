
.PHONY: all setup build run release clean

all: build

setup:
	go get github.com/constabulary/gb/...

bin/libplatform.a: src/app/platform/platform.h src/app/platform/platform.cpp src/app/platform/platform_win32.cpp
	g++ -c src/app/platform/platform.cpp -I src/app/platform -o bin/platform.o -std=c++11 -g -Wall -Wundef
	ar rcs bin/libplatform.a bin/platform.o


build: bin/libplatform.a
	gb build

run: build
	bin/zum.exe

release: bin/libplatform.a
	gb build -ldflags="-s -w -H windowsgui"

clean:
	del bin\libplatform.a
	del bin\zum.exe
