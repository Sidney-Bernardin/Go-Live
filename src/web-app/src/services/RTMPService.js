import { getSessionID } from '../utils'

export default {
  liveURI: (streamName) =>
    `rtmp://${import.meta.env.VITE_RTMP_URL}/live/${streamName}`,
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
