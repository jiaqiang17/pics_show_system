<template>
  <div ref="containerRef" class="nail-discover-public">
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
      <div
        class="tabs-scroll-wrap"
        :class="{
          'is-scrollable': tabsScrollable,
          'at-start': tabsAtStart,
          'at-end': tabsAtEnd,
          'is-hinting': showTabScrollHint
        }"
      >
        <div
          ref="tabsScrollRef"
          class="tabs-scroll"
          @scroll.passive="handleTabsScroll"
          @touchstart.passive="hideTabsHint"
          @wheel.passive="hideTabsHint"
        >
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

        <div v-if="tabsScrollable" class="tabs-scroll-affordance" aria-hidden="true">
          <div v-show="showTabsLeft" class="tabs-edge tabs-edge--left">
            <svg viewBox="0 0 24 24" aria-hidden="true" focusable="false">
              <path
                d="M15 18l-6-6 6-6"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
          <div v-show="showTabsRight" class="tabs-edge tabs-edge--right">
            <svg viewBox="0 0 24 24" aria-hidden="true" focusable="false">
              <path
                d="M9 6l6 6-6 6"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
        </div>

        <div v-if="tabsScrollable && showTabScrollHint" class="tabs-scroll-hint" aria-hidden="true">
          左右滑动查看更多
        </div>
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

      <div v-else class="gallery-waterfall">
        <div
          v-for="nail in nailList"
          :key="nail.ID"
          class="gallery-card"
          :ref="(el) => registerNailCard(el, nail.ID)"
        >
          <button
            v-if="isAdminMode"
            class="card-delete"
            type="button"
            aria-label="删除"
            @click.stop="confirmDeleteStyle(nail)"
          >
            ×
          </button>
          <div class="gallery-media">
            <img
              v-if="isCardVisible(nail.ID) && nail.images && nail.images.length"
              :src="getImageUrl(nail.images[0])"
              :alt="nail.styleName"
              class="gallery-image"
              loading="lazy"
              @error="handleImageError"
            />
            <div v-else class="gallery-placeholder">
              <div class="placeholder-shimmer"></div>
              <div class="placeholder-icon">&#128248;</div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="!loading && nailList.length" class="gallery-footer">
        <div v-if="loadingMore" class="gallery-footer__status">加载中…</div>
        <div v-else-if="!hasMore" class="gallery-footer__status">到底啦</div>
        <div v-else class="gallery-footer__status">继续下滑加载更多</div>
        <div v-if="hasMore" ref="sentinelRef" class="gallery-sentinel"></div>
      </div>
    </section>

    <button
      v-if="canShowAddFab"
      class="fab-add"
      type="button"
      aria-label="添加款式"
      @click="openAddDialog"
    >
      +
    </button>

    <el-dialog
      v-model="addDialogVisible"
      :title="`添加款式${activeTagName ? ` · ${activeTagName}` : ''}`"
      :width="addDialogWidth"
      :fullscreen="isMobile"
      :show-close="false"
      :close-on-click-modal="false"
      destroy-on-close
      append-to-body
      class="add-style-dialog"
    >
      <template #header="{ close, titleId, titleClass }">
        <div class="add-style-header">
          <div class="add-style-header__left">
            <div class="add-style-badge">NEW</div>
            <div class="add-style-titles">
              <div :id="titleId" :class="titleClass" class="add-style-title">
                <span class="add-style-title__text">添加款式</span>
                <button class="add-style-title__close" type="button" aria-label="关闭" @click="close">
                  <span class="add-style-title__close-x">×</span>
                </button>
              </div>
              <div class="add-style-subtitle">
                添加到 <span class="add-style-subtitle__tag">{{ activeTagName || '当前标签' }}</span>
              </div>
            </div>
          </div>
        </div>
      </template>

      <div class="add-style-surface">
        <div class="add-style-meta">
          <div class="add-style-chip">默认：不推荐</div>
          <div class="add-style-chip">状态：启用</div>
          <div class="add-style-chip">标签：{{ activeTagName || '当前标签' }}</div>
        </div>

        <el-form label-position="top" class="add-style-form">
          <el-form-item label="款式名称（可不填）" class="premium-item">
            <el-input
              v-model="addForm.styleName"
              placeholder="给它一个名字（例如：法式、渐变、猫眼…）"
              clearable
              maxlength="50"
              show-word-limit
              class="premium-input"
            />
          </el-form-item>

          <el-form-item label="图片（点击后裁剪上传，默认 1 张）" required class="premium-item cropper-item">
            <div class="cropper-single">
              <button class="cropper-preview" type="button" @click="openCropper">
                <img
                  v-if="addForm.image"
                  :src="getSingleImageUrl(addForm.image)"
                  alt="预览"
                  class="cropper-preview-img"
                />
                <div v-else class="cropper-placeholder">
                  <div class="cropper-orb"></div>
                  <div class="cropper-placeholder__title">一张图，定乾坤</div>
                  <div class="cropper-placeholder__desc">点击上传 · 智能裁剪 · 高清展示</div>
                  <div class="cropper-cta">上传并裁剪</div>
                </div>
              </button>

              <div class="cropper-actions">
                <el-button class="btn-ghost" @click="openCropper">
                  {{ addForm.image ? '重新裁剪/更换' : '选择图片' }}
                </el-button>
                <el-button v-if="addForm.image" class="btn-danger-ghost" @click="removeCropImage">移除</el-button>
              </div>

              <CropperImage
                ref="cropperRef"
                :hide-trigger="true"
                dialog-title="图片裁剪"
                :default-ratio-index="0"
                :default-crop-size="420"
                :default-crop-size-mobile="320"
                :initial-scale-steps="-10"
                :show-preview="false"
                upload-path="/fileUploadAndDownload/uploadPublic"
                @on-success="handleCropSuccess"
              />
            </div>
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <div class="add-style-footer">
          <el-button class="btn-footer-secondary" @click="closeAddDialog">取消</el-button>
          <el-button
            class="btn-footer-primary"
            type="primary"
            :loading="addSubmitting"
            @click="submitAddStyle"
          >
            {{ addSubmitting ? '正在添加…' : '一键添加' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGalleryStore } from '@/pinia/modules/gallery'
import { createNailStyle, deleteNailStyle } from '@/api/nail/nail_style'
import { returnArrImg } from '@/utils/format'
import CropperImage from '@/components/upload/cropper.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const galleryStore = useGalleryStore()

const isAdminMode = computed(() => String(route.query?.isAdmin || '') === '1')

const tags = computed(() =>
  galleryStore.tags.map((item) => ({
    ID: item.value,
    tagName: item.label
  }))
)
const activeTagId = ref(null)
const activeTagName = computed(() => {
  if (!activeTagId.value) return ''
  return tags.value.find((tag) => tag.ID === activeTagId.value)?.tagName || ''
})
const canShowAddFab = computed(() => !!activeTagId.value && isAdminMode.value)

const tabsScrollRef = ref(null)
const tabsScrollable = ref(false)
const tabsAtStart = ref(true)
const tabsAtEnd = ref(false)
const showTabScrollHint = ref(false)
const tabsUserInteracted = ref(false)
const showTabsLeft = computed(
  () => tabsScrollable.value && tabsUserInteracted.value && !tabsAtStart.value
)
const showTabsRight = computed(() => tabsScrollable.value && !tabsAtEnd.value)
let tabsHintTimer = null
let tabsHintSeen = false

const resetTabsScrollToStart = async () => {
  await nextTick()
  const applyStart = () => {
    const el = tabsScrollRef.value
    if (!el) return
    tabsUserInteracted.value = false
    el.scrollLeft = 0
    updateTabsScrollState()
  }

  applyStart()
  if (typeof requestAnimationFrame !== 'undefined') {
    requestAnimationFrame(applyStart)
    requestAnimationFrame(() => requestAnimationFrame(applyStart))
  }
  setTimeout(applyStart, 50)
  setTimeout(applyStart, 200)
}

const updateTabsScrollState = () => {
  const el = tabsScrollRef.value
  if (!el) return
  const scrollable = el.scrollWidth > el.clientWidth + 4
  tabsScrollable.value = scrollable
  tabsAtStart.value = el.scrollLeft <= 2
  tabsAtEnd.value = el.scrollLeft + el.clientWidth >= el.scrollWidth - 2
  if (tabsAtStart.value) {
    tabsUserInteracted.value = false
  }

  if (scrollable && !tabsHintSeen) {
    tabsHintSeen = true
    showTabScrollHint.value = true
    if (tabsHintTimer) clearTimeout(tabsHintTimer)
    tabsHintTimer = setTimeout(() => {
      showTabScrollHint.value = false
    }, 2200)
  }
}

const hideTabsHint = () => {
  showTabScrollHint.value = false
  if (tabsHintTimer) {
    clearTimeout(tabsHintTimer)
    tabsHintTimer = null
  }
}

const handleTabsScroll = () => {
  const el = tabsScrollRef.value
  if (el && el.scrollLeft > 2) tabsUserInteracted.value = true
  updateTabsScrollState()
  if (!tabsAtStart.value) hideTabsHint()
}
const containerRef = ref(null)
const sentinelRef = ref(null)
const nailList = ref([])
const loading = ref(false)
const loadingMore = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const hasMore = computed(() => nailList.value.length < (total.value || 0))

const visibleMap = reactive({})
const resetVisibleMap = () => {
  Object.keys(visibleMap).forEach((key) => {
    delete visibleMap[key]
  })
}
const isCardVisible = (id) => !!visibleMap[String(id)]

let imageObserver = null
let loadMoreObserver = null
const cardElementMap = new Map()
const registerNailCard = (el, id) => {
  const key = String(id)
  const previous = cardElementMap.get(key)

  if (!el) {
    if (previous) imageObserver?.unobserve?.(previous)
    cardElementMap.delete(key)
    return
  }

  if (previous && previous !== el) {
    imageObserver?.unobserve?.(previous)
  }
  el.__nailId = key
  cardElementMap.set(key, el)
  imageObserver?.observe?.(el)
}

const setupObservers = async () => {
  await nextTick()
  const rootEl = containerRef.value
  if (!rootEl || typeof IntersectionObserver === 'undefined') {
    nailList.value.forEach((item) => {
      if (!item?.ID) return
      visibleMap[String(item.ID)] = true
    })
    return
  }

  imageObserver?.disconnect?.()
  imageObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return
        const key = entry.target?.__nailId ? String(entry.target.__nailId) : null
        if (key) visibleMap[key] = true
        imageObserver?.unobserve?.(entry.target)
      })
    },
    {
      root: rootEl,
      rootMargin: '120px 0px',
      threshold: 0.01
    }
  )

  cardElementMap.forEach((value) => {
    if (!value) return
    imageObserver.observe(value)
  })

  loadMoreObserver?.disconnect?.()
  if (sentinelRef.value) {
    loadMoreObserver = new IntersectionObserver(
      (entries) => {
        if (!entries.some((entry) => entry.isIntersecting)) return
        if (!hasMore.value) return
        if (loading.value || loadingMore.value) return
        triggerLoadMore()
      },
      {
        root: rootEl,
        rootMargin: '600px 0px',
        threshold: 0
      }
    )
    loadMoreObserver.observe(sentinelRef.value)
  }
}

