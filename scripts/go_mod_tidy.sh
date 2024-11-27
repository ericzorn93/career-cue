#!/bin/bash

# Update all go modules with the remote dependency
for module_dir in $(find . -name go.mod -type f); do
    cd "$(dirname "$module_dir")"
    echo "Go Tidying $module_dir"
    go mod tidy
    # Pipe to stdout is null
    cd - 1>/dev/null
done