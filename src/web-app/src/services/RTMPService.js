import { getSessionID } from '../utils'

export default {
  liveURI: (streamName, roomName) => {
    const roomKey = encodeURIComponent(
      JSON.stringify({
        session_id: getSessionID(),
        room_settings: {
          name: roomName,
        },
      })
    )

    return `rtmp://${
      import.meta.env.VITE_RTMP_URL
    }/live/${streamName}?key=${roomKey}`
  },
}
