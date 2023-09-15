<script setup lang="ts">
import { useRouter } from "vue-router";
import { useStore } from "vuex";

import NavBar from "./components/NavBar.vue";
import Footer from "./components/Footer.vue";
import Chat from "./components/Chat.vue";
import Explore from "./components/Explore.vue";

import { authenticateUser } from "./requests/users";
import { unexpectedErr, deleteSessionID } from "./utils";

const router = useRouter();
const store = useStore();

router.beforeEach(async () => {
  const res = await authenticateUser(["username"]).catch((err) => {
    if (err.response.data.problem == "unauthorized") deleteSessionID();
    else unexpectedErr(err);
  });

  store.dispatch("setSelf", res);
});
</script>

<template>
  <div class="wrapper">
    <div class="view">
      <NavBar />
      <router-view></router-view>
      <Footer v-if="store.state.self" />
    </div>

    <div class="right">
      <Chat />
      <Explore />
    </div>
  </div>
</template>

<style scoped lang="scss">
@import "./style.scss";

.wrapper {
  display: flex;
  width: 100%;
  height: 100vh;

  .view {
    position: relative;
    background: $green;
    flex: 1;
    overflow: hidden;

    .nav-bar {
      width: 100%;
    }
  }

  .right {
    display: flex;
    width: 50%;
    flex: 1;
    flex-direction: column;

    .chat,
    .explore {
      transition: 0.2s;

      &.primary {
        flex: 3;
      }

      &.secondary {
        flex: 1;

        &:hover {
          flex: 9;
        }
      }

      &.disabled {
        flex: 1;

        &:hover {
          flex: 1;
        }
      }
    }
  }
}
</style>
