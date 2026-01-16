/**
 * ============================================
 * LOCAL STORAGE UTILITIES
 * Quản lý lưu trữ dữ liệu client-side
 * ============================================
 */

// Storage keys (if defined globally, skip this)
if (typeof STORAGE_KEYS === "undefined") {
  window.STORAGE_KEYS = {
    TOKEN: "cinema_auth_token",
    USER_ID: "cinema_user_id",
    USER_EMAIL: "cinema_user_email",
    USER_NAME: "cinema_user_name",
  };
}

const Storage = {
  /**
   * Lưu access token
   */
  setAuthToken(token) {
    localStorage.setItem(STORAGE_KEYS.TOKEN, token);
  },

  /**
   * Lấy access token
   */
  getAuthToken() {
    return localStorage.getItem(STORAGE_KEYS.TOKEN);
  },

  /**
   * Xóa access token
   */
  removeAuthToken() {
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
  },

  /**
   * Lưu thông tin user
   */
  setUserInfo(userId, email, name) {
    localStorage.setItem(STORAGE_KEYS.USER_ID, userId);
    localStorage.setItem(STORAGE_KEYS.USER_EMAIL, email);
    localStorage.setItem(STORAGE_KEYS.USER_NAME, name);
  },

  /**
   * Lấy thông tin user
   */
  getUserInfo() {
    return {
      userId: localStorage.getItem(STORAGE_KEYS.USER_ID),
      email: localStorage.getItem(STORAGE_KEYS.USER_EMAIL),
      name: localStorage.getItem(STORAGE_KEYS.USER_NAME),
    };
  },

  /**
   * Xóa thông tin user
   */
  removeUserInfo() {
    localStorage.removeItem(STORAGE_KEYS.USER_ID);
    localStorage.removeItem(STORAGE_KEYS.USER_EMAIL);
    localStorage.removeItem(STORAGE_KEYS.USER_NAME);
  },

  /**
   * Kiểm tra user đã đăng nhập
   */
  isAuthenticated() {
    return !!this.getAuthToken();
  },

  /**
   * Xóa toàn bộ dữ liệu
   */
  clearAll() {
    this.removeAuthToken();
    this.removeUserInfo();
  }
  /**
   * Lưu dữ liệu tùy chỉnh
   */,
  set(key, value) {
    try {
      const data = typeof value === "object" ? JSON.stringify(value) : value;
      localStorage.setItem(key, data);
      return true;
    } catch (error) {
      console.error("Storage set error:", error);
      return false;
    }
  },

  /**
   * Lấy dữ liệu tùy chỉnh
   */
  get(key, parseJSON = false) {
    try {
      const data = localStorage.getItem(key);
      if (parseJSON && data) {
        return JSON.parse(data);
      }
      return data;
    } catch (error) {
      console.error("Storage get error:", error);
      return null;
    }
  },

  /**
   * Xóa dữ liệu tùy chỉnh
   */
  remove(key) {
    localStorage.removeItem(key);
  },
};

// Make Storage globally available
window.Storage = Storage;
