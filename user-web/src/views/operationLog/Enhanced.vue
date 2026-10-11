<template>
  <div class="op-log">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Enhanced.d978e9a2') }}</span>
          <div>
            <el-button :icon="Download" @click="exportLogs">{{ $t('auto.Enhanced.55405ea6') }}</el-button>
            <el-button :icon="Refresh" @click="load">{{ $t('auto.Enhanced.694fc5ef') }}</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" :model="filters" class="filters">
        <el-form-item :label="$t('auto.Enhanced.f9ac4b2a')">
          <el-input v-model="filters.user" :placeholder="$t('auto.Enhanced.819767ad')" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.de9cc3dd')">
          <el-select v-model="filters.action" clearable style="width: 160px">
            <el-option :label="$t('auto.Enhanced.d9ac9228')" value="create" />
            <el-option :label="$t('auto.Enhanced.32ac152b')" value="update" />
            <el-option :label="$t('auto.Enhanced.2f4aaddd')" value="delete" />
            <el-option :label="$t('auto.Enhanced.402d19e5')" value="login" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.19fcb9eb')">
          <el-date-picker v-model="filters.dateRange" type="daterange" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="load">{{ $t('auto.Enhanced.bee912d7') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="logs" v-loading="loading">
        <el-table-column prop="createdAt" :label="$t('auto.Enhanced.19fcb9eb')" width="170" />
        <el-table-column prop="user" :label="$t('auto.Enhanced.f9ac4b2a')" width="100" />
        <el-table-column prop="action" :label="$t('auto.Enhanced.c500cfab')" width="100">
          <template #default="{ row }">
            <el-tag :type="actionType(row.action)" size="small">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" :label="$t('auto.Enhanced.eee83a92')" width="160" />
        <el-table-column prop="summary" :label="$t('auto.Enhanced.3ae14696')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="$t('auto.Enhanced.8d535108')" width="120">
          <template #default="{ row }">
            <el-button v-if="row.diff" link type="primary" @click="showDiff(row)">{{ $t('auto.Enhanced.0385ba93') }}</el-button>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Enhanced.2b6bc0f2')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="warning" :disabled="!row.rollbackable" @click="rollback(row)">{{ $t('auto.Enhanced.d00b485b') }}</el-button>
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.Enhanced.f26225bd') }}</el-button>
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

    
    <el-dialog v-model="diffVisible" :title="$t('auto.Enhanced.8d535108')" width="800px">
      <el-table :data="diffData">
        <el-table-column prop="field" :label="$t('auto.Enhanced.9caecd93')" width="160" />
        <el-table-column :label="$t('auto.Enhanced.405ac6e2')">
          <template #default="{ row }"><code>{{ row.before }}</code></template>
        </el-table-column>
        <el-table-column :label="$t('auto.Enhanced.8896d6e9')">
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
    await ElMessageBox.confirm(`确认回滚操作 ${row.id}？`, t('auto.Enhanced.d3dc0804'), { type: 'warning' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.post(`/api/operation-logs/${row.id}/rollback`, {})
  ElMessage.success(t('auto.Enhanced.ef6db9e1'))
  await load()
}

function viewDetail(row) {
  ElMessageBox.alert(JSON.stringify(row, null, 2), t('auto.Enhanced.f26225bd'))
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
