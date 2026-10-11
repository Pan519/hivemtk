<template>
  <div class="kuaishou-card-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('快手卡片管理') }}</span>
          <div class="header-actions">
            <el-button type="primary" icon="DataAnalysis" @click="goToStats">
              {{ $t('统计分析') }}
            </el-button>
            <el-button type="primary" icon="Plus" @click="handleAdd">
              {{ $t('添加卡片') }}
            </el-button>
          </div>
        </div>
      </template>

      
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('auto.List.9699a50e')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('auto.List.9761111c')" clearable />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3fea7ca7')">
          <el-select v-model="searchForm.is_active" :placeholder="$t('auto.List.708c9d6d')" clearable>
            <el-option :label="$t('auto.List.83a991d7')" :value="true" />
            <el-option :label="$t('auto.List.710ad08b')" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('auto.List.e5f71fc3') }}</el-button>
          <el-button @click="resetSearch">{{ $t('auto.List.4b9c3271') }}</el-button>
        </el-form-item>
      </el-form>

      
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('auto.List.32c65d8d')" min-width="200">
          <template #default="scope">
            <span class="title-text">{{ scope.row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="image_url" :label="$t('auto.List.20def794')" width="120">
          <template #default="scope">
            <el-image
              v-if="scope.row.image_url"
              :src="scope.row.image_url"
              @error="onImgError"
              :preview-src-list="[scope.row.image_url]"
              fit="cover"
              style="width: 80px; height: 60px; border-radius: 4px;"
            />
            <span v-else class="no-image">{{ $t('auto.List.5423e287') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" :label="$t('auto.List.3bdd08ad')" min-width="300">
          <template #default="scope">
            <el-tooltip :content="scope.row.description" placement="top">
              <div class="description-text">{{ scope.row.description }}</div>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="redirect_url" :label="$t('auto.List.dec2eb71')" min-width="200">
          <template #default="scope">
            <el-tooltip :content="scope.row.redirect_url" placement="top">
              <div class="redirect-url">{{ scope.row.redirect_url }}</div>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" :label="$t('auto.List.88bc4d7a')" width="100" />
        <el-table-column prop="like_count" :label="$t('auto.List.4506d555')" width="100" />
        <el-table-column prop="share_count" :label="$t('auto.List.34688e27')" width="100" />
        <el-table-column prop="is_active" :label="$t('auto.List.3fea7ca7')" width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.is_active"
              active-color="#2684FF"
              inactive-color="#dcdfe6"
              @change="handleStatusChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.List.eca37cb0')" width="160">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button size="small" type="info" @click="handleStats(scope.row)">
              统计
            </el-button>
            <el-button size="small" type="warning" @click="handleCopyLink(scope.row)">
              复制链接
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.21efd88b')" />
        </template>
      </el-table>

      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="1000px"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <div class="form-section">
          <el-form :model="formData" label-width="100px" :rules="rules" ref="formRef">
            <el-form-item :label="$t('auto.List.3944f217')" prop="title">
              <el-input
                v-model="formData.title"
                :placeholder="$t('auto.List.3fd55a4e')"
                maxlength="100"
                show-word-limit
              />
            </el-form-item>
            <el-form-item :label="$t('auto.List.9546fa92')" prop="description">
              <el-input
                v-model="formData.description"
                :placeholder="$t('auto.List.705bdfc9')"
                type="textarea"
                rows="4"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
            <el-form-item :label="$t('auto.List.92b14d1a')" prop="image_url">
              <el-input
                v-model="formData.image_url"
                :placeholder="$t('auto.List.3af187ba')"
                maxlength="500"
                show-word-limit
              />
              <div class="image-preview">
                <el-image
                  v-if="formData.image_url"
                  :src="formData.image_url"
                  fit="cover"
                  style="width: 200px; height: 150px; border-radius: 4px; margin-top: 10px;"
                />
              </div>
            </el-form-item>
            <el-form-item :label="$t('auto.List.dec2eb71')" prop="redirect_url">
              <el-input
                v-model="formData.redirect_url"
                :placeholder="$t('auto.List.c63a5309')"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
            <el-form-item :label="$t('auto.List.0efb0786')" prop="domain_pool_id">
              <el-select 
                v-model="formData.domain_pool_id" 
                :placeholder="$t('auto.List.8ee08e9a')" 
                clearable
                style="width: 100%"
              >
                <el-option
                  v-for="item in domainPoolOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                >
                  <span style="float: left">{{ item.name }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ item.domain }}</span>
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('auto.List.14d34236')" prop="tags">
              <el-input 
                v-model="formData.tags" 
                :placeholder="$t('auto.List.dc1fb70f')" 
                maxlength="255"
                show-word-limit
              />
            </el-form-item>
            <el-form-item :label="$t('auto.List.03c9e27b')">
              <el-switch v-model="formData.is_active" />
            </el-form-item>
          </el-form>
        </div>
        <div class="preview-section">
          <KuaishouCardPreview :card-data="formData" />
        </div>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ $t('auto.List.e83a256e') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>

const onImgError = (e) => {
  const img = e.target
  if (img && !img.dataset.fallback) {
    img.dataset.fallback = '1'
    img.src = 'data:image/svg+xml;base64,' + btoa('<svg xmlns="http://www.w3.org/2000/svg" width="80" height="80"><rect width="80" height="80" fill="%23f0f2f5"/><text x="40" y="46" font-size="12" text-anchor="middle" fill="%2390a0b0">%E6%97%A0%E5%9B%BE</text></svg>')
  }
};
import i18n from '@/i18n'

import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'

import { useRouter } from 'vue-router'
import KuaishouCardPreview from '@/components/KuaishouCardPreview.vue'
import { domainPoolApi } from '@/api/domainPool'
import * as kuaishouCardApi from '@/api/kuaishouCard'

const formatDate = (date, format) => {
  if (!date) return ''
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  
  if (format === 'YYYY-MM-DD') {
    return `${year}-${month}-${day}`
  }
  return new Date(date).toLocaleDateString()
};

const router = useRouter();

const { t } = useI18n()

const loading = ref(false);
const dialogVisible = ref(false)
const formRef = ref(null)
const currentCardId = ref(null)

const domainPoolOptions = ref([]);

const searchForm = reactive({
  keyword: '',
  is_active: null
});

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
});

const tableData = ref([]);

const formData = reactive({
  title: '',
  description: '',
  image_url: '',
  redirect_url: '',
  tags: '',
  is_active: true,
  domain_pool_id: null
});

const rules = {
  title: [
    { required: true, message: i18n.global.t(t('auto.List.3fd55a4e')), trigger: 'blur' },
    { min: 2, max: 100, message: i18n.global.t(t('auto.List.ea55d740')), trigger: 'blur' }
  ],
  description: [
    { required: true, message: i18n.global.t(t('auto.List.705bdfc9')), trigger: 'blur' },
    { min: 5, max: 500, message: i18n.global.t(t('auto.List.ede4cfdc')), trigger: 'blur' }
  ],
  image_url: [
    { required: true, message: i18n.global.t(t('auto.List.3af187ba')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.41368ee8')), trigger: 'blur' }
  ],
  redirect_url: [
    { required: true, message: i18n.global.t(t('auto.List.c63a5309')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.a7f6ca3f')), trigger: 'blur' }
  ],
  tags: [
    { max: 255, message: i18n.global.t(t('auto.List.cfd05181')), trigger: 'blur' }
  ]
};

const dialogTitle = computed(() => {
  return currentCardId.value ? t('auto.List.4c2c1f23') : t('auto.List.b1aed49c')
});

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      ...searchForm,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    }
    const res = await kuaishouCardApi.getKuaishouCardList(params)
    tableData.value = res.list || [];
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.e05c1ca3')))
    console.error(t('auto.List.80e9795e'), error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  pagination.currentPage = 1
  loadData()
}

const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  searchForm.is_active = null
  handleSearch()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  loadData()
}

const handleCurrentChange = (current) => {
  pagination.currentPage = current
  loadData()
}

const handleAdd = () => {
  currentCardId.value = null
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  currentCardId.value = row.id
  formData.title = row.title
  formData.description = row.description
  formData.image_url = row.image_url
  formData.redirect_url = row.redirect_url
  formData.tags = row.tags || ''
  formData.is_active = row.is_active
  formData.domain_pool_id = row.domain_pool_id || null
  dialogVisible.value = true
}

const handleStats = (row) => {
  router.push(`/kuaishou-card-stats/${row.id}`)
}

const handleCopyLink = async (row) => {
  try {
    const cardLink = `${window.location.origin}/kuaishou/card/${row.id}`;
    
    await navigator.clipboard.writeText(cardLink)
    ElMessage.success(i18n.global.t(t('auto.List.deb26c26')))
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.da0ab552')))
    console.error(t('auto.List.fbb8cbdf'), error)
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除卡片「${row.title}」吗？`,
      t('auto.List.50eaf94d'),
      {
        confirmButtonText: t('auto.List.38cf16f2'),
        cancelButtonText: t('auto.List.625fb26b'),
        type: 'warning'
      }
    )
    const res = await kuaishouCardApi.deleteKuaishouCard(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.0007d170')));
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
    }
  }
}

const handleStatusChange = async (row) => {
  try {
    const res = await kuaishouCardApi.updateKuaishouCard({
      id: row.id,
      is_active: row.is_active
    });
  } catch (error) {
    row.is_active = !row.is_active
    ElMessage.error(i18n.global.t(t('auto.List.930442e2')))
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    let res
    if (currentCardId.value) {
      res = await kuaishouCardApi.updateKuaishouCard({
        ...formData,
        id: currentCardId.value
      })
    } else {
      res = await kuaishouCardApi.createKuaishouCard(formData)
    }
    ElMessage.success(currentCardId.value ? t('auto.List.55aa6366') : t('auto.List.04a691b3'));
    dialogVisible.value = false
    loadData()
  } catch (error) {
    if (error === false) return
    ElMessage.error(i18n.global.t(t('auto.List.5fa802be')))
  }
}



const handleImageUrlChange = () => {
  if (formRef.value && formData.image_url) {
    formRef.value.validateField('image_url')
  }
};

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  formData.title = ''
  formData.description = ''
  formData.image_url = ''
  formData.redirect_url = ''
  formData.tags = ''
  formData.is_active = true
  formData.domain_pool_id = null
}

const getDomainPoolList = async () => {
  try {
    const res = await domainPoolApi.getList({ page: 1, page_size: 100 })
    domainPoolOptions.value = res.list || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.1be7af8d')))
    console.error(t('auto.List.61abe0da'), error)
  }
};

const goToStats = () => {
  router.push('/kuaishou/stats')
}

onMounted(async () => {
  await getDomainPoolList();
  loadData()
});
</script>

<style scoped>
.kuaishou-card-container {
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

.search-form {
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.title-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.description-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.redirect-url {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.no-image {
  display: inline-block;
  width: 80px;
  height: 60px;
  background-color: #f0f0f0;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 12px;
}

.dialog-content {
  display: flex;
  gap: 20px;
}

.form-section {
  flex: 1;
}

.preview-section {
  width: 300px;
  flex-shrink: 0;
}

.image-preview {
  margin-top: 10px;
}

/* 保留图片样式 */
.image-preview .el-image {
  border: 1px solid #d9d9d9;
}

.short-link-container {
  padding: 20px 0;
}
</style>