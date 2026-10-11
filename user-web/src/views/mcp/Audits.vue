<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <el-tabs v-model="tab" style="flex: 1">
            <el-tab-pane :label="$t('auto.Audits.3d4100b2')" name="call" />
            <el-tab-pane :label="$t('auto.Audits.b5b410cc')" name="admin" />
          </el-tabs>
          <div v-if="tab === 'call'" class="filters">
            <el-input
              v-model="filters.client_id"
              placeholder="Client ID"
              clearable
              size="small"
              style="width: 180px"
              @change="load"
            />
            <el-input
              v-model="filters.tool"
              :placeholder="$t('auto.Audits.7da4055a')"
              clearable
              size="small"
              style="width: 160px"
              @change="load"
            />
            <el-select v-model="filters.event" :placeholder="$t('auto.Audits.10b2761d')" clearable size="small" style="width: 130px" @change="load">
              <el-option :label="$t('auto.Audits.850b4e4d')" value="call" />
              <el-option :label="$t('auto.Audits.3476d2d4')" value="list" />
              <el-option :label="$t('auto.Audits.c0a3c7a5')" value="auth" />
            </el-select>
            <el-select v-model="filters.allowed" :placeholder="$t('auto.Audits.5ad7f5a8')" clearable size="small" style="width: 110px" @change="load">
              <el-option :label="$t('auto.Audits.c6bc4767')" :value="1" />
              <el-option :label="$t('auto.Audits.c87187d0')" :value="0" />
            </el-select>
            <el-date-picker
              v-model="filters.range"
              type="datetimerange"
              size="small"
              style="width: 320px"
              :start-placeholder="$t('auto.Audits.592c5958')"
              :end-placeholder="$t('auto.Audits.f782779e')"
              value-format="YYYY-MM-DDTHH:mm:ssZ"
              @change="load"
            />
            <el-button type="primary" size="small" @click="load">{{ $t('auto.Audits.694fc5ef') }}</el-button>
            <el-button size="small" @click="exportCSV('mcp/audits/export', 'mcp-call-audits')">{{ $t('auto.Audits.1add12ff') }}</el-button>
          </div>
          <div v-else class="filters">
            <el-select v-model="adminFilters.action" :placeholder="$t('auto.Audits.c500cfab')" clearable filterable size="small" style="width: 180px" @change="loadAdmin">
              <el-option v-for="a in actions" :key="a" :label="a" :value="a" />
            </el-select>
            <el-select v-model="adminFilters.target_type" :placeholder="$t('auto.Audits.b1449413')" clearable size="small" style="width: 120px" @change="loadAdmin">
              <el-option :label="$t('auto.Audits.94171ce2')" value="credential" />
              <el-option :label="$t('auto.Audits.cdb81cf6')" value="key" />
              <el-option label="Skill" value="skill" />
            </el-select>
            <el-input v-model="adminFilters.operator" :placeholder="$t('auto.Audits.f9ac4b2a')" clearable size="small" style="width: 120px" @change="loadAdmin" />
            <el-date-picker
              v-model="adminFilters.range"
              type="datetimerange"
              size="small"
              style="width: 320px"
              :start-placeholder="$t('auto.Audits.592c5958')"
              :end-placeholder="$t('auto.Audits.f782779e')"
              value-format="YYYY-MM-DDTHH:mm:ssZ"
              @change="loadAdmin"
            />
            <el-button type="primary" size="small" @click="loadAdmin">{{ $t('auto.Audits.694fc5ef') }}</el-button>
            <el-button size="small" @click="exportCSV('mcp/admin-audits/export', 'mcp-admin-audits')">{{ $t('auto.Audits.1add12ff') }}</el-button>
          </div>
        </div>
      </template>

      <!-- 调用审计：MCP 入口的使用面留痕 -->
      <template v-if="tab === 'call'">
        <el-alert
          type="info"
          :closable="false"
          show-icon
          :title="$t('auto.Audits.15007d3c')"
          style="margin-bottom: 12px"
        />
        <el-table v-loading="loading" :data="rows" border size="small">
          <el-table-column prop="created_at" :label="$t('auto.Audits.19fcb9eb')" width="170">
            <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
          </el-table-column>
          <el-table-column prop="client_id" label="Client ID" width="150" />
          <el-table-column :label="$t('auto.Audits.10b2761d')" width="90">
            <template #default="{ row }">
              <el-tag v-if="row.event === 'auth'" type="danger" size="small">{{ $t('auto.Audits.c0a3c7a5') }}</el-tag>
              <el-tag v-else-if="row.event === 'list'" type="info" size="small">{{ $t('auto.Audits.56975625') }}</el-tag>
              <el-tag v-else size="small">{{ $t('auto.Audits.97d29d84') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="method" :label="$t('auto.Audits.ea340b9d')" width="100" />
          <el-table-column prop="tool" :label="$t('auto.Audits.20dce2c6')" min-width="140" />
          <el-table-column :label="$t('auto.Audits.02b66eb5')" width="70">
            <template #default="{ row }">
              <span v-if="row.key_id" class="mono">#{{ row.key_id }}</span>
              <span v-else class="muted">—</span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Audits.5ad7f5a8')" width="110">
            <template #default="{ row }">
              <el-tag v-if="row.event === 'auth'" type="danger" size="small">{{ $t('auto.Audits.7173f809') }}</el-tag>
              <el-tag v-else-if="row.allowed && row.ok" type="success" size="small">{{ $t('auto.Audits.330363df') }}</el-tag>
              <el-tag v-else-if="row.allowed" type="warning" size="small">{{ $t('auto.Audits.acd5cb84') }}</el-tag>
              <el-tag v-else type="danger" size="small">{{ $t('auto.Audits.a8a1bf23') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Audits.39f1374d')" width="80">
            <template #default="{ row }">{{ row.event === 'auth' ? '—' : row.latency_ms + 'ms' }}</template>
          </el-table-column>
          <el-table-column prop="ip" :label="$t('auto.Audits.11c50da4')" width="130" />
          <el-table-column prop="args_hash" :label="$t('auto.Audits.436c61c7')" width="150">
            <template #default="{ row }">
              <code class="mono">{{ row.args_hash || '—' }}</code>
            </template>
          </el-table-column>
          <el-table-column prop="error_msg" :label="$t('auto.Audits.7030ff64')" min-width="160" show-overflow-tooltip />
        </el-table>
        <el-pagination
          v-model:current-page="page"
          :page-size="size"
          :total="total"
          layout="total, prev, pager, next"
          small
          style="margin-top: 12px; justify-content: flex-end"
          @current-change="load"
        />
      </template>

      <!-- 管理审计：MCP/Skill 管理面的操作账本 -->
      <template v-else>
        <el-alert
          type="info"
          :closable="false"
          show-icon
          :title="$t('auto.Audits.78ff2ae1')"
          style="margin-bottom: 12px"
        />
        <el-table v-loading="adminLoading" :data="adminRows" border size="small">
          <el-table-column prop="created_at" :label="$t('auto.Audits.19fcb9eb')" width="170">
            <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
          </el-table-column>
          <el-table-column prop="action" :label="$t('auto.Audits.c500cfab')" width="170">
            <template #default="{ row }">
              <el-tag size="small" :type="actionTagType(row.action)">{{ row.action }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Audits.b1449413')" width="90">
            <template #default="{ row }">
              {{ { credential: '凭证', key: '密钥', skill: 'Skill' }[row.target_type] || row.target_type }} #{{ row.target_id }}
            </template>
          </el-table-column>
          <el-table-column prop="target_label" :label="$t('auto.Audits.14d34236')" min-width="140" show-overflow-tooltip>
            <template #default="{ row }">{{ row.target_label || '—' }}</template>
          </el-table-column>
          <el-table-column prop="operator_username" :label="$t('auto.Audits.f9ac4b2a')" width="110">
            <template #default="{ row }">{{ row.operator_username || '—' }}</template>
          </el-table-column>
          <el-table-column prop="ip" :label="$t('auto.Audits.11c50da4')" width="130" />
          <el-table-column prop="detail" :label="$t('auto.Audits.f26225bd')" min-width="240" show-overflow-tooltip>
            <template #default="{ row }">
              <code class="mono">{{ row.detail || '—' }}</code>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="adminPage"
          :page-size="size"
          :total="adminTotal"
          layout="total, prev, pager, next"
          small
          style="margin-top: 12px; justify-content: flex-end"
          @current-change="loadAdmin"
        />
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { mcpAuditApi, mcpAdminAuditApi } from '@/api/mcpCredential'

const { t } = useI18n()

const tab = ref('call')
const size = 20

// —— 调用审计 ——
const rows = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const filters = ref({ client_id: '', tool: '', allowed: undefined, event: undefined, range: null })

// —— 管理审计 ——
const adminRows = ref([])
const adminLoading = ref(false)
const adminTotal = ref(0)
const adminPage = ref(1)
const adminFilters = ref({ action: '', target_type: '', operator: '', range: null })
const actions = [
  'credential.create', 'credential.update', 'credential.delete',
  'key.create', 'key.enable', 'key.disable', 'key.revoke', 'key.delete', 'key.reveal',
  'skill.create', 'skill.update', 'skill.delete'
]

const fmtTime = (v) => (v ? String(v).replace('T', ' ').slice(0, 19) : '—')
const actionTagType = (a) => {
  if (a?.endsWith('.delete')) return 'danger'
  if (a === 'key.revoke' || a === 'key.disable') return 'warning'
  if (a === 'key.reveal') return 'danger'
  if (a?.endsWith('.create')) return 'success'
  return 'info'
}

const load = async () => {
  loading.value = true
  try {
    const params = { page: page.value, size }
    if (filters.value.client_id) params.client_id = filters.value.client_id
    if (filters.value.tool) params.tool = filters.value.tool
    if (filters.value.event) params.event = filters.value.event
    if (filters.value.range?.[0]) params.start = filters.value.range[0]
    if (filters.value.range?.[1]) params.end = filters.value.range[1]
    if (filters.value.allowed !== undefined && filters.value.allowed !== null && filters.value.allowed !== '') {
      params.allowed = filters.value.allowed
    }
    const res = await mcpAuditApi.list({ params })
    const data = res?.data ?? res ?? {}
    rows.value = data.list ?? []
    total.value = data.total ?? 0
  } finally {
    loading.value = false
  }
}

const loadAdmin = async () => {
  adminLoading.value = true
  try {
    const params = { page: adminPage.value, size }
    if (adminFilters.value.action) params.action = adminFilters.value.action
    if (adminFilters.value.target_type) params.target_type = adminFilters.value.target_type
    if (adminFilters.value.operator) params.operator = adminFilters.value.operator
    if (adminFilters.value.range?.[0]) params.start = adminFilters.value.range[0]
    if (adminFilters.value.range?.[1]) params.end = adminFilters.value.range[1]
    const res = await mcpAdminAuditApi.list({ params })
    const data = res?.data ?? res ?? {}
    adminRows.value = data.list ?? []
    adminTotal.value = data.total ?? 0
  } finally {
    adminLoading.value = false
  }
}

watch(tab, (v) => {
  if (v === 'admin' && adminRows.value.length === 0) loadAdmin()
})

// 导出：带 JWT 的 fetch 下载（导出端点在鉴权组内，window.open 不带令牌）
const exportCSV = async (path, name) => {
  try {
    const token = localStorage.getItem('token') || ''
    const qs = new URLSearchParams()
    const src = tab.value === 'call' ? filters.value : adminFilters.value
    Object.entries(src).forEach(([k, v]) => {
      if (k === 'range' || v === undefined || v === null || v === '') return
      qs.set(k, v)
    })
    const resp = await fetch(`/api/${path}?${qs.toString()}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`)
    const blob = await resp.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${name}-${new Date().toISOString().slice(0, 19).replace(/[:T]/g, '')}.csv`
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    ElMessage.error(t('auto.Audits.bccb121f') + e.message)
  }
}

onMounted(load)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}
.filters {
  display: flex;
  gap: 8px;
  align-items: center;
}
.mono {
  font-family: monospace;
  font-size: 12px;
}
.muted {
  color: #999;
  font-size: 12px;
}
</style>
