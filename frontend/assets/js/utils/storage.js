/**
 * ============================================
 * LOCAL STORAGE UTILITIES
 * Quản lý lưu trữ dữ liệu client-side
 * ============================================
 */

const Storage = {
  /**
   * Lưu token xác thực
   */
  setAuthToken(token) {
    localStorage.setItem(STORAGE_KEYS.TOKEN, token);
  },

  /**
   * Lấy token xác thực
   */
  getAuthToken() {
    return localStorage.getItem(STORAGE_KEYS.TOKEN);
  },

  /**
   * Xóa token xác thực
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
  },

  /**
   * Lưu dữ liệu tùy chỉnh
   */
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
