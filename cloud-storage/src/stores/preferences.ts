import {defineStore} from "pinia";

export type Preferences = {
  viewMode: string;
  sortOrder: string;
  sortBy: string;
  name: string
}

export const usePreferencesStore = defineStore('preferences', {
  state: () => ({
    preferences: {
      viewMode: 'list',
      sortOrder: 'desc',
      sortBy: 'name',
      ...(localStorage.getItem('preferences') ? JSON.parse(localStorage.getItem('preferences') as string) : {})
    } as Preferences,

  }),
  actions: {
    setViewMode(viewMode: 'grid' | 'list') {
      this.preferences.viewMode = viewMode
      localStorage.setItem('preferences', JSON.stringify(this.preferences))
    },
    getViewMode() {
      return this.preferences.viewMode
    },
    setSortOrder(sortOrder: 'asc' | 'desc') {
      this.preferences.sortOrder = sortOrder
      localStorage.setItem('preferences', JSON.stringify(this.preferences))
    },
    toggleSortOrder() {
      this.preferences.sortOrder = this.preferences.sortOrder === 'asc' ? 'desc' : 'asc'
      localStorage.setItem('preferences', JSON.stringify(this.preferences))
    },
    getSortOrder() {
      return this.preferences.sortOrder
    },
    setSortBy(sortBy: 'name' | 'mod_time' | 'size') {
      this.preferences.sortBy = sortBy
      localStorage.setItem('preferences', JSON.stringify(this.preferences))
    },
    getSortBy() {
      return this.preferences.sortBy
    },
    setName(name: string) {
      this.preferences.name = name
      localStorage.setItem('preferences', JSON.stringify(this.preferences))
    },
    getName() {
      return this.preferences.name
    },
  }
});
