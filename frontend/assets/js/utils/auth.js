/**
 * ============================================
 * AUTHENTICATION UTILITIES
 * Xử lý xác thực và phân quyền
 * ============================================
 */

const Auth = {
  /**
   * Kiểm tra trạng thái đăng nhập
   */
  isLoggedIn() {
    return Storage.isAuthenticated();
  },

  /**
   * Lấy thông tin user hiện tại
   */
  getCurrentUser() {
    if (!this.isLoggedIn()) {
      return null;
    }
    return Storage.getUserInfo();
  },

  /**
   * Đăng xuất
   */
  logout() {
    Storage.clearAll();
    window.location.href = "/frontend/pages/login.html";
  },

  /**
   * Redirect nếu chưa đăng nhập
   */
  requireAuth() {
    if (!this.isLoggedIn()) {
      window.location.href = "/frontend/pages/login.html";
      return false;
    }
    return true;
  },

  /**
   * Redirect nếu đã đăng nhập
   */
  redirectIfAuthenticated() {
    if (this.isLoggedIn()) {
      window.location.href = "/index.html";
      return true;
    }
    return false;
  },

  /**
   * Lấy Authorization header
   */
  getAuthHeader() {
    const token = Storage.getAuthToken();
    if (!token) {
      return {};
    }
    return {
      Authorization: `Bearer ${token}`,
    };
  },
};

/**
 * Hàm logout global (được gọi từ HTML)
 */
function logout() {
  if (confirm("Bạn có chắc muốn đăng xuất?")) {
    Auth.logout();
  }
}

// Make Auth globally available
window.Auth = Auth;
