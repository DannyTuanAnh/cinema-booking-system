# 🎬 Cinema Booking System - Frontend

Frontend application cho hệ thống đặt vé xem phim, xây dựng bằng HTML/CSS/JavaScript thuần (Vanilla JS).

## 📁 Cấu trúc thư mục

```
frontend/
├── index.html                      # Trang chủ
├── pages/                          # Các trang chức năng
│   ├── login.html                 # Đăng nhập
│   ├── register.html              # Đăng ký
│   ├── movies.html                # Danh sách phim
│   ├── showtimes.html             # Chọn suất chiếu
│   ├── seat-selection.html        # Chọn ghế
│   └── my-tickets.html            # Vé của tôi
│
├── assets/                         # Tài nguyên tĩnh
│   ├── css/                       # Styles
│   │   ├── common.css            # Styles chung
│   │   ├── home.css              # Trang chủ
│   │   ├── auth.css              # Đăng nhập/Đăng ký
│   │   ├── movies.css            # Danh sách phim
│   │   ├── showtimes.css         # Suất chiếu
│   │   ├── seat-selection.css    # Chọn ghế
│   │   └── tickets.css           # Vé của tôi
│   │
│   └── js/                        # JavaScript
│       ├── config.js             # Cấu hình API
│       │
│       ├── utils/                # Tiện ích
│       │   ├── storage.js        # LocalStorage wrapper
│       │   ├── auth.js           # Authentication helpers
│       │   └── api.js            # API client
│       │
│       ├── components/           # Components tái sử dụng
│       │   └── header.js         # Header component
│       │
│       └── modules/              # Module theo nghiệp vụ
│           ├── auth/             # Xác thực
│           │   ├── login.js
│           │   └── register.js
│           ├── movies/           # Phim
│           │   └── movies-list.js
│           ├── showtimes/        # Suất chiếu
│           │   └── showtimes-list.js
│           ├── booking/          # Đặt vé
│           │   └── seat-selection.js
│           └── tickets/          # Vé
│               └── my-tickets.js
│
└── README.md                      # Tài liệu này
```

## 🎯 Tính năng

### 1. **Authentication (Xác thực)**

- ✅ Đăng ký tài khoản mới
- ✅ Đăng nhập
- ✅ Đăng xuất
- ✅ JWT token management
- ✅ Auto redirect khi chưa đăng nhập

### 2. **Movies (Phim)**

- ✅ Xem danh sách phim đang chiếu
- ✅ Xem thông tin chi tiết phim
- ✅ Lọc phim theo thể loại
- ✅ Hiển thị đánh giá và thời lượng

### 3. **Showtimes (Suất chiếu)**

- ✅ Xem danh sách suất chiếu theo phim
- ✅ Nhóm suất chiếu theo ngày
- ✅ Hiển thị thông tin rạp và phòng chiếu

### 4. **Booking (Đặt vé)**

- ✅ Chọn ghế trực quan
- ✅ Hiển thị trạng thái ghế (trống/đã đặt/đang chọn)
- ✅ Chọn nhiều ghế cùng lúc
- ✅ Xử lý concurrency (conflict detection)
- ✅ Booking confirmation

### 5. **My Tickets (Vé của tôi)**

- ✅ Xem danh sách vé đã đặt
- ✅ Hiển thị thông tin chi tiết vé
- ✅ Phân loại vé sắp chiếu/đã chiếu
- ✅ Sắp xếp theo thời gian

## 🚀 Cài đặt và Chạy

### 1. Cấu hình API Key

Mở file `assets/js/config.js` và cập nhật API Key:

```javascript
const API_CONFIG = {
  BASE_URL: "http://localhost:8080/api",
  API_KEY: "YOUR_API_KEY_HERE", // ⚠️ Thay bằng API key thực tế
  // ...
};
```

**Lấy API Key từ backend:**

```bash
# Chạy backend trước
cd ../backend

# Tạo API key cho web client
go run cmd/bff/main.go web 100 60

# Copy API key được generate và paste vào config.js
```

### 2. Chạy Frontend

**Option 1: Sử dụng Live Server (VS Code)**

1. Install extension "Live Server"
2. Right-click vào `index.html`
3. Chọn "Open with Live Server"
4. Browser sẽ tự động mở `http://localhost:5500`

**Option 2: Sử dụng Python HTTP Server**

```bash
# Python 3
python -m http.server 8000

# Truy cập: http://localhost:8000
```

**Option 3: Sử dụng Node.js http-server**

```bash
# Install http-server globally
npm install -g http-server

# Run server
http-server -p 8000

# Truy cập: http://localhost:8000
```

### 3. Đảm bảo Backend đang chạy

```bash
cd ../backend

# Run BFF server (port 8080)
go run cmd/bff/main.go

# Run Core server (port 8081)
go run cmd/core/main.go
```

## 🔧 Cấu hình

### API Endpoints (config.js)

```javascript
ENDPOINTS: {
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    MOVIES: '/movies',
    SHOWS: '/shows',
    SEATS: '/seats',
    BOOK: '/book',
    MY_TICKETS: '/tickets'
}
```

### Storage Keys

```javascript
STORAGE_KEYS: {
    TOKEN: 'cinema_auth_token',
    USER_ID: 'cinema_user_id',
    USER_EMAIL: 'cinema_user_email',
    USER_NAME: 'cinema_user_name'
}
```

### Seat Status

```javascript
SEAT_STATUS: {
    AVAILABLE: 'available',
    BOOKED: 'booked',
    SELECTED: 'selected'  // Client-side only
}
```

## 📱 Responsive Design

