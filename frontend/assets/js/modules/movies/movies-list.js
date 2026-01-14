/**
 * ============================================
 * MOVIES LIST MODULE
 * Hiển thị danh sách phim
 * ============================================
 */

// DOM elements
const loadingContainer = document.getElementById('loadingContainer');
const errorContainer = document.getElementById('errorContainer');
const moviesContainer = document.getElementById('moviesContainer');
const emptyContainer = document.getElementById('emptyContainer');

/**
 * Format duration (minutes to hours:minutes)
 */
function formatDuration(minutes) {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}m`;
}

/**
 * Format date
 */
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('vi-VN', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
    });
}

/**
 * Create movie card HTML
 */
function createMovieCard(movie) {
    return `
        <div class="movie-card" onclick="goToShowtimes(${movie.movie_id})">
            <img 
                src="${movie.url_image || 'https://via.placeholder.com/280x400?text=No+Image'}" 
                alt="${movie.title}" 
                class="movie-poster"
                onerror="this.src='https://via.placeholder.com/280x400?text=No+Image'"
            >
            <div class="movie-info">
                <h3 class="movie-title">${movie.title}</h3>
                <div class="movie-meta">
                    <div class="movie-meta-item">
                        ⏱️ ${formatDuration(movie.duration)}
                    </div>
                    <div class="movie-meta-item">
                        ⭐ ${movie.rate.toFixed(1)}
                    </div>
                </div>
                <span class="movie-genre">${movie.genre}</span>
                <p class="movie-description">${movie.description}</p>
                <div class="movie-actions">
                    <button class="btn btn-primary btn-book" onclick="event.stopPropagation(); goToShowtimes(${movie.movie_id})">
                        Đặt vé ngay
                    </button>
                </div>
            </div>
        </div>
    `;
}

/**
 * Show error message
 */
function showError(message) {
    loadingContainer.classList.add('hidden');
    moviesContainer.classList.add('hidden');
    emptyContainer.classList.add('hidden');
    errorContainer.classList.remove('hidden');
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
    loadingContainer.classList.add('hidden');
    moviesContainer.classList.add('hidden');
    errorContainer.classList.add('hidden');
    emptyContainer.classList.remove('hidden');
}

/**
 * Show movies
 */
function showMovies(movies) {
    loadingContainer.classList.add('hidden');
    errorContainer.classList.add('hidden');
    emptyContainer.classList.add('hidden');
    moviesContainer.classList.remove('hidden');

    moviesContainer.innerHTML = movies.map(movie => createMovieCard(movie)).join('');
}

/**
 * Navigate to showtimes page
 */
function goToShowtimes(movieId) {
    window.location.href = `showtimes.html?movie_id=${movieId}`;
}

/**
 * Load movies from API
 */
async function loadMovies() {
    try {
        // Show loading
        loadingContainer.classList.remove('hidden');
        errorContainer.classList.add('hidden');
        moviesContainer.classList.add('hidden');
        emptyContainer.classList.add('hidden');

        // Fetch movies
        const movies = await API.get(API_CONFIG.ENDPOINTS.MOVIES);

        if (!movies || movies.length === 0) {
            showEmpty();
            return;
        }

        showMovies(movies);

    } catch (error) {
        console.error('Load movies error:', error);
        showError(error.message);
    }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    loadMovies();
});

// Make goToShowtimes globally available for onclick handlers
window.goToShowtimes = goToShowtimes;
