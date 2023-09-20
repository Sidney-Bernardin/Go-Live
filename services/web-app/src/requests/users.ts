import axios from "axios";
import { getSessionID } from "../utils";
import { User, LoginRes } from "./models";

const url = "http://" + import.meta.env.VITE_USERS_URL;
const authHeaders = () => ({ Authorization: `Bearer ${getSessionID()}` });

export const signup = (fd: FormData): Promise<LoginRes> => axios.post(`${url}/signup`, fd).then((res) => res.data);

export const signin = (fd: FormData): Promise<LoginRes> => axios.post(`${url}/signin`, fd).then((res) => res.data);

export const logout = (): Promise<void> => axios.delete(`${url}/auth/logout`, { headers: authHeaders() });

export const authenticateUser = (fiedls: string[]): Promise<User> =>
    axios.get(`${url}/auth/?fields=${fiedls.join(",")}`, { headers: authHeaders() }).then((res) => res.data);

export const searchUsers = (username: string, fields: string[]): Promise<User[]> =>
    axios.get(`${url}/${username ? username : "undefined"}/0/100?fields=${fields.join(",")}`).then((res) => res.data);

export const getUser = (userID: string, fields: string[]): Promise<User> =>
    axios.get(`${url}/${userID}?fields=${fields.join(",")}`).then((res) => res.data);

export const deleteAccount = (): Promise<void> => axios.delete(`${url}/auth`, { headers: authHeaders() });

export const getProfilePictureURL = (userID: string): string => `${url}/${userID}/picture`;
