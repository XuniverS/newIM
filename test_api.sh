#!/bin/bash
# IM 系统 API 测试脚本

BASE_URL="http://localhost:8080/api"
CONTENT_TYPE="Content-Type: application/json"

echo "🧪 IM 系统 API 测试"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数
TESTS_PASSED=0
TESTS_FAILED=0

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_code=$4
    local description=$5

    echo -n "测试: $description ... "

    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "$CONTENT_TYPE")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "$CONTENT_TYPE" \
            -d "$data")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ 通过${NC} (HTTP $http_code)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo "$body"
    else
        echo -e "${RED}✗ 失败${NC} (期望 $expected_code, 实际 $http_code)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        echo "$body"
    fi
    echo ""
}

# 1. 测试注册
echo -e "${YELLOW}1. 测试用户注册${NC}"
REGISTER_DATA='{"username":"testuser1","password":"password123"}'
test_api "POST" "/auth/register" "$REGISTER_DATA" "200" "用户注册"

# 2. 测试重复注册
echo -e "${YELLOW}2. 测试重复注册${NC}"
test_api "POST" "/auth/register" "$REGISTER_DATA" "400" "重复用户名应该失败"

# 3. 测试登录
echo -e "${YELLOW}3. 测试用户登录${NC}"
LOGIN_DATA='{"username":"testuser1","password":"password123"}'
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "$CONTENT_TYPE" \
    -d "$LOGIN_DATA")
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
USER_ID=$(echo $LOGIN_RESPONSE | grep -o '"user_id":[0-9]*' | cut -d':' -f2)

if [ -n "$TOKEN" ]; then
    echo -e "${GREEN}✓ 登录成功${NC}"
    echo "Token: $TOKEN"
    echo "User ID: $USER_ID"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}✗ 登录失败${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
echo ""

# 4. 测试错误密码
echo -e "${YELLOW}4. 测试错误密码${NC}"
WRONG_LOGIN='{"username":"testuser1","password":"wrongpassword"}'
test_api "POST" "/auth/login" "$WRONG_LOGIN" "401" "错误密码应该失败"

# 5. 测试生成密钥
echo -e "${YELLOW}5. 测试生成密钥${NC}"
if [ -n "$TOKEN" ]; then
    KEYS_RESPONSE=$(curl -s -X POST "$BASE_URL/keys/generate" \
        -H "$CONTENT_TYPE" \
        -H "Authorization: Bearer $TOKEN")
    
    if echo "$KEYS_RESPONSE" | grep -q "public_key"; then
        echo -e "${GREEN}✓ 密钥生成成功${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗ 密钥生成失败${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    echo "$KEYS_RESPONSE" | head -c 100
    echo "..."
else
    echo -e "${YELLOW}⊘ 跳过（无有效 Token）${NC}"
fi
echo ""

# 6. 测试获取在线用户
echo -e "${YELLOW}6. 测试获取在线用户${NC}"
if [ -n "$TOKEN" ]; then
    test_api "GET" "/users/online" "" "200" "获取在线用户列表"
else
    echo -e "${YELLOW}⊘ 跳过（无有效 Token）${NC}"
fi
echo ""

# 测试总结
echo "================================"
echo -e "测试总结:"
echo -e "  ${GREEN}通过: $TESTS_PASSED${NC}"
echo -e "  ${RED}失败: $TESTS_FAILED${NC}"
echo "================================"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}✗ 有测试失败${NC}"
    exit 1
fi
