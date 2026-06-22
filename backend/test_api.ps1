# API Test Script
Write-Host "=== Social Network API Tests ===" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:8080/api"

# Test 1: Register a new user
Write-Host "Test 1: Register new user" -ForegroundColor Yellow
$registerData = @{
    email = "test@example.com"
    password = "password123"
    first_name = "John"
    last_name = "Doe"
    date_of_birth = "1990-01-01"
    nickname = "johndoe"
    about_me = "Hello, I'm John!"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "$baseUrl/auth/register" -Method POST -Body $registerData -ContentType "application/json" -UseBasicParsing
    Write-Host "✓ Registration successful: $($response.StatusCode)" -ForegroundColor Green
    Write-Host $response.Content
} catch {
    Write-Host "✗ Registration failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test 2: Login
Write-Host "Test 2: Login with credentials" -ForegroundColor Yellow
$loginData = @{
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri "$baseUrl/auth/login" -Method POST -Body $loginData -ContentType "application/json" -SessionVariable session -UseBasicParsing
    Write-Host "✓ Login successful: $($response.StatusCode)" -ForegroundColor Green
    $user = $response.Content | ConvertFrom-Json
    Write-Host "User ID: $($user.id)"
    Write-Host "Email: $($user.email)"
    Write-Host "Name: $($user.first_name) $($user.last_name)"
    
    # Extract cookies
    $sessionId = $response.Headers['Set-Cookie']
    Write-Host "Session Cookie: $sessionId"
} catch {
    Write-Host "✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test 3: Get current user (requires session)
Write-Host "Test 3: Get current user info" -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/auth/me" -Method GET -WebSession $session -UseBasicParsing
    Write-Host "✓ Get me successful: $($response.StatusCode)" -ForegroundColor Green
    Write-Host $response.Content
} catch {
    Write-Host "✗ Get me failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# Test 4: Get user profile
Write-Host "Test 4: Get user profile" -ForegroundColor Yellow
if ($user) {
    try {
        $response = Invoke-WebRequest -Uri "$baseUrl/users/$($user.id)" -Method GET -WebSession $session -UseBasicParsing
        Write-Host "✓ Get profile successful: $($response.StatusCode)" -ForegroundColor Green
        Write-Host $response.Content
    } catch {
        Write-Host "✗ Get profile failed: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""

# Test 5: Update profile privacy
Write-Host "Test 5: Update profile privacy" -ForegroundColor Yellow
if ($user) {
    $privacyData = @{
        is_private = $true
    } | ConvertTo-Json
    
    try {
        $response = Invoke-WebRequest -Uri "$baseUrl/users/$($user.id)/privacy" -Method PUT -Body $privacyData -ContentType "application/json" -WebSession $session -UseBasicParsing
        Write-Host "✓ Update privacy successful: $($response.StatusCode)" -ForegroundColor Green
        Write-Host $response.Content
    } catch {
        Write-Host "✗ Update privacy failed: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""

# Test 6: Logout
Write-Host "Test 6: Logout" -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/auth/logout" -Method POST -WebSession $session -UseBasicParsing
    Write-Host "✓ Logout successful: $($response.StatusCode)" -ForegroundColor Green
    Write-Host $response.Content
} catch {
    Write-Host "✗ Logout failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "=== Tests Complete ===" -ForegroundColor Cyan
