# Cinema Booking Application

## Overview

The system is designed to address the problem of movie ticket booking in a highly concurrent environment, where a large number of users may access the application and perform actions simultaneously—especially scenarios in which multiple users attempt to book the same seat at the same time.

**Link project**: https://cinema-booking-system-hid9.onrender.com

## Description

The application allows users to:

- Register, log in, and authenticate using JWT;
- Browse movies, showtimes, and seat layouts;
- Perform real-time ticket booking;
- Manage booked tickets (view booking history, seat information, and showtime details);
- Ensure that each seat can be successfully booked by only one user, even when multiple booking requests are submitted concurrently.

In addition, the system is equipped with protective mechanisms such as Rate Limiting and API Keys to ensure backend stability and security.

## Architecture and Technologies

The system follows the **Backend for Frontend (BFF)** architecture, a variation of the API Gateway pattern, where the backend is responsible for handling business logic, security, and data orchestration.

### Core Technologies

#### 1. Go (Golang)

Go is used as the primary programming language due to:

- High performance and low memory overhead;
- Strong support for concurrency through goroutines;
- Built-in synchronization primitives such as atomic operations, mutexes, and channels, which help ensure data integrity in multi-threaded environments;
- Suitability for high-throughput and high-concurrency systems.

#### 2. Gin Framework

Gin is a lightweight and high-performance HTTP framework that provides:

- Clear and structured middleware support (JWT, Rate Limiting, CORS, API Key);
- Easy integration with Clean Architecture principles;
- Efficient request handling suitable for scalable backend services.

#### 3. PostgreSQL

PostgreSQL is used to store the core business data of the system, including users, movies, showtimes, seats, and bookings.

To address the concurrent seat booking problem, the system leverages database transactions combined with `SELECT … FOR UPDATE` to ensure data consistency in a high-contention environment. The mechanism works as follows:

- `SELECT ... FOR UPDATE` explicitly acquires **row-level** locks at read time, preventing other transactions from modifying or locking the same rows until the current transaction completes (**COMMIT** or **ROLLBACK**);
- This approach guarantees strict consistency in high-contention scenarios such as seat booking and prevents concurrency anomalies like **Lost Updates**.
- Fine-grained locking is applied by locking only the rows whose `seat_id` values belong to the `p_seat_ids` array. This allows other transactions to operate on different seats (rows) within the same table without being blocked.

##### Limitations of the application-level approach

In earlier versions of the system, seat booking logic was implemented at the application layer. Although the operations were executed within a single transaction, this approach had several drawbacks:

- Business logic was scattered across multiple functions;
- Refactoring introduced a higher risk of subtle bugs;
- New developers could easily omit a critical step in the workflow.

```
LockSeats()
CountSeatsForUpdate()
BookSeats()
CreateBooking()
```

##### Solution: Moving invariants to the Database

To overcome these issues, all logic related to **concurrent seat booking** was migrated to the **Database Management System (DBMS)** using a **stored function**.

This design ensures that **data invariants** are enforced directly at the database level, rather than relying on correctness in application-layer logic.

Advantages:

- **Strong atomicity**: The entire seat booking process is executed within a single database transaction;
- **Automatic rollback on failure**, independent of application code and error-handling logic;
- **Database-enforced invariants**: No client or service can violate business rules once they are embedded in the database;
- **Elimination of race conditions** by avoiding fragmented “check-then-update” patterns.

Disadvantages

- Requires manual configuration and maintenance of database-level logic;
- Demands a solid understanding of transactions, locking mechanisms, and PL/pgSQL from developers.

```
CREATE OR REPLACE FUNCTION book_seats(
    p_user_id bigint,
    p_seat_ids bigint[]
) RETURNS void AS $$
DECLARE
    p_updated_count int;
BEGIN
    PERFORM 1
    FROM seats
    WHERE seat_id = ANY (p_seat_ids)
    FOR UPDATE;

    UPDATE seats
    SET status = 'booked'
    WHERE seat_id = ANY (p_seat_ids)
      AND status = 'available';

    GET DIAGNOSTICS p_updated_count = ROW_COUNT;

    IF p_updated_count <> array_length(p_seat_ids, 1) THEN
        RAISE EXCEPTION 'Some seats already booked';
    END IF;

    INSERT INTO bookings(user_id, seat_id)
    SELECT p_user_id, unnest(p_seat_ids);
END;
$$ LANGUAGE plpgsql;
```

