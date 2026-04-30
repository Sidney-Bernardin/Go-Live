<script setup lang="ts">
import { ref } from 'vue';
import { router } from '../router';

const signUpOrIn = ref(false)
const profilePictureURL = ref("");
const signupError = ref("");
const signinError = ref("");

async function signup(e: Event) {
    const res = await fetch("http://localhost:8000/users/login/signup", {
        method: "POST",
        body: new FormData(e.target as HTMLFormElement),
        credentials: "include",
    })

    if (res.ok) router.push({ name: "Home" })
    else {
        const err = await res.json()
        signupError.value = err.message
    }
}

async function signin(e: Event) {
    const res = await fetch("http://localhost:8000/users/login/signin", {
        method: "POST",
        body: new FormData(e.target as HTMLFormElement),
        credentials: "include",
    })

    if (res.ok) router.push({ name: "Home" })
    else {
        const err = await res.json()
        signinError.value = err.message
    }
}

function changeProfilePicture(e: Event) {
    const file: File = (e.target as HTMLInputElement).files![0]!;
    profilePictureURL.value = URL.createObjectURL(file);
}
</script>

<template>
    <div class="login page">
        <div class="tabs">
            <button :disabled="signUpOrIn" @click="signUpOrIn = true">sign-in</button>
            <button :disabled="!signUpOrIn" @click="signUpOrIn = false">sign-up</button>
        </div>

        <form v-if="!signUpOrIn" @submit.prevent="signup">
            <label class="profile-picture-label">
                <input type="file" name="profile_picture" @change="changeProfilePicture" />
                <img :src="profilePictureURL" />
                Upload Profile Picture
            </label>

            <label for="username">Username</label>
            <input id="username" type="text" name="username" placeholder="Username" />

            <label for="email">Email</label>
            <input id="email" type="email" name="email" placeholder="Email" />

            <label for="password">Password</label>
            <input id="password" type="password" name="password" placeholder="Password" />

            <input type="submit" value="sign-up" />

            <p class="error" v-if="signupError">{{ signupError }}</p>
        </form>

        <form v-if="signUpOrIn" @submit.prevent="signin">
            <label for="username">Username</label>
            <input id="username" type="text" name="username" placeholder="Username" />

            <label for="password">Password</label>
            <input id="password" type="password" name="password" placeholder="Password" />

            <input type="submit" value="sign-in" />

            <p class="error" v-if="signinError">{{ signinError }}</p>
        </form>
    </div>
</template>

<style scoped>
.login.page {
    display: flex;
    height: 100%;
    justify-content: center;
    align-items: center;
}

.tabs {
    display: flex;
    position: absolute;
    top: 50%;
    left: 2.5rem;
    width: 0;
    gap: 15px;
    rotate: -90deg;
    transform: translate(-50%, 50%);
    justify-content: center;
}

.tabs button {
    display: flex;
    cursor: pointer;
    border: none;
    color: var(--white);
    background: transparent;
    font-size: 1.5rem;
    font-weight: bolder;
    text-transform: uppercase;
    text-decoration: underline;
    text-wrap: nowrap;
}

.tabs button:disabled {
    cursor: auto;
    color: var(--dark-green);
    text-decoration: none;
}

form {
    display: flex;
    border: 2px solid var(--black);
    padding: 2rem;
    flex-direction: column;
}

label {
    font-size: 1.5rem;
    font-weight: bolder;
}

label.profile-picture-label {
    display: flex;
    margin-bottom: 2rem;
    border: 2px solid var(--black);
    padding: 1rem;
    gap: 1rem;
    background: transparent;
    justify-content: center;
    align-items: center;
}

label.profile-picture-label img {
    width: 3rem;
    height: 3rem;
}

input {
    margin-bottom: 2rem;
    border: 2px solid var(--black);
    padding: 1rem;
    color: var(--white);
    background: transparent;
    font-size: 1.5rem;
    font-weight: bolder;
}

input[type="file"] {
    display: none;
}

input::placeholder {
    color: var(--dark-green);
}

input[type="submit"] {
    cursor: pointer;
    margin-bottom: 0;
    border: 2px solid var(--white);
    text-transform: uppercase;
}

.error {
    margin-top: 2rem;
    width: 100%;
    color: var(--yellow);
    font-size: 1.5rem;
    font-weight: bolder;
}
</style>
