#!/bin/bash

# Remove exec.Command calls from Go files
# Usage: ./remove_exec_commands.sh [directory]

TARGET_DIR=${1:-.}

find "$TARGET_DIR" -type f -name "*.go" | while read -r file; do
    echo "Processing $file"

    # Create temporary file
    tmpfile=$(mktemp)

    # Process the file
    awk '
    BEGIN { inExec = 0; braceCount = 0 }
    /exec\.Command/ {
        inExec = 1
        print "// TODO: Replace exec.Command with proper implementation"
        print "// Original line: " $0
        next
    }
    inExec {
        braceCount += gsub(/{/, "{")
        braceCount -= gsub(/}/, "}")
        if (braceCount <= 0) {
            inExec = 0
            braceCount = 0
        }
        next
    }
    { print }
    ' "$file" > "$tmpfile"

    # Overwrite original file if changes were made
    if ! cmp -s "$file" "$tmpfile"; then
        mv "$tmpfile" "$file"
        echo "Updated $file"
    else
        rm "$tmpfile"
    fi
done

echo "Done processing Go files in $TARGET_DIR"