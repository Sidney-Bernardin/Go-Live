<script setup>
import UsersService from '../services/UsersService'

const props = defineProps({ userID: '' })
</script>

<template>
  <button
    class="profile-picture"
    :style="`background-image: url(${UsersService.profilePictureSrc(
      props.userID
    )})`"
  ></button>
</template>

<style scoped>
.profile-picture {
  width: 45px;
  height: 45px;
  border: 3px solid #c1c1c1;
  border-style: inset;
  background-size: 100% 100%;
  padding: 0;
}
</style>
