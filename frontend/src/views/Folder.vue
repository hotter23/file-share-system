<template>
  <div class="folder-container">
    <div class="breadcrumb">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/' }">根目录</el-breadcrumb-item>
        <el-breadcrumb-item v-if="parentFolder">{{ parentFolder.folderName }}</el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    
    <div class="toolbar">
      <el-button type="primary" @click="showUploadDialog = true">
        <el-icon><Upload /></el-icon>
        上传文件
      </el-button>
      <el-button @click="goBack">
        <el-icon><Back /></el-icon>
        返回上级
      </el-button>
    </div>

    <!-- 文件夹和文件列表 -->
    <div class="file-list">
      <!-- 文件夹列表 -->
      <div v-if="folders.length > 0" class="folder-section">
        <div class="section-title">文件夹</div>
        <div class="folder-grid">
          <div 
            v-for="folder in folders" 
            :key="folder.folderUuid"
            class="folder-item"
            @click="enterFolder(folder)"
          >
            <el-icon class="folder-icon"><Folder /></el-icon>
            <span class="folder-name">{{ folder.folderName }}</span>
          </div>
        </div>
      </div>

      <!-- 文件列表 -->
      <div class="file-section">
        <el-table :data="files" style="width: 100%">
          <el-table-column label="文件名" prop="fileName">
            <template #default="{ row }">
              <div class="file-cell">
                <el-icon class="file-icon"><Document /></el-icon>
                {{ row.fileName }}
              </div>
            </template>
          </el-table-column>
          <el-table-column label="大小" width="120" prop="fileSize">
            <template #default="{ row }">
              {{ formatFileSize(row.fileSize) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button size="small" @click="downloadFile(row)">下载</el-button>
              <el-button size="small" type="danger" @click="deleteFile(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 上传对话框 -->
    <el-dialog v-model="showUploadDialog" title="上传文件" width="500px">
      <el-upload
        ref="uploadRef"
        drag
        :action="uploadUrl"
        :headers="{ Authorization: `Bearer ${token}` }"
        :data="{ folderId: folderId }"
        :on-success="handleUploadSuccess"
        multiple
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">拖拽文件到此处 或 <em>点击上传</em></div>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import api from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const token = computed(() => userStore.token)
const folderId = computed(() => route.params.folderId)
const uploadUrl = '/api/file/upload'

const files = ref([])
const folders = ref([])
const parentFolder = ref(null)
const showUploadDialog = ref(false)

const loadData = async () => {
  try {
    // 加载文件夹
    const folderRes = await api.get('/api/file/folder/list', { 
      params: { parentId: folderId.value }
    })
    folders.value = folderRes.data.data || []

    // 加载文件
    const fileRes = await api.get('/api/file/list', { 
      params: { folderId: folderId.value, page: 1, size: 100 }
    })
    files.value = fileRes.data.data.records || []
  } catch (e) {
    console.error('加载数据失败', e)
  }
}

const enterFolder = (folder) => {
  router.push(`/folder/${folder.folderUuid}`)
}

const goBack = () => {
  router.back()
}

const downloadFile = (file) => {
  window.open(`/api/file/download/${file.fileUuid}`, '_blank')
}

const deleteFile = async (file) => {
  try {
    await api.delete(`/api/file/${file.fileUuid}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (e) {
    console.error('删除失败', e)
  }
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  showUploadDialog.value = false
  loadData()
}

const formatFileSize = (size) => {
  if (!size) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return `${size.toFixed(2)} ${units[i]}`
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.folder-container {
  background: white;
  border-radius: 8px;
  padding: 20px;
  min-height: 100%;
}

.breadcrumb {
  margin-bottom: 20px;
}

.toolbar {
  display: flex;
  margin-bottom: 20px;
}

.file-list {
  margin-top: 20px;
}

.section-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 12px;
}

.folder-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.folder-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  cursor: pointer;
}

.folder-item:hover {
  border-color: #409eff;
}

.folder-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 8px;
}

.folder-name {
  font-size: 14px;
  text-align: center;
}

.file-cell {
  display: flex;
  align-items: center;
}

.file-icon {
  margin-right: 8px;
  color: #409eff;
}
</style>
