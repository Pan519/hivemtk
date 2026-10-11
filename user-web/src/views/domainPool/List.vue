<template>
  <div class="domain-pool-container">
    
    <div class="page-header">
      <h2>{{ $t('域名池管理') }}</h2>
      <div class="action-buttons">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          {{ $t('添加域名') }}
        </el-button>
        <el-button type="success" @click="handleCheckAll" :loading="checkingAll">
          <el-icon><Refresh /></el-icon>
          {{ $t('检查所有域名') }}
        </el-button>
      </div>
    </div>

    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="search-form-content">
        <el-form-item ::label="$t('auto.List.1.8abfb3')">
          <el-input v-model="searchForm.domain" ::placeholder="$t('auto.List.2.c52926')" clearable />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.3.1b8578')">
          <el-select v-model="searchForm.status" ::placeholder="$t('auto.List.4.333291')" clearable>
            <el-option ::label="$t('auto.List.5.62f11e')" value="1" />
            <el-option ::label="$t('auto.List.6.8a9c48')" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ $t('搜索') }}
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshRight /></el-icon>
            {{ $t('重置') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    
    <el-table :data="domainList" border style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="domain" ::label="$t('auto.List.1.8abfb3')" min-width="200" />
      <el-table-column prop="port" ::label="$t('auto.List.7.2696c1')" width="100" />
      <el-table-column prop="purpose" ::label="$t('auto.List.8.c55de6')" min-width="150" />
      <el-table-column prop="status_str" ::label="$t('auto.List.3.1b8578')" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
            {{ scope.row.status_str }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_check" :label="$t('auto.List.9.cbe6d4')" width="180" />
      <el-table-column prop="created_at" :label="$t('auto.List.10.eca37c')" width="180" />
      <el-table-column :label="$t('auto.List.11.2b6bc0')" width="200" fixed="right">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">
            编辑
          </el-button>
          <el-button type="success" size="small" @click="handleCheck(scope.row)" :loading="scope.row.checking">
            检查
          </el-button>
          <el-button type="danger" size="small" @click="handleDelete(scope.row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item :label="$t('auto.List.12.190980')" prop="domain">
          <el-input v-model="form.domain" :placeholder="$t('auto.List.13.a6e880')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.14.c76cfe')" prop="port">
          <el-input-number v-model="form.port" :min="1" :max="65535" :placeholder="$t('auto.List.15.bc0eb7')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.16.67044a')" prop="purpose">
          <el-input v-model="form.purpose" type="textarea" :placeholder="$t('auto.List.17.0aa798')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.18.3fea7c')" prop="status" v-if="isEdit">
          <el-select v-model="form.status" :placeholder="$t('auto.List.19.e1c965')">
            <el-option :label="$t('auto.List.20.fd6e80')" :value="1" />
            <el-option :label="$t('auto.List.21.2329a0')" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.22.625fb2') }}</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Search, RefreshRight } from '@element-plus/icons-vue'
import { domainPoolApi } from '@/api/domainPool'

const { t } = useI18n()

const loading = ref(false);
const checkingAll = ref(false)
const domainList = ref([])
const dialogVisible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const searchForm = reactive({
  domain: '',
  status: ''
});

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const form = reactive({
  id: null,
  domain: '',
  port: 80,
  purpose: '',
  status: 1
});

const rules = {
  domain: [
    { required: true, message: i18n.global.t(t('auto.List.13.a6e880')), trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/, message: i18n.global.t(t('auto.List.23.1f5304')), trigger: 'blur' }
  ],
  port: [
    { required: true, message: i18n.global.t(t('auto.List.15.bc0eb7')), trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: i18n.global.t(t('auto.List.24.b6e86d')), trigger: 'blur' }
  ]
};

const dialogTitle = computed(() => {
  return isEdit.value ? t('auto.List.25.b99033') : t('auto.List.26.f605b5')
});

onMounted(() => {
  fetchDomainList()
});

const fetchDomainList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      domain: searchForm.domain,
      status: searchForm.status ? parseInt(searchForm.status) : undefined
    }
    const res = await domainPoolApi.getList(params)
    domainList.value = res.list
    pagination.total = res.total
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.27.3ceb0b')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  pagination.page = 1
  fetchDomainList()
}

const resetSearch = () => {
  searchForm.domain = ''
  searchForm.status = ''
  pagination.page = 1
  fetchDomainList()
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchDomainList()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchDomainList()
}

const handleAdd = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  form.id = row.id
  form.domain = row.domain
  form.port = row.port
  form.purpose = row.purpose
  form.status = row.status
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除域名 "${row.domain}" 吗？`,
    t('auto.List.28.02d981'),
    {
      confirmButtonText: t('auto.List.29.38cf16'),
      cancelButtonText: t('auto.List.22.625fb2'),
      type: 'warning'
    }
  ).then(async () => {
    try {
      await domainPoolApi.delete(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.30.0007d1')))
      fetchDomainList()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.31.acf066')))
      console.error(error)
    }
  }).catch((e) => {
    if (e !== 'cancel' && e !== 'close')
      throw e;
  })
}

const handleCheck = async (row) => {
  row.checking = true;
  try {
    const res = await domainPoolApi.checkDomain(row.id)
    ElMessage.success(`检查完成：${res.msg || 'OK'}`)
    fetchDomainList();
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.32.74bb5c')))
    console.error(error)
  } finally {
    row.checking = false
  }
}

const handleCheckAll = async () => {
  checkingAll.value = true
  try {
    await domainPoolApi.checkAllDomains()
    ElMessage.success(i18n.global.t(t('auto.List.33.6d24e6')))
    fetchDomainList();
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.32.74bb5c')))
    console.error(error)
  } finally {
    checkingAll.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          await domainPoolApi.update(form)
          ElMessage.success(i18n.global.t(t('auto.List.34.55aa63')))
        } else {
          await domainPoolApi.create(form)
          ElMessage.success(i18n.global.t(t('auto.List.35.04a691')))
        }
        
        dialogVisible.value = false
        fetchDomainList()
      } catch (error) {
        ElMessage.error(isEdit.value ? t('auto.List.36.930442') : t('auto.List.37.a88928'))
        console.error(error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDialogClose = () => {
  resetForm()
}

const resetForm = () => {
  form.id = null
  form.domain = ''
  form.port = 80
  form.purpose = ''
  form.status = 1
  if (formRef.value) {
    formRef.value.resetFields()
  }
}
</script>

<style lang="scss" scoped>
.domain-pool-container {
  padding: 20px;
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      font-size: 24px;
      color: #303133;
    }
    
    .action-buttons {
      display: flex;
      gap: 10px;
    }
  }
  
  .search-form {
    margin-bottom: 20px;
    padding: 15px;
    background-color: #f5f7fa;
    border-radius: 4px;
    
    .search-form-content {
      margin: 0;
    }
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
}
</style>