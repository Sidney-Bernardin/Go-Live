<script setup lang="ts">
import UserCard from "./UserCard.vue";
</script>

<template>
  <div class="nav-bar">
    <div class="left">
      <UserCard username="profile" />
      <span>/</span>
      <div class="destination">bob-smith</div>
    </div>

    <div class="right">
      <button class="go-live">Go Live</button>
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

    span {
      font-size: 2rem;
      font-weight: bolder;
    }

    .destination {
      background: transparent;
      font-size: 1.5rem;
      font-weight: bolder;
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
