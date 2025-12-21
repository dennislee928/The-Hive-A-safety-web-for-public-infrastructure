#!/bin/bash
# Run Robot Framework tests for ERH Safety System PoC

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
TEST_SUITE=""
TAGS=""
OUTPUT_DIR="results"
BASE_URL="http://localhost:8080"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -s|--suite)
      TEST_SUITE="$2"
      shift 2
      ;;
    -t|--tags)
      TAGS="$2"
      shift 2
      ;;
    -o|--output)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    -u|--url)
      BASE_URL="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [OPTIONS]"
      echo "Options:"
      echo "  -s, --suite SUITE    Run specific test suite (e.g., Baseline_Test, Performance_Test)"
      echo "  -t, --tags TAGS      Run tests with specific tags (e.g., baseline,performance)"
      echo "  -o, --output DIR     Output directory for test results (default: results)"
      echo "  -u, --url URL        Base URL for API (default: http://localhost:8080)"
      echo "  -h, --help           Show this help message"
      exit 0
      ;;
    *)
      echo -e "${RED}Unknown option: $1${NC}"
      exit 1
      ;;
  esac
done

# Export BASE_URL for Robot Framework
export BASE_URL

# Create output directory
mkdir -p "${OUTPUT_DIR}"

# Build robot command
ROBOT_CMD="robot --outputdir ${OUTPUT_DIR}"

# Add suite or tags
if [ -n "${TEST_SUITE}" ]; then
  ROBOT_CMD="${ROBOT_CMD} tests/Scenarios/${TEST_SUITE}.robot"
elif [ -n "${TAGS}" ]; then
  ROBOT_CMD="${ROBOT_CMD} --include ${TAGS} tests/Scenarios/"
else
  ROBOT_CMD="${ROBOT_CMD} tests/Scenarios/"
fi

echo -e "${GREEN}Running Robot Framework tests...${NC}"
echo -e "${YELLOW}Command: ${ROBOT_CMD}${NC}"
echo ""

# Run tests
eval ${ROBOT_CMD}

# Check exit code
if [ $? -eq 0 ]; then
  echo ""
  echo -e "${GREEN}Tests completed successfully!${NC}"
  echo -e "${GREEN}Results are in: ${OUTPUT_DIR}${NC}"
else
  echo ""
  echo -e "${RED}Tests failed!${NC}"
  echo -e "${RED}Check results in: ${OUTPUT_DIR}${NC}"
  exit 1
fi

