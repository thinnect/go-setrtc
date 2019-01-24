# Makefile for embedding build info into the executable

BUILD_DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_DISTRO = $(shell lsb_release -sd)

all: mist-setrtc

raspberry: export GOOS=linux
raspberry: export GOARCH=arm
raspberry: export GOARM=6
raspberry: build

raspberry2: export GOOS=linux
raspberry2: export GOARCH=arm
raspberry2: export GOARM=7
raspberry2: build

mist-setrtc:
	go build -o mist-setrtc -ldflags "-X 'main.ApplicationBuildDate=$(BUILD_DATE)' -X 'main.ApplicationBuildDistro=$(BUILD_DISTRO)'"

install: mist-setrtc
	install -m0755  mist-setrtc /usr/bin
	install systemd/mist-setrtc.service /etc/systemd/system/mist-setrtc.service
	install -d -m0755 /etc/mist-setrtc
	install -m0664  systemd/mist-setrtc.conf /etc/mist-setrtc/mist-setrtc.conf

clean:
	rm mist-setrtc

.PHONY: mist-setrtc
