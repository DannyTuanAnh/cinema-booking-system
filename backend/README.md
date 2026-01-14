# Cinema Booking Backend - Concurrency Control Demo

## 📚 Educational Purpose

This backend system demonstrates **database-level concurrency control** for a Database Management Systems (DBMS) course. The focus is on showing how PostgreSQL handles concurrent seat booking requests without application-level locking.

## 🎯 Key Learning Objectives

1. **Transaction Management**: Understanding BEGIN, COMMIT, ROLLBACK
2. **Atomicity**: "All or nothing" principle in action
3. **Concurrency Control**: Multiple users competing for the same resource
4. **Optimistic Locking**: Using `UPDATE ... WHERE condition` approach
5. **Conflict Resolution**: Graceful handling of concurrent conflicts

## 🏗️ Architecture

```
backend/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── controller/             # HTTP handlers (Gin)
│   ├── service/               # Business logic & transactions
│   ├── store/                 # Data access layer (SQL)
│   ├── model/                 # Domain models
│   ├── db/                    # Database connection
│   └── middleware/            # Authentication & CORS
├── migrations/cinema.sql       # Database schema
└── go.mod                     # Dependencies
```

## 🔄 Concurrency Control Flow

When multiple users try to book the same seat simultaneously:

```sql
-- Step 1: Start transaction
BEGIN;

-- Step 2: Optimistic locking - only update if seat is AVAILABLE
UPDATE seats
SET status = 'BOOKED'
WHERE seat_id = $1 AND status = 'AVAILABLE';

-- Step 3: Check if update succeeded (rows affected > 0)
-- If no rows affected → seat already booked → ROLLBACK
-- If update succeeded → continue to step 4

-- Step 4: Create booking record
INSERT INTO bookings(user_id, seat_id) VALUES ($1, $2);

-- Step 5: Commit transaction
COMMIT;
```

**Result**: Only ONE user succeeds, others get a conflict error.

## 🚀 Setup Instructions

### 1. Database Setup

```bash
# Install PostgreSQL
# Create database
createdb cinema_db

# Run migrations
psql -d cinema_db -f migrations/cinema.sql
```

### 2. Backend Setup

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

### 3. Environment

Default database configuration in `internal/db/postgres.go`:

- Host: localhost
- Port: 5432
- User: postgres
- Password: password
- Database: cinema_db

## 📋 API Endpoints

| Method | Endpoint               | Description           | Auth Required |
| ------ | ---------------------- | --------------------- | ------------- |
| POST   | `/api/login`           | User authentication   | No            |
| GET    | `/api/movies`          | List all movies       | No            |
| GET    | `/api/shows?movie_id=` | List shows            | No            |
| GET    | `/api/seats?show_id=`  | Get seat availability | No            |
| POST   | `/api/book`            | **Book a seat**       | Yes           |
| GET    | `/api/my-bookings`     | User's bookings       | Yes           |

## 🧪 Testing Concurrency

### Sample Users (from migrations)

- Username: `user1`, Password: `password123`
- Username: `user2`, Password: `password123`
- Username: `user3`, Password: `password123`

### Test Scenario

1. **Login as multiple users**:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user1", "password": "password123"}'
```

2. **Get available seats**:

```bash
curl "http://localhost:8080/api/seats?show_id=1"
```

3. **Simulate concurrent booking** (run simultaneously):

```bash
# Terminal 1 - User 1
curl -X POST http://localhost:8080/api/book \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"seat_id": 1}'

# Terminal 2 - User 2
curl -X POST http://localhost:8080/api/book \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 2" \
  -d '{"seat_id": 1}'
```

**Expected Result**:

- One request: `200 OK` - Booking successful
- Other request: `409 Conflict` - Seat already booked

## 🔍 Key Code Locations

### Concurrency Control Logic

- **File**: `internal/service/booking_service.go`
- **Method**: `BookSeat()`
- **Critical Section**: Transaction with optimistic locking

### Database Transactions

- **File**: `internal/store/seat_store.go`
- **Method**: `UpdateSeatStatusInTx()`
- **Mechanism**: `UPDATE ... WHERE status = 'AVAILABLE'`

## 📊 Database Schema

### Core Tables

- **users**: User accounts
- **movies**: Movie information
- **shows**: Movie showtimes
- **seats**: Seat availability (concurrency target)
- **bookings**: Booking records

### Concurrency Target

The `seats` table is the main focus:

```sql
CREATE TABLE seats (
    seat_id SERIAL PRIMARY KEY,
    show_id INT NOT NULL REFERENCES shows(show_id),
    seat_name VARCHAR(5) NOT NULL,
    status VARCHAR(10) NOT NULL CHECK (status IN ('AVAILABLE', 'BOOKED')),
    UNIQUE (show_id, seat_name)
);
```

## 🎓 Educational Notes

### What This Demo Shows:

✅ **PostgreSQL handles concurrency automatically**
✅ **Transactions provide atomicity**
✅ **Optimistic locking prevents conflicts**
✅ **No application-level synchronization needed**

### What This Demo Doesn't Use:

❌ Go mutexes/channels (application locking)
❌ SERIALIZABLE isolation level
❌ Explicit row locking (SELECT FOR UPDATE)
❌ Complex distributed locking

## 🔧 Troubleshooting

### Common Issues

1. **Database Connection Failed**

   - Check PostgreSQL is running
   - Verify database credentials
   - Ensure `cinema_db` database exists

2. **Import Errors**

   - Run `go mod download`
   - Check Go version (requires Go 1.21+)

3. **Concurrency Not Working**
   - Ensure multiple requests are truly simultaneous
   - Check database transaction isolation level
   - Verify seat exists and is AVAILABLE

## 📖 Further Learning

This demo illustrates basic database concurrency concepts. In production systems, consider:

- Connection pooling
- Deadlock detection
- Performance optimization
- Distributed transactions
- Event-driven architectures

---

**Remember**: The goal is to demonstrate that **PostgreSQL handles concurrency at the database level**, making application code simpler and more reliable.
