import axiosClient from "./axiosConfig";

const seatApi = {
  getSeats: (showId) => {
    if (!showId) {
      return Promise.reject(new Error("showId is required"));
    }

    return axiosClient.get("/seats", {
      params: {
        show_id: showId,
      },
    });
  },
};

export default seatApi;
