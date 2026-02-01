import axiosClient from "./axiosConfig";

const movieApi = {
  // GET /api/movies - Lấy tất cả phim
  getMovies: () => {
    return axiosClient.get("/movies");
  },

  // GET /api/shows?movie_id=X - Lấy lịch chiếu theo phim
  getShows: (movieId) => {
    if (movieId) {
      return axiosClient.get("/shows", {
        params: {
          movie_id: movieId,
        },
      });
    }
    return axiosClient.get("/shows");
  },
};

export default movieApi;