const addDialogVisible = ref(false)
const addSubmitting = ref(false)
const addForm = reactive({
  styleName: '',
  image: ''
})
const cropperRef = ref(null)

const isMobile = ref(false)
let mediaQuery = null
const updateIsMobile = () => {
  isMobile.value = !!mediaQuery?.matches
}
const addDialogWidth = computed(() => (isMobile.value ? '100%' : '460px'))

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
      sort: 'created_at',
      order: 'descending'
    })
    const { list, total: nextTotal } = await galleryStore.fetchStyles({ status: 'enabled' })
    total.value = Number(nextTotal || 0)
    nailList.value = (list || []).map(transformRecord)
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  if (!hasMore.value) return
  loadingMore.value = true
  try {
    const nextPage = page.value + 1
    galleryStore.updateFilters({
      page: nextPage,
      pageSize: pageSize.value,
      status: 'enabled',
      tagIds: activeTagId.value ? [activeTagId.value] : [],
      matchAll: false,
      isRecommended: activeTagId.value ? null : true,
      sort: 'created_at',
      order: 'descending'
    })
    const { list, total: nextTotal } = await galleryStore.fetchStyles({ status: 'enabled' })
    const nextList = (list || []).map(transformRecord)
    total.value = Number(nextTotal || 0)
    page.value = nextPage
    nailList.value = nailList.value.concat(nextList)
  } finally {
    loadingMore.value = false
  }
}

