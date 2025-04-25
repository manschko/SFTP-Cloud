<script setup lang="ts">
import { ref } from 'vue';
const route = useRoute();
const currentPath2 = route.params.pathMatch
  ? Array.isArray(route.params.pathMatch)
    ? '/' + route.params.pathMatch.join('/')
    : '/' + route.params.pathMatch
  : '/';
/*
import { ref, computed, onMounted } from 'vue';
import { useDisplay } from 'vuetify';
import axios, { AxiosInstance } from 'axios';

import List from '../components/fileBrowser/List.vue';
import Tree from '../components/fileBrowser/Tree.vue';
import Upload from '../components/fileBrowser/Upload.vue';
import Toolbar from '../components/fileBrowser/Toolbar.vue';

// Types and interfaces
interface Storage {
  name: string;
  code: string;
  icon: string;
}

interface Endpoint {
  url: string;
  method: string;
}

interface Endpoints {
  list: Endpoint;
  upload: Endpoint;
  mkdir: Endpoint;
  delete: Endpoint;
}

interface Icons {
  [key: string]: string;
}

interface Props {
  storages?: string;
  storage?: string;
  tree?: boolean;
  icons?: Icons;
  endpoints?: Endpoints;
  axios?: () => AxiosInstance;
  axiosConfig?: object;
  maxUploadFilesCount?: number;
  maxUploadFileSize?: number;
}

// Define available storages
const availableStorages: Storage[] = [
  {
    name: "Local",
    code: "local",
    icon: "mdi-folder-multiple-outline"
  }
];

// Default endpoints
const defaultEndpoints: Endpoints = {
  list: { url: "/storage/{storage}/list?path={path}", method: "get" },
  upload: { url: "/storage/{storage}/upload?path={path}", method: "post" },
  mkdir: { url: "/storage/{storage}/mkdir?path={path}", method: "post" },
  delete: { url: "/storage/{storage}/delete?path={path}", method: "post" }
};

// Default file icons


// Define props with default values
const props =null;
/*const props = withDefaults(defineProps<Props>(), {
  storages: () => availableStorages.map(item => item.code).join(","),
  storage: "local",
  tree: true,
  icons: () => fileIcons,
  endpoints: () => defaultEndpoints,
  axiosConfig: () => ({}),
  maxUploadFilesCount: 0,
  maxUploadFileSize: 0
});

// Emits
const emit = defineEmits<{
  (e: 'change', path: string): void
}>();

// Reactive state
const loading = ref(0);
const path = ref("");
const activeStorage = ref<string | null>(null);
const uploadingFiles = ref<File[] | false>(false);
const refreshPending = ref(false);
const axiosInstance = ref<AxiosInstance | null>(null);

// Computed properties
const storagesArray = computed(() => {
  const storageCodes = props.storages.split(",");
  const result: Storage[] = [];
  storageCodes.forEach(code => {
    const storage = availableStorages.find(item => item.code === code);
    if (storage) result.push(storage);
  });
  return result;
});

const display = useDisplay();

// Methods
const loadingChanged = (isLoading: boolean) => {
  if (isLoading) {
    loading.value++;
  } else if (loading.value > 0) {
    loading.value--;
  }
};

const storageChanged = (storage: string) => {
  activeStorage.value = storage;
};

const addUploadingFiles = (files: FileList) => {
  let fileArray = Array.from(files);

  if (props.maxUploadFileSize) {
    fileArray = fileArray.filter(
      file => file.size <= props.maxUploadFileSize
    );
  }

  if (uploadingFiles.value === false) {
    uploadingFiles.value = [];
  }
  
  if (Array.isArray(uploadingFiles.value)) {
    if (props.maxUploadFilesCount && uploadingFiles.value.length + fileArray.length > props.maxUploadFilesCount) {
      fileArray = fileArray.slice(0, props.maxUploadFilesCount - uploadingFiles.value.length);
    }

    uploadingFiles.value.push(...fileArray);
  }
};

const removeUploadingFile = (index: number) => {
  if (uploadingFiles.value !== false) {
    uploadingFiles.value.splice(index, 1);
  }
};

const uploaded = () => {
  uploadingFiles.value = false;
  refreshPending.value = true;
};

const pathChanged = (newPath: string) => {
  path.value = newPath;
  emit('change', newPath);
};

// Lifecycle hooks
activeStorage.value = props.storage;
axiosInstance.value = props.axios ? props.axios() : axios.create(props.axiosConfig);

onMounted(() => {
  if (!path.value && !(props.tree && display.smAndUp.value)) {
    pathChanged("/");
  }
});*/
import Toolbar from '@/components/fileBrowser/Toolbar.vue';
const loading = ref(1);
interface Endpoint {
  url: string;
  method: string;
}
import FileWindow from '@/views/FileWindow.vue';
import Tree from '../components/fileBrowser/Tree.vue';
import Upload from '../components/fileBrowser/Upload.vue';
import { useRoute } from 'vue-router';
const defaultEndpoints: Endpoints = {
  list: { url: "/storage/{storage}/list?path={path}", method: "get" },
  upload: { url: "/storage/{storage}/upload?path={path}", method: "post" },
  mkdir: { url: "/storage/{storage}/mkdir?path={path}", method: "post" },
  delete: { url: "/storage/{storage}/delete?path={path}", method: "post" }
};
const currentPath = ref("/");

</script>

<template>
  <v-card class="mx-auto" :loading="loading > 0">
    <Toolbar
        :path="currentPath"
        :storages="[]"
        :storage="'local'"
        :endpoints="defaultEndpoints"
        :axios="null"
        v-on:storage-changed="null"
        v-on:path-changed="null"
        v-on:add-files="null"
        v-on:folder-created="true"
/>
    <v-row no-gutters>
      <v-col v-if="tree && display.smAndUp.value" sm="auto">
        <Tree
          :path="currentPath"
          :refreshPending="refreshPending"
          @path-changed="pathChanged"
          @loading="loadingChanged"
          @refreshed="refreshPending = false"
        />
      </v-col>
      <v-divider v-if="tree" vertical />
      <v-col>
        <RouterView></RouterView>
      
      </v-col>
    </v-row>
    <!--Upload
      v-if="uploadingFiles !== false"
      :path="path"
      :storage="activeStorage"
      :files="uploadingFiles"
      :icons="icons"
      :axios="axiosInstance"
      :endpoint="endpoints.upload"
      :maxUploadFilesCount="maxUploadFilesCount"
      :maxUploadFileSize="maxUploadFileSize"
      @add-files="addUploadingFiles"
      @remove-file="removeUploadingFile"
      @clear-files="uploadingFiles = []"
      @cancel="uploadingFiles = false"
      @uploaded="uploaded"
    /-->
  </v-card>
</template>

<style lang="scss" scoped>
</style>