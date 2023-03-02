<script setup>
import { ref, reactive, inject } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import UsersService from '../services/UsersService'

import axios from 'axios'
import { setSessionID } from '../utils'

const store = useStore()
const router = useRouter()

const loading = inject('loading')

const signupForm = reactive({
  username: '',
  email: '',
  password: '',
})

const signinForm = reactive({
  username: '',
  password: '',
})

const signupError = ref('')
const signinError = ref('')

const signup = async () => {
  loading.value = true

  try {
    const res1 = await UsersService.signup(signupForm)
    setSessionID(res1.data.session_id)

    const res2 = await UsersService.getSelf(['username'])
    store.dispatch('setSelf', res2.data)

    router.push({ name: 'Home' })
  } catch (err) {
    if (err.response?.data.type == 'invalid_signup_info')
      signupError.value = err.response.data.detail
    else store.dispatch('handleError', err)
  }

  loading.value = false
}

const signin = async () => {
  loading.value = true

  try {
    const res1 = await UsersService.signin(signinForm)
    setSessionID(res1.data.session_id)

    const res2 = await UsersService.getSelf(['username'])
    store.dispatch('setSelf', res2.data)

    router.push({ name: 'Home' })
  } catch (err) {
    if (err.response?.data.type == 'invalid_signin_info')
      signinError.value = err.response.data.detail
    else store.dispatch('handleError', err)
  }

  loading.value = false
}
</script>

<template>
  <div class="login-page">
    <div class="wrapper">
      <form @submit.prevent="signup">
        <h2>Sign Up</h2>

        <label for="username">Username</label>
        <input
          type="text"
          name="username"
          placeholder="Username"
          v-model="signupForm.username"
        />

        <label for="email">Email</label>
        <input
          type="text"
          name="email"
          placeholder="Email"
          v-model="signupForm.email"
        />

        <label for="password">Password</label>
        <input
          type="password"
          name="password"
          placeholder="Password"
          v-model="signupForm.password"
        />

        <input type="submit" value="Go!" />

        <p v-if="signupError">⚠ {{ signupError }} ⚠</p>
      </form>

      <form @submit.prevent="signin">
        <h2>Sign In</h2>

        <label for="username">Username</label>
        <input
          type="text"
          name="username"
          placeholder="Username"
          v-model="signinForm.username"
        />

        <label for="password">Password</label>
        <input
          type="password"
          name="password"
          placeholder="Password"
          v-model="signinForm.password"
        />

        <input type="submit" value="Go!" />

        <p v-if="signinError">⚠ {{ signinError }} ⚠</p>
      </form>
    </div>
  </div>
</template>

<style scoped>
.wrapper {
  display: flex;
  gap: 15px;
  padding-top: 30px;
  justify-content: center;
}

h2 {
  margin: 0;
}

form {
  display: flex;
  width: 200px;
  border: 2px solid #c1c1c1;
  border-style: dashed;
  padding: 10px;
  flex-direction: column;
}

input[type='submit'],
label {
  margin-top: 5px;
}

p {
  color: red;
  font-weight: bolder;
  text-transform: uppercase;
  text-align: center;
  margin: 5px 0 0 0;
  overflow-wrap: break-word;
}
</style>
