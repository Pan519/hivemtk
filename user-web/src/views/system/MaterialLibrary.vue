<template>
  <div class="material-library">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('素材库管理') }}</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleUpload">{{ $t('上传素材') }}</el-button>
            <el-button type="success" @click="handleCreateCategory">{{ $t('新建分类') }}</el-button>
          </div>
        </div>
      </template>

      
      <div class="search-bar">
        <el-form :inline="true" :model="searchForm">
          <el-form-item :label="$t('auto.MaterialLibrary.d0771a42')">
            <el-select v-model="searchForm.categoryId" :placeholder="$t('auto.MaterialLibrary.1a750305')" clearable>
              <el-option :label="$t('auto.MaterialLibrary.1a750305')" value="" />
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('auto.MaterialLibrary.226b0912')">
            <el-select v-model="searchForm.type" :placeholder="$t('auto.MaterialLibrary.c079e7d5')" clearable>
              <el-option :label="$t('auto.MaterialLibrary.c079e7d5')" value="" />
              <el-option :label="$t('auto.MaterialLibrary.20def794')" value="image" />
              <el-option :label="$t('auto.MaterialLibrary.7fcf42ed')" value="video" />
              <el-option :label="$t('auto.MaterialLibrary.726dd5df')" value="audio" />
              <el-option :label="$t('auto.MaterialLibrary.32536950')" value="document" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('auto.MaterialLibrary.9699a50e')">
            <el-input v-model="searchForm.keyword" :placeholder="$t('auto.MaterialLibrary.21724db1')" clearable />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ $t('auto.MaterialLibrary.e5f71fc3') }}</el-button>
            <el-button @click="handleReset">{{ $t('auto.MaterialLibrary.4b9c3271') }}</el-button>
          </el-form-item>
        </el-form>
      </div>

      
      <div class="material-grid" v-loading="loading">
        <div
          v-for="material in materialList"
          :key="material.id"
          class="material-item"
          role="button"
          tabindex="0"
          @click="handleSelect(material)"
          @keydown.enter.prevent="handleSelect(material)"
          @keydown.space.prevent="handleSelect(material)"
        >
          <div class="material-preview">
            <img
              v-if="material.type === 'image'"
              :src="material.url"
              :alt="material.name"
              @error="handleImageError"
            />
            <div v-else class="material-icon">
              <el-icon :size="48">
                <Document v-if="material.type === 'document'" />
                <VideoPlay v-else-if="material.type === 'video'" />
                <Headset v-else-if="material.type === 'audio'" />
                <Document v-else />
              </el-icon>
            </div>
            <div class="material-actions">
              <el-button type="primary" link @click.stop="handleCopyUrl(material.url)">{{ $t('auto.MaterialLibrary.879058ce') }}</el-button>
              <el-button type="danger" link @click.stop="handleDelete(material)">{{ $t('auto.MaterialLibrary.2f4aaddd') }}</el-button>
            </div>
          </div>
          <div class="material-info">
            <div class="material-name">{{ material.name }}</div>
            <div class="material-meta">
              <span>{{ formatFileSize(material.size) }}</span>
              <span>{{ getTypeLabel(material.type) }}</span>
              <span>使用{{ material.usageCount }}次</span>
            </div>
          </div>
        </div>
      </div>

      
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[12, 24, 48, 96]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    
    <el-dialog
      :title="$t('auto.MaterialLibrary.fc5b616e')"
      v-model="uploadDialogVisible"
      width="600px"
      @close="handleUploadClose"
    >
      <el-form ref="uploadFormRef" :model="uploadForm" label-width="100px">
        <el-form-item :label="$t('auto.MaterialLibrary.5fcb1e4d')" prop="categoryId">
          <el-select v-model="uploadForm.categoryId" :placeholder="$t('auto.MaterialLibrary.8bb820b8')">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.MaterialLibrary.a6fc9e3a')">
          <el-upload
            ref="uploadRef"
            drag
            action="#"
            :auto-upload="false"
            :multiple="true"
            :file-list="fileList"
            :on-change="handleFileChange"
            :before-upload="beforeUpload"
            accept="image/*,video/*,audio/*,.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处或 <em>{{ $t('auto.MaterialLibrary.2c808b02') }}</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持图片、视频、音频、文档格式，单个文件不超过10MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="uploadDialogVisible = false">{{ $t('auto.MaterialLibrary.625fb26b') }}</el-button>
        <el-button type="primary" @click="handleUploadSubmit">{{ $t('auto.MaterialLibrary.020f1ecd') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog
      :title="$t('auto.MaterialLibrary.5c274709')"
      v-model="categoryDialogVisible"
      width="500px"
    >
      <div class="category-list">
        <div
          v-for="category in categories"
          :key="category.id"
          class="category-item"
        >
          <span>{{ category.name }}</span>
          <div class="category-actions">
            <el-button type="primary" link @click="handleEditCategory(category)">{{ $t('auto.MaterialLibrary.95b351c8') }}</el-button>
            <el-button type="danger" link @click="handleDeleteCategory(category)">{{ $t('auto.MaterialLibrary.2f4aaddd') }}</el-button>
          </div>
        </div>
      </div>

      <div class="category-form" v-if="showCategoryForm">
        <el-input
          v-model="categoryForm.name"
          :placeholder="$t('auto.MaterialLibrary.68363f18')"
          style="width: 200px; margin-right: 10px"
        />
        <el-button type="primary" @click="handleSaveCategory">{{ $t('auto.MaterialLibrary.be5fbbe3') }}</el-button>
        <el-button @click="handleCancelCategory">{{ $t('auto.MaterialLibrary.625fb26b') }}</el-button>
      </div>

      <div class="category-add" v-else>
        <el-button type="primary" link @click="handleAddCategory">{{ $t('auto.MaterialLibrary.b6a50669') }}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { rtStr, rtNum, rtDur, rtBool } from '@/utils/runtimeConfig'
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, VideoPlay, Headset, UploadFilled } from '@element-plus/icons-vue'
import {
  getMaterialList,
  uploadMaterial,
  deleteMaterial,
  getMaterialCategories,
  createMaterialCategory,
  updateMaterialCategory,
  deleteMaterialCategory
} from '@/api/material'

const { t } = useI18n()

const loading = ref(false)
const materialList = ref([])
const categories = ref([])
const uploadDialogVisible = ref(false)
const categoryDialogVisible = ref(false)
const showCategoryForm = ref(false)
const uploadRef = ref(null)
const uploadFormRef = ref(null)

const emit = defineEmits(['select'])

const searchForm = reactive({
  categoryId: '',
  type: '',
  keyword: ''
});

const uploadForm = reactive({
  categoryId: ''
})

const categoryForm = reactive({
  id: null,
  name: ''
})

const fileList = ref([])

const pagination = reactive({
  current: 1,
  pageSize: 12,
  total: 0
})

const getTypeLabel = (type) => {
  const labels = {
    image: t('auto.MaterialLibrary.20def794'),
    video: t('auto.MaterialLibrary.7fcf42ed'),
    audio: t('auto.MaterialLibrary.726dd5df'),
    document: t('auto.MaterialLibrary.32536950')
  }
  return labels[type] || t('auto.MaterialLibrary.0d98c747')
}

const formatFileSize = (size) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(1) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(1) + ' GB'
}

