#!/usr/bin/env bash
# Check domain availability across multiple TLDs via RDAP.
#
# Usage:
#   ./scripts/domain-check.sh                 # check default name list
#   ./scripts/domain-check.sh foo bar baz     # check specific names
#   TLDS="com app io" ./scripts/domain-check.sh foo
#
# Notes:
# - RDAP returns 404 when a domain has no record (available).
# - 200 means the registry returned data (taken).
# - Some TLDs aren't RDAP-mapped on rdap.org → shows as "??".

set -uo pipefail

DEFAULT_NAMES=(
  bloomly
  smallbloom
  everbloom
  slowbloom
  monthbloom
  kindbloom
  bloomkind
  komorebi
  yutori
  hanami
  iyashi
  nagori
  softdays
  linger
  hush
  inkling
  lilypad
  stamen
  pollen
  calyx
  moodleaf
  leafmood
  paperflower
  softpetals
  thirtypetals
)

# Brand-relevant TLDs. Override with TLDS env var.
DEFAULT_TLDS=(
  com app io co me page studio garden day xyz life land
)

if [[ $# -gt 0 ]]; then
  NAMES=("$@")
else
  NAMES=("${DEFAULT_NAMES[@]}")
fi

if [[ -n "${TLDS:-}" ]]; then
  read -r -a TLDS <<< "$TLDS"
else
  TLDS=("${DEFAULT_TLDS[@]}")
fi

GREEN=$'\033[32m'
RED=$'\033[31m'
DIM=$'\033[2m'
RESET=$'\033[0m'

check() {
  local name="$1" tld="$2"
  local code
  code=$(curl -sL --max-time 10 -o /dev/null -w "%{http_code}" \
    "https://rdap.org/domain/${name}.${tld}")
  case "$code" in
    404)
      printf "  ${GREEN}✓${RESET}  %-22s ${GREEN}available${RESET}\n" "${name}.${tld}"
      ;;
    200)
      printf "  ${RED}✗${RESET}  %-22s ${DIM}taken${RESET}\n" "${name}.${tld}"
      ;;
    400|429)
      printf "  ?  %-22s ${DIM}rate-limited (HTTP ${code})${RESET}\n" "${name}.${tld}"
      ;;
    501|404\ *|"")
      printf "  ?  %-22s ${DIM}no rdap (HTTP ${code:-none})${RESET}\n" "${name}.${tld}"
      ;;
    *)
      printf "  ?  %-22s ${DIM}HTTP ${code}${RESET}\n" "${name}.${tld}"
      ;;
  esac
}

echo "checking ${#NAMES[@]} names × ${#TLDS[@]} TLDs (${#NAMES[@]} × ${#TLDS[@]} = $(( ${#NAMES[@]} * ${#TLDS[@]} )) lookups)"
echo

for n in "${NAMES[@]}"; do
  printf "${DIM}===${RESET} %s ${DIM}===${RESET}\n" "$n"
  for t in "${TLDS[@]}"; do
    check "$n" "$t"
    sleep 0.6  # rdap.org politeness — strict throttle
  done
  echo
done
