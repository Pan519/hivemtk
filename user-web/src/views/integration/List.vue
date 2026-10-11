<template>
  <div class="integration-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('集成管理') }}</h2>
        <p class="subtitle">{{ $t('auto.List.0574f4c5') }}</p>
      </div>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('添加集成') }}
      </el-button>
    </el-card>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('已启用') }}</div>
            <div class="stat-value" style="color: #10B981">{{ stats.enabled }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('已禁用') }}</div>
            <div class="stat-value" style="color: #909399">{{ stats.disabled }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('异常') }}</div>
            <div class="stat-value" style="color: #EF4444">{{ stats.error }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('总调用次数') }}</div>
            <div class="stat-value">{{ stats.totalCalls }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <el-table :data="integrations" v-loading="loading" stripe>
        <el-table-column prop="account_name" ::label="$t('auto.List.41332a14')" min-width="150" />
        <el-table-column prop="platform" ::label="$t('auto.List.a0ee2507')" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.platform }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="toggleIntegration(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="last_sync_at" :label="$t('auto.List.a4f23d5f')" width="180" />
        <el-table-column prop="callCount" :label="$t('auto.List.00d5bdf1')" width="100" />
        <el-table-column prop="errorCount" :label="$t('auto.List.8731f2a8')" width="100" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="300" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="testConnection(row)">{{ $t('auto.List.db06c78d') }}</el-button>
            <el-button link type="primary" @click="viewLogs(row)">{{ $t('auto.List.456d29ef') }}</el-button>
            <el-button link type="primary" @click="editIntegration(row)">{{ $t('auto.List.224e2ccd') }}</el-button>
            <el-button link type="danger" @click="deleteIntegration(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('auto.List.11fb9b6a')" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.204e637e')" prop="type">
          <el-select v-model="form.type" style="width: 100%">
            <el-option label="Webhook" value="webhook" />
            <el-option label="REST API" value="api" />
            <el-option :label="$t('auto.List.68051bf4')" value="database" />
            <el-option :label="$t('auto.List.8bedb7aa')" value="mq" />
            <el-option label="OAuth" value="oauth" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.f682f85e')" v-if="form.type !== 'database'">
          <el-input v-model="form.endpoint" placeholder="https://..." />
        </el-form-item>
        <el-form-item :label="$t('auto.List.b33c7279')" v-if="form.type !== 'database'">
          <el-select v-model="form.authType" style="width: 100%">
            <el-option label="API Key" value="apikey" />
            <el-option label="Bearer Token" value="bearer" />
            <el-option label="Basic Auth" value="basic" />
            <el-option label="OAuth 2.0" value="oauth2" />
            <el-option :label="$t('auto.List.d81bb206')" value="none" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.5f811dd8')" v-if="form.authType !== 'none' && form.type !== 'database'">
          <el-input v-model="form.credentials" type="textarea" :rows="3" placeholder='{"api_key": "xxx"}' />
        </el-form-item>
        <el-form-item :label="$t('auto.List.e45381d5')" v-if="form.type === 'database'">
          <el-input v-model="form.dbConn" placeholder="user:pass@tcp(host)/dbname" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.1709848a')">
          <el-select v-model="form.syncInterval" style="width: 100%">
            <el-option :label="$t('auto.List.2843e2f6')" value="realtime" />
            <el-option :label="$t('auto.List.35cf1686')" value="1m" />
            <el-option :label="$t('auto.List.860c00f4')" value="1h" />
            <el-option :label="$t('auto.List.78623eae')" value="1d" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.57a166bb')">
          <el-checkbox-group v-model="form.events">
            <el-checkbox label="user.created">{{ $t('auto.List.0dec71dc') }}</el-checkbox>
            <el-checkbox label="clue.updated">{{ $t('auto.List.982624d0') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.List.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="logsDialogVisible" :title="`${currentLogName} - 同步日志`" width="900px">
      <el-table :data="logs" v-loading="logsLoading" stripe border max-height="500px">
        <el-table-column prop="timestamp" :label="$t('auto.List.19fcb9eb')" width="180" />
        <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="100">
          <template #default="{ row: logRow }">
            <el-tag :type="logRow.status === 'success' ? 'success' : logRow.status === 'error' ? 'danger' : 'warning'">
              {{ logRow.status === 'success' ? '成功' : logRow.status === 'error' ? '失败' : '处理中' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" :label="$t('auto.List.ff692f04')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="details" :label="$t('auto.List.f26225bd')" min-width="250" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="80" fixed="right">
          <template #default="{ row: logRow }">
            <el-button link type="primary" @click="showLogDetail(logRow)">{{ $t('auto.List.f26225bd') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.ae733366')" />
        </template>
      </el-table>
      <template #footer>
        <el-button @click="logsDialogVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="logDetailVisible" :title="$t('auto.List.e0d9b6f4')" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="$t('auto.List.19fcb9eb')">{{ currentLogDetail?.timestamp }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3fea7ca7')">
          <el-tag :type="currentLogDetail?.status === 'success' ? 'success' : currentLogDetail?.status === 'error' ? 'danger' : 'warning'">
            {{ currentLogDetail?.status === 'success' ? '成功' : currentLogDetail?.status === 'error' ? '失败' : '处理中' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.ff692f04')">{{ currentLogDetail?.message }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.f26225bd')">{{ currentLogDetail?.details }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="logDetailVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { toList } from '@/utils/list'
import {
  getIntegrations,
  createIntegration,
  updateIntegration,
  deleteIntegration as deleteIntegrationApi,
  toggleIntegrationStatus,
  testIntegration,
  getIntegrationLogs
} from '@/api/integration.js'

const { t } = useI18n()

const loading = ref(false)
const integrations = ref([])
const stats = ref({ enabled: 0, disabled: 0, error: 0, totalCalls: 0 })
const dialogVisible = ref(false)
const dialogTitle = ref(t('auto.List.aedf5c22'))
const formRef = ref()
const logsDialogVisible = ref(false)
const logsLoading = ref(false)
const currentLogName = ref('')
const logs = ref([])
const logDetailVisible = ref(false)
const currentLogDetail = ref(null)
const form = ref({
  id: 0,
  name: '',
  type: 'webhook',
  endpoint: '',
  authType: 'none',
  credentials: '',
  dbConn: '',
  syncInterval: '1h',
  events: []
})
const formRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.List.6987d073')), trigger: 'blur' }],
  type: [{ required: true, message: i18n.global.t(t('auto.List.95f11c56')), trigger: 'change' }]
}

const refreshData = async () => {
  loading.value = true
  try {
    const iRes = await getIntegrations()
    const accounts = iRes?.accounts || toList(iRes)
    integrations.value = accounts
    const enabledCount = accounts.filter(a => a.status === 1).length
    stats.value = {
      enabled: enabledCount,
      disabled: accounts.length - enabledCount,
      error: 0,
      totalCalls: 0
    }
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  form.value = { id: 0, name: '', type: 'webhook', endpoint: '', authType: 'none', credentials: '', dbConn: '', syncInterval: '1h', events: [] }
  dialogTitle.value = t('auto.List.aedf5c22')
  dialogVisible.value = true
}

const editIntegration = (row) => {
  form.value = { ...row }
  dialogTitle.value = t('auto.List.dd51ab81')
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (form.value.id) {
        await updateIntegration(form.value.id, form.value)
      } else {
        await createIntegration(form.value)
      }
      ElMessage.success(i18n.global.t(t('auto.List.3b108349')))
      dialogVisible.value = false
      refreshData()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.6de920b4')))
    }
  })
}

const toggleIntegration = async (row) => {
  await toggleIntegrationStatus(row.id, row.status)
  ElMessage.success(i18n.global.t(t('auto.List.7d6e0e1a')))
  refreshData()
}

const testConnection = async (row) => {
  const loading = ElMessage({ message: i18n.global.t(t('auto.List.96ae82ed')), type: 'info', duration: 0 })
  try {
    const res = await testIntegration(row.id)
    loading.close()
    if (res?.status === 'ok') {
      ElMessage.success(i18n.global.t(t('auto.List.b33199d3')))
    } else {
      ElMessage.error(t('auto.List.fcc15a0f') + (res?.error || ''))
    }
  } catch (error) {
    loading.close()
    ElMessage.error(i18n.global.t(t('auto.List.9710d9b9')))
  }
}

const viewLogs = async (row) => {
  currentLogName.value = row.name
  logsDialogVisible.value = true
  logsLoading.value = true
  logs.value = []
  try {
    const res = await getIntegrationLogs(row.id)
    const data = res;
    if (Array.isArray(data)) {
      logs.value = data
    } else if (data && data.logs) {
      logs.value = data.logs
    }
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.c93f334d')))
  } finally {
    logsLoading.value = false
  }
}

const showLogDetail = (row) => {
  currentLogDetail.value = row
  logDetailVisible.value = true
}

const deleteIntegration = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除集成 "${row.name}"？`, t('auto.List.e83a256e'), { type: 'warning' })
    await deleteIntegrationApi(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
    refreshData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
  }
}

onMounted(() => refreshData())
</script>

<style scoped lang="scss">
.integration-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.stats-row {
  margin-bottom: 20px;
  .stat-item {
    text-align: center;
    .stat-label { color: #909399; font-size: 14px; margin-bottom: 10px; }
    .stat-value { font-size: 28px; font-weight: bold; }
  }
}
.logs-content { padding: 0; }
.log-detail-row { margin-bottom: 12px; }
.log-detail-label { font-weight: bold; color: #606266; margin-bottom: 4px; }
</style>
