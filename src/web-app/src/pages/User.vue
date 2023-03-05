<script setup>
import { ref, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

import UsersService from '../services/UsersService'

const store = useStore()
const router = useRouter()

const user = ref({
  id: '',
  username: router.currentRoute.value.params.username,
  email: '',
})

const profilePictureSrc = (userID) => `${import.meta.env.VITE_MICROSERVICES_URL}/users/all/${userID}/picture`

UsersService.getUser(user.value.username, 'username', ['email'])
  .then((res) => {
    user.value.id = res.data.id
    user.value.email = res.data.email
  })
  .catch((err) => store.dispatch('handleError', err))

onMounted(() => {})
</script>

<template>
  <div class="user-page">

    <div class="room">
      <h1>Room Name</h1>

      <video controls autoplay muted />

      <div class="chat">
        <ul>
          <li v-for="(msg, idx) in store.state.chatMessages" :key="idx">
            <button
              :style="`background: url(${profilePictureSrc(msg.id)}) center/100%`"
            ></button>

            <p>{{ msg.text }}</p>
          </li>
        </ul>

        <form>
          <input type="text" placeholder="Hello, World!">
          <input type="submit" value="SEND">
        </form>
      </div>
    </div>

    <div class="info"></div>
  </div>
</template>

<style scoped>
.room {
  position: relative;
  display: grid;
  gap: 10px;
  grid-template-columns: 1fr 1fr 1fr 250px;
  grid-template-rows: auto calc(100vh - 200px);
  grid-template-areas: 
    "header header header header"
    "video video video chat";
  border: 1px solid #c1c1c1;
  border-style: dotted;
  padding: 10px;
}

.room h1 {
  grid-area: header;
  margin: 0;
}

.room video {
  grid-area: video;
  width: 100%;
  max-height: 100%;
}

.room .chat {
  position: relative;
  display: flex;
  grid-area: chat;
  border: 2px solid #c1c1c1;
  border-style: dashed;
  flex-direction: column;
}

.room .chat form {
  display: flex;
  height: 20px;
  margin: 0 5px 5px 5px;
}

.room .chat input[type=text] {
  width: 75%;
}

.room .chat input[type=submit] {
  width: 25%;
}

.room .chat ul {
  height: 100%;
  list-style-type: none;
  margin: 0;
  padding: 5px;
  overflow-y: scroll;
  overscroll-behavior-y: contain;
  scroll-snap-type: y proximity;
}

.room .chat li {
  display: flex;
  gap: 5px;
  border: 1px solid transparent;
  border-style: dotted;
  padding: 5px;
}

.room .chat li:hover {
  border: 1px solid #c1c1c1;
  border-style: dotted;
}

.room .chat li button {
  width: 35px;
  height: 35px;
  border: 2px solid #c1c1c1;
  border-style: inset;
  padding: 0;
}

.room .chat p {
  width: calc(100% - 40px);
  margin: 0;
  overflow-wrap: break-word;
  justify-content: baseline;
}
</style>
