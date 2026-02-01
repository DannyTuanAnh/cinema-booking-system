import { createContext, useContext, useState, useEffect } from "react";
import SessionExpiredModal from "../components/SessionExpiredModal";
import { sessionExpiredEmitter } from "../utils/sessionEmitter";

const SessionContext = createContext(null);

export const useSession = () => {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error("useSession must be used within SessionProvider");
  }
  return context;
};

export const SessionProvider = ({ children }) => {
  const [showExpiredModal, setShowExpiredModal] = useState(false);

  useEffect(() => {
    // Listen for session expired events
    const unsubscribe = sessionExpiredEmitter.subscribe(() => {
      setShowExpiredModal(true);
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const handleSessionExpired = () => {
    setShowExpiredModal(true);
  };

  const handleModalClose = () => {
    setShowExpiredModal(false);
    // Clear localStorage
    localStorage.removeItem("userID");
    localStorage.removeItem("email");
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    // Redirect to login
    window.location.href = "/login";
  };

  const value = {
    handleSessionExpired,
  };

  return (
    <SessionContext.Provider value={value}>
      {children}
      {showExpiredModal && <SessionExpiredModal onClose={handleModalClose} />}
    </SessionContext.Provider>
  );
};
