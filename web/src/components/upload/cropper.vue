<template>
  <el-upload
      ref="uploadRef"
      :action="actionUrl"
      accept="image/*"
      :show-file-list="false"
      :auto-upload="false"
      :data="{'classId': props.classId}"
      :on-success="handleImageSuccess"
      :on-change="handleFileChange"
      :headers="{'x-token': token}"
  >
    <el-button
      v-if="!props.hideTrigger"
      :type="props.triggerType"
      icon="crop"
      :class="props.triggerClass"
    >
      {{ props.triggerText }}
    </el-button>
    <span v-else class="cropper-hidden-trigger"></span>
  </el-upload>

  <el-dialog
    v-model="dialogVisible"
    :title="props.dialogTitle"
    :width="dialogWidth"
    :fullscreen="isMobile"
    append-to-body
    @close="dialogVisible = false"
    :close-on-click-modal="false"
    draggable
  >
    <div :class="dialogBodyClass">
      <!-- 左侧编辑区 -->
      <div class="flex flex-col flex-1">
        <div class="bg-[#f8f8f8] rounded-lg overflow-hidden" :style="cropAreaStyle">
          <VueCropper
              ref="cropperRef"
              :img="imgSrc"
              outputType="jpeg"
              :autoCrop="true"
              :autoCropWidth="cropWidth"
              :autoCropHeight="cropHeight"
              :fixedBox="false"
              :fixed="fixedRatio"
              :fixedNumber="fixedNumber"
              :centerBox="true"
              :canMoveBox="true"
              :full="false"
              :maxImgSize="1200"
              :original="true"
              @realTime="handleRealTime"
              @imgLoad="handleImgLoad"
          ></VueCropper>
        </div>

        <!-- 工具栏 -->
        <div class="mt-[16px] flex flex-wrap items-center gap-3 p-[10px] bg-white rounded-lg shadow-[0_2px_12px_rgba(0,0,0,0.1)]">
          <el-button-group>
            <el-tooltip content="向左旋转">
              <el-button @click="rotate(-90)" :icon="RefreshLeft" />
            </el-tooltip>
            <el-tooltip content="向右旋转">
              <el-button @click="rotate(90)" :icon="RefreshRight" />
            </el-tooltip>
            <el-button :icon="Plus" @click="changeScale(1)"></el-button>
            <el-button :icon="Minus" @click="changeScale(-1)"></el-button>
          </el-button-group>


          <el-select v-model="currentRatio" placeholder="选择比例" class="w-32" @change="onCurrentRatio">
            <el-option v-for="(item, index) in ratioOptions" :key="index" :label="item.label" :value="index" />
          </el-select>
        </div>
      </div>

      <!-- 右侧预览区 -->
      <div :class="previewPanelClass">
        <div v-if="props.showPreview" class="bg-white p-5 rounded-lg shadow-[0_2px_12px_rgba(0,0,0,0.1)]">
          <div class="mb-[15px] text-gray-600">裁剪预览</div>
          <div
            class="bg-white p-5 rounded-lg shadow-[0_2px_12px_rgba(0,0,0,0.1)]"
            :style="{'width': previews.w + 'px', 'height': previews.h + 'px'}"
          >
            <div class="w-full h-full relative overflow-hidden">
              <img
                :src="previews.url"
                :style="previews.img"
                alt=""
                class="max-w-none absolute transition-all duration-300 ease-in-out image-render-pixelated origin-[0_0]"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleUpload" :loading="uploading"> {{ uploading ? '上传中...' : '上 传' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, getCurrentInstance, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { RefreshLeft, RefreshRight, Plus, Minus } from '@element-plus/icons-vue'
import 'vue-cropper/dist/index.css'
import { VueCropper } from 'vue-cropper'
import { getBaseUrl } from '@/utils/format'
import { useUserStore } from '@/pinia'

defineOptions({
  name: 'CropperImage'
})

const emit = defineEmits(['on-success'])

const props = defineProps({
  classId: {
    type: Number,
    default: 0
  },
  uploadPath: {
    type: String,
    default: '/fileUploadAndDownload/upload'
  },
  defaultRatioIndex: {
    type: Number,
    default: 4
  },
  defaultCropSize: {
    type: Number,
    default: 300
  },
  defaultCropSizeMobile: {
    type: Number,
    default: 300
  },
  initialScaleSteps: {
    type: Number,
    default: 0
  },
  hideTrigger: {
    type: Boolean,
    default: false
  },
  triggerText: {
    type: String,
    default: '裁剪上传'
  },
  triggerType: {
    type: String,
    default: 'primary'
  },
  triggerClass: {
    type: String,
    default: ''
  },
  dialogTitle: {
    type: String,
    default: '图片裁剪'
  },
  showPreview: {
    type: Boolean,
    default: true
  }
})

const uploadRef = ref(null)
// 响应式数据
const dialogVisible = ref(false)
const imgSrc = ref('')
const cropperRef = ref(null)
const { proxy } = getCurrentInstance()
const previews = ref({})
const uploading = ref(false)
const autoScaled = ref(false)

const userStore = useUserStore()
const token = userStore.token

const actionUrl = computed(() => {
  const base = String(getBaseUrl() || '')
  const baseNormalized = base.endsWith('/') ? base.slice(0, -1) : base
  const path = String(props.uploadPath || '/fileUploadAndDownload/upload')
  const pathNormalized = path.startsWith('/') ? path : `/${path}`
  return `${baseNormalized}${pathNormalized}`
})

const isMobile = ref(false)
let mediaQuery = null
const updateIsMobile = () => {
  isMobile.value = !!mediaQuery?.matches
}

