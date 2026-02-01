import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const API_KEY = import.meta.env.VITE_API_KEY;

const axiosClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

axiosClient.interceptors.request.use(
  (config) => {
    if (API_KEY) {
      config.headers["X-API-Key"] = API_KEY;
    }

    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

axiosClient.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Chỉ logout và redirect nếu:
      // 1. Đã có token (user đã login trước đó)
      // 2. Không phải request đăng nhập/đăng ký
      const isLoginRequest = error.config?.url?.includes("/auth/login");
      const isRegisterRequest = error.config?.url?.includes("/auth/register");
      const hasToken = localStorage.getItem("token");

      if (hasToken && !isLoginRequest && !isRegisterRequest) {
        // Token hết hạn hoặc không hợp lệ -> logout
        localStorage.removeItem("userID");
        localStorage.removeItem("email");
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        window.location.href = "/login";
      }

      // Nếu là login/register request, để error được throw về component
    }

    // Handle rate limiting
    if (error.response?.status === 429) {
      alert("⚠️ Bạn đã gửi quá nhiều request. Vui lòng thử lại sau.");
    }

    return Promise.reject(error);
  },
);

export default axiosClient;
