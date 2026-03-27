#!/bin/bash
set -e

# -------- Input --------
LANGUAGE=$1
CODE_FILE=$2
INPUT_FILE=$3

# -------- Validation --------
if [ $# -lt 3 ]; then
  echo "INVALID_INPUT"
  exit 1
fi

if [ ! -f "$CODE_FILE" ]; then
  echo "CODE_FILE_NOT_FOUND"
  exit 1
fi

if [ ! -f "$INPUT_FILE" ]; then
  echo "INPUT_FILE_NOT_FOUND"
  exit 1
fi

# -------- Absolute paths --------
CODE_FILE=$(realpath "$CODE_FILE")
INPUT_FILE=$(realpath "$INPUT_FILE")

# -------- Setup workspace --------
WORKDIR=$(mktemp -d)
cp "$CODE_FILE" "$WORKDIR/"
cd "$WORKDIR" || exit 1

FILENAME=$(basename "$CODE_FILE")

# -------- Compile / Setup --------
case "$LANGUAGE" in

  # -------- C --------
  c)
    gcc "$FILENAME" -O2 -o main || { echo "CE"; exit 1; }
    RUN_CMD="./main"
    ;;

  # -------- C++ --------
  cpp)
    g++ "$FILENAME" -O2 -o main || { echo "CE"; exit 1; }
    RUN_CMD="./main"
    ;;

  # -------- Java --------
  java)
    javac "$FILENAME" || { echo "CE"; exit 1; }
    CLASS_NAME=$(basename "$FILENAME" .java)
    RUN_CMD="java $CLASS_NAME"
    ;;

  # -------- Python --------
  python)
    RUN_CMD="python3 $FILENAME"
    ;;

  # -------- Node --------
  node)
    RUN_CMD="node $FILENAME"
    ;;

  # -------- TypeScript --------
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
timeout 5s bash -c "$RUN_CMD < \"$INPUT_FILE\""
STATUS=$?

if [ $STATUS -eq 124 ]; then
  echo "TLE"
  exit 124
elif [ $STATUS -ne 0 ]; then
  echo "RE"
  exit $STATUS
fi