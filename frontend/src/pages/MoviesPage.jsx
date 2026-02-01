import { useState, useEffect } from "react";
import movieApi from "../api/movieApi";
import { unwrapResponse } from "../utils/errorHandler";
import MovieCard from "../components/MovieCard";
import Loading from "../components/Loading";
import ErrorMessage from "../components/ErrorMessage";
import "../styles/moviesPage.css";

const MoviesPage = () => {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadMovies = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await movieApi.getMovies();
      const moviesData = unwrapResponse(data);
      setMovies(Array.isArray(moviesData) ? moviesData : []);
    } catch (err) {
      console.error("Load movies error:", err);
      setError("Không thể tải danh sách phim. Vui lòng thử lại.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadMovies();
  }, []);

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={loadMovies} />;
  }

  return (
    <div className="movies-page">
      <div className="container">
        <h1 className="page-title">Phim đang chiếu</h1>

        {movies.length === 0 ? (
          <div className="empty-state">
            <p>Hiện chưa có phim nào</p>
          </div>
        ) : (
          <div className="movies-grid">
            {movies.map((movie) => (
              <MovieCard key={movie.movie_id} movie={movie} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default MoviesPage;
