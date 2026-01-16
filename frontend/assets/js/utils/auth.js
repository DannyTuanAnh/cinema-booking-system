/**
 * ============================================
 * AUTHENTICATION UTILITIES
 * Xử lý xác thực và phân quyền
 * ============================================
 */

const Auth = {
  /**
   * Flag to prevent multiple refresh attempts
   */
  _isRefreshing: false,
  _refreshPromise: null,
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

  /**
   * Refresh access token bằng refresh token từ HttpOnly cookie
   * @returns {Promise<boolean>} - true nếu refresh thành công, false nếu thất bại
   */
  async refreshAccessToken() {
    // Prevent multiple simultaneous refresh attempts
    if (this._isRefreshing) {
      return this._refreshPromise;
    }

    this._isRefreshing = true;

    this._refreshPromise = (async () => {
      try {
        console.log("🔄 Attempting to refresh access token...");

        // Call refresh endpoint
        // Backend sẽ tự động đọc refresh_token từ HttpOnly cookie (do backend set khi login)
        // Frontend KHÔNG lưu refresh_token ở localStorage vì lý do bảo mật
        const response = await fetch(
          `${API_CONFIG.BASE_URL}${API_CONFIG.ENDPOINTS.REFRESH}`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              "X-API-Key": API_CONFIG.API_KEY,
            },
            credentials: "include", // Critical: Gửi kèm HttpOnly cookies chứa refresh_token
          }
        );

        if (!response.ok) {
          // Nếu refresh cũng thất bại (401), nghĩa là refresh token hết hạn
          if (response.status === 401) {
            console.error("❌ Refresh token expired");
            return false;
          }

          const errorData = await response.json().catch(() => ({}));
          throw new Error(errorData.error || "Failed to refresh token");
        }

        const data = await response.json();
        const result = data.response || data;

        // Lưu access token mới
        if (result.access_token) {
          Storage.setAuthToken(result.access_token);
          console.log("✅ Access token refreshed successfully");
          return true;
        } else {
          console.error("❌ No access token in refresh response");
          return false;
        }
      } catch (error) {
        console.error("❌ Refresh token error:", error);
        return false;
      } finally {
        this._isRefreshing = false;
        this._refreshPromise = null;
      }
    })();

    return this._refreshPromise;
  },

  /**
   * Thử refresh token hoặc logout nếu thất bại
   * Sử dụng khi nhận 401 từ API
   */
  async refreshTokenOrLogout() {
    const refreshSuccess = await this.refreshAccessToken();

    if (!refreshSuccess) {
      // Refresh thất bại -> logout và redirect về login
      console.log("🚪 Refresh failed, logging out...");
      Storage.clearAll();
      window.location.href = "/public/frontend/pages/login.html";
      return false;
    }

    return true;
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
