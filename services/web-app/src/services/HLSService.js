import { getSessionID } from '../utils'

const url = import.meta.env.VITE_HLS_URL

export default {
  hlsURI: (streamName) =>
    `http://${url}/hls/${streamName}.m3u8?session_id=${getSessionID()}`,
}
