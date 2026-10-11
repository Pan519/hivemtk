<template>
  <div class="openapi-integration">
    
    <el-alert
      ::title="$t('auto.OpenAPIIntegration.ab958c0f')"
      type="info"
      :closable="false"
      show-icon
    >
      <template #default>
        通过配置外部 OpenAPI 接口,自动将返回的数据导入到知识库。支持 GET / POST / RESTful 风格,支持 Bearer / API Key / HMAC / Basic 认证。
      </template>
    </el-alert>

    
    <el-card class="filter-card">
      <div class="filter-bar">
        <el-select v-model="filter.product_id" :placeholder="$t('auto.OpenAPIIntegration.86705f46')" clearable style="width: 220px" @change="loadSources">
          <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
        </el-select>
        <el-button :icon="Refresh" @click="loadSources">{{ $t('auto.OpenAPIIntegration.694fc5ef') }}</el-button>
        <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">{{ $t('auto.OpenAPIIntegration.2b4af7e4') }}</el-button>
      </div>
    </el-card>

    
    <el-card>
      <el-table :data="sources" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('auto.OpenAPIIntegration.d7ec2d3f')" min-width="160" show-overflow-tooltip />
        <el-table-column :label="$t('auto.OpenAPIIntegration.226b0912')" width="80">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type?.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.OpenAPIIntegration.ea340b9d')" width="80">
          <template #default="{ row }">
            <el-tag :type="getMethodTagType(row.method)" size="small">{{ getMethodLabel(row.method) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="endpoint" :label="$t('auto.OpenAPIIntegration.cf28e8ea')" min-width="280" show-overflow-tooltip>
          <template #default="{ row }">
            <el-text size="small" style="font-family: monospace">{{ row.endpoint }}</el-text>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.OpenAPIIntegration.b7158a42')" width="100">
          <template #default="{ row }">{{ authTypeLabel(row.auth_type) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.OpenAPIIntegration.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-switch v-model="row._enabled" @change="(val) => handleToggle(row, val)" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.OpenAPIIntegration.63ce1b73')" width="180">
          <template #default="{ row }">
            <div>{{ row.last_sync_at ? formatDate(row.last_sync_at) : '从未同步' }}</div>
            <el-tag :type="getLastStatusTagType(row.last_status)" size="small">
              {{ getLastStatusLabel(row.last_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_synced" :label="$t('auto.OpenAPIIntegration.a5373155')" width="100" />
        <el-table-column :label="$t('auto.OpenAPIIntegration.2b6bc0f2')" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="showEditDialog(row)">{{ $t('auto.OpenAPIIntegration.95b351c8') }}</el-button>
            <el-button link type="success" size="small" @click="handleTest(row)">{{ $t('auto.OpenAPIIntegration.db06c78d') }}</el-button>
            <el-button link type="primary" size="small" @click="handleSync(row)">{{ $t('auto.OpenAPIIntegration.5f71b2b2') }}</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">{{ $t('auto.OpenAPIIntegration.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && sources.length === 0" :description="$t('auto.OpenAPIIntegration.1bf1eae0')" />
    </el-card>

    
    <el-dialog v-model="showCreateDialog" :title="editingId ? '编辑数据源' : '新建数据源'" width="720px" :close-on-click-modal="false">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('auto.OpenAPIIntegration.d7ec2d3f')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.OpenAPIIntegration.819e4d1d')" />
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.aa3fdfcb')" prop="product_id">
          <el-select v-model="form.product_id" :placeholder="$t('auto.OpenAPIIntegration.cbf1cfff')" style="width: 100%">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.bcfbd51b')" prop="endpoint">
          <el-input v-model="form.endpoint" placeholder="https://api.example.com/v1/docs" />
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.ec49329b')">
          <el-radio-group v-model="form.method">
            <el-radio-button label="GET">GET</el-radio-button>
            <el-radio-button label="POST">POST</el-radio-button>
            <el-radio-button label="PUT">PUT</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.b33c7279')">
          <el-select v-model="form.auth_type" style="width: 100%">
            <el-option :label="$t('auto.OpenAPIIntegration.7870c3f1')" value="none" />
            <el-option label="Bearer Token" value="bearer" />
            <el-option label="API Key" value="api_key" />
            <el-option :label="$t('auto.OpenAPIIntegration.baffb04e')" value="hmac" />
            <el-option label="Basic Auth" value="basic" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.auth_type === 'bearer'" label="Token">
          <el-input v-model="form.auth_token" placeholder="Bearer Token" show-password />
        </el-form-item>
        <el-form-item v-if="form.auth_type === 'api_key'" label="API Key">
          <el-input v-model="form.auth_token" :placeholder="$t('auto.OpenAPIIntegration.ed3b95c9')" show-password />
        </el-form-item>
        <el-form-item v-if="form.auth_type === 'hmac'" label="HMAC Secret">
          <el-input v-model="form.auth_token" :placeholder="$t('auto.OpenAPIIntegration.eca22544')" show-password />
        </el-form-item>
        <template v-if="form.auth_type === 'basic'">
          <el-form-item :label="$t('auto.OpenAPIIntegration.819767ad')">
            <el-input v-model="form.auth_username" :placeholder="$t('auto.OpenAPIIntegration.b82c1bf0')" />
          </el-form-item>
          <el-form-item :label="$t('auto.OpenAPIIntegration.a8105204')">
            <el-input v-model="form.auth_token" type="password" :placeholder="$t('auto.OpenAPIIntegration.31e80542')" show-password />
          </el-form-item>
        </template>
        <el-form-item :label="$t('auto.OpenAPIIntegration.f0bfdcf2')" v-if="form.method === 'POST' || form.method === 'PUT'">
          <el-input v-model="form.request_template" type="textarea" :rows="4" placeholder="支持模板变量: {{ now }}, {{ timestamp }}, {{ date }}, {{ last_sync_at }}" />
          <div class="form-tip">{{ $t('auto.OpenAPIIntegration.29a277e3') }}<code v-pre>{{.last_sync_at}}</code>{{ $t('auto.OpenAPIIntegration.b54a5788') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.4fcfe747')">
          <el-input v-model="form.response_path" placeholder="data.items" />
          <div class="form-tip">{{ $t('auto.OpenAPIIntegration.902a9878') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.e5b2ad1c')">
          <el-input v-model="form.field_mapping" type="textarea" :rows="3" placeholder='{"title": "name", "content": "body", "ref": "id"}' />
          <div class="form-tip">{{ $t('auto.OpenAPIIntegration.334c7f06') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.a1bd9760')">
          <el-input v-model="form.schedule" :placeholder="$t('auto.OpenAPIIntegration.e840c2c6')" />
        </el-form-item>
        <el-form-item :label="$t('auto.OpenAPIIntegration.7854b52a')">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">{{ $t('auto.OpenAPIIntegration.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">{{ $t('auto.OpenAPIIntegration.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="showTestResultDialog" :title="$t('auto.OpenAPIIntegration.92440ed5')" width="600px">
      <el-descriptions :column="1" border v-if="testResult">
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.5e9f2ded')">
          <el-tag :type="testResult.success ? 'success' : 'danger'">
            {{ testResult.success ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="testResult.status_code !== undefined" :label="$t('auto.OpenAPIIntegration.ac7f0f8b')">
          <el-tag :type="testResult.status_code < 400 ? 'success' : 'danger'">{{ testResult.status_code }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="testResult.latency_ms !== undefined" :label="$t('auto.OpenAPIIntegration.e6d0202a')">
          {{ testResult.latency_ms }} ms
        </el-descriptions-item>
        <el-descriptions-item v-if="testResult.body_size !== undefined" :label="$t('auto.OpenAPIIntegration.9dc5e1ab')">
          {{ formatFileSize(testResult.body_size) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="testResult.error" :label="$t('auto.OpenAPIIntegration.4604d502')">
          <el-text type="danger">{{ testResult.error }}</el-text>
        </el-descriptions-item>
        <el-descriptions-item v-if="testResult.body_sample" :label="$t('auto.OpenAPIIntegration.bbd61587')">
          <pre class="response-preview">{{ testResult.body_sample }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    
    <el-dialog v-model="showSyncResultDialog" :title="$t('auto.OpenAPIIntegration.c6db98ed')" width="600px">
      <el-descriptions :column="2" border v-if="syncResult">
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.5455357a')">{{ syncResult.source_id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.3fea7ca7')">
          <el-tag :type="syncResult.status === 'success' ? 'success' : 'danger'">{{ syncResult.status === 'success' ? '成功' : '失败' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.f529de30')">{{ syncResult.total_items }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.ef1bcebd')">{{ syncResult.imported_num }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.4a085478')">{{ syncResult.skipped_num }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.acd5cb84')">{{ syncResult.failed_num }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.OpenAPIIntegration.39f1374d')" :span="2">{{ syncResult.duration_ms }} ms</el-descriptions-item>
        <el-descriptions-item v-if="syncResult.error_msg" :label="$t('auto.OpenAPIIntegration.4604d502')" :span="2">
          <el-text type="danger">{{ syncResult.error_msg }}</el-text>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { knowledgeAPI } from '@/api/knowledge'
import { ragProductConfigAPI } from '@/api/ragProductConfig'
import { getAuthTypeLabel } from '@/constants/authType';
import { PASS_FAIL_STATUS, getStatusLabel, getStatusTagType } from '@/constants/status'

const HTTP_METHOD_TAG = { GET: '', POST: 'success', PUT: 'warning', DELETE: 'danger', PATCH: 'warning' };
const getMethodLabel = (m) => (m || '').toUpperCase()
const getMethodTagType = (m) => HTTP_METHOD_TAG[(m || '').toUpperCase()] || 'info'
const getLastStatusLabel = (s) => (s ? getStatusLabel(s, PASS_FAIL_STATUS) : i18n.global.t('auto.OpenAPIIntegration.a0bf20e8'));
const getLastStatusTagType = (s) => (s ? getStatusTagType(s, PASS_FAIL_STATUS) : 'info')

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const showCreateDialog = ref(false)
const showTestResultDialog = ref(false)
const showSyncResultDialog = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const productList = ref([])
const sources = ref([])
const testResult = ref(null)
const syncResult = ref(null)

const filter = reactive({ product_id: '' })

const form = reactive({
  name: '',
  product_id: '',
  endpoint: '',
  method: 'GET',
  auth_type: 'none',
  auth_token: '',
  auth_username: '',
  request_template: '',
  response_path: '',
  field_mapping: '',
  schedule: '',
  enabled: true
})

const rules = {
  name: [{ required: true, message: i18n.global.t(t('auto.OpenAPIIntegration.06e2f88f')), trigger: 'blur' }],
  product_id: [{ required: true, message: i18n.global.t(t('auto.OpenAPIIntegration.3b8e12b1')), trigger: 'change' }],
  endpoint: [{ required: true, message: i18n.global.t(t('auto.OpenAPIIntegration.894fdc2e')), trigger: 'blur' }]
}

const loadAll = async () => {
  await Promise.all([loadProducts(), loadSources()])
}

const loadProducts = async () => {
  try {
    const res = await ragProductConfigAPI.listProducts()
    if (Array.isArray(res)) {
      productList.value = res
    } else if (Array.isArray(res?.items)) {
      productList.value = res.items
    }
  } catch (e) {
    console.error(t('auto.OpenAPIIntegration.4c46f3b1'), e)
  }
}

const loadSources = async () => {
  loading.value = true
  try {
    const res = await knowledgeAPI.listOpenAPISources({ product_id: filter.product_id })
    const list = res?.items || []
    list.forEach(s => { s._enabled = s.enabled === 1 })
    sources.value = list
  } catch (e) {
    ElMessage.error(t('auto.OpenAPIIntegration.a63c8153') + (e.message || ''))
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  Object.assign(form, {
    name: '',
    product_id: '',
    endpoint: '',
    method: 'GET',
    auth_type: 'none',
    auth_token: '',
    auth_username: '',
    request_template: '',
    response_path: '',
    field_mapping: '',
    schedule: '',
    enabled: true
  })
  editingId.value = null
  formRef.value?.clearValidate()
}

const showEditDialog = (row) => {
  Object.assign(form, {
    name: row.name,
    product_id: row.product_id,
    endpoint: row.endpoint,
    method: row.method,
    auth_type: row.auth_type,
    auth_token: '',
    auth_username: '',
    request_template: row.request_template,
    response_path: row.response_path,
    field_mapping: row.field_mapping,
    schedule: row.schedule,
    enabled: row.enabled === 1
  })
  try {
    const authCfg = row.auth_config ? JSON.parse(row.auth_config) : {}
    if (authCfg.token) form.auth_token = authCfg.token
    if (authCfg.secret) form.auth_token = authCfg.secret
    if (authCfg.username) form.auth_username = authCfg.username
  } catch (e) { console.warn("[request] 后台调用失败(已忽略):", e) }
  editingId.value = row.id
  showCreateDialog.value = true
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return;
  }
  saving.value = true
  try {
    const authConfig = {};
    if (form.auth_type === 'bearer' || form.auth_type === 'api_key' || form.auth_type === 'hmac') {
      authConfig.token = form.auth_token
      authConfig.secret = form.auth_token
    } else if (form.auth_type === 'basic') {
      authConfig.username = form.auth_username
      authConfig.password = form.auth_token
    }
    const payload = {
      name: form.name,
      product_id: form.product_id,
      endpoint: form.endpoint,
      method: form.method,
      type: 'rest',
      auth_type: form.auth_type,
      auth_config: JSON.stringify(authConfig),
      request_template: form.request_template,
      response_path: form.response_path,
      field_mapping: form.field_mapping,
      schedule: form.schedule,
      enabled: form.enabled ? 1 : 0
    }
    if (editingId.value) {
      await knowledgeAPI.updateOpenAPISource(editingId.value, payload)
      ElMessage.success(i18n.global.t(t('auto.OpenAPIIntegration.55aa6366')))
    } else {
      await knowledgeAPI.createOpenAPISource(payload)
      ElMessage.success(i18n.global.t(t('auto.OpenAPIIntegration.04a691b3')))
    }
    showCreateDialog.value = false
    resetForm()
    loadSources()
  } catch (e) {
    if (e !== false) {
      ElMessage.error(t('auto.OpenAPIIntegration.40f90217') + (e.message || ''))
    }
  } finally {
    saving.value = false
  }
}

const handleTest = async (row) => {
  try {
    const payload = {
      name: row.name,
      endpoint: row.endpoint,
      method: row.method,
      auth_type: row.auth_type,
      auth_config: row.auth_config,
      request_template: row.request_template,
      response_path: row.response_path
    }
    const res = await knowledgeAPI.testOpenAPISource(row.id, payload)
    testResult.value = res
    showTestResultDialog.value = true
  } catch (e) {
    ElMessage.error(t('auto.OpenAPIIntegration.9d9b6ba9') + (e.message || ''))
  }
}

const handleSync = async (row) => {
  try {
    await ElMessageBox.confirm(`确认立即同步数据源「${row.name}」吗?`, t('auto.OpenAPIIntegration.5f71b2b2'), { type: 'info' })
    const res = await knowledgeAPI.syncOpenAPISource(row.id, { product_id: row.product_id })
    syncResult.value = res
    showSyncResultDialog.value = true
    loadSources()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(t('auto.OpenAPIIntegration.46df3629') + (e.message || ''))
  }
}

const handleToggle = async (row, enabled) => {
  try {
    await knowledgeAPI.toggleOpenAPISource(row.id, { enabled }, { product_id: row.product_id })
    ElMessage.success(enabled ? t('auto.OpenAPIIntegration.53ace430') : t('auto.OpenAPIIntegration.1c1ed981'))
  } catch (e) {
    row._enabled = !enabled
    ElMessage.error(t('auto.OpenAPIIntegration.979ab6ba') + (e.message || ''))
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除数据源「${row.name}」吗?已同步的数据不会被删除`, t('auto.OpenAPIIntegration.cd72cca9'), { type: 'warning' })
    await knowledgeAPI.deleteOpenAPISource(row.id, { product_id: row.product_id })
    ElMessage.success(i18n.global.t(t('auto.OpenAPIIntegration.0007d170')))
    loadSources()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(t('auto.OpenAPIIntegration.ad23f072') + (e.message || ''))
  }
}

const authTypeLabel = (t) => getAuthTypeLabel(t)
const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'
const formatFileSize = (b) => {
  if (!b || b <= 0) return '-'
  if (b < 1024) return b + ' B'
  if (b < 1024 * 1024) return (b / 1024).toFixed(1) + ' KB'
  return (b / 1024 / 1024).toFixed(1) + ' MB'
}

onMounted(() => {
  loadAll()
})
</script>

<style scoped lang="scss">
.openapi-integration {
  .el-alert {
    margin-bottom: 16px;
  }
  .filter-card {
    margin-bottom: 16px;
  }
  .filter-bar {
    display: flex;
    gap: 12px;
    align-items: center;
  }
  .form-tip {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }
  .response-preview {
    background: #f5f7fa;
    padding: 8px;
    border-radius: 4px;
    font-size: 12px;
    max-height: 200px;
    overflow: auto;
    white-space: pre-wrap;
    word-break: break-all;
  }
}
</style>
