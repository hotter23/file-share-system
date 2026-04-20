<template>
  <div class="home-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button type="primary" @click="showUploadDialog = true">
        <el-icon><Upload /></el-icon>
        上传文件
      </el-button>
      <el-button @click="showFolderDialog = true">
        <el-icon><FolderAdd /></el-icon>
        新建文件夹
      </el-button>
      <el-input
        v-model="keyword"
        placeholder="搜索文件"
        style="width: 200px; margin-left: auto"
        clearable
        @keyup.enter="searchFiles"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 文件列表 -->
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
            <el-dropdown trigger="click" @command="handleFolderCommand($event, folder)">
              <el-icon class="folder-action"><MoreFilled /></el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="delete">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>

      <!-- 文件列表 -->
      <div class="file-section">
        <div class="section-title">文件</div>
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
          <el-table-column label="上传时间" width="180" prop="createdAt">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button size="small" @click="downloadFile(row)">下载</el-button>
              <el-button size="small" @click="shareFile(row)">分享</el-button>
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
        :data="{ folderId: currentFolderId }"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :before-upload="beforeUpload"
        multiple
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">拖拽文件到此处 或 <em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">单文件最大5GB</div>
        </template>
      </el-upload>
    </el-dialog>

    <!-- 新建文件夹对话框 -->
    <el-dialog v-model="showFolderDialog" title="新建文件夹" width="400px">
      <el-form @submit.prevent="createFolder">
        <el-form-item label="文件夹名称">
          <el-input v-model="newFolderName" placeholder="请输入文件夹名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFolderDialog = false">取消</el-button>
        <el-button type="primary" @click="createFolder">创建</el-button>
      </template>
    </el-dialog>

    <!-- 分享对话框 -->
    <el-dialog v-model="showShareDialog" title="分享文件" width="400px">
      <el-form>
        <el-form-item label="分享有效期">
          <el-select v-model="shareExpireType" style="width: 100%">
            <el-option label="1小时" value="1H" />
            <el-option label="24小时" value="24H" />
            <el-option label="7天" value="7D" />
            <el-option label="永久" value="PERMANENT" />
          </el-select>
        </el-form-item>
        <el-form-item label="访问密码（可选）">
          <el-input v-model="sharePassword" placeholder="不设置则无需密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showShareDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmShare">创建分享</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import api from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const token = computed(() => userStore.token)
const uploadUrl = '/api/file/upload'

const keyword = ref('')
const currentFolderId = ref(null)
const files = ref([])
const folders = ref([])

const showUploadDialog = ref(false)
const showFolderDialog = ref(false)
const showShareDialog = ref(false)
const newFolderName = ref('')
const shareExpireType = ref('24H')
const sharePassword = ref('')
const currentShareFile = ref(null)

const uploadRef = ref(null)

// 加载文件列表
const loadFiles = async () => {
  try {
    const params = { page: 1, size: 100 }
    if (keyword.value) params.keyword = keyword.value
    if (currentFolderId.value) params.folderId = currentFolderId.value
    
    const res = await api.get('/api/file/list', { params })
    files.value = res.data.data.records || []
  } catch (e) {
    console.error('加载文件列表失败', e)
  }
}

// 加载文件夹列表
const loadFolders = async () => {
  try {
    const params = {}
    if (currentFolderId.value) params.parentId = currentFolderId.value
    
    const res = await api.get('/api/file/folder/list', { params })
    folders.value = res.data.data || []
  } catch (e) {
    console.error('加载文件夹列表失败', e)
  }
}

// 搜索文件
const searchFiles = () => {
  loadFiles()
}

// 进入文件夹
const enterFolder = (folder) => {
  router.push(`/folder/${folder.folderUuid}`)
}

// 创建文件夹
const createFolder = async () => {
  if (!newFolderName.value.trim()) {
    ElMessage.warning('请输入文件夹名称')
    return
  }
  try {
    const data = { folderName: newFolderName.value }
    if (currentFolderId.value) data.parentId = currentFolderId.value
    await api.post('/api/file/folder', data)
    ElMessage.success('创建成功')
    showFolderDialog.value = false
    newFolderName.value = ''
    loadFolders()
  } catch (e) {
    console.error('创建文件夹失败', e)
  }
}

// 文件夹操作
const handleFolderCommand = async (command, folder) => {
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm('确定要删除此文件夹吗？', '提示', { type: 'warning' })
      await api.delete(`/api/file/folder/${folder.folderUuid}`)
      ElMessage.success('删除成功')
      loadFolders()
    } catch (e) {
      if (e !== 'cancel') console.error('删除失败', e)
    }
  }
}

// 上传前校验
const beforeUpload = (file) => {
  const maxSize = 5 * 1024 * 1024 * 1024 // 5GB
  if (file.size > maxSize) {
    ElMessage.error('文件大小超出限制（最大5GB）')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = () => {
  ElMessage.success('上传成功')
  showUploadDialog.value = false
  loadFiles()
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error('上传失败')
}

// 下载文件
const downloadFile = async (file) => {
  window.open(`/api/file/download/${file.fileUuid}`, '_blank')
}

// 分享文件
const shareFile = (file) => {
  currentShareFile.value = file
  showShareDialog.value = true
}

// 确认分享
const confirmShare = async () => {
  try {
    const res = await api.post('/api/share', {
      fileUuid: currentShareFile.value.fileUuid,
      expireType: shareExpireType.value,
      password: sharePassword.value || null
    })
    const shareUrl = `${window.location.origin}/share/${res.data.data.shareUuid}`
    await ElMessageBox.alert(`分享链接：${shareUrl}`, '分享成功', {
      confirmButtonText: '复制链接',
      callback: () => {
        navigator.clipboard.writeText(shareUrl)
        ElMessage.success('链接已复制')
      }
    })
    showShareDialog.value = false
  } catch (e) {
    console.error('分享失败', e)
  }
}

// 删除文件
const deleteFile = async (file) => {
  try {
    await ElMessageBox.confirm('确定要删除此文件吗？', '提示', { type: 'warning' })
    await api.delete(`/api/file/${file.fileUuid}`)
    ElMessage.success('删除成功')
    loadFiles()
  } catch (e) {
    if (e !== 'cancel') console.error('删除失败', e)
  }
}

// 格式化文件大小
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

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  loadFiles()
  loadFolders()
  userStore.fetchUserInfo()
})
</script>

<style scoped>
.home-container {
  background: white;
  border-radius: 8px;
  padding: 20px;
  min-height: 100%;
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
  transition: all 0.3s;
}

.folder-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.2);
}

.folder-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 8px;
}

.folder-name {
  font-size: 14px;
  color: #303133;
  text-align: center;
  word-break: break-all;
}

.folder-action {
  position: absolute;
  top: 8px;
  right: 8px;
  cursor: pointer;
  color: #909399;
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
