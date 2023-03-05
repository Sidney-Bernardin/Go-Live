<script setup>
import { ref, onMounted } from 'vue'

import Hls from 'hls.js/dist/hls.min'
import { getSessionID } from '../utils'

const videoElem = ref(null)

onMounted(() => {
  if (Hls.isSupported()) {
    var hls = new Hls({ debug: true })
    hls.loadSource(`http://localhost:8003/hls/64037bc46117a8fdc471c823.m3u8?session_id=${getSessionID()}`)
    hls.attachMedia(videoElem.value)
    hls.on(Hls.Events.MANIFEST_PARSED, () => videoElem.value.play())
  }
})
</script>

<template>
  <div class="home-page">
    <video controls ref="videoElem" />
  </div>
</template>

<style scoped>
video {
  width: 300px;
  height: 300px;
}
</style>
