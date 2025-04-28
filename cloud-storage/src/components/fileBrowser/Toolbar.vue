<script setup lang="ts">
import {computed, onMounted, onUnmounted, ref, watch} from 'vue';
import {usePreferencesStore} from "@/stores/preferences.ts";
import {useRoute, useRouter} from "vue-router";

interface Storage {
  code: string;
  name: string;
  icon: string;
}

interface PathSegment {
  name: string;
  path: string;
}

interface Endpoints {
  mkdir: {
    url: string;
    method?: string;
  };

  [key: string]: any;
}

interface Props {
  path: string;
  endpoints: Endpoints;
  axios: any; // Consider using a more specific type for axios
}

const props = defineProps<Props>();
const emit = defineEmits<{
  'storage-changed': [code: string];
  'path-changed': [path: string];
  'add-files': [files: FileList];
  'loading': [isLoading: boolean];
  'folder-created': [folderName: string];
}>();

const inputUpload = ref<HTMLInputElement | null>(null);
const newFolderPopper = ref(false);
const newFolderName = ref('');

const route = useRoute();

const pathSegments = ref<PathSegment[]>([]);
const hiddenSegments = ref<PathSegment[]>([]);

function updateSegments() {
  const segments = route.fullPath.split('/').filter(item => item);
  segments.shift(); // Remove the first segment (base path)
  const basePath = "/files";

  const width = document.documentElement.clientWidth;
  const slots = Math.floor((width - 150) / 175) || 1;

  const allSegments = segments.map((item, index) => {
    const partialPath = segments.slice(0, index + 1).join('/');
    return {
      name: item,
      path: `${basePath}/${partialPath}`,
    };
  });

  if (allSegments.length <= slots) {
    hiddenSegments.value = [];
    pathSegments.value = allSegments;
  } else {
    hiddenSegments.value
    hiddenSegments.value = [
      {
        name: preferencesStore.getName(),
        path: `${basePath}`,
      },
      ...allSegments.slice(0, -slots)];
    pathSegments.value = allSegments.slice(-slots);
  }
}
onMounted(() => {
  updateSegments();
  window.addEventListener('resize', updateSegments);
});

onUnmounted(() => {
  window.removeEventListener('resize', updateSegments);
});
watch(() => route.fullPath, updateSegments);
/*const pathSegments = computed((): PathSegment[] => {
  let path = '/';
  const isFolder = props.path[props.path.length - 1] === '/';
  const segments = props.path.split('/').filter(item => item);

  return segments.map((item, index) => {
    path += item + (index < segments.length - 1 || isFolder ? '/' : '');
    return {
      name: item,
      path
    };
  });
});

const storageObject = computed((): Storage => {
  return props.storages.find(item => item.code === props.storage) as Storage;
});

function changeStorage(code: string): void {
  if (props.storage !== code) {
    emit('storage-changed', code);
    emit('path-changed', '');
  }
}*/
const router = useRouter();

function changePath(path: string): void {
  router.push(path)
}

function goUp(): void {
  const segments = pathSegments.values;
  const path = segments.length === 1 ? '/' : segments[segments.length - 2].path;
  changePath(path);
}

function addFiles(event: Event): void {
  const target = event.target as HTMLInputElement;
  if (target.files) {
    emit('add-files', target.files);
    if (inputUpload.value) {
      inputUpload.value.value = '';
    }
  }
}

async function mkdir(): Promise<void> {
  emit('loading', true);

  let url = props.endopints.mkdir.url
    .replace(new RegExp("{storage}", "g"), props.storage)
    .replace(new RegExp("{path}", "g"), props.path + newFolderName.value);

  const config = {
    url,
    method: props.endpoints.mkdir.method || 'post'
  };

  await props.axios.request(config);
  emit('folder-created', newFolderName.value);
  newFolderPopper.value = false;
  newFolderName.value = '';
  emit('loading', false);
}

const preferencesStore = usePreferencesStore();
const viewMode = computed({
  get: () => preferencesStore.getViewMode(),
  set: (value) => preferencesStore.setViewMode(value as 'grid' | 'list'),
});
const showHiddenSegments =  ref<boolean>(false);
</script>

<template>
  <v-toolbar flat density="compact" color="blue-grey-lighten-5">
    <v-toolbar-items>
      <v-btn
        variant="text"
        @click="changePath('/files')"
        v-if="!hiddenSegments.length"
      >
        <v-icon class="mr-2">{{ 'home' }}</v-icon>
        {{ preferencesStore.getName() }}
      </v-btn>
      <v-btn
        variant="text"
        @click="showHiddenSegments = !showHiddenSegments"
        v-else
        >
        <v-icon>more_horiz</v-icon>
        <v-menu v-model="showHiddenSegments" location="bottom">
          <v-list>
            <v-list-item
              v-for="(segment, index) in hiddenSegments"
              :key="index"
              @click="changePath(segment.path)"
            >
              <v-icon class="mr-2">{{ index === 0 ? 'home': 'folder' }}</v-icon>
              {{ segment.name }}
            </v-list-item>
          </v-list>
        </v-menu>
      </v-btn>
      <template v-for="(segment, index) in pathSegments" :key="index">
        <v-icon class="h-100">chevron_right</v-icon>
        <v-btn
          class="segment-button"
          variant="text"
          @click="changePath(segment.path)"
        >
          {{ segment.name }}
        </v-btn>
      </template>
    </v-toolbar-items>
    <v-spacer></v-spacer>

    <template v-slot:append>
      <v-btn-toggle rounded="xl" v-model="viewMode" mandatory color="primary" variant="outlined"
                    density="compact">
        <v-btn value="list">

          <v-icon v-show="viewMode == 'list'">check</v-icon>
          <v-icon>menu</v-icon>
        </v-btn>

        <v-btn value="grid">
          <v-icon v-show="viewMode == 'grid'">check</v-icon>
          <v-icon>grid_view</v-icon>
        </v-btn>

      </v-btn-toggle>
    </template>
  </v-toolbar>
</template>

<style scoped>
.segment-button {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap
}

.storage-select-list :deep(.v-list-item--disabled) {
  background-color: rgba(0, 0, 0, 0.25);
  color: #fff;
}

.storage-select-list :deep(.v-list-item--disabled .v-icon) {
  color: #fff;
}
</style>

<style lang="scss">
.segment-button {
  .v-btn__content{
    display: block;
    max-width: 150px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap
  }
}
</style>
