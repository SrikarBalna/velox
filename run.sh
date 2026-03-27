#!/bin/bash
set -e

# -------- Input --------
LANGUAGE=$1
CODE_FILE=$2
INPUT_FILE=$3   # optional

# -------- Validation --------
if [ $# -lt 2 ]; then
  echo "INVALID_INPUT"
  exit 1
fi

if [ ! -f "$CODE_FILE" ]; then
  echo "CODE_FILE_NOT_FOUND"
  exit 1
fi

# -------- Absolute paths --------
CODE_FILE=$(realpath "$CODE_FILE")

if [ -n "$INPUT_FILE" ] && [ -f "$INPUT_FILE" ]; then
  INPUT_FILE=$(realpath "$INPUT_FILE")
else
  INPUT_FILE=""
fi

# -------- Setup workspace --------
WORKDIR=$(mktemp -d)
cp "$CODE_FILE" "$WORKDIR/"
cd "$WORKDIR" || exit 1

FILENAME=$(basename "$CODE_FILE")

# -------- Compile / Setup --------
case "$LANGUAGE" in

  c)
    gcc "$FILENAME" -O2 -o main || { echo "CE"; exit 1; }
    RUN_CMD="./main"
    ;;

  cpp)
    g++ "$FILENAME" -O2 -o main || { echo "CE"; exit 1; }
    RUN_CMD="./main"
    ;;

  java)
    javac "$FILENAME" || { echo "CE"; exit 1; }
    CLASS_NAME=$(basename "$FILENAME" .java)
    RUN_CMD="java $CLASS_NAME"
    ;;

  python)
    RUN_CMD="python3 $FILENAME"
    ;;

  node)
    RUN_CMD="node $FILENAME"
    ;;

  ts)
    command -v tsc >/dev/null 2>&1 || { echo "TS_NOT_INSTALLED"; exit 1; }
    tsc "$FILENAME" || { echo "CE"; exit 1; }
    JS_FILE=$(basename "$FILENAME" .ts).js
    RUN_CMD="node $JS_FILE"
    ;;

  *)
    echo "UNSUPPORTED"
    exit 1
    ;;
esac

# -------- Execute --------

if [ -n "$INPUT_FILE" ]; then
  timeout 5s bash -c "$RUN_CMD < \"$INPUT_FILE\""
else
  timeout 5s bash -c "$RUN_CMD"
fi

STATUS=$?

if [ $STATUS -eq 124 ]; then
  echo "TLE"
  exit 124
elif [ $STATUS -ne 0 ]; then
  echo "RE"
  exit $STATUS
fi