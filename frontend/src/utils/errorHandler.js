/**
 * Parse error message from API response
 * @param {Object} error - Axios error object
 * @param {string} defaultMessage - Default error message
 * @returns {string} - Formatted error message
 */
export const parseErrorMessage = (error, defaultMessage = "Có lỗi xảy ra") => {
  if (!error.response?.data?.error) {
    return defaultMessage;
  }

  const errorData = error.response.data.error;

  // Nếu error là string
  if (typeof errorData === "string") {
    return errorData;
  }

  // Nếu error là object (validation errors từ backend)
  if (typeof errorData === "object" && !Array.isArray(errorData)) {
    // Lấy tất cả error messages và join thành 1 string
    const messages = Object.values(errorData);
    if (messages.length > 0) {
      // Nếu có nhiều lỗi, hiển thị lỗi đầu tiên
      return messages[0];
    }
  }

  // Nếu error là array
  if (Array.isArray(errorData)) {
    return errorData[0] || defaultMessage;
  }

  return defaultMessage;
};

/**
 * Unwrap API response
 * Backend trả về: {response: data} hoặc {data: data}
 * Function này extract ra data thật
 * @param {Object} data - Response data từ axios
 * @returns {*} - Unwrapped data
 */
export const unwrapResponse = (data) => {
  // Nếu có key "response", unwrap nó
  if (data && typeof data === "object" && "response" in data) {
    return data.response;
  }
  // Nếu không có, return data gốc
  return data;
};

/**
 * Validate strong password
 * @param {string} password
 * @returns {Object} - { valid: boolean, message: string }
 */
export const validatePassword = (password) => {
  if (!password) {
    return {
      valid: false,
      message: "Mật khẩu không được để trống",
    };
  }

  if (password.length < 8) {
    return {
      valid: false,
      message: "Mật khẩu phải có ít nhất 8 ký tự",
    };
  }

  const hasLowercase = /[a-z]/.test(password);
  const hasUppercase = /[A-Z]/.test(password);
  const hasNumber = /[0-9]/.test(password);
  const hasSpecial = /[@$!%*?&]/.test(password);

  if (!hasLowercase) {
    return {
      valid: false,
      message: "Mật khẩu phải có ít nhất 1 chữ thường (a-z)",
    };
  }

  if (!hasUppercase) {
    return {
      valid: false,
      message: "Mật khẩu phải có ít nhất 1 chữ hoa (A-Z)",
    };
  }

  if (!hasNumber) {
    return {
      valid: false,
      message: "Mật khẩu phải có ít nhất 1 số (0-9)",
    };
  }

  if (!hasSpecial) {
    return {
      valid: false,
      message: "Mật khẩu phải có ít nhất 1 ký tự đặc biệt (@$!%*?&)",
    };
  }

  return {
    valid: true,
    message: "",
  };
};

/**
 * Validate email format
 * @param {string} email
 * @returns {boolean}
 */
export const validateEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};
