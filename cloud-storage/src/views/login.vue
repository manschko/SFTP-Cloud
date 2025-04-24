<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const username = ref('')
const password = ref('')
const visible = ref(false)
const pending = ref(false)
const error = ref('')
const router = useRouter()

const login = async () => {
    pending.value = true
    error.value = ''
    try {
        const res = await fetch('http://localhost:8000/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: username.value, password: password.value })
        })
        if (!res.ok) throw new Error('Login failed')
        const data = await res.json()
        localStorage.setItem('jwt', data.token)
        router.push({ name: 'Home' }) // Change 'Home' to your main route name
    } catch (e: any) {
        error.value = e.message || 'Login failed'
    } finally {
        pending.value = false
    }
}
</script>

<template>
    <v-row no-gutters align="start" class="gradient-background">
        <v-col cols="12" md="4" class="d-flex flex-column align-center justify-center h-screen login-container">
            <v-container>
                <h1>
                    Login into Cloud Storage
                </h1>
                <v-text-field v-model="username" label="Username" variant="outlined"></v-text-field>
                <v-text-field v-model="password" :append-inner-icon="visible ? 'visibility_off' : 'visibility'"
                    :type="visible ? 'text' : 'password'" label="Password" variant="outlined"
                    @click:append-inner="visible = !visible"></v-text-field>
                <v-btn variant="tonal" :loading="pending" :disabled="pending" @click="login">
                    Login
                </v-btn>
            </v-container>
        </v-col>
    </v-row>
</template>
<style scoped lang="scss">
.login-container {
    background-color: rgb(var(--v-theme-background));
}
</style>