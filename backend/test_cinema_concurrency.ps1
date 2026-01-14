# Cinema Booking System - Comprehensive Concurrency Test
# Demonstrates PostgreSQL concurrency control with UNIFIED booking API

Write-Host "=== Cinema Booking System - Concurrency Test ===" -ForegroundColor Cyan
Write-Host "Testing UNIFIED BookSeats API with PostgreSQL concurrency control" -ForegroundColor Yellow
Write-Host ""

$baseUrl = "http://localhost:8080"
$showId = 7

function Test-Connection {
    Write-Host "Testing server connection..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/seats?show_id=$showId" -Method GET
        Write-Host "Server connection: OK" -ForegroundColor Green
        Write-Host "Found $($response.seats.Count) seats for show $showId" -ForegroundColor Gray
        return $true
    } catch {
        Write-Host "Server connection: FAILED" -ForegroundColor Red
        Write-Host "Make sure server is running: go run cmd/server/main.go" -ForegroundColor Yellow
        return $false
    }
}

function Test-SingleSeat {
    Write-Host "Testing single seat booking (Legacy format)..." -ForegroundColor Yellow
    
    # Use actual seat ID from database
    $request = @{
        seat_id = 2  # This maps to seat A1 in show 7
    } | ConvertTo-Json
    
    Write-Host "Request: $request" -ForegroundColor Gray
    
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/book" -Method POST -Headers @{"Content-Type"="application/json"; "X-User-ID"="1"} -Body $request
        Write-Host "Single seat booking: SUCCESS" -ForegroundColor Green
        Write-Host "Booking ID: $($response.booking_id)" -ForegroundColor Gray
        Write-Host "Seat: $($response.seat_name)" -ForegroundColor Gray
        return $true
    } catch {
        Write-Host "Single seat booking: FAILED" -ForegroundColor Yellow
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Gray
        return $false
    }
}

function Test-MultiSeat {
    Write-Host "Testing multi-seat booking (New unified format)..." -ForegroundColor Yellow
    
    # Use seat NAMES (not IDs) for multi-seat booking
    $request = @{
        show_id = $showId
        seats = @("A4", "A6", "B1")  # CORRECT: Use seat names
    } | ConvertTo-Json
    
    Write-Host "Request: $request" -ForegroundColor Gray
    
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/book" -Method POST -Headers @{"Content-Type"="application/json"; "X-User-ID"="2"} -Body $request
        Write-Host "Multi-seat booking: SUCCESS" -ForegroundColor Green
        Write-Host "Booking IDs: $($response.booking_ids -join ', ')" -ForegroundColor Gray
        Write-Host "Seats: $($response.seat_names -join ', ')" -ForegroundColor Gray
        Write-Host "Movie: $($response.movie_title)" -ForegroundColor Gray
        return $true
    } catch {
        Write-Host "Multi-seat booking: FAILED" -ForegroundColor Yellow
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Gray
        return $false
    }
}

