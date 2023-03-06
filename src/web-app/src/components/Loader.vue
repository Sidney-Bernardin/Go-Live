<script setup></script>

<template>
  <div class="loader">
    <h1>Loading, please wait...</h1>
  </div>
</template>

<style scoped>
.loader {
  z-index: 1;
  position: fixed;
  top: 0;
  left: 0;
  display: flex;
  width: 100%;
  height: 100%;
  color: #fff;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(2px);
  justify-content: center;
  align-items: center;
}
</style>
