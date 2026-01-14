/**
 * ============================================
 * REGISTER MODULE
 * Xử lý đăng ký user mới
 * ============================================
 */

// Redirect if already logged in
Auth.redirectIfAuthenticated();

// Form elements
const registerForm = document.getElementById("registerForm");
const fullNameInput = document.getElementById("fullName");
const emailInput = document.getElementById("email");
const passwordInput = document.getElementById("password");
const confirmPasswordInput = document.getElementById("confirmPassword");
const registerBtn = document.getElementById("registerBtn");
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
  const fullName = fullNameInput.value.trim();
  const email = emailInput.value.trim();
  const password = passwordInput.value.trim();
  const confirmPassword = confirmPasswordInput.value.trim();

  if (!fullName) {
    showAlert("Vui lòng nhập họ và tên");
    fullNameInput.focus();
    return false;
  }

  if (fullName.length < VALIDATION.NAME_MIN_LENGTH) {
    showAlert(`Họ và tên phải có ít nhất ${VALIDATION.NAME_MIN_LENGTH} ký tự`);
    fullNameInput.focus();
    return false;
  }

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

  if (!confirmPassword) {
    showAlert("Vui lòng xác nhận mật khẩu");
    confirmPasswordInput.focus();
    return false;
  }

  if (password !== confirmPassword) {
    showAlert("Mật khẩu xác nhận không khớp");
    confirmPasswordInput.focus();
    return false;
  }

  return true;
}

/**
 * Xử lý đăng ký
 */
async function handleRegister(e) {
  e.preventDefault();

  if (!validateForm()) {
    return;
  }

  // Disable button
  registerBtn.disabled = true;
  registerBtn.textContent = "Đang đăng ký...";

  try {
    const fullName = fullNameInput.value.trim();
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();

    // Call register API
    await API.post(API_CONFIG.ENDPOINTS.REGISTER, {
      full_name: fullName,
      email,
      password,
    });

    // Show success message
    showAlert(SUCCESS_MESSAGES.REGISTER, "success");

    // Redirect to login after short delay
    setTimeout(() => {
      window.location.href = "login.html";
    }, 2000);
  } catch (error) {
    console.error("Register error:", error);
    showAlert(error.message);

    // Re-enable button
    registerBtn.disabled = false;
    registerBtn.textContent = "Đăng ký";
  }
}

// Event listeners
registerForm.addEventListener("submit", handleRegister);

// Auto focus full name input
fullNameInput.focus();