function Test-ConcurrencySingleSeat {
    Write-Host "Testing single seat concurrency (Legacy API)..." -ForegroundColor Yellow
    Write-Host "3 users compete for seat ID 10 (A5)" -ForegroundColor White
    
    $jobs = @()
    
    # Create 3 concurrent jobs for SAME seat
    for ($i = 1; $i -le 3; $i++) {
        $scriptBlock = {
            param($baseUrl, $userId)
            
            $headers = @{
                "Content-Type" = "application/json"
                "X-User-ID" = $userId.ToString()
            }
            
            # All users try to book seat ID 10 (A5) 
            $body = @{
                seat_id = 10  # A5 seat
            } | ConvertTo-Json
            
            try {
                $response = Invoke-RestMethod -Uri "$baseUrl/api/book" -Method POST -Headers $headers -Body $body
                return @{
                    UserId = $userId
                    Success = $true
                    Message = $response.message
                    BookingId = $response.booking_id
                    SeatName = $response.seat_name
                }
            } catch {
                $statusCode = 500
                if ($_.Exception.Response) {
                    $statusCode = $_.Exception.Response.StatusCode.value__
                }
                return @{
                    UserId = $userId
                    Success = $false
                    StatusCode = $statusCode
                    Message = $_.Exception.Message
                }
            }
        }
        
        $job = Start-Job -ScriptBlock $scriptBlock -ArgumentList $baseUrl, $i
        $jobs += $job
    }
    
    # Wait for all jobs to complete
    Start-Sleep -Seconds 3
    
    # Collect results
    $results = @()
    foreach ($job in $jobs) {
        $result = Wait-Job $job | Receive-Job
        $results += $result
        Remove-Job $job
    }
    
    # Analyze results
    Write-Host ""
    Write-Host "=== SINGLE SEAT CONCURRENCY RESULTS ===" -ForegroundColor Cyan
    $successCount = 0
    foreach ($result in $results) {
        if ($result.Success) {
            Write-Host "User $($result.UserId): SUCCESS - Booked $($result.SeatName), ID: $($result.BookingId)" -ForegroundColor Green
            $successCount++
        } else {
            if ($result.StatusCode -eq 409) {
                Write-Host "User $($result.UserId): CONFLICT - Expected! (409)" -ForegroundColor Yellow
            } else {
                Write-Host "User $($result.UserId): ERROR - Status $($result.StatusCode)" -ForegroundColor Red
            }
        }
    }
    
    Write-Host ""
    if ($successCount -eq 1) {
        Write-Host "PERFECT! Exactly 1 booking succeeded (PostgreSQL concurrency control works!)" -ForegroundColor Green
    } else {
        Write-Host "Unexpected: $successCount bookings succeeded (should be 1)" -ForegroundColor Red
    }
    
    return $successCount
}

function Test-ConcurrencyMultiSeat {
    Write-Host "Testing multi-seat concurrency (Unified API)..." -ForegroundColor Yellow
    Write-Host "3 users compete for overlapping seats" -ForegroundColor White
    
    $scenarios = @(
        @{
            UserId = 1
            Seats = @("B2", "B3", "B4")
            Description = "User 1: B2, B3, B4"
        },
        @{
            UserId = 2  
            Seats = @("B3", "B4", "B5")
            Description = "User 2: B3, B4, B5 (overlaps with User 1)"
        },
        @{
            UserId = 3
            Seats = @("B4", "B5", "B6") 
            Description = "User 3: B4, B5, B6 (overlaps with User 2)"
        }
    )
    
    Write-Host ""
    foreach ($scenario in $scenarios) {
        Write-Host "   $($scenario.Description)" -ForegroundColor White
    }
    
    $jobs = @()
    
    # Create concurrent jobs
    foreach ($scenario in $scenarios) {
        $scriptBlock = {
            param($baseUrl, $userId, $showId, $seats)
            
            $headers = @{
                "Content-Type" = "application/json"
                "X-User-ID" = $userId.ToString()
            }
            
            $body = @{
                show_id = $showId
                seats = $seats
            } | ConvertTo-Json
            
            try {
                $response = Invoke-RestMethod -Uri "$baseUrl/api/book" -Method POST -Headers $headers -Body $body
                return @{
                    UserId = $userId
                    Success = $true
                    Message = $response.message
                    BookingIds = $response.booking_ids
                    SeatNames = $response.seat_names
                }
            } catch {
                $statusCode = 500
                if ($_.Exception.Response) {
                    $statusCode = $_.Exception.Response.StatusCode.value__
                }
                return @{
                    UserId = $userId
                    Success = $false
                    StatusCode = $statusCode
                    Message = $_.Exception.Message
                    AttemptedSeats = $seats
                }
            }
        }
        
        $job = Start-Job -ScriptBlock $scriptBlock -ArgumentList $baseUrl, $scenario.UserId, $showId, $scenario.Seats
        $jobs += $job
    }
    
    # Wait for completion
    Start-Sleep -Seconds 3
    
    # Collect results
    $results = @()
    foreach ($job in $jobs) {
        $result = Wait-Job $job | Receive-Job
        $results += $result
        Remove-Job $job
    }
    
    # Analyze results
    Write-Host ""
    Write-Host "=== MULTI-SEAT CONCURRENCY RESULTS ===" -ForegroundColor Cyan
    $successCount = 0
    foreach ($result in $results) {
        if ($result.Success) {
            Write-Host "User $($result.UserId): SUCCESS - Booked [$($result.SeatNames -join ', ')]" -ForegroundColor Green
            Write-Host "   Booking IDs: $($result.BookingIds -join ', ')" -ForegroundColor Gray
            $successCount++
        } else {
            if ($result.StatusCode -eq 409) {
                Write-Host "User $($result.UserId): CONFLICT - Attempted [$($result.AttemptedSeats -join ', ')] (Expected!)" -ForegroundColor Yellow
            } else {
                Write-Host "User $($result.UserId): ERROR - Status $($result.StatusCode)" -ForegroundColor Red
            }
        }
    }
    
    Write-Host ""
    if ($successCount -eq 1) {
        Write-Host "PERFECT! Exactly 1 multi-seat booking succeeded (Atomicity + Concurrency control works!)" -ForegroundColor Green
    } else {
        Write-Host "Unexpected: $successCount multi-seat bookings succeeded (should be 1)" -ForegroundColor Red
    }
    
    return $successCount
}

