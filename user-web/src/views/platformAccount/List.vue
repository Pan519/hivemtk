<template>
  <div class="platform-account-container">
    
    <div class="page-header">
      <h2>{{ $t('平台账号管理') }}</h2>
      <div class="action-buttons">
        <el-button type="info" @click="handleViewPlatforms">
          <el-icon><Platform /></el-icon>
          {{ $t('支持平台') }}
        </el-button>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          {{ $t('新增账号') }}
        </el-button>
      </div>
    </div>

    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="search-form-content">
        <el-form-item ::label="$t('auto.List.1.80a5d0')">
          <el-select v-model="searchForm.platform" ::placeholder="$t('auto.List.2.9b247d')" clearable>
            <el-option
              v-for="platform in platformOptions"
              :key="platform.value"
              :label="platform.label"
              :value="platform.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.3.4a56d3')">
          <el-input v-model="searchForm.account" ::placeholder="$t('auto.List.4.5a2cef')" clearable />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.5.1b8578')">
          <el-select v-model="searchForm.status" ::placeholder="$t('auto.List.6.333291')" clearable>
            <el-option ::label="$t('auto.List.7.62f11e')" value="1" />
            <el-option ::label="$t('auto.List.8.38b899')" value="2" />
            <el-option ::label="$t('auto.List.9.c66e0a')" value="3" />
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

    
    <el-table :data="accountList" border style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="platform" ::label="$t('auto.List.1.80a5d0')" width="120" />
      <el-table-column prop="account_name" ::label="$t('auto.List.10.5f0720')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="account_id" ::label="$t('auto.List.11.61c2f5')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="status" ::label="$t('auto.List.5.1b8578')" width="100">
        <template #default="scope">
          <el-tag :type="getStatusType(scope.row.status)">
            {{ getStatusText(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_login_at" :label="$t('auto.List.12.19fb21')" width="180" />
      <el-table-column prop="created_at" :label="$t('auto.List.13.eca37c')" width="180" />
      <el-table-column :label="$t('auto.List.14.2b6bc0')" width="280" fixed="right">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-button type="success" size="small" @click="handleLogin(scope.row)">
            <el-icon><SwitchButton /></el-icon>
            登录
          </el-button>
          <el-button type="info" size="small" @click="handleCheckStatus(scope.row)">
            <el-icon><Check /></el-icon>
            检查
          </el-button>
          <el-button type="danger" size="small" @click="handleDelete(scope.row)">
            <el-icon><Delete /></el-icon>
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
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item :label="$t('auto.List.15.498128')" prop="platform">
          <el-select v-model="form.platform" :placeholder="$t('auto.List.16.72b562')" style="width: 100%">
            <el-option
              v-for="platform in platformOptions"
              :key="platform.value"
              :label="platform.label"
              :value="platform.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.17.17327f')" prop="account_name">
          <el-input v-model="form.account_name" :placeholder="$t('auto.List.18.1f3797')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.19.14d171')" prop="account_id">
          <el-input v-model="form.account_id" :placeholder="$t('auto.List.20.cf5e16')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.21.ecf09a')" prop="credentials">
          <el-input
            v-model="form.credentials"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.List.22.763626')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.23.3fea7c')" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">{{ $t('auto.List.34.fd6e80') }}</el-radio>
            <el-radio :label="2">{{ $t('auto.List.35.c195df') }}</el-radio>
            <el-radio :label="3">{{ $t('auto.List.36.63e85d') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.37.625fb2') }}</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="platformsDialogVisible"
      :title="$t('auto.List.24.39e4a5')"
      width="600px"
    >
      <div v-loading="platformsLoading">
        <el-table :data="platformList" border style="width: 100%">
          <el-table-column prop="code" :label="$t('auto.List.25.cb4e0f')" width="150" />
          <el-table-column prop="name" :label="$t('auto.List.26.00c1f0')" min-width="150" />
          <el-table-column prop="description" :label="$t('auto.List.27.f411d0')" min-width="200" show-overflow-tooltip />
        </el-table>
      </div>
    </el-dialog>

    
    <el-dialog
      v-model="loginDialogVisible"
      :title="$t('auto.List.28.efae77')"
      width="500px"
      @close="handleLoginDialogClose"
    >
      <div v-loading="loginLoading">
        <el-form :model="loginForm" label-width="100px">
          <el-form-item :label="$t('auto.List.29.7035c6')">
            <span>{{ currentAccount.account_name }}</span>
          </el-form-item>
          <el-form-item :label="$t('auto.List.15.498128')">
            <span>{{ currentAccount.platform }}</span>
          </el-form-item>
          <el-form-item :label="$t('auto.List.30.983f59')">
            <el-input v-model="loginForm.code" :placeholder="$t('auto.List.31.09998b')" />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="loginDialogVisible = false">{{ $t('auto.List.37.625fb2') }}</el-button>
          <el-button type="primary" @click="submitLogin" :loading="loginLoading">
            确认登录
          </el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="statusDialogVisible"
      :title="$t('auto.List.32.d05d90')"
      width="500px"
    >
      <div v-loading="statusLoading">
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="$t('auto.List.29.7035c6')">{{ currentAccount.account_name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.15.498128')">{{ currentAccount.platform }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.23.3fea7c')">
            <el-tag :type="getStatusType(statusResult.status)">
              {{ getStatusText(statusResult.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.33.44c422')">{{ statusResult.checked_at }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.27.f411d0')">{{ statusResult.message || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Search, RefreshRight, Edit, Delete, Platform, SwitchButton, Check
} from '@element-plus/icons-vue'
import { platformAccountApi } from '@/api/platformAccount'
import { getAccountStatusLabel, getAccountStatusTagType } from '@/constants/accountType';

const { t } = useI18n()

const loading = ref(false);
const submitting = ref(false)
const platformsLoading = ref(false)
const loginLoading = ref(false)
const statusLoading = ref(false)
const accountList = ref([])
const platformList = ref([])
const platformOptions = ref([])
const dialogVisible = ref(false)
const platformsDialogVisible = ref(false)
const loginDialogVisible = ref(false)
const statusDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const currentAccount = ref({})
const statusResult = ref({})

const searchForm = reactive({
  platform: '',
  account: '',
  status: ''
});

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const form = reactive({
  id: null,
  platform: '',
  account_name: '',
  account_id: '',
  credentials: '',
  status: 1
});

const loginForm = reactive({
  code: ''
});

const rules = {
  platform: [
    { required: true, message: i18n.global.t(t('auto.List.16.72b562')), trigger: 'change' }
  ],
  account_name: [
    { required: true, message: i18n.global.t(t('auto.List.18.1f3797')), trigger: 'blur' }
  ],
  account_id: [
    { required: true, message: i18n.global.t(t('auto.List.20.cf5e16')), trigger: 'blur' }
  ]
};

const dialogTitle = computed(() => {
  return isEdit.value ? t('auto.List.38.656e89') : t('auto.List.39.34374b')
});

onMounted(() => {
  fetchAccountList()
  fetchPlatforms()
});

const fetchAccountList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      platform: searchForm.platform,
      account: searchForm.account,
      status: searchForm.status ? parseInt(searchForm.status) : undefined
    }
    const res = await platformAccountApi.getAccounts(params)
    accountList.value = res || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  pagination.page = 1
  fetchAccountList()
}

const resetSearch = () => {
  searchForm.platform = ''
  searchForm.account = ''
  searchForm.status = ''
  pagination.page = 1
  fetchAccountList()
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchAccountList()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchAccountList()
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
  form.platform = row.platform
  form.account_name = row.account_name
  form.account_id = row.account_id
  form.credentials = typeof row.credentials === 'string' ? row.credentials : JSON.stringify(row.credentials || {})
  form.status = row.status
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除平台账号 "${row.account_name}" 吗？`,
    t('auto.List.40.02d981'),
    {
      confirmButtonText: t('auto.List.41.38cf16'),
      cancelButtonText: t('auto.List.37.625fb2'),
      type: 'warning'
    }
  ).then(async () => {
    try {
      await platformAccountApi.deleteAccount(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.42.0007d1')))
      fetchAccountList()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.43.acf066')))
      console.error(error)
    }
  }).catch((e) => {
    if (e !== 'cancel' && e !== 'close')
      throw e;
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const submitData = { ...form }
        if (submitData.credentials) {
          try {
            submitData.credentials = JSON.parse(submitData.credentials)
          } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
        }
        if (isEdit.value) {
          await platformAccountApi.updateAccount(form.id, submitData)
        } else {
          await platformAccountApi.createAccount(submitData)
        }
        ElMessage.success(isEdit.value ? t('auto.List.44.55aa63') : t('auto.List.45.04a691'))
        dialogVisible.value = false
        fetchAccountList()
      } catch (error) {
        ElMessage.error(isEdit.value ? t('auto.List.46.930442') : t('auto.List.47.a88928'))
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
  form.platform = ''
  form.account_name = ''
  form.account_id = ''
  form.credentials = ''
  form.status = 1
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const fetchPlatforms = async () => {
  try {
    const res = await platformAccountApi.getPlatforms()
    platformList.value = res.list || res || []
    platformOptions.value = platformList.value.map(item => ({
      label: item.name || item.code,
      value: item.code
    }))
  } catch (error) {
    console.error(error)
  }
};

const handleViewPlatforms = () => {
  platformsDialogVisible.value = true
  fetchPlatforms()
}

const handleLogin = (row) => {
  currentAccount.value = row
  loginDialogVisible.value = true
  loginForm.code = ''
};

const handleLoginDialogClose = () => {
  loginForm.code = ''
}

const submitLogin = async () => {
  loginLoading.value = true
  try {
    await platformAccountApi.loginAccount(currentAccount.value.id, loginForm)
    ElMessage.success(i18n.global.t(t('auto.List.48.71fa3b')))
    loginDialogVisible.value = false
    fetchAccountList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.49.b6076a')))
    console.error(error)
  } finally {
    loginLoading.value = false
  }
}

const handleCheckStatus = async (row) => {
  currentAccount.value = row
  statusDialogVisible.value = true
  statusLoading.value = true
  try {
    const res = await platformAccountApi.checkStatus(row.id)
    statusResult.value = res || {}
  } catch (error) {
    console.error(error)
  } finally {
    statusLoading.value = false
  }
};

const getStatusType = (status) => getAccountStatusTagType(status);
const getStatusText = (status) => getAccountStatusLabel(status)
</script>

<style lang="scss" scoped>
.platform-account-container {
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
