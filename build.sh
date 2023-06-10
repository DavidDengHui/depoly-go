set -euxo pipefail

mkdir -p functions
GOBIN=functions go install ./...
chmod +x functions/*
go env