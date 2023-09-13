<script setup lang="ts">
import { computed } from "vue";
import { useStore } from "vuex";

const store = useStore();

const state = computed(() => (store.state.room ? "secondary" : "primary"));
</script>

<template>
  <div :class="`explore ${state}`">
    <div class="cover">
      <h1>Go Explore</h1>
      <input type="text" placeholder="Search for users!" />
    </div>

    <p>
      Lorem ipsum, dolor sit amet consectetur adipisicing elit. Amet, aliquid
      accusamus iste veniam optio laborum ipsam, ab harum dicta, quis labore
      impedit dolores reiciendis at ducimus a ad quo. Modi?
    </p>
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
      flex: 2;

      &::placeholder {
        color: $dark-white;
      }
    }
  }

  p {
    border: 2px solid $black;
    margin: 0 30px 30px 30px;
    padding: 15px;
    font-size: 5rem;
    flex: 5;
  }
}
</style>
