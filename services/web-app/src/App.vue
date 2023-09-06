<script setup lang="ts">
import NavBar from "./components/NavBar.vue";
import Chat from "./components/Chat.vue";
import Explore from "./components/Explore.vue";
</script>

<template>
  <div class="wrapper">
    <div class="view">
      <NavBar />
      <router-view></router-view>
    </div>

    <div class="right">
      <Chat />
      <Explore class="secondary" />
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

    .nav-bar {
      width: 100%;
    }
  }

  .right {
    display: flex;
    width: 50%;
    flex: 1;
    flex-direction: column;

    .chat {
      flex: 3;
      transition: 0.2s;

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

    .explore {
      flex: 3;
      transition: 0.2s;

      &.secondary {
        flex: 1;

        &:hover {
          flex: 9;
        }
      }
    }
  }
}
</style>
