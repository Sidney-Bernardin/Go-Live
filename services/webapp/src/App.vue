<script setup lang="ts">
import TopBar from "./components/TopBar.vue"
import BottomBar from "./components/BottomBar.vue"
import Chat from "./components/Chat.vue"
import Search from "./components/Search.vue"
</script>

<template>
    <div class="app">
        <div class="left">
            <TopBar />
            <router-view v-slot="{ Component }">
                <keep-alive>
                    <component :is="Component"></component>
                </keep-alive>
            </router-view>
            <BottomBar />
        </div>

        <div class="right">
            <Chat />
            <Search />
        </div>
    </div>
</template>

<style scoped>
.app {
    display: flex;
    width: 100%;
    height: 100vh;
}

.left {
    flex: 1;
    background: var(--green);
}

.left .page {
    height: 100%;
}

.right {
    flex: 1;
    display: flex;
    flex-direction: column;
    width: 50%;
}

.chat {
    flex: 1;
    transition: 0.2s;
}

.search {
    flex: 4;
    transition: 0.2s;
}
</style>
