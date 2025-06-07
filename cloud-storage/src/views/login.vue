<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePreferencesStore } from "@/stores/preferences.ts";
import { useFileStore } from '@/stores/fileStore';
const preferencesStore = usePreferencesStore()
const fileStore = useFileStore()
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
        preferencesStore.setName(username.value)
        fileStore.fetchFiles("/");
        router.push({ name: 'files-root' })
    } catch (e: any) {
        console.error(e)
        error.value = 'Anmeldung fehlgeschlagen bitte überprüfen sie ihre Eingaben'
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
                <v-form>
                    <v-text-field v-model="username" :error="!!error" label="Username"
                        variant="outlined"></v-text-field>
                    <v-text-field v-model="password" :error="!!error"
                        :append-inner-icon="visible ? 'visibility_off' : 'visibility'"
                        :type="visible ? 'text' : 'password'" label="Password" variant="outlined"
                        @click:append-inner="visible = !visible"></v-text-field>
                    <v-slide-y-transition>
                        <v-alert v-if="error" type="error" class="mb-4" density="compact" border="start">
                            {{ error }}
                        </v-alert>
                    </v-slide-y-transition>
                    <v-btn variant="tonal" :loading="pending" :disabled="pending" @click="login" type="submit">
                        Login
                    </v-btn>
                </v-form>
            </v-container>
        </v-col>
    </v-row>
</template>
<style scoped lang="scss">
.login-container {
    background-color: rgb(var(--v-theme-background));
}
</style>
