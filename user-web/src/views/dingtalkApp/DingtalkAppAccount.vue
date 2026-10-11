<template>
  <div class="dingtalk-app">
    <el-card class="card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.DingtalkAppAccount.17.1ef4b3') }}</span>
          <el-button type="primary" @click="openCreate">{{ $t('auto.DingtalkAppAccount.18.34374b') }}</el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="account_name" :label="$t('auto.DingtalkAppAccount.1.d7ec2d')" min-width="120" />
        <el-table-column prop="app_key" label="AppKey" min-width="130" show-overflow-tooltip />
        <el-table-column prop="agent_id" label="AgentId" min-width="100" show-overflow-tooltip />
        <el-table-column :label="$t('auto.DingtalkAppAccount.2.ab343c')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.inbound_enabled ? 'success' : 'info'" size="small">
              {{ row.inbound_enabled ? '已启用' : '未启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.DingtalkAppAccount.3.c9e87d')" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag :type="row.ai_agent_id ? 'success' : 'info'" size="small">
              {{ row.ai_agent_id || '未绑定' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.DingtalkAppAccount.4.2b6bc0')" width="320" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openTest(row)">{{ $t('auto.DingtalkAppAccount.19.db06c7') }}</el-button>
            <el-button link type="primary" @click="openBind(row)">{{ $t('auto.DingtalkAppAccount.3.c9e87d') }}</el-button>
            <el-button link type="primary" @click="copyCallbackUrl(row)">{{ $t('auto.DingtalkAppAccount.20.e73c73') }}</el-button>
            <el-button link type="primary" @click="openEdit(row)">{{ $t('auto.DingtalkAppAccount.21.95b351') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ $t('auto.DingtalkAppAccount.22.2f4aad') }}</el-button>
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

    
    <el-dialog v-model="dialogVisible" ::title="$t('auto.DingtalkAppAccount.5.ab3ad6')" width="680px">
      <el-form :model="form" label-width="150px">
        <el-form-item :label="$t('auto.DingtalkAppAccount.1.d7ec2d')" prop="account_name">
          <el-input v-model="form.account_name" :placeholder="$t('auto.DingtalkAppAccount.6.17327f')" />
        </el-form-item>
        <el-form-item label="AppKey" prop="app_key">
          <el-input v-model="form.app_key" :placeholder="$t('auto.DingtalkAppAccount.7.57ba6e')" />
        </el-form-item>
        <el-form-item label="AppSecret" prop="app_secret">
          <el-input v-model="form.app_secret" type="password" show-password :placeholder="$t('auto.DingtalkAppAccount.8.a9798e')" />
        </el-form-item>
        <el-form-item label="AgentId" prop="agent_id">
          <el-input v-model="form.agent_id" :placeholder="$t('auto.DingtalkAppAccount.9.93abc0')" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('auto.DingtalkAppAccount.23.f96de9') }}</el-divider>
        <el-form-item :label="$t('auto.DingtalkAppAccount.10.ecbc51')" prop="token">
          <el-input v-model="form.token" :placeholder="$t('auto.DingtalkAppAccount.11.8b7ea3')" />
          <div class="form-hint">{{ $t('auto.DingtalkAppAccount.24.61e4a5') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.12.67f586')" prop="aes_key">
          <el-input v-model="form.aes_key" :placeholder="$t('auto.DingtalkAppAccount.13.ea4845')" />
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.14.5c7f1a')" v-if="dialogMode === 'edit'">
          <el-input :model-value="callbackUrl(form.id)" readonly>
            <template #append>
              <el-button @click="copyText(callbackUrl(form.id))">{{ $t('auto.DingtalkAppAccount.25.79d3ab') }}</el-button>
            </template>
          </el-input>
          <div class="form-hint">{{ $t('auto.DingtalkAppAccount.26.5e2eda') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.2.ab343c')">
          <el-switch v-model="form.inbound_enabled" />
          <span class="form-hint-inline">{{ $t('auto.DingtalkAppAccount.27.b51e0d') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.DingtalkAppAccount.15.eed590')">
          <el-input v-model="form.ai_agent_id" :placeholder="$t('auto.DingtalkAppAccount.16.231dd0')" />
          <div class="form-hint">{{ $t('auto.DingtalkAppAccount.28.fc6d37') }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.DingtalkAppAccount.29.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="submit">{{ $t('auto.DingtalkAppAccount.30.be5fbb') }}</el-button>
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
    ElMessage.error(t('auto.DingtalkAppAccount.31.c4fa9b') + (e.message || e))
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
    ElMessage.error(t('auto.DingtalkAppAccount.32.db870b') + (e.message || e))
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
    ElMessage.success(t('auto.DingtalkAppAccount.33.3b1083'))
    dialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.34.78d235') + (e.message || e))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除账号「${row.account_name}」？`, t('auto.DingtalkAppAccount.35.02d981'), { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteDingtalkApp(row.id)
    ElMessage.success(t('auto.DingtalkAppAccount.36.5cc232'))
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.37.bc42bb') + (e.message || e))
  }
}

async function openTest(row) {
  try {
    const res = await testDingtalkApp(row.id)
    ElMessage.success(t('auto.DingtalkAppAccount.38.24f202') + (res?.ok ? 'AppKey/AppSecret/Token 齐全' : 'OK'))
  } catch (e) {
    ElMessage.error(t('auto.DingtalkAppAccount.39.4b3cb9') + (e.message || e))
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
    () => ElMessage.success(t('auto.DingtalkAppAccount.40.4fb42e')),
    () => ElMessage.error(t('auto.DingtalkAppAccount.41.41da1c'))
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
