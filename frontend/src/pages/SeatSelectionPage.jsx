import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import seatApi from "../api/seatApi";
import bookingApi from "../api/bookingApi";
import { unwrapResponse } from "../utils/errorHandler";
import Loading from "../components/Loading";
import ErrorMessage from "../components/ErrorMessage";
import "../styles/seatSelectionPage.css";

const SEAT_STATUS = {
  AVAILABLE: "available",
  BOOKED: "booked",
};

const SeatSelectionPage = () => {
  const { showId } = useParams();
  const navigate = useNavigate();
  const [seats, setSeats] = useState([]);
  const [selectedSeats, setSelectedSeats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [booking, setBooking] = useState(false);

  const loadSeats = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await seatApi.getSeats(showId);
      const seatsData = unwrapResponse(data);
      setSeats(Array.isArray(seatsData) ? seatsData : []);
    } catch (err) {
      console.error("Load seats error:", err);
      setError("Không thể tải danh sách ghế. Vui lòng thử lại.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSeats();
  }, [showId]);

  const groupSeatsByRow = (seats) => {
    const rows = {};
    seats.forEach((seat) => {
      const row = seat.seat_name.charAt(0);
      if (!rows[row]) {
        rows[row] = [];
      }
      rows[row].push(seat);
    });
    return rows;
  };

  const toggleSeat = (seat) => {
    if (seat.status === SEAT_STATUS.BOOKED || booking) {
      return;
    }

    const isSelected = selectedSeats.some((s) => s.seat_id === seat.seat_id);

    if (isSelected) {
      setSelectedSeats(selectedSeats.filter((s) => s.seat_id !== seat.seat_id));
    } else {
      setSelectedSeats([...selectedSeats, seat]);
    }
  };

  const handleBooking = async () => {
    if (selectedSeats.length === 0 || booking) {
      return;
    }

    const seatNames = selectedSeats.map((s) => s.seat_name).join(", ");
    if (
      !window.confirm(
        `Bạn có chắc muốn đặt ${selectedSeats.length} ghế: ${seatNames}?`,
      )
    ) {
      return;
    }

    setBooking(true);

    try {
      const seatIds = selectedSeats.map((s) => s.seat_id);
      await bookingApi.bookSeats(seatIds);

      alert("Đặt vé thành công!");
      navigate("/my-tickets");
    } catch (err) {
      console.error("Booking error:", err);

      if (err.response?.status === 409) {
        alert("Một hoặc nhiều ghế đã được đặt. Vui lòng chọn lại.");
        setSelectedSeats([]);
        loadSeats();
      } else {
        alert(
          err.response?.data?.error || "Đặt vé thất bại. Vui lòng thử lại.",
        );
      }
    } finally {
      setBooking(false);
    }
  };

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={loadSeats} />;
  }

  const groupedSeats = groupSeatsByRow(seats);
  const sortedRows = Object.keys(groupedSeats).sort();

  return (
    <div className="seat-selection-page">
      <div className="container">
        <h1 className="page-title">Chọn ghế</h1>

        <div className="cinema-info">
          <p>📽️ Suất chiếu: {showId}</p>
          <p>📍 Phòng chiếu</p>
        </div>

        <div className="screen">
          <div className="screen-label">Màn hình</div>
        </div>

        <div className="seats-container">
          {sortedRows.map((row) => (
            <div key={row} className="seat-row">
              <div className="seat-label">{row}</div>
              <div className="seats">
                {groupedSeats[row]
                  .sort((a, b) => a.seat_name.localeCompare(b.seat_name))
                  .map((seat) => {
                    const isSelected = selectedSeats.some(
                      (s) => s.seat_id === seat.seat_id,
                    );
                    const isBooked = seat.status === SEAT_STATUS.BOOKED;

                    return (
                      <button
                        key={seat.seat_id}
                        className={`seat ${isBooked ? "seat-booked" : isSelected ? "seat-selected" : "seat-available"}`}
                        onClick={() => toggleSeat(seat)}
                        disabled={isBooked || booking}
                      >
                        {seat.seat_name}
                      </button>
                    );
                  })}
              </div>
            </div>
          ))}
        </div>

        <div className="seat-legend">
          <div className="legend-item">
            <div className="seat seat-available"></div>
            <span>Còn trống</span>
          </div>
          <div className="legend-item">
            <div className="seat seat-selected"></div>
            <span>Đang chọn</span>
          </div>
          <div className="legend-item">
            <div className="seat seat-booked"></div>
            <span>Đã đặt</span>
          </div>
        </div>

        <div className="booking-summary">
          <div className="summary-info">
            <p>
              Ghế đã chọn:{" "}
              <strong>
                {selectedSeats.length > 0
                  ? selectedSeats.map((s) => s.seat_name).join(", ")
                  : "Chưa chọn"}
              </strong>
            </p>
            <p>
              Số lượng: <strong>{selectedSeats.length}</strong>
            </p>
          </div>
          <button
            className="btn btn-primary btn-lg"
            onClick={handleBooking}
            disabled={selectedSeats.length === 0 || booking}
          >
            {booking ? "Đang đặt vé..." : "Đặt vé ngay"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default SeatSelectionPage;