const triggerLoadMore = async () => {
  const sentinelEl = sentinelRef.value
  if (sentinelEl) loadMoreObserver?.unobserve?.(sentinelEl)
  await loadMore()
  await nextTick()
  if (sentinelRef.value && hasMore.value) {
    loadMoreObserver?.observe?.(sentinelRef.value)
  }
}

const resetAndLoad = async () => {
  page.value = 1
  total.value = 0
  nailList.value = []
  resetVisibleMap()
  await loadNailList()
  await setupObservers()
}

const parseTagId = (value) => {
  const numberValue = Number(value)
  return Number.isFinite(numberValue) && numberValue > 0 ? numberValue : null
}

const syncTagQuery = (tagId) => {
  const nextQuery = { ...route.query }
  if (tagId) {
    nextQuery.tagId = String(tagId)
  } else {
    delete nextQuery.tagId
  }
  router.replace({ query: nextQuery })
}

const selectTag = (tagId) => {
  if (activeTagId.value === tagId) return
  activeTagId.value = tagId
  syncTagQuery(tagId)
  resetAndLoad()
}

const getImageUrl = (image) => {
  const list = returnArrImg(image)
  return list.length ? list[0] : ''
}

const getSingleImageUrl = (image) => {
  const list = returnArrImg(image)
  return list.length ? list[0] : ''
}

