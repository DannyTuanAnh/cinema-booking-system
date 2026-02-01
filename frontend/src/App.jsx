import { BrowserRouter } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { SessionProvider } from "./context/SessionContext";
import AppRoutes from "./routes/AppRoutes";
import "./styles/App.css";

function App() {
  return (
    <BrowserRouter>
      <SessionProvider>
        <AuthProvider>
          <AppRoutes />
        </AuthProvider>
      </SessionProvider>
    </BrowserRouter>
  );
}

export default App;
