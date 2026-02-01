import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import "../styles/movieCard.css";

const MovieCard = ({ movie }) => {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const handleBooking = () => {
    if (!isAuthenticated()) {
      navigate("/login");
      return;
    }
    navigate(`/showtimes/${movie.movie_id}`);
  };

  return (
    <div className="movie-card">
      <div className="movie-poster">
        <img src={movie.url_image} alt={movie.title} />
        <div className="movie-overlay">
          <button onClick={handleBooking} className="btn btn-primary">
            Đặt vé
          </button>
        </div>
      </div>
      <div className="movie-info">
        <h3 className="movie-title">{movie.title}</h3>
        <p className="movie-genre">{movie.genre}</p>
        <p className="movie-duration">⏱️ {movie.duration} phút</p>
      </div>
    </div>
  );
};

export default MovieCard;