onMounted(() => {
  if (typeof window === 'undefined' || !window.matchMedia) return
  mediaQuery = window.matchMedia('(max-width: 768px)')
  updateIsMobile()
  if (mediaQuery.addEventListener) {
    mediaQuery.addEventListener('change', updateIsMobile)
  } else if (mediaQuery.addListener) {
    mediaQuery.addListener(updateIsMobile)
  }
})

onBeforeUnmount(() => {
  if (!mediaQuery) return
  if (mediaQuery.removeEventListener) {
    mediaQuery.removeEventListener('change', updateIsMobile)
  } else if (mediaQuery.removeListener) {
    mediaQuery.removeListener(updateIsMobile)
  }
})

const dialogWidth = computed(() => (isMobile.value ? '100%' : '1200px'))
const dialogBodyClass = computed(() =>
  isMobile.value ? 'flex flex-col gap-4' : 'flex gap-[30px] h-[600px]'
)
const cropAreaStyle = computed(() => (isMobile.value ? { height: '55vh' } : { flex: 1 }))
const previewPanelClass = computed(() => {
  if (!props.showPreview) return 'hidden'
  return isMobile.value ? 'w-full' : 'w-[340px]'
})

const open = () => {
  const input = uploadRef.value?.$el?.querySelector?.('input[type="file"]')
  input?.click?.()
}

defineExpose({
  open
})

// 缩放控制
const changeScale = (value) => {
  proxy.$refs.cropperRef.changeScale(value)
}

// 比例预设
const ratioOptions = ref([
  { label: '1:1', value: [1, 1] },
  { label: '16:9', value: [16, 9] },
  { label: '9:16', value: [9, 16] },
  { label: '4:3', value: [4, 3] },
  { label: '自由比例', value: [] }
])

const fixedNumber = ref([1, 1])
const cropWidth = ref(props.defaultCropSize)
const cropHeight = ref(props.defaultCropSize)

const fixedRatio = ref(false)
const currentRatio = ref(props.defaultRatioIndex)
const onCurrentRatio = () => {
  const baseSize = isMobile.value ? props.defaultCropSizeMobile : props.defaultCropSize
  fixedNumber.value = ratioOptions.value[currentRatio.value].value
  switch (currentRatio.value) {
    case 0:
      cropWidth.value = baseSize
      cropHeight.value = baseSize
      fixedRatio.value = true
      break
    case 1:
      cropWidth.value = baseSize
      cropHeight.value = baseSize * 9 / 16
      fixedRatio.value = true
      break
    case 2:
      cropWidth.value = baseSize * 9 / 16
      cropHeight.value = baseSize
      fixedRatio.value = true
      break
    case 3:
      cropWidth.value = baseSize
      cropHeight.value = baseSize * 3 / 4
      fixedRatio.value = true
      break
    default:
      cropWidth.value = baseSize
      cropHeight.value = baseSize
      fixedRatio.value = false
  }
}

const applyInitialScale = async () => {
  const steps = Number(props.initialScaleSteps || 0)
  if (!steps) return
  await nextTick()

  const changeScaleBy = (delta) => proxy?.$refs?.cropperRef?.changeScale?.(delta)

  if (steps > 0) {
    const absSteps = Math.abs(steps)
    for (let i = 0; i < absSteps; i++) {
      changeScaleBy(1)
    }
    return
  }

  // 缩小时：尽量缩到最小（达到约束边界后 changeScale 会返回 false）
  const absSteps = Math.abs(steps)
  const maxAttempts = Math.max(absSteps, 60)
  for (let i = 0; i < maxAttempts; i++) {
    const ok = changeScaleBy(-1)
    if (ok === false) break
  }
}

const handleImgLoad = async (msg) => {
  if (msg !== 'success') return
  if (autoScaled.value) return
  autoScaled.value = true
  await applyInitialScale()
}

// 文件处理
const handleFileChange = (file) => {
  const isImage = file.raw.type.includes('image')
  if (!isImage) {
    ElMessage.error('请选择图片文件')
    return
  }

  if (file.raw.size / 1024 / 1024 > 8) {
    ElMessage.error('文件大小不能超过8MB!')
    return false
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    imgSrc.value = e.target.result
    currentRatio.value = props.defaultRatioIndex
    onCurrentRatio()
    autoScaled.value = false
    dialogVisible.value = true
  }
  reader.readAsDataURL(file.raw)
}

// 旋转控制
const rotate = (degree) => {
  if (degree === -90) {
    proxy.$refs.cropperRef.rotateLeft()
  } else {
    proxy.$refs.cropperRef.rotateRight()
  }
}

// 实时预览
const handleRealTime = (data) => {
  previews.value = data
  //console.log(data)
}

// 上传处理
const handleUpload = () => {
  uploading.value = true
  proxy.$refs.cropperRef.getCropBlob((blob) => {
    try {
      const file = new File([blob], `${Date.now()}.jpg`, { type: 'image/jpeg' })
      uploadRef.value.clearFiles()
      uploadRef.value.handleStart(file)
      uploadRef.value.submit()

    } catch (error) {
      uploading.value = false
      ElMessage.error('上传失败: ' + error.message)
    }
  })
}

const handleImageSuccess = (res) => {
  const url = res?.data?.file?.url || res?.data?.url || ''
  if (!url) {
    uploading.value = false
    ElMessage.error('上传失败：未获取到图片地址')
    return
  }

  setTimeout(() => {
    uploading.value = false
    dialogVisible.value = false
    previews.value = {}
    ElMessage.success('上传成功')
    emit('on-success', url)
  }, 300)
}

watch(
  () => isMobile.value,
  () => {
    if (!dialogVisible.value) return
    onCurrentRatio()
  }
)

</script>

<style scoped>
:deep(.vue-cropper) {
  background: transparent;
}

.cropper-hidden-trigger {
  display: none;
}
</style>
