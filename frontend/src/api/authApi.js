import axiosClient from "./axiosConfig";

const authApi = {
  login: (credentials) => {
    // Đăng nhập bằng email
    return axiosClient.post("/auth/login", {
      email: credentials.email,
      password: credentials.password,
    });
  },

  register: (full_name, email, password) => {
    // Đăng ký tài khoản mới
    return axiosClient.post("/auth/register", {
      full_name,
      email,
      password,
    });
  },

  logout: () => {
    // Xóa tất cả thông tin đăng nhập
    localStorage.removeItem("userID");
    localStorage.removeItem("email");
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  },

  isAuthenticated: () => {
    // Kiểm tra có token hoặc userID
    return !!(localStorage.getItem("token") || localStorage.getItem("userID"));
  },

  getStoredUser: () => {
    try {
      // Ưu tiên lấy từ user object đã lưu
      const savedUser = localStorage.getItem("user");
      if (savedUser) {
        return JSON.parse(savedUser);
      }

      // Fallback về cách cũ
      const userID = localStorage.getItem("userID");
      const email = localStorage.getItem("email");
      if (userID && email) {
        return {
          user_id: parseInt(userID),
          userID: parseInt(userID),
          email: email,
        };
      }
    } catch (error) {
      console.error("Error parsing stored user:", error);
    }

    // Nếu thiếu thông tin thì return null (user chưa login)
    return null;
  },
};

export default authApi;