const loadMaterialList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      limit: pagination.pageSize,
      ...searchForm
    }
    const response = await getMaterialList(params)
    materialList.value = (response && response.list) || []
    pagination.total = (response && response.total) || 0
  } catch (error) {
    if (error.response?.status !== 401) {
      ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.934bcb9c')))
    }
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const response = await getMaterialCategories()
    categories.value = (response && response.list) || []
  } catch (error) {
    if (error.response?.status !== 401) {
      ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.a8205d01')))
    }
  }
}

const handleSearch = () => {
  pagination.current = 1
  loadMaterialList()
}

const handleReset = () => {
  searchForm.categoryId = ''
  searchForm.type = ''
  searchForm.keyword = ''
  handleSearch()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  loadMaterialList()
}

const handleCurrentChange = (page) => {
  pagination.current = page
  loadMaterialList()
}

const handleUpload = () => {
  uploadDialogVisible.value = true
}

const handleCreateCategory = () => {
  categoryDialogVisible.value = true
}

const handleFileChange = (file, files) => {
  fileList.value = files
}

const beforeUpload = (file) => {
  // 上限改由参数中心下发（frontend_ui.material_max_size_mb）。
  // 注意这是**前端提示用**的副本，真正生效的是后端 misc.upload_max_size_mb（阶段 2a）——
  // 两边不一致会出现「前端放过、后端拦死」。默认取 frontend_ui 的 10MB，
  // 与后端 misc.upload_max_size_mb 的默认 0（=不覆盖、沿用 10MB 兜底）对齐。
  const maxSize = rtNum('frontend_ui', 'material_max_size_mb', 10) * 1024 * 1024
  if (file.size > maxSize) {
    ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.e84c269b')))
    return false
  }
  return true
}

const handleUploadSubmit = async () => {
  if (!uploadForm.categoryId) {
    ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.8bb820b8')))
    return
  }
  
  if (fileList.value.length === 0) {
    ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.9febf311')))
    return
  }

  try {
    for (const file of fileList.value) {
      const formData = new FormData()
      formData.append('file', file.raw)
      formData.append('category_id', uploadForm.categoryId)
      
      await uploadMaterial(formData)
    }
    
    ElMessage.success(i18n.global.t(t('auto.MaterialLibrary.a7699ba7')))
    uploadDialogVisible.value = false
    loadMaterialList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.54e5de42')))
  }
}

