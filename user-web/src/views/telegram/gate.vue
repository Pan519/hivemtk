<template>
  <div class="tg-gate-management">
    <el-alert type="info" :closable="false" style="margin-bottom: 12px">
      <template #title>
        入群管控解决 Telegram 官方限制（Bot 无法主动私聊从未交互的用户）：
        <b>{{ $t('auto.gate.eae2e538') }}</b>私密群开启「申请加入」，用户申请后引导其向 Bot 发送 /start 验证，通过后自动批准进群；
        <b>{{ $t('auto.gate.4d62d576') }}</b>公开群进群即禁言，用户点击群内验证链接在 Bot 私聊完成 /start 后自动解禁。超时未验证自动清理。
      </template>
    </el-alert>

    <div class="toolbar">
      <el-select v-model="filterAccount" :placeholder="$t('auto.gate.4a0c18d4')" clearable style="width: 220px" @change="fetchGates">
        <el-option v-for="acc in accounts" :key="acc.id" :label="acc.account_name" :value="acc.id" />
      </el-select>
      <el-button type="primary" @click="handleAdd">{{ $t('auto.gate.c7882471') }}</el-button>
      <el-button :loading="loading" @click="fetchGates">
        <el-icon><Refresh /></el-icon>
        <span>{{ $t('auto.gate.694fc5ef') }}</span>
      </el-button>
    </div>

    <el-table v-loading="loading" :data="gates" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column :label="$t('auto.gate.4b0ab7ba')" min-width="180">
        <template #default="scope">
          {{ scope.row.chat_title || scope.row.chat_id }}
        </template>
      </el-table-column>
      <el-table-column prop="chat_id" label="Chat ID" width="160" />
      <el-table-column :label="$t('auto.gate.c0fd5bc1')" width="180">
        <template #default="scope">
          <el-tag :type="scope.row.mode === 'join_request' ? 'warning' : 'primary'" effect="plain">
            {{ scope.row.mode === 'join_request' ? '方案 A 申请审批' : '方案 B 禁言解锁' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.gate.3fea7ca7')" width="100" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.enabled ? 'success' : 'info'" effect="plain">
            {{ scope.row.enabled ? '已启用' : '未启用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="verify_ttl_min" :label="$t('auto.gate.e7f3e3e8')" width="130" align="center" />
      <el-table-column :label="$t('auto.gate.2b6bc0f2')" width="260" fixed="right">
        <template #default="scope">
          <el-button size="small" type="primary" @click="openMembers(scope.row)">{{ $t('auto.gate.255b5078') }}</el-button>
          <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.gate.95b351c8') }}</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.gate.2f4aaddd') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '添加群组管控' : '编辑群组管控'" width="620px" :close-on-click-modal="false">
      <el-form ref="gateFormRef" :model="gateForm" :rules="rules" label-width="130px">
        <el-form-item :label="$t('auto.gate.f8a61f3f')" prop="account_id">
          <el-select v-model="gateForm.account_id" :placeholder="$t('auto.gate.30b05276')" :disabled="dialogType === 'edit'" style="width: 100%">
            <el-option v-for="acc in accounts" :key="acc.id" :label="acc.account_name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.bacd1e85')" prop="chat_id">
          <el-input v-model="gateForm.chat_id" :placeholder="$t('auto.gate.f1b44ea4')" :disabled="dialogType === 'edit'" />
        </el-form-item>
        <el-form-item :label="$t('auto.gate.acb4c1a1')">
          <el-input v-model="gateForm.chat_title" :placeholder="$t('auto.gate.ff075714')" />
        </el-form-item>
        <el-form-item :label="$t('auto.gate.c0fd5bc1')" prop="mode">
          <el-radio-group v-model="gateForm.mode">
            <el-radio value="join_request">{{ $t('auto.gate.bc97c67a') }}</el-radio>
            <el-radio value="mute_unlock">{{ $t('auto.gate.9570e252') }}</el-radio>
          </el-radio-group>
          <div class="form-hint">
            方案 A 需在群管理开启「申请加入」；方案 B 需要 Bot 拥有「限制成员」管理员权限
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.1ca3f4c9')">
          <el-switch v-model="gateForm.enabled" />
          <span class="form-hint">{{ $t('auto.gate.ae862e43') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.e7f3e3e8')">
          <el-input-number v-model="gateForm.verify_ttl_min" :min="1" :max="1440" />
          <span class="form-hint">{{ $t('auto.gate.85940a96') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.cc7752ec')">
          <el-input v-model="gateForm.welcome_msg" type="textarea" :rows="3" :placeholder="$t('auto.gate.3c4d4577')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.gate.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">{{ $t('auto.gate.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="membersVisible" :title="`成员验证台账 - ${currentGate ? (currentGate.chat_title || currentGate.chat_id) : ''}`" width="860px">
      <el-tabs v-model="memberStatusTab" @tab-change="fetchMembers">
        <el-tab-pane :label="$t('auto.gate.a8b0c204')" name="" />
        <el-tab-pane :label="$t('auto.gate.7c8332aa')" name="pending" />
        <el-tab-pane :label="$t('auto.gate.995dcd67')" name="restricted" />
        <el-tab-pane :label="$t('auto.gate.75e96bf7')" name="approved" />
        <el-tab-pane :label="$t('auto.gate.3d65a49f')" name="kicked" />
      </el-tabs>
      <el-table v-loading="membersLoading" :data="members" border size="small">
        <el-table-column prop="user_id" :label="$t('auto.gate.e53376fa')" width="150" />
        <el-table-column prop="full_name" :label="$t('auto.gate.23eb0e60')" width="140" />
        <el-table-column prop="username" :label="$t('auto.gate.819767ad')" width="140">
          <template #default="scope">{{ scope.row.username ? '@' + scope.row.username : '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.gate.3fea7ca7')" width="110" align="center">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.join_status)" effect="plain">{{ statusLabel(scope.row.join_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.gate.c1a5fbd0')" width="110" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.authorized ? 'success' : 'info'" effect="plain">{{ scope.row.authorized ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.gate.8653cd3d')" min-width="160" />
        <el-table-column :label="$t('auto.gate.2b6bc0f2')" width="120" fixed="right">
          <template #default="scope">
            <el-button v-if="!scope.row.authorized" size="small" type="success" @click="handleAuthorize(scope.row)">{{ $t('auto.gate.befab2f3') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listAccounts, listGates, createGate, updateGate, deleteGate, listGateMembers, authorizeGateMember } from '@/api/telegram'

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const gates = ref([])
const accounts = ref([])
const filterAccount = ref(null)
const dialogVisible = ref(false)
const dialogType = ref('add')

const membersVisible = ref(false)
const membersLoading = ref(false)
const members = ref([])
const currentGate = ref(null)
const memberStatusTab = ref('')

const gateForm = reactive({
  id: null,
  account_id: null,
  chat_id: '',
  chat_title: '',
  mode: 'mute_unlock',
  enabled: true,
  verify_ttl_min: 10,
  welcome_msg: ''
})

const rules = {
  account_id: [{ required: true, message: t('auto.gate.feecf454'), trigger: 'change' }],
  chat_id: [{ required: true, message: t('auto.gate.e46a5909'), trigger: 'blur' }],
  mode: [{ required: true, message: t('auto.gate.060a3338'), trigger: 'change' }]
}

const statusLabel = (s) => ({ pending: t('auto.gate.7c8332aa'), restricted: t('auto.gate.995dcd67'), approved: t('auto.gate.75e96bf7'), kicked: t('auto.gate.3d65a49f') }[s] || s)
const statusTagType = (s) => ({ pending: 'warning', restricted: 'danger', approved: 'success', kicked: 'info' }[s] || 'info')

async function fetchAccounts() {
  const res = await listAccounts()
  accounts.value = (res && (res.list || res.data?.list)) || []
}

async function fetchGates() {
  loading.value = true
  try {
    const params = filterAccount.value ? { account_id: filterAccount.value } : {}
    const res = await listGates(params)
    gates.value = (res && (res.list || res.data?.list)) || []
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  dialogType.value = 'add'
  Object.assign(gateForm, { id: null, account_id: filterAccount.value || null, chat_id: '', chat_title: '', mode: 'mute_unlock', enabled: true, verify_ttl_min: 10, welcome_msg: '' })
  dialogVisible.value = true
}

function handleEdit(row) {
  dialogType.value = 'edit'
  Object.assign(gateForm, { id: row.id, account_id: row.account_id, chat_id: row.chat_id, chat_title: row.chat_title, mode: row.mode, enabled: row.enabled, verify_ttl_min: row.verify_ttl_min, welcome_msg: row.welcome_msg })
  dialogVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    const payload = { ...gateForm }
    if (dialogType.value === 'add') {
      await createGate(payload)
      ElMessage.success(t('auto.gate.04a691b3'))
    } else {
      await updateGate(gateForm.id, payload)
      ElMessage.success(t('auto.gate.55aa6366'))
    }
    dialogVisible.value = false
    fetchGates()
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定删除群组「${row.chat_title || row.chat_id}」的管控配置？`, t('auto.gate.02d9819d'), { type: 'warning' })
  await deleteGate(row.id)
  ElMessage.success(t('auto.gate.5cc23262'))
  fetchGates()
}

function openMembers(row) {
  currentGate.value = row
  memberStatusTab.value = ''
  membersVisible.value = true
  fetchMembers()
}

async function fetchMembers() {
  if (!currentGate.value) return
  membersLoading.value = true
  try {
    const params = { limit: 100 }
    if (memberStatusTab.value) params.status = memberStatusTab.value
    const res = await listGateMembers(currentGate.value.id, params)
    members.value = (res && (res.list || res.data?.list)) || []
  } finally {
    membersLoading.value = false
  }
}

async function handleAuthorize(row) {
  await authorizeGateMember(currentGate.value.id, row.id)
  ElMessage.success(t('auto.gate.75e96bf7'))
  fetchMembers()
}

onMounted(() => {
  fetchAccounts()
  fetchGates()
})
</script>

<style scoped>
.tg-gate-management { padding: 16px; }
.toolbar { display: flex; gap: 12px; margin-bottom: 12px; align-items: center; }
.form-hint { font-size: 12px; color: #909399; margin-left: 8px; }
</style>
