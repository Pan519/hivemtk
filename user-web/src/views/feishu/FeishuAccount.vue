<template>
  <div class="feishu-account-management">
    <div class="account-tabs">
      <div class="account-list-container">
        <div class="account-search">
          <el-input
            v-model="searchKeyword"
            ::placeholder="$t('auto.FeishuAccount.9d9a5fe2')"
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
            <el-table-column prop="account_name" ::label="$t('auto.FeishuAccount.5f0720d4')" min-width="160" />
            <el-table-column prop="app_id" label="App ID" min-width="200" />
            <el-table-column ::label="$t('auto.FeishuAccount.1b857887')" width="100" align="center">
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
            <el-table-column :label="$t('auto.FeishuAccount.72fbb9b5')" width="110" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.ai_agent_enabled ? 'success' : 'info'" effect="plain">
                  {{ scope.row.ai_agent_enabled ? '已启用' : '未启用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="Token" width="110" align="center">
              <template #default="scope">
                <el-tooltip ::content="$t('auto.FeishuAccount.caf951e0')" placement="top">
                  <el-tag :type="scope.row.token_cached ? 'success' : 'warning'" effect="plain">
                    {{ scope.row.token_cached ? '已缓存' : '未缓存' }}
                  </el-tag>
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column prop="last_error_msg" :label="$t('auto.FeishuAccount.a2b83df0')" min-width="180" show-overflow-tooltip />
            <el-table-column :label="$t('auto.FeishuAccount.2b6bc0f2')" width="420" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.FeishuAccount.95b351c8') }}</el-button>
                <el-button size="small" type="success" @click="handleTestSend(scope.row)">{{ $t('auto.FeishuAccount.e1c652bc') }}</el-button>
                <el-button size="small" type="warning" @click="handleRefreshToken(scope.row)">{{ $t('auto.FeishuAccount.742befe5') }}</el-button>
                <el-button size="small" type="primary" @click="openBindingDialog(scope.row)">{{ $t('auto.FeishuAccount.c9e87d56') }}</el-button>
                <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.FeishuAccount.2f4aaddd') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="!loading && filteredAccounts.length === 0" class="empty-data">
          <el-empty :description="$t('auto.FeishuAccount.277c85b3')" />
        </div>
      </div>
    </div>

    
    <el-dialog
      v-model="dialogVisible"
      ::title="$t('auto.FeishuAccount.ec5f1b30')"
      width="640px"
      :close-on-click-modal="false"
    >
      <el-form ref="accountFormRef" :model="accountForm" :rules="rules" label-width="140px">
        <el-form-item :label="$t('auto.FeishuAccount.17327f9c')" prop="account_name">
          <el-input v-model="accountForm.account_name" :placeholder="$t('auto.FeishuAccount.28265a1e')" />
        </el-form-item>
        <el-form-item label="App ID" prop="app_id">
          <el-input v-model="accountForm.app_id" :placeholder="$t('auto.FeishuAccount.f44b31a1')" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input
            v-model="accountForm.app_secret"
            ::placeholder="$t('auto.FeishuAccount.6a205ea0')"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="Verification Token" prop="verification_token">
          <el-input
            v-model="accountForm.verification_token"
            :placeholder="$t('auto.FeishuAccount.4a9969bd')"
          />
        </el-form-item>
        <el-form-item label="Encrypt Key" prop="encrypt_key">
          <el-input
            v-model="accountForm.encrypt_key"
            :placeholder="$t('auto.FeishuAccount.16dd58c9')"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.685aab5e')" prop="webhook_enabled">
          <el-switch v-model="accountForm.webhook_enabled" />
          <span class="form-hint">{{ $t('auto.FeishuAccount.f46a63ed') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.7063676a')" prop="ai_agent_enabled">
          <el-switch v-model="accountForm.ai_agent_enabled" />
          <span class="form-hint">{{ $t('auto.FeishuAccount.78e4f38c') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.3fea7ca7')" prop="status">
          <el-radio-group v-model="accountForm.status">
            <el-radio :label="1">{{ $t('auto.FeishuAccount.fd6e80f1') }}</el-radio>
            <el-radio :label="0">{{ $t('auto.FeishuAccount.5c56a889') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.FeishuAccount.625fb26b') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{ $t('auto.FeishuAccount.38cf16f2') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="testSendVisible"
      :title="$t('auto.FeishuAccount.88007514')"
      width="500px"
    >
      <el-form :model="testSendForm" label-width="100px">
        <el-form-item :label="$t('auto.FeishuAccount.88340ead')">
          <el-input v-model="testSendForm.open_id" :placeholder="$t('auto.FeishuAccount.2ca9d3ed')" />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.b87b7756')">
          <el-input v-model="testSendForm.content" type="textarea" :rows="4" :placeholder="$t('auto.FeishuAccount.dfbeaf44')" />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.6d00710a')">
          <el-select v-model="testSendForm.msg_type" :placeholder="$t('auto.FeishuAccount.b15a6054')">
            <el-option label="text" value="text" />
            <el-option label="post" value="post" />
            <el-option label="interactive" value="interactive" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testSendVisible = false">{{ $t('auto.FeishuAccount.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitTestSend" :loading="testSending">{{ $t('auto.FeishuAccount.1535fcfa') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="bindingVisible"
      :title="$t('auto.FeishuAccount.253b0d24')"
      width="500px"
    >
      <el-form :model="bindingForm" label-width="100px">
        <el-form-item :label="$t('auto.FeishuAccount.7035c62f')">
          <el-input :value="bindingForm.account_name" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.FeishuAccount.72fbb9b5')">
          <el-switch v-model="bindingForm.ai_agent_enabled" />
        </el-form-item>
        <el-form-item label="Webhook">
          <el-switch v-model="bindingForm.webhook_enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bindingVisible = false">{{ $t('auto.FeishuAccount.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitBinding" :loading="bindingSubmitting">{{ $t('auto.FeishuAccount.be5fbbe3') }}</el-button>
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
  account_name: [{ required: true, message: i18n.global.t(t('auto.FeishuAccount.1f3797f1')), trigger: 'blur' }],
  app_id: [{ required: true, message: i18n.global.t(t('auto.FeishuAccount.4288f98f')), trigger: 'blur' }],
  app_secret: [{ required: true, message: i18n.global.t(t('auto.FeishuAccount.5eace474')), trigger: 'blur' }]
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
    ElMessage.error(t('auto.FeishuAccount.d36a5748') + e.message)
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
      ElMessage.success(isEdit ? t('auto.FeishuAccount.55aa6366') : t('auto.FeishuAccount.04a691b3'))
      dialogVisible.value = false
      await fetchAccounts()
    } catch (e) {
      ElMessage.error(t('auto.FeishuAccount.d36a5748') + e.message)
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除账号 "${row.account_name}" 吗？`,
      t('auto.FeishuAccount.631cd220'),
      { type: 'warning' }
    )
    await deleteAccount(row.id)
    ElMessage.success(i18n.global.t(t('auto.FeishuAccount.0007d170')))
    await fetchAccounts()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(t('auto.FeishuAccount.79585f92') + e.message)
    }
  }
}

const handleTestSend = (row) => {
  testSendForm.value = {
    account_id: row.id,
    open_id: '',
    content: i18n.global.t(t('auto.FeishuAccount.8a145fb5')) + row.account_name,
    msg_type: 'text'
  }
  testSendVisible.value = true
}

const submitTestSend = async () => {
  if (!testSendForm.value.open_id) {
    ElMessage.warning(i18n.global.t(t('auto.FeishuAccount.3d80fae7')))
    return
  }
  if (!testSendForm.value.content) {
    ElMessage.warning(i18n.global.t(t('auto.FeishuAccount.fa967bd0')))
    return
  }
  testSending.value = true
  try {
    await testSend(testSendForm.value.account_id, {
      open_id: testSendForm.value.open_id,
      content: testSendForm.value.content,
      msg_type: testSendForm.value.msg_type
    })
    ElMessage.success(i18n.global.t(t('auto.FeishuAccount.8c6e34ef')))
    testSendVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.FeishuAccount.d36a5748') + e.message)
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
    ElMessage.error(t('auto.FeishuAccount.d36a5748') + e.message)
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
    ElMessage.success(i18n.global.t(t('auto.FeishuAccount.3b108349')))
    bindingVisible.value = false
    await fetchAccounts()
  } catch (e) {
    ElMessage.error(t('auto.FeishuAccount.d36a5748') + e.message)
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
