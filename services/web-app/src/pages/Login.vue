<script setup lang="ts">
import { ref } from "vue";

const showSignup = ref(true);
</script>

<template>
  <div class="login-page">
    <div class="tabs">
      <button :disabled="!showSignup" @click="showSignup = false">
        sign-in
      </button>

      <button :disabled="showSignup" @click="showSignup = true">sign-up</button>
    </div>

    <form v-if="showSignup">
      <label class="file-label">
        <input type="file" />
        <img
          src="https://b1512865.smushcdn.com/1512865/wp-content/uploads/2020/04/stormlight-simplifed_v02_dk.jpg?lossy=1&strip=1&webp=1"
          alt=""
        />
        Upload Profile Picture
      </label>

      <label for="username">Username</label>
      <input type="text" name="username" placeholder="Username" />

      <label for="email">Email</label>
      <input type="email" placeholder="Email" />

      <label for="password">Password</label>
      <input type="password" placeholder="Password" />

      <input type="submit" value="sign-up" />

      <p>Bad username or whatever...</p>
    </form>

    <form v-else>
      <label for="username">Username</label>
      <input type="text" name="username" placeholder="Username" />

      <label for="password">Password</label>
      <input type="password" placeholder="Password" />

      <input type="submit" value="sign-in" />

      <p>Bad username or whatever...</p>
    </form>
  </div>
</template>

<style scoped lang="scss">
@import "../style.scss";

.login-page {
  display: flex;
  height: 100%;
  justify-content: center;
  align-items: center;

  .tabs {
    position: absolute;
    display: flex;
    top: 50%;
    left: 30px;
    gap: 15px;
    width: 0;
    rotate: -90deg;
    transform: translate(-50%, 50%);
    justify-content: center;

    button {
      @include basic-button;

      font-size: 1.5rem;
      text-wrap: nowrap;
    }
  }

  form {
    display: flex;
    border: 1px solid $black;
    padding: 30px;
    flex-direction: column;

    p {
      color: $yellow;
      font-size: 1.5rem;
      font-weight: bolder;
      margin: 30px 0 0 0;
    }

    label {
      font-size: 1.5rem;
      font-weight: bolder;
      margin-top: 15px;

      &.file-label {
        @include basic-button;
        border: 1px solid $black;
        background: transparent;
        margin-top: 0;
        padding: 15px;
      }
    }

    input {
      border: 1px solid $black;
      color: $white;
      background: transparent;
      font-size: 1.5rem;
      font-weight: bolder;
      padding: 15px;

      &::placeholder {
        color: $dark-green;
      }

      &[type="file"] {
        display: none;
      }

      &[type="submit"] {
        cursor: pointer;
        border: 1px solid $white;
        margin-top: 30px;
        text-transform: uppercase;

        &[disabled="true"] {
          border: 1px solid $dark-green;
          color: $dark-green;
          text-decoration: none;
        }
      }
    }
  }
}
</style>
