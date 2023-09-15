<script setup lang="ts">
import { ref, computed } from "vue";
import { useStore } from "vuex";

import UserCard from "./UserCard.vue";

import { User } from "../requests/models";
import { searchUsers } from "../requests/users";
import { unexpectedErr } from "../utils";

const store = useStore();

const state = computed(() => (store.state.room ? "secondary" : "primary"));
const users = ref<User[]>([]);

const onExplore = (e: Event) =>
  searchUsers((e.target as HTMLInputElement).value, ["username"])
    .then((res) => (users.value = res))
    .catch((err) => unexpectedErr(err));
</script>

<template>
  <div :class="`explore ${state}`">
    <div class="cover">
      <h1>Go Explore</h1>
      <input type="text" placeholder="Search for users!" @input="onExplore" />
    </div>

    <ul>
      <li v-for="user in users">
        <UserCard :user="user" />
      </li>
    </ul>
  </div>
</template>

<style scoped lang="scss">
@import "../style.scss";

.explore {
  position: relative;
  overflow: scroll;
  color: $black;
  background: $white;

  &.disabled {
    pointer-events: none;

    input {
      display: none;
    }
  }

  .cover {
    z-index: 1;
    position: sticky;
    display: flex;
    top: 0;
    left: 0;
    gap: 30px;
    height: calc(100vh / 4);
    background: $white;
    padding: 0 30px;
    align-items: center;

    h1 {
      margin: 0;
      color: $dark-white;
      font-size: 3.5rem;
      text-wrap: nowrap;
      flex: 1;
    }

    input {
      width: 100%;
      height: 2rem;
      border: 2px solid $black;
      color: $black;
      font-size: 2rem;
      font-weight: bolder;
      background: transparent;
      padding: 15px;
  
      &::placeholder {
        color: $dark-white;
      }
    }
  }

  ul {
    position: relative;
    display: flex;
    gap: 15px;
    border: 2px solid $black;
    margin: 0 30px 30px 30px;
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
        color: $black;
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