const handleImageError = (event) => {
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200"><rect fill="#f5f5f5" width="200" height="200"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-size="60">&#128133;</text></svg>`
  event.target.src = `data:image/svg+xml,${encodeURIComponent(svg)}`
}

const openCropper = () => {
  cropperRef.value?.open?.()
}

const handleCropSuccess = (url) => {
  addForm.image = url || ''
}

const confirmDeleteStyle = async (nail) => {
  if (!isAdminMode.value) return
  try {
    await ElMessageBox.confirm('确认删除该款式？删除后不可恢复。', '提示', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
      confirmButtonClass: 'btn-danger-confirm'
    })
  } catch {
    return
  }

  try {
    const res = await deleteNailStyle({ ID: nail.ID })
    if (!res || res.code !== 0) {
      throw new Error(res?.msg || '删除失败')
    }
    nailList.value = nailList.value.filter((item) => item.ID !== nail.ID)
    total.value = Math.max(0, (total.value || 0) - 1)
    if (hasMore.value && !loadingMore.value) {
      await loadMore()
    }
    ElMessage.success('已删除')
  } catch (error) {
    console.error('删除失败', error)
    ElMessage.error('删除失败，请重试')
  }
}

const removeCropImage = () => {
  addForm.image = ''
}

const resetAddForm = () => {
  addForm.styleName = ''
  addForm.image = ''
}

const openAddDialog = () => {
  if (!isAdminMode.value) {
    ElMessage.warning('无权限操作')
    return
  }
  if (!activeTagId.value) {
    ElMessage.warning('请先选择一个标签')
    return
  }
  addDialogVisible.value = true
}

const closeAddDialog = () => {
  addDialogVisible.value = false
  resetAddForm()
}

const submitAddStyle = async () => {
  if (!activeTagId.value) return
  if (!addForm.image) {
    ElMessage.warning('请先上传图片')
    return
  }

  addSubmitting.value = true
  try {
    const payload = {
      styleName: addForm.styleName?.trim() || '',
      images: [addForm.image],
      isRecommended: false,
      status: 'enabled',
      tagIds: [activeTagId.value]
    }

    const res = await createNailStyle(payload)
    if (res.code !== 0) {
      ElMessage.error(res.msg || '添加失败')
      return
    }
    ElMessage.success('添加成功')
    closeAddDialog()
    await resetAndLoad()
  } finally {
    addSubmitting.value = false
  }
}

onMounted(async () => {
  if (typeof window !== 'undefined' && window.matchMedia) {
    mediaQuery = window.matchMedia('(max-width: 768px)')
    updateIsMobile()
    if (mediaQuery.addEventListener) {
      mediaQuery.addEventListener('change', updateIsMobile)
    } else if (mediaQuery.addListener) {
      mediaQuery.addListener(updateIsMobile)
    }
  }
  await loadTags()
  await nextTick()
  await resetTabsScrollToStart()
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', updateTabsScrollState)
  }
  const initialTagId = parseTagId(route.query.tagId)
  if (initialTagId) {
    activeTagId.value = initialTagId
  }
  await resetAndLoad()
})

onBeforeUnmount(() => {
  imageObserver?.disconnect?.()
  loadMoreObserver?.disconnect?.()
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', updateTabsScrollState)
  }
  hideTabsHint()
  if (!mediaQuery) return
  if (mediaQuery.removeEventListener) {
    mediaQuery.removeEventListener('change', updateIsMobile)
  } else if (mediaQuery.removeListener) {
    mediaQuery.removeListener(updateIsMobile)
  }
})
</script>

