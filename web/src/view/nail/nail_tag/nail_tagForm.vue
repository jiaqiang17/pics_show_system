
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="标签名称:" prop="tagName">
    <el-input v-model="formData.tagName" :clearable="true" placeholder="请输入标签名称" />
</el-form-item>
        <el-form-item label="标签图标:" prop="tagIcon">
    <el-input v-model="formData.tagIcon" :clearable="true" placeholder="请输入标签图标" />
</el-form-item>
        <el-form-item label="标签颜色:" prop="tagColor">
    <el-input v-model="formData.tagColor" :clearable="true" placeholder="请输入标签颜色" />
</el-form-item>
        <el-form-item label="标签描述:" prop="description">
    <el-input v-model="formData.description" :clearable="true" placeholder="请输入标签描述" />
</el-form-item>
        <el-form-item label="排序号:" prop="sort">
    <el-input v-model.number="formData.sort" :clearable="true" placeholder="请输入排序号" />
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择状态" :data="nail_tag_statusOptions" style="width:100%" filterable :clearable="true" check-strictly></el-tree-select>
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
  createNailTag,
  updateNailTag,
  findNailTag
} from '@/api/nail/nail_tag'

defineOptions({
    name: 'NailTagForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const nail_tag_statusOptions = ref([])
const formData = ref({
            tagName: '',
            tagIcon: '',
            tagColor: '',
            description: '',
            sort: 0,
            status: '',
        })
// 验证规则
const rule = reactive({
               tagName : [{
                   required: true,
                   message: '请输入标签名称',
                   trigger: ['input','blur'],
               }],
               status : [{
                   required: true,
                   message: '请选择状态',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findNailTag({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    nail_tag_statusOptions.value = await getDictFunc('nail_tag_status')
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
               res = await createNailTag(formData.value)
               break
             case 'update':
               res = await updateNailTag(formData.value)
               break
             default:
               res = await createNailTag(formData.value)
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
