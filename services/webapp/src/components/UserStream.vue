<script setup lang="ts">
import HLS from "hls.js"
import { onMounted, ref } from "vue";


const videoElem = ref<HTMLMediaElement>()

onMounted(() => {
    if (!HLS.isSupported()) {
        alert("HLS isn't supported")
        return
    }

    try {
        const hls = new HLS({ debug: true })
        hls.loadSource("http://localhost:8001/hls/abc.m3u8")
        hls.attachMedia(videoElem.value!)
    } catch (e) {
        console.log(e)
    }
})
</script>

<template>
    <video controls ref="videoElem" />
</template>

<style scoped></style>
