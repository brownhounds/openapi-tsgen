#!/usr/bin/env sh

set -eu

root_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
fixtures_dir="$root_dir/tests/fixtures"
snapshots_dir="$root_dir/tests/snapshots"

mkdir -p "$snapshots_dir"
rm -f "$snapshots_dir"/*.snapshot.ts

for fixture in "$fixtures_dir"/*.fixture.yml; do
  [ -e "$fixture" ] || continue
  base=$(basename "$fixture" .fixture.yml)
  out="$snapshots_dir/$base.yml.snapshot.ts"
  go run . -s "$fixture" -o "$out"
done

for fixture in "$fixtures_dir"/*.fixture.json; do
  [ -e "$fixture" ] || continue
  base=$(basename "$fixture" .fixture.json)
  out="$snapshots_dir/$base.json.snapshot.ts"
  go run . -s "$fixture" --input-json -o "$out"
done
