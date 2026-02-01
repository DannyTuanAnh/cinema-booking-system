import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { validatePassword } from "../utils/errorHandler";
import Toast from "../components/Toast";
import "../styles/authPage.css";

const RegisterPage = () => {
  const [full_name, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [toast, setToast] = useState(null);
  const [loading, setLoading] = useState(false);

  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setToast(null);

    // Validation
    if (!full_name || !email || !password || !confirmPassword) {
      setToast({ message: "Vui lòng nhập đầy đủ thông tin", type: "error" });
      return;
    }

    if (password !== confirmPassword) {
      setToast({ message: "Mật khẩu xác nhận không khớp", type: "error" });
      return;
    }

    // Validate strong password
    const passwordValidation = validatePassword(password);
    if (!passwordValidation.valid) {
      setToast({ message: passwordValidation.message, type: "error" });
      return;
    }

    setLoading(true);

    try {
      const result = await register(full_name, email, password);

      if (result.success) {
        setToast({
          message: "Đăng ký thành công! Đang chuyển tới trang đăng nhập...",
          type: "success",
        });
        setTimeout(() => {
          navigate("/login");
        }, 1500);
      } else {
        setToast({ message: result.message, type: "error" });
      }
    } catch (err) {
      setToast({
        message: "Đăng ký thất bại. Vui lòng thử lại.",
        type: "error",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}

      <div className="auth-card">
        <h2>Đăng ký</h2>
        <p className="auth-subtitle">Tạo tài khoản mới</p>

        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="name">Họ và tên</label>
            <input
              type="text"
              id="name"
              className="form-input"
              placeholder="Nguyễn Văn A"
              value={full_name}
              onChange={(e) => setFullName(e.target.value)}
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              className="form-input"
              placeholder="example@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Mật khẩu</label>
            <input
              type="password"
              id="password"
              className="form-input"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={loading}
            />
            <small className="form-hint">
              Tối thiểu 8 ký tự, có chữ hoa, chữ thường, số và ký tự đặc biệt
              (@$!%*?&)
            </small>
          </div>

          <div className="form-group">
            <label htmlFor="confirmPassword">Xác nhận mật khẩu</label>
            <input
              type="password"
              id="confirmPassword"
              className="form-input"
              placeholder="••••••••"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              disabled={loading}
            />
          </div>

          <button
            type="submit"
            className="btn btn-primary btn-block"
            disabled={loading}
          >
            {loading ? "Đang đăng ký..." : "Đăng ký"}
          </button>
        </form>

        <div className="auth-footer">
          <p>
            Đã có tài khoản?{" "}
            <Link to="/login" className="auth-link">
              Đăng nhập ngay
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
