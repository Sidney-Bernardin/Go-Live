<script setup>
import UsersService from '../services/UsersService'

const props = defineProps({ userID: '' })
</script>

<template>
  <button
    class="profile-picture"
    :style="`background: url(${UsersService.profilePictureSrc(
      props.userID
    )}) center/100%`"
  ></button>
</template>

<style scoped>
.profile-picture {
  width: 45px;
  height: 45px;
  border: 3px solid #c1c1c1;
  border-style: inset;
  padding: 0;
}
</style>
