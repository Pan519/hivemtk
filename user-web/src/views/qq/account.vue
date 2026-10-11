<template>
  <div class="qq-bot-management">
    <div class="account-tabs">
      <div class="account-list-container">
        <div class="qq-webhook-hint">
          <el-alert
            type="info"
            :closable="false"
            show-icon
            :title="$t('auto.account.1.024611')"
            description="在 QQ 开放平台（q.qq.com）创建机器人后，将 AppID / AppSecret / BotSecret 填入下方表单，并在开放平台管理端把回调地址配置为：{你的域名}/api/webhook/qq/{账号ID}。群聊场景用户 @机器人 即可触发智能体回复（被动回复 5 分钟内限 5 条）。"
          />
        </div>
        <div class="account-search">
          <el-input
            v-model="searchKeyword"
            ::placeholder="$t('auto.account.2.9d9a5f')"
            prefix-icon="Search"
            clearable
            @clear="handleSearch"
            @input="handleSearch"
          />
          <div class="action-buttons">
            <el-button type="primary" @click="handleAdd">{{ $t('添加机器人') }}</el-button>
            <el-button type="info" :loading="loading" @click="fetchAccounts">
              <el-icon><Refresh /></el-icon>
              <span>{{ $t('刷新') }}</span>
            </el-button>
          </div>
        </div>

        <div v-loading="loading" class="account-list">
          <el-table :data="filteredAccounts" style="width: 100%" border>
            <el-table-column prop="account_name" ::label="$t('auto.account.3.5f0720')" min-width="150" />
            <el-table-column prop="app_id" label="AppID" min-width="140" />
            <el-table-column prop="app_secret_masked" label="AppSecret" min-width="150" />
            <el-table-column prop="webhook_url" label="Webhook URL" min-width="230" show-overflow-tooltip />
            <el-table-column label="Webhook" width="100" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.webhook_enabled ? 'success' : 'info'" effect="plain">
                  {{ scope.row.webhook_enabled ? '已配置' : '未配置' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.account.4.72fbb9')" width="100" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.ai_agent_enabled ? 'success' : 'info'" effect="plain">
                  {{ scope.row.ai_agent_enabled ? '已启用' : '未启用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.account.5.3fea7c')" width="90" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'" effect="plain">
                  {{ scope.row.status === 1 ? '正常' : '停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="last_error_msg" :label="$t('auto.account.6.a2b83d')" min-width="160" show-overflow-tooltip />
            <el-table-column :label="$t('auto.account.7.2b6bc0')" width="380" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.account.24.95b351') }}</el-button>
                <el-button size="small" type="warning" @click="handleTestSend(scope.row)">{{ $t('auto.account.25.e1c652') }}</el-button>
                <el-button size="small" type="success" :loading="verifyingId === scope.row.id" @click="handleVerifyCallback(scope.row)">{{ $t('auto.account.26.bb81f2') }}</el-button>
                <el-button size="small" type="primary" @click="openBindingDialog(scope.row)">{{ $t('auto.account.27.c9e87d') }}</el-button>
                <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.account.28.2f4aad') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="!loading && filteredAccounts.length === 0" class="empty-data">
          <el-empty :description="$t('auto.account.8.a315d2')" />
        </div>
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      ::title="$t('auto.account.9.801f34')"
      width="640px"
      :close-on-click-modal="false"
    >
      <el-form ref="accountFormRef" :model="accountForm" :rules="rules" label-width="130px">
        <el-form-item :label="$t('auto.account.10.17327f')" prop="account_name">
          <el-input v-model="accountForm.account_name" :placeholder="$t('auto.account.11.db1e20')" />
        </el-form-item>
        <el-form-item label="AppID" prop="app_id">
          <el-input v-model="accountForm.app_id" :placeholder="$t('auto.account.12.c14f2e')" />
        </el-form-item>
        <el-form-item label="AppSecret" prop="app_secret">
          <el-input
            v-model="accountForm.app_secret"
            ::placeholder="$t('auto.account.13.1c9593')"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="BotSecret" prop="webhook_secret">
          <el-input
            v-model="accountForm.webhook_secret"
            ::placeholder="$t('auto.account.14.9edfaf')"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="Webhook URL" prop="webhook_url">
          <el-input
            v-model="accountForm.webhook_url"
            ::placeholder="$t('auto.account.15.b18f51')"
            :disabled="webhookUrlLocked"
          >
            <template v-if="accountForm.webhook_url" #append>
              <el-button @click="copyText(accountForm.webhook_url)">{{ $t('auto.account.29.79d3ab') }}</el-button>
            </template>
          </el-input>
          <div class="form-hint-block">
            <el-text type="info" size="small">
              {{ webhookUrlHint }}
            </el-text>
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.account.16.685aab')" prop="webhook_enabled">
          <el-switch v-model="accountForm.webhook_enabled" />
          <span class="form-hint">{{ $t('auto.account.30.6d2142') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.account.17.706367')" prop="ai_agent_enabled">
          <el-switch v-model="accountForm.ai_agent_enabled" />
          <span class="form-hint">{{ $t('auto.account.31.ff27ae') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.account.5.3fea7c')" prop="status">
          <el-radio-group v-model="accountForm.status">
            <el-radio :label="1">{{ $t('auto.account.32.fd6e80') }}</el-radio>
            <el-radio :label="2">{{ $t('auto.account.33.5c56a8') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.account.34.625fb2') }}</el-button>
          <el-button type="primary" :loading="submitting" @click="submitForm">{{ $t('auto.account.35.e83a25') }}</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="testDialogVisible" :title="$t('auto.account.18.880075')" width="520px">
      <el-form :model="testForm" label-width="110px">
        <el-form-item :label="$t('auto.account.19.b02b8b')">
          <el-radio-group v-model="testForm.target_type">
            <el-radio label="group">{{ $t('auto.account.36.df9abd') }}</el-radio>
            <el-radio label="c2c">{{ $t('auto.account.37.f994ae') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.account.20.fb48b9')">
          <el-input v-model="testForm.target_openid" :placeholder="$t('auto.account.21.5bc291')" />
        </el-form-item>
        <el-form-item :label="$t('auto.account.22.b87b77')">
          <el-input v-model="testForm.text" type="textarea" :rows="3" :placeholder="$t('auto.account.23.2477fb')" />
        </el-form-item>
      </el-form>
      <div class="test-hint">
        <el-text type="info" size="small">
          注意：QQ 主动消息需对方开启「允许主动发送」且受平台频控；推荐在群里 @机器人 通过被动回复验证链路。
        </el-text>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="testDialogVisible = false">{{ $t('auto.account.34.625fb2') }}</el-button>
          <el-button type="primary" :loading="testing" @click="submitTestSend">{{ $t('auto.account.38.1535fc') }}</el-button>
        </span>
      </template>
    </el-dialog>

    <AgentBindingDialog
      v-model="bindingDialogVisible"
      channel-type="qq"
      :account-id="bindingAccountId"
      :account-label="bindingAccountLabel"
      :account-enabled="bindingAccountEnabled"
    />
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listAccounts,
  createAccount,
  updateAccount,
  deleteAccount,
  testSend,
  verifyCallback
} from '@/api/qqBot'
import AgentBindingDialog from '@/components/AgentBindingDialog.vue'

const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const testing = ref(false)
const accounts = ref([])
const searchKeyword = ref('')
const verifyingId = ref(null)

const bindingDialogVisible = ref(false);
const bindingAccountId = ref('')
const bindingAccountLabel = ref('')
const bindingAccountEnabled = ref(true)

const openBindingDialog = (row) => {
  bindingAccountId.value = String(row.id)
  bindingAccountLabel.value = row.account_name || ''
  bindingAccountEnabled.value = row.status === 1
  bindingDialogVisible.value = true
}

const filteredAccounts = computed(() => {
  if (!searchKeyword.value) return accounts.value
  const kw = searchKeyword.value.toLowerCase()
  return accounts.value.filter(a => (a.account_name || '').toLowerCase().includes(kw))
})

const fetchAccounts = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await listAccounts()
    accounts.value = (res && res.list) || []
  } catch (e) {
    console.error(t('auto.account.39.5b5f6e'), e)
    ElMessage.error(i18n.global.t(t('auto.account.40.2701a1')))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {}

const dialogVisible = ref(false);
const dialogType = ref('add')
const accountFormRef = ref(null)
const accountForm = reactive({
  id: null,
  account_name: '',
  app_id: '',
  app_secret: '',
  webhook_secret: '',
  webhook_url: '',
  webhook_url_suggested: '',
  webhook_enabled: true,
  ai_agent_enabled: true,
  status: 1
})

// Webhook URL 锁定逻辑：公网 base 可用且推导值合法时只读展示推导值；否则可编辑
const webhookUrlLocked = computed(() => {
  return !!accountForm.webhook_url_suggested && !accountForm.id_dirty_url
})
const webhookUrlHint = computed(() => {
  if (accountForm.webhook_url_suggested) {
    return t('auto.account.41.ea94ea')
  }
  return t('auto.account.42.940181')
})

const rules = {
  account_name: [{ required: true, message: i18n.global.t(t('auto.account.43.1f3797')), trigger: 'blur' }],
  app_id: [{ required: true, message: i18n.global.t(t('auto.account.44.88a083')), trigger: 'blur' }],
  app_secret: [
    {
      validator: (rule, value, callback) => {
        if (dialogType.value === 'add' && !value) {
          return callback(new Error(t('auto.account.45.0e2227')))
        }
        callback()
      },
      trigger: 'blur'
    }
  ]
}

const handleAdd = () => {
  dialogType.value = 'add'
  Object.assign(accountForm, {
    id: null,
    account_name: '',
    app_id: '',
    app_secret: '',
    webhook_secret: '',
    webhook_url: '',
    webhook_url_suggested: '',
    webhook_enabled: true,
    ai_agent_enabled: true,
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogType.value = 'edit'
  Object.assign(accountForm, {
    id: row.id,
    account_name: row.account_name,
    app_id: row.app_id,
    app_secret: '',
    webhook_secret: '',
    webhook_url: row.webhook_url,
    webhook_url_suggested: row.webhook_url_suggested || '',
    webhook_enabled: row.webhook_enabled,
    ai_agent_enabled: row.ai_agent_enabled,
    status: row.status
  })
  dialogVisible.value = true
}

const submitForm = () => {
  accountFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const payload = {
        account_name: accountForm.account_name,
        app_id: accountForm.app_id,
        app_secret: accountForm.app_secret,
        webhook_secret: accountForm.webhook_secret,
        // 公网域名自动推导场景下不回传（空值=后端推导/保留原值）
        webhook_url: webhookUrlLocked.value ? '' : accountForm.webhook_url,
        webhook_enabled: accountForm.webhook_enabled,
        ai_agent_enabled: accountForm.ai_agent_enabled,
        status: accountForm.status
      }
      if (accountForm.id) {
        await updateAccount(accountForm.id, payload)
        ElMessage.success(i18n.global.t(t('auto.account.46.55aa63')))
      } else {
        await createAccount(payload)
        ElMessage.success(i18n.global.t(t('auto.account.47.04a691')))
      }
      dialogVisible.value = false
      fetchAccounts()
    } catch (e) {
      console.error(t('auto.account.48.4d5731'), e)
      ElMessage.error(t('auto.account.49.78d235') + (e.message || e))
    } finally {
      submitting.value = false
    }
  })
}

const copyText = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(t('auto.account.50.4fb42e'))
  } catch {
    ElMessage.warning(t('auto.account.51.5268f8'))
  }
}

const handleVerifyCallback = async (row) => {
  verifyingId.value = row.id
  try {
    const res = await verifyCallback(row.id)
    const sig = (res && res.signature) || ''
    ElMessageBox.alert(
      `签名前 32 位：${sig.slice(0, 32)}…\n\n${(res && res.hint) || t('auto.account.52.517b34')}`,
      t('auto.account.53.c7250e'),
      { confirmButtonText: t('auto.account.54.ce2695'), type: 'success' }
    ).catch(() => {})
  } catch (e) {
    ElMessage.error(t('auto.account.55.69a73b') + (e.message || e))
  } finally {
    verifyingId.value = null
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除机器人 ${row.account_name} 吗？`, t('auto.account.56.900c70'), {
    confirmButtonText: t('auto.account.57.38cf16'),
    cancelButtonText: t('auto.account.34.625fb2'),
    type: 'warning'
  })
    .then(async () => {
      try {
        await deleteAccount(row.id)
        ElMessage.success(i18n.global.t(t('auto.account.58.0007d1')))
        fetchAccounts()
      } catch (e) {
        ElMessage.error(t('auto.account.59.bc42bb') + (e.message || e))
      }
    })
    .catch(() => {})
}

const testDialogVisible = ref(false);
const testForm = reactive({ target_type: 'group', target_openid: '', text: '' })
let currentTestAccountId = null

const handleTestSend = (row) => {
  currentTestAccountId = row.id
  testForm.target_type = 'group'
  testForm.target_openid = ''
  testForm.text = t('auto.account.60.5ae844')
  testDialogVisible.value = true
}

const submitTestSend = async () => {
  if (!testForm.target_openid) {
    ElMessage.warning(t('auto.account.61.7e5c59'))
    return
  }
  if (!testForm.text.trim()) {
    ElMessage.warning(t('auto.account.62.fa967b'))
    return
  }
  testing.value = true
  try {
    await testSend(currentTestAccountId, {
      target_type: testForm.target_type,
      target_openid: testForm.target_openid,
      text: testForm.text
    })
    ElMessage.success(t('auto.account.63.9db9a7'))
    testDialogVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.account.64.9704ea') + (e.message || e))
  } finally {
    testing.value = false
  }
}

onMounted(fetchAccounts)
</script>

<style scoped>
.qq-bot-management {
  padding: 16px;
}
.account-search {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;
}
.action-buttons {
  display: flex;
  gap: 8px;
  white-space: nowrap;
}
.qq-webhook-hint {
  margin-bottom: 16px;
}
.form-hint {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
.form-hint-block {
  margin-top: 4px;
  line-height: 1.4;
}
.test-hint {
  padding: 0 8px 8px;
}
.empty-data {
  padding: 32px 0;
}
</style>
