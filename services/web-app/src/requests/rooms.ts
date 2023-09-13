import axios from "axios";
import { Room } from "./models";

const url = "http://" + import.meta.env.VITE_ROOMS_URL;

export const getRoom = (roomID: string): Promise<Room> =>
  axios.get(`${url}/${roomID}`).then((res) => res.data);
