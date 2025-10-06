<script setup lang="ts">
import TopBar from "./components/TopBar.vue"
import BottomBar from "./components/BottomBar.vue"
import Chat from "./components/Chat.vue"
import Search from "./components/Search.vue"
</script>

<template>
    <div class="left">
        <TopBar />
        <RouterView />
        <BottomBar />
    </div>

    <div class="right">
        <Chat />
        <Search />
    </div>
</template>

<style scoped></style>
