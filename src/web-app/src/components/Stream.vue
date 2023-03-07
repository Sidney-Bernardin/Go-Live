<script setup>
import { ref, onMounted } from 'vue'
import { useStore } from 'vuex'

import Hls from 'hls.js/dist/hls.min'
import HLSService from '../services/HLSService'

const store = useStore()

const videoElem = ref(null)

const goLive = () => (videoElem.value.currentTime = videoElem.value.duration)

onMounted(() => {
  if (!Hls.isSupported()) return

  var hls = new Hls({ debug: true })
  hls.loadSource(HLSService.hlsSrc(store.state.currentRoomID))
  hls.attachMedia(videoElem.value)
  hls.on(Hls.Events.MANIFEST_PARSED, () => videoElem.value?.play())
})
</script>

<template>
  <video class="stream" controls muted @play="goLive" ref="videoElem"></video>
</template>

<style scoped></style>
