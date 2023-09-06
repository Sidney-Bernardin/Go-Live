import axios from "axios";

const url = "http://" + import.meta.env.VITE_ROOMS_URL;

export const getRoom = (roomID: string) => axios.get(`${url}/${roomID}`);
