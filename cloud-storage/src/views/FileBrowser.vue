<script setup lang="ts">
import Toolbar from '@/components/fileBrowser/Toolbar.vue';
import { useFileStore } from '@/stores/fileStore';

const fileStore = useFileStore();
</script>

<template>
  <v-card class="mx-auto" :loading="fileStore.isLoading">
    <Toolbar></Toolbar>
  
    <v-row no-gutters>
      <v-col>
        <RouterView></RouterView>
      </v-col>
    </v-row>
  </v-card>
</template>

<style lang="scss" scoped>
</style>