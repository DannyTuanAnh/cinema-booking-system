import { Link } from "react-router-dom";
import "../styles/homePage.css";

const HomePage = () => {
  return (
    <div className="home-page">
      <section className="hero">
        <div className="container">
          <h1 className="hero-title">Trải nghiệm điện ảnh đỉnh cao</h1>
          <p className="hero-subtitle">
            Đặt vé xem phim online - Nhanh chóng, Tiện lợi, An toàn
          </p>
          <div className="hero-actions">
            <Link to="/movies" className="btn btn-primary btn-lg">
              Xem phim ngay
            </Link>
          </div>
        </div>
      </section>

      <section className="features">
        <div className="container">
          <h2 className="section-title">Tại sao chọn chúng tôi?</h2>
          <div className="features-grid">
            <div className="feature-card">
              <div className="feature-icon">🎬</div>
              <h3>Phim mới nhất</h3>
              <p>Cập nhật liên tục các bộ phim hot nhất</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">💺</div>
              <h3>Chọn ghế dễ dàng</h3>
              <p>Giao diện trực quan, chọn ghế nhanh chóng</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">🔒</div>
              <h3>An toàn bảo mật</h3>
              <p>Thông tin được bảo mật tuyệt đối</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">⚡</div>
              <h3>Thanh toán nhanh</h3>
              <p>Đặt vé và thanh toán chỉ trong vài phút</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default HomePage;
