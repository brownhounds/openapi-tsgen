#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <new_version>"
  exit 1
fi

new_version="$1"

version_file="./ver/VERSION"

if [ ! -f "$version_file" ]; then
  echo "Error: $version_file not found!"
  exit 1
fi

old_version=$(cat "$version_file")

sed -i "s/$old_version/$new_version/g" "./ver/VERSION"
sed -i "s/$old_version/$new_version/g" "./README.md"
sed -i "s/$old_version/$new_version/g" "./tools/install-linux-amd64.sh"
sed -i "s/$old_version/$new_version/g" "./tools/install-linux-arm64.sh"

echo "$new_version" > "$version_file"
