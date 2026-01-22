#!/bin/bash
set -e

fmt_out=$(gofmt -l .)
if [ -n "$fmt_out" ]; then
  echo "These files are not properly formatted:"
  echo "$fmt_out"
  exit 1
fi
