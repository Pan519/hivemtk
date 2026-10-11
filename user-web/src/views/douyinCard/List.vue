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
        <el-table-column prop="title" :label="$t('auto.List.32c65d8d')" width="180" />
        <el-table-column prop="description" :label="$t('auto.List.3bdd08ad')" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.20def794')" width="120">
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
        <el-table-column prop="view_count" :label="$t('auto.List.f8ff5260')" width="100" />
        <el-table-column prop="is_active" :label="$t('auto.List.3fea7ca7')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_active ? 'success' : 'danger'">
              {{ scope.row.is_active ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.List.eca37cb0')" width="180" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
            <el-button size="small" type="primary" @click="handleStats(scope.row)">{{ $t('auto.List.d7656a65') }}</el-button>
            <el-button size="small" type="info" @click="handleCopyLink(scope.row)">{{ $t('auto.List.879058ce') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.21efd88b')" />
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
            <el-form-item :label="$t('auto.List.32c65d8d')" prop="title">
              <el-input v-model="form.title" :placeholder="$t('auto.List.3fd55a4e')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.3bdd08ad')" prop="description">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="3"
                :placeholder="$t('auto.List.705bdfc9')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.List.844bf1a8')" prop="image_url">
              <el-input v-model="form.image_url" :placeholder="$t('auto.List.3af187ba')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.dec2eb71')" prop="redirect_url">
              <el-input v-model="form.redirect_url" :placeholder="$t('auto.List.c63a5309')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.0efb0786')" prop="domain_pool_id">
              <el-select 
                v-model="form.domain_pool_id" 
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
              <el-input v-model="form.tags" :placeholder="$t('auto.List.dc1fb70f')" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.3fea7ca7')" prop="is_active">
              <el-switch v-model="form.is_active" :active-text="$t('auto.List.83a991d7')" :inactive-text="$t('auto.List.710ad08b')" />
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
          <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
          <el-button type="primary" @click="handleSubmit">{{ $t('auto.List.38cf16f2') }}</el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="detailVisible"
      :title="$t('auto.List.4c4bc19f')"
      width="600px"
    >
      <el-descriptions v-if="currentCard" :column="2" border>
        <el-descriptions-item label="ID">{{ currentCard.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.32c65d8d')">{{ currentCard.title }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3bdd08ad')" :span="2">{{ currentCard.description }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.4981287f')">{{ $t('auto.List.277a984f') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3fea7ca7')">
          <el-tag :type="currentCard.is_active ? 'success' : 'danger'">
            {{ currentCard.is_active ? '激活' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.eca37cb0')" :span="2">{{ currentCard.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.dec2eb71')" :span="2">
          <el-link v-if="currentCard.redirect_url" :href="currentCard.redirect_url" target="_blank" type="primary">
            {{ currentCard.redirect_url }}
          </el-link>
          <span v-else>{{ $t('auto.List.d81bb206') }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.d098e4f7')" :span="2">
          <template v-if="currentCard.short_link_url">
            {{ window.location.origin }}{{ currentCard.short_link_url }}
          </template>
          <span v-else>{{ $t('auto.List.d81bb206') }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.14d34236')">{{ currentCard.tags || '无' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.f8ff5260')">{{ currentCard.view_count }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.20def794')" :span="2">
          <el-image
            v-if="currentCard.image_url"
            :src="currentCard.image_url"
            style="width: 200px; height: 200px"
            fit="cover"
            :preview-src-list="[currentCard.image_url]"
          />
          <span v-else>{{ $t('auto.List.5423e287') }}</span>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="detailVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
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
const dialogTitle = ref(t('auto.List.a5cd9836'))
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
    { required: true, message: i18n.global.t(t('auto.List.3fd55a4e')), trigger: 'blur' },
    { max: 255, message: i18n.global.t(t('auto.List.e35fdd46')), trigger: 'blur' }
  ],
  description: [
    { max: 500, message: i18n.global.t(t('auto.List.2092c7b5')), trigger: 'blur' }
  ],
  image_url: [
    { required: true, message: i18n.global.t(t('auto.List.3af187ba')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.a7f6ca3f')), trigger: 'blur' }
  ],
  redirect_url: [
    { required: true, message: i18n.global.t(t('auto.List.c63a5309')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.a7f6ca3f')), trigger: 'blur' }
  ],
  tags: [
    { max: 255, message: i18n.global.t(t('auto.List.cfd05181')), trigger: 'blur' }
  ]
};

const getDomainPoolList = async () => {
  try {
    const res = await domainPoolApi.getList({ page: 1, page_size: 100 })
    domainPoolOptions.value = res.list || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.1be7af8d')))
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
    ElMessage.error(i18n.global.t(t('auto.List.fe9d24d5')))
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
  dialogTitle.value = t('auto.List.a5cd9836')
  resetForm()
  dialogVisible.value = true
};

const handleEdit = async (row) => {
  dialogTitle.value = t('auto.List.cc8b41e6')
  try {
    const res = await getDouyinCard(row.id)
    Object.assign(form, res);
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.181421dd')))
    console.error(error)
  }
};

const handleDelete = (row) => {
  ElMessageBox.confirm(t('auto.List.5731eff5'), t('auto.List.02d9819d'), {
    confirmButtonText: t('auto.List.38cf16f2'),
    cancelButtonText: t('auto.List.625fb26b'),
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteDouyinCard(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.0007d170')));
      getList()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
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

        ElMessage.success(form.id ? t('auto.List.55aa6366') : t('auto.List.04a691b3'));
        dialogVisible.value = false
        getList()
      } catch (error) {
        ElMessage.error(form.id ? t('auto.List.930442e2') : t('auto.List.a889286a'))
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
        console.error(t('auto.List.49d98a4d'), apiError)
        linkUrl = row.redirect_url || window.location.origin + '/douyin/card/' + row.id;
      }
    }
    
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(linkUrl)
      ElMessage.success(i18n.global.t(t('auto.List.deb26c26')))
    } else {
      const textArea = document.createElement('textarea');
      textArea.value = linkUrl
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      
      try {
        document.execCommand('copy')
        ElMessage.success(i18n.global.t(t('auto.List.deb26c26')))
      } catch (copyError) {
        ElMessage.error(i18n.global.t(t('auto.List.9fb49918')))
        console.error(t('auto.List.b3895887'), copyError)
      } finally {
        document.body.removeChild(textArea)
      }
    }
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.da0ab552')))
    console.error(t('auto.List.0109739b'), error)
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