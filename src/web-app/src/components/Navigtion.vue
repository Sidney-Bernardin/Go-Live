<script setup></script>

<template>
  <div class="navigation">
    <ul class="links">
      <li><router-link to="home">Home</router-link></li>
    </ul>

    <input class="search" type="text" placeholder="Search for Users" />

    <button
      class="profile-picture"
      style="background: url('https://i.stack.imgur.com/ttV8S.jpg') center/100%"
    ></button>

    <div class="dropdown">
      <ul>
        <li><router-link to="sidney">Profile</router-link></li>
      </ul>

      <ul>
        <li><button>Logout</button></li>
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
  width: 100%;
  height: 60px;
  border-bottom: 6px solid #c1c1c1;
  border-bottom-style: double;
  justify-content: space-between;
  align-items: center;
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
