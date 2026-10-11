<template>
  <div class="clue-list-page">
    
    <el-card class="header-card" shadow="never">
      <div class="page-header">
        <div class="header-text">
          <h2>{{ $t('线索列表') }}</h2>
          <p class="subtitle">{{ $t('auto.List.19f53c26') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="fetchCluelist">
            <el-icon><Refresh /></el-icon>
            {{ $t('刷新') }}
          </el-button>
          <el-button type="primary" @click="openImportDialog">
            <el-icon><Upload /></el-icon>
            {{ $t('导入线索') }}
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-label">{{ $t('auto.List.fa01bc43') }}</div>
          <div class="stat-value">{{ stats.total || cluetotal || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-label">{{ $t('auto.List.876caf44') }}</div>
          <div class="stat-value" style="color: #10B981">{{ stats.verified || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-label">{{ $t('auto.List.b928ea24') }}</div>
          <div class="stat-value" style="color: #F59E0B">{{ stats.unverified || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-label">{{ $t('auto.List.b29a54ea') }}</div>
          <div class="stat-value" style="color: #4F46E5">{{ stats.today || 0 }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item :label="$t('auto.List.cfb5f18c')">
          <el-input
            v-model="filterForm.keyword"
            :placeholder="$t('auto.List.cebd95a0')"
            clearable
            style="width: 240px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.1cc291d3')">
          <el-select v-model="filterForm.type" :placeholder="$t('auto.List.c079e7d5')" clearable style="width: 160px">
            <el-option v-for="t in clueTypeMap" :key="t.value" :label="t.label" :value="String(t.value)" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.644603e5')">
          <el-select v-model="filterForm.verified" :placeholder="$t('auto.List.a8b0c204')" clearable style="width: 140px">
            <el-option :label="$t('auto.List.876caf44')" :value="1" />
            <el-option :label="$t('auto.List.b928ea24')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.19fcb9eb')">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            :range-separator="$t('auto.List.981cbe31')"
            :start-placeholder="$t('auto.List.a3e3b883')"
            :end-placeholder="$t('auto.List.12f1d7ef')"
            value-format="YYYY-MM-DD"
            style="width: 280px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshRight /></el-icon>重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card class="table-card" shadow="never">
      <el-table
        :data="filteredList"
        v-loading="loading"
        stripe
        ::empty-text="$t('auto.List.086bbe59')"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('auto.List.d7ec2d3f')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="account" :label="$t('auto.List.7035c62f')" min-width="160" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.226b0912')" min-width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getClueType(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="city" :label="$t('auto.List.f7d29dfa')" min-width="100" show-overflow-tooltip />
        <el-table-column prop="address" :label="$t('auto.List.7650487a')" min-width="180" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.26ca20b1')" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag :type="getChannelTagType(row.source)" effect="plain" size="small">
              {{ getChannelLabel(row.source) || row.source || '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.644603e5')" min-width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="(row.is_verify ?? row.verified) >= 1 ? 'success' : 'info'" size="small">
              {{ (row.is_verify ?? row.verified) >= 1 ? '已验证' : '未验证' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.List.eca37cb0')" min-width="170" show-overflow-tooltip>
          <template #default="{ row }">{{ formatTime(row.created_at || row.createdAt) }}</template>
        </el-table-column>
        <el-table-column prop="last_follow_at" :label="$t('auto.List.3ed5b8a2')" min-width="170" show-overflow-tooltip>
          <template #default="{ row }">{{ formatTime(row.last_follow_at || row.lastFollowAt) }}</template>
        </el-table-column>
        <el-table-column prop="assignee" :label="$t('auto.List.3b22c931')" min-width="120" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" min-width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewClue(row)">{{ $t('auto.List.607e7a4f') }}</el-button>
            <el-button link type="primary" size="small" @click="assignClue(row)">{{ $t('auto.List.e39de337') }}</el-button>
            <el-button link type="warning" size="small" @click="editClue(row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button link type="danger" size="small" @click="deleteClue(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          :current-page="cluepage"
          :page-size="cluelimit"
          :page-sizes="[10, 20, 50, 100]"
          :total="cluetotal"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    
    <el-dialog v-model="importDialogVisible" ::title="$t('auto.List.9b7bd118')" width="600px">
      <el-form ref="importFormRef" :model="importForm" label-width="80px">
        <el-form-item ::label="$t('auto.List.4d6ae67f')" prop="type">
          <el-select v-model="importForm.type" ::placeholder="$t('auto.List.24487a58')" value-key="value" @change="handleTypeChange" style="width: 100%">
            <el-option v-for="item in clueTypeMap" :key="String(item.value)" :label="item.label" :value="String(item.value)"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.512165ef')" prop="content">
          <el-input
            type="textarea"
            v-model="importForm.content"
            ::placeholder="$t('auto.List.3db0dfe7')"
            :rows="10"
            class="import-textarea"
            @input="handleTextareaInput"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importDialogVisible = false">{{ $t('取消') }}</el-button>
        <el-button type="primary" @click="submitImport">{{ $t('导入') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailVisible" :title="$t('auto.List.f0e2686b')" width="640px">
      <el-descriptions v-if="current" :column="2" border>
        <el-descriptions-item label="ID">{{ current.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.d7ec2d3f')">{{ current.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.7035c62f')">{{ current.account }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.226b0912')">{{ getClueType(current.type) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.f7d29dfa')">{{ current.city || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.7650487a')">{{ current.address || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.26ca20b1')">{{ getChannelLabel(current.source) || current.source || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.cd8992b6')">{{ (current.is_verify ?? current.verified) >= 1 ? '已验证' : '未验证' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.eca37cb0')">{{ formatTime(current.created_at || current.createdAt) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.3b22c931')">{{ current.assignee || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, RefreshRight, Upload } from '@element-plus/icons-vue'
import { clueApi } from '@/api/clue'
import { getChannelLabel, getChannelTagType } from '@/constants/channel'
import { getClueTypeLabel, CLUE_TYPE_OPTIONS } from '@/constants/cardPlatform';

const { t } = useI18n()

const cluelist = ref([]);
const cluetotal = ref(0)
const cluepage = ref(1)
const cluelimit = ref(10)
const loading = ref(false)
const clueTypeMap = ref(CLUE_TYPE_OPTIONS.slice())
const stats = ref({ total: 0, verified: 0, unverified: 0, today: 0 })

const filterForm = reactive({
  keyword: '',
  type: '',
  verified: ''
});
const dateRange = ref([])

const filteredList = computed(() => {
  let arr = cluelist.value.slice()
  const k = (filterForm.keyword || '').trim().toLowerCase()
  if (k) {
    arr = arr.filter(
      (c) =>
        (c.name || '').toLowerCase().includes(k) ||
        (c.account || '').toLowerCase().includes(k) ||
        (c.city || '').toLowerCase().includes(k) ||
        (c.address || '').toLowerCase().includes(k)
    )
  }
  if (filterForm.type) {
    arr = arr.filter((c) => String(c.type) === String(filterForm.type))
  }
  if (filterForm.verified !== '' && filterForm.verified !== null) {
    arr = arr.filter((c) => Number(c.is_verify ?? c.verified) === Number(filterForm.verified))
  }
  if (Array.isArray(dateRange.value) && dateRange.value.length === 2) {
    const [s, e] = dateRange.value
    arr = arr.filter((c) => {
      const t = (c.created_at || c.createdAt || '').slice(0, 10)
      return t >= s && t <= e
    })
  }
  return arr
});

const getClueType = (type) => getClueTypeLabel(type)

const handlePageChange = (val) => {
  if (cluepage.value === val) return
  cluepage.value = val
  fetchCluelist()
}

const handleSizeChange = (val) => {
  cluelimit.value = val
  cluepage.value = 1
  fetchCluelist()
}

const handleSearch = () => {
  cluepage.value = 1
  fetchCluelist()
}

const resetSearch = () => {
  filterForm.keyword = ''
  filterForm.type = ''
  filterForm.verified = ''
  dateRange.value = []
  cluepage.value = 1
  fetchCluelist()
}

const fetchCluelist = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await clueApi.list(cluepage.value, cluelimit.value)
    cluelist.value = res.list || []
    cluetotal.value = res.total || 0
    refreshStats()
  } catch (error) {
    console.error(t('auto.List.0c189758'), error)
    cluelist.value = []
    cluetotal.value = 0
  } finally {
    loading.value = false
  }
}

const refreshStats = () => {
  const list = cluelist.value || []
  const today = new Date().toISOString().slice(0, 10)
  stats.value = {
    total: cluetotal.value || list.length,
    verified: list.filter((c) => (c.is_verify ?? c.verified) >= 1).length,
    unverified: list.filter((c) => (c.is_verify ?? c.verified) < 1).length,
    today: list.filter((c) => (c.created_at || c.createdAt || '').slice(0, 10) === today).length
  }
}

const deleteClue = async (clue) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除线索 ${clue.name} 吗？`,
      t('auto.List.900c70fa'),
      { confirmButtonText: t('auto.List.38cf16f2'), cancelButtonText: t('auto.List.625fb26b'), type: 'warning' }
    )
  } catch {
    return
  }
  try {
    await clueApi.delete(clue.id)
    ElMessage.success(i18n.global.t(t('auto.List.22ea51a4')))
    fetchCluelist()
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
  }
}

const viewClue = (row) => {
  current.value = row
  detailVisible.value = true
}
const assignClue = (row) => {
  ElMessage.info(i18n.global.t(t('auto.List.ecaf351f')) + (row.name || row.id))
}
const editClue = (row) => {
  ElMessage.info(i18n.global.t(t('auto.List.3c685a52')) + (row.name || row.id))
}

const detailVisible = ref(false);
const current = ref(null)

const importDialogVisible = ref(false);
const importForm = ref({ type: '', content: '' })
const importFormRef = ref()

const openImportDialog = () => {
  importForm.value = { type: '', content: '' }
  importDialogVisible.value = true
}

const submitImport = async () => {
  if (!importForm.value.type) {
    ElMessage.error(i18n.global.t(t('auto.List.bd144be3')))
    return
  }
  if (!importForm.value.content) {
    ElMessage.error(i18n.global.t(t('auto.List.0183aa26')))
    return
  }
  const lines = importForm.value.content.split('\n')
  const clues = []
  for (const line of lines) {
    if (!line.trim()) continue
    const parts = line.split(',')
    if (parts.length < 2) {
      ElMessage.warning(`跳过无效格式行: ${line}`)
      continue
    }
    clues.push({
      name: parts[0].trim(),
      account: parts[1].trim(),
      city: parts.length > 2 ? parts[2].trim() : '',
      address: parts.length > 3 ? parts[3].trim() : '',
      type: importForm.value.type
    })
  }
  if (clues.length === 0) {
    ElMessage.warning(i18n.global.t(t('auto.List.a0779b51')))
    return
  }
  try {
    const res = await clueApi.import(clues)
    ElMessage.success(`成功导入 ${res.successCount} 条线索，跳过 ${res.skipCount} 条重复数据`)
    importDialogVisible.value = false
    importForm.value = { type: '', content: '' }
    fetchCluelist()
  } catch (error) {
    console.error(t('auto.List.43326d5b'), error)
    ElMessage.error(i18n.global.t(t('auto.List.4451a3ca')))
  }
}

const handleTextareaInput = (value) => {
  importForm.value.content = value
}
const handleTypeChange = (value) => {
  importForm.value.type = String(value)
}

const formatTime = (val) => {
  if (!val) return '-'
  try {
    const d = new Date(val)
    if (Number.isNaN(d.getTime())) return val
    const pad = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
  } catch (e) {
    return val
  }
};

onMounted(() => {
  fetchCluelist()
})
</script>

<style lang="scss" scoped>
.clue-list-page {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .header-card {
    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 16px;
      .header-text {
        flex: 1;
        h2 {
          margin: 0;
          font-size: 22px;
          color: #303133;
        }
        .subtitle {
          margin: 6px 0 0;
          font-size: 13px;
          color: #909399;
        }
      }
      .header-actions {
        display: flex;
        gap: 8px;
      }
    }
  }

  .stats-row {
    .stat-card {
      .stat-label {
        font-size: 13px;
        color: #909399;
      }
      .stat-value {
        font-size: 24px;
        font-weight: 600;
        margin-top: 8px;
        color: #303133;
      }
    }
  }

  .filter-card {
    .filter-form {
      margin: 0;
    }
  }

  .table-card {
    .pagination-container {
      margin-top: 16px;
      display: flex;
      justify-content: center;
    }
  }
}
</style>
