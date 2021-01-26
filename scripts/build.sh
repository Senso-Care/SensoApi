#!/bin/bash

if [ "$#" -lt 1 ] || ! [ -d "./cmd/$1" ] > /dev/null; then
  echo "Usage: $0 CMDNAME [OS/PLATFORM]" >&2
  echo "Check os/platform available by running: \"go tool dist list\""
  exit 1
fi

os="linux"
platform="amd64"

if [ ! -z "$2" ] && go tool dist list | grep "$2" > /dev/null; then
  os=$(echo "$2" | cut -d"/" -f1)
  platform=$(echo "$2" | cut -d"/" -f2)
fi

if [ "$2" == "linux/arm/v7" ]; then
  os="linux"
  platform="arm"
fi

cmd_path="./cmd/$1"
echo "Compiling $1 on $os/$platform"
CGO_ENABLED=0 GOOS=$os GOARCH=$platform go build -o="bin/$2/$1" -x "./$cmd_path"
