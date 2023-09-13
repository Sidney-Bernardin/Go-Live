<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useStore } from "vuex";
import Hls from 'hls.js'
import { getHlsUrl } from "../requests/video-stream";

const store = useStore();

const videoElem = ref<HTMLMediaElement | null>(null)

onMounted(() => {
  if (!Hls.isSupported()) return

  const hls = new Hls({ debug: true })
  hls.loadSource(getHlsUrl(store.state.room.id))
  hls.attachMedia(videoElem.value!)
  hls.on(Hls.Events.MANIFEST_PARSED, () => videoElem.value?.play())
})
</script>

<template>
  <div class="video-stream">
    <h2>{{ store.state.room?.name }}</h2>
    <video controls ref="videoElem" />
  </div>
</template>

<style scoped lang="scss">
.video-stream {
  h2 {
    position: absolute;
    font-size: 2.5rem;
    margin: 0;
    transform: translateY(calc(-100% - 5px));
  }

  video {
    width: 100%;
    height: 100%;
  }
}
</style>
