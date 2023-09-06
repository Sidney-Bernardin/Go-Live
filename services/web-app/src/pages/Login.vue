<script setup lang="ts"></script>

<template>
  <div class="login-page">
    <div class="tabs">
      <button class="signin-btn">sign-in</button>
      <button class="signup-btn active">sign-up</button>
    </div>

    <form>
      <label for="username">Username</label>
      <input type="text" name="username" placeholder="Username" />

      <label for="email">Email</label>
      <input type="email" placeholder="Email" />

      <label for="password">Password</label>
      <input type="password" placeholder="Password" />

      <input type="submit" value="sign-up" />
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
    transform: translate(-50%, 50%);
    justify-content: center;
    rotate: -90deg;

    button {
      cursor: pointer;
      border: none;
      color: $white;
      background: transparent;
      font-size: 1.5rem;
      font-weight: bolder;
      padding: 0;
      text-wrap: nowrap;
      text-transform: uppercase;
      text-decoration: underline;

      &.active {
        cursor: auto;
        color: $dark-green;
        text-decoration: none;
      }
    }
  }

  form {
    display: flex;
    border: 1px solid $black;
    padding: 30px;
    flex-direction: column;

    label {
      font-size: 1.25rem;
      font-weight: bolder;
      margin-top: 15px;
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

      &[type="submit"] {
        cursor: pointer;
        border: 1px solid $white;
        margin-top: 30px;
        text-transform: uppercase;
        text-decoration: underline;

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
