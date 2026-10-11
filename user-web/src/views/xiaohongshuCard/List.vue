<template>
  <div class="xiaohongshu-card-container">
    
    <el-card class="list-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('小红书卡片管理') }}</span>
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
      
      
      <el-form :model="searchForm" inline class="search-form">
        <el-form-item :label="$t('auto.List.1.9699a5')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('auto.List.2.cd11be')" clearable />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3.3fea7c')">
          <el-select v-model="searchForm.isActive" :placeholder="$t('auto.List.4.e1c965')" clearable>
            <el-option :label="$t('auto.List.5.7854b5')" :value="true" />
            <el-option :label="$t('auto.List.6.710ad0')" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('auto.List.28.e5f71f') }}</el-button>
          <el-button @click="handleReset">{{ $t('auto.List.29.4b9c32') }}</el-button>
        </el-form-item>
      </el-form>
      
      
      <el-table :data="tableData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('auto.List.7.32c65d')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="description" :label="$t('auto.List.8.3bdd08')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.9.20def7')" width="120">
          <template #default="scope">
            <el-image
              v-if="scope.row.image_url"
              :src="scope.row.image_url"
              @error="onImgError"
              :preview-src-list="[scope.row.image_url]"
              style="width: 80px; height: 60px"
              fit="cover"
            />
            <span v-else>{{ $t('auto.List.30.5423e2') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.10.dec2eb')" width="150">
          <template #default="scope">
            <el-link v-if="scope.row.redirect_url" :href="scope.row.redirect_url" target="_blank" type="primary">
              访问链接
            </el-link>
            <span v-else>{{ $t('auto.List.31.7b8af9') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.11.439a19')" width="100">
          <template #default="scope">
            <div>
              <span>浏览: {{ scope.row.view_count }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.3.3fea7c')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_active ? 'success' : 'danger'">
              {{ scope.row.is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.List.12.eca37c')" width="180" />
        <el-table-column :label="$t('auto.List.13.2b6bc0')" width="280" fixed="right">
            <template #default="scope">
              <el-button size="small" type="primary" @click="handleEdit(scope.row)">{{ $t('auto.List.32.95b351') }}</el-button>
              <el-button size="small" type="info" @click="handleStats(scope.row)">{{ $t('auto.List.33.d7656a') }}</el-button>
              <el-button size="small" type="info" @click="handleCopyLink(scope.row)">{{ $t('auto.List.34.879058') }}</el-button>
              <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.List.35.2f4aad') }}</el-button>
            </template>
          </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.14.21efd8')" />
        </template>
      </el-table>

      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
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
              <el-input v-model="form.title" :placeholder="$t('auto.List.15.96641a')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.8.3bdd08')" prop="description">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="3"
                :placeholder="$t('auto.List.16.11956a')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.List.17.844bf1')" prop="image_url">
              <el-input v-model="form.image_url" :placeholder="$t('auto.List.18.3af187')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.10.dec2eb')" prop="redirect_url">
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
              <el-switch v-model="form.is_active" />
            </el-form-item>
          </el-form>
        </div>
        
        
        <div class="preview-section">
          <XiaohongshuCardPreview 
            :card-data="form" 
            @view="handlePreviewView"
          />
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.36.625fb2') }}</el-button>
          <el-button type="primary" @click="handleSubmit">{{ $t('auto.List.37.38cf16') }}</el-button>
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
        <el-descriptions-item :label="$t('auto.List.25.498128')">{{ $t('auto.List.38.ae26bb') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3.3fea7c')">
          <el-tag :type="currentCard.is_active ? 'success' : 'danger'">
            {{ currentCard.is_active ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.12.eca37c')" :span="2">{{ currentCard.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.10.dec2eb')" :span="2">
          <el-link v-if="currentCard.redirect_url" :href="currentCard.redirect_url" target="_blank" type="primary">
            {{ currentCard.redirect_url }}
          </el-link>
          <span v-else>{{ $t('auto.List.39.d81bb2') }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.26.d098e4')" :span="2">
          <template v-if="currentCard.short_link_url">
            {{ window.location.origin }}{{ currentCard.short_link_url }}
          </template>
          <span v-else>{{ $t('auto.List.39.d81bb2') }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.22.14d342')">{{ currentCard.tags || '无' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.27.f8ff52')">{{ currentCard.view_count }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.9.20def7')" :span="2">
          <el-image
            v-if="currentCard.image_url"
            :src="currentCard.image_url"
            style="width: 200px; height: 200px"
            fit="cover"
            :preview-src-list="[currentCard.image_url]"
          />
          <span v-else>{{ $t('auto.List.30.5423e2') }}</span>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="detailVisible = false">{{ $t('auto.List.40.b15d91') }}</el-button>
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import XiaohongshuCardPreview from '@/components/XiaohongshuCardPreview.vue'
import { domainPoolApi } from '@/api/domainPool'
import { getXiaohongshuCardList, createXiaohongshuCard, updateXiaohongshuCard, deleteXiaohongshuCard, generateShortLink } from '@/api/xiaohongshuCard'

const router = useRouter();

const { t } = useI18n()

const domainPoolOptions = ref([]);

const searchForm = reactive({
  keyword: '',
  isActive: null
});

const tableData = ref([]);
const loading = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const dialogVisible = ref(false);
const dialogTitle = ref('')
const formRef = ref()

const detailVisible = ref(false);
const currentCard = ref(null)

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
    { required: true, message: i18n.global.t(t('auto.List.15.96641a')), trigger: 'blur' },
    { max: 100, message: i18n.global.t(t('auto.List.41.5f27c0')), trigger: 'blur' }
  ],
  description: [
    { required: true, message: i18n.global.t(t('auto.List.16.11956a')), trigger: 'blur' },
    { max: 500, message: i18n.global.t(t('auto.List.42.2092c7')), trigger: 'blur' }
  ],
  image_url: [
    { required: true, message: i18n.global.t(t('auto.List.18.3af187')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.43.ddb160')), trigger: 'blur' }
  ],
  redirect_url: [
    { type: 'url', message: i18n.global.t(t('auto.List.43.ddb160')), trigger: 'blur' }
  ],
  tags: [
    { max: 200, message: i18n.global.t(t('auto.List.44.1c8395')), trigger: 'blur' }
  ]
};



const getDomainPoolList = async () => {
  try {
    const res = await domainPoolApi.getList({ page: 1, page_size: 100 })
    domainPoolOptions.value = res.list || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.45.1be7af')))
    console.error(error)
  }
};

const getList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword,
      is_active: searchForm.isActive
    }
    const res = await getXiaohongshuCardList(params)
    tableData.value = res.list || [];
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.46.fe9d24')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  pagination.page = 1
  getList()
};

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.isActive = null
  pagination.page = 1
  getList()
};

const handleSizeChange = (val) => {
  pagination.pageSize = val
  pagination.page = 1
  getList()
};

const handleCurrentChange = (val) => {
  pagination.page = val
  getList()
};

const handleAdd = () => {
  dialogTitle.value = t('auto.List.47.cf4621')
  dialogVisible.value = true
  resetForm()
};

const handleEdit = (row) => {
  dialogTitle.value = t('auto.List.48.dde3a3')
  dialogVisible.value = true
  Object.assign(form, {
    id: row.id,
    title: row.title,
    description: row.description,
    image_url: row.image_url,
    redirect_url: row.redirect_url,
    tags: row.tags,
    is_active: row.is_active,
    domain_pool_id: row.domain_pool_id
  })
};

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除卡片"${row.title}"吗？`, t('auto.List.49.02d981'), {
    confirmButtonText: t('auto.List.37.38cf16'),
    cancelButtonText: t('auto.List.36.625fb2'),
    type: 'warning'
  })
    .then(async () => {
      try {
    const res = await deleteXiaohongshuCard(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.50.0007d1')));
    getList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.51.acf066')))
    console.error(error)
  }
    })
    .catch(() => {})
};

const handlePreviewView = (cardData) => {
  currentCard.value = cardData
  detailVisible.value = true
};

const handleSubmit = () => {
  formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        let res
        if (form.id) {
          res = await updateXiaohongshuCard(form);
        } else {
          res = await createXiaohongshuCard(form);
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

const resetForm = () => {
  form.id = null
  form.title = ''
  form.description = ''
  form.image_url = ''
  form.redirect_url = 'https://www.xiaohongshu.com'
  form.tags = ''
  form.is_active = true
  form.domain_pool_id = null
};

const handleDialogClose = () => {
  resetForm()
  if (formRef.value) {
    formRef.value.clearValidate()
  }
};

const goToStats = () => {
  router.push('/xiaohongshu/stats')
};

const handleStats = (row) => {
  router.push(`/xiaohongshu-card-stats/${row.id}`)
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
        linkUrl = row.redirect_url || window.location.origin + '/xiaohongshu/card/' + row.id;
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
.xiaohongshu-card-container {
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

.search-card {
  margin-bottom: 20px;
}

.list-card {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.dialog-content {
  display: flex;
  gap: 20px;
}

.form-section {
  flex: 1;
  padding-right: 20px;
  border-right: 1px solid #ebeef5;
}

.preview-section {
  flex: 1;
  padding-left: 20px;
  min-width: 400px;
}
</style>