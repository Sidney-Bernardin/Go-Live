<script setup lang="ts">
import { useRouter } from "vue-router";
import { logout } from "../requests/users";
import { loader, deleteSessionID } from "../utils";

const router = useRouter();
const { loading, wrapLoad } = loader();

const onLogout = (): Promise<void> =>
  wrapLoad(
    logout().then(() => {
      deleteSessionID();
      router.push({ name: "Login" });
    }),
  );
</script>

<template>
  <div class="footer">
    <div class="left">
      <button @click="onLogout" :disabled="loading">log-out</button>
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
