import axios from "axios";
import { defineStore } from "pinia";
import { useRouter } from "vue-router";

export const useFileStore = defineStore('file',{
  state: () => ({
    folders: [] as any[],
    files: [] as any[],
    loading: false,
    currentPath: '',
    cache: {} as Record<string, any>,
    error: null
  }),
  getters: {
    getFiles: (state) => state.files,
    getFolders: (state) => state.folders,
    isLoading: (state) => state.loading,
    getCurrentPath: (state) => state.currentPath,
    getCache: (state) => state.cache,
    getError: (state) => state.error
  },

  actions: {
    async fetchFiles(path: string) {
      this.loading = true
      this.currentPath = path
      const router = useRouter()

      try {
        // Check if we have cached data
        if (path in this.cache) {
          console.log('Serving from cache:', path)
          if(!this.cache[path]) {
            this.folders = []
            this.files = []
          }else{
            this.folders = this.cache[path].filter((item: any) => item.is_dir)
            this.files = this.cache[path].filter((item: any) => !item.is_dir)
          }
        } else {
          console.log('Fetching from API:', path)
          const token = localStorage.getItem('jwt');
          const headers = token ? { Authorization: `Bearer ${token}`} : {};
          const response = await axios.get(`/api/files${path}` , { headers });
          // Store in cache
          this.cache[path] = response.data

          if(response.data) {
            this.folders = response.data.filter((item: any) => item.is_dir)
            this.files = response.data.filter((item: any) => !item.is_dir)
          }else{
            this.folders = []
            this.files = []
          }
        }
        this.loading = false
      } catch (error: any) {
        if (error.response && error.response.status === 401) {
          router.push({name: 'login'})
        }
        console.error('Error fetching files:', error)
        this.loading = false
        throw error
      }
    },

    clearCache(path?:string) {
      if (path) {
        delete this.cache[path]
      } else {
        this.cache = {}
      }
    },

    // Optional: Method to refresh data, bypassing cache
    async refreshFiles(path: string) {
      this.clearCache(path)
      await this.fetchFiles(path)
    }
  }
});
