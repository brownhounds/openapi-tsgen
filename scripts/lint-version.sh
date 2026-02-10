#!/bin/bash

CHANGELOG_FILE="CHANGELOG.md"
VERSION_FILE="VERSION"

if [[ ! -f "$CHANGELOG_FILE" ]]; then
  echo "Changelog file not found!"
  exit 1
fi

if [[ ! -f "$VERSION_FILE" ]]; then
  echo "VERSION file not found!"
  exit 1
fi

latest_version=$(grep -E "^## \[[0-9]+\.[0-9]+\.[0-9]+\]" "$CHANGELOG_FILE" | head -n 1 | awk -F'[][]' '{print $2}')

if [[ -z "$latest_version" ]]; then
  echo "No version found in the changelog."
  exit 1
fi

version_file=$(cat "$VERSION_FILE")

if [[ -z "$version_file" ]]; then
  echo "The VERSION file is empty."
  exit 1
fi

if [[ "$latest_version" == "$version_file" ]]; then
  echo "The versions match: $latest_version"
else
  echo "The versions do not match!"
  echo "Changelog version: $latest_version"
  echo "VERSION file: $version_file"
  exit 1
fi
