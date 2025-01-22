<template>
  <el-form
    :model="taskForm"
    :rules="rules"
    ref="formRef"
    label-width="100px"
    @submit.prevent
  >
    <el-form-item label="任务名称" prop="name">
      <el-input
        v-model="taskForm.name"
        placeholder="请输入任务名称"
        clearable
      />
    </el-form-item>
    <el-form-item label="任务描述" prop="description">
      <el-input
        v-model="taskForm.description"
        type="textarea"
        :rows="4"
        placeholder="请输入任务描述"
        resize="none"
      />
    </el-form-item>
    <el-form-item label="依赖任务" prop="dependencies">
      <el-select
        v-model="taskForm.dependencies"
        multiple
        placeholder="请选择依赖任务"
        style="width: 100%"
        collapse-tags
        collapse-tags-tooltip
      >
        <el-option
          v-for="task in availableTasks"
          :key="task.id"
          :label="task.name"
          :value="task.id"
        />
      </el-select>
    </el-form-item>
    <el-form-item class="form-buttons">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">
        保存
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

const props = defineProps<{
  availableTasks: Array<{ id: number; name: string }>
}>()

const emit = defineEmits<{
  (e: 'submit', task: any): void
  (e: 'cancel'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const taskForm = ref({
  name: '',
  description: '',
  dependencies: [] as number[]
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '描述不能超过 500 个字符', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    submitting.value = true
    await formRef.value.validate()
    emit('submit', { ...taskForm.value })
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    submitting.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.el-form {
  max-width: 100%;
}

.form-buttons {
  margin-bottom: 0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input__wrapper),
:deep(.el-textarea__inner),
:deep(.el-select__wrapper) {
  box-shadow: none;
  border: 1px solid #dcdfe6;
}

:deep(.el-input__wrapper:hover),
:deep(.el-textarea__inner:hover),
:deep(.el-select__wrapper:hover) {
  border-color: #409eff;
}

:deep(.el-form-item.is-error .el-input__wrapper),
:deep(.el-form-item.is-error .el-textarea__inner) {
  border-color: #f56c6c;
}
</style> 