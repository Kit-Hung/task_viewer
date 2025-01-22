<template>
  <div class="task-detail">
    <el-container>
      <el-header height="64px">
        <div class="header-content">
          <div class="header-left">
            <el-button link @click="router.back()">
              <el-icon class="el-icon--left"><Back /></el-icon>
              返回
            </el-button>
            <h1>任务详情</h1>
          </div>
          <div class="header-right">
            <el-button type="danger" @click="handleDelete">删除任务</el-button>
            <el-button type="primary" @click="handleEdit">编辑任务</el-button>
          </div>
        </div>
      </el-header>
      <el-main>
        <div class="main-content">
          <el-row :gutter="24">
            <el-col :span="16">
              <el-card class="detail-card">
                <template #header>
                  <div class="card-header">
                    <span>基本信息</span>
                    <el-tag :type="task?.status === 'completed' ? 'success' : 'warning'">
                      {{ task?.status === 'completed' ? '已完成' : '进行中' }}
                    </el-tag>
                  </div>
                </template>
                <div class="task-info">
                  <div class="info-item">
                    <label>任务名称：</label>
                    <span>{{ task?.name }}</span>
                  </div>
                  <div class="info-item">
                    <label>创建时间：</label>
                    <span>{{ formatDate(task?.createdAt) }}</span>
                  </div>
                  <div class="info-item description">
                    <label>任务描述：</label>
                    <p>{{ task?.description || '暂无描述' }}</p>
                  </div>
                </div>
              </el-card>

              <el-card class="dependency-card">
                <template #header>
                  <div class="card-header">
                    <span>依赖关系图</span>
                  </div>
                </template>
                <div class="graph-container">
                  <dependency-graph
                    :tasks="relatedTasks"
                    :dependencies="relatedDependencies"
                    :selected-tasks="[Number(id)]"
                    @select="handleTaskSelect"
                  />
                </div>
              </el-card>
            </el-col>

            <el-col :span="8">
              <el-card class="dependencies-card">
                <template #header>
                  <div class="card-header">
                    <span>依赖任务</span>
                  </div>
                </template>
                <div v-if="dependencies.length > 0" class="dependencies-list">
                  <div class="list-section">
                    <h3>前置依赖</h3>
                    <el-empty v-if="upstreamDependencies.length === 0" description="无前置依赖" />
                    <el-tag
                      v-for="taskId in upstreamDependencies"
                      :key="taskId"
                      class="dependency-tag"
                      @click="navigateToTask(taskId)"
                    >
                      {{ getTaskName(taskId) }}
                    </el-tag>
                  </div>
                  <el-divider />
                  <div class="list-section">
                    <h3>后置依赖</h3>
                    <el-empty v-if="downstreamDependencies.length === 0" description="无后置依赖" />
                    <el-tag
                      v-for="taskId in downstreamDependencies"
                      :key="taskId"
                      class="dependency-tag"
                      @click="navigateToTask(taskId)"
                    >
                      {{ getTaskName(taskId) }}
                    </el-tag>
                  </div>
                </div>
                <el-empty v-else description="暂无依赖关系" />
              </el-card>
            </el-col>
          </el-row>
        </div>
      </el-main>
    </el-container>

    <el-dialog
      v-model="showEditForm"
      title="编辑任务"
      width="500px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <task-form
        :available-tasks="availableTasks"
        @submit="handleTaskSubmit"
        @cancel="showEditForm = false"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Back } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import DependencyGraph from '../components/DependencyGraph/index.vue'
import TaskForm from '../components/TaskForm/index.vue'

