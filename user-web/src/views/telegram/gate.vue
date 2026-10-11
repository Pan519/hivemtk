<template>
  <div class="tg-gate-management">
    <el-alert type="info" :closable="false" style="margin-bottom: 12px">
      <template #title>
        入群管控解决 Telegram 官方限制（Bot 无法主动私聊从未交互的用户）：
        <b>{{ $t('auto.gate.27.eae2e5') }}</b>私密群开启「申请加入」，用户申请后引导其向 Bot 发送 /start 验证，通过后自动批准进群；
        <b>{{ $t('auto.gate.28.4d62d5') }}</b>公开群进群即禁言，用户点击群内验证链接在 Bot 私聊完成 /start 后自动解禁。超时未验证自动清理。
      </template>
    </el-alert>

    <div class="toolbar">
      <el-select v-model="filterAccount" :placeholder="$t('auto.gate.1.4a0c18')" clearable style="width: 220px" @change="fetchGates">
        <el-option v-for="acc in accounts" :key="acc.id" :label="acc.account_name" :value="acc.id" />
      </el-select>
      <el-button type="primary" @click="handleAdd">{{ $t('auto.gate.29.c78824') }}</el-button>
      <el-button :loading="loading" @click="fetchGates">
        <el-icon><Refresh /></el-icon>
        <span>{{ $t('auto.gate.30.694fc5') }}</span>
      </el-button>
    </div>

    <el-table v-loading="loading" :data="gates" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column :label="$t('auto.gate.2.4b0ab7')" min-width="180">
        <template #default="scope">
          {{ scope.row.chat_title || scope.row.chat_id }}
        </template>
      </el-table-column>
      <el-table-column prop="chat_id" label="Chat ID" width="160" />
      <el-table-column :label="$t('auto.gate.3.c0fd5b')" width="180">
        <template #default="scope">
          <el-tag :type="scope.row.mode === 'join_request' ? 'warning' : 'primary'" effect="plain">
            {{ scope.row.mode === 'join_request' ? '方案 A 申请审批' : '方案 B 禁言解锁' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.gate.4.3fea7c')" width="100" align="center">
        <template #default="scope">
          <el-tag :type="scope.row.enabled ? 'success' : 'info'" effect="plain">
            {{ scope.row.enabled ? '已启用' : '未启用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="verify_ttl_min" :label="$t('auto.gate.5.e7f3e3')" width="130" align="center" />
      <el-table-column :label="$t('auto.gate.6.2b6bc0')" width="260" fixed="right">
        <template #default="scope">
          <el-button size="small" type="primary" @click="openMembers(scope.row)">{{ $t('auto.gate.31.255b50') }}</el-button>
          <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.gate.32.95b351') }}</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.gate.33.2f4aad') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" ::title="$t('auto.gate.7.75b507')" width="620px" :close-on-click-modal="false">
      <el-form ref="gateFormRef" :model="gateForm" :rules="rules" label-width="130px">
        <el-form-item :label="$t('auto.gate.8.f8a61f')" prop="account_id">
          <el-select v-model="gateForm.account_id" :placeholder="$t('auto.gate.9.30b052')" :disabled="dialogType === 'edit'" style="width: 100%">
            <el-option v-for="acc in accounts" :key="acc.id" :label="acc.account_name" :value="acc.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.10.bacd1e')" prop="chat_id">
          <el-input v-model="gateForm.chat_id" :placeholder="$t('auto.gate.11.f1b44e')" :disabled="dialogType === 'edit'" />
        </el-form-item>
        <el-form-item :label="$t('auto.gate.12.acb4c1')">
          <el-input v-model="gateForm.chat_title" :placeholder="$t('auto.gate.13.ff0757')" />
        </el-form-item>
        <el-form-item :label="$t('auto.gate.3.c0fd5b')" prop="mode">
          <el-radio-group v-model="gateForm.mode">
            <el-radio value="join_request">{{ $t('auto.gate.34.bc97c6') }}</el-radio>
            <el-radio value="mute_unlock">{{ $t('auto.gate.35.9570e2') }}</el-radio>
          </el-radio-group>
          <div class="form-hint">
            方案 A 需在群管理开启「申请加入」；方案 B 需要 Bot 拥有「限制成员」管理员权限
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.14.1ca3f4')">
          <el-switch v-model="gateForm.enabled" />
          <span class="form-hint">{{ $t('auto.gate.36.ae862e') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.5.e7f3e3')">
          <el-input-number v-model="gateForm.verify_ttl_min" :min="1" :max="1440" />
          <span class="form-hint">{{ $t('auto.gate.37.85940a') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.gate.15.cc7752')">
          <el-input v-model="gateForm.welcome_msg" type="textarea" :rows="3" :placeholder="$t('auto.gate.16.3c4d45')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.gate.38.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">{{ $t('auto.gate.39.be5fbb') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="membersVisible" :title="`成员验证台账 - ${currentGate ? (currentGate.chat_title || currentGate.chat_id) : ''}`" width="860px">
      <el-tabs v-model="memberStatusTab" @tab-change="fetchMembers">
        <el-tab-pane :label="$t('auto.gate.17.a8b0c2')" name="" />
        <el-tab-pane :label="$t('auto.gate.18.7c8332')" name="pending" />
        <el-tab-pane :label="$t('auto.gate.19.995dcd')" name="restricted" />
        <el-tab-pane :label="$t('auto.gate.20.75e96b')" name="approved" />
        <el-tab-pane :label="$t('auto.gate.21.3d65a4')" name="kicked" />
      </el-tabs>
      <el-table v-loading="membersLoading" :data="members" border size="small">
        <el-table-column prop="user_id" :label="$t('auto.gate.22.e53376')" width="150" />
        <el-table-column prop="full_name" :label="$t('auto.gate.23.23eb0e')" width="140" />
        <el-table-column prop="username" :label="$t('auto.gate.24.819767')" width="140">
          <template #default="scope">{{ scope.row.username ? '@' + scope.row.username : '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.gate.4.3fea7c')" width="110" align="center">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.join_status)" effect="plain">{{ statusLabel(scope.row.join_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.gate.25.c1a5fb')" width="110" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.authorized ? 'success' : 'info'" effect="plain">{{ scope.row.authorized ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.gate.26.8653cd')" min-width="160" />
        <el-table-column :label="$t('auto.gate.6.2b6bc0')" width="120" fixed="right">
          <template #default="scope">
            <el-button v-if="!scope.row.authorized" size="small" type="success" @click="handleAuthorize(scope.row)">{{ $t('auto.gate.40.befab2') }}</el-button>
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
  account_id: [{ required: true, message: t('auto.gate.41.feecf4'), trigger: 'change' }],
  chat_id: [{ required: true, message: t('auto.gate.42.e46a59'), trigger: 'blur' }],
  mode: [{ required: true, message: t('auto.gate.43.060a33'), trigger: 'change' }]
}

const statusLabel = (s) => ({ pending: t('auto.gate.18.7c8332'), restricted: t('auto.gate.19.995dcd'), approved: t('auto.gate.20.75e96b'), kicked: t('auto.gate.21.3d65a4') }[s] || s)
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
      ElMessage.success(t('auto.gate.44.04a691'))
    } else {
      await updateGate(gateForm.id, payload)
      ElMessage.success(t('auto.gate.45.55aa63'))
    }
    dialogVisible.value = false
    fetchGates()
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确定删除群组「${row.chat_title || row.chat_id}」的管控配置？`, t('auto.gate.46.02d981'), { type: 'warning' })
  await deleteGate(row.id)
  ElMessage.success(t('auto.gate.47.5cc232'))
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
  ElMessage.success(t('auto.gate.20.75e96b'))
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
