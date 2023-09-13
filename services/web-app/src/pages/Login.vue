<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useStore } from "vuex";

import { loader, setSessionID, unexpectedErr } from "../utils";
import { LoginRes } from "../requests/models";
import { signup, signin, getUser } from "../requests/users";

const store = useStore();
const router = useRouter();
const { loading, wrapLoad } = loader();

const showSignup = ref(true);
const profilePictureURL = ref("");
const signupProblem = ref("");
const signinProblem = ref("");

const onProfilePictureChange = (e: Event): void => {
  const file: File = (e.target as HTMLInputElement).files![0];
  profilePictureURL.value = URL.createObjectURL(file);
};

const finishLogin = ({ session_id, user_id }: LoginRes): Promise<void> =>
  getUser(user_id, ["username"])
    .then((res) => {
      store.dispatch("setSelf", res)
      setSessionID(session_id);
      router.push({ name: "User", params: { username: res.username } });
    })
    .catch((err) => unexpectedErr(err))

const onSubmitSignup = (e: Event): Promise<void> =>
  wrapLoad(
    signup(new FormData(e.target as HTMLFormElement))
      .then((res) => finishLogin(res))
      .catch((err) => {
        if (err.response.data.problem == "invalid_signup_info")
          signupProblem.value = err.response.data.detail;
        else unexpectedErr(err);
      }),
  );

const onSubmitSignin = (e: Event): Promise<void> =>
  wrapLoad(
    signin(new FormData(e.target as HTMLFormElement))
      .then((res) => finishLogin(res))
      .catch((err) => {
        if (err.response.data.problem == "invalid_signin_info")
          signinProblem.value = err.response.data.detail;
        else unexpectedErr(err);
      }),
  );
</script>

<template>
  <div class="login-page">
    <div class="tabs">
      <button :disabled="!showSignup" @click="showSignup = false">
        sign-in
      </button>

      <button :disabled="showSignup" @click="showSignup = true">sign-up</button>
    </div>

    <form v-if="showSignup" @submit.prevent="onSubmitSignup">
      <label class="file-label">
        <input type="file" name="profile_picture" @change="onProfilePictureChange" />
        <img :src="profilePictureURL" />
        Upload Profile Picture
      </label>

      <label for="username">Username</label>
      <input type="text" name="username" placeholder="Username" />

      <label for="email">Email</label>
      <input type="email" name="email" placeholder="Email" />

      <label for="password">Password</label>
      <input type="password" name="password" placeholder="Password" />

      <input type="submit" value="sign-up" :disabled="loading" />

      <p v-if="signupProblem">{{ signupProblem }}</p>
    </form>

    <form v-else @submit.prevent="onSubmitSignin">
      <label for="username">Username</label>
      <input type="text" name="username" placeholder="Username" />

      <label for="password">Password</label>
      <input type="password" name="password" placeholder="Password" />

      <input type="submit" value="sign-in" :disabled="loading" />

      <p v-if="signinProblem">{{ signinProblem }}</p>
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
    width: 450px;
    border: 1px solid $black;
    padding: 30px;
    flex-direction: column;

    p {
      width: 100%;
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

        &[disabled] {
          cursor: auto;
          border: 1px solid $dark-green;
          color: $dark-green;
          text-decoration: none;
        }
      }
    }
  }
}
</style>
