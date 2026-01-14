/**
 * ============================================
 * LOGIN MODULE
 * Xử lý đăng nhập user
 * ============================================
 */

// Redirect if already logged in
Auth.redirectIfAuthenticated();

// Form elements
const loginForm = document.getElementById("loginForm");
const emailInput = document.getElementById("email");
const passwordInput = document.getElementById("password");
const loginBtn = document.getElementById("loginBtn");
const alertContainer = document.getElementById("alertContainer");

/**
 * Hiển thị thông báo
 */
function showAlert(message, type = "error") {
  alertContainer.innerHTML = `
        <div class="alert alert-${type}">
            ${message}
        </div>
    `;

  // Auto hide after 5 seconds
  setTimeout(() => {
    alertContainer.innerHTML = "";
  }, 5000);
}

/**
 * Validate form
 */
function validateForm() {
  const email = emailInput.value.trim();
  const password = passwordInput.value.trim();

  if (!email) {
    showAlert("Vui lòng nhập email");
    emailInput.focus();
    return false;
  }

  if (!VALIDATION.EMAIL_REGEX.test(email)) {
    showAlert("Email không hợp lệ");
    emailInput.focus();
    return false;
  }

  if (!password) {
    showAlert("Vui lòng nhập mật khẩu");
    passwordInput.focus();
    return false;
  }

  if (password.length < VALIDATION.PASSWORD_MIN_LENGTH) {
    showAlert(
      `Mật khẩu phải có ít nhất ${VALIDATION.PASSWORD_MIN_LENGTH} ký tự`
    );
    passwordInput.focus();
    return false;
  }

  return true;
}

/**
 * Xử lý đăng nhập
 */
async function handleLogin(e) {
  e.preventDefault();

  if (!validateForm()) {
    return;
  }

  // Disable button
  loginBtn.disabled = true;
  loginBtn.textContent = "Đang đăng nhập...";

  try {
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();

    // Call login API
    const response = await API.post(API_CONFIG.ENDPOINTS.LOGIN, {
      email,
      password,
    });

    // Extract response data (BFF wraps it in "response" field)
    const data = response.response || response;

    // Save auth data
    Storage.setAuthToken(data.access_token);
    Storage.setUserInfo(data.user_id, data.email, data.name);

    // Show success message
    showAlert(SUCCESS_MESSAGES.LOGIN, "success");

    // Redirect after short delay
    setTimeout(() => {
      window.location.href = "../index.html";
    }, 1000);
  } catch (error) {
    console.error("Login error:", error);
    showAlert(error.message);

    // Re-enable button
    loginBtn.disabled = false;
    loginBtn.textContent = "Đăng nhập";
  }
}

// Event listeners
loginForm.addEventListener("submit", handleLogin);

// Auto focus email input
emailInput.focus();

// Handle Enter key on password field
passwordInput.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    handleLogin(e);
  }
});
