import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', () => {
  const count = ref(0)
  const doubleCount = computed(() => count.value * 2)
  function increment() {
    count.value++
  }

  return { count, doubleCount, increment }
})

export const useFileStore = defineStore('file',{
  state: () => ({
    files: [],
    loading: false,
    currentPath: '',
    cache: {} as Record<string, any>,
    error: null
  }),
  getters: {
    getFiles: (state) => state.files,
    isLoading: (state) => state.loading,
    getCurrentPath: (state) => state.currentPath,
    getCache: (state) => state.cache,
    getError: (state) => state.error
  },

  actions: {
    async fetchFiles(path: string) {
      this.loading = true
      this.currentPath = path
      this.error = null
      
      try {
        // Check if we have cached data
        if (path in this.cache) {
          console.log('Serving from cache:', path)
          this.files = this.cache[path]
        } else {
          console.log('Fetching from API:', path)
          const response = await axios.get(`/api/files${path}`)
          // Store in cache
          this.cache[path] = response.data
          this.files = response.data
        }
      } catch (error) {
        console.error('Error fetching files:', error)
        this.error = error.message || 'Failed to fetch files'
        this.files = []
      } finally {
        this.loading = false
      }
    },
    
    clearCache(path = null) {
      if (path) {
        delete this.cache[path]
      } else {
        this.cache = {}
      }
    },
    
    invalidateCache(path) {
      if (path in this.cache) {
        delete this.cache[path]
      }
    },
    
    // Optional: Method to refresh data, bypassing cache
    async refreshFiles(path) {
      this.invalidateCache(path)
      await this.fetchFiles(path)
    }
  }
});