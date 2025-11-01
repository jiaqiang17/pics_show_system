import { defineStore } from 'pinia'
import { getNailStyleDataSource, getNailStyleList } from '@/api/nail/nail_style'

export const useGalleryStore = defineStore('gallery', {
  state: () => ({
    tags: [],
    list: [],
    total: 0,
    loading: false,
    filters: {
      tagIds: [],
      matchAll: false,
      isRecommended: null,
      status: '',
      page: 1,
      pageSize: 20,
      sort: '',
      order: ''
    }
  }),
  getters: {
    tagMap(state) {
      return state.tags.reduce((map, item) => {
        map.set(item.value, item.label)
        return map
      }, new Map())
    }
  },
  actions: {
    async fetchTags(force = false) {
      if (this.tags.length && !force) {
        return this.tags
      }
      const res = await getNailStyleDataSource()
      if (res.code === 0) {
        const tagList = res.data?.tagIds || []
        this.tags = tagList.map((item) => ({
          label: item.label,
          value: item.value
        }))
      }
      return this.tags
    },
    async fetchStyles(extraFilters = {}) {
      this.loading = true
      try {
        const params = {
          ...this.filters,
          ...extraFilters
        }
        if (!params.tagIds?.length) {
          delete params.matchAll
        }
        if (params.isRecommended === null) {
          delete params.isRecommended
        }
        const res = await getNailStyleList(params)
        if (res.code === 0) {
          this.list = res.data?.list || []
          this.total = res.data?.total || 0
        }
      } finally {
        this.loading = false
      }
      return {
        list: this.list,
        total: this.total
      }
    },
    updateFilters(patch) {
      this.filters = {
        ...this.filters,
        ...patch
      }
    }
  }
})
