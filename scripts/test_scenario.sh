#!/bin/bash

# Configuration
API_URL="http://localhost:8080/api/v1"
TIMESTAMP=$(date +%s)

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper Functions
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

check_jq() {
    if ! command -v jq &> /dev/null; then
        log_error "jq is required but not installed."
        exit 1
    fi
}

register() {
    local email="user$1_${TIMESTAMP}@test.com"
    local username="user$1_${TIMESTAMP}"
    local password="password123"
    
    # Register
    curl -s -X POST "$API_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$email\",\"username\":\"$username\",\"password\":\"$password\"}" > /dev/null
    
    # Login
    local resp=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$email\",\"password\":\"$password\"}")
        
    local token=$(echo "$resp" | jq -r '.data.token')
    
    if [ -z "$token" ] || [ "$token" == "null" ]; then
        log_error "Login failed for $email"
        echo "$resp"
        exit 1
    fi
    
    echo "$token"
}

create_quiz() {
    local token=$1
    local resp=$(curl -s -X POST "$API_URL/quizzes" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d '{"title":"Realtime Quiz Test","description":"Automated Test Quiz"}')
        
    local id=$(echo "$resp" | jq -r '.data.id')
    if [ -z "$id" ] || [ "$id" == "null" ]; then
        log_error "Failed to create quiz"
        echo "$resp"
        exit 1
    fi
    echo "$id"
}

add_question() {
    local token=$1
    local quiz_id=$2
    local text=$3
    local correct=$4
    local resp=$(curl -s -X POST "$API_URL/quizzes/$quiz_id/questions" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{\"text\":\"$text\",\"options\":[\"A\",\"B\"],\"correct_answer\":\"$correct\",\"points\":100,\"time_limit\":30}")
        
    local id=$(echo "$resp" | jq -r '.data.id')
    if [ -z "$id" ] || [ "$id" == "null" ]; then
        log_error "Failed to add question"
        echo "$resp"
        exit 1
    fi
    echo "$id"
}

join_quiz() {
    local token=$1
    local code=$2
    curl -s -X POST "$API_URL/quizzes/join" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{\"code\":\"$code\"}" > /dev/null
}

submit_answer() {
    local token=$1
    local quiz_id=$2
    local question_id=$3
    local answer=$4
    
    curl -s -X POST "$API_URL/quizzes/$quiz_id/submit" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{\"question_id\":\"$question_id\",\"answer\":\"$answer\"}"
}

get_leaderboard() {
    local token=$1
    local quiz_id=$2
    curl -s -X GET "$API_URL/quizzes/$quiz_id/leaderboard" \
        -H "Authorization: Bearer $token"
}

# Main Execution
check_jq

log_info "Initializing Test Scenario..."
log_info "API URL: $API_URL"

# 1. Register Users
log_info "Registering Users..."
TOKEN_OWNER=$(register "owner")
TOKEN_P1=$(register "p1")
TOKEN_P2=$(register "p2")
TOKEN_P3=$(register "p3")

log_info "Users Registered & Logged In."

# 2. Create Quiz
log_info "Creating Quiz..."
QUIZ_ID=$(create_quiz "$TOKEN_OWNER")
# Get Quiz Code (need to fetch quiz details)
QUIZ_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN_OWNER" "$API_URL/quizzes/$QUIZ_ID")
QUIZ_CODE=$(echo "$QUIZ_DETAILS" | jq -r '.data.code')
log_info "Quiz Created: ID=$QUIZ_ID, Code=$QUIZ_CODE"

# 3. Add Questions
log_info "Adding Questions..."
Q1_ID=$(add_question "$TOKEN_OWNER" "$QUIZ_ID" "Is this valid?" "A")
Q2_ID=$(add_question "$TOKEN_OWNER" "$QUIZ_ID" "Is 1+1=2?" "A")
log_info "Questions Added."

# 4. Join Quiz
log_info "Players Joining Quiz..."
join_quiz "$TOKEN_P1" "$QUIZ_CODE"
join_quiz "$TOKEN_P2" "$QUIZ_CODE"
join_quiz "$TOKEN_P3" "$QUIZ_CODE"
log_info "Players Joined."

log_warn "📢  Please open WebSocket connections in separate terminals to see real-time updates!"
log_warn "    Run: wscat -c \"ws://localhost:8080/api/v1/ws?token=<TOKEN>\""
echo "    Owner Token: $TOKEN_OWNER"
echo "    P1 Token:    $TOKEN_P1"
echo "    P2 Token:    $TOKEN_P2"
echo "    P3 Token:    $TOKEN_P3"
echo "    Send Join:   {\"type\":\"join_quiz\",\"payload\":{\"quiz_id\":\"$QUIZ_ID\"}}"
echo ""

read -p "Press ENTER to start verify submission flow..."

# 5. Question 1 Flow
log_info "--- Question 1 Flow ---"
# P1: Correct (Fast)
log_info "P1 answering Correct..."
submit_answer "$TOKEN_P1" "$QUIZ_ID" "$Q1_ID" "A" > /dev/null

# P2: Correct (Slow)
sleep 1
log_info "P2 answering Correct (Slower)..."
submit_answer "$TOKEN_P2" "$QUIZ_ID" "$Q1_ID" "A" > /dev/null

# P3: Incorrect
log_info "P3 answering Incorrect..."
submit_answer "$TOKEN_P3" "$QUIZ_ID" "$Q1_ID" "B" > /dev/null

# Verify Leaderboard
log_info "Verifying Leaderboard Q1..."
LB_Q1=$(get_leaderboard "$TOKEN_OWNER" "$QUIZ_ID")
echo "$LB_Q1" | jq .

# 6. Question 2 Flow
log_info "--- Question 2 Flow ---"
# P3: Correct (Redemption)
log_info "P3 answering Correct..."
submit_answer "$TOKEN_P3" "$QUIZ_ID" "$Q2_ID" "A" > /dev/null

# P2: Incorrect
log_info "P2 answering Incorrect..."
submit_answer "$TOKEN_P2" "$QUIZ_ID" "$Q2_ID" "B" > /dev/null

# P1: Correct
log_info "P1 answering Correct..."
submit_answer "$TOKEN_P1" "$QUIZ_ID" "$Q2_ID" "A" > /dev/null

# Final Leaderboard
log_info "--- Final Leaderboard ---"
LB_FINAL=$(get_leaderboard "$TOKEN_OWNER" "$QUIZ_ID")
echo "$LB_FINAL" | jq .

log_info "Test Scenario Completed."
