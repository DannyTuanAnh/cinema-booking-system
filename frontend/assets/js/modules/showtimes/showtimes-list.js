/**
 * ============================================
 * SHOWTIMES LIST MODULE
 * Hiển thị danh sách suất chiếu theo phim
 * ============================================
 */

// DOM elements
const loadingContainer = document.getElementById("loadingContainer");
const errorContainer = document.getElementById("errorContainer");
const showtimesContainer = document.getElementById("showtimesContainer");
const emptyContainer = document.getElementById("emptyContainer");
const movieInfoElement = document.getElementById("movieInfo");
const showtimesListElement = document.getElementById("showtimesList");

// Get movie_id from URL
const urlParams = new URLSearchParams(window.location.search);
const movieId = urlParams.get("movie_id");

// Store movie data
let currentMovie = null;

/**
 * Format time
 */
function formatTime(dateString) {
  const date = new Date(dateString);
  return date.toLocaleTimeString("vi-VN", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
}

/**
 * Format date
 */
function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString("vi-VN", {
    weekday: "long",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

/**
 * Format duration
 */
function formatDuration(minutes) {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return `${hours}h ${mins}m`;
}

/**
 * Create movie info HTML
 */
function createMovieInfo(show) {
  return `
        <div class="movie-info-content">
            <img 
                src="https://via.placeholder.com/200x280?text=${encodeURIComponent(
                  show.movie_title
                )}" 
                alt="${show.movie_title}" 
                class="movie-info-poster"
            >
            <div class="movie-info-details">
                <h2>${show.movie_title}</h2>
                <div class="movie-info-meta">
                    <span>🎬 Phim</span>
                    <span>ID: ${show.movie_id}</span>
                </div>
                <p class="movie-info-description">
                    Chọn suất chiếu phù hợp với bạn để tiếp tục đặt vé.
                </p>
            </div>
        </div>
    `;
}

/**
 * Group showtimes by date
 */
function groupShowtimesByDate(showtimes) {
  const groups = {};

  showtimes.forEach((show) => {
    const date = new Date(show.show_time).toDateString();
    if (!groups[date]) {
      groups[date] = [];
    }
    groups[date].push(show);
  });

  return groups;
}

/**
 * Create showtime card HTML
 */
function createShowtimeCard(show) {
  return `
        <div class="showtime-card" onclick="selectShowtime(${show.show_id})">
            <div class="showtime-time">${formatTime(show.show_time)}</div>
            <div class="showtime-info">Suất ${show.show_id}</div>
        </div>
    `;
}

/**
 * Create showtimes list HTML
 */
function createShowtimesList(showtimes) {
  const groups = groupShowtimesByDate(showtimes);

  return Object.entries(groups)
    .map(
      ([date, shows]) => `
        <div class="showtime-group">
            <h3 class="showtime-date">${formatDate(shows[0].show_time)}</h3>
            <div class="showtime-slots">
                ${shows.map((show) => createShowtimeCard(show)).join("")}
            </div>
        </div>
    `
    )
    .join("");
}

/**
 * Show error message
 */
function showError(message) {
  loadingContainer.classList.add("hidden");
  showtimesContainer.classList.add("hidden");
  emptyContainer.classList.add("hidden");
  errorContainer.classList.remove("hidden");
  errorContainer.innerHTML = `
        <div class="alert alert-error">
            ${message}
        </div>
    `;
}

/**
 * Show empty state
 */
function showEmpty() {
  loadingContainer.classList.add("hidden");
  showtimesContainer.classList.add("hidden");
  errorContainer.classList.add("hidden");
  emptyContainer.classList.remove("hidden");
}

/**
 * Show showtimes
 */
function showShowtimes(showtimes) {
  loadingContainer.classList.add("hidden");
  errorContainer.classList.add("hidden");
  emptyContainer.classList.add("hidden");
  showtimesContainer.classList.remove("hidden");

  // Display movie info
  movieInfoElement.innerHTML = createMovieInfo(showtimes[0]);

  // Display showtimes list
  showtimesListElement.innerHTML = createShowtimesList(showtimes);
}

/**
 * Select showtime and navigate to seat selection
 */
function selectShowtime(showId) {
  window.location.href = `seat-selection.html?show_id=${showId}`;
}

/**
 * Load showtimes from API
 */
async function loadShowtimes() {
  if (!movieId) {
    showError("Không tìm thấy thông tin phim");
    return;
  }

  try {
    // Show loading
    loadingContainer.classList.remove("hidden");
    errorContainer.classList.add("hidden");
    showtimesContainer.classList.add("hidden");
    emptyContainer.classList.add("hidden");

    // Fetch showtimes
    const showtimes = await API.get(
      `${API_CONFIG.ENDPOINTS.SHOWS}?movie_id=${movieId}`
    );

    if (!showtimes || showtimes.length === 0) {
      showEmpty();
      return;
    }

    showShowtimes(showtimes);
  } catch (error) {
    console.error("Load showtimes error:", error);
    showError(error.message);
  }
}

// Initialize
document.addEventListener("DOMContentLoaded", () => {
  loadShowtimes();
});

// Make selectShowtime globally available for onclick handlers
window.selectShowtime = selectShowtime;
