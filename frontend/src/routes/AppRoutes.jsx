import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import MainLayout from "../layout/MainLayout";
import AuthLayout from "../layout/AuthLayout";

// Pages
import HomePage from "../pages/HomePage";
import LoginPage from "../pages/LoginPage";
import RegisterPage from "../pages/RegisterPage";
import MoviesPage from "../pages/MoviesPage";
import ShowtimesPage from "../pages/ShowtimesPage";
import SeatSelectionPage from "../pages/SeatSelectionPage";
import MyTicketsPage from "../pages/MyTicketsPage";

// Protected Route wrapper
const ProtectedRoute = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated() ? children : <Navigate to="/login" replace />;
};

// Public Only Route (redirect if authenticated)
const PublicOnlyRoute = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return !isAuthenticated() ? children : <Navigate to="/" replace />;
};

const AppRoutes = () => {
  return (
    <Routes>
      {/* Public routes with MainLayout */}
      <Route element={<MainLayout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/movies" element={<MoviesPage />} />
      </Route>

      {/* Auth routes with AuthLayout */}
      <Route element={<AuthLayout />}>
        <Route
          path="/login"
          element={
            <PublicOnlyRoute>
              <LoginPage />
            </PublicOnlyRoute>
          }
        />
        <Route
          path="/register"
          element={
            <PublicOnlyRoute>
              <RegisterPage />
            </PublicOnlyRoute>
          }
        />
      </Route>

      {/* Protected routes with MainLayout */}
      <Route element={<MainLayout />}>
        <Route
          path="/showtimes/:movieId"
          element={
            <ProtectedRoute>
              <ShowtimesPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/seats/:showId"
          element={
            <ProtectedRoute>
              <SeatSelectionPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/my-tickets"
          element={
            <ProtectedRoute>
              <MyTicketsPage />
            </ProtectedRoute>
          }
        />
      </Route>

      {/* 404 */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
};

export default AppRoutes;
