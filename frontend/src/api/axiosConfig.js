import axios from "axios";
import { sessionExpiredEmitter } from "../utils/sessionEmitter";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const API_KEY = import.meta.env.VITE_API_KEY;

// Flag để tránh nhiều refresh request đồng thời
let isRefreshing = false;
// Queue các request đang chờ token mới
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

const axiosClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true, // Cho phép gửi cookies
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
  async (error) => {
    const originalRequest = error.config;

    // Handle 401 - Token hết hạn
    if (error.response?.status === 401 && !originalRequest._retry) {
      const isLoginRequest = originalRequest?.url?.includes("/auth/login");
      const isRegisterRequest =
        originalRequest?.url?.includes("/auth/register");
      const isRefreshRequest = originalRequest?.url?.includes("/auth/refresh");

      // Nếu là login/register request, throw error về component
      if (isLoginRequest || isRegisterRequest) {
        return Promise.reject(error);
      }

      // Nếu refresh request bị 401, nghĩa là refresh token cũng hết hạn
      if (isRefreshRequest) {
        console.log("🔴 Refresh token expired - showing session expired modal");
        sessionExpiredEmitter.emit();
        return Promise.reject(error);
      }

      // Nếu đang refresh, đưa request vào queue
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return axiosClient(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      // Bắt đầu refresh token
      originalRequest._retry = true;
      isRefreshing = true;

      console.log("🔄 Access token expired - refreshing...");

      return new Promise((resolve, reject) => {
        axios
          .post(
            `${API_BASE_URL}/auth/refresh`,
            {},
            {
              headers: {
                "X-API-Key": API_KEY,
              },
              withCredentials: true, // Gửi cookie refresh_token
            },
          )
          .then((response) => {
            // Backend trả về: {response: {access_token: "..."}}
            const responseData = response.data.response || response.data;
            const newAccessToken = responseData.access_token;

            console.log("✅ Token refreshed successfully");

            // Lưu token mới
            localStorage.setItem("token", newAccessToken);

            // Update header của request gốc
            originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;

            // Process queue
            processQueue(null, newAccessToken);

            // Retry request gốc
            resolve(axiosClient(originalRequest));
          })
          .catch((err) => {
            console.error("❌ Refresh token failed:", err);

            // Refresh token hết hạn hoặc invalid
            processQueue(err, null);

            // Show session expired modal
            sessionExpiredEmitter.emit();

            reject(err);
          })
          .finally(() => {
            isRefreshing = false;
          });
      });
    }

    // Handle rate limiting
    if (error.response?.status === 429) {
      alert("⚠️ Bạn đã gửi quá nhiều request. Vui lòng thử lại sau.");
    }

    return Promise.reject(error);
  },
);

export default axiosClient;
