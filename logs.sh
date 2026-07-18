#!/bin/bash
# Stream logs from the production server.
#
# Usage:
#   ./logs.sh               # tail follow
#   ./logs.sh -n 50         # last 50 lines then follow
#   ./logs.sh --since "10 min ago"   # recent logs, then follow
#   ./logs.sh --no-pager -n 100      # last 100 lines, no follow
#
# For one-shot (no follow), add --no-pager.
# See journalctl(1) for all options.

args=()
for a in "$@"; do
	args+=("$(printf '%q' "$a")")
done
ssh mhc "journalctl -u realworld-fiber-clean -f ${args[*]}"
