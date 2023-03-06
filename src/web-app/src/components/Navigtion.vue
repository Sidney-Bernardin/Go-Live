<script setup>
import { ref, computed, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

import SearchBar from './SearchBar.vue'

import UsersService from '../services/UsersService'
import { removeSessionID } from '../utils'

const router = useRouter()
const store = useStore()

const loading = inject('loading')

const logout = async () => {
  loading.value = true

  await UsersService.logout()
    .then(() => {
      removeSessionID()
      store.dispatch('setSelf', null)
      router.push({ name: 'Login' })
    })
    .catch((err) => store.dispatch('handleError', err))

  loading.value = false
}
</script>

<template>
  <div class="navigation">
    <p class="error-message" v-if="store.state.errorMessage">
      ⚠ {{ store.state.errorMessage }} ⚠
    </p>

    <ul class="links">
      <li><router-link to="/">Home</router-link></li>
    </ul>

    <SearchBar v-if="store.state.self" />

    <button
      v-if="store.state.self"
      class="profile-picture"
      :style="`background: url(${UsersService.profilePictureSrc(store.state.self.id)}) center/100%`"
    ></button>

    <div class="dropdown" v-if="store.state.self">
      <ul>
        <li>
          <router-link :to="store.state.self.username">Profile</router-link>
        </li>
      </ul>

      <ul>
        <li><button @click="logout">Logout</button></li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.navigation {
  position: fixed;
  display: flex;
  top: 0;
  left: 0;
  z-index: 1;
  width: 100%;
  height: 60px;
  border-bottom: 6px solid #c1c1c1;
  border-bottom-style: double;
  background-color: #fff;
  justify-content: space-between;
  align-items: center;
}

.error-message {
  position: absolute;
  top: 66px;
  width: 100%;
  color: red;
  font-size: 12px;
  font-weight: bolder;
  text-transform: uppercase;
  text-align: center;
  line-height: 15px;
  margin: 0;
}

ul.links {
  display: flex;
  width: 51px;
  gap: 5px;
  list-style-type: none;
  margin: 0 0 0 10px;
  padding: 0;
}

button.profile-picture {
  width: 45px;
  height: 45px;
  border: 3px solid #c1c1c1;
  border-style: inset;
  margin-right: 6px;
  padding: 0;
}

.dropdown {
  position: fixed;
  display: none;
  top: 75px;
  right: 0;
  border: 6px solid #c1c1c1;
  border-style: double;
  border-right: none;
  gap: 30px;
  background-color: #fff;
  padding: 10px;
  flex-direction: column;
  justify-content: space-between;
}

button.profile-picture:focus + .dropdown,
.dropdown:hover {
  display: flex;
}

.dropdown ul {
  display: flex;
  gap: 5px;
  list-style-type: none;
  margin: 0;
  padding: 0;
  flex-direction: column;
}

.dropdown ul button {
  width: 100%;
}
</style>
