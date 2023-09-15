<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter, RouteLocationNormalizedLoaded } from "vue-router";
import { useStore } from "vuex";

import VideoStream from "../components/VideoStream.vue";

import { User } from "../requests/models";
import { searchUsers } from "../requests/users";
import { getRoom } from "../requests/rooms";
import { unexpectedErr } from "../utils";

const router = useRouter();
const store = useStore();

const user = ref<User | null>(null);

const onRouteChange = async (newRoute: RouteLocationNormalizedLoaded): Promise<void> => {
  if (newRoute.params.username == "_") {
    router.push({ path: `/${store.state.self.username}` });
    return;
  }

  user.value = null;
  store.dispatch("setRoom", null);

  try {
    const users = await searchUsers(newRoute.params.username as string, ["username"]);
    if (users.length == 0) return;
    else user.value = users[0];

    const room = await getRoom(user.value!.id);
    store.dispatch("setRoom", room);
  } catch (err: any) {
    if (err.response?.data.problem == "room_doesnt_exist") return;
    unexpectedErr(err);
  }
};

watch(router.currentRoute, onRouteChange, { immediate: true });
</script>

<template>
  <div class="user-page">
    <h1 v-if="!user">User Not Found</h1>
    <h1 v-else-if="!store.state.room">Offline</h1>
    <VideoStream v-else />
  </div>
</template>

<style scoped lang="scss">
@import "../style.scss";

.user-page {
  display: flex;
  height: 100%;
  justify-content: center;
  align-items: center;

  h1 {
    max-width: 450px;
    border: 2px solid $dark-green;
    color: $dark-green;
    font-size: 5rem;
    margin: 0;
    padding: 60px;
    text-align: center;
  }

  .video-stream {
    width: calc(100% - 60px);
  }
}
</style>
