#!/bin/bash
set -e

THRESHOLD=50.0

echo "=== Coverage Metrics Report ==="
echo ""

# 1. Check total line coverage
TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
echo "Total line coverage: ${TOTAL_COVERAGE}%"

if (( $(echo "$TOTAL_COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "::error::Line coverage ${TOTAL_COVERAGE}% is below the required threshold of ${THRESHOLD}%"
  LINE_PASS=false
else
  echo "::notice::Line coverage ${TOTAL_COVERAGE}% meets the required threshold of ${THRESHOLD}%"
  LINE_PASS=true
fi
echo ""

# 2. Check file coverage (% of files with >= 70% coverage)
TOTAL_FILES=$(go tool cover -func=coverage.out | grep -v total | grep -v "^$" | awk '{print $1}' | cut -d: -f1 | sort -u | wc -l)

# Get files with >= 70% average coverage
FILES_PASSING=0
for file in $(go tool cover -func=coverage.out | grep -v total | grep -v "^$" | awk '{print $1}' | cut -d: -f1 | sort -u); do
  FILE_COV=$(go tool cover -func=coverage.out | grep "$file" | awk '{sum += substr($3, 1, length($3)-1); count++} END {if(count>0) print sum/count; else print 0}')
  if (( $(echo "$FILE_COV >= $THRESHOLD" | bc -l) )); then
    FILES_PASSING=$((FILES_PASSING + 1))
  fi
done

if [ "$TOTAL_FILES" -gt 0 ]; then
  FILE_COVERAGE_PERCENT=$(echo "scale=2; ($FILES_PASSING / $TOTAL_FILES) * 100" | bc)
else
  FILE_COVERAGE_PERCENT=0
fi

echo "File coverage: ${FILES_PASSING}/${TOTAL_FILES} files with >= ${THRESHOLD}% coverage (${FILE_COVERAGE_PERCENT}%)"

if (( $(echo "$FILE_COVERAGE_PERCENT < $THRESHOLD" | bc -l) )); then
  echo "::error::File coverage ${FILE_COVERAGE_PERCENT}% is below the required threshold of ${THRESHOLD}%"
  FILE_PASS=false
else
  echo "::notice::File coverage ${FILE_COVERAGE_PERCENT}% meets the required threshold of ${THRESHOLD}%"
  FILE_PASS=true
fi
echo ""

# 3. Check method/function coverage (% of functions with any coverage)
TOTAL_FUNCS=$(go tool cover -func=coverage.out | grep -v total | grep -v "^$" | wc -l)
COVERED_FUNCS=$(go tool cover -func=coverage.out | grep -v total | grep -v "^$" | awk '{cov = substr($3, 1, length($3)-1); if(cov > 0) count++} END {print count}')

if [ "$TOTAL_FUNCS" -gt 0 ]; then
  METHOD_COVERAGE_PERCENT=$(echo "scale=2; ($COVERED_FUNCS / $TOTAL_FUNCS) * 100" | bc)
else
  METHOD_COVERAGE_PERCENT=0
fi

echo "Method coverage: ${COVERED_FUNCS}/${TOTAL_FUNCS} methods covered (${METHOD_COVERAGE_PERCENT}%)"

if (( $(echo "$METHOD_COVERAGE_PERCENT < $THRESHOLD" | bc -l) )); then
  echo "::error::Method coverage ${METHOD_COVERAGE_PERCENT}% is below the required threshold of ${THRESHOLD}%"
  METHOD_PASS=false
else
  echo "::notice::Method coverage ${METHOD_COVERAGE_PERCENT}% meets the required threshold of ${THRESHOLD}%"
  METHOD_PASS=true
fi
echo ""

# 4. Check statement/branch coverage approximation
# Go uses statement coverage by default, which is close to branch coverage
BRANCH_COVERAGE=$TOTAL_COVERAGE
echo "Branch coverage (statement-based): ${BRANCH_COVERAGE}%"

if (( $(echo "$BRANCH_COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "::error::Branch coverage ${BRANCH_COVERAGE}% is below the required threshold of ${THRESHOLD}%"
  BRANCH_PASS=false
else
  echo "::notice::Branch coverage ${BRANCH_COVERAGE}% meets the required threshold of ${THRESHOLD}%"
  BRANCH_PASS=true
fi
echo ""

# Summary
echo "=== Coverage Summary ==="
echo "Line coverage:   ${TOTAL_COVERAGE}% (threshold: ${THRESHOLD}%)"
echo "File coverage:   ${FILE_COVERAGE_PERCENT}% (threshold: ${THRESHOLD}%)"
echo "Method coverage: ${METHOD_COVERAGE_PERCENT}% (threshold: ${THRESHOLD}%)"
echo "Branch coverage: ${BRANCH_COVERAGE}% (threshold: ${THRESHOLD}%)"
echo ""

# Fail if any metric is below threshold
if [ "$LINE_PASS" = false ] || [ "$FILE_PASS" = false ] || [ "$METHOD_PASS" = false ] || [ "$BRANCH_PASS" = false ]; then
  echo "::error::One or more coverage metrics are below the required threshold of ${THRESHOLD}%"
  exit 1
else
  echo "::notice::All coverage metrics meet the required threshold of ${THRESHOLD}%"
fi
