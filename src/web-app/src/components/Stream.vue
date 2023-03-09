<script setup>
import { ref, onMounted } from 'vue'
import { useStore } from 'vuex'

import Hls from 'hls.js/dist/hls.min'
import HLSService from '../services/HLSService'

const store = useStore()

const videoElem = ref(null)
const live = ref(false)

const setLive = () =>
  (live.value = videoElem.value?.currentTime > videoElem.value?.duration - 10)
const goLive = () =>
  (videoElem.value.currentTime = videoElem.value.duration - 1)

onMounted(() => {
  if (!Hls.isSupported()) return

  var hls = new Hls({ debug: true })
  hls.loadSource(HLSService.hlsURI(store.state.currentRoomID))
  hls.attachMedia(videoElem.value)
  hls.on(Hls.Events.MANIFEST_PARSED, () => videoElem.value?.play())
})
</script>

<template>
  <div class="stream">
    <video controls muted @timeupdate="setLive" ref="videoElem" />

    <div class="controls">
      <button v-if="live" @click="goLive">🔴 live</button>
      <button v-else @click="goLive">⚫ live</button>
    </div>
  </div>
</template>

<style scoped>
.stream {
  position: relative;
  background: #000;
}

.stream:hover .controls {
  opacity: 1;
}

video {
  width: 100%;
  height: 100%;
}

.controls {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  opacity: 0;
  padding: 15px;
  transition: 0.2s;
}

.controls * {
  pointer-events: auto;
}

button {
  text-transform: uppercase;
}
</style>
