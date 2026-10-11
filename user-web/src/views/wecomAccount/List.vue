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
            <div class="summary-sub">{{ $t('auto.List.28.f32efb') }}</div>
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
        <el-table-column prop="account.corp_id" :label="$t('auto.List.1.5a83cb')" min-width="160" />
        <el-table-column prop="account.agent_id" :label="$t('auto.List.2.3751f6')" width="100" />
        <el-table-column :label="$t('auto.List.3.959aca')" width="110">
          <template #default="{ row }">
            <el-tag :type="loginStateTag(row.account.login_state)" size="small">
              {{ loginStateText(row.account.login_state) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.4.25eb6b')" width="110">
          <template #default="{ row }">
            <el-tag :type="riskLevelTag(row.account.risk_level)" size="small" effect="dark">
              {{ riskLevelText(row.account.risk_level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.5.1661a7')" width="100">
          <template #default="{ row }">
            <span :class="healthScoreTextClass(row.health?.health_score)">
              {{ row.health?.health_score?.toFixed(0) ?? '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.6.2ce834')" width="120">
          <template #default="{ row }">
            {{ row.account.friend_count ?? 0 }} / {{ row.account.group_count ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.7.f74c1d')" width="130">
          <template #default="{ row }">
            {{ row.account.daily_msg_used ?? 0 }} / {{ row.account.daily_msg_quota ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.8.7f2846')" min-width="160">
          <template #default="{ row }">
            {{ row.account.last_active_at ? formatTime(row.account.last_active_at) : '从未' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.9.2b6bc0')" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">{{ $t('auto.List.29.f26225') }}</el-button>
            <el-button size="small" type="primary" @click="openEditDialog(row.account)">{{ $t('auto.List.30.95b351') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.account)">{{ $t('auto.List.31.2f4aad') }}</el-button>
            <el-dropdown trigger="click" :disabled="!!actionLoading" @command="(cmd) => handleMoreCommand(cmd, row)">
              <el-button size="small" type="warning">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="data" :icon="DataLine">{{ $t('auto.List.32.b12af0') }}</el-dropdown-item>
                  <el-dropdown-item command="bind" :icon="Link">{{ $t('auto.List.33.c9e87d') }}</el-dropdown-item>
                  <el-dropdown-item command="syncCustomers" :icon="Refresh">{{ $t('auto.List.34.763037') }}</el-dropdown-item>
                  <el-dropdown-item command="syncGroups" :icon="Refresh">{{ $t('auto.List.35.0e5548') }}</el-dropdown-item>
                  <el-dropdown-item command="syncTags" :icon="Refresh">{{ $t('auto.List.36.17caa1') }}</el-dropdown-item>
                  <el-dropdown-item command="send" :icon="Promotion">{{ $t('auto.List.37.e1c652') }}</el-dropdown-item>
                  <el-dropdown-item command="refresh" :icon="RefreshRight" divided>{{ $t('auto.List.38.e0846f') }}</el-dropdown-item>
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
          <el-input v-model="createForm.corp_id" :placeholder="$t('auto.List.10.30af5c')" />
        </el-form-item>
        <el-form-item label="CorpSecret" prop="corp_secret">
          <el-input v-model="createForm.corp_secret" type="password" show-password :placeholder="$t('auto.List.11.ff2ae6')" />
        </el-form-item>
        <el-form-item label="AgentID" prop="agent_id">
          <el-input v-model="createForm.agent_id" :placeholder="$t('auto.List.12.e2efba')" />
        </el-form-item>
        <el-form-item label="AgentSecret" prop="agent_secret">
          <el-input v-model="createForm.agent_secret" type="password" show-password :placeholder="$t('auto.List.13.00718a')" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('auto.List.39.045c35') }}</el-divider>
        <el-form-item :label="$t('auto.List.14.5c7f1a')" v-if="dialogMode === 'edit'">
          <el-input :model-value="wecomCallbackUrl" readonly>
            <template #append>
              <el-button @click="copyCallbackUrl">{{ $t('auto.List.40.79d3ab') }}</el-button>
            </template>
          </el-input>
          <div class="form-hint">{{ $t('auto.List.41.373018') }}</div>
        </el-form-item>
        <el-form-item label="Token">
          <el-input v-model="createForm.callback_token" :placeholder="$t('auto.List.15.22d91b')" />
        </el-form-item>
        <el-form-item label="EncodingAESKey">
          <el-input v-model="createForm.encoding_aes_key" :placeholder="$t('auto.List.16.9678f7')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.17.5f1e0e')">
          <el-switch v-model="createForm.webhook_enabled" />
          <span class="form-hint-inline">{{ $t('auto.List.42.1aad39') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.List.18.bfdd28')">
          <el-switch v-model="createForm.ai_agent_enabled" />
          <span class="form-hint-inline">{{ $t('auto.List.43.809880') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">{{ $t('auto.List.44.625fb2') }}</el-button>
        <el-button type="primary" :loading="createLoading" @click="submitForm">{{ $t('auto.List.45.38cf16') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailDialogVisible" :title="$t('auto.List.19.9a1d81')" width="640px">
      <el-descriptions v-if="detailAccount" :column="2" border size="small">
        <el-descriptions-item :label="$t('auto.List.20.14d171')">{{ detailAccount.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.1.5a83cb')">{{ detailAccount.corp_id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.2.3751f6')">{{ detailAccount.agent_id || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3.959aca')">
          <el-tag :type="loginStateTag(detailAccount.login_state)" size="small">
            {{ loginStateText(detailAccount.login_state) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.4.25eb6b')">
          <el-tag :type="riskLevelTag(detailAccount.risk_level)" size="small" effect="dark">
            {{ riskLevelText(detailAccount.risk_level) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.5.1661a7')">
          {{ detailHealth?.health_score?.toFixed(0) ?? '-' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.21.fe59f4')">{{ detailAccount.friend_count ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.22.3d7d80')">{{ detailAccount.group_count ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.7.f74c1d')">
          {{ detailAccount.daily_msg_used ?? 0 }} / {{ detailAccount.daily_msg_quota ?? 0 }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.23.eef65d')">{{ detailAccount.total_sent ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.24.b6b5a5')">{{ detailAccount.total_received ?? 0 }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.8.7f2846')">
          {{ detailAccount.last_active_at ? formatTime(detailAccount.last_active_at) : '从未' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.25.eca37c')">
          {{ detailAccount.created_at ? formatTime(detailAccount.created_at) : '-' }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailHealth" :label="$t('auto.List.26.4604d5')" :span="2">
          {{ detailHealth.error_message || '无' }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailHealth" :label="$t('auto.List.27.735be2')" :span="2">
          {{ detailHealth.reported_at ? formatTime(detailHealth.reported_at) : '暂无' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">{{ $t('auto.List.46.b15d91') }}</el-button>
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
  corp_id: [{ required: true, message: t('auto.List.47.46e5f8'), trigger: 'blur' }],
  corp_secret: [{ required: true, message: t('auto.List.48.b34a26'), trigger: 'blur' }]
}
const dialogTitle = computed(() => (dialogMode.value === 'edit' ? t('auto.List.49.e03831') : t('auto.List.50.ec0a93')))
const wecomCallbackUrl = computed(() => {
  if (dialogMode.value !== 'edit' || editingId.value == null) return ''
  return `${window.location.origin}/api/webhook/wecom/${editingId.value}`
})
const copyCallbackUrl = async () => {
  try {
    await navigator.clipboard.writeText(wecomCallbackUrl.value)
    ElMessage.success(t('auto.List.51.94f2c6'))
  } catch {
    ElMessage.error(t('auto.List.52.41da1c'))
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
      ElMessage.success(t('auto.List.53.8a225e'))
    } else {
      await wecomAccountApi.createAccount(payload)
      ElMessage.success(t('auto.List.54.482ecf'))
    }
    createDialogVisible.value = false
    loadData()
  } catch (err) {
    ElMessage.error(t('auto.List.55.40f902') + (err.message || t('auto.List.56.974e74')))
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
    await ElMessageBox.confirm(`确认删除企业微信账号「${account.corp_id}」？该操作不可恢复。`, t('auto.List.57.50eaf9'), {
      type: 'warning',
      confirmButtonText: t('auto.List.31.2f4aad'),
      cancelButtonText: t('auto.List.44.625fb2')
    })
  } catch {
    return
  }
  try {
    await wecomAccountApi.deleteAccount(account.id)
    ElMessage.success(t('auto.List.58.0007d1'))
    loadData()
  } catch (err) {
    ElMessage.error(t('auto.List.59.ad23f0') + (err.message || t('auto.List.56.974e74')))
  }
};

const healthScoreClass = computed(() => {
  const score = summary.value?.avg_score ?? 0
  if (score >= 80) return 'text-success'
  if (score >= 60) return 'text-warning'
  return 'text-danger'
})

const loginStateText = (state) => {
  const map = { online: t('auto.List.60.68905c'), offline: t('auto.List.61.50d4a8'), banned: t('auto.List.62.b08d38') }
  return map[state] || state
}

const loginStateTag = (state) => {
  const map = { online: 'success', offline: 'info', banned: 'danger' }
  return map[state] || 'info'
}

const riskLevelText = (level) => {
  const map = { normal: t('auto.List.63.fd6e80'), warning: t('auto.List.64.900c70'), critical: t('auto.List.65.540f84'), banned: t('auto.List.62.b08d38') }
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
    ElMessage.error(t('auto.List.66.779379') + (err.message || t('auto.List.56.974e74')))
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
    ElMessage.error(t('auto.List.67.4aad39') + (err.message || t('auto.List.56.974e74')))
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
    ElMessage.success(t('auto.List.68.02f667') + (count ? `，共 ${count} 条` : ''))
  } catch (e) {
    ElMessage.error(t('auto.List.69.a0f67c') + (e.message || e))
  } finally {
    actionLoading.value = ''
  }
};

const refreshAccount = async (account) => {
  actionLoading.value = 'refresh'
  try {
    await wecomAccountApi.refreshAccount(account.id)
    ElMessage.success(t('auto.List.70.9b353c'))
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.List.71.df570e') + (e.message || e))
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
