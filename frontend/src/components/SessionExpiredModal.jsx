import { useEffect } from "react";
import "../styles/sessionExpiredModal.css";

const SessionExpiredModal = ({ onClose }) => {
  useEffect(() => {
    // Auto close after 3 seconds
    const timer = setTimeout(() => {
      onClose();
    }, 3000);

    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <div className="session-expired-overlay">
      <div className="session-expired-modal">
        <div className="session-expired-icon">
          <svg
            width="64"
            height="64"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
          >
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
        </div>

        <h2 className="session-expired-title">Phiên đăng nhập đã hết hạn</h2>

        <p className="session-expired-message">
          Vui lòng đăng nhập lại để tiếp tục sử dụng dịch vụ
        </p>

        <div className="session-expired-footer">
          <button className="btn btn-primary" onClick={onClose}>
            Đăng nhập lại
          </button>
        </div>

        <div className="session-expired-timer">
          Tự động chuyển hướng sau 3 giây...
        </div>
      </div>
    </div>
  );
};

export default SessionExpiredModal;
