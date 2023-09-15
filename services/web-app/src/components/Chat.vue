<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useStore } from "vuex";

import UserCard from "./UserCard.vue";

import { wsMessage, RoomEvent } from "../requests/models";
import { joinRoom } from "../requests/rooms";
import { unexpectedErr } from "../utils";

const store = useStore();

const state = computed(() => (store.state.room ? "primary" : "disabled"));
const room = computed(() => store.state.room);
const chatMessages = ref<RoomEvent[]>([]);

const onChat = (e: Event): void =>
  ws.value?.send(
    JSON.stringify({
      type: "CHAT",
      user_id: store.state.self.id,
      username: store.state.self.username,
      message: Object.fromEntries(new FormData(e.target as HTMLFormElement)).message,
    } as RoomEvent),
  );

var ws = ref<WebSocket | null>(null);

watch(room, (newRoom) => {
  if (!newRoom) return;

  ws.value = joinRoom(store.state.room.id);
  ws.value.onerror = (err) => unexpectedErr(err);
  ws.value.onclose = () => store.dispatch("setRoom", null)
  ws.value.onmessage = (msg) => {
    const wsMsg = JSON.parse(msg.data) as wsMessage<RoomEvent>;
    if (wsMsg.content.type == "CHAT") chatMessages.value.push(wsMsg.content);
  };
});
</script>

<template>
  <div :class="`chat ${state}`">
    <ul>
      <li v-for="msg in chatMessages">
        <UserCard :user="{id: msg.user_id, username: msg.username}" />        
        <p>{{ msg.message }}</p>
      </li>
    </ul>

    <div class="cover">
      <h1>Go Chat</h1>

      <form @submit.prevent="onChat">
        <input type="text" name="message" placeholder="Send a message!" />
        <input type="submit" />
      </form>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import "../style.scss";

.chat {
  position: relative;
  overflow: scroll;
  color: $white;
  background: $black;

  &.disabled {
    pointer-events: none;

    form {
      display: none;
    }
  }

  .cover {
    position: sticky;
    display: flex;
    left: 0;
    bottom: 0;
    gap: 30px;
    height: calc(100vh / 4);
    background: $black;
    padding: 0 30px;
    align-items: center;

    h1 {
      margin: 0;
      color: $dark-black;
      font-size: 3.5rem;
      text-wrap: nowrap;
    }

    form {
      width: 100%;

      input {
        box-sizing: border-box;
        width: 100%;
        height: 2rem;
        border: 2px solid $white;
        color: $white;
        font-size: 2rem;
        font-weight: bolder;
        background: transparent;
        padding: 30px 15px;

        &::placeholder {
          color: $dark-black;
        }

        &[type="submit"] {
          display: none;
        }
      }
    }
  }

  ul {
    position: relative;
    display: flex;
    gap: 15px;
    border: 2px solid $white;
    margin: 30px 30px 0 30px;
    padding: 15px;
    flex: 5;
    flex-direction: column;
    list-style-type: none;

    li {
      display: flex;
      gap: 15px;
      font-size: 1.75rem;
      align-items: normal;

      .user-card {
        font-size: 1.75rem;
        align-self: start;  
      }

      p {
        margin: 0;
        word-break: break-all;
      }
    }
  }
}
</style>
