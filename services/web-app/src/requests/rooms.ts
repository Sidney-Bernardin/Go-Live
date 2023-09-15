import axios from "axios";
import { getSessionID } from "../utils";
import { Room } from "./models";

const httpURL = "http://" + import.meta.env.VITE_ROOMS_URL;
const wsURL = "ws://" + import.meta.env.VITE_ROOMS_URL;

export const getRoom = (roomID: string): Promise<Room> => axios.get(`${httpURL}/${roomID}`).then((res) => res.data);

export const joinRoom = (roomID: string) => new WebSocket(`${wsURL}/join/${roomID}?sid=${getSessionID()}`);
