import { getSessionID } from "../utils";

const rtmpURL = import.meta.env.VITE_RTMP_URL;
const hlsURL = import.meta.env.VITE_HLS_URL;

export const getRtmpUrl = (userID: string) =>
  `rtmp://${rtmpURL}/live/${userID}?key=${getSessionID()}`;

export const getHlsUrl = (userID: string) =>
  `http://${hlsURL}/hls/${userID}.m3u8`;
