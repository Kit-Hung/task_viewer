<template>
  <div class="home">
    <div class="header">
      <div class="header-left">
        <h1>任务依赖管理系统</h1>
      </div>
      <div class="header-right">
        <el-button type="primary" size="large" @click="showTaskForm = true">
          <el-icon class="el-icon--left"><Plus /></el-icon>
          添加任务
        </el-button>
      </div>
    </div>

    <div class="graph-container">
      <dependency-graph
        :tasks="tasks"
        :dependencies="dependencies"
        :selected-tasks="selectedTasks"
        @select="handleTaskSelect"
        @node-dblclick="handleNodeDblClick"
      />
    </div>
    
    <div v-if="selectedTasks.length > 0" class="action-bar">
      <div class="selected-count">
        已选择 {{ selectedTasks.length }} 个任务
      </div>
      <div class="action-buttons">
        <el-button @click="selectedTasks = []">
          清除选择
        </el-button>
        <el-button type="danger" @click="handleDeleteTasks">
          <el-icon class="el-icon--left"><Delete /></el-icon>
          删除任务
        </el-button>
        <el-button type="primary" @click="showDependencies = true">
          <el-icon class="el-icon--left"><Connection /></el-icon>
          查看依赖关系
        </el-button>
      </div>
    </div>

    <el-dialog
      v-model="showTaskForm"
      title="添加任务"
      width="500px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <task-form
        :available-tasks="tasks"
        @submit="handleTaskSubmit"
        @cancel="showTaskForm = false"
      />
    </el-dialog>

    <el-dialog
      v-model="showDependencies"
      title="任务依赖关系"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="dependencies-view">
        <div class="selected-section">
          <h3>已选择的任务</h3>
          <div class="tags-wrapper">
            <el-tag
              v-for="taskId in selectedTasks"
              :key="taskId"
              class="task-tag"
              type="primary"
              effect="light"
            >
              {{ getTaskName(taskId) }}
            </el-tag>
          </div>
        </div>
        <el-divider>
          <el-icon><ArrowDown /></el-icon>
        </el-divider>
        <div class="dependencies-section">
          <h3>所需前置依赖</h3>
          <div class="tags-wrapper">
            <el-tag
              v-for="taskId in requiredDependencies"
              :key="taskId"
              class="task-tag"
              type="success"
              effect="light"
            >
              {{ getTaskName(taskId) }}
            </el-tag>
          </div>
          <div v-if="requiredDependencies.length === 0" class="empty-text">
            无需前置依赖
          </div>
        </div>
        <div class="execution-order">
          <h3>建议执行顺序</h3>
          <div class="order-list">
            <div
              v-for="(taskId, index) in executionOrder"
              :key="taskId"
              class="order-item"
            >
              <span class="order-number">{{ index + 1 }}</span>
              <el-tag size="large">{{ getTaskName(taskId) }}</el-tag>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Plus, Refresh, ArrowDown, Connection, Delete } from '@element-plus/icons-vue'
import DependencyGraph from '../components/DependencyGraph/index.vue'
import TaskForm from '../components/TaskForm/index.vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'

// 定义接口
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

const tasks = ref<Task[]>([])
const dependencies = ref<{ source: number; target: number }[]>([])
const selectedTasks = ref<number[]>([])
const showTaskForm = ref(false)
const showDependencies = ref(false)

const router = useRouter()

// 获取任务名称
const getTaskName = (taskId: number) => {
  const task = tasks.value.find(t => t.id === taskId)
  return task ? task.name : `任务 ${taskId}`
}

