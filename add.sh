#!/bin/bash

# ./add.sh example.com example2.com
# ./add.sh example.com example2.com --name thename
# ./add.sh example.com example2.com -n thename

cd $(dirname ${BASH_SOURCE[0]})

if [[ $# -eq 0 ]]; then
    echo "args empty"
    exit 0
fi

urls=()
name=""

# Parse options
while [[ $# -gt 0 ]]; do
    case "$1" in
    -n | --name | -name)
        name="$2"
        shift 2
        ;;
    *)
        urls+=("$1")
        shift
        ;;
    esac
done

if [ ! -n "$name" ]; then
    name="other"
fi

for url in "${urls[@]}"; do
    # Remove the protocol (http://, https://, ftp://, etc.)
    url=${url#*://}

    # Remove the path and query string (if any)
    url=${url%%/*}

    # Remove www. prefix (if any)
    url=$(echo "$url" | awk -F "." '{print $(NF-1)"."$NF}')

    # if exist
    grep -nr "\.$url" ./ >/dev/null
    if [[ $? -ne 1 ]]; then
        echo "$url exist"
        continue
    fi

    echo "add .$url to ./rules/$name"
    echo ".$url" >>./rules/$name
done

