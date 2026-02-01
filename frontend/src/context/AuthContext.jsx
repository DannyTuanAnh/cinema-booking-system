import { createContext, useContext, useState, useEffect } from "react";
import authApi from "../api/authApi";
import { parseErrorMessage } from "../utils/errorHandler";

const AuthContext = createContext(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Check if user is logged in on mount
  useEffect(() => {
    const storedUser = authApi.getStoredUser();
    if (storedUser) {
      setUser(storedUser);
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    try {
      const data = await authApi.login({ email, password });

      // Backend trả về: {response: {token, user_id, email, name}}
      const responseData = data.response || data;

      // Save to localStorage
      localStorage.setItem("token", responseData.token);
      localStorage.setItem("userID", responseData.user_id);
      localStorage.setItem("email", responseData.email);
      localStorage.setItem(
        "user",
        JSON.stringify({
          user_id: responseData.user_id,
          userID: responseData.user_id,
          email: responseData.email,
          name: responseData.name,
        }),
      );

      // Update state
      setUser({
        user_id: responseData.user_id,
        userID: responseData.user_id,
        email: responseData.email,
        name: responseData.name,
      });

      return { success: true };
    } catch (error) {
      console.error("Login error:", error);

      return {
        success: false,
        message: parseErrorMessage(error, "Đăng nhập thất bại"),
      };
    }
  };

  const register = async (name, email, password) => {
    try {
      await authApi.register(name, email, password);
      return { success: true };
    } catch (error) {
      console.error("Register error:", error);

      return {
        success: false,
        message: parseErrorMessage(error, "Đăng ký thất bại"),
      };
    }
  };

  const logout = () => {
    authApi.logout();
    setUser(null);
  };

  const isAuthenticated = () => {
    return authApi.isAuthenticated();
  };

  const value = {
    user,
    loading,
    login,
    register,
    logout,
    isAuthenticated,
  };

  return (
    <AuthContext.Provider value={value}>
      {!loading && children}
    </AuthContext.Provider>
  );
};
