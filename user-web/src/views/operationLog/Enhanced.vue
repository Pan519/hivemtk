<template>
  <div class="op-log">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Enhanced.17.d978e9') }}</span>
          <div>
            <el-button :icon="Download" @click="exportLogs">{{ $t('auto.Enhanced.18.55405e') }}</el-button>
            <el-button :icon="Refresh" @click="load">{{ $t('auto.Enhanced.19.694fc5') }}</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" :model="filters" class="filters">
        <el-form-item :label="$t('auto.Enhanced.1.f9ac4b')">
          <el-input v-model="filters.user" :placeholder="$t('auto.Enhanced.2.819767')" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.3.de9cc3')">
          <el-select v-model="filters.action" clearable style="width: 160px">
            <el-option :label="$t('auto.Enhanced.4.d9ac92')" value="create" />
            <el-option :label="$t('auto.Enhanced.5.32ac15')" value="update" />
            <el-option :label="$t('auto.Enhanced.6.2f4aad')" value="delete" />
            <el-option :label="$t('auto.Enhanced.7.402d19')" value="login" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.8.19fcb9')">
          <el-date-picker v-model="filters.dateRange" type="daterange" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="load">{{ $t('auto.Enhanced.20.bee912') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="logs" v-loading="loading">
        <el-table-column prop="createdAt" :label="$t('auto.Enhanced.8.19fcb9')" width="170" />
        <el-table-column prop="user" :label="$t('auto.Enhanced.1.f9ac4b')" width="100" />
        <el-table-column prop="action" :label="$t('auto.Enhanced.9.c500cf')" width="100">
          <template #default="{ row }">
            <el-tag :type="actionType(row.action)" size="small">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" :label="$t('auto.Enhanced.10.eee83a')" width="160" />
        <el-table-column prop="summary" :label="$t('auto.Enhanced.11.3ae146')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="$t('auto.Enhanced.12.8d5351')" width="120">
          <template #default="{ row }">
            <el-button v-if="row.diff" link type="primary" @click="showDiff(row)">{{ $t('auto.Enhanced.21.0385ba') }}</el-button>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Enhanced.13.2b6bc0')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="warning" :disabled="!row.rollbackable" @click="rollback(row)">{{ $t('auto.Enhanced.22.d00b48') }}</el-button>
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.Enhanced.23.f26225') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next, jumper"
        @current-change="load"
      />
    </el-card>

    
    <el-dialog v-model="diffVisible" :title="$t('auto.Enhanced.12.8d5351')" width="800px">
      <el-table :data="diffData">
        <el-table-column prop="field" :label="$t('auto.Enhanced.14.9caecd')" width="160" />
        <el-table-column :label="$t('auto.Enhanced.15.405ac6')">
          <template #default="{ row }"><code>{{ row.before }}</code></template>
        </el-table-column>
        <el-table-column :label="$t('auto.Enhanced.16.8896d6')">
          <template #default="{ row }"><code style="color: #10B981">{{ row.after }}</code></template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Refresh } from '@element-plus/icons-vue'
import { http } from '@/utils/request'

const { t } = useI18n()

const filters = reactive({ user: '', action: '', dateRange: [] })
const logs = ref([])
const loading = ref(false)
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })
const diffVisible = ref(false)
const diffData = ref([])

async function load() {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      user: filters.user,
      action: filters.action,
      start: filters.dateRange?.[0],
      end: filters.dateRange?.[1]
    }
    const res = await http.get('/api/operation-logs', params)
    logs.value = res.items || []
    pagination.total = res.total || 0
  } finally {
    loading.value = false
  }
}

function actionType(action) {
  return { create: 'success', update: 'primary', delete: 'danger', login: 'info' }[action] || ''
}

function showDiff(row) {
  diffData.value = row.diff || []
  diffVisible.value = true
}

async function rollback(row) {
  try {
    await ElMessageBox.confirm(`确认回滚操作 ${row.id}？`, t('auto.Enhanced.24.d3dc08'), { type: 'warning' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.post(`/api/operation-logs/${row.id}/rollback`, {})
  ElMessage.success(t('auto.Enhanced.25.ef6db9'))
  await load()
}

function viewDetail(row) {
  ElMessageBox.alert(JSON.stringify(row, null, 2), t('auto.Enhanced.23.f26225'))
}

function exportLogs() {
  window.open(`/api/operation-logs/export?${new URLSearchParams(filters)}`, '_blank')
}

onMounted(load)
</script>

<style scoped>
.op-log { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.filters { margin-bottom: 12px; }
</style>
