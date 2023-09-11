<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { useStore } from "vuex";

import { searchUsers } from "../requests/users";
import { unexpectedErr } from "../utils";

const router = useRouter();
const store = useStore();

const user = ref({});

watch(
  router.currentRoute,
  (newRoute) =>
    searchUsers(newRoute.params.username as string, ["username"])
      .then((res) => {
        if (res.data[0] != newRoute.params.username)
          router.replace({ path: `/${store.state.self.username}` });
        else user.value = res.data[0];
      })
      .catch((err) => unexpectedErr(err)),
  { immediate: true },
);
</script>

<template>
  <div class="user-page">
    <h1>User</h1>
  </div>
</template>

<style scoped lang="scss">
@import "../style.scss";

.user-page {
  background: $green;
}
</style>
