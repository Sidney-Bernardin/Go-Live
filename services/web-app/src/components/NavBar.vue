<script setup lang="ts">
import { useRouter } from "vue-router";
import { useStore } from "vuex";
import UserCard from "./UserCard.vue";
import { getRtmpUrl } from "../requests/video-stream";

const router = useRouter();
const store = useStore();

const onGoLive = (): void => {
  const url = getRtmpUrl(store.state.self.id)
  navigator.clipboard.writeText(url);
}
</script>

<template>
  <div class="nav-bar">
    <div class="left">
      <UserCard v-if="store.state.self" :user="store.state.self" />
      <div class="logo" v-else>Go Live</div>

      <div class="destination" v-if="store.state.self">
        <span>/</span>
        {{ router.currentRoute.value.params.username }}
      </div>
    </div>

    <div class="right" v-if="store.state.self">
      <button class="go-live" @click="onGoLive">go live</button>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import "../style.scss";

.nav-bar {
  position: absolute;
  display: flex;
  padding: 15px 0;
  justify-content: space-between;
  align-items: center;

  button {
    font-size: 1.75rem;
  }

  .left {
    display: flex;
    gap: 10px;
    margin-left: 15px;
    align-items: center;

    .logo {
      @include basic-button;
      color: $dark-green;
      font-size: 2rem;
      text-decoration: none;
    }

    .destination {
      display: flex;
      gap: 10px;
      font-size: 1.5rem;
      font-weight: bolder;
      align-items: center;

      span {
        font-size: 2rem;
        font-weight: bolder;
      }
    }
  }

  .right {
    margin-right: 15px;

    .go-live {
      @include basic-button;
    }
  }
}
</style>
