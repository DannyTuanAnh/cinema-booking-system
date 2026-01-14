/**
 * ============================================
 * HEADER COMPONENT
 * Quản lý hiển thị header (user info, logout)
 * ============================================
 */

document.addEventListener("DOMContentLoaded", function () {
  updateHeaderUI();
});

/**
 * Cập nhật UI header dựa trên trạng thái đăng nhập
 */
function updateHeaderUI() {
  const userInfo = document.getElementById("userInfo");
  const authButtons = document.getElementById("authButtons");

  if (!userInfo || !authButtons) {
    return; // Not on a page with header
  }

  if (Auth.isLoggedIn()) {
    const user = Auth.getCurrentUser();

    // Show user info
    userInfo.classList.remove("hidden");
    authButtons.classList.add("hidden");

    // Set user name
    const userNameElement = userInfo.querySelector(".user-name");
    if (userNameElement) {
      userNameElement.textContent = user.name || user.email;
    }
  } else {
    // Show auth buttons
    userInfo.classList.add("hidden");
    authButtons.classList.remove("hidden");
  }
}

// Update header when storage changes (e.g., login from another tab)
window.addEventListener("storage", function (e) {
  if (e.key === STORAGE_KEYS.TOKEN || e.key === null) {
    updateHeaderUI();
  }
});
