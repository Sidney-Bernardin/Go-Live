<script setup lang="ts">
import { useRouter } from "vue-router";
import { User } from "../requests/models";
import { getProfilePictureURL } from "../requests/users";

const router = useRouter();

const props = defineProps<{ user: User }>();
</script>

<template>
  <button class="user-card" @click="router.push({ path: `/${props.user.username}` })">
    <img :src="getProfilePictureURL(props.user.id)" />
    {{ props.user.username }}
  </button>
</template>

<style scoped lang="scss">
@import "../style.scss";

.user-card {
  @include basic-button;
}
</style>
