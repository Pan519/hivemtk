<template>
  <div class="douyin-card-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('抖音卡片管理') }}</span>
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
        <el-form-item :label="$t('auto.List.1.9699a5')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('auto.List.2.976111')" clearable />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3.3fea7c')">
          <el-select v-model="searchForm.is_active" :placeholder="$t('auto.List.4.708c9d')" clearable>
            <el-option :label="$t('auto.List.5.83a991')" :value="true" />
            <el-option :label="$t('auto.List.6.710ad0')" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('auto.List.27.e5f71f') }}</el-button>
          <el-button @click="resetSearch">{{ $t('auto.List.28.4b9c32') }}</el-button>
        </el-form-item>
      </el-form>

      
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('auto.List.7.32c65d')" width="180" />
        <el-table-column prop="description" :label="$t('auto.List.8.3bdd08')" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.9.20def7')" width="120">
          <template #default="scope">
            <el-image 
              style="width: 100px; height: 100px" 
              :src="scope.row.image_url"
              @error="onImgError" 
              :preview-src-list="[scope.row.image_url]"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column prop="view_count" :label="$t('auto.List.10.f8ff52')" width="100" />
        <el-table-column prop="is_active" :label="$t('auto.List.3.3fea7c')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_active ? 'success' : 'danger'">
              {{ scope.row.is_active ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.List.11.eca37c')" width="180" />
        <el-table-column :label="$t('auto.List.12.2b6bc0')" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.List.29.95b351') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.List.30.2f4aad') }}</el-button>
            <el-button size="small" type="primary" @click="handleStats(scope.row)">{{ $t('auto.List.31.d7656a') }}</el-button>
            <el-button size="small" type="info" @click="handleCopyLink(scope.row)">{{ $t('auto.List.32.879058') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.13.21efd8')" />
        </template>
      </el-table>

      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
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
      @close="handleDialogClose"
    >
      <div class="dialog-content">
        
        <div class="form-section">
          <el-form
            ref="formRef"
            :model="form"
            :rules="rules"
            label-width="100px"
          >
            <el-form-item :label="$t('auto.List.7.32c65d')" prop="title">
              <el-input v-model="form.title" :placeholder="$t('auto.List.14.3fd55a')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.8.3bdd08')" prop="description">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="3"
                :placeholder="$t('auto.List.15.705bdf')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.List.16.844bf1')" prop="image_url">
              <el-input v-model="form.image_url" :placeholder="$t('auto.List.17.3af187')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.18.dec2eb')" prop="redirect_url">
              <el-input v-model="form.redirect_url" :placeholder="$t('auto.List.19.c63a53')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.20.0efb07')" prop="domain_pool_id">
              <el-select 
                v-model="form.domain_pool_id" 
                :placeholder="$t('auto.List.21.8ee08e')" 
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
            <el-form-item :label="$t('auto.List.22.14d342')" prop="tags">
              <el-input v-model="form.tags" :placeholder="$t('auto.List.23.dc1fb7')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.3.3fea7c')" prop="is_active">
              <el-switch v-model="form.is_active" :active-text="$t('auto.List.5.83a991')" :inactive-text="$t('auto.List.6.710ad0')" />
            </el-form-item>
          </el-form>
        </div>
        
        
        <div class="preview-section">
          <DouyinCardPreview 
          :card-data="form" 
          @view="handlePreviewView"
        />
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.33.625fb2') }}</el-button>
          <el-button type="primary" @click="handleSubmit">{{ $t('auto.List.34.38cf16') }}</el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="detailVisible"
      :title="$t('auto.List.24.4c4bc1')"
      width="600px"
    >
      <el-descriptions v-if="currentCard" :column="2" border>
        <el-descriptions-item label="ID">{{ currentCard.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.7.32c65d')">{{ currentCard.title }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.8.3bdd08')" :span="2">{{ currentCard.description }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.25.498128')">{{ $t('auto.List.35.277a98') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3.3fea7c')">
          <el-tag :type="currentCard.is_active ? 'success' : 'danger'">
            {{ currentCard.is_active ? '激活' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.11.eca37c')" :span="2">{{ currentCard.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.18.dec2eb')" :span="2">
          <el-link v-if="currentCard.redirect_url" :href="currentCard.redirect_url" target="_blank" type="primary">
            {{ currentCard.redirect_url }}
          </el-link>
          <span v-else>{{ $t('auto.List.36.d81bb2') }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.26.d098e4')" :span="2">
          <template v-if="currentCard.short_link_url">
            {{ window.location.origin }}{{ currentCard.short_link_url }}
          </template>
          <span v-else>{{ $t('auto.List.36.d81bb2') }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.22.14d342')">{{ currentCard.tags || '无' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.10.f8ff52')">{{ currentCard.view_count }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.9.20def7')" :span="2">
          <el-image
            v-if="currentCard.image_url"
            :src="currentCard.image_url"
            style="width: 200px; height: 200px"
            fit="cover"
            :preview-src-list="[currentCard.image_url]"
          />
          <span v-else>{{ $t('auto.List.37.5423e2') }}</span>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="detailVisible = false">{{ $t('auto.List.38.b15d91') }}</el-button>
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

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import DouyinCardPreview from '@/components/DouyinCardPreview.vue'
import { domainPoolApi } from '@/api/domainPool'
import {
  getDouyinCardList,
  getDouyinCard,
  createDouyinCard,
  updateDouyinCard,
  deleteDouyinCard,
  generateShortLink
} from '@/api/douyinCard'

const { t } = useI18n()

const loading = ref(false);
const tableData = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref(t('auto.List.39.a5cd98'))
const formRef = ref(null)
const router = useRouter()

const detailVisible = ref(false);
const currentCard = ref(null)

const domainPoolOptions = ref([]);

const searchForm = reactive({
  keyword: '',
  is_active: null
});

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
});

const form = reactive({
  id: null,
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
    { required: true, message: i18n.global.t(t('auto.List.14.3fd55a')), trigger: 'blur' },
    { max: 255, message: i18n.global.t(t('auto.List.40.e35fdd')), trigger: 'blur' }
  ],
  description: [
    { max: 500, message: i18n.global.t(t('auto.List.41.2092c7')), trigger: 'blur' }
  ],
  image_url: [
    { required: true, message: i18n.global.t(t('auto.List.17.3af187')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.42.a7f6ca')), trigger: 'blur' }
  ],
  redirect_url: [
    { required: true, message: i18n.global.t(t('auto.List.19.c63a53')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.42.a7f6ca')), trigger: 'blur' }
  ],
  tags: [
    { max: 255, message: i18n.global.t(t('auto.List.43.cfd051')), trigger: 'blur' }
  ]
};

const getDomainPoolList = async () => {
  try {
    const res = await domainPoolApi.getList({ page: 1, page_size: 100 })
    domainPoolOptions.value = res.list || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.44.1be7af')))
    console.error(error)
  }
};

const getList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: searchForm.keyword,
      is_active: searchForm.is_active
    }
    const res = await getDouyinCardList(params)
    tableData.value = res.list || [];
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.45.fe9d24')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  pagination.page = 1
  getList()
};

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.is_active = null
  pagination.page = 1
  getList()
};

const handleSizeChange = (size) => {
  pagination.page_size = size
  getList()
};

const handleCurrentChange = (page) => {
  pagination.page = page
  getList()
};

const handleAdd = () => {
  dialogTitle.value = t('auto.List.39.a5cd98')
  resetForm()
  dialogVisible.value = true
};

const handleEdit = async (row) => {
  dialogTitle.value = t('auto.List.46.cc8b41')
  try {
    const res = await getDouyinCard(row.id)
    Object.assign(form, res);
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.47.181421')))
    console.error(error)
  }
};

const handleDelete = (row) => {
  ElMessageBox.confirm(t('auto.List.48.5731ef'), t('auto.List.49.02d981'), {
    confirmButtonText: t('auto.List.34.38cf16'),
    cancelButtonText: t('auto.List.33.625fb2'),
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteDouyinCard(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.50.0007d1')));
      getList()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.51.acf066')))
      console.error(error)
    }
  }).catch(() => {})
};

const handleSubmit = () => {
  formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        let res
        if (form.id) {
          res = await updateDouyinCard(form);
        } else {
          res = await createDouyinCard(form);
        }

        ElMessage.success(form.id ? t('auto.List.52.55aa63') : t('auto.List.53.04a691'));
        dialogVisible.value = false
        getList()
      } catch (error) {
        ElMessage.error(form.id ? t('auto.List.54.930442') : t('auto.List.55.a88928'))
        console.error(error)
      }
    }
  })
};

const handleDialogClose = () => {
  resetForm()
};

const resetForm = () => {
  form.id = null
  form.title = ''
  form.description = ''
  form.image_url = ''
  form.redirect_url = 'https://www.douyin.com'
  form.tags = ''
  form.is_active = true
  form.domain_pool_id = null
  if (formRef.value) {
    formRef.value.resetFields()
  }
};

const handlePreviewView = (cardData) => {
  currentCard.value = cardData
  detailVisible.value = true
};

const handleStats = (row) => {
  router.push(`/douyin-card-stats/${row.id}`)
};

const goToStats = () => {
  router.push('/douyin/stats')
};

const handleCopyLink = async (row) => {
  try {
    let linkUrl = '';
    if (row.short_link_url) {
      linkUrl = window.location.origin + row.short_link_url;
    } else {
      try {
        const res = await generateShortLink(row.id)
        linkUrl = window.location.origin + res.short_link_url;
        row.short_link_url = res.short_link_url;
        row.short_code = res.short_code
      } catch (apiError) {
        console.error(t('auto.List.56.49d98a'), apiError)
        linkUrl = row.redirect_url || window.location.origin + '/douyin/card/' + row.id;
      }
    }
    
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(linkUrl)
      ElMessage.success(i18n.global.t(t('auto.List.57.deb26c')))
    } else {
      const textArea = document.createElement('textarea');
      textArea.value = linkUrl
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      
      try {
        document.execCommand('copy')
        ElMessage.success(i18n.global.t(t('auto.List.57.deb26c')))
      } catch (copyError) {
        ElMessage.error(i18n.global.t(t('auto.List.58.9fb499')))
        console.error(t('auto.List.59.b38958'), copyError)
      } finally {
        document.body.removeChild(textArea)
      }
    }
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.60.da0ab5')))
    console.error(t('auto.List.61.010973'), error)
  }
};

onMounted(() => {
  getList()
  getDomainPoolList()
});
</script>

<style scoped>
.douyin-card-container {
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
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.dialog-content {
  display: flex;
  gap: 20px;
}

.form-section {
  flex: 1;
  max-width: 500px;
}

.preview-section {
  flex: 1;
  max-width: 400px;
  border-left: 1px solid #ebeef5;
  padding-left: 20px;
}
</style>