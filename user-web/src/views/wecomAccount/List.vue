<template>
  <div class="wecom-account-list">
    
    <el-row :gutter="16" class="summary-row" v-loading="summaryLoading">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-card">
            <div class="summary-label">{{ $t('账号总数') }}</div>
            <div class="summary-value">{{ summary?.total_accounts ?? 0 }}</div>
            <div class="summary-sub">在线 {{ summary?.online_count ?? 0 }} / 离线 {{ summary?.offline_count ?? 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-card">
            <div class="summary-label">{{ $t('平均健康分') }}</div>
            <div class="summary-value" :class="healthScoreClass">{{ summary?.avg_score?.toFixed(1) ?? '-' }}</div>
            <div class="summary-sub">{{ $t('auto.List.f32efbf5') }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-card">
            <div class="summary-label">{{ $t('配额使用') }}</div>
            <div class="summary-value">{{ summary?.total_used ?? 0 }} / {{ summary?.total_quota ?? 0 }}</div>
            <div class="summary-sub">{{ $t('日配额') }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="summary-card">
            <div class="summary-label">{{ $t('风险账号') }}</div>
            <div class="summary-value" :class="{ 'text-danger': (summary?.risk_accounts?.length ?? 0) > 0 }">
              {{ summary?.risk_accounts?.length ?? 0 }}
            </div>
            <div class="summary-sub">{{ $t('需关注') }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('企微账号管理') }}</span>
          <div>
            <el-button type="primary" :icon="Plus" @click="openCreateDialog">{{ $t('新增账号') }}</el-button>
            <el-button type="primary" :icon="Refresh" @click="loadData" :loading="listLoading">{{ $t('刷新') }}</el-button>
          </div>
        </div>
      </template>

      <el-table :data="accountList" v-loading="listLoading" border style="width: 100%">
        <el-table-column prop="account.corp_id" :label="$t('auto.List.5a83cb3a')" min-width="160" />
        <el-table-column prop="account.agent_id" :label="$t('auto.List.3751f65f')" width="100" />
        <el-table-column :label="$t('auto.List.959aca6b')" width="110">
          <template #default="{ row }">
            <el-tag :type="loginStateTag(row.account.login_state)" size="small">
              {{ loginStateText(row.account.login_state) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.25eb6ba3')" width="110">
          <template #default="{ row }">
            <el-tag :type="riskLevelTag(row.account.risk_level)" size="small" effect="dark">
              {{ riskLevelText(row.account.risk_level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.1661a791')" width="100">
          <template #default="{ row }">
            <span :class="healthScoreTextClass(row.health?.health_score)">
              {{ row.health?.health_score?.toFixed(0) ?? '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.2ce83433')" width="120">
          <template #default="{ row }">
            {{ row.account.friend_count ?? 0 }} / {{ row.account.group_count ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.f74c1da4')" width="130">
          <template #default="{ row }">
            {{ row.account.daily_msg_used ?? 0 }} / {{ row.account.daily_msg_quota ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.7f2846ef')" min-width="160">
          <template #default="{ row }">
            {{ row.account.last_active_at ? formatTime(row.account.last_active_at) : '从未' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
            <el-button size="small" type="primary" @click="openEditDialog(row.account)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.account)">{{ $t('auto.List.2f4aaddd') }}</el-button>
            <el-dropdown trigger="click" :disabled="!!actionLoading" @command="(cmd) => handleMoreCommand(cmd, row)">
              <el-button size="small" type="warning">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="data" :icon="DataLine">{{ $t('auto.List.b12af094') }}</el-dropdown-item>
                  <el-dropdown-item command="bind" :icon="Link">{{ $t('auto.List.c9e87d56') }}</el-dropdown-item>
                  <el-dropdown-item command="syncCustomers" :icon="Refresh">{{ $t('auto.List.76303737') }}</el-dropdown-item>
                  <el-dropdown-item command="syncGroups" :icon="Refresh">{{ $t('auto.List.0e55482c') }}</el-dropdown-item>
                  <el-dropdown-item command="syncTags" :icon="Refresh">{{ $t('auto.List.17caa1a4') }}</el-dropdown-item>
                  <el-dropdown-item command="send" :icon="Promotion">{{ $t('auto.List.e1c652bc') }}</el-dropdown-item>
                  <el-dropdown-item command="refresh" :icon="RefreshRight" divided>{{ $t('auto.List.e0846f76') }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <AgentBindingDialog
      v-model="bindingDialogVisible"
      channel-type="wecom"
      :account-id="bindingAccountId"
      :account-label="bindingAccountLabel"
      :account-enabled="bindingAccountEnabled"
    />

    
    <WeComSendDialog v-model:visible="sendDialogVisible" :account-id="sendAccountId" />

    
    <el-dialog v-model="createDialogVisible" :title="dialogTitle" width="520px" @closed="resetCreateForm">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="110px">
        <el-form-item label="CorpID" prop="corp_id">
          <el-input v-model="createForm.corp_id" :placeholder="$t('auto.List.30af5c40')" />
        </el-form-item>
        <el-form-item label="CorpSecret" prop="corp_secret">
          <el-input v-model="createForm.corp_secret" type="password" show-password :placeholder="$t('auto.List.ff2ae665')" />
        </el-form-item>
        <el-form-item label="AgentID" prop="agent_id">
          <el-input v-model="createForm.agent_id" :placeholder="$t('auto.List.e2efbaa2')" />
        </el-form-item>
        <el-form-item label="AgentSecret" prop="agent_secret">
          <el-input v-model="createForm.agent_secret" type="password" show-password :placeholder="$t('auto.List.00718aa0')" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('auto.List.045c3540') }}</el-divider>
        <el-form-item :label="$t('auto.List.5c7f1a99')" v-if="dialogMode === 'edit'">
          <el-input :model-value="wecomCallbackUrl" readonly>
            <template #append>
              <el-button @click="copyCallbackUrl">{{ $t('auto.List.79d3abe9') }}</el-button>
            </template>
          </el-input>
          <div class="form-hint">{{ $t('auto.List.3730189a') }}</div>
        </el-form-item>
        <el-form-item label="Token">
          <el-input v-model="createForm.callback_token" :placeholder="$t('auto.List.22d91bb6')" />
        </el-form-item>
        <el-form-item label="EncodingAESKey">
          <el-input v-model="createForm.encoding_aes_key" :placeholder="$t('auto.List.9678f744')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.5f1e0e72')">
          <el-switch v-model="createForm.webhook_enabled" />
          <span class="form-hint-inline">{{ $t('auto.List.1aad39ea') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.List.bfdd28c7')">
          <el-switch v-model="createForm.ai_agent_enabled" />
          <span class="form-hint-inline">{{ $t('auto.List.809880b9') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="createLoading" @click="submitForm">{{ $t('auto.List.38cf16f2') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailDialogVisible" :title="$t('auto.List.9a1d8118')" width="640px">
      <el-descriptions v-if="detailAccount" :column="2" border size="small">
        <el-descriptions-item :label="$t('auto.List.14d1718c')">{{ detailAccount.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.5a83cb3a')">{{ detailAccount.corp_id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3751f65f')">{{ detailAccount.agent_id || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.959aca6b')">
          <el-tag :type="loginStateTag(detailAccount.login_state)" size="small">
            {{ loginStateText(detailAccount.login_state) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.25eb6ba3')">
          <el-tag :type="riskLevelTag(detailAccount.risk_level)" size="small" effect="dark">
            {{ riskLevelText(detailAccount.risk_level) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.1661a791')">
          {{ detailHealth?.health_score?.toFixed(0) ?? '-' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.fe59f4e5')">{{ detailAccount.friend_count ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3d7d8031')">{{ detailAccount.group_count ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.f74c1da4')">
          {{ detailAccount.daily_msg_used ?? 0 }} / {{ detailAccount.daily_msg_quota ?? 0 }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.eef65de8')">{{ detailAccount.total_sent ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.b6b5a566')">{{ detailAccount.total_received ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.7f2846ef')">
          {{ detailAccount.last_active_at ? formatTime(detailAccount.last_active_at) : '从未' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.eca37cb0')">
          {{ detailAccount.created_at ? formatTime(detailAccount.created_at) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailHealth" :label="$t('auto.List.4604d502')" :span="2">
          {{ detailHealth.error_message || '无' }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailHealth" :label="$t('auto.List.735be255')" :span="2">
          {{ detailHealth.reported_at ? formatTime(detailHealth.reported_at) : '暂无' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus, ArrowDown, Link, Promotion, RefreshRight, DataLine } from '@element-plus/icons-vue'
import { wecomAccountApi } from '@/api/wecomAccount.js'
import AgentBindingDialog from '@/components/AgentBindingDialog.vue'
import WeComSendDialog from '@/components/WeComSendDialog.vue'
import { getAccountStatusLabel, getAccountStatusTagType } from '@/constants/accountType';

const router = useRouter()

const { t } = useI18n()

const listLoading = ref(false)
const summaryLoading = ref(false)
const accountList = ref([])
const summary = ref(null)

const createDialogVisible = ref(false);
const createLoading = ref(false)
const createFormRef = ref(null)
const dialogMode = ref('create');
const editingId = ref(null)
const createForm = reactive({
  corp_id: '',
  corp_secret: '',
  agent_id: '',
  agent_secret: '',
  callback_token: '',
  encoding_aes_key: '',
  webhook_enabled: false,
  ai_agent_enabled: false
})
const createRules = {
  corp_id: [{ required: true, message: t('auto.List.46e5f855'), trigger: 'blur' }],
  corp_secret: [{ required: true, message: t('auto.List.b34a268e'), trigger: 'blur' }]
}
const dialogTitle = computed(() => (dialogMode.value === 'edit' ? t('auto.List.e03831ae') : t('auto.List.ec0a9389')))
const wecomCallbackUrl = computed(() => {
  if (dialogMode.value !== 'edit' || editingId.value == null) return ''
  return `${window.location.origin}/api/webhook/wecom/${editingId.value}`
})
const copyCallbackUrl = async () => {
  try {
    await navigator.clipboard.writeText(wecomCallbackUrl.value)
    ElMessage.success(t('auto.List.94f2c63c'))
  } catch {
    ElMessage.error(t('auto.List.41da1ca4'))
  }
}
const openCreateDialog = () => {
  dialogMode.value = 'create'
  editingId.value = null
  resetCreateForm()
  createDialogVisible.value = true
}
const openEditDialog = (account) => {
  dialogMode.value = 'edit'
  editingId.value = account.id
  createForm.corp_id = account.corp_id || ''
  createForm.corp_secret = account.corp_secret || ''
  createForm.agent_id = account.agent_id ? String(account.agent_id) : ''
  createForm.agent_secret = account.agent_secret || ''
  createForm.callback_token = account.callback_token || ''
  createForm.encoding_aes_key = account.encoding_aes_key || ''
  createForm.webhook_enabled = !!account.webhook_enabled
  createForm.ai_agent_enabled = !!account.ai_agent_enabled
  createDialogVisible.value = true
}
const resetCreateForm = () => {
  createForm.corp_id = ''
  createForm.corp_secret = ''
  createForm.agent_id = ''
  createForm.agent_secret = ''
  createForm.callback_token = ''
  createForm.encoding_aes_key = ''
  createForm.webhook_enabled = false
  createForm.ai_agent_enabled = false
  createFormRef.value?.clearValidate?.()
}
const submitForm = async () => {
  if (!createFormRef.value) return
  try {
    await createFormRef.value.validate()
  } catch {
    return
  }
  createLoading.value = true
  try {
    const payload = {
      corp_id: createForm.corp_id,
      corp_secret: createForm.corp_secret,
      agent_id: createForm.agent_id ? Number(createForm.agent_id) : undefined,
      agent_secret: createForm.agent_secret || undefined,
      callback_token: createForm.callback_token || undefined,
      encoding_aes_key: createForm.encoding_aes_key || undefined,
      webhook_enabled: createForm.webhook_enabled,
      ai_agent_enabled: createForm.ai_agent_enabled
    }
    if (dialogMode.value === 'edit' && editingId.value != null) {
      await wecomAccountApi.updateAccount(editingId.value, payload)
      ElMessage.success(t('auto.List.8a225eba'))
    } else {
      await wecomAccountApi.createAccount(payload)
      ElMessage.success(t('auto.List.482ecffa'))
    }
    createDialogVisible.value = false
    loadData()
  } catch (err) {
    ElMessage.error(t('auto.List.40f90217') + (err.message || t('auto.List.974e7484')))
  } finally {
    createLoading.value = false
  }
}

const detailDialogVisible = ref(false);
const detailAccount = ref(null)
const detailHealth = ref(null)
const viewDetail = (row) => {
  detailAccount.value = row.account
  detailHealth.value = row.health
  detailDialogVisible.value = true
}

const handleDelete = async (account) => {
  try {
    await ElMessageBox.confirm(`确认删除企业微信账号「${account.corp_id}」？该操作不可恢复。`, t('auto.List.50eaf94d'), {
      type: 'warning',
      confirmButtonText: t('auto.List.2f4aaddd'),
      cancelButtonText: t('auto.List.625fb26b')
    })
  } catch {
    return
  }
  try {
    await wecomAccountApi.deleteAccount(account.id)
    ElMessage.success(t('auto.List.0007d170'))
    loadData()
  } catch (err) {
    ElMessage.error(t('auto.List.ad23f072') + (err.message || t('auto.List.974e7484')))
  }
};

const healthScoreClass = computed(() => {
  const score = summary.value?.avg_score ?? 0
  if (score >= 80) return 'text-success'
  if (score >= 60) return 'text-warning'
  return 'text-danger'
})

const loginStateText = (state) => {
  const map = { online: t('auto.List.68905cf3'), offline: t('auto.List.50d4a850'), banned: t('auto.List.b08d3870') }
  return map[state] || state
}

const loginStateTag = (state) => {
  const map = { online: 'success', offline: 'info', banned: 'danger' }
  return map[state] || 'info'
}

const riskLevelText = (level) => {
  const map = { normal: t('auto.List.fd6e80f1'), warning: t('auto.List.900c70fa'), critical: t('auto.List.540f84dd'), banned: t('auto.List.b08d3870') }
  return map[level] || level
}

const riskLevelTag = (level) => {
  const map = { normal: 'success', warning: 'warning', critical: 'danger', banned: 'danger' }
  return map[level] || 'info'
}

const healthScoreTextClass = (score) => {
  if (score == null) return ''
  if (score >= 80) return 'text-success'
  if (score >= 60) return 'text-warning'
  return 'text-danger'
}

const formatTime = (iso) => {
  if (!iso) return '-'
  return new Date(iso).toLocaleString('zh-CN')
}

const loadAccountList = async () => {
  listLoading.value = true
  try {
    const data = await wecomAccountApi.listAccounts()
    accountList.value = Array.isArray(data) ? data : []
  } catch (err) {
    ElMessage.error(t('auto.List.77937932') + (err.message || t('auto.List.974e7484')))
    accountList.value = []
  } finally {
    listLoading.value = false
  }
}

const loadSummary = async () => {
  summaryLoading.value = true
  try {
    const data = await wecomAccountApi.getHealthSummary()
    summary.value = data
  } catch (err) {
    ElMessage.error(t('auto.List.4aad39ca') + (err.message || t('auto.List.974e7484')))
    summary.value = null
  } finally {
    summaryLoading.value = false
  }
}

const loadData = () => {
  loadAccountList()
  loadSummary()
}

onMounted(() => {
  loadData()
})

const bindingDialogVisible = ref(false);
const bindingAccountId = ref('')
const bindingAccountLabel = ref('')
const bindingAccountEnabled = ref(true)

const openBindingDialog = (account) => {
  bindingAccountId.value = String(account.id)
  bindingAccountLabel.value = account.corp_id || ''
  bindingAccountEnabled.value = account.status === 1
  bindingDialogVisible.value = true
}

const actionLoading = ref('');

const sendDialogVisible = ref(false);
const sendAccountId = ref(null)

const openSendDialog = (account) => {
  sendAccountId.value = account.id
  sendDialogVisible.value = true
}

const syncAction = async (account, type) => {
  const apiMap = {
    syncCustomers: wecomAccountApi.syncCustomers,
    syncGroups: wecomAccountApi.syncGroups,
    syncTags: wecomAccountApi.syncTags
  }
  const api = apiMap[type]
  if (!api) return
  actionLoading.value = type
  try {
    const res = await api(account.id)
    const count = res && res.count != null ? res.count : 0
    ElMessage.success(t('auto.List.02f667de') + (count ? `，共 ${count} 条` : ''))
  } catch (e) {
    ElMessage.error(t('auto.List.a0f67c55') + (e.message || e))
  } finally {
    actionLoading.value = ''
  }
};

const refreshAccount = async (account) => {
  actionLoading.value = 'refresh'
  try {
    await wecomAccountApi.refreshAccount(account.id)
    ElMessage.success(t('auto.List.9b353cfb'))
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.List.df570e03') + (e.message || e))
  } finally {
    actionLoading.value = ''
  }
};

const handleMoreCommand = (cmd, row) => {
  const account = row.account
  switch (cmd) {
    case 'bind':
      openBindingDialog(account)
      break
    case 'syncCustomers':
      syncAction(account, 'syncCustomers')
      break
    case 'syncGroups':
      syncAction(account, 'syncGroups')
      break
    case 'syncTags':
      syncAction(account, 'syncTags')
      break
    case 'send':
      openSendDialog(account)
      break
    case 'data':
      router.push({ name: 'WecomAccountData' })
      break
    case 'refresh':
      refreshAccount(account)
      break
  }
};
</script>

<style scoped>
.wecom-account-list {
  padding: 4px;
}

.summary-row {
  margin-bottom: 16px;
}

.summary-card {
  text-align: center;
  padding: 8px 0;
}

.summary-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}

.summary-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  line-height: 1.2;
}

.summary-sub {
  font-size: 12px;
  color: #c0c4cc;
  margin-top: 6px;
}

.table-card {
  margin-top: 4px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-hint {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}

.form-hint-inline {
  font-size: 12px;
  color: #909399;
  margin-left: 10px;
}

.text-success {
  color: #10B981;
}

.text-warning {
  color: #F59E0B;
}

.text-danger {
  color: #EF4444;
}
</style>
