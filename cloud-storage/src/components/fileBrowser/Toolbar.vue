<script setup lang="ts">
import { ref, computed } from 'vue';

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

const pathSegments: PathSegment[] = [{name: "folder1", path: "/folder1"}]
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

function changePath(path: string): void {
  emit('path-changed', path);
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
</script>

<template>
  <v-toolbar flat density="compact" color="blue-grey-lighten-5">
    <v-toolbar-items>
      <v-btn
        variant="text"
        :active="path === '/'"
        @click="changePath('/')"
      >
        <v-icon class="mr-2">{{ 'home'/*storageObject.icon*/ }}</v-icon>
        {{  'test' /*storageObject.name*/ }}
      </v-btn>
      <template v-for="(segment, index) in pathSegments" :key="index">
        <v-icon>chevron_right</v-icon>
        <v-btn
          variant="text"
          :active="index === pathSegments.length - 1"
          @click="changePath(segment.path)"
        >
          {{ segment.name }}
        </v-btn>
      </template>
    </v-toolbar-items>
    <v-spacer></v-spacer>

    <template v-slot:append>
      <!--v-tooltip location="bottom" v-if="pathSegments.length > 0">
        <template v-slot:activator="{ props: tooltipProps }">
          <v-btn icon v-bind="tooltipProps" @click="goUp">
            <v-icon>arrow_upward</v-icon>
          </v-btn>
        </template>
        <span v-if="pathSegments.length === 1">Up to "root"</span>
        <span v-else>Up to "{{ pathSegments[pathSegments.length - 2].name }}"</span>
      </v-tooltip>
      <v-menu
        v-model="newFolderPopper"
        :close-on-content-click="false"
        :width="200"
        location="bottom"
      >
        <template #activator="{ props: menuProps }">
          <v-btn v-if="path" icon v-bind="menuProps" title="Create Folder">
            <v-icon>create_new_folder</v-icon>
          </v-btn>
        </template>
        <v-card>
          <v-card-text>
            <v-text-field
              label="Name"
              v-model="newFolderName"
              hide-details
            ></v-text-field>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn @click="newFolderPopper = false" variant="text">Cancel</v-btn>
            <v-btn
              color="success"
              :disabled="!newFolderName"
              variant="text"
              @click="mkdir"
            >
              Create Folder
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-menu>
      <v-btn v-if="path" icon @click="$refs.inputUpload.click()" title="Upload Files">
        <v-icon>upload_file</v-icon>
        <input v-show="false" ref="inputUpload" type="file" multiple @change="addFiles" />
      </v-btn-->
      <v-btn-toggle v-model="viewMode" mandatory color="primary" variant="outlined">
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
.storage-select-list :deep(.v-list-item--disabled) {
  background-color: rgba(0, 0, 0, 0.25);
  color: #fff;
}

.storage-select-list :deep(.v-list-item--disabled .v-icon) {
  color: #fff;
}
</style>