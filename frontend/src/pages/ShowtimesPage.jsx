import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import movieApi from "../api/movieApi";
import { unwrapResponse } from "../utils/errorHandler";
import Loading from "../components/Loading";
import ErrorMessage from "../components/ErrorMessage";
import "../styles/showtimesPage.css";

const ShowtimesPage = () => {
  const { movieId } = useParams();
  const navigate = useNavigate();
  const [movie, setMovie] = useState(null);
  const [shows, setShows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Load shows by movieId - API returns shows array
      const data = await movieApi.getShows(movieId);
      const showsData = unwrapResponse(data);

      setShows(Array.isArray(showsData) ? showsData : []);
    } catch (err) {
      console.error("Load showtimes error:", err);
      setError("Không thể tải danh sách suất chiếu. Vui lòng thử lại.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [movieId]);

  const groupShowsByDate = (shows) => {
    const grouped = {};
    shows.forEach((show) => {
      const date = new Date(show.show_time).toLocaleDateString("vi-VN");
      if (!grouped[date]) {
        grouped[date] = [];
      }
      grouped[date].push(show);
    });
    return grouped;
  };

  const formatTime = (dateString) => {
    return new Date(dateString).toLocaleTimeString("vi-VN", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const handleSelectShow = (showId) => {
    navigate(`/seats/${showId}`);
  };

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={loadData} />;
  }

  const groupedShows = groupShowsByDate(shows);

  return (
    <div className="showtimes-page">
      <div className="container">
        {movie && (
          <div className="movie-header">
            <img
              src={movie.poster_url || "/placeholder-movie.jpg"}
              alt={movie.title}
              className="movie-poster-small"
            />
            <div className="movie-details">
              <h1>{movie.title}</h1>
              <p className="movie-meta">
                <span>{movie.genre}</span>
                <span>⏱️ {movie.duration} phút</span>
              </p>
            </div>
          </div>
        )}

        <h2 className="section-title">Chọn suất chiếu</h2>

        {shows.length === 0 ? (
          <div className="empty-state">
            <p>Hiện chưa có suất chiếu nào</p>
          </div>
        ) : (
          <div className="showtimes-list">
            {Object.entries(groupedShows).map(([date, dateShows]) => (
              <div key={date} className="showtime-group">
                <h3 className="showtime-date">📅 {date}</h3>
                <div className="showtime-buttons">
                  {dateShows.map((show) => (
                    <button
                      key={show.show_id}
                      className="showtime-btn"
                      onClick={() => handleSelectShow(show.show_id)}
                    >
                      {formatTime(show.show_time)}
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default ShowtimesPage;
