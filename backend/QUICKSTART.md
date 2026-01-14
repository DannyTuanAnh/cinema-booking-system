# 🚀 Cinema Booking System - Complete Guide

## 🎯 **What This Demonstrates**

Advanced **PostgreSQL concurrency control** for cinema seat booking with:

- ✅ **Atomicity**: ALL seats booked or NONE (no partial booking)
- ✅ **Row-level Locking**: `SELECT ... FOR UPDATE`
- ✅ **Optimistic Concurrency**: `UPDATE ... WHERE status = 'available'`
- ✅ **Conflict Detection**: Check `rowsAffected` count
- ✅ **Transaction Isolation**: One transaction sees consistent state

## Prerequisites

1. **PostgreSQL** installed and running
2. **Go 1.21+** installed
3. **Database created**: `createdb cinema_db`

## 🚀 1-Minute Setup

```bash
# 1. Setup database
psql -d cinema_db -f migrations/cinema.sql

# 2. Run server (from backend directory)
go run cmd/server/main.go
```

Server will start at: http://localhost:8080

## 🧪 Testing Concurrency Control

### **Option 1: Automated Test (Recommended):**

```powershell
.\test_cinema_concurrency.ps1
```

### **Option 2: Manual Testing**

#### **Single Seat Booking:**

```bash
curl -X POST http://localhost:8080/api/book \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"seat_id": 1}'
```

#### **Multi-Seat Booking:**

```bash
curl -X POST http://localhost:8080/api/book \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"show_id": 7, "seats": ["A1", "A2", "A3"]}'
```

#### **Concurrency Test (Key Demo!):**

```bash
# Terminal 1 - User 1
curl -X POST http://localhost:8080/api/book \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"show_id": 7, "seats": ["A1", "A2"]}'

# Terminal 2 - User 2 (run simultaneously)
curl -X POST http://localhost:8080/api/book \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 2" \
  -d '{"show_id": 7, "seats": ["A2", "A3"]}'
```

## 📊 **Expected Results**

### ✅ **Success Response (200 OK):**

```json
{
  "booking_ids": [15, 16, 17],
  "seat_names": ["A1", "A2", "A3"],
  "movie_title": "Avengers: Endgame",
  "show_time": "2025-12-17 20:00",
  "message": "Successfully booked 3 seats for Avengers: Endgame at 20:00"
}
```

### ❌ **Conflict Response (409 Conflict):**

```json
{
  "error": "Seats no longer available",
  "message": "One or more seats have been booked by other users. Please select different seats.",
  "details": "seats no longer available: [A2] - already booked by other users"
}
```

**Key Point:** ✅ Only ONE request succeeds, others get 409 Conflict - this proves PostgreSQL's concurrency control works! 🎉

## 📋 **API Reference**

### **Unified Booking Endpoint:**

- **URL**: `POST /api/book`
- **Auth**: Required (`X-User-ID` header)
- **Supports**: Both single and multi-seat booking

### **Request Formats:**

**Single Seat (Legacy Support):**

```json
{ "seat_id": 123 }
```

**Multi-Seat (Primary Format):**

```json
{
  "show_id": 7,
  "seats": ["A1", "A2", "A3"]
}
```

### **Other Endpoints:**

| Method | Endpoint               | Auth | Description            |
| ------ | ---------------------- | ---- | ---------------------- |
| GET    | `/api/movies`          | No   | List movies            |
| GET    | `/api/shows?movie_id=` | No   | List shows             |
| GET    | `/api/seats?show_id=1` | No   | Show seat availability |
| GET    | `/api/my-bookings`     | Yes  | User's bookings        |

## 🎓 **Educational Value & DBMS Concepts**

### **What the Code Demonstrates:**

1. **Transaction Scope**: Everything in one `BEGIN...COMMIT`
2. **Row Locking**: `FOR UPDATE` serializes access
3. **Bulk Operations**: Update multiple rows atomically
4. **Conflict Detection**: Compare expected vs actual `rowsAffected`
5. **Rollback Logic**: Automatic cleanup on any failure

### **Concurrency Test Scenarios:**

| Scenario         | User 1 Seats | User 2 Seats | Expected Result         |
| ---------------- | ------------ | ------------ | ----------------------- |
| Complete Overlap | A1, A2, A3   | A1, A2, A3   | 1 succeeds, 1 conflicts |
| Partial Overlap  | A1, A2       | A2, A3       | 1 succeeds, 1 conflicts |
| No Overlap       | A1, A2       | B1, B2       | Both succeed            |

### **Key Teaching Points:**

1. **PostgreSQL does ALL concurrency control** - not Go application
2. **Transactions provide atomicity** across multiple operations
3. **Row-level locks prevent race conditions** automatically
4. **Optimistic concurrency** performs better than pessimistic locking
5. **Database guarantees consistency** even under high concurrent load

## 🔍 **Server Logs to Watch**

When testing concurrency, monitor server logs for:

```
User 1 attempting to book multiple seats: [A1 A2] for show 7
Concurrency conflict detected: Seats [A2] already booked by other users
SUCCESS: User 1 booked 2 seats [A1 A2] for show 7, booking IDs: [15 16]
```

## 🗄️ **Sample Data Setup**

If you need additional test data:

```sql
psql -U postgres -d cinema_db

-- Add more seats for testing
INSERT INTO seats (show_id, seat_name, status) VALUES
(7, 'A1', 'available'), (7, 'A2', 'available'), (7, 'A3', 'available'),
(7, 'A4', 'available'), (7, 'A5', 'available'), (7, 'B1', 'available'),
(7, 'B2', 'available'), (7, 'B3', 'available'), (7, 'B4', 'available');
```

**Default Data:**

- **Users**: user1, user2, user3 (password: password123)
- **Movies**: Avatar: The Way of Water, Top Gun: Maverick
- **Shows**: Multiple showtimes available
- **Seats**: Various seats per show

---

**Perfect for demonstrating DBMS principles in action! 🏆**

This system showcases how PostgreSQL handles concurrent access to shared resources without any application-level locking mechanisms.
