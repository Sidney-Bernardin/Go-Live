<script setup lang="ts">
const props = defineProps({ username: String })
</script>

<template>
  <button class="user-card">
    <img
      src="https://b1512865.smushcdn.com/1512865/wp-content/uploads/2020/04/stormlight-simplifed_v02_dk.jpg?lossy=1&strip=1&webp=1">
    {{ props.username }}
  </button>
</template>

<style scoped lang="scss">
@import "../style.scss";

.user-card {
  position: relative;
  display: flex;
  cursor: pointer;
  gap: 5px;
  border: none;
  background: transparent;
  font-size: 1.75rem;
  font-weight: bolder;
  padding: 0;
  align-items: center;

  img {
    height: 1.75rem;
  }
}
</style>
