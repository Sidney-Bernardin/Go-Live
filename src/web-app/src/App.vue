<script setup>
import { ref, provide, onMounted } from 'vue'
import Hls from 'hls.js/dist/hls.min'
import Loader from './components/Loader.vue'
import Navigation from './components/Navigtion.vue'

const loading = ref(false)
const videoElem = ref(null)

provide('loading', loading)

onMounted(() => {
  if (Hls.isSupported()) {
    var hls = new Hls({ debug: true })
    hls.loadSource('http://localhost:8003/hls/stream.m3u8')
    hls.attachMedia(videoElem.value)
    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      console.log('Hello, World!')
      videoElem.value.play()
    })
  }
})
</script>

<template>
  <div>
    <Navigation />

    <!--video controls ref="videoElem" /-->
    <div class="wrapper">
      <router-view></router-view>
    </div>

    <Loader v-if="loading" />
  </div>
</template>

<style scoped>
.wrapper {
  margin-top: 60px;
}
</style>
