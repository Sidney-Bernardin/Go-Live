<script setup lang="ts">
import UserStream from "../components/UserStream.vue";

import { onMounted, ref } from 'vue';
import { onBeforeRouteUpdate, useRoute } from 'vue-router';
import { useStore } from 'vuex';

const store = useStore()
const route = useRoute()

const user = ref<any>(undefined)

async function setUser(userID: string) {
    if (userID === store.state.self.id) {
        user.value = store.state.self
        return
    }

    const res = await fetch(`http://localhost:8000/users/users/${userID}`)
    const json = await res.json()
    if (!res.ok) {
        if (res.status !== 404 && json.type !== "uuid_invalid") {
            console.error(json)
            alert("Something went wrong!")
        }

        user.value = undefined
        return
    }

    user.value = json
}

onMounted(() => setUser(route.params.userID as string))
onBeforeRouteUpdate((to) => setUser(to.params.userID as string))
</script>

<template>
    <div class="user page">
        <h1 class="error wrapper" v-if="!user">user doesn't exist</h1>
        <div class="wrapper" v-else>
            <UserStream />
            <div class="info">
                <h2 class="username">{{ user.username }}</h2>
                <p class="viewers">viewers: -1000000</p>
            </div>
        </div>
    </div>
</template>

<style scoped>
.user.page {
    display: flex;
    justify-content: center;
    align-items: center;
}

.wrapper {
    border: 2px solid var(--dark-green);
    padding: 3rem;
}

.wrapper.error {
    border: 2px solid var(--dark-green);
    color: var(--dark-green);
    font-size: 5rem;
}

video {
    min-width: 30rem;
    min-height: 20rem;
}

.info {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.info .username {
    font-size: 2rem;
}

.info .viewers {
    font-size: 2rem;
}
</style>
