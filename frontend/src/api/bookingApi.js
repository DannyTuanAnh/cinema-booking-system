import axiosClient from "./axiosConfig";

const bookingApi = {
  // POST /api/book - Đặt ghế (1 hoặc nhiều ghế)
  // Backend expects: { "seats": [6, 9] } (array of seat IDs)
  // JWT token tự động gửi qua header Authorization
  bookSeats: (seatIds) => {
    // Accept cả số đơn hoặc array
    const seats = Array.isArray(seatIds) ? seatIds : [seatIds];

    if (seats.length === 0) {
      return Promise.reject(new Error("At least 1 seat must be selected"));
    }

    return axiosClient.post("/book", {
      seats: seats, // Array of seat IDs: [6, 9]
    });
  },

  // GET /api/tickets - Lấy danh sách vé đã đặt (protected route)
  getMyBookings: () => {
    return axiosClient.get("/tickets");
  },

  // DELETE /api/cancel/:bookingId - Hủy booking
  cancelBooking: (bookingId) => {
    if (!bookingId) {
      return Promise.reject(new Error("bookingId is required"));
    }
    return axiosClient.delete(`/cancel/${bookingId}`);
  },
};

export default bookingApi;
