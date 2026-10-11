<template>
  <div class="operation-log-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('操作日志') }}</h2>
        <p class="subtitle">{{ $t('审计追踪所有用户操作记录') }}</p>
      </div>
      <div>
        <el-button @click="$router.push('/operationLog/enhanced')">{{ $t('auto.List.1b596eec') }}</el-button>
        <el-button @click="exportLogs">
          <el-icon><Download /></el-icon>
          {{ $t('导出') }}
        </el-button>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          {{ $t('刷新') }}
        </el-button>
      </div>
    </el-card>

    <el-card>
      <div class="filter-bar">
        <el-input v-model="searchKeyword" ::placeholder="$t('auto.List.ebd8b174')" clearable style="width: 200px" />
        <el-select v-model="filterModule" ::placeholder="$t('auto.List.a2a3afcd')" clearable style="width: 150px">
          <el-option ::label="$t('auto.List.9cfb9a37')" value="user" />
          <el-option ::label="$t('auto.List.f35f9813')" value="clue" />
          <el-option ::label="$t('auto.List.4f29ef53')" value="marketing" />
          <el-option ::label="$t('auto.List.b67742c2')" value="system" />
        </el-select>
        <el-select v-model="filterAction" ::placeholder="$t('auto.List.1fc34045')" clearable style="width: 150px">
          <el-option ::label="$t('auto.List.f40c87a1')" value="create" />
          <el-option ::label="$t('auto.List.7ba167cc')" value="update" />
          <el-option ::label="$t('auto.List.ae9b918b')" value="delete" />
          <el-option ::label="$t('auto.List.51ae8247')" value="login" />
          <el-option ::label="$t('auto.List.4c0be02c')" value="logout" />
        </el-select>
        <el-select v-model="filterUser" ::placeholder="$t('auto.List.a110721d')" clearable style="width: 150px">
          <el-option label="admin" value="admin" />
          <el-option label="manager" value="manager" />
        </el-select>
        <el-date-picker v-model="dateRange" type="daterange" :range-separator="$t('auto.List.981cbe31')" :start-placeholder="$t('auto.List.a3e3b883')" :end-placeholder="$t('auto.List.12f1d7ef')" value-format="YYYY-MM-DD" />
      </div>

      <el-table :data="filteredLogs" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" ::label="$t('auto.List.a110721d')" width="100" />
        <el-table-column prop="module" ::label="$t('auto.List.a2a3afcd')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ row.module }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" :label="$t('auto.List.2b6bc0f2')" width="100">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" size="small">{{ getActionText(row.action) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" :label="$t('auto.List.49e81f5e')" min-width="150">
          <template #default="{ row }">
            <span v-if="row.resource_id">{{ row.resource }} #{{ row.resource_id }}</span>
            <span v-else>{{ row.resource }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="detail" :label="$t('auto.List.f3ffb05b')" min-width="250" show-overflow-tooltip />
        <el-table-column prop="ip" :label="$t('auto.List.2dc9105c')" width="130" />
        <el-table-column prop="user_agent" :label="$t('auto.List.912302cb')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('auto.List.19fcb9eb')" width="180" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        layout="total, prev, pager, next, jumper"
        @current-change="loadLogs"
        style="margin-top: 15px; text-align: right"
      />
    </el-card>

    <el-dialog v-model="detailVisible" :title="$t('auto.List.e0d9b6f4')" width="700px">
      <el-descriptions :column="1" border v-if="currentLog">
        <el-descriptions-item :label="$t('auto.List.f9ac4b2a')">{{ currentLog.username }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.fac54c34')">{{ currentLog.module }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.de9cc3dd')">{{ currentLog.action }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.49e81f5e')">
          <span v-if="currentLog.resource_id">{{ currentLog.resource }} #{{ currentLog.resource_id }}</span>
          <span v-else>{{ currentLog.resource }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.f3ffb05b')">{{ currentLog.detail }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.86ae3444')">
          <pre>{{ currentLog.old_value }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.79bf6bc3')">
          <pre>{{ currentLog.new_value }}</pre>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.2dc9105c')">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.912302cb')">{{ currentLog.user_agent }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.19fcb9eb')">{{ currentLog.created_at }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Download, Refresh } from '@element-plus/icons-vue'
import { getOperationLogs, exportOperationLogs, getOperationLogDetail } from '@/api/operationLog.js'

const { t } = useI18n()

const loading = ref(false)
const searchKeyword = ref('')
const filterModule = ref('')
const filterAction = ref('')
const filterUser = ref('')
const dateRange = ref([])
const logs = ref([])
const pagination = ref({ page: 1, size: 20, total: 0 })
const detailVisible = ref(false)
const currentLog = ref(null)

const filteredLogs = computed(() => {
  let result = logs.value
  if (searchKeyword.value) result = result.filter(l => l.description?.includes(searchKeyword.value))
  if (filterModule.value) result = result.filter(l => l.module === filterModule.value)
  if (filterAction.value) result = result.filter(l => l.action === filterAction.value)
  if (filterUser.value) result = result.filter(l => l.username === filterUser.value)
  return result
})

const getActionType = (action) => {
  const map = { create: 'success', update: 'primary', delete: 'danger', login: 'warning', logout: 'info' }
  return map[action]}
const getActionText = (action) => {
  const map = { create: t('auto.List.d9ac9228'), update: t('auto.List.32ac152b'), delete: t('auto.List.2f4aaddd'), login: t('auto.List.402d19e5'), logout: t('auto.List.75c856b1') }
  return map[action] || action
}

const loadLogs = async () => {
  loading.value = true
  try {
    const res = await getOperationLogs({ page: pagination.value.page, size: pagination.value.size })
    logs.value = res?.list || [];
    pagination.value.total = res?.total || 0
  } finally {
    loading.value = false
  }
}

const refreshData = () => loadLogs()

const exportLogs = async () => {
  try {
    const params= {}
    if (filterModule.value) params.module = filterModule.value
    if (filterAction.value) params.action = filterAction.value
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const blob = await exportOperationLogs(params)
    const url = window.URL.createObjectURL(blob)
    const anchor = document.createElement('a')
    anchor.href = url
    anchor.download = `operation-logs-${new Date().toISOString().slice(0, 10)}.csv`
    document.body.appendChild(anchor)
    anchor.click()
    document.body.removeChild(anchor)
    window.URL.revokeObjectURL(url)
    ElMessage.success(i18n.global.t(t('auto.List.105c8a51')))
  } catch (error) {
    ElMessage.error(t('auto.List.2e0d8c60') + (error.message))
  }
}

const viewDetail = async (row) => {
  const res = await getOperationLogDetail(row.id)
  currentLog.value = res
  detailVisible.value = true
}

onMounted(() => loadLogs())
</script>

<style scoped lang="scss">
.operation-log-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
  flex-wrap: wrap;
}
</style>
