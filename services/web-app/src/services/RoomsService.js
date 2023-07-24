import axios from 'axios'
import { getSessionID } from '../utils'

const url = import.meta.env.VITE_ROOMS_URL

export default {
  getRoom: (roomID) => axios.get('http://' + url + `/all/${roomID}`),
  joinRoom: (roomID) =>
    new WebSocket(
      `ws://${url}/join?session_id=${getSessionID()}&room_id=${roomID}`
    ),
}
