import { getSessionID } from '../utils'

const url = import.meta.env.VITE_RTMP_URL

export default {
  liveURI: (streamName) => `rtmp://${url}/live/${streamName}`,
  liveKey: (roomName) =>
    encodeURIComponent(
      JSON.stringify({
        session_id: getSessionID(),
        room_settings: {
          name: roomName,
        },
      })
    ),
}