interface Task {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

interface Dependency {
  id: number
  source_id: number
  target_id: number
  created_at: string
}

const props = defineProps<{
  id: string
}>()

const router = useRouter()
const showEditForm = ref(false)

const task = ref<Task | null>(null)
const tasks = ref<Task[]>([])
const dependencies = ref<{ source: number; target: number }[]>([])

// 加载任务详情
const loadTaskDetail = async () => {
  try {
    // 加载任务详情
    const response = await fetch(`http://localhost:8080/api/tasks/${props.id}`)
    if (!response.ok) {
      throw new Error('Task not found')
    }
    const data = await response.json()
    task.value = data.task
    
    // 加载相关任务
    const tasksResponse = await fetch('http://localhost:8080/api/tasks')
    const tasksData = await tasksResponse.json()  // 修复：使用 tasksResponse 而不是 response
    tasks.value = tasksData.tasks
    dependencies.value = tasksData.dependencies.map((dep: Dependency) => ({
      source: dep.source_id,
      target: dep.target_id
    }))
  } catch (error) {
    console.error('Failed to load task:', error)
    ElMessage.error('加载任务失败')
    router.push('/')
  }
}

// 计算相关任务和依赖
const relatedTasks = computed(() => {
  const taskIds = new Set<number>()
  taskIds.add(Number(props.id))
  
  dependencies.value.forEach(dep => {
    if (dep.source === Number(props.id) || dep.target === Number(props.id)) {
      taskIds.add(dep.source)
      taskIds.add(dep.target)
    }
  })
  
  return tasks.value.filter(task => taskIds.has(task.id))
})

const relatedDependencies = computed(() => {
  return dependencies.value.filter(dep => 
    dep.source === Number(props.id) || 
    dep.target === Number(props.id)
  )
})

// 计算上下游依赖
const upstreamDependencies = computed(() => {
  const result = new Set<number>()
  dependencies.value.forEach(dep => {
    if (dep.target === Number(props.id)) {
      result.add(dep.source)
    }
  })
  return Array.from(result)
})

const downstreamDependencies = computed(() => {
  const result = new Set<number>()
  dependencies.value.forEach(dep => {
    if (dep.source === Number(props.id)) {
      result.add(dep.target)
    }
  })
  return Array.from(result)
})

// 可选的依赖任务（排除自己和已有依赖）
const availableTasks = computed(() => {
  return tasks.value.filter(t => t.id !== Number(props.id))
})

// 获取任务名称
const getTaskName = (taskId: number) => {
  const task = tasks.value.find(t => t.id === taskId)
  return task ? task.name : `任务 ${taskId}`
}

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 处理任务选择
const handleTaskSelect = (taskIds: number[]) => {
  if (taskIds.length === 1 && taskIds[0] !== Number(props.id)) {
    navigateToTask(taskIds[0])
  }
}

// 跳转到任务详情
const navigateToTask = (taskId: number) => {
  router.push(`/task/${taskId}`)
}

// 处理编辑
const handleEdit = () => {
  showEditForm.value = true
}

// 处理删除
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该任务吗？删除后将无法恢复。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const response = await fetch(`http://localhost:8080/api/tasks/${props.id}`, {
      method: 'DELETE'
    })
    
    if (!response.ok) {
      throw new Error('Failed to delete task')
    }
    
    ElMessage.success('任务已删除')
    router.push('/')
  } catch (error) {
    if (error instanceof Error && error.message === 'Failed to delete task') {
      ElMessage.error('删除任务失败')
    }
  }
}

// 处理任务提交
const handleTaskSubmit = async (formData: any) => {
  try {
    const response = await fetch(`http://localhost:8080/api/tasks/${props.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: formData.name,
        description: formData.description
      })
    })

    if (!response.ok) {
      throw new Error('Failed to update task')
    }

    const updatedTask = await response.json()
    task.value = updatedTask
    showEditForm.value = false
    ElMessage.success('任务已更新')
  } catch (error) {
    console.error('Failed to update task:', error)
    ElMessage.error('更新任务失败')
  }
}

// 在组件挂载时加载任务数据
onMounted(() => {
  loadTaskDetail()
})
</script>

<style scoped>
.task-detail {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
}

.el-container {
  height: 100%;
  min-height: 100vh;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.header-content {
  max-width: 1400px;
  height: 100%;
  margin: 0 auto;
  padding: 0 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-right {
  display: flex;
  gap: 12px;
}

.main-content {
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.detail-card,
.dependency-card,
.dependencies-card {
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 4px;
}

.card-header span {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: flex-start;
}

.info-item label {
  width: 100px;
  color: #606266;
  font-weight: 500;
}

.info-item.description {
  align-items: flex-start;
}

.info-item.description p {
  margin: 0;
  white-space: pre-wrap;
}

.graph-container {
  height: 400px;
  border-radius: 4px;
  overflow: hidden;
}

.dependencies-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.list-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.list-section h3 {
  margin: 0;
  font-size: 14px;
  color: #606266;
}

.dependency-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.dependency-tag:hover {
  transform: translateY(-2px);
}

h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.el-main {
  padding: 24px;
}

@media (max-width: 1440px) {
  .main-content {
    max-width: 1200px;
  }
}

@media (max-width: 1200px) {
  .main-content {
    max-width: 960px;
  }
}
</style> 