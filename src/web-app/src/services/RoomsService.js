import axios from 'axios'
import { getSessionID } from '../utils'

export default {
  getRoom: (roomID) => axios.get(`rooms/all/${roomID}`),
  joinRoom: (roomID) =>
    new WebSocket(
      `ws://${import.meta.env.VITE_MICROSERVICES_URL
      }/rooms/join?session_id=${getSessionID()}&room_id=${roomID}`
    ),
}
