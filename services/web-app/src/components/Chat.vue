<script setup lang="ts">
import { computed } from "vue";
import { useStore } from "vuex";

const store = useStore();

const state = computed(() => (store.state.room ? "primary" : "disabled"));
</script>

<template>
  <div :class="`chat ${state}`">
    <p>
      Lorem ipsum, dolor sit amet consectetur adipisicing elit. Amet, aliquid
      accusamus iste veniam optio laborum ipsam, ab harum dicta, quis labore
      impedit dolores reiciendis at ducimus a ad quo. Modi?
    </p>

    <div class="cover">
      <h1>Go Chat</h1>
      <input type="text" placeholder="Send a message!" />
    </div>
  </div>
</template>

<style lang="scss">
@import "../style.scss";

.chat {
  position: relative;
  overflow: scroll;
  color: $white;
  background: $black;

  &.disabled {
    pointer-events: none;

    input {
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
      flex: 1;
    }

    input {
      width: 100%;
      height: 2rem;
      border: 2px solid $white;
      color: $white;
      font-size: 2rem;
      font-weight: bolder;
      background: transparent;
      padding: 15px;
      flex: 2;

      &::placeholder {
        color: $dark-black;
      }
    }
  }

  p {
    border: 2px solid $white;
    margin: 30px 30px 0 30px;
    padding: 15px;
    font-size: 5rem;
    flex: 5;
  }
}
</style>
