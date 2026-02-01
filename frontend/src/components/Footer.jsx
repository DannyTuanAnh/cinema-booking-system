import "../styles/footer.css";

const Footer = () => {
  return (
    <footer className="footer">
      <div className="container">
        <div className="footer-content">
          <div className="footer-section">
            <h3>Cinema Booking</h3>
            <p>Hệ thống đặt vé xem phim trực tuyến</p>
          </div>
          <div className="footer-section">
            <h4>Liên hệ</h4>
            <p>Email: support@cinema.com</p>
            <p>Hotline: 1900-xxxx</p>
          </div>
          <div className="footer-section">
            <h4>Theo dõi</h4>
            <p>Facebook | Instagram | Twitter</p>
          </div>
        </div>
        <div className="footer-bottom">
          <p>&copy; 2026 Cinema Booking System. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
