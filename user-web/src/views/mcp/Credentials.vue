<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Credentials.36c8f50b') }}</span>
          <el-button type="primary" size="small" @click="openCreate">{{ $t('auto.Credentials.97b0823e') }}</el-button>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('auto.Credentials.c02c86a1')"
        style="margin-bottom: 12px"
      />

      <el-table v-loading="loading" :data="rows" border size="small">
        <el-table-column label="Client ID" width="200">
          <template #default="{ row }">
            <span class="mono">{{ row.client_id }}</span>
            <el-button link size="small" @click="copyText(row.client_id, 'Client ID')">{{ $t('auto.Credentials.79d3abe9') }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="name" :label="$t('auto.Credentials.2432b575')" min-width="120" />
        <el-table-column :label="$t('auto.Credentials.3fea7ca7')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.1fa23f4d')" width="170">
          <template #default="{ row }">
            <el-tag v-if="row.expires_at" :type="isExpired(row.expires_at) ? 'danger' : 'warning'" size="small">
              {{ fmtTime(row.expires_at) }}{{ isExpired(row.expires_at) ? '（已过期）' : '' }}
            </el-tag>
            <el-tag v-else type="info" size="small">{{ $t('auto.Credentials.718ab6a0') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.adf2d72a')" width="110">
          <template #default="{ row }">
            <el-tag v-if="row.owner_user_id" size="small">用户 #{{ row.owner_user_id }}</el-tag>
            <span v-else class="muted">{{ $t('auto.Credentials.906ad18b') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.e16723d5')" min-width="160">
          <template #default="{ row }">
            <el-tag v-for="t in (row.allowed_tools || '').split(',').filter(Boolean)" :key="t" size="small" style="margin-right: 4px">
              {{ t }}
            </el-tag>
            <span v-if="!row.allowed_tools" class="muted">{{ $t('auto.Credentials.c9fc9997') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_used_at" :label="$t('auto.Credentials.7354e633')" width="160">
          <template #default="{ row }">{{ fmtTime(row.last_used_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.2b6bc0f2')" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" plain @click="openKeys(row)">{{ $t('auto.Credentials.cdb81cf6') }}</el-button>
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.Credentials.95b351c8') }}</el-button>
            <el-button size="small" @click="toggle(row)">{{ row.enabled ? '停用' : '启用' }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.Credentials.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 签发弹窗 -->
    <el-dialog v-model="createVisible" :title="$t('auto.Credentials.b9df5a0e')" width="520px">
      <el-form label-width="110px">
        <el-form-item :label="$t('auto.Credentials.2432b575')">
          <el-input v-model="createForm.name" :placeholder="$t('auto.Credentials.bf9eae08')" maxlength="100" />
        </el-form-item>
        <el-form-item label="Client ID">
          <el-input v-model="createForm.client_id" :placeholder="$t('auto.Credentials.0861d929')" maxlength="64" />
        </el-form-item>
        <el-form-item :label="$t('auto.Credentials.adf2d72a')">
          <el-input v-model="createForm.owner_user_id" :placeholder="$t('auto.Credentials.7804d1a3')" style="width: 200px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Credentials.1fa23f4d')">
          <el-radio-group v-model="createForm.expiryMode">
            <el-radio value="default">{{ $t('auto.Credentials.a0beab2c') }}</el-radio>
            <el-radio value="never">{{ $t('auto.Credentials.718ab6a0') }}</el-radio>
            <el-radio value="custom">{{ $t('auto.Credentials.f1d4ff50') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="createForm.expiryMode === 'custom'" :label="$t('auto.Credentials.1fa23f4d')">
          <el-date-picker
            v-model="createForm.expires_at"
            type="datetime"
            :placeholder="$t('auto.Credentials.0b52bcc5')"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ $t('auto.Credentials.625fb26b') }}</el-button>
        <el-button type="primary" :loading="creating" @click="doCreate">{{ $t('auto.Credentials.8a681e4f') }}</el-button>
      </template>
    </el-dialog>

    <!-- 编辑弹窗（改 ID / 续期 / 工具范围） -->
    <el-dialog v-model="editVisible" :title="$t('auto.Credentials.10aab08b')" width="520px">
      <el-form label-width="110px">
        <el-form-item label="Client ID">
          <el-input v-model="editForm.client_id" maxlength="64" />
        </el-form-item>
        <el-form-item :label="$t('auto.Credentials.1fa23f4d')">
          <el-radio-group v-model="editForm.expiryMode">
            <el-radio value="keep">{{ $t('auto.Credentials.5cfe2a94') }}</el-radio>
            <el-radio value="custom">{{ $t('auto.Credentials.ad13ef5d') }}</el-radio>
            <el-radio value="never">{{ $t('auto.Credentials.718ab6a0') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="editForm.expiryMode === 'custom'" :label="$t('auto.Credentials.ad13ef5d')">
          <el-date-picker
            v-model="editForm.expires_at"
            type="datetime"
            :placeholder="$t('auto.Credentials.952d9daf')"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Credentials.e16723d5')">
          <el-input v-model="editForm.allowed_tools" :placeholder="$t('auto.Credentials.da5711d8')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Credentials.adf2d72a')">
          <el-input v-model="editForm.owner_user_id" :placeholder="$t('auto.Credentials.7804d1a3')" style="width: 200px" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">{{ $t('auto.Credentials.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="doEdit">{{ $t('auto.Credentials.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    <!-- 签发结果（明文只显示一次） -->
    <el-dialog v-model="issuedVisible" :title="$t('auto.Credentials.81321fea')" width="520px" :close-on-click-modal="false">
      <el-alert type="warning" :closable="false" show-icon :title="$t('auto.Credentials.38569600')" style="margin-bottom: 12px" />
      <el-descriptions :column="1" border>
        <el-descriptions-item label="Client ID">{{ issued.client_id }}</el-descriptions-item>
        <el-descriptions-item label="API Key">
          <span class="mono">{{ issued.api_key }}</span>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="copyAll">{{ $t('auto.Credentials.c4519cdc') }}</el-button>
      </template>
    </el-dialog>

    <!-- 钥匙管理抽屉（1:N）：一个 Client ID 名下多把 API Key -->
    <el-drawer v-model="keysVisible" :title="`密钥管理 — ${keysCred.client_id || ''}`" size="720px">
      <div style="margin-bottom: 12px">
        <el-button type="primary" size="small" @click="openKeyCreate">{{ $t('auto.Credentials.0274e86e') }}</el-button>
        <span class="muted" style="margin-left: 8px">{{ $t('auto.Credentials.f82fc56e') }}</span>
      </div>
      <el-table v-loading="keysLoading" :data="keys" border size="small">
        <el-table-column prop="name" :label="$t('auto.Credentials.d7ec2d3f')" min-width="110" show-overflow-tooltip />
        <el-table-column label="API Key" min-width="230">
          <template #default="{ row }">
            <span v-if="revealed[row.id]" class="mono">{{ revealed[row.id] }}</span>
            <span v-else class="mono">{{ row.key_masked }}</span>
            <el-button
              v-if="row.can_reveal"
              link
              size="small"
              @click="toggleReveal(row)"
            >{{ revealed[row.id] ? '隐藏' : '显示' }}</el-button>
            <el-tag v-else size="small" type="info">{{ $t('auto.Credentials.6847524c') }}</el-tag>
            <el-button link size="small" @click="copyKey(row)">{{ $t('auto.Credentials.79d3abe9') }}</el-button>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.3fea7ca7')" width="90">
          <template #default="{ row }">
            <el-tag v-if="row.expired" type="danger" size="small">{{ $t('auto.Credentials.4d5ccd01') }}</el-tag>
            <el-tag v-else-if="row.status === 'active'" type="success" size="small">{{ $t('auto.Credentials.7854b52a') }}</el-tag>
            <el-tag v-else-if="row.status === 'disabled'" type="warning" size="small">{{ $t('auto.Credentials.710ad08b') }}</el-tag>
            <el-tag v-else-if="row.status === 'suspended'" type="warning" size="small">{{ $t('auto.Credentials.a2d930fd') }}</el-tag>
            <el-tag v-else type="info" size="small">{{ $t('auto.Credentials.69b0f684') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.1fa23f4d')" width="160">
          <template #default="{ row }">
            <span v-if="row.expires_at">{{ fmtTime(row.expires_at) }}</span>
            <span v-else class="muted">{{ $t('auto.Credentials.0f83a732') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.7354e633')" width="150">
          <template #default="{ row }">{{ fmtTime(row.last_used_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.Credentials.2b6bc0f2')" width="210" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'disabled'" size="small" type="success" plain @click="keyEnable(row)">{{ $t('auto.Credentials.7854b52a') }}</el-button>
            <el-button v-if="row.status === 'active'" size="small" @click="keyDisable(row)">{{ $t('auto.Credentials.710ad08b') }}</el-button>
            <el-button v-if="row.status === 'active' || row.status === 'disabled'" size="small" type="warning" plain @click="keySuspend(row)">{{ $t('auto.Credentials.8d63ef38') }}</el-button>
            <el-button v-if="row.status !== 'revoked'" size="small" type="warning" plain @click="keyRevoke(row)">{{ $t('auto.Credentials.5c56a889') }}</el-button>
            <el-button size="small" type="danger" plain @click="keyRemove(row)">{{ $t('auto.Credentials.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <!-- 新建密钥弹窗 -->
    <el-dialog v-model="keyCreateVisible" :title="$t('auto.Credentials.22615349')" width="480px" append-to-body>
      <el-form label-width="90px">
        <el-form-item :label="$t('auto.Credentials.d7ec2d3f')">
          <el-input v-model="keyForm.name" :placeholder="$t('auto.Credentials.df446027')" maxlength="100" />
        </el-form-item>
        <el-form-item :label="$t('auto.Credentials.1fa23f4d')">
          <el-radio-group v-model="keyForm.expiryMode">
            <el-radio value="default">{{ $t('auto.Credentials.a0beab2c') }}</el-radio>
            <el-radio value="never">{{ $t('auto.Credentials.718ab6a0') }}</el-radio>
            <el-radio value="custom">{{ $t('auto.Credentials.f1d4ff50') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="keyForm.expiryMode === 'custom'" :label="$t('auto.Credentials.1fa23f4d')">
          <el-date-picker
            v-model="keyForm.expires_at"
            type="datetime"
            :placeholder="$t('auto.Credentials.0b52bcc5')"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="keyCreateVisible = false">{{ $t('auto.Credentials.625fb26b') }}</el-button>
        <el-button type="primary" :loading="keyCreating" @click="doKeyCreate">{{ $t('auto.Credentials.8a681e4f') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { mcpCredentialApi, mcpKeyApi } from '@/api/mcpCredential'

const { t } = useI18n()

const rows = ref([])
const loading = ref(false)
const createVisible = ref(false)
const editVisible = ref(false)
const issuedVisible = ref(false)
const creating = ref(false)
const saving = ref(false)
const createForm = ref({})
const editForm = ref({})
const issued = ref({})

// —— 钥匙抽屉（1:N）——
const keysVisible = ref(false)
const keysLoading = ref(false)
const keysCred = ref({})
const keys = ref([])
const keyCreateVisible = ref(false)
const keyCreating = ref(false)
const keyForm = ref({})
// revealed[keyId] 仅存于内存：显示/隐藏是"拉明文↔收起"的开关，
// 刷新抽屉即清空；服务端对每次 reveal 都落管理审计。
const revealed = reactive({})

const fmtTime = (v) => (v ? String(v).replace('T', ' ').slice(0, 19) : '—')
const isExpired = (v) => v && new Date(v).getTime() < Date.now()

const toISO = (v) => {
  if (!v) return undefined
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return undefined
  return d.toISOString()
}

const load = async () => {
  loading.value = true
  try {
    const res = await mcpCredentialApi.list()
    rows.value = res?.data ?? res ?? []
  } finally {
    loading.value = false
  }
}

const copyText = async (text, label) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(`${label || t('auto.Credentials.2d711b09')}已复制`)
  } catch {
    ElMessage.info(t('auto.Credentials.063c173b'))
  }
}

const openCreate = () => {
  createForm.value = { name: '', client_id: '', owner_user_id: '', expiryMode: 'default', expires_at: undefined }
  createVisible.value = true
}

const openEdit = (row) => {
  editForm.value = {
    id: row.id,
    client_id: row.client_id,
    allowed_tools: row.allowed_tools || '',
    owner_user_id: row.owner_user_id ?? '',
    expiryMode: row.expires_at ? 'keep' : 'never',
    expires_at: undefined
  }
  editVisible.value = true
}

const doCreate = async () => {
  creating.value = true
  try {
    const payload = { name: createForm.value.name }
    if (createForm.value.client_id) payload.client_id = createForm.value.client_id.trim()
    if (createForm.value.owner_user_id) payload.owner_user_id = Number(createForm.value.owner_user_id)
    if (createForm.value.expiryMode === 'never') {
      // 永不过期：显式 no_expiry，后端不落默认 90 天
      payload.no_expiry = true
    } else if (createForm.value.expiryMode === 'custom' && createForm.value.expires_at) {
      payload.expires_at = toISO(createForm.value.expires_at)
    }
    const res = await mcpCredentialApi.create(payload)
    // 注意：@/utils/request 的响应拦截器已经把 {code,data} 里的 data 拆出来了，
    // res 本身就是 {id, client_id, api_key, api_hint}。这里再写 res?.data ?? {} 会
    // 在 res.data 为 undefined 时兜成空对象 ⇒ 弹窗与「复制并关闭」印出
    // "Client ID: undefined / API Key: undefined"（真实报障）。保留 res 兜底。
    issued.value = res?.data ?? res ?? {}
    createVisible.value = false
    issuedVisible.value = true
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Credentials.c1e45431'))
  } finally {
    creating.value = false
  }
}

const doEdit = async () => {
  saving.value = true
  try {
    const payload = {}
    if (editForm.value.client_id) payload.client_id = editForm.value.client_id.trim()
    if (editForm.value.owner_user_id) payload.owner_user_id = Number(editForm.value.owner_user_id)
    payload.allowed_tools = editForm.value.allowed_tools || ''
    if (editForm.value.expiryMode === 'custom' && editForm.value.expires_at) {
      payload.expires_at = toISO(editForm.value.expires_at)
    } else if (editForm.value.expiryMode === 'never') {
      // 显式清除过期（永不过期）；与"保持当前"（keep）由后端布尔区分
      payload.clear_expires = true
    }
    const res = await mcpCredentialApi.update(editForm.value.id, payload)
    ElMessage.success(res?.message || t('auto.Credentials.f8dfedcd'))
    editVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Credentials.6de920b4'))
  } finally {
    saving.value = false
  }
}

const copyAll = async () => {
  const text = `Client ID: ${issued.value.client_id}\nAPI Key: ${issued.value.api_key}`
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(t('auto.Credentials.52e6abbe'))
  } catch {
    ElMessage.info(t('auto.Credentials.6170576c'))
  }
  issuedVisible.value = false
}

const toggle = async (row) => {
  await mcpCredentialApi.update(row.id, { enabled: !row.enabled })
  ElMessage.success(row.enabled ? t('auto.Credentials.69b0f684') : t('auto.Credentials.53ace430'))
  await load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确认删除凭证 ${row.client_id}？名下全部 API Key 将一并删除，使用它的 Skill/Agent 立即失效。`, t('auto.Credentials.50eaf94d'), { type: 'warning' })
  await mcpCredentialApi.remove(row.id)
  ElMessage.success(t('auto.Credentials.5cc23262'))
  await load()
}

// —— 钥匙管理 ——

const loadKeys = async () => {
  keysLoading.value = true
  try {
    const res = await mcpKeyApi.list(keysCred.value.id)
    const data = res?.data ?? res ?? {}
    keys.value = data.list ?? []
  } finally {
    keysLoading.value = false
  }
}

const openKeys = (row) => {
  keysCred.value = row
  Object.keys(revealed).forEach((k) => delete revealed[k])
  keysVisible.value = true
  loadKeys()
}

const toggleReveal = async (row) => {
  if (revealed[row.id]) {
    delete revealed[row.id]
    return
  }
  try {
    const res = await mcpKeyApi.reveal(keysCred.value.id, row.id)
    const data = res?.data ?? res ?? {}
    if (data.api_key) {
      revealed[row.id] = data.api_key
    } else {
      ElMessage.error(t('auto.Credentials.e1290069'))
    }
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Credentials.3792273b'))
  }
}

const copyKey = async (row) => {
  let plain = revealed[row.id]
  if (!plain) {
    if (!row.can_reveal) {
      ElMessage.info(t('auto.Credentials.dd502bf1'))
      return
    }
    try {
      const res = await mcpKeyApi.reveal(keysCred.value.id, row.id)
      plain = (res?.data ?? res ?? {}).api_key
    } catch (e) {
      ElMessage.error(e?.response?.data?.message || t('auto.Credentials.e1290069'))
      return
    }
  }
  await copyText(plain, 'API Key')
}

const openKeyCreate = () => {
  keyForm.value = { name: '', expiryMode: 'default', expires_at: undefined }
  keyCreateVisible.value = true
}

const doKeyCreate = async () => {
  keyCreating.value = true
  try {
    const payload = { name: keyForm.value.name }
    if (keyForm.value.expiryMode === 'never') {
      payload.no_expiry = true
    } else if (keyForm.value.expiryMode === 'custom' && keyForm.value.expires_at) {
      payload.expires_at = toISO(keyForm.value.expires_at)
    }
    const res = await mcpKeyApi.create(keysCred.value.id, payload)
    const data = res?.data ?? res ?? {}
    keyCreateVisible.value = false
    // 明文只完整显示这一次（此后走列表的显示/复制）。
    await ElMessageBox.alert(
      `<div style="word-break:break-all;font-family:monospace">${data.api_key || ''}</div>`,
      t('auto.Credentials.83eedb9e'),
      { dangerouslyUseHTMLString: true, confirmButtonText: t('auto.Credentials.87707098') }
    )
    await loadKeys()
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Credentials.c1e45431'))
  } finally {
    keyCreating.value = false
  }
}

const keyStatusAction = async (row, action, confirmText) => {
  await ElMessageBox.confirm(confirmText, t('auto.Credentials.e83a256e'), { type: 'warning' })
  await mcpKeyApi[action](keysCred.value.id, row.id)
  ElMessage.success(t('auto.Credentials.1b92092d'))
  await loadKeys()
}

const keyEnable = (row) => keyStatusAction(row, 'enable', `确认启用密钥「${row.name}」？`)
const keyDisable = (row) => keyStatusAction(row, 'disable', `确认禁用密钥「${row.name}」？禁用后该钥匙立即无法调用（可再启用）。`)

const keySuspend = (row) => keyStatusAction(row, 'suspend', `确认暂停密钥「${row.name}」？暂停后该钥匙立即无法调用（可再启用）。`)
const keyRevoke = (row) => keyStatusAction(row, 'revoke', `确认停用（吊销）密钥「${row.name}」？吊销为终态，不可重新启用，仅可删除。`)

const keyRemove = async (row) => {
  await ElMessageBox.confirm(`确认删除密钥「${row.name}」？该操作不可恢复（管理审计保留记录）。`, t('auto.Credentials.50eaf94d'), { type: 'warning' })
  await mcpKeyApi.remove(keysCred.value.id, row.id)
  ElMessage.success(t('auto.Credentials.5cc23262'))
  await loadKeys()
}

onMounted(load)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.mono {
  font-family: monospace;
  user-select: all;
}
.muted {
  color: #999;
  font-size: 12px;
}
</style>
