<template>
  <div class="backup-page">
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('auto.Enhanced.15.5355b3') }}</div>
          <div class="stat-value">{{ stats.total }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card success">
          <div class="stat-label">{{ $t('auto.Enhanced.16.b2f3fe') }}</div>
          <div class="stat-value">{{ stats.lastSuccess }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card warning">
          <div class="stat-label">{{ $t('auto.Enhanced.17.ea9ecb') }}</div>
          <div class="stat-value">{{ formatSize(stats.totalSize) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card primary">
          <div class="stat-label">{{ $t('auto.Enhanced.18.9fe679') }}</div>
          <div class="stat-value" style="font-size: 16px">{{ stats.nextRun }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Enhanced.19.fe9575') }}</span>
          <el-button type="primary" @click="saveStrategy">{{ $t('auto.Enhanced.20.2d3f52') }}</el-button>
        </div>
      </template>
      <el-form :model="strategy" label-width="120px">
        <el-form-item :label="$t('auto.Enhanced.1.590e5b')">
          <el-switch v-model="strategy.enabled" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.2.ce2210')">
          <el-time-picker v-model="strategy.daily" placeholder="02:00" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.3.09b03b')">
          <el-select v-model="strategy.weeklyDay" style="width: 120px">
            <el-option :label="$t('auto.Enhanced.4.1603b0')" :value="1" />
            <el-option :label="$t('auto.Enhanced.5.562d74')" :value="0" />
          </el-select>
          <el-time-picker v-model="strategy.weeklyTime" placeholder="03:00" style="margin-left: 8px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.6.824607')">
          <el-input-number v-model="strategy.retentionDays" :min="7" :max="365" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.7.ff59a8')">
          <el-switch v-model="strategy.checksum" />
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="restore-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Enhanced.21.9014d6') }}</span>
          <el-button type="danger" :loading="backing" @click="createBackup">{{ $t('auto.Enhanced.22.68c4a9') }}</el-button>
        </div>
      </template>
      <el-table :data="backups" v-loading="loading">
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="createdAt" :label="$t('auto.Enhanced.8.ae0fd9')" width="170" />
        <el-table-column prop="size" :label="$t('auto.Enhanced.9.58f966')" width="100">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="checksum" label="SHA256" width="160" />
        <el-table-column prop="status" :label="$t('auto.Enhanced.10.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'ok' ? 'success' : 'danger'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Enhanced.11.2b6bc0')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="showPreview(row)">{{ $t('auto.Enhanced.23.9caf9d') }}</el-button>
            <el-button link type="warning" @click="restore(row)">{{ $t('auto.Enhanced.24.c7db6d') }}</el-button>
            <el-button link type="danger" @click="remove(row)">{{ $t('auto.Enhanced.25.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="previewVisible" :title="$t('auto.Enhanced.12.8552e7')" width="600px">
      <el-alert type="warning" :closable="false">
        即将恢复此备份：会影响以下数据（行数估算）
      </el-alert>
      <el-table :data="previewData">
        <el-table-column prop="table" :label="$t('auto.Enhanced.13.7ecd01')" />
        <el-table-column prop="rows" :label="$t('auto.Enhanced.14.f5d74e')" width="160" />
      </el-table>
      <template #footer>
        <el-button @click="previewVisible = false">{{ $t('auto.Enhanced.26.625fb2') }}</el-button>
        <el-button type="danger" @click="confirmRestore">{{ $t('auto.Enhanced.27.841e5e') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatFileSize } from '@/utils/format'
import { http } from '@/utils/request'

const { t } = useI18n()

const stats = ref({ total: 0, lastSuccess: '-', totalSize: 0, nextRun: '-' })
const strategy = reactive({
  enabled: true,
  daily: new Date(2024, 0, 1, 2, 0),
  weeklyDay: 0,
  weeklyTime: new Date(2024, 0, 1, 3, 0),
  retentionDays: 30,
  checksum: true
})
const backups = ref([])
const loading = ref(false)
const backing = ref(false)
const previewVisible = ref(false)
const previewData = ref([])
let _selectedBackup = null

const formatSize = formatFileSize

async function load() {
  loading.value = true
  try {
    const [s, b] = await Promise.all([
      http.get('/api/backup/stats'),
      http.get('/api/backup/list')
    ])
    stats.value = s || stats.value
    backups.value = b || []
  } finally {
    loading.value = false
  }
}

async function saveStrategy() {
  const d = strategy.daily instanceof Date ? strategy.daily : new Date(strategy.daily);
  const w = strategy.weeklyTime instanceof Date ? strategy.weeklyTime : new Date(strategy.weeklyTime)
  await http.put('/api/backup/strategy', {
    enabled: strategy.enabled,
    daily_hour: d.getHours(),
    daily_minute: d.getMinutes(),
    weekly_day: strategy.weeklyDay,
    weekly_hour: w.getHours(),
    weekly_minute: w.getMinutes(),
    retention_days: strategy.retentionDays,
    checksum: strategy.checksum
  })
  ElMessage.success(t('auto.Enhanced.28.2864c2'))
}

async function createBackup() {
  backing.value = true
  try {
    await http.post('/api/backup/create', {})
    ElMessage.success(t('auto.Enhanced.29.3eb08b'))
    setTimeout(load, 3000)
  } finally {
    backing.value = false
  }
}

async function showPreview(row) {
  _selectedBackup = row
  previewData.value = await http.get(`/api/backup/${row.id}/preview`)
  previewVisible.value = true
}

async function confirmRestore() {
  if (!_selectedBackup) return
  try {
    await ElMessageBox.confirm(t('auto.Enhanced.30.7a1fee'), t('auto.Enhanced.31.653fc4'), { type: 'error' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.post(`/api/backup/${_selectedBackup.id}/restore`, {})
  ElMessage.success(t('auto.Enhanced.32.0e7fa8'))
  previewVisible.value = false
}

function restore(row) {
  _selectedBackup = row
  showPreview(row)
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确认删除备份 ${row.id}？`, t('auto.Enhanced.33.50eaf9'), { type: 'warning' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.delete(`/api/backup/${row.id}`)
  ElMessage.success(t('auto.Enhanced.34.5cc232'))
  await load()
}

onMounted(load)
</script>

<style scoped>
.backup-page { padding: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; padding: 8px; }
.stat-label { color: #64748B; font-size: 12px; }
.stat-value { font-size: 28px; font-weight: 700; margin: 8px 0; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.restore-card { margin-top: 16px; }
</style>
