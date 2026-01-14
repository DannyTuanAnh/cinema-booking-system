/**
 * ============================================
 * SEAT SELECTION MODULE
 * Chọn ghế và đặt vé
 * ============================================
 */

// Check authentication
if (!Auth.requireAuth()) {
  // Will redirect to login
}

// DOM elements
const loadingContainer = document.getElementById("loadingContainer");
const errorContainer = document.getElementById("errorContainer");
const seatContainer = document.getElementById("seatContainer");
const movieInfoElement = document.getElementById("movieInfo");
const seatsGridElement = document.getElementById("seatsGrid");
const selectedSeatsDisplay = document.getElementById("selectedSeatsDisplay");
const selectedCountElement = document.getElementById("selectedCount");
const bookBtn = document.getElementById("bookBtn");

// Get show_id from URL
const urlParams = new URLSearchParams(window.location.search);
const showId = urlParams.get("show_id");

// State
let seats = [];
let selectedSeats = [];
let isBooking = false;

/**
 * Format time
 */
function formatTime(dateString) {
  const date = new Date(dateString);
  return date.toLocaleString("vi-VN", {
    hour: "2-digit",
    minute: "2-digit",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

/**
 * Create movie info HTML
 */
function createMovieInfo(seats) {
  const firstSeat = seats[0];
  return `
        <div class="info-content">
            <h3>Thông tin suất chiếu</h3>
            <div class="info-meta">
                <span>📽️ Suất chiếu: ${showId}</span>
                <span>📍 Phòng chiếu</span>
            </div>
        </div>
    `;
}

/**
 * Group seats by row
 */
function groupSeatsByRow(seats) {
  const rows = {};

  seats.forEach((seat) => {
    const row = seat.seat_name.charAt(0); // First character (A, B, C, etc.)
    if (!rows[row]) {
      rows[row] = [];
    }
    rows[row].push(seat);
  });

  return rows;
}

/**
 * Toggle seat selection
 */
function toggleSeat(seatId, seatName, status) {
  if (status === SEAT_STATUS.BOOKED || isBooking) {
    return;
  }

  const index = selectedSeats.findIndex((s) => s.id === seatId);

  if (index > -1) {
    // Deselect
    selectedSeats.splice(index, 1);
  } else {
    // Select
    selectedSeats.push({ id: seatId, name: seatName });
  }

  updateUI();
}

/**
 * Update UI after selection changes
 */
function updateUI() {
  // Update seat visuals
  seats.forEach((seat) => {
    const seatElement = document.querySelector(
      `[data-seat-id="${seat.seat_id}"]`
    );
    if (!seatElement) return;

    const isSelected = selectedSeats.some((s) => s.id === seat.seat_id);

    seatElement.classList.remove(
      "seat-available",
      "seat-selected",
      "seat-booked"
    );

    if (seat.status === SEAT_STATUS.BOOKED) {
      seatElement.classList.add("seat-booked");
    } else if (isSelected) {
      seatElement.classList.add("seat-selected");
    } else {
      seatElement.classList.add("seat-available");
    }
  });

  // Update summary
  if (selectedSeats.length > 0) {
    const seatNames = selectedSeats.map((s) => s.name).join(", ");
    selectedSeatsDisplay.textContent = seatNames;
    selectedCountElement.textContent = selectedSeats.length;
    bookBtn.disabled = false;
  } else {
    selectedSeatsDisplay.textContent = "Chưa chọn";
    selectedCountElement.textContent = "0";
    bookBtn.disabled = true;
  }
}

/**
 * Create seats grid HTML
 */
function createSeatsGrid(seats) {
  const rows = groupSeatsByRow(seats);
  const sortedRows = Object.keys(rows).sort();

  return sortedRows
    .map((row) => {
      const rowSeats = rows[row].sort((a, b) =>
        a.seat_name.localeCompare(b.seat_name)
      );

      return `
            <div class="seat-row">
                <div class="seat-label">${row}</div>
                ${rowSeats
                  .map((seat) => {
                    const statusClass =
                      seat.status === SEAT_STATUS.BOOKED
                        ? "seat-booked"
                        : "seat-available";

                    return `
                        <div 
                            class="seat ${statusClass}" 
                            data-seat-id="${seat.seat_id}"
                            data-seat-name="${seat.seat_name}"
                            onclick="toggleSeat(${seat.seat_id}, '${seat.seat_name}', '${seat.status}')"
                        >
                            ${seat.seat_name}
                        </div>
                    `;
                  })
                  .join("")}
            </div>
        `;
    })
    .join("");
}

/**
 * Show error message
 */
function showError(message) {
  loadingContainer.classList.add("hidden");
  seatContainer.classList.add("hidden");
  errorContainer.classList.remove("hidden");
  errorContainer.innerHTML = `
        <div class="alert alert-error">
            ${message}
        </div>
    `;
}

/**
 * Show seats
 */
function showSeats(seatsData) {
  seats = seatsData;

  loadingContainer.classList.add("hidden");
  errorContainer.classList.add("hidden");
  seatContainer.classList.remove("hidden");

  // Display movie info
  movieInfoElement.innerHTML = createMovieInfo(seatsData);

  // Display seats grid
  seatsGridElement.innerHTML = createSeatsGrid(seatsData);
}

/**
 * Load seats from API
 */
async function loadSeats() {
  if (!showId) {
    showError("Không tìm thấy thông tin suất chiếu");
    return;
  }

  try {
    // Show loading
    loadingContainer.classList.remove("hidden");
    errorContainer.classList.add("hidden");
    seatContainer.classList.add("hidden");

    // Fetch seats
    const seatsData = await API.get(
      `${API_CONFIG.ENDPOINTS.SEATS}?show_id=${showId}`
    );

    if (!seatsData || seatsData.length === 0) {
      showError("Không có ghế nào cho suất chiếu này");
      return;
    }

    showSeats(seatsData);
  } catch (error) {
    console.error("Load seats error:", error);
    showError(error.message);
  }
}

/**
 * Handle booking
 */
async function handleBooking() {
  if (selectedSeats.length === 0 || isBooking) {
    return;
  }

  if (
    !confirm(
      `Bạn có chắc muốn đặt ${selectedSeats.length} ghế: ${selectedSeats
        .map((s) => s.name)
        .join(", ")}?`
    )
  ) {
    return;
  }

  isBooking = true;
  bookBtn.disabled = true;
  bookBtn.textContent = "Đang đặt vé...";

  try {
    // Get seat IDs
    const seatIds = selectedSeats.map((s) => s.id);

    // Call booking API
    await API.post(
      API_CONFIG.ENDPOINTS.BOOK,
      {
        seats: seatIds,
      },
      true
    ); // Include auth token

    // Show success message
    alert(SUCCESS_MESSAGES.BOOKING);

    // Redirect to my tickets
    window.location.href = "my-tickets.html";
  } catch (error) {
    console.error("Booking error:", error);
    alert(error.message);

    // Reset state
    isBooking = false;
    bookBtn.disabled = false;
    bookBtn.textContent = "Đặt vé ngay";

    // Reload seats to get updated status
    selectedSeats = [];
    loadSeats();
  }
}

// Initialize
document.addEventListener("DOMContentLoaded", () => {
  loadSeats();
});

// Make functions globally available
window.toggleSeat = toggleSeat;
window.handleBooking = handleBooking;