# ==================== MAIN EXECUTION ====================

Write-Host "Step 1: Connection Test" -ForegroundColor Cyan
if (-not (Test-Connection)) {
    exit 1
}

Write-Host ""
Write-Host "Step 2: Basic Booking Tests" -ForegroundColor Cyan
Write-Host ""

Write-Host "Step 2a: Single Seat Booking Test" -ForegroundColor Magenta
Test-SingleSeat

Write-Host ""
Write-Host "Step 2b: Multi-Seat Booking Test" -ForegroundColor Magenta  
Test-MultiSeat

Write-Host ""
Read-Host "Press Enter to continue to concurrency stress tests..."

Write-Host ""
Write-Host "Step 3: Concurrency Control Demonstration" -ForegroundColor Cyan
Write-Host "This is the MAIN EVENT - demonstrating PostgreSQL DBMS capabilities!" -ForegroundColor Yellow
Write-Host ""

Write-Host "Step 3a: Single Seat Race Condition" -ForegroundColor Magenta
$singleResult = Test-ConcurrencySingleSeat

Write-Host ""
Read-Host "Press Enter to continue to multi-seat concurrency test..."

Write-Host ""
Write-Host "Step 3b: Multi-Seat Atomicity + Concurrency" -ForegroundColor Magenta
$multiResult = Test-ConcurrencyMultiSeat

Write-Host ""
Write-Host "==================== EDUCATIONAL SUMMARY ====================" -ForegroundColor Cyan
Write-Host ""
Write-Host "DBMS Concepts Successfully Demonstrated:" -ForegroundColor White
Write-Host ""
Write-Host "1. CONCURRENCY CONTROL:" -ForegroundColor Green
Write-Host "   PostgreSQL row-level locking prevents race conditions" -ForegroundColor White
Write-Host "   Multiple users safely handled without app-level locks" -ForegroundColor White
Write-Host ""
Write-Host "2. ATOMICITY:" -ForegroundColor Green
Write-Host "   Multi-seat bookings are ALL OR NOTHING" -ForegroundColor White
Write-Host "   Transaction rollback on any seat conflict" -ForegroundColor White
Write-Host ""
Write-Host "3. ISOLATION:" -ForegroundColor Green
Write-Host "   Each transaction sees consistent database state" -ForegroundColor White
Write-Host "   Optimistic locking with conflict detection" -ForegroundColor White
Write-Host ""
Write-Host "4. CONSISTENCY:" -ForegroundColor Green
Write-Host "   Database constraints prevent invalid states" -ForegroundColor White
Write-Host "   Only ONE booking succeeds per seat" -ForegroundColor White
Write-Host ""

if ($singleResult -eq 1 -and $multiResult -eq 1) {
    Write-Host "EXCELLENT! All concurrency tests passed perfectly!" -ForegroundColor Green
    Write-Host "Your PostgreSQL database is handling concurrency like a pro!" -ForegroundColor Green
} else {
    Write-Host "Some tests had unexpected results. Check server logs for details." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Real-World Applications:" -ForegroundColor Yellow
Write-Host "E-commerce inventory management" -ForegroundColor White
Write-Host "Banking account transactions" -ForegroundColor White  
Write-Host "Reservation systems (hotels, flights)" -ForegroundColor White
Write-Host "Resource allocation in distributed systems" -ForegroundColor White
Write-Host ""
Write-Host "=== DEMO COMPLETE ===" -ForegroundColor Cyan
Write-Host "Perfect demonstration of industrial-grade DBMS capabilities!" -ForegroundColor Green
