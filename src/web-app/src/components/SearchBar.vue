<script setup>
import { ref } from 'vue'
import UsersService from '../services/UsersService'

const input = ref('')
const results = ref([])

const search = () =>
  UsersService.searchUsers(input.value, 0, 20, ['username'])
    .then((res) => (results.value = res.data))
    .catch((err) => store.dispatch('handleError', err))
</script>

<template>
  <div class="search-bar">
    <input
      type="text"
      placeholder="Search for Users"
      @input="search"
      v-model="input"
    />

    <ul>
      <li v-for="(res, idx) in results" :key="idx">
        <router-link :to="res.username">{{ res.username }}</router-link>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.search-bar {
  position: relative;
  display: flex;
  background-color: #fff;
  flex-direction: column;
  align-items: center;
}

ul {
  position: absolute;
  display: none;
  top: 150%;
  min-width: 100%;
  border: 6px solid #c1c1c1;
  border-style: double;
  gap: 5px;
  background: #fff;
  margin: 0;
  padding: 10px;
  list-style-type: none;
  flex-direction: column;
  justify-content: space-between;
}

input:focus + ul,
ul:hover {
  display: flex;
}

li {
  margin: 0;
  padding: 0;
}
</style>
