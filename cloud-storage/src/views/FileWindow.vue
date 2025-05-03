<script setup lang="ts">
import {useFileStore} from '@/stores/fileStore';
import {computed, ref, watch} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import {usePreferencesStore} from "@/stores/preferences.ts";

interface Icons {
  [key: string]: string;
}

const fileStore = useFileStore();
const route = useRoute();
const router = useRouter();
const currentPath = computed(() => {
  return route.params.pathMatch
    ? Array.isArray(route.params.pathMatch)
      ? '/' + route.params.pathMatch.join('/')
      : '/' + route.params.pathMatch
    : '/';
});

// Add a watcher to fetch files when the path changes
watch(currentPath, (newPath) => {
  fileStore.fetchFiles(newPath).catch(() => {
    window.history.back();
  });
}, { immediate: true });

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

// Track selected items
const selectedItems = ref<string[]>([]);

// Selection functions
function toggleItemSelection(item: any, event: MouseEvent) {
  const itemId = item.name;

  // Handle multi-select with Ctrl/Cmd key
  if (event.ctrlKey || event.metaKey) {
    if (selectedItems.value.includes(itemId)) {
      selectedItems.value = selectedItems.value.filter(id => id !== itemId);
    } else {
      selectedItems.value.push(itemId);
    }
  }
  // Handle range selection with Shift key
  else if (event.shiftKey && selectedItems.value.length > 0) {
    const allItems = sortedAll.value;
    const lastSelectedId = selectedItems.value[selectedItems.value.length - 1];
    const lastSelectedIndex = allItems.findIndex(i => i.name === lastSelectedId);
    const currentIndex = allItems.findIndex(i => i.name === itemId);

    const start = Math.min(lastSelectedIndex, currentIndex);
    const end = Math.max(lastSelectedIndex, currentIndex);

    const newSelection = [...selectedItems.value];
    for (let i = start; i <= end; i++) {
      const id = allItems[i].name;
      if (!newSelection.includes(id)) {
        newSelection.push(id);
      }
    }
    selectedItems.value = newSelection;
  }
  else if (item.is_dir && isItemSelected(item)) {
    const newPath = currentPath.value + (currentPath.value.endsWith('/') ? '' : '/') + item.name;
    router.push({ path: `/files${newPath}` });
    selectedItems.value = selectedItems.value.filter(id => id !== itemId);
  } else {
    selectedItems.value = [itemId];
  }
  

}

// Handle double-click for folder navigation
/*function handleDoubleClick(item: any) {
  if (item.is_dir) {
    const newPath = currentPath.value + (currentPath.value.endsWith('/') ? '' : '/') + item.name;
    router.push({ path: `/files${newPath}` });
  }
  // You could add file handling logic here (preview, download, etc.)
}*/

// Check if an item is selected
function isItemSelected(item: any): boolean {
  return selectedItems.value.includes(item.name);
}

// Clear selection when path changes
watch(() => currentPath, () => {
  selectedItems.value = [];
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
    <v-card-text v-if="fileStore.files.length || fileStore.folders.length" class="grow">
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
                  {{
                    preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down'
                  }}
                </v-icon>
              </v-btn>
            </th>
            <th class="text-left">
              <v-btn variant="text" @click="() => changeSort('mod_time')" class="sort-btn">
                modified
                <v-icon :style="{ opacity: preferencesStore.getSortBy() === 'mod_time' ? 1 : 0 }">
                  {{
                    preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down'
                  }}
                </v-icon>
              </v-btn>
            </th>
            <th class="text-left">
              <v-btn variant="text" @click="() => changeSort('size')" class="sort-btn">
                size
                <v-icon :style="{ opacity: preferencesStore.getSortBy() === 'size' ? 1 : 0 }">
                  {{
                    preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down'
                  }}
                </v-icon>
              </v-btn>
            </th>
          </tr>
          </thead>
          <tbody>
          <tr
            v-for="item in sortedAll"
            :key="item.name"
            @click="(event) => toggleItemSelection(item, event)"
            @dblclick="() => handleDoubleClick(item)"
          >
            <td>
              <v-icon class="icon-space">{{ item.is_dir ? 'folder' : getFileIcon(item.name) }}
              </v-icon>
              {{ item.name }}
            </td>
            <td>{{ formatDate(item.mod_time) }}</td>
            <td>{{ item.is_dir ? '' : formatBytes(item.size) }}</td>
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

          <p v-if="fileStore.folders.length" class="text-subtitle-1 mb-2 mt-2">Ordner</p>
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
            >
              <v-card
                ripple
                variant="tonal"
                :class="{ 'selection-indicator': isItemSelected(folder), 'mb-2': true }"
                @click="(event) => toggleItemSelection(folder, event)"
                @dblclick="() => handleDoubleClick(folder)">
                <v-card-title class="d-flex align-center">
                  <v-icon class="mr-2">folder</v-icon>
                  <span class="folder-title text-truncate">{{ folder.name }}</span>
                </v-card-title>
              </v-card>
            </v-col>
          </v-row>
          <p v-if="fileStore.files.length" class="text-subtitle-1 mb-2 mt-2">Dateien</p>
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
            >
              <v-card
                variant="tonal"
                ripple
                :class="{ 'selection-indicator': isItemSelected(file), 'mb-2': true }"
                @click="(event) => toggleItemSelection(file, event)"
                @dblclick="() => handleDoubleClick(file)"
              >
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
    >The folder is empty
    </v-card-text>
  </v-card>
</template>

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

.sort-btn {
  border-radius: 18px;
}

.drop-container {
  position: relative;
}

.drop-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.85);
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

.upload-drag-icon {
  animation: bounce 1s infinite;
}

.sort-container {
  padding-top: 1rem;
}


@keyframes bounce {

  0%, 50%, 100% {
    transform: translateY(0);
  }

  25%, 75% {
    transform: translateY(5px);
  }

}

.v-table__selected {
  background-color: rgb(var(--v-theme-primary), 0.1) !important;
}
.selection-indicator{
  color: rgb(var(--v-theme-primary)) !important;
  background-color: rgba(var(--v-theme-primary), calc((var(--v-activated-opacity) + var(--v-hover-opacity)) * var(--v-theme-overlay-multiplier)))!important;
}
</style>
<style lang="scss">
.sort-select {
  .v-field__outline {
    display: none;
  }

  .v-field__overlay {
    opacity: 0 !important;
    border-radius: 18px;
  }

  .v-field--variant-filled:hover .v-field__overlay {
    opacity: calc((0.04 + var(--v-hover-opacity)) * var(--v-theme-overlay-multiplier)) !important;
  }
}
</style>
