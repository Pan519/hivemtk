<template>
  <div class="xianyu-card-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('闲鱼卡片管理') }}</span>
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
          <el-button type="primary" @click="handleSearch">{{ $t('auto.List.23.e5f71f') }}</el-button>
          <el-button @click="resetSearch">{{ $t('auto.List.24.4b9c32') }}</el-button>
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
        <el-table-column prop="click_count" :label="$t('auto.List.11.3824fe')" width="100" />
        <el-table-column prop="share_count" :label="$t('auto.List.12.70a39a')" width="100" />
        <el-table-column :label="$t('auto.List.3.3fea7c')" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.is_active ? 'success' : 'danger'">
              {{ scope.row.is_active ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.13.2b6bc0')" width="280" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.List.25.95b351') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.List.26.2f4aad') }}</el-button>
            <el-button size="small" type="primary" @click="handleStats(scope.row)">{{ $t('auto.List.27.d7656a') }}</el-button>
            <el-button size="small" type="info" @click="handleCopyLink(scope.row)">{{ $t('auto.List.28.879058') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item :label="$t('auto.List.7.32c65d')" prop="title">
          <el-input v-model="form.title" :placeholder="$t('auto.List.14.96641a')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.8.3bdd08')" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.List.15.11956a')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.9.20def7')" prop="image_url">
          <el-input v-model="form.image_url" :placeholder="$t('auto.List.16.3af187')">
            <template #append>
              <el-button @click="selectImage">{{ $t('auto.List.29.153fa6') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('auto.List.17.dec2eb')" prop="redirect_url">
          <el-input v-model="form.redirect_url" :placeholder="$t('auto.List.18.c63a53')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.19.0efb07')" prop="domain_pool_id">
          <el-select 
            v-model="form.domain_pool_id" 
            :placeholder="$t('auto.List.20.8ee08e')" 
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
        <el-form-item :label="$t('auto.List.21.14d342')" prop="tags">
          <el-input v-model="form.tags" :placeholder="$t('auto.List.22.dc1fb7')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3.3fea7c')" prop="is_active">
          <el-switch v-model="form.is_active" :active-text="$t('auto.List.5.83a991')" :inactive-text="$t('auto.List.6.710ad0')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.30.625fb2') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ $t('auto.List.31.38cf16') }}</el-button>
        </span>
      </template>
    </el-dialog>

    
    <MaterialSelectDialog
      v-model="imageDialogVisible"
      @confirm="handleImageSelected"
    />
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
import MaterialSelectDialog from '@/components/dialogs/MaterialSelectDialog.vue'
import { domainPoolApi } from '@/api/domainPool'
import {
  getXianyuCardList,
  getXianyuCard,
  createXianyuCard,
  updateXianyuCard,
  deleteXianyuCard,
  generateXianyuShortLink
} from '@/api/xianyuCard'

const router = useRouter()

const { t } = useI18n()

const tableData = ref([]);
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const domainPoolOptions = ref([]);

const searchForm = reactive({
  keyword: '',
  is_active: null
});

const dialogVisible = ref(false);
const dialogTitle = ref('')
const formRef = ref()
const form = reactive({
  id: null,
  title: '',
  description: '',
  image_url: '',
  redirect_url: '',
  tags: '',
  is_active: true,
  domain_pool_id: null
})

const imageDialogVisible = ref(false);

const rules = {
  title: [
    { required: true, message: i18n.global.t(t('auto.List.32.3fd55a')), trigger: 'blur' },
    { max: 255, message: i18n.global.t(t('auto.List.33.e35fdd')), trigger: 'blur' }
  ],
  description: [
    { required: true, message: i18n.global.t(t('auto.List.34.705bdf')), trigger: 'blur' },
    { max: 500, message: i18n.global.t(t('auto.List.35.2092c7')), trigger: 'blur' }
  ],
  image_url: [
    { required: true, message: i18n.global.t(t('auto.List.16.3af187')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.36.a7f6ca')), trigger: 'blur' }
  ],
  redirect_url: [
    { required: true, message: i18n.global.t(t('auto.List.18.c63a53')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.36.a7f6ca')), trigger: 'blur' }
  ],
  tags: [
    { max: 255, message: i18n.global.t(t('auto.List.37.cfd051')), trigger: 'blur' }
  ]
};

const getDomainPoolList = async () => {
  try {
    const res = await domainPoolApi.getList({ page: 1, page_size: 100 })
    domainPoolOptions.value = res.list || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.38.1be7af')))
    console.error(error)
  }
};

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchForm.keyword,
      is_active: searchForm.is_active
    }
    const res = await getXianyuCardList(params)
    tableData.value = res.list || [];
    total.value = res.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.39.fe9d24')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  currentPage.value = 1
  fetchData()
};

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.is_active = null
  handleSearch()
};

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchData()
};

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchData()
}

const resetForm = () => {
  form.id = null
  form.title = ''
  form.description = ''
  form.image_url = ''
  form.redirect_url = 'https://www.xianyu.com'
  form.tags = ''
  form.is_active = true
  form.domain_pool_id = null
  if (formRef.value) {
    formRef.value.resetFields()
  }
};

const handleAdd = () => {
  dialogTitle.value = t('auto.List.40.5c2e85')
  resetForm()
  dialogVisible.value = true
};

const handleEdit = async (row) => {
  dialogTitle.value = t('auto.List.41.b3ddda')
  try {
    const res = await getXianyuCard(row.id)
    Object.assign(form, res);
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.42.181421')))
    console.error(error)
  }
};

const submitForm = async () => {
  formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        let res
        if (form.id) {
          res = await updateXianyuCard(form);
        } else {
          res = await createXianyuCard(form);
        }

        ElMessage.success(form.id ? t('auto.List.43.55aa63') : t('auto.List.44.04a691'));
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        ElMessage.error(form.id ? t('auto.List.45.930442') : t('auto.List.46.a88928'))
        console.error(error)
      }
    }
  })
};

