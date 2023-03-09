<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

import Stream from '../components/Stream.vue'
import Chat from '../components/Chat.vue'
import ProfilePicture from '../components/ProfilePicture.vue'

import UsersService from '../services/UsersService'
import RoomsService from '../services/RoomsService'

const router = useRouter()
const store = useStore()

const user = ref({
  id: '',
  username: '',
  email: '',
})

const room = ref({
  id: '',
  name: '',
})

onMounted(async () => {
  try {
    const res1 = await UsersService.getUser(
      router.currentRoute.value.params.username,
      'username',
      ['username', 'email']
    )
    user.value.id = res1.data.id
    user.value.username = res1.data.username
    user.value.email = res1.data.email

    const res2 = await RoomsService.getRoom(user.value.id)
    room.value.id = res2.data.id
    room.value.name = res2.data.name
  } catch (err) {
    if (err.response?.data.type == 'room_doesnt_exist') return
    if (err.response?.data.type == 'user_doesnt_exist')
      err = "User doesn't exixt."

    store.dispatch('handleError', err)
    return
  }

  store.dispatch('joinRoom', user.value.id)
})

onUnmounted(() => store.dispatch('leaveRoom'))
</script>

<template>
  <div class="user-page">
    <div class="room" v-if="room.id">
      <h1>{{ room.name }}</h1>

      <Stream v-if="store.state.currentRoomID" />

      <div class="video-placeholder" v-else>
        <h2>Couldn't join room</h2>
      </div>

      <Chat />
    </div>

    <div class="info" v-if="user.id">
      <ProfilePicture :userID="user.id" />
      {{ user.username }} - {{ user.email }}
    </div>
  </div>
</template>

<style scoped>
.room {
  position: relative;
  display: grid;
  gap: 10px;
  grid-template-columns: 1fr 1fr 1fr 250px;
  grid-template-rows: auto calc(100vh - 250px);
  grid-template-areas:
    'header header header header'
    'video video video chat';
  border: 1px solid #c1c1c1;
  border-style: dotted;
  padding: 10px;
  margin-bottom: 15px;
}

.room h1 {
  grid-area: header;
  height: 40px;
  line-height: 40px;
  margin: 0;
}

.room .stream {
  grid-area: video;
  width: 100%;
  max-height: 100%;
}

.room .video-placeholder {
  display: flex;
  grid-area: video;
  border: 2px solid #c1c1c1;
  border-style: dashed;
  justify-content: center;
  align-items: center;
}

.room .chat {
  grid-area: chat;
}

.info {
  display: flex;
  gap: 10px;
  border: 1px solid #c1c1c1;
  border-style: dotted;
  padding: 10px;
  align-items: center;
}
</style>