Giao diện được thiết kế responsive cho:

- 🖥️ Desktop (> 768px)
- 📱 Tablet (768px)
- 📱 Mobile (< 768px)

## 🎨 Design System

### Colors

```css
--primary-color: #e50914; /* Netflix Red */
--secondary-color: #564d4d;
--background-dark: #141414;
--background-light: #1a1a1a;
--background-card: #2a2a2a;
--text-primary: #ffffff;
--text-secondary: #b3b3b3;
```

### Typography

- Font family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto'
- Base font size: 16px
- Line height: 1.6

### Spacing

```css
--spacing-xs: 0.5rem;
--spacing-sm: 1rem;
--spacing-md: 1.5rem;
--spacing-lg: 2rem;
--spacing-xl: 3rem;
```

## 🔒 Security Features

### 1. **Triple-Layer Security**

- API Key authentication
- JWT token authentication
- Rate limiting

### 2. **XSS Protection**

- Input validation
- Output encoding
- CSP headers (nếu deploy production)

### 3. **CSRF Protection**

- Token-based authentication
- SameSite cookies (nếu dùng cookies)

## 🧪 Testing

### Manual Testing Checklist

**Authentication:**

- [ ] Đăng ký tài khoản mới
- [ ] Đăng nhập với tài khoản hợp lệ
- [ ] Đăng nhập với sai mật khẩu
- [ ] Đăng xuất

**Movies:**

- [ ] Load danh sách phim
- [ ] Click vào phim để xem suất chiếu

**Booking:**

- [ ] Chọn suất chiếu
- [ ] Chọn ghế
- [ ] Đặt vé thành công
- [ ] Xử lý conflict khi ghế đã được đặt

**My Tickets:**

- [ ] Xem danh sách vé
- [ ] Phân loại vé sắp chiếu/đã chiếu

### Concurrency Testing

**Test 2 users đặt cùng ghế:**

1. Mở 2 browser khác nhau (Chrome + Firefox)
2. Đăng nhập 2 tài khoản khác nhau
3. Cùng chọn 1 suất chiếu
4. Cùng chọn 1 ghế
5. Click "Đặt vé" đồng thời

**Kết quả mong đợi:**

- ✅ User 1: Đặt vé thành công
- ❌ User 2: Lỗi "Ghế đã được đặt"

## 📖 API Integration

### Request Flow

```
1. User Action (Click button)
   ↓
2. Validation (Client-side)
   ↓
3. API Call (api.js)
   ↓
4. BFF Server (port 8080)
   ├─ API Key check
   ├─ Rate limit check
   └─ JWT validation
   ↓
5. Core Server (port 8081)
   ├─ Business logic
   └─ Database
   ↓
6. Response
   ↓
7. Update UI
```

### Example: Booking API

```javascript
// Client code
const response = await API.post(API_CONFIG.ENDPOINTS.BOOK, {
    seats: [1, 2, 3]
}, true); // Include auth token

// HTTP Request
POST http://localhost:8080/api/book
Headers:
  Content-Type: application/json
  X-API-Key: web_xxx...
  Authorization: Bearer eyJhbGc...
Body:
  { "seats": [1, 2, 3] }
```

## 🐛 Troubleshooting

### Lỗi CORS

**Triệu chứng:** `Access to fetch at 'http://localhost:8080' from origin 'http://localhost:5500' has been blocked by CORS policy`

**Giải pháp:**

1. Kiểm tra backend đã enable CORS middleware
2. Đảm bảo `Access-Control-Allow-Origin` header được set
3. Check `Access-Control-Allow-Headers` includes `X-API-Key`

### Lỗi 401 Unauthorized

**Triệu chứng:** Luôn bị redirect về trang login

**Giải pháp:**

1. Check API Key trong `config.js`
2. Check JWT token trong localStorage
3. Kiểm tra token expiration
4. Clear localStorage và login lại

### Lỗi 429 Too Many Requests

**Triệu chứng:** `Bạn đã gửi quá nhiều yêu cầu`

**Giải pháp:**

1. Đợi 60 giây (rate limit window)
2. Hoặc tăng rate limit trong backend
3. Check Redis đang chạy

### UI không load

**Giải pháp:**

1. Check browser console (F12) for errors
2. Verify API endpoints trong `config.js`
3. Check backend servers đang chạy
4. Clear browser cache

## 🚀 Deployment

### Production Checklist

- [ ] Update API_CONFIG.BASE_URL to production URL
- [ ] Update API_KEY with production key
- [ ] Minify CSS/JS files
- [ ] Enable CSP headers
- [ ] Add analytics (Google Analytics, etc.)
- [ ] Test on multiple browsers
- [ ] Test on mobile devices
- [ ] Setup CDN for static assets
- [ ] Configure caching headers

### Build for Production

```bash
# Minify CSS (using clean-css-cli)
npm install -g clean-css-cli
cleancss -o assets/css/styles.min.css assets/css/*.css

# Minify JS (using uglify-js)
npm install -g uglify-js
uglifyjs assets/js/**/*.js -o assets/js/bundle.min.js
```

## 📚 Technologies Used

- **HTML5** - Markup
- **CSS3** - Styling (Flexbox, Grid, Custom Properties)
- **JavaScript (ES6+)** - Logic
  - Async/Await
  - Fetch API
  - LocalStorage API
  - URLSearchParams
- **No frameworks/libraries** - Vanilla JS only

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This is an educational project for DBMS course.

## 👥 Credits

Developed by: [Your Name]
Course: Database Management Systems
Year: 2026

---

**Happy Coding! 🎬🍿**
