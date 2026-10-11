<template>
  <div class="backup-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('数据备份') }}</h2>
        <p class="subtitle">{{ $t('系统全量备份、自动调度、文件归档管理') }}</p>
      </div>
      <div>
        <el-button @click="loadData">
          <el-icon><Refresh /></el-icon>
          {{ $t('刷新') }}
        </el-button>
        <el-button type="primary" :loading="creating" @click="openCreateDialog">
          <el-icon><Plus /></el-icon>
          {{ $t('立即备份') }}
        </el-button>
      </div>
    </el-card>

    <el-row :gutter="20" class="stat-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('备份总数') }}</div>
            <div class="stat-value">{{ pagination.total }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('已完成') }}</div>
            <div class="stat-value success">{{ stats.completed }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('进行中') }}</div>
            <div class="stat-value warning">{{ stats.running }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('失败') }}</div>
            <div class="stat-value danger">{{ stats.failed }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      
      <el-tab-pane ::label="$t('auto.List.f5273794')" name="backup">
        <el-card>
          <div class="filter-bar">
            <el-input v-model="searchKeyword" ::placeholder="$t('auto.List.f6cd07b1')" clearable style="width: 220px" />
            <el-select v-model="filterType" ::placeholder="$t('auto.List.46041e44')" clearable style="width: 150px">
              <el-option ::label="$t('auto.List.ea46a347')" value="full" />
              <el-option ::label="$t('auto.List.d4b04d75')" value="incremental" />
              <el-option ::label="$t('auto.List.a2d10b79')" value="differential" />
            </el-select>
            <el-select v-model="filterStatus" ::placeholder="$t('auto.List.1b857887')" clearable style="width: 150px">
              <el-option ::label="$t('auto.List.3a3efebb')" value="pending" />
              <el-option ::label="$t('auto.List.9365602e')" value="running" />
              <el-option ::label="$t('auto.List.c98dc133')" value="completed" />
              <el-option ::label="$t('auto.List.f8d3aedb')" value="failed" />
            </el-select>
          </div>

          <el-table :data="filteredBackups" v-loading="loading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="backup_name" ::label="$t('auto.List.517d9787')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="backup_type" ::label="$t('auto.List.a0ee2507')" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="getTypeTag(row.backup_type)" size="small">
                  {{ getTypeText(row.backup_type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="file_size" :label="$t('auto.List.396b7d3f')" width="120" align="right">
              <template #default="{ row }">
                {{ formatSize(row.file_size) }}
              </template>
            </el-table-column>
            <el-table-column prop="started_at" :label="$t('auto.List.592c5958')" min-width="170">
              <template #default="{ row }">
                {{ formatTime(row.started_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="completed_at" :label="$t('auto.List.7cbadbbe')" min-width="170">
              <template #default="{ row }">
                {{ formatTime(row.completed_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="error_message" :label="$t('auto.List.4604d502')" min-width="180" show-overflow-tooltip />
            <el-table-column :label="$t('auto.List.2b6bc0f2')" width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
                <el-button
                  link
                  type="success"
                  :disabled="row.status !== 'completed'"
                  @click="openRestoreDialog(row)"
                >{{ $t('auto.List.c7db6d4f') }}</el-button>
                <el-button link type="danger" @click="doDeleteBackup(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.1f00e1f0')" />
            </template>
          </el-table>

          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="loadData"
            @size-change="loadData"
            style="margin-top: 15px; text-align: right"
          />
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.078ebd10')" name="restore">
        <el-card>
          <div class="filter-bar">
            <el-input v-model="restoreKeyword" :placeholder="$t('auto.List.8d67fe82')" clearable style="width: 220px" />
            <el-select v-model="restoreStatus" :placeholder="$t('auto.List.3fea7ca7')" clearable style="width: 150px">
              <el-option :label="$t('auto.List.6139a699')" value="pending" />
              <el-option :label="$t('auto.List.fb852fc6')" value="running" />
              <el-option :label="$t('auto.List.fad5222c')" value="completed" />
              <el-option :label="$t('auto.List.acd5cb84')" value="failed" />
            </el-select>
          </div>

          <el-table :data="filteredRestores" v-loading="restoreLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="backup_id" :label="$t('auto.List.d105690c')" width="100" />
            <el-table-column prop="backup_name" :label="$t('auto.List.77b9ecc8')" min-width="220" show-overflow-tooltip />
            <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="120" align="center">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="restored_at" :label="$t('auto.List.592c5958')" min-width="170">
              <template #default="{ row }">
                {{ formatTime(row.restored_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="error_message" :label="$t('auto.List.4604d502')" min-width="220" show-overflow-tooltip />
            <el-table-column :label="$t('auto.List.2b6bc0f2')" width="100" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="viewRestoreDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.da024546')" />
            </template>
          </el-table>

          <el-pagination
            v-model:current-page="restorePagination.page"
            v-model:page-size="restorePagination.size"
            :total="restorePagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="loadRestores"
            @size-change="loadRestores"
            style="margin-top: 15px; text-align: right"
          />
        </el-card>
      </el-tab-pane>
    </el-tabs>

    
    <el-dialog v-model="createVisible" :title="$t('auto.List.a4006e5c')" width="520px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item :label="$t('auto.List.77b9ecc8')" prop="backup_name">
          <el-input v-model="createForm.backup_name" :placeholder="$t('auto.List.1e66d8d9')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.8c61c92b')" prop="backup_type">
          <el-select v-model="createForm.backup_type" style="width: 100%">
            <el-option :label="$t('auto.List.6d68bd54')" value="full" />
            <el-option :label="$t('auto.List.43213875')" value="incremental" />
            <el-option :label="$t('auto.List.b23a912b')" value="differential" />
          </el-select>
        </el-form-item>
        <el-alert type="info" :closable="false" show-icon>
          <template #title>{{ $t('auto.List.bcbabc37') }}</template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="creating" @click="submitCreate">{{ $t('auto.List.68c4a94d') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="restoreVisible" :title="$t('auto.List.554d4999')" width="520px">
      <el-descriptions :column="1" border v-if="restoreTarget">
        <el-descriptions-item :label="$t('auto.List.d105690c')">{{ restoreTarget.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.77b9ecc8')">{{ restoreTarget.backup_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.8c61c92b')">{{ getTypeText(restoreTarget.backup_type) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.396b7d3f')">{{ formatSize(restoreTarget.file_size) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.7cbadbbe')">{{ formatTime(restoreTarget.completed_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-alert type="warning" :closable="false" show-icon style="margin-top: 12px">
        <template #title>{{ $t('auto.List.6aa24539') }}</template>
      </el-alert>
      <template #footer>
        <el-button @click="restoreVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="restoring" @click="submitRestore">{{ $t('auto.List.841e5e89') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailVisible" :title="$t('auto.List.4f9230d0')" width="600px">
      <el-descriptions :column="1" border v-if="detailRecord">
        <el-descriptions-item :label="$t('auto.List.d105690c')">{{ detailRecord.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.77b9ecc8')">{{ detailRecord.backup_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.8c61c92b')">{{ getTypeText(detailRecord.backup_type) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3fea7ca7')">
          <el-tag :type="getStatusType(detailRecord.status)" size="small">
            {{ getStatusText(detailRecord.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.de6104ca')">{{ detailRecord.file_path || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.396b7d3f')">{{ formatSize(detailRecord.file_size) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.95a43eaa')">{{ detailRecord.created_by || '系统' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.592c5958')">{{ formatTime(detailRecord.started_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.7cbadbbe')">{{ formatTime(detailRecord.completed_at) }}</el-descriptions-item>
        <el-descriptions-item v-if="detailRecord.error_message" :label="$t('auto.List.4604d502')">
          <span class="error-text">{{ detailRecord.error_message }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import {
  getBackupList,
  createBackup,
  deleteBackup,
  restoreBackup as callRestoreApi,
  getRestoreList
} from '@/api/backup.js'
import {
  BACKUP_STATUS,
  getStatusLabel as getStatusLabelFromArr,
  getStatusTagType as getStatusTagTypeFromArr,
} from '@/constants/status';

const { t } = useI18n()

const activeTab = ref('backup')

const loading = ref(false);
const backups = ref([])
const pagination = reactive({ page: 1, size: 20, total: 0 })
const searchKeyword = ref('')
const filterType = ref('')
const filterStatus = ref('')

const restoreLoading = ref(false);
const restores = ref([])
const restorePagination = reactive({ page: 1, size: 20, total: 0 })
const restoreKeyword = ref('')
const restoreStatus = ref('')

const createVisible = ref(false);
const creating = ref(false)
const createFormRef = ref()
const createForm = reactive({
  backup_name: '',
  backup_type: 'full'
})
const createRules = {
  backup_type: [{ required: true, message: i18n.global.t(t('auto.List.43ebf364')), trigger: 'change' }]
}

const restoreVisible = ref(false);
const restoring = ref(false)
const restoreTarget = ref(null)

const detailVisible = ref(false);
const detailRecord = ref(null)

const stats = computed(() => {
  const result = { completed: 0, running: 0, failed: 0, pending: 0 }
  backups.value.forEach(b => {
    if (b.status === 'completed') result.completed++
    else if (b.status === 'running') result.running++
    else if (b.status === 'failed') result.failed++
    else if (b.status === 'pending') result.pending++
  })
  return result
});

const filteredBackups = computed(() => {
  let result = backups.value
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(b => b.backup_name?.toLowerCase().includes(kw))
  }
  if (filterType.value) result = result.filter(b => b.backup_type === filterType.value)
  if (filterStatus.value) result = result.filter(b => b.status === filterStatus.value)
  return result
})

const filteredRestores = computed(() => {
  let result = restores.value
  if (restoreKeyword.value) {
    const kw = restoreKeyword.value.toLowerCase()
    result = result.filter(r => r.backup_name?.toLowerCase().includes(kw))
  }
  if (restoreStatus.value) result = result.filter(r => r.status === restoreStatus.value)
  return result
})

const getTypeText = (type) => {
  const map = { full: t('auto.List.c60ad696'), incremental: t('auto.List.36daebe0'), differential: t('auto.List.5fdfbade') }
  return map[type] || type || '-'
}
const getTypeTag = (type) => {
  const map = { full: 'primary', incremental: 'success', differential: 'warning' }
  return map[type] || 'info'
}
const getStatusText = (status) => getStatusLabelFromArr(status, BACKUP_STATUS);
const getStatusType = (status) => getStatusTagTypeFromArr(status, BACKUP_STATUS)
const formatSize = (bytes) => {
  if (!bytes || bytes <= 0) return '-'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return `${size.toFixed(2)} ${units[i]}`
}
const formatTime = (val) => {
  if (!val) return '-'
  try {
    const d = new Date(val)
    if (isNaN(d.getTime())) return val
    return d.toLocaleString('zh-CN', { hour12: false })
  } catch (e) {
    return val
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getBackupList({
      page: pagination.page,
      page_size: pagination.size
    })
    const data = res || {}
    backups.value = data.list || []
    pagination.total = data.total || 0
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.e29eced0')))
    backups.value = []
  } finally {
    loading.value = false
  }
}

const loadRestores = async () => {
  restoreLoading.value = true
  try {
    const res = await getRestoreList({
      page: restorePagination.page,
      page_size: restorePagination.size
    })
    const data = res || {}
    restores.value = data.list || []
    restorePagination.total = data.total || 0
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.b9ab2bdf')))
    restores.value = []
  } finally {
    restoreLoading.value = false
  }
}

const onTabChange = (name) => {
  if (name === 'restore' && restores.value.length === 0) loadRestores()
}

const openCreateDialog = () => {
  createForm.backup_name = ''
  createForm.backup_type = 'full'
  createVisible.value = true
}

const submitCreate = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    creating.value = true
    try {
      await createBackup({
        backup_name: createForm.backup_name,
        backup_type: createForm.backup_type
      })
      ElMessage.success(i18n.global.t(t('auto.List.72f1633b')))
      createVisible.value = false
      await loadData()
    } catch (e) {
      ElMessage.error(e?.message || t('auto.List.c33cb565'))
    } finally {
      creating.value = false
    }
  })
}

const openRestoreDialog = (row) => {
  restoreTarget.value = row
  restoreVisible.value = true
}

const submitRestore = async () => {
  if (!restoreTarget.value) return
  restoring.value = true
  try {
    await callRestoreApi({ backup_id: restoreTarget.value.id })
    ElMessage.success(i18n.global.t(t('auto.List.c740083a')))
    restoreVisible.value = false
    activeTab.value = 'restore'
    await loadRestores()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.List.2e827d23'))
  } finally {
    restoring.value = false
  }
}

const doDeleteBackup = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定删除备份 "${row.backup_name}"？删除后不可恢复`,
      t('auto.List.900c70fa'),
      { type: 'warning' }
    )
    await deleteBackup(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
    await loadData()
  } catch (e) {
    if (e !== 'cancel' && e !== 'close') ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
  }
}

const viewDetail = (row) => {
  detailRecord.value = row
  detailVisible.value = true
}

const viewRestoreDetail = (row) => {
  detailRecord.value = row
  detailVisible.value = true
}

onMounted(() => loadData())
</script>

<style scoped lang="scss">
.backup-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.stat-row {
  margin-bottom: 20px;
  .stat-item {
    text-align: center;
    .stat-label { color: #909399; font-size: 14px; margin-bottom: 10px; }
    .stat-value { font-size: 28px; font-weight: bold; }
    .stat-value.success { color: #10B981; }
    .stat-value.warning { color: #F59E0B; }
    .stat-value.danger { color: #EF4444; }
  }
}
.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
  flex-wrap: wrap;
}
.error-text { color: #EF4444; }
</style>
