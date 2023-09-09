<script setup lang="ts"></script>

<template>
  <div class="footer">
    <div class="left">
      <button>log-out</button>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import "../style.scss";

.footer {
  position: absolute;
  display: flex;
  bottom: 0;
  padding: 15px 0;
  justify-content: space-between;
  align-items: center;

  .left {
    margin-left: 15px;

    button {
      @include basic-button;

      color: $yellow;
      font-size: 1.75rem;
    }
  }
}
</style>
