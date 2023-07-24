<script setup>
import { ref, computed, inject } from 'vue'
import { useStore } from 'vuex'

import RTMPService from '../services/RTMPService'
import { getSessionID } from '../utils'

const store = useStore()

const showURICreator = inject('show_URI_creator')

const uriInfo = ref({
  name: store.state.self.username + "'s Room",
})

const uri = computed(() => RTMPService.liveURI(store.state.self.id))
const key = computed(() => RTMPService.liveKey(uriInfo.value.name))
</script>

<template>
  <div class="uri-creator">
    <div class="overlay" @click="showURICreator = false"></div>

    <div class="main">

      <div class="output">
        <div class="uri">
          <h3>URI:</h3>
          <p>{{ uri }}</p>
        </div>

        <div class="key">
          <h3>KEY:</h3>
          <p>{{ key }}</p>
        </div>
      </div>

      <p>
        Fill out the settings for your room, then begin streaming with the above
        URI and key.
      </p>

      <form @submit.prevent>
        <h2>Room Settings</h2>

        <label for="name">Name</label>
        <input
          type="text"
          name="name"
          placeholder="Name"
          v-model="uriInfo.name"
        />
      </form>
    </div>
  </div>
</template>

<style scoped>
.uri-creator {
  z-index: 1;
  position: fixed;
  display: flex;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  justify-content: center;
  align-items: center;
}

.overlay {
  position: absolute;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
}

.main {
  z-index: 1;
  position: relative;
  display: flex;
  width: 400px;
  gap: 15px;
  border: 6px solid #c1c1c1;
  border-style: outset;
  background: #fff;
  padding: 15px;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.output {
  width: 350px;
  border: 2px solid #c1c1c1;
  border-style: dotted;
  padding: 10px;
}

.output h3 {
  margin: 0;
}

.output p {
  margin: 0;
  overflow: scroll;
}

.output div {
  display: flex;
  gap: 10px;
}

p {
  width: 300px;
  text-align: center;
}

form {
  display: flex;
  border: 2px solid #c1c1c1;
  border-style: dashed;
  padding: 10px;
  flex-direction: column;
}

h2 {
  margin: 0;
}

label {
  margin-top: 5px;
}
</style>
