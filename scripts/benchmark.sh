#!/bin/bash
# ============================================================
# MachPay CLI - Performance Benchmarks
# ============================================================
#
# Measures CLI performance for key operations.
#
# Usage: ./scripts/benchmark.sh
#
# ============================================================

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m'

echo ""
echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  ⏱️  MachPay CLI - Performance Benchmarks                 ${BLUE}║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Build if needed
if [ ! -f "./machpay" ]; then
    echo "Building CLI..."
    go build -o machpay ./cmd/machpay
fi

# ============================================================
# Startup Time
# ============================================================
echo -e "${BLUE}## Startup Time${NC}"
echo ""

benchmark_cmd() {
    local name="$1"
    local cmd="$2"
    local iterations=5
    local total=0
    
    for i in $(seq 1 $iterations); do
        start=$(python3 -c "import time; print(int(time.time() * 1000))")
        eval "$cmd" >/dev/null 2>&1
        end=$(python3 -c "import time; print(int(time.time() * 1000))")
        elapsed=$((end - start))
        total=$((total + elapsed))
    done
    
    avg=$((total / iterations))
    printf "  %-35s %dms\n" "$name" "$avg"
}

benchmark_cmd "machpay --help" "./machpay --help"
benchmark_cmd "machpay version" "./machpay version"
benchmark_cmd "machpay status (no config)" "./machpay status"

echo ""

# ============================================================
# Go Benchmarks
# ============================================================
echo -e "${BLUE}## Go Benchmarks${NC}"
echo ""

# Run Go benchmarks if they exist
if go test -list 'Benchmark' ./... 2>/dev/null | grep -q Benchmark; then
    go test -bench=. -benchmem ./... 2>/dev/null | grep -E "^Benchmark|ns/op|allocs/op" | head -20
else
    echo "  No Go benchmarks defined yet."
fi

echo ""

# ============================================================
# Binary Size
# ============================================================
echo -e "${BLUE}## Binary Size${NC}"
echo ""

if [ -f "./machpay" ]; then
    size=$(ls -lh ./machpay | awk '{print $5}')
    echo "  CLI binary: $size"
fi

echo ""

# ============================================================
# Summary
# ============================================================
echo -e "${GREEN}Benchmark complete!${NC}"
echo ""
echo "Target metrics:"
echo "  - Startup time: <100ms"
echo "  - Binary size: <30MB"
echo ""

