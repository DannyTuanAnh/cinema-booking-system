import { Outlet } from "react-router-dom";
import "../styles/authLayout.css";

const AuthLayout = () => {
  return (
    <div className="auth-layout">
      <div className="auth-container">
        <div className="auth-brand">
          <h1>🎬 Cinema Booking</h1>
          <p>Đặt vé xem phim dễ dàng, nhanh chóng</p>
        </div>
        <div className="auth-content">
          <Outlet />
        </div>
      </div>
    </div>
  );
};

export default AuthLayout;
