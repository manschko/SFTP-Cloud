<script setup lang="ts">
import { useFileStore } from '@/stores/fileStore';
import { computed, ref } from 'vue';
import { useRoute}  from 'vue-router';
import {usePreferencesStore} from "@/stores/preferences.ts";


const fileStore = useFileStore();
const route = useRoute();
const currentPath = route.params.pathMatch
  ? Array.isArray(route.params.pathMatch)
    ? '/' + route.params.pathMatch.join('/')
    : '/' + route.params.pathMatch
  : '/';
fileStore.fetchFiles(currentPath);

const preferencesStore = usePreferencesStore()

// Sorting state
const sortItems = [
                    {
                        text: 'Name',
                        value: 'name'
                    },
                    {
                        text: 'Größe',
                        value: 'size'
                    },
                    {
                        text: 'Änderungsdatum',
                        value: 'mod_time'
                    }]

const fileIcons: Icons = {
    zip: "folder_zip",
  rar: "folder_zip",
  '7z': "folder_zip",
  htm: "html",
  html: "html",
  js: "javascript",
  json: "data_object",
  md: "article",
  pdf: "picture_as_pdf",
  png: "image",
  jpg: "image",
  jpeg: "image",
    webp: "image",
  mp4: "movie",
  mkv: "movie",
  avi: "movie",
  wmv: "movie",
  mov: "movie",
  txt: "description",
  xls: "table_chart",
  other: "insert_drive_file"
};
function getFileIcon(fileName: string) {
  const ext = fileName.split('.').pop()
  return fileIcons[ext] || fileIcons['other']
}

function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('de-DE', {
    day: '2-digit',
    month: 'long',
    year: 'numeric'
  })
}

// Format bytes
function formatBytes(bytes: number) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const sortedAll = computed(() => {
  return [...fileStore.folders, ...fileStore.files].sort((a, b) => {
    let result = 0
    const sortBy = preferencesStore.getSortBy();
    const sortOrder = preferencesStore.getSortOrder();

    if (a[sortBy] < b[sortBy]) result = -1;
    if (a[sortBy] > b[sortBy]) result = 1;

    return sortOrder === 'asc' ? result : -result
  })
})

const sortedFolders = computed(() => sortedAll.value.filter(item => item.is_dir))
const sortedFiles = computed(() => sortedAll.value.filter(item => !item.is_dir))

function changeSort(key: 'name' | 'mod_time' | 'size') {
  if (preferencesStore.getSortBy() === key) {
    preferencesStore.toggleSortOrder();
  } else {
    preferencesStore.setSortOrder('asc');
    preferencesStore.setSortBy(key);
  }
}

const isDragging = ref(false)