<style scoped lang="scss">
.nail-discover-public {
  box-sizing: border-box;
  height: 100vh;
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
  padding: 32px 20px 48px;
  background: linear-gradient(180deg, #ffffff 0%, #f5f5f7 100%);
  color: #202124;
  font-family: 'PingFang SC', 'Helvetica Neue', Arial, sans-serif;
  --gallery-gap: clamp(12px, 3vw, 16px);

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

    .tabs-scroll-wrap {
      position: relative;
      border-radius: 26px;
      padding: 6px;
      background: rgba(255, 255, 255, 0.72);
      border: 1px solid rgba(15, 23, 42, 0.06);
      box-shadow: 0 12px 28px rgba(15, 23, 42, 0.06);
      backdrop-filter: blur(14px);
      overflow: hidden;

      &::before,
      &::after {
        content: '';
        position: absolute;
        top: 0;
        bottom: 0;
        width: 46px;
        pointer-events: none;
        opacity: 0;
        transition: opacity 0.22s ease;
        z-index: 2;
      }

      &.is-scrollable::before {
        left: 0;
        opacity: 1;
        background: linear-gradient(90deg, rgba(255, 255, 255, 0.98) 0%, rgba(255, 255, 255, 0) 88%);
      }

      &.is-scrollable::after {
        right: 0;
        opacity: 1;
        background: linear-gradient(270deg, rgba(255, 255, 255, 0.98) 0%, rgba(255, 255, 255, 0) 88%);
      }

      &.at-start::before {
        opacity: 0;
      }

      &.at-end::after {
        opacity: 0;
      }
    }

    .tabs-scroll {
      display: flex;
      gap: 12px;
      overflow-x: auto;
      padding: 6px 0;
      scroll-snap-type: x proximity;
      -webkit-overflow-scrolling: touch;

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
        scroll-snap-align: start;

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

    .tabs-scroll-affordance {
      position: absolute;
      top: 6px;
      bottom: 6px;
      left: 0;
      right: 0;
      z-index: 3;
      pointer-events: none;
    }

    .tabs-edge {
      position: absolute;
      top: 50%;
      transform: translateY(-50%);
      width: 28px;
      height: 28px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: rgba(15, 23, 42, 0.5);
      background: rgba(255, 255, 255, 0.78);
      border: 1px solid rgba(15, 23, 42, 0.08);
      box-shadow: 0 12px 26px rgba(15, 23, 42, 0.1);
      opacity: 0.9;
      transition: opacity 0.2s ease, transform 0.2s ease;
      user-select: none;
    }

    .tabs-edge svg {
      width: 16px;
      height: 16px;
      display: block;
    }

    .tabs-edge--left {
      left: 8px;
    }

    .tabs-edge--right {
      right: 8px;
    }

    .tabs-scroll-wrap.at-start .tabs-edge--left {
      opacity: 0;
      transform: translate(-6px, -50%);
    }

    .tabs-scroll-wrap.at-end .tabs-edge--right {
      opacity: 0;
      transform: translate(6px, -50%);
    }

    .tabs-scroll-wrap.is-hinting:not(.at-end) .tabs-edge--right {
      animation: tabsNudge 1.2s ease-in-out infinite;
    }

    .tabs-scroll-hint {
      position: absolute;
      right: 14px;
      bottom: 16px;
      z-index: 4;
      padding: 8px 10px;
      border-radius: 999px;
      font-size: 12px;
      font-weight: 700;
      letter-spacing: 0.2px;
      color: rgba(15, 23, 42, 0.72);
      background: rgba(255, 255, 255, 0.85);
      border: 1px solid rgba(15, 23, 42, 0.08);
      box-shadow: 0 16px 34px rgba(15, 23, 42, 0.12);
      backdrop-filter: blur(12px);
      pointer-events: none;
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
      gap: var(--gallery-gap);

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

    .gallery-waterfall {
      column-count: 2;
      column-gap: var(--gallery-gap);
    }

    @media (min-width: 920px) {
      .gallery-waterfall {
        column-count: 3;
      }
    }

    @media (min-width: 1280px) {
      .gallery-waterfall {
        column-count: 4;
        column-gap: var(--gallery-gap);
      }
    }

    .gallery-card {
      position: relative;
      border-radius: 30px;
      padding: 6px;
      background: linear-gradient(
        135deg,
        rgba(255, 97, 210, 0.22) 0%,
        rgba(254, 144, 144, 0.18) 42%,
        rgba(46, 211, 183, 0.18) 100%
      );
      box-shadow:
        0 20px 48px rgba(15, 23, 42, 0.12),
        0 10px 26px rgba(254, 110, 168, 0.08);
      display: block;
        break-inside: avoid;
        margin-bottom: var(--gallery-gap);
        transform: translateZ(0);
        transition: transform 0.25s ease, box-shadow 0.25s ease, filter 0.25s ease;

        .card-delete {
          position: absolute;
          top: 12px;
          right: 12px;
          width: 26px;
          height: 26px;
          border-radius: 10px;
          border: 1px solid rgba(15, 23, 42, 0.06);
          background: rgba(255, 255, 255, 0.82);
          color: rgba(15, 23, 42, 0.65);
          font-size: 16px;
          line-height: 1;
          display: inline-flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          backdrop-filter: blur(10px);
          opacity: 0.35;
          z-index: 2;
          box-shadow: 0 8px 18px rgba(15, 23, 42, 0.12);
          transition: opacity 0.2s ease, transform 0.15s ease, box-shadow 0.2s ease, background 0.2s ease;

          &:hover {
            opacity: 0.95;
            transform: translateY(-1px);
            box-shadow: 0 10px 24px rgba(15, 23, 42, 0.2);
            background: rgba(255, 255, 255, 0.92);
          }
        }

        &::after {
          content: '';
          position: absolute;
          inset: 0;
          border-radius: inherit;
        pointer-events: none;
        background: radial-gradient(
          120% 120% at 15% 0%,
          rgba(255, 255, 255, 0.55) 0%,
          rgba(255, 255, 255, 0) 48%
        );
        opacity: 0.8;
      }

      .gallery-media {
        width: 100%;
        border-radius: 24px;
        overflow: hidden;
        position: relative;
        background: rgba(255, 255, 255, 0.9);
        border: 1px solid rgba(255, 255, 255, 0.85);
        box-shadow:
          inset 0 1px 0 rgba(255, 255, 255, 0.75),
          0 10px 24px rgba(15, 23, 42, 0.08);

        &::before {
          content: '';
          position: absolute;
          inset: -20%;
          background: radial-gradient(
            60% 60% at 20% 12%,
            rgba(255, 97, 210, 0.16) 0%,
            rgba(255, 255, 255, 0) 60%
          );
          opacity: 0.6;
          pointer-events: none;
        }

        &::after {
          content: '';
          position: absolute;
          inset: 0;
          border-radius: inherit;
          pointer-events: none;
          border: 1px solid rgba(15, 23, 42, 0.06);
        }
      }

      .gallery-image {
        width: 100%;
        height: auto;
        object-fit: cover;
        display: block;
        transform: scale(1.01);
        transition: transform 0.35s ease, filter 0.35s ease;
      }

      .gallery-placeholder {
        width: 100%;
        aspect-ratio: 1 / 1;
        background: linear-gradient(135deg, rgba(245, 245, 247, 0.92) 0%, rgba(238, 242, 255, 0.92) 100%);
        color: rgba(15, 23, 42, 0.25);
        font-size: 34px;
        display: grid;
        place-items: center;
        position: relative;
      }

      .placeholder-shimmer {
        position: absolute;
        inset: 0;
        background: linear-gradient(
          110deg,
          rgba(255, 255, 255, 0) 0%,
          rgba(255, 255, 255, 0.65) 35%,
          rgba(255, 255, 255, 0) 70%
        );
        transform: translateX(-60%);
        animation: shimmer 1.4s ease-in-out infinite;
      }

      .placeholder-icon {
        position: relative;
        z-index: 1;
      }
    }

    @media (hover: hover) and (pointer: fine) {
      .gallery-card:hover {
        transform: translateY(-6px) scale(1.01);
        filter: saturate(1.05);
        box-shadow:
          0 26px 60px rgba(15, 23, 42, 0.16),
          0 14px 34px rgba(254, 110, 168, 0.14);
      }

      .gallery-card:hover .gallery-image {
        transform: scale(1.045);
        filter: contrast(1.03) saturate(1.06);
      }
    }

    .gallery-footer {
      padding: 14px 0 18px;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 10px;

      .gallery-footer__status {
        font-size: 12px;
        font-weight: 650;
        color: rgba(15, 23, 42, 0.55);
        padding: 8px 12px;
        border-radius: 999px;
        border: 1px solid rgba(15, 23, 42, 0.08);
        background: rgba(255, 255, 255, 0.75);
        backdrop-filter: blur(12px);
      }

      .gallery-sentinel {
        width: 1px;
        height: 1px;
      }
    }
  }
}

.fab-add {
  position: fixed;
  right: 22px;
  bottom: 28px;
  width: 56px;
  height: 56px;
  border-radius: 18px;
  border: none;
  background: linear-gradient(120deg, #ff61d2 0%, #fe9090 100%);
  color: #ffffff;
  font-size: 34px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 16px 32px rgba(254, 110, 168, 0.35);
  z-index: 9;
}

.add-style-form {
  :deep(.el-form-item__label) {
    padding-bottom: 8px;
    font-weight: 600;
    color: #202124;
  }

  .cropper-single {
    display: flex;
    flex-direction: column;
    gap: 14px;
    align-items: center;
  }

  .cropper-preview {
    width: min(360px, 100%);
    margin: 0 auto;
    padding: 0;
    border: 1px solid rgba(255, 255, 255, 0.65);
    border-radius: 22px;
    background: radial-gradient(120% 120% at 10% 0%, rgba(255, 97, 210, 0.16) 0%, rgba(254, 144, 144, 0.14) 35%, rgba(255, 255, 255, 0.9) 100%);
    cursor: pointer;
    overflow: hidden;
    aspect-ratio: 1 / 1;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 14px 38px rgba(15, 23, 42, 0.12);
    position: relative;

    &::before {
      content: '';
      position: absolute;
      inset: 0;
      background: linear-gradient(135deg, rgba(255, 97, 210, 0.32) 0%, rgba(254, 144, 144, 0.28) 40%, rgba(46, 211, 183, 0.18) 100%);
      opacity: 0;
      transition: opacity 0.25s ease;
      pointer-events: none;
    }

    &:hover::before {
      opacity: 1;
    }
  }

  .cropper-preview-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }

  .cropper-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 10px;
    padding: 22px 18px;
    text-align: center;
    color: #525866;
  }

  .cropper-orb {
    width: 56px;
    height: 56px;
    border-radius: 18px;
    background: linear-gradient(120deg, rgba(255, 97, 210, 0.95) 0%, rgba(254, 144, 144, 0.95) 55%, rgba(46, 211, 183, 0.9) 100%);
    box-shadow: 0 16px 36px rgba(254, 110, 168, 0.28);
    position: relative;

    &::after {
      content: '';
      position: absolute;
      inset: 10px;
      border-radius: 14px;
      background: rgba(255, 255, 255, 0.22);
      backdrop-filter: blur(10px);
    }
  }

  .cropper-placeholder__title {
    font-size: 16px;
    font-weight: 700;
    letter-spacing: 0.2px;
    color: #1f2937;
  }

  .cropper-placeholder__desc {
    font-size: 12px;
    color: rgba(82, 88, 102, 0.9);
    line-height: 1.4;
  }

  .cropper-cta {
    margin-top: 4px;
    padding: 10px 14px;
    border-radius: 999px;
    color: #ffffff;
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.3px;
    background: linear-gradient(120deg, #ff61d2 0%, #fe9090 55%, #2ed3b7 110%);
    box-shadow: 0 16px 32px rgba(254, 110, 168, 0.25);
  }

  .cropper-actions {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
    width: min(360px, 100%);
    margin: 0 auto;
    justify-content: center;

    :deep(.el-button) {
      border-radius: 14px;
      padding: 10px 14px;
      font-weight: 650;
      letter-spacing: 0.2px;
      flex: 1;
    }

    .btn-ghost {
      border: 1px solid rgba(31, 41, 55, 0.12);
      background: rgba(255, 255, 255, 0.75);
      color: #1f2937;
      backdrop-filter: blur(14px);
      transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;

      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 12px 28px rgba(15, 23, 42, 0.12);
        border-color: rgba(255, 97, 210, 0.25);
      }
    }

    .btn-danger-ghost {
      border: 1px solid rgba(239, 68, 68, 0.18);
      background: rgba(255, 255, 255, 0.75);
      color: #ef4444;
      backdrop-filter: blur(14px);

      &:hover {
        border-color: rgba(239, 68, 68, 0.28);
      }
    }
  }

  .add-style-hint {
    margin-top: -6px;
    font-size: 12px;
    color: #8b9097;
  }
}

@media (max-width: 480px) {
  .add-style-form {
    .cropper-actions {
      flex-direction: column;
    }
  }
}

@media (min-width: 1280px) {
  .nail-discover-public {
    --gallery-gap: 18px;
  }
}

/* 适配手机端弹窗：避免内容溢出屏幕 */
.add-style-dialog {
  :deep(.el-dialog) {
    max-width: calc(100% - 24px);
    border-radius: 26px;
    overflow: hidden;
    border: 1px solid rgba(255, 255, 255, 0.65);
    background: radial-gradient(130% 130% at 0% 0%, rgba(255, 97, 210, 0.18) 0%, rgba(254, 144, 144, 0.12) 38%, rgba(255, 255, 255, 0.92) 100%);
    box-shadow: 0 28px 70px rgba(15, 23, 42, 0.18);
  }

  :deep(.el-dialog__body) {
    padding: 18px 18px 0;
    max-height: calc(100vh - 190px);
    overflow: auto;
  }

  :deep(.el-dialog__footer) {
    position: sticky;
    bottom: 0;
    background: linear-gradient(180deg, rgba(255, 255, 255, 0.6) 0%, rgba(255, 255, 255, 0.92) 35%);
    backdrop-filter: blur(12px);
    padding: 14px 18px 18px;
  }

  :deep(.el-dialog__header) {
    padding: 0;
    margin-right: 0;
  }

  :deep(.el-dialog__headerbtn) {
    display: none;
  }
}

.add-style-header {
  padding: 18px 18px 16px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  background: linear-gradient(120deg, rgba(255, 97, 210, 0.24) 0%, rgba(254, 144, 144, 0.18) 40%, rgba(46, 211, 183, 0.14) 100%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.6);
}

.add-style-header__left {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.add-style-badge {
  height: 28px;
  padding: 0 10px;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.7px;
  color: #ffffff;
  background: linear-gradient(120deg, #ff61d2 0%, #fe9090 55%, #2ed3b7 110%);
  box-shadow: 0 14px 30px rgba(254, 110, 168, 0.24);
}

.add-style-titles {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}

.add-style-title {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 18px;
  font-weight: 800;
  letter-spacing: 0.2px;
  color: #0f172a;
}

.add-style-title__text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.add-style-title__close {
  flex: 0 0 auto;
  width: 34px;
  height: 34px;
  border-radius: 14px;
  border: 1px solid rgba(31, 41, 55, 0.12);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s ease, box-shadow 0.2s ease, border-color 0.2s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 14px 28px rgba(15, 23, 42, 0.12);
    border-color: rgba(255, 97, 210, 0.22);
  }
}

.add-style-title__close-x {
  font-size: 20px;
  line-height: 1;
  color: rgba(15, 23, 42, 0.72);
}

.add-style-subtitle {
  font-size: 12px;
  color: rgba(15, 23, 42, 0.62);
}

.add-style-subtitle__tag {
  font-weight: 800;
  color: rgba(15, 23, 42, 0.78);
}

.add-style-surface {
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.7);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(14px);
  padding: 14px;
  box-shadow: 0 18px 44px rgba(15, 23, 42, 0.1);
}

.add-style-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 2px 2px 12px;
}