const handleDelete = (row) => {
  ElMessageBox.confirm(t('auto.List.47.383acb'), t('auto.List.48.02d981'), {
    confirmButtonText: t('auto.List.31.38cf16'),
    cancelButtonText: t('auto.List.30.625fb2'),
    type: 'warning'
  }).then(async () => {
    try {
      await deleteXianyuCard(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.49.0007d1')))
      fetchData()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.50.acf066')))
    }
  }).catch(() => {})
};

const handleCopyLink = async (row) => {
  try {
    let linkUrl = '';
    if (row.short_link_url) {
      linkUrl = window.location.origin + row.short_link_url;
    } else {
      try {
        const res = await generateXianyuShortLink(row.id)
        linkUrl = window.location.origin + res.short_link_url;
        row.short_link_url = res.short_link_url;
        row.short_code = res.short_code
      } catch (apiError) {
        console.error(t('auto.List.51.49d98a'), apiError)
        linkUrl = row.redirect_url || window.location.origin + '/xianyu/card/' + row.id;
      }
    }
    
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(linkUrl)
      ElMessage.success(i18n.global.t(t('auto.List.52.deb26c')))
    } else {
      const textArea = document.createElement('textarea');
      textArea.value = linkUrl
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      
      try {
        document.execCommand('copy')
        ElMessage.success(i18n.global.t(t('auto.List.52.deb26c')))
      } catch (copyError) {
        ElMessage.error(i18n.global.t(t('auto.List.53.9fb499')))
        console.error(t('auto.List.54.b38958'), copyError)
      } finally {
        document.body.removeChild(textArea)
      }
    }
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.55.da0ab5')))
    console.error(t('auto.List.56.010973'), error)
  }
};

const handleStats = (row) => {
  router.push(`/xianyu-card-stats/${row.id}`)
};

const goToStats = () => {
  router.push('/xianyu/stats')
};

const selectImage = () => {
  imageDialogVisible.value = true
};

const handleImageSelected = (material) => {
  form.image_url = material.url
  imageDialogVisible.value = false
};

onMounted(() => {
  fetchData()
  getDomainPoolList()
})
</script>

<style scoped>
.xianyu-card-container {
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
</style>