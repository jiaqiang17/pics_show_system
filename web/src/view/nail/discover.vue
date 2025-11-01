<template>
  <div class="nail-discover">
    <!-- 顶部标题 -->
    <div class="header">
      <h1 class="title">发现</h1>
    </div>

    <!-- 标签筛选栏 -->
    <div class="tags-container">
      <div class="tags-scroll">
        <div
          v-for="tag in tags"
          :key="tag.id"
          class="tag-item"
          :class="{ active: activeTagId === tag.id }"
          @click="selectTag(tag.id)"
        >
          {{ tag.tagName }}
        </div>
      </div>
    </div>

    <!-- 美甲图片网格 -->
    <div class="nail-grid" v-loading="loading">
      <div
        v-for="nail in nailList"
        :key="nail.ID"
        class="nail-card"
        @click="viewDetail(nail)"
      >
        <div class="image-container">
          <el-image
            :src="nail.images && nail.images[0] ? nail.images[0].url : ''"
            fit="cover"
            class="nail-image"
            :preview-src-list="nail.images ? nail.images.map(img => img.url) : []"
          >
            <template #error>
              <div class="image-error">
                <el-icon><Picture /></el-icon>
              </div>
            </template>
          </el-image>
        </div>
        <div class="nail-info">
          <div class="nail-title">{{ nail.styleName }}</div>
          <div class="nail-stats">
            <span class="stat-item">
              <el-icon><View /></el-icon>
              {{ nail.viewCount || 0 }}
            </span>
            <span class="stat-item">
              <el-icon><Star /></el-icon>
              {{ nail.likeCount || 0 }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载更多 -->
    <div class="load-more" v-if="total > nailList.length">
      <el-button @click="loadMore" :loading="loading">加载更多</el-button>
    </div>

    <!-- 暂无数据 -->
    <el-empty v-if="!loading && nailList.length === 0" description="暂无美甲款式" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getNailStyleList } from '@/api/nail/nail_style'
import { getNailTagList } from '@/api/nail/nail_tag'
import { ElMessage } from 'element-plus'
import { Picture, View, Star } from '@element-plus/icons-vue'

// 响应式数据
const tags = ref([])
const activeTagId = ref(null)
const nailList = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 获取标签列表
const getTags = async () => {
  try {
    const res = await getNailTagList({
      page: 1,
      pageSize: 100,
      status: 'enabled'
    })
    if (res.code === 0) {
      // 添加"为你推荐"标签
      tags.value = [
        { id: null, tagName: '为你推荐' },
        ...res.data.list
      ]
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

    const res = await getNailStyleList(params)
    if (res.code === 0) {
      if (loadMore) {
        nailList.value = [...nailList.value, ...res.data.list]
      } else {
        nailList.value = res.data.list
      }
      total.value = res.data.total
    }
  } catch (error) {
    console.error('获取美甲列表失败:', error)
    ElMessage.error('获取美甲列表失败')
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
  page.value++
  getNailList(true)
}

// 查看详情
const viewDetail = (nail) => {
  // 这里可以跳转到详情页面或打开对话框
  console.log('查看详情:', nail)
}

// 页面加载时获取数据
onMounted(() => {
  getTags()
  getNailList()
})
</script>

<style scoped lang="scss">
.nail-discover {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 20px;

  .header {
    background: white;
    padding: 20px;
    text-align: center;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

    .title {
      font-size: 24px;
      font-weight: bold;
      margin: 0;
      color: #333;
    }
  }

  .tags-container {
    background: white;
    padding: 15px 0;
    margin-bottom: 10px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;

    .tags-scroll {
      display: flex;
      gap: 12px;
      padding: 0 20px;
      min-width: max-content;

      .tag-item {
        padding: 8px 20px;
        background: #f5f5f5;
        border-radius: 20px;
        font-size: 14px;
        color: #666;
        cursor: pointer;
        white-space: nowrap;
        transition: all 0.3s;

        &:hover {
          background: #e8e8e8;
        }

        &.active {
          background: linear-gradient(135deg, #FF6B9D 0%, #C06C84 100%);
          color: white;
          font-weight: 500;
        }
      }
    }

    &::-webkit-scrollbar {
      display: none;
    }
  }

  .nail-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 10px;
    padding: 0 10px;

    @media (min-width: 768px) {
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 15px;
      padding: 0 20px;
    }

    .nail-card {
      background: white;
      border-radius: 12px;
      overflow: hidden;
      cursor: pointer;
      transition: transform 0.3s, box-shadow 0.3s;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
      }

      .image-container {
        position: relative;
        width: 100%;
        padding-bottom: 133%; // 3:4 比例
        overflow: hidden;
        background: #f5f5f5;

        .nail-image {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;

          :deep(.el-image__inner) {
            width: 100%;
            height: 100%;
            object-fit: cover;
          }
        }

        .image-error {
          display: flex;
          align-items: center;
          justify-content: center;
          height: 100%;
          color: #ccc;
          font-size: 40px;
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

            .el-icon {
              font-size: 14px;
            }
          }
        }
      }
    }
  }

  .load-more {
    text-align: center;
    padding: 20px;
  }
}
</style>
