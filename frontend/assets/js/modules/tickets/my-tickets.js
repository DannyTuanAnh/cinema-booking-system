/**
 * ============================================
 * MY TICKETS MODULE
 * Hiển thị danh sách vé đã đặt
 * ============================================
 */

// Check authentication
if (!Auth.requireAuth()) {
  // Will redirect to login
}

// DOM elements
const loadingContainer = document.getElementById("loadingContainer");
const errorContainer = document.getElementById("errorContainer");
const ticketsContainer = document.getElementById("ticketsContainer");
const emptyContainer = document.getElementById("emptyContainer");

/**
 * Format date time
 */
function formatDateTime(dateString) {
  const date = new Date(dateString);
  return date.toLocaleString("vi-VN", {
    weekday: "long",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

/**
 * Format date only
 */
function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString("vi-VN", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

/**
 * Check if show is upcoming
 */
function isUpcoming(dateString) {
  const showDate = new Date(dateString);
  const now = new Date();
  return showDate > now;
}

/**
 * Create ticket card HTML
 */
function createTicketCard(ticket) {
  const upcoming = isUpcoming(ticket.show_time);
  const statusClass = upcoming ? "status-upcoming" : "status-past";
  const statusText = upcoming ? "Sắp chiếu" : "Đã chiếu";

  return `
        <div class="ticket-card">
            <div class="ticket-header">
                <div>
                    <h3 class="ticket-title">${ticket.title}</h3>
                    <p class="ticket-id">Booking ID: #${ticket.booking_id}</p>
                </div>
                <div class="ticket-status ${statusClass}">
                    ${statusText}
                </div>
            </div>
            <div class="ticket-body">
                <div class="ticket-info-group">
                    <span class="ticket-info-label">Suất chiếu</span>
                    <span class="ticket-info-value">
                        <span class="icon">📽️</span>
                        Suất ${ticket.show_id}
                    </span>
                </div>
                <div class="ticket-info-group">
                    <span class="ticket-info-label">Thời gian chiếu</span>
                    <span class="ticket-info-value">
                        <span class="icon">🕐</span>
                        ${formatDateTime(ticket.show_time)}
                    </span>
                </div>
                <div class="ticket-info-group">
                    <span class="ticket-info-label">Ghế ngồi</span>
                    <span class="ticket-info-value">
                        <span class="icon">💺</span>
                        ${ticket.seat_name}
                    </span>
                </div>
                <div class="ticket-info-group">
                    <span class="ticket-info-label">Ngày đặt</span>
                    <span class="ticket-info-value">
                        <span class="icon">📅</span>
                        ${formatDate(ticket.book_at)}
                    </span>
                </div>
            </div>
        </div>
    `;
}

/**
 * Show error message
 */
function showError(message) {
  loadingContainer.classList.add("hidden");
  ticketsContainer.classList.add("hidden");
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
  ticketsContainer.classList.add("hidden");
  errorContainer.classList.add("hidden");
  emptyContainer.classList.remove("hidden");
}

/**
 * Show tickets
 */
function showTickets(tickets) {
  loadingContainer.classList.add("hidden");
  errorContainer.classList.add("hidden");
  emptyContainer.classList.add("hidden");
  ticketsContainer.classList.remove("hidden");

  // Sort tickets by show time (upcoming first, then by date desc)
  const sortedTickets = tickets.sort((a, b) => {
    const aUpcoming = isUpcoming(a.show_time);
    const bUpcoming = isUpcoming(b.show_time);

    if (aUpcoming && !bUpcoming) return -1;
    if (!aUpcoming && bUpcoming) return 1;

    return new Date(b.show_time) - new Date(a.show_time);
  });

  ticketsContainer.innerHTML = sortedTickets
    .map((ticket) => createTicketCard(ticket))
    .join("");
}

/**
 * Load tickets from API
 */
async function loadTickets() {
  try {
    // Show loading
    loadingContainer.classList.remove("hidden");
    errorContainer.classList.add("hidden");
    ticketsContainer.classList.add("hidden");
    emptyContainer.classList.add("hidden");

    // Fetch tickets (requires authentication)
    const tickets = await API.get(API_CONFIG.ENDPOINTS.MY_TICKETS, true);

    if (!tickets || tickets.length === 0) {
      showEmpty();
      return;
    }

    showTickets(tickets);
  } catch (error) {
    console.error("Load tickets error:", error);

    // If 401, Auth.requireAuth will redirect
    if (error.message === ERROR_MESSAGES.UNAUTHORIZED) {
      return;
    }

    showError(error.message);
  }
}

// Initialize
document.addEventListener("DOMContentLoaded", () => {
  loadTickets();
});
