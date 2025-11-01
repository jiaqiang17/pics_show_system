<template>
  <div class="nail-discover-public">
    <header class="discover-header">
      <button class="header-icon sparkle" type="button" aria-label="featured">
        &#10023;
      </button>
      <h1 class="discover-title">发现</h1>
      <button class="header-icon menu" type="button" aria-label="菜单">
        <span class="menu-line"></span>
        <span class="menu-line"></span>
        <span class="menu-line"></span>
      </button>
    </header>

    <section class="category-tabs">
      <div class="tabs-scroll">
        <button
          class="tab-item"
          type="button"
          :class="{ active: activeTagId === null }"
          @click="selectTag(null)"
        >
          为你推荐
        </button>
        <button
          v-for="tag in tags"
          :key="tag.ID"
          class="tab-item"
          type="button"
          :class="{ active: activeTagId === tag.ID }"
          @click="selectTag(tag.ID)"
        >
          {{ tag.tagName }}
        </button>
      </div>
    </section>

    <div class="tabs-divider"></div>

    <section class="gallery-section">
      <div v-if="loading" class="gallery-loading">
        <div v-for="index in 6" :key="index" class="skeleton-card"></div>
      </div>

      <div v-else-if="!loading && nailList.length === 0" class="gallery-empty">
        <span class="empty-icon">&#128133;</span>
        <p class="empty-text">暂无美甲款式</p>
      </div>

      <div v-else class="gallery-grid">
        <div
          v-for="nail in nailList"
          :key="nail.ID"
          class="gallery-card"
        >
          <img
            v-if="nail.images && nail.images.length"
            :src="getImageUrl(nail.images[0])"
            :alt="nail.styleName"
            class="gallery-image"
            loading="lazy"
            @error="handleImageError"
          />
          <div v-else class="gallery-fallback">&#128248;</div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useGalleryStore } from '@/pinia/modules/gallery'
import { returnArrImg } from '@/utils/format'

const galleryStore = useGalleryStore()

const tags = computed(() =>
  galleryStore.tags.map((item) => ({
    ID: item.value,
    tagName: item.label
  }))
)
const activeTagId = ref(null)
const nailList = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)

const transformRecord = (item) => {
  const parseArray = (value) => {
    if (!value) return []
    if (Array.isArray(value)) return value
    if (typeof value === 'string') {
      try {
        const parsed = JSON.parse(value)
        return Array.isArray(parsed) ? parsed : []
      } catch (error) {
        return []
      }
    }
    return []
  }
  return {
    ...item,
    images: parseArray(item.images)
  }
}

const loadTags = async () => {
  await galleryStore.fetchTags()
}

const loadNailList = async () => {
  loading.value = true
  try {
    galleryStore.updateFilters({
      page: page.value,
      pageSize: pageSize.value,
      status: 'enabled',
      tagIds: activeTagId.value ? [activeTagId.value] : [],
      matchAll: false,
      isRecommended: activeTagId.value ? null : true,
      sort: 'sort',
      order: ''
    })
    const { list } = await galleryStore.fetchStyles({ status: 'enabled' })
    nailList.value = (list || []).map(transformRecord)
  } finally {
    loading.value = false
  }
}

const selectTag = (tagId) => {
  if (activeTagId.value === tagId) return
  activeTagId.value = tagId
  page.value = 1
  loadNailList()
}

const getImageUrl = (image) => {
  const list = returnArrImg(image)
  return list.length ? list[0] : ''
}

const handleImageError = (event) => {
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200"><rect fill="#f5f5f5" width="200" height="200"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-size="60">&#128133;</text></svg>`
  event.target.src = `data:image/svg+xml,${encodeURIComponent(svg)}`
}

onMounted(async () => {
  await loadTags()
  loadNailList()
})
</script>

<style scoped lang="scss">
.nail-discover-public {
  box-sizing: border-box;
  min-height: 100vh;
  padding: 32px 20px 48px;
  background: linear-gradient(180deg, #ffffff 0%, #f5f5f7 100%);
  color: #202124;
  font-family: 'PingFang SC', 'Helvetica Neue', Arial, sans-serif;

  .discover-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;

    .header-icon {
      width: 44px;
      height: 44px;
      border-radius: 16px;
      border: none;
      background: #ffffff;
      box-shadow: 0 10px 24px rgba(0, 0, 0, 0.08);
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      padding: 0;

      &.sparkle {
        color: #ff4f9a;
        font-size: 22px;
        font-weight: 600;
      }

      &.menu {
        flex-direction: column;
        gap: 6px;

        .menu-line {
          width: 18px;
          height: 2px;
          border-radius: 999px;
          background: #202124;
        }
      }
    }

    .discover-title {
      font-size: 22px;
      font-weight: 700;
      letter-spacing: 2px;
      color: #111111;
    }
  }

  .category-tabs {
    margin-bottom: 12px;

    .tabs-scroll {
      display: flex;
      gap: 12px;
      overflow-x: auto;
      padding-bottom: 6px;

      &::-webkit-scrollbar {
        display: none;
      }

      .tab-item {
        flex-shrink: 0;
        min-width: 96px;
        padding: 10px 18px;
        border-radius: 22px;
        border: none;
        background: #ffffff;
        color: #4a4a4a;
        font-size: 14px;
        font-weight: 500;
        box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
        cursor: pointer;
        transition: transform 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;

        &:active {
          transform: scale(0.98);
        }

        &.active {
          color: #ffffff;
          background: linear-gradient(120deg, #ff61d2 0%, #fe9090 100%);
          box-shadow: 0 10px 28px rgba(254, 110, 168, 0.35);
        }
      }
    }
  }

  .tabs-divider {
    height: 1px;
    background: #e4e6eb;
    border-radius: 999px;
    margin-bottom: 24px;
  }

  .gallery-section {
    .gallery-loading {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 18px;

      .skeleton-card {
        border-radius: 28px;
        aspect-ratio: 1 / 1;
        background: linear-gradient(135deg, #f3f4f7 0%, #e8ebf2 100%);
        animation: pulse 1.5s ease-in-out infinite;
      }
    }

    .gallery-empty {
      margin-top: 60px;
      text-align: center;
      color: #9aa0a6;

      .empty-icon {
        font-size: 48px;
        display: block;
        margin-bottom: 12px;
      }

      .empty-text {
        font-size: 16px;
      }
    }

    .gallery-grid {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 18px;
    }

    .gallery-card {
      background: #ffffff;
      border-radius: 28px;
      padding: 14px;
      box-shadow: 0 16px 32px rgba(15, 23, 42, 0.12);
      display: flex;
      align-items: center;
      justify-content: center;
      transition: transform 0.25s ease, box-shadow 0.25s ease;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 20px 40px rgba(15, 23, 42, 0.16);
      }

      .gallery-image {
        width: 100%;
        border-radius: 22px;
        aspect-ratio: 1 / 1;
        object-fit: cover;
        display: block;
      }

      .gallery-fallback {
        width: 100%;
        aspect-ratio: 1 / 1;
        border-radius: 22px;
        background: #f5f5f7;
        color: #c5c5c5;
        font-size: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
  }
}

@keyframes pulse {
  0% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.6;
  }
}
</style>
