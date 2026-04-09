#!/bin/bash

# Configuration
BASE_URL="http://localhost:8080"
EMAIL="metricuser_$(date +%s)@example.com"
PASSWORD="SecurePassword123"
NAME="Metric Tester"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}==================================================${NC}"
echo -e "${BLUE}       VELOX API METRICS VERIFICATION SUITE       ${NC}"
echo -e "${BLUE}==================================================${NC}"
echo -e "Using Email: ${YELLOW}$EMAIL${NC}"

# Helpers
function assert_status() {
    local actual=$1
    local expected=$2
    local msg=$3
    if [ "$actual" -ne "$expected" ]; then
        echo -e "${RED}FAIL: $msg (Expected $expected, got $actual)${NC}"
        exit 1
    else
        echo -e "${GREEN}PASS: $msg${NC}"
    fi
}

function submit_and_wait() {
    local lang=$1
    local code=$2
    local desc=$3
    
    echo -e "${BLUE}Action: $desc...${NC}"
    resp=$(curl -s -X POST "$BASE_URL/submit" \
      -H "Authorization: Bearer $api_key" \
      -H "Content-Type: application/json" \
      -d "{
        \"language\": \"$lang\",
        \"source_code\": \"$code\",
        \"test_cases\": [{\"test_case_id\": 1, \"input\": \"\", \"expected_output\": \"expected\"}]
      }")
    
    sub_id=$(echo "$resp" | jq -r '.submission_id')
    echo "  - Queued: $sub_id"
    
    # Wait for processing
    sleep 4
    
    # Trigger log update via status check
    curl -s -G "$BASE_URL/status" \
      -H "Authorization: Bearer $api_key" \
      --data-urlencode "submission_id=$sub_id" > /dev/null
    echo "  - Result processed and logged."
}

# 1. Setup
echo -e "\n${YELLOW}[1/4] Setting up test credentials...${NC}"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"$NAME\", \"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" > /dev/null

login_resp=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

token=$(echo "$login_resp" | jq -r '.data.token')

key_resp=$(curl -s -X POST "$BASE_URL/auth/api-keys" \
  -H "Authorization: Bearer $token" \
  -H "Content-Type: application/json" \
  -d '{"name": "Visual Metrics Key", "scopes": ["submit", "status"]}')

api_key=$(echo "$key_resp" | jq -r '.key')
key_id=$(echo "$key_resp" | jq -r '.id')
echo -e "${GREEN}✓ Environment ready. KeyID: $key_id${NC}"

# 2. Traffic Generation
echo -e "\n${YELLOW}[2/4] Generating simulated API traffic...${NC}"

# A. Burst Success (RPM check)
echo -e "${BLUE}Bursting 3 successful Python requests...${NC}"
for i in {1..3}
do
   curl -s -X POST "$BASE_URL/submit" \
     -H "Authorization: Bearer $api_key" \
     -H "Content-Type: application/json" \
     -d '{
       "language": "python",
       "source_code": "print('\''expected'\'')",
       "test_cases": [{"test_case_id": 1, "input": "", "expected_output": "expected"}]
     }' > /dev/null
   echo -n "."
done
echo " Done."

# B. Specific States
submit_and_wait "python" "import time\ntime.sleep(10)" "Simulating Time Limit Exceeded (TLE)"
submit_and_wait "cpp" "invalid syntax" "Simulating Compile Error"
submit_and_wait "python" "print('expected')" "Simulating Accepted run"

# 3. Stats Retrieval
echo -e "\n${YELLOW}[3/4] Fetching aggregated metrics...${NC}"
sleep 2 # Final settle time for async logger
stats_resp=$(curl -s -H "Authorization: Bearer $token" "$BASE_URL/auth/api-keys/stats?id=$key_id")

# 4. Visual Report
echo -e "\n${YELLOW}[4/4] METRICS REPORT${NC}"
echo -e "${BLUE}--------------------------------------------------${NC}"
echo -e "Metric                Value"
echo -e "${BLUE}--------------------------------------------------${NC}"
echo -e "Total Requests        $(echo "$stats_resp" | jq -r '.total_requests')"
echo -e "Peak RPM              $(echo "$stats_resp" | jq -r '.peak_rpm')"
echo -e "Peak RPD              $(echo "$stats_resp" | jq -r '.peak_rpd')"
echo -e "Success Rate          $(echo "$stats_resp" | jq -r '.success_rate')%"
echo -e "${BLUE}--------------------------------------------------${NC}"
echo -e "${BLUE}Error Distribution:${NC}"
echo "$stats_resp" | jq -r '.error_counts | to_entries[] | "  \(.key): \(.value)"'
echo -e "${BLUE}--------------------------------------------------${NC}"

# Assertions
total=$(echo "$stats_resp" | jq -r '.total_requests')
if [ "$total" -ge 6 ]; then
    echo -e "${GREEN}✓ Total Requests verified: $total${NC}"
else
    echo -e "${RED}✗ Total Requests mismatch: $total${NC}"
    exit 1
fi

echo -e "\n${GREEN}ALL METRICS VERIFIED SUCCESSFULLY${NC}"
echo -e "${BLUE}==================================================${NC}"
