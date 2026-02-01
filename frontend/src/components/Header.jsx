import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import "../styles/header.css";

const Header = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    if (window.confirm("Bạn có chắc muốn đăng xuất?")) {
      logout();
      navigate("/login");
    }
  };

  return (
    <header className="header">
      <div className="container">
        <Link to="/" className="logo">
          <span className="logo-icon">🎬</span>
          <span className="logo-text">Cinema</span>
        </Link>

        <nav className="nav">
          <Link to="/" className="nav-link">
            Trang chủ
          </Link>
          <Link to="/movies" className="nav-link">
            Phim
          </Link>
          {isAuthenticated() && (
            <Link to="/my-tickets" className="nav-link">
              Vé của tôi
            </Link>
          )}
        </nav>

        <div className="header-actions">
          {isAuthenticated() ? (
            <>
              <span className="user-name">👤 {user?.name}</span>
              <button onClick={handleLogout} className="btn btn-outline">
                Đăng xuất
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="btn btn-outline">
                Đăng nhập
              </Link>
              <Link to="/register" className="btn btn-primary">
                Đăng ký
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;
