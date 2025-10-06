<script setup lang="ts">
import HLS from "hls.js"
import { onMounted, ref } from "vue";


const videoElem = ref<HTMLMediaElement>()

onMounted(() => {
    if (!HLS.isSupported()) {
        alert("HLS isn't supported")
        return
    }

    const hls = new HLS({ debug: true })
    hls.loadSource("http://localhost:3011/hls/abc.m3u8")
    hls.attachMedia(videoElem.value!)
})
</script>

<template>
    <video controls ref="videoElem" />
</template>

<style scoped></style>