.add-style-chip {
  padding: 8px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 650;
  color: rgba(15, 23, 42, 0.78);
  background: rgba(15, 23, 42, 0.05);
  border: 1px solid rgba(15, 23, 42, 0.06);
}

.premium-item {
  margin-bottom: 16px;
}

.cropper-item {
  :deep(.el-form-item__content) {
    display: flex;
    justify-content: center;
  }
}

.premium-input {
  :deep(.el-input__wrapper) {
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.72);
    backdrop-filter: blur(14px);
    box-shadow: 0 10px 26px rgba(15, 23, 42, 0.08);
    border: 1px solid rgba(15, 23, 42, 0.06);
    transition: box-shadow 0.2s ease, border-color 0.2s ease;
  }

  :deep(.el-input__wrapper.is-focus) {
    box-shadow: 0 16px 34px rgba(254, 110, 168, 0.16);
    border-color: rgba(255, 97, 210, 0.26);
  }
}

.add-style-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-footer-secondary {
  border-radius: 16px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.75);
  border: 1px solid rgba(15, 23, 42, 0.1);
  color: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(14px);
}

.add-style-dialog :deep(.el-button.btn-footer-primary) {
  border-radius: 16px;
  padding: 12px 18px;
  border: none;
  background: linear-gradient(120deg, #ff61d2 0%, #fe9090 55%, #2ed3b7 110%);
  box-shadow: 0 22px 44px rgba(254, 110, 168, 0.28);
  font-weight: 800;
  letter-spacing: 0.3px;
}

.add-style-dialog :deep(.el-button.btn-footer-primary .el-button__loading) {
  color: rgba(255, 255, 255, 0.9);
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

@keyframes shimmer {
  0% {
    transform: translateX(-60%);
  }
  100% {
    transform: translateX(60%);
  }
}

@keyframes tabsNudge {
  0% {
    transform: translate(0, -50%);
  }
  50% {
    transform: translate(4px, -50%);
  }
  100% {
    transform: translate(0, -50%);
  }
}
</style>
