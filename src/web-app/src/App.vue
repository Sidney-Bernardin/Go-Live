<script setup>
import { ref, provide, onMounted } from 'vue'
import { useStore } from 'vuex'

import Loader from './components/Loader.vue'
import Navigation from './components/Navigtion.vue'

import Hls from 'hls.js/dist/hls.min'
import UsersService from './services/UsersService'
import { removeSessionID } from './utils'

const store = useStore()

const loading = ref(false)
const videoElem = ref(null)

provide('loading', loading)

UsersService.getSelf(['username'])
  .then((res) => store.dispatch('setSelf', res.data))
  .catch((err) => {
    if (err.response?.data.type == 'unauthorized') {
      removeSessionID()
      store.dispatch('setSelf', null)
      return
    }

    store.dispatch('handleError', err)
  })

onMounted(() => {
  if (Hls.isSupported()) {
    var hls = new Hls({ debug: true })
    hls.loadSource(
      'http://localhost:8003/hls/stream.m3u8?session_id=6400e4784ada72808818b2bc'
    )
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

    <video controls ref="videoElem" style="width: 500px; height: 500px" />

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
