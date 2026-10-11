<template>
  <div class="backup-page">
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('auto.Enhanced.5355b3e6') }}</div>
          <div class="stat-value">{{ stats.total }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card success">
          <div class="stat-label">{{ $t('auto.Enhanced.b2f3fed4') }}</div>
          <div class="stat-value">{{ stats.lastSuccess }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card warning">
          <div class="stat-label">{{ $t('auto.Enhanced.ea9ecb93') }}</div>
          <div class="stat-value">{{ formatSize(stats.totalSize) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card primary">
          <div class="stat-label">{{ $t('auto.Enhanced.9fe67963') }}</div>
          <div class="stat-value" style="font-size: 16px">{{ stats.nextRun }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Enhanced.fe95755b') }}</span>
          <el-button type="primary" @click="saveStrategy">{{ $t('auto.Enhanced.2d3f52b4') }}</el-button>
        </div>
      </template>
      <el-form :model="strategy" label-width="120px">
        <el-form-item :label="$t('auto.Enhanced.590e5b46')">
          <el-switch v-model="strategy.enabled" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.ce221098')">
          <el-time-picker v-model="strategy.daily" placeholder="02:00" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.09b03b6d')">
          <el-select v-model="strategy.weeklyDay" style="width: 120px">
            <el-option :label="$t('auto.Enhanced.1603b069')" :value="1" />
            <el-option :label="$t('auto.Enhanced.562d7476')" :value="0" />
          </el-select>
          <el-time-picker v-model="strategy.weeklyTime" placeholder="03:00" style="margin-left: 8px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.824607be')">
          <el-input-number v-model="strategy.retentionDays" :min="7" :max="365" />
        </el-form-item>
        <el-form-item :label="$t('auto.Enhanced.ff59a874')">
          <el-switch v-model="strategy.checksum" />
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="restore-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Enhanced.9014d6d2') }}</span>
          <el-button type="danger" :loading="backing" @click="createBackup">{{ $t('auto.Enhanced.68c4a94d') }}</el-button>
        </div>
      </template>
      <el-table :data="backups" v-loading="loading">
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="createdAt" :label="$t('auto.Enhanced.ae0fd9b9')" width="170" />
        <el-table-column prop="size" :label="$t('auto.Enhanced.58f96667')" width="100">
          <template #default="{ row }">{{ formatSize(row.size) }}</template>
        </el-table-column>
        <el-table-column prop="checksum" label="SHA256" width="160" />
        <el-table-column prop="status" :label="$t('auto.Enhanced.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'ok' ? 'success' : 'danger'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Enhanced.2b6bc0f2')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="showPreview(row)">{{ $t('auto.Enhanced.9caf9d2a') }}</el-button>
            <el-button link type="warning" @click="restore(row)">{{ $t('auto.Enhanced.c7db6d4f') }}</el-button>
            <el-button link type="danger" @click="remove(row)">{{ $t('auto.Enhanced.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="previewVisible" :title="$t('auto.Enhanced.8552e7b2')" width="600px">
      <el-alert type="warning" :closable="false">
        即将恢复此备份：会影响以下数据（行数估算）
      </el-alert>
      <el-table :data="previewData">
        <el-table-column prop="table" :label="$t('auto.Enhanced.7ecd01a0')" />
        <el-table-column prop="rows" :label="$t('auto.Enhanced.f5d74e6a')" width="160" />
      </el-table>
      <template #footer>
        <el-button @click="previewVisible = false">{{ $t('auto.Enhanced.625fb26b') }}</el-button>
        <el-button type="danger" @click="confirmRestore">{{ $t('auto.Enhanced.841e5e89') }}</el-button>
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
  ElMessage.success(t('auto.Enhanced.2864c254'))
}

async function createBackup() {
  backing.value = true
  try {
    await http.post('/api/backup/create', {})
    ElMessage.success(t('auto.Enhanced.3eb08b6b'))
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
    await ElMessageBox.confirm(t('auto.Enhanced.7a1feec4'), t('auto.Enhanced.653fc421'), { type: 'error' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.post(`/api/backup/${_selectedBackup.id}/restore`, {})
  ElMessage.success(t('auto.Enhanced.0e7fa81a'))
  previewVisible.value = false
}

function restore(row) {
  _selectedBackup = row
  showPreview(row)
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确认删除备份 ${row.id}？`, t('auto.Enhanced.50eaf94d'), { type: 'warning' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.delete(`/api/backup/${row.id}`)
  ElMessage.success(t('auto.Enhanced.5cc23262'))
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
