---
description: "Run command and log it and it's output into a file"
tags: []
---
#!/usr/bin/env sh

FILE="$1"
shift

DATE="$(date --rfc-3339=seconds)"

echo "${DATE} ❯ $@" > "$FILE"
$@ | tee -a "$FILE"
