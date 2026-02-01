import { useState, useEffect } from "react";
import bookingApi from "../api/bookingApi";
import { unwrapResponse } from "../utils/errorHandler";
import Loading from "../components/Loading";
import ErrorMessage from "../components/ErrorMessage";
import "../styles/myTicketsPage.css";

const MyTicketsPage = () => {
  const [tickets, setTickets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadTickets = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await bookingApi.getMyBookings();
      const ticketsData = unwrapResponse(data);
      setTickets(Array.isArray(ticketsData) ? ticketsData : []);
    } catch (err) {
      console.error("Load tickets error:", err);
      setError("Không thể tải danh sách vé. Vui lòng thử lại.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTickets();
  }, []);

  const formatDateTime = (dateString) => {
    return new Date(dateString).toLocaleString("vi-VN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const isUpcoming = (showTime) => {
    return new Date(showTime) > new Date();
  };

  const upcomingTickets = tickets.filter((t) => isUpcoming(t.show_time));
  const pastTickets = tickets.filter((t) => !isUpcoming(t.show_time));

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={loadTickets} />;
  }

  return (
    <div className="my-tickets-page">
      <div className="container">
        <h1 className="page-title">Vé của tôi</h1>

        {tickets.length === 0 ? (
          <div className="empty-state">
            <p>Bạn chưa có vé nào</p>
          </div>
        ) : (
          <>
            {upcomingTickets.length > 0 && (
              <section className="tickets-section">
                <h2 className="section-title">
                  Sắp chiếu ({upcomingTickets.length})
                </h2>
                <div className="tickets-grid">
                  {upcomingTickets.map((ticket) => (
                    <div
                      key={ticket.ticket_id}
                      className="ticket-card upcoming"
                    >
                      <div className="ticket-header">
                        <span className="ticket-status">Sắp chiếu</span>
                        <span className="ticket-id">#{ticket.ticket_id}</span>
                      </div>
                      <div className="ticket-body">
                        <h3 className="ticket-movie">
                          🎬 Suất chiếu: {ticket.show_id}
                        </h3>
                        <p className="ticket-info">
                          <span>📅 {formatDateTime(ticket.show_time)}</span>
                        </p>
                        <p className="ticket-info">
                          <span>💺 Ghế: {ticket.seat_name}</span>
                        </p>
                        <p className="ticket-info">
                          <span>📍 Phòng chiếu</span>
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </section>
            )}

            {pastTickets.length > 0 && (
              <section className="tickets-section">
                <h2 className="section-title">Đã xem ({pastTickets.length})</h2>
                <div className="tickets-grid">
                  {pastTickets.map((ticket) => (
                    <div key={ticket.ticket_id} className="ticket-card past">
                      <div className="ticket-header">
                        <span className="ticket-status">Đã xem</span>
                        <span className="ticket-id">#{ticket.ticket_id}</span>
                      </div>
                      <div className="ticket-body">
                        <h3 className="ticket-movie">
                          🎬 Suất chiếu: {ticket.show_id}
                        </h3>
                        <p className="ticket-info">
                          <span>📅 {formatDateTime(ticket.show_time)}</span>
                        </p>
                        <p className="ticket-info">
                          <span>💺 Ghế: {ticket.seat_name}</span>
                        </p>
                        <p className="ticket-info">
                          <span>📍 Phòng chiếu</span>
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </section>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default MyTicketsPage;
