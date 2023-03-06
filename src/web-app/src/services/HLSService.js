import { getSessionID } from '../utils'

export default {
  hlsSrc: (streamName) =>
    `http://${import.meta.env.VITE_HLS_URL}/hls/${streamName}.m3u8?session_id=${getSessionID()}`
}
