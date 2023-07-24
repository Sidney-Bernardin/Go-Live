<script setup>
import { ref, provide, onMounted } from 'vue'
import { useStore } from 'vuex'

import Loader from './components/Loader.vue'
import URICreator from './components/URICreator.vue'
import Navigation from './components/Navigtion.vue'

import UsersService from './services/UsersService'
import { removeSessionID } from './utils'

const store = useStore()

const loading = ref(false)
const showURICreator = ref(false)

provide('loading', loading)
provide('show_URI_creator', showURICreator)

UsersService.getSelf(['username'])
  .then((res) => store.dispatch('setSelf', res.data))
  .catch((err) => {
    if (err.response?.data.type == 'unauthorized') {
      removeSessionID()
      store.dispatch('setSelf', null)
      return
    }

    store.dispatch('handleError', err)
  })
</script>

<template>
  <div>
    <Navigation />

    <div class="wrapper">
      <router-view></router-view>
    </div>

    <Loader v-if="loading" />
    <URICreator v-if="showURICreator" />
  </div>
</template>

<style scoped>
.wrapper {
  margin: 81px 15px 15px 15px;
}
</style>
