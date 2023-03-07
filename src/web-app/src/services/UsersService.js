import axios from 'axios'
import { getSessionID } from '../utils'

export default {
  signup: (signupInfo) => axios.post('users/signup', signupInfo),
  signin: (signinInfo) => axios.post('users/signin', signinInfo),
  logout: () =>
    axios.get('users/self/logout', {
      headers: { Authorization: `Bearer ${getSessionID()}` },
    }),

  getSelf: (fields) =>
    axios.get(`users/self?fields=${fields?.join()}`, {
      headers: { Authorization: `Bearer ${getSessionID()}` },
    }),

  getUser: (search, searchBy, fields) =>
    axios.get(`users/all/${search}?by=${searchBy}&fields=${fields?.join()}`),
  searchUsers: (username, offset, limit, fields) =>
    axios.get(
      `users/all?username=${username}&offset=${offset}&limit=${limit}&fields=${fields?.join()}`
    ),
  profilePictureSrc: (userID) =>
    `http://${
      import.meta.env.VITE_MICROSERVICES_URL
    }/users/all/${userID}/picture`,
}