// 加载任务数据
const loadTasks = async () => {
  try {
    console.log('开始加载任务数据...')
    const response = await fetch('http://localhost:8080/api/tasks', {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    })
    
    console.log('收到响应:', response.status, response.statusText)
    if (!response.ok) {
      const text = await response.text()
      console.error('响应错误:', text)
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const data = await response.json()
    console.log('解析数据:', data)
    
    tasks.value = data.tasks || []
    dependencies.value = (data.dependencies || []).map((dep: Dependency) => ({
      source: dep.source_id,
      target: dep.target_id
    }))
    
    console.log('数据加载完成:', { tasks: tasks.value, dependencies: dependencies.value })
  } catch (error) {
    console.error('加载任务失败:', error)
    ElMessage.error('加载任务失败: ' + (error instanceof Error ? error.message : String(error)))
  }
}

// 处理任务选择
const handleTaskSelect = async (taskIds: number[]) => {
  selectedTasks.value = taskIds
}

// 添加删除任务的方法
const handleDeleteTasks = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedTasks.value.length} 个任务吗？删除后将无法恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    for (const taskId of selectedTasks.value) {
      const response = await fetch(`http://localhost:8080/api/tasks/${taskId}`, {
        method: 'DELETE'
      })

      if (!response.ok) {
        throw new Error(`Failed to delete task ${taskId}`)
      }
    }

    await loadTasks() // 重新加载任务列表
    selectedTasks.value = [] // 清空选择
    ElMessage.success('任务已删除')
  } catch (error) {
    if (error instanceof Error && error.message.startsWith('Failed to delete task')) {
      ElMessage.error('删除任务失败')
    }
  }
}

// 处理任务提交
const handleTaskSubmit = async (task: any) => {
  try {
    console.log('提交任务数据:', task)
    const response = await fetch('http://localhost:8080/api/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      body: JSON.stringify({
        name: task.name,
        description: task.description,
        dependencies: task.dependencies
      })
    })

    if (!response.ok) {
      const text = await response.text()
      console.error('创建任务失败:', text)
      throw new Error('Failed to create task')
    }

    await loadTasks() // 重新加载任务列表
    showTaskForm.value = false
    ElMessage.success('任务创建成功')
  } catch (error) {
    console.error('创建任务失败:', error)
    ElMessage.error('创建任务失败')
  }
}

// 添加处理双击事件的方法
const handleNodeDblClick = (taskId: number) => {
  router.push(`/task/${taskId}`)
}

// 计算所需的前置依赖（排除已选任务）
const requiredDependencies = computed(() => {
  const allDeps = new Set<number>()
  
  const addDependencies = (taskId: number) => {
    dependencies.value.forEach(dep => {
      if (dep.target === taskId) {
        allDeps.add(dep.source)
        // 递归查找更上层的依赖
        addDependencies(dep.source)
      }
    })
  }

  // 从选中的任务开始查找所有前置依赖
  selectedTasks.value.forEach(taskId => {
    addDependencies(taskId)
  })

  // 排除已选择的任务
  selectedTasks.value.forEach(taskId => {
    allDeps.delete(taskId)
  })

  return Array.from(allDeps)
})

// 计算执行顺序
const executionOrder = computed(() => {
  const result: number[] = []
  const visited = new Set<number>()
  
  const visit = (taskId: number) => {
    if (visited.has(taskId)) return
    
    // 先处理前置依赖
    dependencies.value.forEach(dep => {
      if (dep.target === taskId) {
        visit(dep.source)
      }
    })
    
    if (!visited.has(taskId)) {
      visited.add(taskId)
      result.push(taskId)
    }
  }
  
  // 先添加所有前置依赖
  requiredDependencies.value.forEach(taskId => {
    visit(taskId)
  })
  
  // 再添加选中的任务
  selectedTasks.value.forEach(taskId => {
    if (!visited.has(taskId)) {
      visited.add(taskId)
      result.push(taskId)
    }
  })
  
  return result
})

// 在组件挂载时加载任务数据
onMounted(() => {
  loadTasks()
})
</script>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #fff;
}

.header {
  height: 64px;
  padding: 0 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #dcdfe6;
}

.header-left {
  flex: 1;
}

.header-right {
  padding-left: 32px;
}

.graph-container {
  flex: 1;
  height: calc(100vh - 64px);
  padding: 24px;
  position: relative;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #fff;
  padding: 16px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 -2px 12px 0 rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.selected-count {
  color: #606266;
  font-size: 14px;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.dependencies-view {
  padding: 0 20px;
}

.selected-section,
.dependencies-section {
  margin-bottom: 24px;
}

h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #303133;
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.task-tag {
  font-size: 14px;
}

.empty-text {
  color: #909399;
  font-size: 14px;
}

.execution-order {
  margin-top: 32px;
}

.order-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.order-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.order-number {
  width: 24px;
  height: 24px;
  border-radius: 12px;
  background-color: #f5f7fa;
  color: #909399;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}
</style> 