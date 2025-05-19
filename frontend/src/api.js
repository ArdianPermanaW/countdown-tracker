import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080/api/events",
});

export default api;