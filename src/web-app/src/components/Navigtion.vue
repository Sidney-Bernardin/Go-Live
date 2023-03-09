<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

import SearchBar from './SearchBar.vue'
import ProfilePicture from './ProfilePicture.vue'

import UsersService from '../services/UsersService'
import { removeSessionID } from '../utils'

const router = useRouter()
const store = useStore()

const loading = inject('loading')
const showURICreator = inject('show_URI_creator')

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

    <ProfilePicture v-if="store.state.self" :userID="store.state.self.id" />

    <div class="dropdown" v-if="store.state.self">
      <ul>
        <li>
          <router-link :to="store.state.self.username">Profile</router-link>
        </li>
      </ul>

      <ul>
        <li><button @click="showURICreator = true">Go Live</button></li>
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
  background: #fff;
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

.profile-picture {
  margin-right: 6px;
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
  background: #fff;
  padding: 10px;
  flex-direction: column;
  justify-content: space-between;
}

.profile-picture:focus + .dropdown,
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