const handleUploadClose = () => {
  fileList.value = []
  uploadForm.categoryId = ''
}

const handleCopyUrl = (url) => {
  navigator.clipboard.writeText(url).then(() => {
    ElMessage.success(i18n.global.t(t('auto.MaterialLibrary.deb26c26')))
  }).catch(() => {
    ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.5154ae17')))
  })
}

const handleDelete = async (material) => {
  try {
    await ElMessageBox.confirm(t('auto.MaterialLibrary.6fe54905'), t('auto.MaterialLibrary.02d9819d'), {
      type: 'warning'
    })
    
    await deleteMaterial(material.id)
    ElMessage.success(i18n.global.t(t('auto.MaterialLibrary.0007d170')))
    loadMaterialList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.acf0664a')))
    }
  }
}

const handleImageError = (event) => {
  event.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDgiIGhlaWdodD0iNDgiIHZpZXdCb3g9IjAgMCA0OCA0OCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjQ4IiBoZWlnaHQ9IjQ4IiBmaWxsPSIjRjBGMkY1Ii8+CjxwYXRoIGQ9Ik0yNCAxNkMxOS41ODE3IDE2IDE2IDE5LjU4MTcgMTYgMjRDMTYgMjguNDE4MyAxOS41ODE3IDMyIDI0IDMyQzI4LjQxODMgMzIgMzIgMjguNDE4MyAzMiAyNEMzMiAxOS41ODE3IDI4LjQxODMgMTYgMjQgMTZaIiBmaWxsPSIjQ0NDQ0NDIi8+CjxwYXRoIGQ9Ik0yNCAxOVYyOU0xOSAyNEgyOU0yNCAxNkMxOS41ODE3IDE2IDE2IDE5LjU4MTcgMTYgMjRDMTYgMjguNDE4MyAxOS41ODE3IDMyIDI0IDMyQzI4LjQxODMgMzIgMzIgMjguNDE4MyAzMiAyNEMzMiAxOS41ODE3IDI4LjQxODMgMTYgMjQgMTZaIiBzdHJva2U9IiM5OTk5OTkiIHN0cm9rZS13aWR0aD0iMiIvPgo8L3N2Zz4K'
}

const handleSelect = (material) => {
  emit('select', material);
}

const handleEditCategory = (category) => {
  categoryForm.id = category.id
  categoryForm.name = category.name
  showCategoryForm.value = true
}

const handleDeleteCategory = async (category) => {
  try {
    await ElMessageBox.confirm(t('auto.MaterialLibrary.cc0ca398'), t('auto.MaterialLibrary.02d9819d'), {
      type: 'warning'
    })
    
    await deleteMaterialCategory(category.id)
    ElMessage.success(i18n.global.t(t('auto.MaterialLibrary.0007d170')))
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.acf0664a')))
    }
  }
}

const handleAddCategory = () => {
  categoryForm.id = null
  categoryForm.name = ''
  showCategoryForm.value = true
}

const handleSaveCategory = async () => {
  if (!categoryForm.name.trim()) {
    ElMessage.error(i18n.global.t(t('auto.MaterialLibrary.68363f18')))
    return
  }
  
  try {
    if (categoryForm.id) {
      await updateMaterialCategory(categoryForm.id, categoryForm)
      ElMessage.success(i18n.global.t(t('auto.MaterialLibrary.55aa6366')))
    } else {
      await createMaterialCategory(categoryForm)
      ElMessage.success(i18n.global.t(t('auto.MaterialLibrary.04a691b3')))
    }
    
    showCategoryForm.value = false
    loadCategories()
  } catch (error) {
    ElMessage.error(categoryForm.id ? t('auto.MaterialLibrary.930442e2') : t('auto.MaterialLibrary.a889286a'))
  }
}

const handleCancelCategory = () => {
  showCategoryForm.value = false
  categoryForm.id = null
  categoryForm.name = ''
}

onMounted(() => {
  loadCategories()
  loadMaterialList()
})
</script>

<style scoped>
.material-library {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.search-bar {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.material-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.material-item {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
}

.material-item:hover {
  border-color: #4F46E5;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.material-item:hover .material-actions {
  opacity: 1;
}

.material-preview {
  position: relative;
  width: 100%;
  height: 150px;
  background-color: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.material-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.material-icon {
  color: #909399;
}

.material-actions {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  opacity: 0;
  transition: opacity 0.3s;
}

.material-info {
  padding: 12px;
}

.material-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.material-meta {
  font-size: 12px;
  color: #909399;
  display: flex;
  justify-content: space-between;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.category-list {
  max-height: 300px;
  overflow-y: auto;
  margin-bottom: 20px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #e4e7ed;
}

.category-actions {
  display: flex;
  gap: 10px;
}

.category-form {
  display: flex;
  align-items: center;
  margin-top: 20px;
}

.category-add {
  margin-top: 20px;
  text-align: center;
}
</style>
