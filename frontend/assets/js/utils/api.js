/**
 * ============================================
 * API CLIENT
 * HTTP request wrapper với error handling
 * ============================================
 */

const API = {
  /**
   * Tạo headers cho request
   */
  _getHeaders(includeAuth = false) {
    const headers = {
      "Content-Type": "application/json",
      "X-API-Key": API_CONFIG.API_KEY,
    };

    if (includeAuth) {
      const authHeader = Auth.getAuthHeader();
      Object.assign(headers, authHeader);
    }

    return headers;
  },

  /**
   * Xử lý response
   */
  async _handleResponse(response, retryCallback = null) {
    const contentType = response.headers.get("content-type");
    const isJSON = contentType && contentType.includes("application/json");

    if (!response.ok) {
      let errorMessage = ERROR_MESSAGES.SERVER_ERROR;

      if (isJSON) {
        const errorData = await response.json();
        errorMessage = errorData.error || errorData.message || errorMessage;

        // Handle validation errors
        if (typeof errorMessage === "object") {
          const errors = Object.values(errorMessage).join(", ");
          errorMessage = errors || ERROR_MESSAGES.VALIDATION_ERROR;
        }
      }

      // Handle specific status codes
      switch (response.status) {
        case 401:
          // Nếu có callback để retry (sau khi refresh token)
          if (
            retryCallback &&
            typeof Auth !== "undefined" &&
            typeof Auth.refreshTokenOrLogout === "function"
          ) {
            console.log("🔄 Received 401, attempting token refresh...");
            const refreshed = await Auth.refreshTokenOrLogout();
            if (refreshed) {
              // Token đã được refresh, retry request
              console.log("✅ Token refreshed, retrying request...");
              return await retryCallback();
            }
            // Nếu refresh thất bại, Auth.refreshTokenOrLogout đã redirect về login
            return;
          }

          errorMessage = ERROR_MESSAGES.UNAUTHORIZED;
          // Fallback: Auto logout on 401 nếu không có Auth.refreshTokenOrLogout
          Storage.clearAll();
          if (!window.location.pathname.includes("login")) {
            window.location.href = "/pages/login.html";
          }
          break;
        case 403:
          errorMessage = ERROR_MESSAGES.FORBIDDEN;
          break;
        case 404:
          errorMessage = ERROR_MESSAGES.NOT_FOUND;
          break;
        case 429:
          errorMessage = ERROR_MESSAGES.RATE_LIMIT;
          break;
        case 409:
          errorMessage = ERROR_MESSAGES.BOOKING_CONFLICT;
          break;
      }

      const error = new Error(errorMessage);
      error.status = response.status;
      throw error;
    }

    if (isJSON) {
      return await response.json();
    }

    return await response.text();
  },

  /**
   * GET request
   */
  async get(endpoint, includeAuth = false) {
    const makeRequest = async () => {
      const response = await fetch(`${API_CONFIG.BASE_URL}${endpoint}`, {
        method: "GET",
        headers: this._getHeaders(includeAuth),
        credentials: includeAuth ? "include" : "same-origin",
        timeout: API_CONFIG.TIMEOUT,
      });

      // Pass retry callback cho _handleResponse
      return await this._handleResponse(
        response,
        includeAuth ? makeRequest : null
      );
    };

    try {
      return await makeRequest();
    } catch (error) {
      if (error.name === "TypeError" || error.message === "Failed to fetch") {
        throw new Error(ERROR_MESSAGES.NETWORK_ERROR);
      }
      throw error;
    }
  },

  /**
   * POST request
   */
  async post(endpoint, data, includeAuth = false) {
    const makeRequest = async () => {
      const response = await fetch(`${API_CONFIG.BASE_URL}${endpoint}`, {
        method: "POST",
        headers: this._getHeaders(includeAuth),
        body: JSON.stringify(data),
        credentials: includeAuth ? "include" : "same-origin",
        timeout: API_CONFIG.TIMEOUT,
      });

      // Pass retry callback cho _handleResponse
      return await this._handleResponse(
        response,
        includeAuth ? makeRequest : null
      );
    };

    try {
      return await makeRequest();
    } catch (error) {
      if (error.name === "TypeError" || error.message === "Failed to fetch") {
        throw new Error(ERROR_MESSAGES.NETWORK_ERROR);
      }
      throw error;
    }
  },

  /**
   * PUT request
   */
  async put(endpoint, data, includeAuth = false) {
    try {
      const response = await fetch(`${API_CONFIG.BASE_URL}${endpoint}`, {
        method: "PUT",
        headers: this._getHeaders(includeAuth),
        body: JSON.stringify(data),
        timeout: API_CONFIG.TIMEOUT,
      });

      return await this._handleResponse(response);
    } catch (error) {
      if (error.name === "TypeError" || error.message === "Failed to fetch") {
        throw new Error(ERROR_MESSAGES.NETWORK_ERROR);
      }
      throw error;
    }
  },

  /**
   * DELETE request
   */
  async delete(endpoint, includeAuth = false) {
    try {
      const response = await fetch(`${API_CONFIG.BASE_URL}${endpoint}`, {
        method: "DELETE",
        headers: this._getHeaders(includeAuth),
        timeout: API_CONFIG.TIMEOUT,
      });

      return await this._handleResponse(response);
    } catch (error) {
      if (error.name === "TypeError" || error.message === "Failed to fetch") {
        throw new Error(ERROR_MESSAGES.NETWORK_ERROR);
      }
      throw error;
    }
  },
};

// Make API globally available
window.API = API;
