
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="款式名称:" prop="styleName">
    <el-input v-model="formData.styleName" :clearable="true" placeholder="请输入款式名称" />
</el-form-item>
        <el-form-item label="美甲图片集:" prop="images">
    <SelectImage
     multiple
     v-model="formData.images"
     file-type="image"
     />
</el-form-item>
        <el-form-item label="款式介绍:" prop="description">
    <RichEdit v-model="formData.description"/>
</el-form-item>
        <el-form-item label="是否推荐:" prop="isRecommended">
    <el-switch v-model="formData.isRecommended" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
        <el-form-item label="排序号:" prop="sort">
    <el-input v-model.number="formData.sort" :clearable="true" placeholder="请输入排序号" />
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择状态" :data="nail_style_statusOptions" style="width:100%" filterable :clearable="true" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="关联的标签ID列表:" prop="tagIds">
    <el-select multiple v-model="formData.tagIds" placeholder="请选择关联的标签ID列表" filterable style="width:100%" :clearable="true">
        <el-option v-for="(item,key) in dataSource.tagIds" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
    getNailStyleDataSource,
  createNailStyle,
  updateNailStyle,
  findNailStyle
} from '@/api/nail/nail_style'

defineOptions({
    name: 'NailStyleForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'
// 数组控制组件
import ArrayCtrl from '@/components/arrayCtrl/arrayCtrl.vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const nail_style_statusOptions = ref([])
const formData = ref({
            styleName: '',
            images: [],
            description: '',
            isRecommended: false,
            sort: 0,
            status: '',
            tagIds: [],
        })
// 验证规则
const rule = reactive({
               styleName : [{
                   required: true,
                   message: '请输入款式名称',
                   trigger: ['input','blur'],
               }],
               images : [{
                   required: true,
                   message: '请上传美甲图片',
                   trigger: ['input','blur'],
               }],
               status : [{
                   required: true,
                   message: '请选择状态',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getNailStyleDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findNailStyle({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    nail_style_statusOptions.value = await getDictFunc('nail_style_status')
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createNailStyle(formData.value)
               break
             case 'update':
               res = await updateNailStyle(formData.value)
               break
             default:
               res = await createNailStyle(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
