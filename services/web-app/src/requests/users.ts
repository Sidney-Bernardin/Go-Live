import axios from "axios";
import { getSessionID } from "../utils";

const url = "http://" + import.meta.env.VITE_USERS_URL;
const authHeaders = () => ({ Authorization: `Bearer ${getSessionID()}` });

export const signup = (info: object) => axios.post(`${url}/signup`, info);
export const signin = (info: object) => axios.post(`${url}/signin`, info);
export const logout = () =>
  axios.delete(`${url}/auth/logout`, {
    headers: authHeaders(),
  });

export const authenticateUser = (fiedls: string[]) =>
  axios.get(`${url}/auth/?fields=${fiedls.join(",")}`, {
    headers: authHeaders(),
  });

export const searchUsers = (username: string, fields: string[]) =>
  axios.get(`${url}/${username}/0/100?fields=${fields.join(",")}`);

export const getUser = (userID: string) => axios.get(`${url}/${userID}`);

export const updateProfilePicture = (profilePicture: string | Blob) => {
  const formData = new FormData();
  formData.append("profile_picture", profilePicture);

  return axios.put(url + "/auth/picture/", formData, {
    headers: { Authorization: `Bearer ${getSessionID()}` },
  });
};

export const getProfilePictureURL = (userID: string) =>
  `${url}/${userID}/picture`;
