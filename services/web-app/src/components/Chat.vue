<script setup>
import { ref } from 'vue'
import { useStore } from 'vuex'

import ProfilePicture from './ProfilePicture.vue'
import UsersService from '../services/UsersService'

const store = useStore()

const messageText = ref('')

const chat = () =>
  store.dispatch('sendMessage', {
    type: 'CHAT',
    text: messageText.value,
  })
</script>

<template>
  <div class="chat">
    <ul>
      <li v-for="(msg, idx) in store.state.chatMessages" :key="idx">
        <ProfilePicture :userID="msg.userID" />

        <p>
          <router-link :to="msg.username">{{ msg.username }}</router-link>
          <br />
          {{ msg.text }}
        </p>
      </li>
    </ul>

    <form v-if="store.state.currentRoomID" @submit.prevent="chat">
      <input type="text" placeholder="Say Hello!" v-model="messageText" />
      <input type="submit" value="SEND" />
    </form>

    <form v-else>
      <input type="text" placeholder="Say Hello!" disabled />
      <input type="submit" value="SEND" disabled />
    </form>
  </div>
</template>

<style scoped>
.chat {
  position: relative;
  display: flex;
  border: 2px solid #c1c1c1;
  border-style: dashed;
  flex-direction: column;
}

form {
  display: flex;
  height: 20px;
  margin: 0 5px 5px 5px;
}

input[type='text'] {
  width: 75%;
}

input[type='submit'] {
  cursor: pointer;
  width: 25%;
}

ul {
  height: 100%;
  list-style-type: none;
  margin: 0;
  padding: 5px;
  overflow-y: scroll;
  overscroll-behavior-y: contain;
  scroll-snap-type: y proximity;
}

li {
  display: flex;
  gap: 5px;
  border: 1px solid transparent;
  border-style: dotted;
  padding: 5px;
}

li:hover {
  border: 1px solid #c1c1c1;
  border-style: dotted;
}

.profile-picture {
  width: 35px;
  height: 35px;
  border: 2px solid #c1c1c1;
  border-style: inset;
}

p {
  width: calc(100% - 40px);
  margin: 0;
  overflow-wrap: break-word;
  justify-content: baseline;
}
</style>
