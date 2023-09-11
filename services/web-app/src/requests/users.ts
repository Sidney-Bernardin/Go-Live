import axios, { AxiosResponse } from "axios";
import { getSessionID } from "../utils";

const url = "http://" + import.meta.env.VITE_USERS_URL;
const authHeaders = () => ({ Authorization: `Bearer ${getSessionID()}` });

export const signup = (fd: FormData): Promise<AxiosResponse> =>
  axios.post(`${url}/signup`, fd);

export const signin = (fd: FormData): Promise<AxiosResponse> =>
  axios.post(`${url}/signin`, fd);

export const logout = (): Promise<AxiosResponse> =>
  axios.delete(`${url}/auth/logout`, {
    headers: authHeaders(),
  });

export const authenticateUser = (fiedls: string[]): Promise<AxiosResponse> =>
  axios.get(`${url}/auth/?fields=${fiedls.join(",")}`, {
    headers: authHeaders(),
  });

export const searchUsers = (
  username: string,
  fields: string[],
): Promise<AxiosResponse> =>
  axios.get(`${url}/${username}/0/100?fields=${fields.join(",")}`);

export const getUser = (userID: string): Promise<AxiosResponse> =>
  axios.get(`${url}/${userID}`);

export const deleteAccount = (): Promise<AxiosResponse> =>
  axios.delete(`${url}/auth`, { headers: authHeaders() });

export const getProfilePictureURL = (userID: string): string =>
  `${url}/${userID}/picture`;
