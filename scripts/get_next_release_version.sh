#!/usr/bin/env bash

# This script returns the next release version depending on the given release-major.minor version

# standard bash error handling
set -o nounset  # treat unset variables as an error and exit immediately.
set -o errexit  # exit immediately when a command fails.
set -E          # needs to be set if we want the ERR trap
set -o pipefail # prevents errors in a pipeline from being masked

BRANCH_NAME="$1"

MAJOR_MINOR_VERSION=${BRANCH_NAME#"release-"}

RESPONSE=$(curl -s "https://api.github.com/repos/kyma-project/eventing-publisher-proxy/releases")

LATEST_PATCH_VERSION=$(echo "$RESPONSE" | jq -r --arg version "$MAJOR_MINOR_VERSION" '
  map(select(.tag_name | startswith("v" + $version + ".") or startswith($version + ".")))
  | map(.tag_name | ltrimstr("v") | ltrimstr($version + "."))
  | map(select(test("^[0-9]+$")))
  | map(. | tonumber)
  | max
  | if . then "\($version).\(.)" else null end
')

# If no version found, set the patch version to 0
NEXT_PATCH_VERSION=$(echo "$LATEST_PATCH_VERSION" | awk -F'.' '{print ($3 == "" || $3 == "unset") ? 0 : $3 + 1}')

# Print the next release version
NEXT_RELEASE_VERSION="$MAJOR_MINOR_VERSION.$NEXT_PATCH_VERSION"
echo "Next release version: $NEXT_RELEASE_VERSION"