#### 4. Redis

Redis is primarily used for Rate Limiting:

- Reduces load on the backend and database during traffic spikes;
- Ensures accurate rate limiting across multiple application instances;
- Provides an automatic fallback to in-memory rate limiting when Redis is unavailable.

#### 5. JWT (JSON Web Token)

JWT is used for user authentication and authorization, with a clear separation between:

- A lightweight JWT used to identify logged-in users (mainly for user-based rate limiting);
- A full JWT used to protect APIs that require authentication.

## Rate Limiting Strategy

The system applies rate limiting at multiple levels:

- For unauthenticated users → rate limiting based on IP address;
- For authenticated users → rate limiting based on `user_id`;
- Redis is preferred to ensure consistency in distributed environments;
- If Redis becomes unavailable, the system automatically falls back to in-memory rate limiting.

This strategy helps:

- Mitigate DDoS attacks, brute-force attempts, and spam requests;
- Prevent server overload during traffic spikes;
- Avoid scenarios where a single client monopolizes system resources and degrades the experience for others.

## Authentication and User Management

After logging in, users receive an **Access Token** and a **Refresh Token**, which are stored in cookies with security attributes such as `HttpOnly`, `Secure`, and `SameSite`.

Since JWT is inherently stateless, a leaked token cannot be revoked until it expires. To address this limitation, the system combines Refresh Tokens with **Refresh Token Rotation** to enhance security.

Without refresh tokens, access tokens would need a long expiration time (e.g., 7 days) to avoid frequent re-authentication. If such a token were compromised, an attacker would have full access for the entire duration. With refresh tokens in place, access tokens can be kept very short-lived (5–15 minutes), significantly reducing the impact of token leakage.

### Refresh Token Rotation Flow

Each time a client uses an old refresh token to obtain a new access token, the server invalidates the old refresh token and issues a completely new token pair.

Benefits include:

- **Account takeover detection**: If an attacker uses a stolen refresh token first, the legitimate user's subsequent request will reveal that the token has already been used;
- **Self-revocation mechanism**: Upon detecting refresh token reuse, the system invalidates all active tokens associated with the user and forces re-authentication;
- **Reduced token lifetime exposure**: Even if a refresh token is stolen, it remains usable only for a very short window before becoming invalid.

```
Login
 ├─ Access Token (15m)
 └─ Refresh Token (RT₁)

Refresh (RT₁)
 ├─ New Access Token
 └─ New Refresh Token (RT₂)
    └─ RT₁ revoked

Refresh (RT₂)
 ├─ New Access Token
 └─ New Refresh Token (RT₃)
    └─ RT₂ revoked
```

JWT is used to:

- Authenticate users;
- Attach `user_id` to the request context;
- Control access between public APIs and authenticated APIs.

Public APIs (e.g., browsing movies and showtimes) remain accessible even when the user is logged in.

## Implementation Challenges

During system development, several key challenges were encountered:

- **Concurrent ticket booking**: Ensuring that no duplicate seat bookings occur when multiple users select the same seat simultaneously;
- **Middleware ordering**: Determining the correct execution order for Rate Limiting, JWT, and API Key middleware to maintain both security and correct business logic;
- **Rate limiting in distributed environments**: Combining Redis-based and in-memory rate limiting with Redis health checks to ensure resilience;
- **Development vs. production environments**: Handling differences between local environments (self-signed HTTPS) and production environments (Render, AWS) while maintaining strong security guarantees.

These challenges highlight real-world issues involved in building backend systems with high requirements for concurrency, security, and scalability.

## Future Enhancements

In future iterations, the system may be extended with:

- Online payment integration;
- Advanced role-based access control (admin, staff);
- Centralized monitoring and logging (Prometheus, Grafana);
- Migration toward a microservices architecture as the system scales;
- Enhanced JWT security by combining tokens with `session_id` for tighter session management.