function onDragOver(e: DragEvent) {
  e.preventDefault()
  isDragging.value = true
}
function onDragLeave(e: DragEvent) {
  e.preventDefault()
  // Check if the mouse is outside the container
  const current = e.currentTarget as Node
  const related = e.relatedTarget as Node | null
  if (!current?.contains(related)) {
    isDragging.value = false
  }
}
function onDrop(e: DragEvent) {
  e.preventDefault()
  isDragging.value = false
  // handle files here, e.dataTransfer.files
}
const sortKey = computed({
  get: () => preferencesStore.getSortBy(),
  set: (value) => preferencesStore.setSortBy(value as 'name' | 'mod_time' | 'size'),
});
</script>
<template>
    <v-card flat tile min-height="380" class="d-flex flex-column ">
        <!--v-btn-toggle v-model="viewMode" mandatory color="primary" variant="outlined">
          <v-btn value="list">

            <v-icon v-show="viewMode == 'list'">check</v-icon>
            <v-icon>menu</v-icon>
          </v-btn>

          <v-btn value="grid">
            <v-icon v-show="viewMode == 'grid'">check</v-icon>
            <v-icon>grid_view</v-icon>
          </v-btn>

        </v-btn-toggle-->
        <v-card-text v-if="fileStore.files || fileStore.folders" class="grow">
            <div
            class="drop-container"
            @dragover="onDragOver"
            @dragleave="onDragLeave"
            @drop="onDrop">
            <div v-if="isDragging" class="drop-overlay">
                <div class="drop-message">
                    <v-icon large color="primary" class="upload-drag-icon">cloud_upload</v-icon>
                    <div>Dateien zum Hochladen ablegen</div>
                </div>
            </div>
            <v-table v-if="preferencesStore.getViewMode() == 'list'">
                <thead>
                <tr>
                    <th class="text-left">
                    <v-btn variant="text" @click="() => changeSort('name')" class="sort-btn">
                    Name
                    <v-icon :style="{ opacity: preferencesStore.getSortBy() === 'name' ? 1 : 0 }">
                        {{ preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down' }}
                    </v-icon>
                    </v-btn>
                    </th>
                    <th class="text-left">
                    <v-btn variant="text" @click="() => changeSort('mod_time')" class="sort-btn">
                    modified
                    <v-icon :style="{ opacity: preferencesStore.getSortBy() === 'mod_time' ? 1 : 0 }">
                        {{ preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down' }}
                    </v-icon>
                    </v-btn>
                    </th>
                    <th class="text-left">
                    <v-btn variant="text" @click="() => changeSort('size')" class="sort-btn">
                    size
                    <v-icon :style="{ opacity: preferencesStore.getSortBy() === 'size' ? 1 : 0 }">
                        {{ preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down' }}
                    </v-icon>
                    </v-btn>
                    </th>
                </tr>
                </thead>
                <tbody>
                <tr
                    v-for="item in sortedAll"
                    :key="item.name"
                >
                    <td><v-icon class="icon-space">{{item.is_dir ? 'folder' :getFileIcon(item.name)}}</v-icon>  {{ item.name }}</td>
                    <td>{{ formatDate(item.mod_time) }}</td>
                    <td>{{ item.is_dir ? '' :formatBytes(item.size) }}</td>
                </tr>
                </tbody>
            </v-table>
            <v-card-text v-else>
                <div class="d-flex justify-end position-absolute top-0 right-0 p-2 sort-container">
                    <v-btn
                    variant="text"
                    :icon="preferencesStore.getSortOrder() ==='asc'? 'arrow_upward': 'arrow_downward'"
                    @click="() => preferencesStore.toggleSortOrder()"
                    style="height:2.5rem;"
                    >

                    </v-btn>
                    <v-select
                    v-model="sortKey"
                    :items="sortItems"
                    item-title="text"
                    item-value="value"
                    class="sort-select"
                    density="compact"
                    ></v-select>
                </div>

                <p class="text-subtitle-1 mb-2 mt-2">Ordner</p>
                <v-row
                align="start">
                    <v-col
                    cols="12"
                    sm="4"
                    md="3"
                    lg="2"
                    xl="1"
                        v-for="folder in sortedFolders"
                        :key="folder.name"
                        class=""
                        @click="(console.log(currentPath ,folder.path))"
                    >
                        <v-card
                        ripple
                        color="grey-lighten-3"
                        class="mb-2">
                            <v-card-title class="d-flex align-center">
                                <v-icon class="mr-2">folder</v-icon><span class="folder-title text-truncate">{{ folder.name }}</span>
                            </v-card-title>
                        </v-card>
                    </v-col>
                </v-row>
                <p class="text-subtitle-1 mb-2 mt-2">Dateien</p>
                <v-row
                align="start">
                    <v-col
                    cols="12"
                    sm="4"
                    md="3"
                    lg="2"
                    xl="1"
                        v-for="file in sortedFiles"
                        :key="file.name"
                        class=""
                        @click="(console.log(currentPath ,file.path))"
                    >
                        <v-card
                        ripple
                        color="grey-lighten-4"
                        class="mb-2">
                        <v-card-title class="d-flex align-center">
                            <v-icon large class="mr-4 file-icon">{{ getFileIcon(file.name) }}</v-icon>
                            <div class="d-flex flex-column">
                                <span class="file-name text-h6">{{ file.name }}</span>
                                <span class="file-size text-body-2">{{ formatBytes(file.size) }}</span>
                                <span class="file-date text-caption">{{ formatDate(file.mod_time) }}</span>
                            </div>
                            </v-card-title>
                        </v-card>
                    </v-col>
                </v-row>
            </v-card-text>

        </div>
        </v-card-text>
        <v-card-text
            v-else
            class="grow d-flex justify-center align-center grey--text py-5"
        >The folder is empty</v-card-text>
    </v-card>
</template>

<!--script lang="ts">
import { formatBytes } from "@/components/fileBrowser/util";
//import Confirm from "@/components/Confirm.vue";

export default {
    props: {
        icons: Object,
        storage: String,
        path: String,
        endpoints: Object,
        axios: Function,
        refreshPending: Boolean
    },
    components: {
        //Confirm
    },
    data() {
        return {
            items: [],
            filter: ""
        };
    },
    computed: {
        dirs() {
            return this.items.filter(
                item =>
                    item.type === "dir" && item.basename.includes(this.filter)
            );
        },
        files() {
            return this.items.filter(
                item =>
                    item.type === "file" && item.basename.includes(this.filter)
            );
        },
        isDir() {
            return this.path[this.path.length - 1] === "/";
        },
        isFile() {
            return !this.isDir;
        }
    },
    methods: {
        formatBytes,
        changePath(path) {
            this.$emit("path-changed", path);
        },
        async load() {
            this.$emit("loading", true);
            if (this.isDir) {
                let url = this.endpoints.list.url
                    .replace(new RegExp("{storage}", "g"), this.storage)
                    .replace(new RegExp("{path}", "g"), this.path);

                let config = {
                    url,
                    method: this.endpoints.list.method || "get"
                };

                let response = await this.axios.request(config);
                this.items = response.data;
            } else {
                // TODO: load file
            }
            this.$emit("loading", false);
        },
        async deleteItem(item) {
            let confirmed = await this.$refs.confirm.open(
                "Delete",
                `Are you sure<br>you want to delete this ${
                    item.type === "dir" ? "folder" : "file"
                }?<br><em>${item.basename}</em>`
            );

            if (confirmed) {
                this.$emit("loading", true);
                let url = this.endpoints.delete.url
                    .replace(new RegExp("{storage}", "g"), this.storage)
                    .replace(new RegExp("{path}", "g"), item.path);

                let config = {
                    url,
                    method: this.endpoints.delete.method || "post"
                };

                await this.axios.request(config);
                this.$emit("file-deleted");
                this.$emit("loading", false);
            }
        }
    },
    watch: {
        async path() {
            this.items = [];
            await this.load();
        },
        async refreshPending() {
            if (this.refreshPending) {
                await this.load();
                this.$emit("refreshed");
            }
        }
    }
};
</script-->

<style lang="scss" scoped>
.v-card {
    height: 100%;
}
.icon-space {
    margin-right: 16px;
}
.file-icon {
    font-size: 3.75rem;
}
.folder-title {
  max-width: calc(100% - 30px); /* Adjust as needed */
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  vertical-align: middle;
}
.sort-btn{
    border-radius: 18px;
}

.drop-container {
  position: relative;
}
.drop-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(255,255,255,0.85);
  border: 2px solid rgba(var(--v-theme-primary));
  border-radius: 5px;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}
.drop-message {
  text-align: center;
  font-size: 1.5rem;
  color: #1976d2;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.upload-drag-icon{
    animation: bounce 1s infinite;
}
.sort-container{
    padding-top: 1rem;
}


@keyframes bounce {

0% ,50%, 100% {
    transform: translateY(0);
}

25%, 75% {
    transform: translateY(5px);
}

}
</style>
<style lang="scss">
.sort-select{
    .v-field__outline {
        display: none;
    }
    .v-field__overlay{
        opacity: 0!important;
        border-radius: 18px;
    }
    .v-field--variant-filled:hover .v-field__overlay {
        opacity: calc((0.04 + var(--v-hover-opacity)) * var(--v-theme-overlay-multiplier)) !important;
    }
}
</style>
