<template>
  <div class="dingtalk-app">
    <el-card class="card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.DingtalkAppAccount.1ef4b354') }}</span>
          <el-button type="primary" @click="openCreate">{{ $t('auto.DingtalkAppAccount.34374b71') }}</el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="account_name" :label="$t('auto.DingtalkAppAccount.d7ec2d3f')" min-width="120" />
        <el-table-column prop="app_key" label="AppKey" min-width="130" show-overflow-tooltip />
        <el-table-column prop="agent_id" label="AgentId" min-width="100" show-overflow-tooltip />
        <el-table-column :label="$t('auto.DingtalkAppAccount.ab343c20')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.inbound_enabled ? 'success' : 'info'" size="small">
              {{ row.inbound_enabled ? '已启用' : '未启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.DingtalkAppAccount.c9e87d56')" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag :type="row.ai_agent_id ? 'success' : 'info'" size="small">
              {{ row.ai_agent_id || '未绑定' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.DingtalkAppAccount.2b6bc0f2')" width="320" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openTest(row)">{{ $t('auto.DingtalkAppAccount.db06c78d') }}</el-button>
            <el-button link type="primary" @click="openBind(row)">{{ $t('auto.DingtalkAppAccount.c9e87d56') }}</el-button>
            <el-button link type="primary" @click="copyCallbackUrl(row)">{{ $t('auto.DingtalkAppAccount.e73c73c1') }}</el-button>
            <el-button link type="primary" @click="openEdit(row)">{{ $t('auto.DingtalkAppAccount.95b351c8') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ $t('auto.DingtalkAppAccount.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > pageSize"
        class="pager"
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        @current-change="handlePage"
      />
    </el-card>

    
    <el-dialog v-model="dialogVisible" ::title="$t('auto.DingtalkAppAccount.ab3ad6e0')" width="680px">
      <el-form :model="form" label-width="150px">
        <el-form-item :label="$t('auto.DingtalkAppAccount.d7ec2d3f')" prop="account_name">
          <el-input v-model="form.account_name" :placeholder="$t('auto.DingtalkAppAccount.17327f9c')" />
        </el-form-item>
        <el-form-item label="AppKey" prop="app_key">
          <el-input v-model="form.app_key" :placeholder="$t('auto.DingtalkAppAccount.57ba6e91')" />
        </el-form-item>
        <el-form-item label="AppSecret" prop="app_secret">
          <el-input v-model="form.app_secret" type="password" show-password :placeholder="$t('auto.DingtalkAppAccount.a9798ef9')" />
        </el-form-item>
        <el-form-item label="AgentId" prop="agent_id">
          <el-input v-model="form.agent_id" :placeholder="$t('auto.DingtalkAppAccount.93abc007')" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('auto.DingtalkAppAccount.f96de929') }}</el-divider>
        <el-form-item :label="$t('auto.DingtalkAppAccount.ecbc51e8')" prop="token">
          <el-input v-model="form.token" :placeholder="$t('auto.DingtalkAppAccount.8b7ea3f3')" />
          <div class="form-hint">{{ $t('auto.DingtalkAppAccount.61e4a51f') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.67f58612')" prop="aes_key">
          <el-input v-model="form.aes_key" :placeholder="$t('auto.DingtalkAppAccount.ea4845bc')" />
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.5c7f1a99')" v-if="dialogMode === 'edit'">
          <el-input :model-value="callbackUrl(form.id)" readonly>
            <template #append>
              <el-button @click="copyText(callbackUrl(form.id))">{{ $t('auto.DingtalkAppAccount.79d3abe9') }}</el-button>
            </template>
          </el-input>
          <div class="form-hint">{{ $t('auto.DingtalkAppAccount.5e2eda38') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.ab343c20')">
          <el-switch v-model="form.inbound_enabled" />
          <span class="form-hint-inline">{{ $t('auto.DingtalkAppAccount.b51e0d8e') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.eed5907c')">
          <el-input v-model="form.ai_agent_id" :placeholder="$t('auto.DingtalkAppAccount.231dd00b')" />
          <div class="form-hint">{{ $t('auto.DingtalkAppAccount.fc6d37d8') }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.DingtalkAppAccount.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="submit">{{ $t('auto.DingtalkAppAccount.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    
    <agent-binding-dialog
      v-model:visible="bindVisible"
      channel-type="dingtalk"
      :account-id="String(bindRow?.id)"
      :account-name="bindRow?.account_name"
      :current-agent-id="bindRow?.ai_agent_id ? String(bindRow.ai_agent_id) : ''"
      @saved="onBindSaved"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import AgentBindingDialog from '@/components/AgentBindingDialog.vue'
import {
  listDingtalkApp,
  getDingtalkApp,
  createDingtalkApp,
  updateDingtalkApp,
  deleteDingtalkApp,
  testDingtalkApp
} from '@/api/dingtalkApp'

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogMode = ref('create')
const form = reactive({
  id: null,
  account_name: '',
  app_key: '',
  app_secret: '',
  agent_id: '',
  token: '',
  aes_key: '',
  inbound_enabled: false,
  ai_agent_id: ''
})

const bindVisible = ref(false)
const bindRow = ref(null)

function callbackUrl(id) {
  if (!id) return ''
  return `${window.location.origin}/api/webhook/dingtalk/${id}`
}

async function loadData() {
  loading.value = true
  try {
    const res = await listDingtalkApp({ page: page.value, page_size: pageSize.value })
    const arr = Array.isArray(res) ? res : (res.data || [])
    tableData.value = arr
    total.value = res && res.total != null ? res.total : arr.length
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.c4fa9b74') + (e.message || e))
  } finally {
    loading.value = false
  }
}

function resetForm() {
  Object.assign(form, {
    id: null,
    account_name: '',
    app_key: '',
    app_secret: '',
    agent_id: '',
    token: '',
    aes_key: '',
    inbound_enabled: false,
    ai_agent_id: ''
  })
}

function openCreate() {
  dialogMode.value = 'create'
  resetForm()
  dialogVisible.value = true
}

async function openEdit(row) {
  dialogMode.value = 'edit'
  resetForm()
  try {
    const detail = await getDingtalkApp(row.id)
    Object.assign(form, {
      id: detail.id,
      account_name: detail.account_name || '',
      app_key: detail.app_key || '',
      app_secret: '',
      agent_id: detail.agent_id || '',
      token: detail.token || '',
      aes_key: '',
      inbound_enabled: !!detail.inbound_enabled,
      ai_agent_id: detail.ai_agent_id || ''
    })
    dialogVisible.value = true
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.db870bff') + (e.message || e))
  }
}

async function submit() {
  const payload = {
    account_name: form.account_name || undefined,
    app_key: form.app_key || undefined,
    app_secret: form.app_secret || undefined,
    agent_id: form.agent_id || undefined,
    token: form.token || undefined,
    aes_key: form.aes_key || undefined,
    inbound_enabled: form.inbound_enabled,
    ai_agent_id: form.ai_agent_id || undefined
  }
  saving.value = true
  try {
    if (dialogMode.value === 'edit') {
      await updateDingtalkApp(form.id, payload)
    } else {
      await createDingtalkApp(payload)
    }
    ElMessage.success(t('auto.DingtalkAppAccount.3b108349'))
    dialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.78d235cc') + (e.message || e))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除账号「${row.account_name}」？`, t('auto.DingtalkAppAccount.02d9819d'), { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteDingtalkApp(row.id)
    ElMessage.success(t('auto.DingtalkAppAccount.5cc23262'))
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.bc42bb8f') + (e.message || e))
  }
}

async function openTest(row) {
  try {
    const res = await testDingtalkApp(row.id)
    ElMessage.success(t('auto.DingtalkAppAccount.24f2022c') + (res?.ok ? 'AppKey/AppSecret/Token 齐全' : 'OK'))
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.4b3cb9fe') + (e.message || e))
  }
}

function openBind(row) {
  bindRow.value = row
  bindVisible.value = true
}

function onBindSaved() {
  ElMessage.success('AI 绑定已保存')
  loadData()
}

function copyText(text) {
  if (!text) return
  navigator.clipboard.writeText(text).then(
    () => ElMessage.success(t('auto.DingtalkAppAccount.4fb42e6e')),
    () => ElMessage.error(t('auto.DingtalkAppAccount.41da1ca4'))
  )
}

function copyCallbackUrl(row) {
  copyText(callbackUrl(row.id))
}

function handlePage(p) {
  page.value = p
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.pager {
  margin-top: 16px;
  justify-content: flex-end;
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
</style>
