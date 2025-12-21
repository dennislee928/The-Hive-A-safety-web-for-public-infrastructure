# Makefile for Robot Framework tests

.PHONY: install test test-baseline test-performance test-security test-evaluation test-integration clean help

# Variables
PYTHON := python3
PIP := $(PYTHON) -m pip
ROBOT := robot
OUTPUT_DIR := results
BASE_URL := http://localhost:8080

# Default target
help:
	@echo "Available targets:"
	@echo "  install          - Install Python dependencies"
	@echo "  test             - Run all tests"
	@echo "  test-baseline    - Run baseline tests"
	@echo "  test-performance - Run performance tests"
	@echo "  test-security    - Run security tests"
	@echo "  test-evaluation  - Run evaluation metrics tests"
	@echo "  test-integration - Run integration tests"
	@echo "  clean            - Clean test results"
	@echo ""
	@echo "Variables:"
	@echo "  BASE_URL         - API base URL (default: http://localhost:8080)"
	@echo "  OUTPUT_DIR       - Output directory (default: results)"

# Install dependencies
install:
	$(PIP) install -r requirements.txt

# Run all tests
test:
	BASE_URL=$(BASE_URL) $(ROBOT) --outputdir $(OUTPUT_DIR) Scenarios/

# Run baseline tests
test-baseline:
	BASE_URL=$(BASE_URL) $(ROBOT) --outputdir $(OUTPUT_DIR) --include baseline Scenarios/

# Run performance tests
test-performance:
	BASE_URL=$(BASE_URL) $(ROBOT) --outputdir $(OUTPUT_DIR) --include performance Scenarios/

# Run security tests
test-security:
	BASE_URL=$(BASE_URL) $(ROBOT) --outputdir $(OUTPUT_DIR) --include security Scenarios/

# Run evaluation tests
test-evaluation:
	BASE_URL=$(BASE_URL) $(ROBOT) --outputdir $(OUTPUT_DIR) --include evaluation Scenarios/

# Run integration tests
test-integration:
	BASE_URL=$(BASE_URL) $(ROBOT) --outputdir $(OUTPUT_DIR) --include integration Scenarios/

# Clean test results
clean:
	rm -rf $(OUTPUT_DIR)
	find . -name "*.log" -delete
	find . -name "*.html" -delete
	find . -name "*.xml" -delete

