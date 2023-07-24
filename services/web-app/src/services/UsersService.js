import axios from 'axios'
import { getSessionID } from '../utils'

const url = 'http://' + import.meta.env.VITE_USERS_URL

export default {
  signup: (signupInfo) => axios.post(url + '/signup', signupInfo),
  signin: (signinInfo) => axios.post(url + '/signin', signinInfo),
  logout: () =>
    axios.get(url + '/self/logout', {
      headers: { Authorization: `Bearer ${getSessionID()}` },
    }),

  getSelf: (fields) =>
    axios.get(url + `/self?fields=${fields?.join()}`, {
      headers: { Authorization: `Bearer ${getSessionID()}` },
    }),

  getUser: (search, searchBy, fields) =>
    axios.get(
      url + `/all/${search}?by=${searchBy}&fields=${fields?.join()}`
    ),
  searchUsers: (username, offset, limit, fields) =>
    axios.get(
      url +
        `/all?username=${username}&offset=${offset}&limit=${limit}&fields=${fields?.join()}`
    ),
  setProfilePicture: (profilePicture) => {
    const formData = new FormData()
    formData.append('profile_picture', profilePicture)

    return axios.post(url + '/self/picture', formData, {
      headers: { Authorization: `Bearer ${getSessionID()}` },
    })
  },
  profilePictureSrc: (userID) => url + `/all/${userID}/picture`,
}
