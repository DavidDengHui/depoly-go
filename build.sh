set -euxo pipefail

mkdir -p "$(pwd)/functions"
rm -rf "$(pwd)/api"
GOBIN=$(pwd)/functions go install ./...
chmod +x "$(pwd)"/functions/*
go env