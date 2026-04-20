<template>
  <div class="share-container">
    <el-card class="share-card" v-if="!loading">
      <!-- 需要密码 -->
      <div v-if="requirePassword" class="password-form">
        <div class="share-title">此分享需要访问密码</div>
        <el-input 
          v-model="password" 
          type="password" 
          placeholder="请输入访问密码"
          @keyup.enter="accessShare"
        />
        <el-button type="primary" style="width: 100%; margin-top: 16px" @click="accessShare">
          验证密码
        </el-button>
      </div>

      <!-- 显示文件信息 -->
      <div v-else class="file-info">
        <div class="share-title">分享文件</div>
        <div class="file-details">
          <el-icon class="file-icon"><Document /></el-icon>
          <div class="file-meta">
            <div class="file-name">{{ fileInfo.fileName }}</div>
            <div class="file-size">{{ formatFileSize(fileInfo.fileSize) }}</div>
          </div>
        </div>
        <el-button type="primary" style="width: 100%; margin-top: 20px" @click="downloadFile">
          <el-icon><Download /></el-icon>
          下载文件
        </el-button>
      </div>
    </el-card>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../api'
import { ElMessage } from 'element-plus'

const route = useRoute()

const loading = ref(true)
const requirePassword = ref(false)
const password = ref('')
const fileInfo = ref({})

const shareUuid = route.params.shareUuid

const accessShare = async () => {
  try {
    const res = await api.get(`/api/share/${shareUuid}`, {
      params: { password: password.value }
    })
    if (res.data.data.requirePassword) {
      requirePassword.value = true
      if (password.value) {
        ElMessage.error('密码错误')
      }
    } else {
      requirePassword.value = false
      fileInfo.value = res.data.data
    }
  } catch (e) {
    console.error('访问分享失败', e)
    ElMessage.error('分享不存在或已过期')
  } finally {
    loading.value = false
  }
}

const downloadFile = () => {
  const url = `/api/share/${shareUuid}/download`
  const fullUrl = password.value ? `${url}?password=${password.value}` : url
  window.open(fullUrl, '_blank')
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
  accessShare()
})
</script>

<style scoped>
.share-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.share-card {
  width: 400px;
}

.share-title {
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 20px;
}

.password-form {
  padding: 10px;
}

.file-info {
  padding: 10px;
}

.file-details {
  display: flex;
  align-items: center;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.file-icon {
  font-size: 48px;
  color: #409eff;
  margin-right: 16px;
}

.file-meta {
  flex: 1;
}

.file-name {
  font-size: 16px;
  font-weight: 500;
  word-break: break-all;
}

.file-size {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: white;
  font-size: 16px;
}
</style>
