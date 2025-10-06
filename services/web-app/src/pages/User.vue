<script setup lang="ts">
import { onBeforeRouteUpdate } from 'vue-router';
import { useStore } from 'vuex';

const store = useStore()

onBeforeRouteUpdate(async (to, _from) => {
    await store.dispatch("loadUser", to.params.userID)
})
</script>

<template>
    {{ store.state.user }}
</template>
