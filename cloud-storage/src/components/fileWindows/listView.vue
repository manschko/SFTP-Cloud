<script setup lang="ts">
import { usePreferencesStore } from '@/stores/preferences';
import { formatBytes, formatDate, getFileIcon } from '@/components/fileBrowser/util';
import { useFileStore } from '@/stores/fileStore';
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const preferencesStore = usePreferencesStore();
const fileStore = useFileStore();
const route = useRoute();
const router = useRouter();

function changeSort(key: 'name' | 'mod_time' | 'size') {
    if (preferencesStore.getSortBy() === key) {
        preferencesStore.toggleSortOrder();
    } else {
        preferencesStore.setSortOrder('asc');
        preferencesStore.setSortBy(key);
    }
}

function toggleItemSelection(item: any, event: MouseEvent) {
    const itemId = item.name;

    // Handle multi-select with Ctrl/Cmd key
    if (event.ctrlKey || event.metaKey) {
        if (fileStore.selectedItems.includes(itemId)) {
            fileStore.selectedItems = fileStore.selectedItems.filter(id => id !== itemId);
        } else {
            fileStore.selectedItems.push(itemId);
        }
    }
    // Handle range selection with Shift key
    else if (event.shiftKey && fileStore.selectedItems.length > 0) {
        const allItems = sortedAll.value;
        const lastSelectedId = fileStore.selectedItems[fileStore.selectedItems.length - 1];
        const lastSelectedIndex = allItems.findIndex(i => i.name === lastSelectedId);
        const currentIndex = allItems.findIndex(i => i.name === itemId);

        const start = Math.min(lastSelectedIndex, currentIndex);
        const end = Math.max(lastSelectedIndex, currentIndex);

        const newSelection = [...fileStore.selectedItems];
        for (let i = start; i <= end; i++) {
            const id = allItems[i].name;
            if (!newSelection.includes(id)) {
                newSelection.push(id);
            }
        }
        fileStore.selectedItems = newSelection;
    }
    else if (item.is_dir && fileStore.isItemSelected(item)) {
        //double click to open folder
        router.push(`${route.path}/${item.name}`);
    } else {
        fileStore.selectedItems = [itemId];
    }


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

</script>
<template>
    <v-table>
        <thead>
            <tr>
                <th class="text-left" scope="col">
                    <v-btn variant="text" @click="() => changeSort('name')" class="sort-btn">
                        Name
                        <v-icon :style="{ opacity: preferencesStore.getSortBy() === 'name' ? 1 : 0 }">
                            {{
                                preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down'
                            }}
                        </v-icon>
                    </v-btn>
                </th>
                <th class="text-left" scope="col">
                    <v-btn variant="text" @click="() => changeSort('mod_time')" class="sort-btn">
                        modified
                        <v-icon :style="{ opacity: preferencesStore.getSortBy() === 'mod_time' ? 1 : 0 }">
                            {{
                                preferencesStore.getSortOrder() === 'asc' ? 'arrow_drop_up' : 'arrow_drop_down'
                            }}
                        </v-icon>
                    </v-btn>
                </th>
                <th class="text-left" scope="col">
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
            <tr v-for="item in sortedAll" :key="item.name" @click="(event) => toggleItemSelection(item, event)">
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
</template>