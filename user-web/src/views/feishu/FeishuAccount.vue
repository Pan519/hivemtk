<template>
  <div class="feishu-account-management">
    <div class="account-tabs">
      <div class="account-list-container">
        <div class="account-search">
          <el-input
            v-model="searchKeyword"
            ::placeholder="$t('auto.FeishuAccount.1.9d9a5f')"
            prefix-icon="Search"
            clearable
            @clear="handleSearch"
            @input="handleSearch"
          />
          <div class="action-buttons">
            <el-button type="primary" @click="handleAdd">{{ $t('添加飞书账号') }}</el-button>
            <el-button type="info" :loading="loading" @click="fetchAccounts">
              <el-icon><Refresh /></el-icon>
              <span>{{ $t('刷新') }}</span>
            </el-button>
          </div>
        </div>

        <div v-loading="loading" class="account-list">
          <el-table :data="filteredAccounts" style="width: 100%" border>
            <el-table-column prop="account_name" ::label="$t('auto.FeishuAccount.2.5f0720')" min-width="160" />
            <el-table-column prop="app_id" label="App ID" min-width="200" />
            <el-table-column ::label="$t('auto.FeishuAccount.3.1b8578')" width="100" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'" effect="plain">
                  {{ scope.row.status === 1 ? '正常' : '停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="Webhook" width="110" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.webhook_enabled ? 'success' : 'info'" effect="plain">
                  {{ scope.row.webhook_enabled ? '已启用' : '未启用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.FeishuAccount.4.72fbb9')" width="110" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.ai_agent_enabled ? 'success' : 'info'" effect="plain">
                  {{ scope.row.ai_agent_enabled ? '已启用' : '未启用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="Token" width="110" align="center">
              <template #default="scope">
                <el-tooltip ::content="$t('auto.FeishuAccount.5.caf951')" placement="top">
                  <el-tag :type="scope.row.token_cached ? 'success' : 'warning'" effect="plain">
                    {{ scope.row.token_cached ? '已缓存' : '未缓存' }}
                  </el-tag>
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column prop="last_error_msg" :label="$t('auto.FeishuAccount.6.a2b83d')" min-width="180" show-overflow-tooltip />
            <el-table-column :label="$t('auto.FeishuAccount.7.2b6bc0')" width="420" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.FeishuAccount.28.95b351') }}</el-button>
                <el-button size="small" type="success" @click="handleTestSend(scope.row)">{{ $t('auto.FeishuAccount.29.e1c652') }}</el-button>
                <el-button size="small" type="warning" @click="handleRefreshToken(scope.row)">{{ $t('auto.FeishuAccount.30.742bef') }}</el-button>
                <el-button size="small" type="primary" @click="openBindingDialog(scope.row)">{{ $t('auto.FeishuAccount.31.c9e87d') }}</el-button>
                <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.FeishuAccount.32.2f4aad') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="!loading && filteredAccounts.length === 0" class="empty-data">
          <el-empty :description="$t('auto.FeishuAccount.8.277c85')" />
        </div>
      </div>
    </div>

    
    <el-dialog
      v-model="dialogVisible"
      ::title="$t('auto.FeishuAccount.9.ec5f1b')"
      width="640px"
      :close-on-click-modal="false"
    >
      <el-form ref="accountFormRef" :model="accountForm" :rules="rules" label-width="140px">
        <el-form-item :label="$t('auto.FeishuAccount.10.17327f')" prop="account_name">
          <el-input v-model="accountForm.account_name" :placeholder="$t('auto.FeishuAccount.11.28265a')" />
        </el-form-item>
        <el-form-item label="App ID" prop="app_id">
          <el-input v-model="accountForm.app_id" :placeholder="$t('auto.FeishuAccount.12.f44b31')" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input
            v-model="accountForm.app_secret"
            ::placeholder="$t('auto.FeishuAccount.13.6a205e')"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="Verification Token" prop="verification_token">
          <el-input
            v-model="accountForm.verification_token"
            :placeholder="$t('auto.FeishuAccount.14.4a9969')"
          />
        </el-form-item>
        <el-form-item label="Encrypt Key" prop="encrypt_key">
          <el-input
            v-model="accountForm.encrypt_key"
            :placeholder="$t('auto.FeishuAccount.15.16dd58')"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.16.685aab')" prop="webhook_enabled">
          <el-switch v-model="accountForm.webhook_enabled" />
          <span class="form-hint">{{ $t('auto.FeishuAccount.33.f46a63') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.17.706367')" prop="ai_agent_enabled">
          <el-switch v-model="accountForm.ai_agent_enabled" />
          <span class="form-hint">{{ $t('auto.FeishuAccount.34.78e4f3') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.18.3fea7c')" prop="status">
          <el-radio-group v-model="accountForm.status">
            <el-radio :label="1">{{ $t('auto.FeishuAccount.35.fd6e80') }}</el-radio>
            <el-radio :label="0">{{ $t('auto.FeishuAccount.36.5c56a8') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.FeishuAccount.37.625fb2') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{ $t('auto.FeishuAccount.38.38cf16') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="testSendVisible"
      :title="$t('auto.FeishuAccount.19.880075')"
      width="500px"
    >
      <el-form :model="testSendForm" label-width="100px">
        <el-form-item :label="$t('auto.FeishuAccount.20.88340e')">
          <el-input v-model="testSendForm.open_id" :placeholder="$t('auto.FeishuAccount.21.2ca9d3')" />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.22.b87b77')">
          <el-input v-model="testSendForm.content" type="textarea" :rows="4" :placeholder="$t('auto.FeishuAccount.23.dfbeaf')" />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.24.6d0071')">
          <el-select v-model="testSendForm.msg_type" :placeholder="$t('auto.FeishuAccount.25.b15a60')">
            <el-option label="text" value="text" />
            <el-option label="post" value="post" />
            <el-option label="interactive" value="interactive" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testSendVisible = false">{{ $t('auto.FeishuAccount.37.625fb2') }}</el-button>
        <el-button type="primary" @click="submitTestSend" :loading="testSending">{{ $t('auto.FeishuAccount.39.1535fc') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="bindingVisible"
      :title="$t('auto.FeishuAccount.26.253b0d')"
      width="500px"
    >
      <el-form :model="bindingForm" label-width="100px">
        <el-form-item :label="$t('auto.FeishuAccount.27.7035c6')">
          <el-input :value="bindingForm.account_name" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.4.72fbb9')">
          <el-switch v-model="bindingForm.ai_agent_enabled" />
        </el-form-item>
        <el-form-item label="Webhook">
          <el-switch v-model="bindingForm.webhook_enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindingVisible = false">{{ $t('auto.FeishuAccount.37.625fb2') }}</el-button>
        <el-button type="primary" @click="submitBinding" :loading="bindingSubmitting">{{ $t('auto.FeishuAccount.40.be5fbb') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { toList } from '@/utils/list'
import {
  listAccounts,
  createAccount,
  updateAccount,
  deleteAccount,
  testSend,
  refreshToken
} from '@/api/feishu'

const { t } = useI18n()

const loading = ref(false);
const submitting = ref(false)
const testSending = ref(false)
const bindingSubmitting = ref(false)
const accounts = ref([])
const searchKeyword = ref('')

const dialogVisible = ref(false);
const dialogType = ref('add');
const accountFormRef = ref(null)
const accountForm = ref({
  account_name: '',
  app_id: '',
  app_secret: '',
  verification_token: '',
  encrypt_key: '',
  webhook_enabled: false,
  ai_agent_enabled: false,
  status: 1
})

const testSendVisible = ref(false);
const testSendForm = ref({
  account_id: null,
  open_id: '',
  content: '',
  msg_type: 'text'
})

const bindingVisible = ref(false);
const bindingForm = ref({
  account_id: null,
  account_name: '',
  ai_agent_enabled: false,
  webhook_enabled: false
})

const rules = {
  account_name: [{ required: true, message: i18n.global.t(t('auto.FeishuAccount.41.1f3797')), trigger: 'blur' }],
  app_id: [{ required: true, message: i18n.global.t(t('auto.FeishuAccount.42.4288f9')), trigger: 'blur' }],
  app_secret: [{ required: true, message: i18n.global.t(t('auto.FeishuAccount.43.5eace4')), trigger: 'blur' }]
};

const filteredAccounts = computed(() => {
  if (!searchKeyword.value) return accounts.value
  return accounts.value.filter(a =>
    a.account_name?.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
});

const fetchAccounts = async () => {
  loading.value = true
  try {
    const res = await listAccounts()
    accounts.value = toList(res)
  } catch (e) {
    ElMessage.error(t('auto.FeishuAccount.44.d36a57') + e.message)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {}

const resetForm = () => {
  accountForm.value = {
    account_name: '',
    app_id: '',
    app_secret: '',
    verification_token: '',
    encrypt_key: '',
    webhook_enabled: false,
    ai_agent_enabled: false,
    status: 1
  }
}

const handleAdd = () => {
  dialogType.value = 'add'
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogType.value = 'edit'
  accountForm.value = {
    id: row.id,
    account_name: row.account_name,
    app_id: row.app_id,
    app_secret: '',
    verification_token: row.verification_token || '',
    encrypt_key: '',
    webhook_enabled: row.webhook_enabled || false,
    ai_agent_enabled: row.ai_agent_enabled || false,
    status: row.status
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!accountFormRef.value) return
  await accountFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const isEdit = dialogType.value === 'edit'
      await (isEdit
        ? updateAccount(accountForm.value.id, accountForm.value)
        : createAccount(accountForm.value))
      ElMessage.success(isEdit ? t('auto.FeishuAccount.45.55aa63') : t('auto.FeishuAccount.46.04a691'))
      dialogVisible.value = false
      await fetchAccounts()
    } catch (e) {
      ElMessage.error(t('auto.FeishuAccount.44.d36a57') + e.message)
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除账号 "${row.account_name}" 吗？`,
      t('auto.FeishuAccount.47.631cd2'),
      { type: 'warning' }
    )
    await deleteAccount(row.id)
    ElMessage.success(i18n.global.t(t('auto.FeishuAccount.48.0007d1')))
    await fetchAccounts()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(t('auto.FeishuAccount.49.79585f') + e.message)
    }
  }
}

const handleTestSend = (row) => {
  testSendForm.value = {
    account_id: row.id,
    open_id: '',
    content: i18n.global.t(t('auto.FeishuAccount.50.8a145f')) + row.account_name,
    msg_type: 'text'
  }
  testSendVisible.value = true
}

const submitTestSend = async () => {
  if (!testSendForm.value.open_id) {
    ElMessage.warning(i18n.global.t(t('auto.FeishuAccount.51.3d80fa')))
    return
  }
  if (!testSendForm.value.content) {
    ElMessage.warning(i18n.global.t(t('auto.FeishuAccount.52.fa967b')))
    return
  }
  testSending.value = true
  try {
    await testSend(testSendForm.value.account_id, {
      open_id: testSendForm.value.open_id,
      content: testSendForm.value.content,
      msg_type: testSendForm.value.msg_type
    })
    ElMessage.success(i18n.global.t(t('auto.FeishuAccount.53.8c6e34')))
    testSendVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.FeishuAccount.44.d36a57') + e.message)
  } finally {
    testSending.value = false
  }
}

const handleRefreshToken = async (row) => {
  try {
    await refreshToken(row.id)
    ElMessage.success(i18n.global.t('Token 刷新成功'))
    await fetchAccounts()
  } catch (e) {
    ElMessage.error(t('auto.FeishuAccount.44.d36a57') + e.message)
  }
}

const openBindingDialog = (row) => {
  bindingForm.value = {
    account_id: row.id,
    account_name: row.account_name,
    ai_agent_enabled: row.ai_agent_enabled || false,
    webhook_enabled: row.webhook_enabled || false
  }
  bindingVisible.value = true
}

const submitBinding = async () => {
  bindingSubmitting.value = true
  try {
    await updateAccount(bindingForm.value.account_id, {
      ai_agent_enabled: bindingForm.value.ai_agent_enabled,
      webhook_enabled: bindingForm.value.webhook_enabled
    })
    ElMessage.success(i18n.global.t(t('auto.FeishuAccount.54.3b1083')))
    bindingVisible.value = false
    await fetchAccounts()
  } catch (e) {
    ElMessage.error(t('auto.FeishuAccount.44.d36a57') + e.message)
  } finally {
    bindingSubmitting.value = false
  }
}

onMounted(() => {
  fetchAccounts()
});
</script>

<style scoped>
.feishu-account-management {
  padding: 16px;
}
.account-search {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
}
.action-buttons {
  display: flex;
  gap: 8px;
}
.empty-data {
  padding: 40px 0;
  text-align: center;
}
.form-hint {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
</style>
