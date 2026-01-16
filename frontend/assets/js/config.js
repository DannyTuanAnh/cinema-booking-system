/**
 * ============================================
 * CONFIGURATION FILE
 * Cấu hình API endpoints và constants
 * ============================================
 */

const API_CONFIG = {
  // BFF API Base URL
  BASE_URL: "https://localhost:8443/api",

  // API Key (generated from backend)
  API_KEY: "web_40a58cd8-5182-47d1-9e69-e7f01b07bc9a_crLdf4Wm3Z5bAWeX",

  // API Endpoints
  ENDPOINTS: {
    // Auth
    LOGIN: "/auth/login",
    REGISTER: "/auth/register",
    REFRESH: "/auth/refresh",

    // Movies
    MOVIES: "/movies",

    // Shows
    SHOWS: "/shows",

    // Seats
    SEATS: "/seats",

    // Booking
    BOOK: "/book",

    // Tickets
    MY_TICKETS: "/tickets",
  },

  // Request timeout
  TIMEOUT: 10000, // 10 seconds

  // Rate limiting info
  RATE_LIMIT: {
    MAX_REQUESTS: 100,
    WINDOW_SECONDS: 60,
  },
};

// Storage keys
const STORAGE_KEYS = {
  TOKEN: "cinema_auth_token",
  USER_ID: "cinema_user_id",
  USER_EMAIL: "cinema_user_email",
  USER_NAME: "cinema_user_name",
};

// Seat status
const SEAT_STATUS = {
  AVAILABLE: "available",
  BOOKED: "booked",
  SELECTED: "selected", // Client-side only
};

// Error messages
const ERROR_MESSAGES = {
  NETWORK_ERROR:
    "Không thể kết nối đến server. Vui lòng kiểm tra kết nối mạng.",
  UNAUTHORIZED: "Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.",
  FORBIDDEN: "Bạn không có quyền truy cập.",
  NOT_FOUND: "Không tìm thấy dữ liệu.",
  SERVER_ERROR: "Lỗi server. Vui lòng thử lại sau.",
  RATE_LIMIT: "Bạn đã gửi quá nhiều yêu cầu. Vui lòng thử lại sau.",
  VALIDATION_ERROR: "Dữ liệu không hợp lệ.",
  BOOKING_CONFLICT:
    "Một hoặc nhiều ghế đã được đặt bởi người khác. Vui lòng chọn ghế khác.",
};

// Success messages
const SUCCESS_MESSAGES = {
  LOGIN: "Đăng nhập thành công!",
  REGISTER: "Đăng ký thành công! Vui lòng đăng nhập.",
  BOOKING: "Đặt vé thành công!",
  LOGOUT: "Đăng xuất thành công!",
};

// Validation rules
const VALIDATION = {
  EMAIL_REGEX: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PASSWORD_MIN_LENGTH: 6,
  NAME_MIN_LENGTH: 2,
};

// Make config globally available
window.API_CONFIG = API_CONFIG;
window.STORAGE_KEYS = STORAGE_KEYS;
window.SEAT_STATUS = SEAT_STATUS;
window.ERROR_MESSAGES = ERROR_MESSAGES;
window.SUCCESS_MESSAGES = SUCCESS_MESSAGES;
window.VALIDATION = VALIDATION;
