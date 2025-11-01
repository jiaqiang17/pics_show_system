<template>
  <div class="nail-style-page">
    <div class="gva-search-box">
      <el-form
        ref="searchFormRef"
        :inline="true"
        :model="filters"
        class="search-form"
        @keyup.enter="handleSearch"
      >
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="filters.createdAtRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            class="w-320"
            clearable
          />
        </el-form-item>
        <el-form-item label="款式名称">
          <el-input
            v-model="filters.styleName"
            placeholder="输入关键字搜索"
            clearable
          />
        </el-form-item>
        <el-form-item label="标签">
          <el-select
            v-model="filters.tagIds"
            multiple
            filterable
            clearable
            collapse-tags
            collapse-tags-tooltip
            placeholder="选择标签"
            class="w-240"
          >
            <el-option
              v-for="tag in tagOptions"
              :key="tag.value"
              :label="tag.label"
              :value="tag.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="filters.tagIds.length" label="匹配方式">
          <el-switch
            v-model="filters.matchAll"
            active-text="全部匹配"
            inactive-text="任意匹配"
            inline-prompt
          />
        </el-form-item>
        <el-form-item label="是否推荐">
          <el-select v-model="filters.isRecommended" clearable placeholder="全部">
            <el-option label="是" :value="true" />
            <el-option label="否" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="filters.status"
            clearable
            filterable
            placeholder="全部"
            class="w-160"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="RefreshRight" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" :icon="Plus" @click="openCreate">新增</el-button>
        <el-button
          type="success"
          :icon="Tickets"
          :disabled="!selectedIds.length"
          @click="openTagDialog"
        >
          批量调整标签
        </el-button>
        <el-button
          type="danger"
          plain
          :icon="Delete"
          :disabled="!selectedIds.length"
          @click="handleBatchDelete"
        >
          删除
        </el-button>
        <ExportTemplate template-id="nail_NailStyle" />
        <ExportExcel template-id="nail_NailStyle" filterDeleted />
        <ImportExcel template-id="nail_NailStyle" @on-success="fetchTableData" />
      </div>

      <el-table
        v-loading="tableLoading"
        :data="tableData"
        row-key="ID"
        style="width: 100%"
        tooltip-effect="dark"
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column
          label="创建时间"
          prop="CreatedAt"
          sortable="custom"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column
          label="款式名称"
          prop="styleName"
          min-width="140"
          show-overflow-tooltip
        />
        <el-table-column label="封面" width="120">
          <template #default="{ row }">
            <el-image
              v-if="row.images.length"
              :src="resolveImage(row.images[0])"
              :preview-src-list="resolvePreviewList(row.images)"
              fit="cover"
              class="thumb"
            />
            <span v-else class="image-placeholder">无</span>
          </template>
        </el-table-column>
        <el-table-column
          label="标签"
          min-width="200"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <el-space wrap size="small">
              <el-tag v-for="tagId in row.tagIds" :key="tagId" type="info">
                {{ resolveTagName(tagId) }}
              </el-tag>
            </el-space>
            <span v-if="!row.tagIds.length" class="text-disabled">未打标签</span>
          </template>
        </el-table-column>
        <el-table-column
          v-if="filters.tagIds.length"
          label="命中标签"
          prop="tagMatchCount"
          width="100"
          sortable="custom"
        />
        <el-table-column
          label="是否推荐"
          prop="isRecommended"
          width="110"
        >
          <template #default="{ row }">
            <el-tag :type="row.isRecommended ? 'success' : 'info'">
              {{ formatBoolean(row.isRecommended) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="浏览次数"
          prop="viewCount"
          sortable="custom"
          width="100"
        />
        <el-table-column
          label="点赞数"
          prop="likeCount"
          sortable="custom"
          width="100"
        />
        <el-table-column
          label="排序值"
          prop="sort"
          sortable="custom"
          width="90"
        />
        <el-table-column label="状态" prop="status" width="120">
          <template #default="{ row }">
            <el-tag>{{ filterDict(row.status, statusOptions) || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          fixed="right"
          :min-width="appStore.operateMinWith"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              :icon="View"
              @click="openDetail(row)"
            >
              查看
            </el-button>
            <el-button
              type="primary"
              link
              :icon="EditPen"
              @click="openEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              :icon="Delete"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="pagination.page"
          :page-size="pagination.pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="pagination.total"
          @current-change="handlePageChange"
          @size-change="handlePageSizeChange"
        />
      </div>
    </div>

    <el-dialog
      v-model="tagDialog.visible"
      width="460px"
      :close-on-click-modal="false"
      title="批量调整标签"
    >
      <el-form :model="tagDialog.form" label-width="96px">
        <el-form-item label="新增标签">
          <el-select
            v-model="tagDialog.form.addTagIds"
            multiple
            filterable
            placeholder="选择要添加的标签"
            clearable
          >
            <el-option
              v-for="tag in tagOptions"
              :key="`add-${tag.value}`"
              :label="tag.label"
              :value="tag.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="移除标签">
          <el-select
            v-model="tagDialog.form.removeTagIds"
            multiple
            filterable
            placeholder="选择要移除的标签"
            clearable
          >
            <el-option
              v-for="tag in tagOptions"
              :key="`remove-${tag.value}`"
              :label="tag.label"
              :value="tag.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="操作备注">
          <el-input
            v-model="tagDialog.form.operatorNote"
            type="textarea"
            maxlength="200"
            show-word-limit
            placeholder="可选，填写此次批量调整的说明"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tagDialog.visible = false">取 消</el-button>
        <el-button type="primary" :loading="tagDialog.loading" @click="submitTagAdjust">
          确 定
        </el-button>
      </template>
    </el-dialog>

    <el-drawer
      v-model="formDrawer.visible"
      :size="appStore.drawerSize"
      destroy-on-close
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="drawer-header">
          <span>{{ formDrawer.mode === 'create' ? '新增款式' : '编辑款式' }}</span>
          <div>
            <el-button @click="formDrawer.visible = false">取 消</el-button>
            <el-button type="primary" :loading="formDrawer.loading" @click="submitForm">
              保 存
            </el-button>
          </div>
        </div>
      </template>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-position="top"
      >
        <el-form-item label="款式名称" prop="styleName">
          <el-input v-model="formData.styleName" placeholder="请输入款式名称" />
        </el-form-item>
        <el-form-item label="展示图片" prop="images">
          <SelectImage v-model="formData.images" multiple file-type="image" />
        </el-form-item>
        <el-form-item label="款式介绍" prop="description">
          <RichEdit v-model="formData.description" />
        </el-form-item>
        <el-form-item label="是否推荐" prop="isRecommended">
          <el-switch
            v-model="formData.isRecommended"
            active-text="是"
            inactive-text="否"
          />
        </el-form-item>
        <el-form-item label="排序值" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="formData.status" placeholder="请选择状态">
            <el-option
              v-for="item in statusOptions"
              :key="`status-${item.value}`"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关联标签" prop="tagIds">
          <el-select
            v-model="formData.tagIds"
            multiple
            filterable
            clearable
            placeholder="选择关联标签"
          >
            <el-option
              v-for="tag in tagOptions"
              :key="`form-tag-${tag.value}`"
              :label="tag.label"
              :value="tag.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-model="detailDrawer.visible"
      title="款式详情"
      :size="appStore.drawerSize"
      destroy-on-close
    >
      <div v-if="detailDrawer.data" class="detail-wrapper">
        <section class="detail-section">
          <h3>基础信息</h3>
          <div class="detail-row">
            <span class="label">款式名称</span>
            <span>{{ detailDrawer.data.styleName || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="label">是否推荐</span>
            <span>{{ formatBoolean(detailDrawer.data.isRecommended) }}</span>
          </div>
          <div class="detail-row">
            <span class="label">状态</span>
            <span>{{ filterDict(detailDrawer.data.status, statusOptions) || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="label">排序值</span>
            <span>{{ detailDrawer.data.sort ?? '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="label">浏览 / 点赞</span>
            <span>{{ detailDrawer.data.viewCount ?? 0 }} / {{ detailDrawer.data.likeCount ?? 0 }}</span>
          </div>
          <div class="detail-row">
            <span class="label">关联标签</span>
            <el-space wrap size="small">
              <el-tag v-for="tagId in detailDrawer.data.tagIds" :key="`detail-${tagId}`">
                {{ resolveTagName(tagId) }}
              </el-tag>
            </el-space>
            <span v-if="!detailDrawer.data.tagIds.length">未打标签</span>
          </div>
        </section>

        <section class="detail-section">
          <h3>展示图片</h3>
          <el-space wrap>
            <el-image
              v-for="(img, index) in detailDrawer.data.images"
              :key="`detail-img-${index}`"
              :src="resolveImage(img)"
              :preview-src-list="resolvePreviewList(detailDrawer.data.images)"
              fit="cover"
              class="detail-image"
            />
            <span v-if="!detailDrawer.data.images.length">暂无图片</span>
          </el-space>
        </section>

        <section class="detail-section">
          <h3>款式介绍</h3>
          <div class="detail-rich" v-html="detailDrawer.data.description || '暂无介绍'" />
        </section>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useAppStore } from '@/pinia/modules/app'
import {
  createNailStyle,
  deleteNailStyle,
  deleteNailStyleByIds,
  findNailStyle,
  getNailStyleList,
  updateNailStyle,
  getNailStyleDataSource,
  batchUpdateNailStyleTags
} from '@/api/nail/nail_style'
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
import ImportExcel from '@/components/exportExcel/importExcel.vue'
import SelectImage from '@/components/selectImage/selectImage.vue'
import RichEdit from '@/components/richtext/rich-edit.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, EditPen, View, Tickets, RefreshRight, Search } from '@element-plus/icons-vue'
import { filterDict, formatBoolean, formatDate, getDictFunc, returnArrImg } from '@/utils/format'

const appStore = useAppStore()

const searchFormRef = ref()
const formRef = ref()

const filters = reactive({
  createdAtRange: [],
  styleName: '',
  tagIds: [],
  matchAll: false,
  isRecommended: null,
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const sortState = reactive({
  prop: '',
  order: ''
})

const tableData = ref([])
const tableLoading = ref(false)
const selectedRows = ref([])

const tagOptions = ref([])
const statusOptions = ref([])

const tagDialog = reactive({
  visible: false,
  loading: false,
  form: {
    addTagIds: [],
    removeTagIds: [],
    operatorNote: ''
  }
})

const formDrawer = reactive({
  visible: false,
  mode: 'create',
  loading: false
})

const detailDrawer = reactive({
  visible: false,
  data: null
})

const formData = reactive({
  ID: undefined,
  styleName: '',
  images: [],
  description: '',
  isRecommended: false,
  sort: 0,
  status: '',
  tagIds: []
})

const formRules = {
  styleName: [{ required: true, message: '请输入款式名称', trigger: 'blur' }],
  images: [{ required: true, message: '请上传展示图片', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const selectedIds = computed(() => selectedRows.value.map((item) => item.ID))

const tagMap = computed(() => {
  const map = new Map()
  tagOptions.value.forEach((item) => {
    map.set(item.value, item.label)
  })
  return map
})

const resolveTagName = (id) => tagMap.value.get(id) || `#${id}`

const resolveImage = (image) => {
  const list = returnArrImg(image)
  return list.length ? list[0] : ''
}

const resolvePreviewList = (images) => {
  if (!images?.length) return []
  return images.map((img) => resolveImage(img))
}

const transformRecord = (item) => {
  const normalizeJSON = (value) => {
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
    images: normalizeJSON(item.images),
    tagIds: normalizeJSON(item.tagIds)
  }
}

const buildQueryParams = () => {
  const params = {
    page: pagination.page,
    pageSize: pagination.pageSize
  }

  if (filters.createdAtRange?.length === 2) {
    params.createdAtRange = filters.createdAtRange
  }
  if (filters.styleName) {
    params.styleName = filters.styleName
  }
  if (filters.tagIds.length) {
    params.tagIds = filters.tagIds
    params.matchAll = filters.matchAll
  }
  if (filters.isRecommended !== null && filters.isRecommended !== undefined) {
    params.isRecommended = filters.isRecommended
  }
  if (filters.status) {
    params.status = filters.status
  }
  if (sortState.prop) {
    params.sort = sortState.prop
    params.order = sortState.order
  }
  return params
}

const fetchTableData = async () => {
  tableLoading.value = true
  try {
    const res = await getNailStyleList(buildQueryParams())
    if (res.code === 0) {
      const { list = [], total = 0 } = res.data || {}
      tableData.value = list.map(transformRecord)
      pagination.total = total
    }
  } finally {
    tableLoading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchTableData()
}

const handleReset = () => {
  filters.createdAtRange = []
  filters.styleName = ''
  filters.tagIds = []
  filters.matchAll = false
  filters.isRecommended = null
  filters.status = ''
  sortState.prop = ''
  sortState.order = ''
  pagination.page = 1
  fetchTableData()
}

const handleSelectionChange = (rows) => {
  selectedRows.value = rows
}

const handleSortChange = ({ prop, order }) => {
  sortState.prop = order ? prop : ''
  sortState.order = order
  fetchTableData()
}

const handlePageChange = (page) => {
  pagination.page = page
  fetchTableData()
}

const handlePageSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchTableData()
}

const openCreate = () => {
  formDrawer.mode = 'create'
  Object.assign(formData, {
    ID: undefined,
    styleName: '',
    images: [],
    description: '',
    isRecommended: false,
    sort: 0,
    status: '',
    tagIds: []
  })
  formDrawer.visible = true
}

const openEdit = async (row) => {
  formDrawer.mode = 'update'
  const res = await findNailStyle({ ID: row.ID })
  if (res.code === 0) {
    const data = transformRecord(res.data)
    Object.assign(formData, data)
    formDrawer.visible = true
  }
}

const submitForm = () => {
  formRef.value?.validate(async (valid) => {
    if (!valid) return
    formDrawer.loading = true
    try {
      const payload = { ...formData }
      let res
      if (formDrawer.mode === 'create') {
        res = await createNailStyle(payload)
      } else {
        res = await updateNailStyle(payload)
      }
      if (res.code === 0) {
        ElMessage.success('保存成功')
        formDrawer.visible = false
        fetchTableData()
      }
    } finally {
      formDrawer.loading = false
    }
  })
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm(`确定删除「${row.styleName}」吗？`, '提示', {
    type: 'warning'
  })
  const res = await deleteNailStyle({ ID: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && pagination.page > 1) {
      pagination.page -= 1
    }
    fetchTableData()
  }
}

const handleBatchDelete = async () => {
  await ElMessageBox.confirm('确定删除选中的款式吗？', '批量删除', {
    type: 'warning'
  })
  const res = await deleteNailStyleByIds({ IDs: selectedIds.value })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === selectedIds.value.length && pagination.page > 1) {
      pagination.page -= 1
    }
    fetchTableData()
  }
}

const openTagDialog = () => {
  tagDialog.form.addTagIds = []
  tagDialog.form.removeTagIds = []
  tagDialog.form.operatorNote = ''
  tagDialog.visible = true
}

const submitTagAdjust = async () => {
  if (!tagDialog.form.addTagIds.length && !tagDialog.form.removeTagIds.length) {
    ElMessage.warning('请至少选择需要新增或移除的标签')
    return
  }
  tagDialog.loading = true
  try {
    const payload = {
      styleIds: selectedIds.value,
      addTagIds: tagDialog.form.addTagIds,
      removeTagIds: tagDialog.form.removeTagIds,
      operatorNote: tagDialog.form.operatorNote || undefined
    }
    const res = await batchUpdateNailStyleTags(payload)
    if (res.code === 0) {
      const result = res.data || {}
      let message = `成功更新 ${result.updated || 0} 条记录`
      if (result.skipped?.length) {
        message += `，跳过 ${result.skipped.length} 条`
      }
      if (result.invalidTagIds?.length) {
        message += `。无效标签：${result.invalidTagIds.join(', ')}`
      }
      ElMessage.success(message)
      tagDialog.visible = false
      fetchTableData()
    }
  } finally {
    tagDialog.loading = false
  }
}

const openDetail = async (row) => {
  const res = await findNailStyle({ ID: row.ID })
  if (res.code === 0) {
    detailDrawer.data = transformRecord(res.data)
    detailDrawer.visible = true
  }
}

const initOptions = async () => {
  const [dictRes, dataSourceRes] = await Promise.all([
    getDictFunc('nail_style_status'),
    getNailStyleDataSource()
  ])
  statusOptions.value = dictRes || []
  const tagList = dataSourceRes?.data?.tagIds || []
  tagOptions.value = tagList.map((item) => ({
    label: item.label,
    value: item.value
  }))
}

onMounted(async () => {
  await initOptions()
  fetchTableData()
})
</script>

<style scoped>
.nail-style-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-form .el-form-item {
  margin-right: 16px;
  margin-bottom: 12px;
}

.w-320 {
  width: 320px;
}

.w-240 {
  width: 240px;
}

.w-160 {
  width: 160px;
}

.thumb {
  width: 80px;
  height: 80px;
  border-radius: 12px;
}

.image-placeholder {
  color: var(--el-color-info);
  font-size: 12px;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.detail-wrapper {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-section > h3 {
  margin: 0 0 12px;
  font-weight: 600;
  font-size: 16px;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.detail-row .label {
  width: 100px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.detail-image {
  width: 120px;
  height: 120px;
  border-radius: 12px;
  object-fit: cover;
}

.detail-rich {
  padding: 12px;
  border-radius: 12px;
  background: var(--el-fill-color-lighter);
  line-height: 1.7;
  color: var(--el-text-color-primary);
}
</style>
