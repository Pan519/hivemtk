<template>
  <div class="tg-bot-management">

    <!-- ===== 操作说明卡片（首次使用必读）===== -->
    <el-card shadow="never" class="setup-guide-card">
      <template #header>
        <div class="guide-header">
          <span class="guide-title">{{ $t('auto.account.2a60fde4') }}</span>
          <el-tag effect="plain" size="small">{{ $t('auto.account.7177383c') }}</el-tag>
        </div>
      </template>

      <el-steps :active="0" direction="horizontal" finish-status="success" simple>
        <el-step :title="$t('auto.account.8dbf13a5')" :description="$t('auto.account.204d1d67')" />
        <el-step :title="$t('auto.account.5f8e7728')" :description="$t('auto.account.ff5b7773')" />
        <el-step :title="$t('auto.account.77363fbf')" :description="$t('auto.account.17e4b67c')" />
        <el-step :title="$t('auto.account.a8d90732')" :description="$t('auto.account.08049a53')" />
      </el-steps>

      <el-collapse class="guide-collapse" accordion>

        <!-- Step 1: 创建 Bot -->
        <el-collapse-item :title="$t('auto.account.df422480')" name="s1">
          <div class="guide-content">
            <ol>
              <li>{{ $t('auto.account.6949961b') }}<code>@BotFather</code></li>
              <li>{{ $t('auto.account.1535fcfa') }}<code>/newbot</code>{{ $t('auto.account.da97ca90') }}</li>
              <li>{{ $t('auto.account.1535fcfa') }}<code>/setdescription</code> <code>/setabouttext</code> <code>/setuserpic</code>{{ $t('auto.account.575b0896') }}</li>
              <li><strong>{{ $t('auto.account.d1c88548') }}<code>/setjoingroups</code>{{ $t('auto.account.52478bc3') }}<code>Enable</code></strong>
                <br>{{ $t('auto.account.c78f5dc8') }}</li>
              <li><strong>{{ $t('auto.account.d1c88548') }}<code>/setprivacy</code>{{ $t('auto.account.52478bc3') }}<code>Disable</code></strong>
                <br>{{ $t('auto.account.b859e187') }}</li>
              <li>{{ $t('auto.account.1535fcfa') }}<code>/token</code>{{ $t('auto.account.609af849') }}</li>
              <li>{{ $t('auto.account.1535fcfa') }}<code>/mybots</code>{{ $t('auto.account.18613ede') }}</li>
            </ol>
            <el-alert type="warning" :closable="false" show-icon>
              <template #title>
                Token 泄露了？马上到 @BotFather → /mybots → 选 Bot → 「Revoke current token」→ 重新生成即可。旧 Token 立刻失效。
              </template>
            </el-alert>
          </div>
        </el-collapse-item>

        <!-- Step 2: 拉进群+权限 -->
        <el-collapse-item :title="$t('auto.account.88a86bbd')" name="s2">
          <div class="guide-content">
            <ol>
              <li>{{ $t('auto.account.1f2c0718') }}</li>
              <li>{{ $t('auto.account.e95e7e8d') }}</li>
              <li>{{ $t('auto.account.d254245e') }}</li>
              <li><strong>{{ $t('auto.account.38e8089d') }}</strong>
                <ul>
                  <li>✅ <code>can_invite_users</code>{{ $t('auto.account.8349fc23') }}</li>
                  <li>✅ <code>can_restrict_members</code>{{ $t('auto.account.3dc150ee') }}</li>
                  <li>✅ <code>can_delete_messages</code>{{ $t('auto.account.565caed8') }}</li>
                  <li>✅ <code>can_send_messages</code>{{ $t('auto.account.4d1929b7') }}</li>
                </ul>
              </li>
            </ol>
            <el-alert type="info" :closable="false" show-icon>
              <template #title>
                如果是「方案 B 进群后禁言解锁」的公开群，建议群设置里开启「慢速模式 slow mode = 10 秒」，防止用户进群后 Bot 禁言前先发消息。
              </template>
            </el-alert>
          </div>
        </el-collapse-item>

        <!-- Step 3: Webhook vs Polling -->
        <el-collapse-item :title="$t('auto.account.9de60509')" name="s3">
          <div class="guide-content">
            <table class="guide-table">
              <thead>
                <tr><th>{{ $t('auto.account.f0789e79') }}</th><th>{{ $t('auto.account.62b37e99') }}</th><th>{{ $t('auto.account.d68f381b') }}</th><th>{{ $t('auto.account.c3318eaa') }}</th></tr>
              </thead>
              <tbody>
                <tr><td>🔗 Webhook</td><td>{{ $t('auto.account.5c9b494a') }}</td><td>{{ $t('auto.account.b591a123') }}</td><td>{{ $t('auto.account.638cab54') }}</td></tr>
                <tr><td>🔄 Polling</td><td>{{ $t('auto.account.ef9f9b2a') }}</td><td>{{ $t('auto.account.7b051f94') }}</td><td>{{ $t('auto.account.ceb1437f') }}</td></tr>
              </tbody>
            </table>
            <p><strong>{{ $t('auto.account.afef8faf') }}</strong>{{ $t('auto.account.d0db4e3a') }}<code>PUBLIC_BASE_URL</code>{{ $t('auto.account.90808bb3') }}<code>https://your-domain.com</code>{{ $t('auto.account.cab1e7e9') }}</p>
            <p>{{ $t('auto.account.f16ed611') }}</p>
          </div>
        </el-collapse-item>

        <!-- Step 4: 绑定智能体 -->
        <el-collapse-item :title="$t('auto.account.e704fe6c')" name="s4">
          <div class="guide-content">
            <p>{{ $t('auto.account.f8d4a30e') }}<strong>{{ $t('auto.account.e8f1834d') }}</strong>{{ $t('auto.account.165a45dc') }}<strong>{{ $t('auto.account.673c2033') }}</strong>{{ $t('auto.account.072c2726') }}</p>
            <p>{{ $t('auto.account.691814a7') }}</p>
            <ul>
              <li><code>hivemtk_bot</code>{{ $t('auto.account.427df66d') }}</li>
              <li><code>lovesanrenxing_bot</code>{{ $t('auto.account.6f22f2d2') }}</li>
              <li>{{ $t('auto.account.ed6a16b9') }}</li>
            </ul>
          </div>
        </el-collapse-item>

        <!-- 常见问题 -->
        <el-collapse-item :title="$t('auto.account.89a14970')" name="s5">
          <div class="guide-content">
            <div class="faq-item">
              <strong>{{ $t('auto.account.de361f70') }}</strong>
              <p>{{ $t('auto.account.e7609bcc') }}</p>
            </div>
            <div class="faq-item">
              <strong>{{ $t('auto.account.118bec06') }}</strong>
              <p>{{ $t('auto.account.9c0af4eb') }}</p>
            </div>
            <div class="faq-item">
              <strong>{{ $t('auto.account.cf2e5f66') }}</strong>
              <p>{{ $t('auto.account.b9c7ec18') }}</p>
            </div>
            <div class="faq-item">
              <strong>{{ $t('auto.account.2f228e78') }}</strong>
              <p>{{ $t('auto.account.4cab040a') }}<br>{{ $t('auto.account.52b84de5') }}<br>{{ $t('auto.account.1e60ad25') }}</p>
            </div>
          </div>
        </el-collapse-item>

      </el-collapse>
    </el-card>

    <div class="account-tabs">
      <div class="account-list-container">
        <div class="account-search">
          <el-input
            v-model="searchKeyword"
            ::placeholder="$t('auto.account.9d9a5fe2')"
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
            <el-table-column prop="account_name" ::label="$t('auto.account.5f0720d4')" min-width="160" />
            <el-table-column prop="bot_token_masked" label="Bot Token" min-width="200" />
            <el-table-column :label="$t('auto.account.fb7dd803')" width="130" align="center">
              <template #default="scope">
                <el-tooltip v-if="scope.row.last_error_msg && (scope.row.last_error_msg.includes('401') || scope.row.last_error_msg.includes('404'))" placement="top">
                  <template #content>{{ $t('auto.account.1d65374d') }}<br/>{{ $t('auto.account.403e67ff') }}<br/>{{ $t('auto.account.374279f0') }}</template>
                  <el-tag type="danger" effect="plain">{{ $t('auto.account.93fc3f61') }}</el-tag>
                </el-tooltip>
                <el-tag v-else-if="scope.row.last_error_msg" type="warning" effect="plain">{{ scope.row.last_error_msg.includes('429') ? '限流中' : '有错误' }}</el-tag>
                <el-tag v-else type="success" effect="plain">{{ $t('auto.account.fd6e80f1') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="webhook_url" label="Webhook URL" min-width="240" show-overflow-tooltip />
            <el-table-column label="Webhook" width="110" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.webhook_enabled ? 'success' : 'info'" effect="plain">
                  {{ scope.row.webhook_enabled ? '已注册' : '未注册' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.account.72fbb9b5')" width="110" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.ai_agent_enabled ? 'success' : 'info'" effect="plain">
                  {{ scope.row.ai_agent_enabled ? '已启用' : '未启用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.account.3fea7ca7')" width="100" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'" effect="plain">
                  {{ scope.row.status === 1 ? '正常' : '停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="last_error_msg" :label="$t('auto.account.a2b83df0')" min-width="180" show-overflow-tooltip />
            <el-table-column :label="$t('auto.account.2b6bc0f2')" width="400" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.account.95b351c8') }}</el-button>
                <el-button size="small" type="success" @click="handleRegisterWebhook(scope.row)">{{ $t('auto.account.fdb9b44e') }}</el-button>
                <el-button size="small" type="warning" @click="handleTestSend(scope.row)">{{ $t('auto.account.e1c652bc') }}</el-button>
                <el-button size="small" type="primary" @click="openBindingDialog(scope.row)">{{ $t('auto.account.c9e87d56') }}</el-button>
                <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.account.2f4aaddd') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="!loading && filteredAccounts.length === 0" class="empty-data">
          <el-empty :description="$t('auto.account.a4859af2')" />
        </div>
      </div>
    </div>

    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加机器人' : '编辑机器人'"
      width="640px"
      :close-on-click-modal="false"
    >
      <el-form ref="accountFormRef" :model="accountForm" :rules="rules" label-width="140px">
        <el-form-item :label="$t('auto.account.17327f9c')" prop="account_name">
          <el-input v-model="accountForm.account_name" :placeholder="$t('auto.account.db1e20cb')" />
        </el-form-item>
        <el-form-item label="Bot Token" prop="bot_token">
          <el-input
            v-model="accountForm.bot_token"
            :placeholder="dialogType === 'edit' ? '留空则保持原 Token 不变' : '形如 123456789:AAxxxx...（共 35 位）'"
            type="password"
            show-password
          />
          <div class="form-hint" style="font-size:12px;margin-top:4px">
            获取方式：Telegram 搜索 @BotFather → 发送 /mybots → 选择 Bot → API Token → 复制完整 Token（数字:35位字符）
          </div>
        </el-form-item>
        <!-- 测试 Token 按钮（创建模式必显；编辑模式仅当填了新 Token 时显） -->
        <el-form-item v-if="showTestToken" label=" ">
          <el-button type="success" :loading="testingToken" @click="handleTestToken">
            {{ testingToken ? '验证中...' : '🔍 测试 Token 有效性' }}
          </el-button>
          <span v-if="tokenTestResult" :style="{ marginLeft: '12px', fontSize: '13px', color: tokenTestResult.ok ? '#67c23a' : '#f56c6c' }">
            {{ tokenTestResult.ok ? `✅ Token 有效！Bot @${tokenTestResult.username || '(无username)'} 就绪` : '❌ ' + tokenTestResult.msg }}
          </span>
        </el-form-item>

        <!-- ========== 以下为"编辑模式只读展示区"（创建模式不显示） ========== -->
        <template v-if="dialogType === 'edit'">
          <el-divider content-position="left">{{ $t('auto.account.be6c194f') }}</el-divider>
          <el-form-item :label="$t('auto.account.735df6da')">
            <el-input ::model-value="$t('auto.account.8b7bb766')" disabled />
          </el-form-item>
          <el-form-item label="Webhook URL">
            <el-input ::model-value="$t('auto.account.bd191d58')" disabled>
              <template v-if="accountForm.webhook_url" #append>
                <el-button @click="copyText(accountForm.webhook_url)">{{ $t('auto.account.79d3abe9') }}</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Webhook Secret">
            <el-input :model-value="accountForm.webhook_secret ? '*** (已隐藏)' : '(后端自动生成)'" disabled />
          </el-form-item>
          <el-form-item :label="$t('auto.account.d2ac778b')">
            <el-tag :type="accountForm.webhook_enabled ? 'success' : 'info'" effect="plain">
              {{ accountForm.webhook_enabled ? '已注册' : '未注册' }}
            </el-tag>
          </el-form-item>
          <el-form-item :label="$t('auto.account.a1194f7e')">
            <el-tag :type="accountForm.ai_agent_enabled ? 'success' : 'info'" effect="plain">
              {{ accountForm.ai_agent_enabled ? '已启用' : '未启用' }}
            </el-tag>
          </el-form-item>
          <el-form-item :label="$t('auto.account.3fea7ca7')">
            <el-tag :type="accountForm.status === 1 ? 'success' : 'danger'" effect="plain">
              {{ accountForm.status === 1 ? '正常' : '停用' }}
            </el-tag>
          </el-form-item>
          <div class="form-hint" style="color:#909399;font-size:12px;margin-left:140px;margin-top:-8px">
            以上字段由后端自动管理，如需修改请使用列表上方的「注册Webhook」等按钮
          </div>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.account.625fb26b') }}</el-button>
          <el-button type="primary" :loading="submitting" @click="submitForm">
            {{ dialogType === 'add' ? '创建（自动配置）' : '保存' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog v-model="testDialogVisible" :title="$t('auto.account.88007514')" width="500px">
      <el-form :model="testForm" label-width="100px">
        <el-form-item :label="$t('auto.account.35f1ba0f')">
          <el-input v-model.number="testForm.chat_id" :placeholder="$t('auto.account.3dc3ac43')" />
        </el-form-item>
        <el-form-item :label="$t('auto.account.b87b7756')">
          <el-input v-model="testForm.text" type="textarea" :rows="3" :placeholder="$t('auto.account.2477fb6a')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="testDialogVisible = false">{{ $t('auto.account.625fb26b') }}</el-button>
          <el-button type="primary" :loading="testing" @click="submitTestSend">{{ $t('auto.account.1535fcfa') }}</el-button>
        </span>
      </template>
    </el-dialog>

    
    <AgentBindingDialog
      v-model="bindingDialogVisible"
      channel-type="telegram"
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
  registerWebhook,
  testSend
} from '@/api/telegram'
import AgentBindingDialog from '@/components/AgentBindingDialog.vue'

const { t } = useI18n()

const loading = ref(false)
const submitting = ref(false)
const testing = ref(false)
const testingToken = ref(false)
const tokenTestResult = ref(null) // { ok: bool, username: str, msg: str }
const accounts = ref([])
const searchKeyword = ref('')

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
    console.error(t('auto.account.07109179'), e)
    ElMessage.error(i18n.global.t(t('auto.account.2701a199')))
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
  bot_token: '',
  bot_username: '',
  webhook_url: '',
  webhook_secret: '',
  webhook_enabled: false,
  ai_agent_enabled: false,
  status: 1
})

// 是否显示"测试 Token"按钮：创建时始终显示；编辑时仅当填了新 Token 才显示
const showTestToken = computed(() => {
  if (dialogType.value === 'add') return !!accountForm.bot_token
  return !!accountForm.bot_token // 编辑时只要填了 token 就可测
})

const rules = {
  account_name: [{ required: true, message: i18n.global.t(t('auto.account.1f3797f1')), trigger: 'blur' }],
  bot_token: [
    {
      validator: (rule, value, callback) => {
        if (dialogType.value === 'add' && !value) {
          return callback(new Error(t('auto.account.25d7586d')))
        }
        if (value && !/^\d{6,10}:[A-Za-z0-9_-]{35}$/.test(value.trim())) {
          return callback(new Error('Token 格式不对：应为「数字:35位字母数字」，请到 @BotFather → /mybots → API Token 完整复制'))
        }
        callback()
      },
      trigger: 'blur'
    }
  ]
}

const handleAdd = () => {
  dialogType.value = 'add'
  tokenTestResult.value = null
  Object.assign(accountForm, {
    id: null,
    account_name: '',
    bot_token: '',
    bot_username: '',
    webhook_url: '',
    webhook_secret: '',
    webhook_enabled: false,
    ai_agent_enabled: false,
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogType.value = 'edit'
  tokenTestResult.value = null
  Object.assign(accountForm, {
    id: row.id,
    account_name: row.account_name,
    bot_token: '',
    bot_username: row.bot_username || '',
    webhook_url: row.webhook_url || '',
    webhook_secret: row.webhook_secret || '',
    webhook_enabled: !!row.webhook_enabled,
    ai_agent_enabled: !!row.ai_agent_enabled,
    status: row.status
  })
  dialogVisible.value = true
}

/**
 * 测试 Bot Token 有效性
 * 直接前端调用 Telegram getMe API（公开接口，无需后端代理）
 */
const handleTestToken = async () => {
  const token = (accountForm.bot_token || '').trim()
  if (!token) {
    ElMessage.warning(t('auto.account.8df062e2'))
    return
  }
  testingToken.value = true
  tokenTestResult.value = null
  try {
    // Telegram getMe 是公开 REST API，直接 fetch（CORS 友好）
    const resp = await fetch(`https://api.telegram.org/bot${token}/getMe`)
    const data = await resp.json()
    if (data && data.ok && data.result) {
      const u = data.result
      tokenTestResult.value = {
        ok: true,
        username: u.username,
        firstName: u.first_name,
        msg: ''
      }
      // 自动填充 bot_username 到表单
      if (u.username) {
        accountForm.bot_username = '@' + u.username
      }
      ElMessage.success(`Token 验证通过：@${u.username || '(无)'}`)
    } else {
      const desc = data?.description || 'Token 无效'
      tokenTestResult.value = { ok: false, msg: desc.replace(/^[^:]+:\s*/, '') }
      ElMessage.error(`Token 无效: ${desc}`)
    }
  } catch (err) {
    tokenTestResult.value = { ok: false, msg: t('auto.account.749733dd') }
    ElMessage.error(t('auto.account.98d1eff6') + err.message)
  } finally {
    testingToken.value = false
  }
}

const copyText = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(t('auto.account.4fb42e6e'))
  } catch {
    ElMessage.warning(t('auto.account.5268f875'))
  }
}

const submitForm = () => {
  accountFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      let payload
      if (accountForm.id) {
        // Update：只传需要改的字段；webhook_secret 留空=保留原值
        payload = {
          account_name: accountForm.account_name
        }
        if (accountForm.bot_token) {
          payload.bot_token = accountForm.bot_token
        }
        // webhook_enabled / ai_agent_enabled / status 由后端自动管理，
        // 编辑对话框不允许修改这些字段，所以 Update 时不传它们
      } else {
        // Create：只传 2 个必填字段，其余后端自动处理
        payload = {
          account_name: accountForm.account_name,
          bot_token: accountForm.bot_token
        }
      }
      if (accountForm.id) {
        await updateAccount(accountForm.id, payload)
        ElMessage.success(i18n.global.t(t('auto.account.55aa6366')))
      } else {
        const res = await createAccount(payload)
        ElMessage.success(i18n.global.t(t('auto.account.464f3ad4')))
        // 创建成功后立刻刷新列表，展示后端自动填充的字段
        fetchAccounts()
      }
      dialogVisible.value = false
      if (!accountForm.id) {
        fetchAccounts()
      }
    } catch (e) {
      console.error(t('auto.account.4d57314a'), e)
      ElMessage.error(t('auto.account.78d235cc') + (e.message || e))
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除机器人 ${row.account_name} 吗？`, t('auto.account.900c70fa'), {
    confirmButtonText: t('auto.account.38cf16f2'),
    cancelButtonText: t('auto.account.625fb26b'),
    type: 'warning'
  })
    .then(async () => {
      try {
        await deleteAccount(row.id)
        ElMessage.success(i18n.global.t(t('auto.account.0007d170')))
        fetchAccounts()
      } catch (e) {
        ElMessage.error(t('auto.account.bc42bb8f') + (e.message || e))
      }
    })
    .catch(() => {})
}

const handleRegisterWebhook = (row) => {
  // 前置自检：给出可操作的提示，避免打到 Telegram 才收到晦涩报错
  const problems = []
  if (row.last_error_msg && row.last_error_msg.includes('401')) {
    problems.push(t('auto.account.9f7c67cf'))
  }
  if (row.webhook_url && row.webhook_url.startsWith('http://')) {
    problems.push('Webhook URL 是 http:// 开头，Telegram 强制要求 https://，请在「编辑」里修改')
  }
  if (!row.webhook_url) {
    problems.push('Webhook URL 为空：请先在「编辑」对话框填入（或等待后端自动推导完成后再点注册）')
  }
  if (problems.length) {
    ElMessageBox.alert(problems.join('\n'), t('auto.account.22398359'), { type: 'warning', confirmButtonText: t('auto.account.ce26955a') })
    return
  }
  ElMessageBox.confirm(`确定要为机器人 ${row.account_name} 注册 Webhook 吗？`, t('auto.account.e83a256e'), {
    confirmButtonText: t('auto.account.38cf16f2'),
    cancelButtonText: t('auto.account.625fb26b'),
    type: 'info'
  })
    .then(async () => {
      try {
        await registerWebhook(row.id, { webhook_url: row.webhook_url })
        ElMessage.success(i18n.global.t('Webhook 注册成功'))
        fetchAccounts()
      } catch (e) {
        ElMessage.error('Webhook 注册失败：' + (e.message || e))
      }
    })
    .catch(() => {})
}

const testDialogVisible = ref(false);
const testForm = reactive({ chat_id: null, text: '' })
let currentTestAccountId = null

const handleTestSend = (row) => {
  currentTestAccountId = row.id
  testForm.chat_id = null
  testForm.text = t('auto.account.5ae8445d')
  testDialogVisible.value = true
}

const submitTestSend = async () => {
  if (!testForm.chat_id) {
    ElMessage.warning(i18n.global.t(t('auto.account.6d3eaa82')))
    return
  }
  if (!testForm.text) {
    ElMessage.warning(i18n.global.t(t('auto.account.fa967bd0')))
    return
  }
  testing.value = true
  try {
    await testSend(currentTestAccountId, { chat_id: testForm.chat_id, text: testForm.text })
    ElMessage.success(i18n.global.t(t('auto.account.9db9a7e3')))
    testDialogVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.account.9704ea75') + (e.message || e))
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  fetchAccounts()
})
</script>

<style lang="scss" scoped>
// 变量由 vite additionalData 全局注入(@import)，此处不得再 @use（dart-sass 要求 @use 先于一切规则）

.tg-bot-management {

  // ===== 操作说明卡片样式 =====
  .setup-guide-card {
    margin-bottom: 20px;
    border-left: 4px solid #409eff;

    .guide-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 12px;

      .guide-title {
        font-size: 15px;
        font-weight: 600;
        color: #303133;
      }
    }

    .guide-collapse {
      margin-top: 10px;
      border: none;

      .el-collapse-item__header {
        font-weight: 500;
        font-size: 14px;
        border-bottom: none;
        padding: 8px 0;
      }

      .el-collapse-item__wrap {
        border-bottom: none;
      }
    }

    .guide-content {
      font-size: 13px;
      line-height: 1.8;
      color: #606266;
      padding: 4px 0 12px;

      ol {
        padding-left: 20px;

        li {
          margin-bottom: 6px;
          code {
            background: #f4f4f5;
            border-radius: 3px;
            padding: 1px 5px;
            font-size: 12px;
            color: #e6a23c;
          }
        }

        ul {
          margin: 6px 0 6px 10px;
          padding-left: 18px;

          li {
            margin-bottom: 3px;
            code {
              color: #409eff;
              font-size: 12px;
            }
          }
        }
      }

      .faq-item {
        margin-bottom: 12px;
        padding: 8px 12px;
        background: #fafafa;
        border-radius: 6px;
        border-left: 3px solid #e6a23c;

        strong {
          font-size: 13px;
          color: #303133;
        }

        p {
          margin: 4px 0 0;
          font-size: 12px;
          color: #606266;
        }
      }

      .guide-table {
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 10px;
        font-size: 12px;

        th, td {
          border: 1px solid #ebeef5;
          padding: 6px 10px;
          text-align: left;
        }

        th {
          background: #f5f7fa;
          font-weight: 600;
        }
      }

      p {
        margin: 6px 0;

        code {
          background: #ecf5ff;
          color: #409eff;
          border-radius: 3px;
          padding: 1px 5px;
          font-size: 12px;
        }
      }
    }
  }

  .account-tabs {
    background-color: #fff;
    border-radius: 4px;
    box-shadow: $box-shadow-light;
  }

  .account-list-container {
    padding: 20px;

    .account-search {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;

      .el-input {
        width: 320px;
      }

      .action-buttons {
        display: flex;
        gap: 10px;
      }
    }

    .account-list {
      margin-bottom: 20px;
    }

    .empty-data {
      padding: 40px 0;
    }
  }
}

.form-hint {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}
</style>
