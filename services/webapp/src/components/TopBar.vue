<script setup lang="ts">
import { useStore } from 'vuex';

const store = useStore()
</script>

<template>
    <div>
        <router-link to="/login" v-if="!store.state.self">/login</router-link>
        <div v-else>
            <img :src="`http://localhost:4566/users/profiles/pictures/${store.state.self.id}.png`" alt="">
            <router-link :to="`/${store.state.self.id}`">/{{ store.state.self.username }}</router-link>
        </div>

    </div>
</template>
