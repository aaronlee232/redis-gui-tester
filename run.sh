#!/usr/bin/env bash
set -euo pipefail

# Development
npx @tailwindcss/cli -i ./ui/static/input.css -o ./ui/static/output.css

go tool templ generate
go run *.go