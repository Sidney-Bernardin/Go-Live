import axios from 'axios'
import { getSessionID } from '../utils'

export default {
  getSelf: (fields) =>
    axios.get(`users/self?fields=${fields?.join()}`, {
      headers: { Authorization: `Bearer ${getSessionID()}` },
    }),
  signup: (signupInfo) => axios.post('users/self/signup', signupInfo),
  signin: (signinInfo) => axios.post('users/self/signin', signinInfo),
  logout: () =>
    axios.get('users/self/logout', {
      headers: { Authorization: `Bearer ${getSessionID()}` },
    }),

  getUser: (search, searchBy, fields) =>
    axios.get(
      `users/all/${search}?search_by=${searchBy}&fields=${fields?.join()}`
    ),
  searchUsers: (username, offset, limit, fields) =>
    axios.get(
      `users/all?username=${username}&offset=${offset}&limit=${limit}&fields=${fields?.join()}`
    ),
  profilePictureSrc: (userID) =>
    `http://${import.meta.env.VITE_MICROSERVICES_URL}/users/all/${userID}/picture`,
}
