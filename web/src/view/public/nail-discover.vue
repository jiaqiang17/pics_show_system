<template>
  <div class="nail-discover-public">
    <!-- 顶部标题 -->
    <div class="header">
      <h1 class="title">✨ 发现美甲灵感</h1>
      <p class="subtitle">探索最新最美的美甲款式</p>
    </div>

    <!-- 标签筛选栏 -->
    <div class="tags-container">
      <div class="tags-scroll">
        <div
          class="tag-item"
          :class="{ active: activeTagId === null }"
          @click="selectTag(null)"
        >
          <span class="tag-icon">✨</span>
          为你推荐
        </div>
        <div
          v-for="tag in tags"
          :key="tag.ID"
          class="tag-item"
          :class="{ active: activeTagId === tag.ID }"
          :style="{ '--tag-color': tag.tagColor || '#FF6B9D' }"
          @click="selectTag(tag.ID)"
        >
          <span v-if="tag.tagIcon" class="tag-icon">{{ tag.tagIcon }}</span>
          {{ tag.tagName }}
        </div>
      </div>
    </div>

    <!-- 美甲图片网格 -->
    <div class="content-wrapper">
      <div v-if="loading && page === 1" class="loading-wrapper">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="!loading && nailList.length === 0" class="empty-wrapper">
        <div class="empty-icon">💅</div>
        <p class="empty-text">暂无美甲款式</p>
        <p class="empty-hint">敬请期待更多精彩内容</p>
      </div>

      <div v-else class="nail-grid">
        <div
          v-for="nail in nailList"
          :key="nail.ID"
          class="nail-card"
          @click="viewDetail(nail)"
        >
          <div class="image-container">
            <img
              v-if="nail.images && nail.images.length > 0"
              :src="getImageUrl(nail.images[0])"
              :alt="nail.styleName"
              class="nail-image"
              @error="handleImageError"
            />
            <div v-else class="image-placeholder">
              <span class="placeholder-icon">🖼️</span>
            </div>

            <!-- 多图角标 -->
            <div v-if="nail.images && nail.images.length > 1" class="image-count">
              <span class="count-icon">📷</span>
              {{ nail.images.length }}
            </div>
          </div>

          <div class="nail-info">
            <div class="nail-title">{{ nail.styleName }}</div>
            <div class="nail-stats">
              <span class="stat-item">
                <span class="stat-icon">👁️</span>
                {{ formatNumber(nail.viewCount || 0) }}
              </span>
              <span class="stat-item">
                <span class="stat-icon">❤️</span>
                {{ formatNumber(nail.likeCount || 0) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="nailList.length > 0 && total > nailList.length" class="load-more">
        <button class="load-more-btn" @click="loadMore" :disabled="loading">
          <span v-if="loading">加载中...</span>
          <span v-else>加载更多</span>
        </button>
      </div>

      <!-- 已加载全部 -->
      <div v-if="nailList.length > 0 && total <= nailList.length" class="load-end">
        <span class="end-line"></span>
        <span class="end-text">已加载全部</span>
        <span class="end-line"></span>
      </div>
    </div>

    <!-- 详情弹窗 -->
    <div v-if="detailVisible" class="detail-modal" @click="closeDetail">
      <div class="modal-content" @click.stop>
        <button class="close-btn" @click="closeDetail">✕</button>

        <div class="detail-images">
          <div class="image-slider">
            <img
              v-if="currentDetail.images && currentDetail.images[currentImageIndex]"
              :src="getImageUrl(currentDetail.images[currentImageIndex])"
              :alt="currentDetail.styleName"
              class="detail-image"
            />

            <!-- 左右切换按钮 -->
            <button
              v-if="currentDetail.images && currentDetail.images.length > 1"
              class="nav-btn prev-btn"
              @click="prevImage"
            >
              ‹
            </button>
            <button
              v-if="currentDetail.images && currentDetail.images.length > 1"
              class="nav-btn next-btn"
              @click="nextImage"
            >
              ›
            </button>

            <!-- 图片指示器 -->
            <div
              v-if="currentDetail.images && currentDetail.images.length > 1"
              class="image-indicators"
            >
              <span
                v-for="(img, idx) in currentDetail.images"
                :key="idx"
                class="indicator"
                :class="{ active: idx === currentImageIndex }"
                @click="currentImageIndex = idx"
              ></span>
            </div>
          </div>
        </div>

        <div class="detail-info">
          <h2 class="detail-title">{{ currentDetail.styleName }}</h2>
          <div class="detail-stats">
            <span class="detail-stat">
              <span class="stat-icon">👁️</span>
              {{ formatNumber(currentDetail.viewCount || 0) }} 浏览
            </span>
            <span class="detail-stat">
              <span class="stat-icon">❤️</span>
              {{ formatNumber(currentDetail.likeCount || 0) }} 喜欢
            </span>
          </div>
          <div v-if="currentDetail.description" class="detail-description" v-html="currentDetail.description"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

// 创建一个不需要认证的 axios 实例
const publicApi = axios.create({
  baseURL: import.meta.env.VITE_BASE_API || '/api',
  timeout: 30000
})

// 响应式数据
const tags = ref([])
const activeTagId = ref(null)
const nailList = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 详情弹窗
const detailVisible = ref(false)
const currentDetail = ref({})
const currentImageIndex = ref(0)

// 获取标签列表
const getTags = async () => {
  try {
    const res = await publicApi.get('/nailTag/getNailTagList', {
      params: {
        page: 1,
        pageSize: 100,
        status: 'enabled'
      }
    })
    if (res.data.code === 0) {
      tags.value = res.data.data.list || []
    }
  } catch (error) {
    console.error('获取标签失败:', error)
  }
}

// 获取美甲列表
const getNailList = async (loadMore = false) => {
  try {
    loading.value = true
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      status: 'enabled'
    }

    // 如果有选中的标签，添加标签筛选
    if (activeTagId.value) {
      params.tagId = activeTagId.value
    } else {
      // 如果是"为你推荐"，显示被推荐的款式
      params.isRecommended = true
    }

    const res = await publicApi.get('/nailStyle/getNailStyleList', { params })
    if (res.data.code === 0) {
      const newList = res.data.data.list || []
      if (loadMore) {
        nailList.value = [...nailList.value, ...newList]
      } else {
        nailList.value = newList
      }
      total.value = res.data.data.total || 0
    }
  } catch (error) {
    console.error('获取美甲列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 选择标签
const selectTag = (tagId) => {
  activeTagId.value = tagId
  page.value = 1
  getNailList()
}

// 加载更多
const loadMore = () => {
  if (loading.value) return
  page.value++
  getNailList(true)
}

// 查看详情
const viewDetail = (nail) => {
  currentDetail.value = nail
  currentImageIndex.value = 0
  detailVisible.value = true
  document.body.style.overflow = 'hidden'
}

// 关闭详情
const closeDetail = () => {
  detailVisible.value = false
  document.body.style.overflow = ''
}

// 上一张图片
const prevImage = () => {
  if (currentImageIndex.value > 0) {
    currentImageIndex.value--
  } else {
    currentImageIndex.value = currentDetail.value.images.length - 1
  }
}

// 下一张图片
const nextImage = () => {
  if (currentImageIndex.value < currentDetail.value.images.length - 1) {
    currentImageIndex.value++
  } else {
    currentImageIndex.value = 0
  }
}

// 获取图片URL
const getImageUrl = (image) => {
  if (typeof image === 'string') {
    return image
  }
  return image?.url || ''
}

// 图片加载错误处理
const handleImageError = (e) => {
  e.target.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200"%3E%3Crect fill="%23f5f5f5" width="200" height="200"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle" font-size="60"%3E🖼️%3C/text%3E%3C/svg%3E'
}

// 格式化数字
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num
}

// 页面加载时获取数据
onMounted(() => {
  getTags()
  getNailList()
})

// 键盘事件（图片切换）
const handleKeyDown = (e) => {
  if (!detailVisible.value) return
  if (e.key === 'ArrowLeft') {
    prevImage()
  } else if (e.key === 'ArrowRight') {
    nextImage()
  } else if (e.key === 'Escape') {
    closeDetail()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

// 清理事件监听
import { onUnmounted } from 'vue'
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  document.body.style.overflow = ''
})
</script>

<style scoped lang="scss">
.nail-discover-public {
  min-height: 100vh;
  background: linear-gradient(180deg, #fff5f7 0%, #f5f5f5 100%);
  padding-bottom: 40px;

  .header {
    background: linear-gradient(135deg, #FF6B9D 0%, #C06C84 100%);
    padding: 40px 20px 30px;
    text-align: center;
    color: white;
    box-shadow: 0 4px 12px rgba(255, 107, 157, 0.2);

    .title {
      font-size: 32px;
      font-weight: bold;
      margin: 0 0 10px 0;
      text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .subtitle {
      font-size: 16px;
      opacity: 0.9;
      margin: 0;
    }
  }

  .tags-container {
    background: white;
    padding: 20px 0;
    margin-bottom: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;

    &::-webkit-scrollbar {
      display: none;
    }

    .tags-scroll {
      display: flex;
      gap: 12px;
      padding: 0 20px;
      min-width: max-content;

      .tag-item {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 10px 24px;
        background: #f5f5f5;
        border-radius: 24px;
        font-size: 15px;
        color: #666;
        cursor: pointer;
        white-space: nowrap;
        transition: all 0.3s;
        user-select: none;

        .tag-icon {
          font-size: 16px;
        }

        &:hover {
          background: #e8e8e8;
          transform: translateY(-2px);
        }

        &.active {
          background: var(--tag-color, linear-gradient(135deg, #FF6B9D 0%, #C06C84 100%));
          color: white;
          font-weight: 500;
          box-shadow: 0 4px 12px rgba(255, 107, 157, 0.3);
        }
      }
    }
  }

  .content-wrapper {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 10px;

    @media (min-width: 768px) {
      padding: 0 20px;
    }
  }

  .loading-wrapper {
    text-align: center;
    padding: 60px 20px;
    color: #999;

    .loading-spinner {
      width: 50px;
      height: 50px;
      margin: 0 auto 20px;
      border: 4px solid #f3f3f3;
      border-top: 4px solid #FF6B9D;
      border-radius: 50%;
      animation: spin 1s linear infinite;
    }

    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
  }

  .empty-wrapper {
    text-align: center;
    padding: 80px 20px;

    .empty-icon {
      font-size: 80px;
      margin-bottom: 20px;
      opacity: 0.5;
    }

    .empty-text {
      font-size: 18px;
      color: #666;
      margin: 0 0 10px 0;
    }

    .empty-hint {
      font-size: 14px;
      color: #999;
      margin: 0;
    }
  }

  .nail-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;

    @media (min-width: 640px) {
      grid-template-columns: repeat(3, 1fr);
    }

    @media (min-width: 768px) {
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;
    }

    @media (min-width: 1024px) {
      grid-template-columns: repeat(5, 1fr);
    }

    @media (min-width: 1280px) {
      grid-template-columns: repeat(6, 1fr);
    }

    .nail-card {
      background: white;
      border-radius: 24px;
      overflow: hidden;
      cursor: pointer;
      transition: all 0.3s;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);

      &:hover {
        transform: translateY(-6px);
        box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
      }

      .image-container {
        position: relative;
        width: 100%;
        padding-bottom: 133%;
        overflow: hidden;
        background: #f5f5f5;

        .nail-image {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          object-fit: cover;
        }

        .image-placeholder {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: center;
          background: #f5f5f5;

          .placeholder-icon {
            font-size: 48px;
            opacity: 0.3;
          }
        }

        .image-count {
          position: absolute;
          top: 8px;
          right: 8px;
          background: rgba(0, 0, 0, 0.6);
          color: white;
          padding: 4px 10px;
          border-radius: 12px;
          font-size: 12px;
          display: flex;
          align-items: center;
          gap: 4px;

          .count-icon {
            font-size: 14px;
          }
        }
      }

      .nail-info {
        padding: 12px;

        .nail-title {
          font-size: 14px;
          color: #333;
          margin-bottom: 8px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          font-weight: 500;
        }

        .nail-stats {
          display: flex;
          gap: 12px;
          font-size: 12px;
          color: #999;

          .stat-item {
            display: flex;
            align-items: center;
            gap: 4px;

            .stat-icon {
              font-size: 14px;
            }
          }
        }
      }
    }
  }

  .load-more {
    text-align: center;
    padding: 30px 20px;

    .load-more-btn {
      padding: 12px 40px;
      background: linear-gradient(135deg, #FF6B9D 0%, #C06C84 100%);
      color: white;
      border: none;
      border-radius: 24px;
      font-size: 15px;
      cursor: pointer;
      transition: all 0.3s;
      box-shadow: 0 4px 12px rgba(255, 107, 157, 0.3);

      &:hover:not(:disabled) {
        transform: translateY(-2px);
        box-shadow: 0 6px 16px rgba(255, 107, 157, 0.4);
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }
    }
  }

  .load-end {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 30px 20px;
    color: #999;
    font-size: 14px;

    .end-line {
      flex: 1;
      max-width: 60px;
      height: 1px;
      background: #e0e0e0;
    }
  }

  // 详情弹窗
  .detail-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.85);
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    animation: fadeIn 0.3s;

    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }

    .modal-content {
      position: relative;
      background: white;
      border-radius: 20px;
      max-width: 900px;
      width: 100%;
      max-height: 90vh;
      overflow-y: auto;
      animation: slideUp 0.3s;

      @keyframes slideUp {
        from {
          opacity: 0;
          transform: translateY(50px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }

      .close-btn {
        position: absolute;
        top: 20px;
        right: 20px;
        width: 40px;
        height: 40px;
        background: rgba(0, 0, 0, 0.5);
        color: white;
        border: none;
        border-radius: 50%;
        font-size: 24px;
        cursor: pointer;
        z-index: 10;
        transition: all 0.3s;

        &:hover {
          background: rgba(0, 0, 0, 0.7);
          transform: rotate(90deg);
        }
      }

      .detail-images {
        .image-slider {
          position: relative;
          width: 100%;
          padding-bottom: 75%;
          background: #f5f5f5;
          border-radius: 20px 20px 0 0;
          overflow: hidden;

          .detail-image {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            object-fit: contain;
          }

          .nav-btn {
            position: absolute;
            top: 50%;
            transform: translateY(-50%);
            width: 50px;
            height: 50px;
            background: rgba(0, 0, 0, 0.5);
            color: white;
            border: none;
            font-size: 32px;
            cursor: pointer;
            transition: all 0.3s;
            display: flex;
            align-items: center;
            justify-content: center;

            &:hover {
              background: rgba(0, 0, 0, 0.7);
            }

            &.prev-btn {
              left: 20px;
              border-radius: 0 8px 8px 0;
            }

            &.next-btn {
              right: 20px;
              border-radius: 8px 0 0 8px;
            }
          }

          .image-indicators {
            position: absolute;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            display: flex;
            gap: 8px;

            .indicator {
              width: 8px;
              height: 8px;
              background: rgba(255, 255, 255, 0.5);
              border-radius: 50%;
              cursor: pointer;
              transition: all 0.3s;

              &.active {
                background: white;
                width: 24px;
                border-radius: 4px;
              }
            }
          }
        }
      }

      .detail-info {
        padding: 30px;

        .detail-title {
          font-size: 24px;
          font-weight: bold;
          color: #333;
          margin: 0 0 16px 0;
        }

        .detail-stats {
          display: flex;
          gap: 24px;
          margin-bottom: 20px;
          font-size: 14px;
          color: #666;

          .detail-stat {
            display: flex;
            align-items: center;
            gap: 6px;

            .stat-icon {
              font-size: 18px;
            }
          }
        }

        .detail-description {
          font-size: 15px;
          line-height: 1.8;
          color: #666;

          :deep(p) {
            margin: 0 0 10px 0;
          }

          :deep(img) {
            max-width: 100%;
            height: auto;
            border-radius: 8px;
            margin: 10px 0;
          }
        }
      }
    }
  }
}
</style>
