import "../styles/errorMessage.css";

const ErrorMessage = ({ message, onRetry }) => {
  return (
    <div className="error-container">
      <div className="error-icon">⚠️</div>
      <p className="error-message">{message}</p>
      {onRetry && (
        <button onClick={onRetry} className="btn btn-primary">
          Thử lại
        </button>
      )}
    </div>
  );
};

export default ErrorMessage;
