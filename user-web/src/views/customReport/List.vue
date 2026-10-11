<template>
  <div class="custom-report-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('自定义报表') }}</h2>
        <p class="subtitle">{{ $t('自定义维度、指标生成专属业务报表') }}</p>
      </div>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('新建报表') }}
      </el-button>
    </el-card>

    <el-card>
      <el-table :data="reports" v-loading="loading" stripe>
        <el-table-column prop="name" ::label="$t('auto.List.1.ba6d34')" min-width="180" />
        <el-table-column prop="data_source" :label="$t('auto.List.2.c11322')" width="120">
          <template #default="{ row }">
            <el-tag>{{ dataSourceLabel(row.data_source) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="chart_type" :label="$t('auto.List.3.7aed24')" width="110">
          <template #default="{ row }">
            <el-tag type="info">{{ chartTypeLabel(row.chart_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.4.f29c54')" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="d in parseDimensions(row.dimensions)" :key="d.field" class="tag-item" size="small">{{ d.label }}</el-tag>
            <span v-if="!parseDimensions(row.dimensions).length" class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.5.7e6875')" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="m in parseMetrics(row.metrics)" :key="m.field" class="tag-item" size="small" type="warning">{{ m.label }}</el-tag>
            <span v-if="!parseMetrics(row.metrics).length" class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_by" :label="$t('auto.List.6.95a43e')" width="90" />
        <el-table-column prop="created_at" :label="$t('auto.List.7.eca37c')" min-width="170" />
        <el-table-column :label="$t('auto.List.8.2b6bc0')" width="300" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewReport(row)">{{ $t('auto.List.39.607e7a') }}</el-button>
            <el-button link type="primary" @click="exportReport(row)">{{ $t('auto.List.40.55405e') }}</el-button>
            <el-button link type="primary" @click="editReport(row)">{{ $t('auto.List.41.95b351') }}</el-button>
            <el-button link type="danger" @click="deleteReport(row)">{{ $t('auto.List.42.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="820px">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item :label="$t('auto.List.9.a2e063')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.List.10.948e17')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.2.c11322')" prop="dataSource">
          <el-select v-model="form.dataSource" style="width: 100%">
            <el-option :label="$t('auto.List.11.987711')" value="sessions" />
            <el-option :label="$t('auto.List.12.878466')" value="messages" />
            <el-option :label="$t('auto.List.13.69a3cf')" value="orders" />
            <el-option :label="$t('auto.List.14.68262f')" value="clues" />
            <el-option :label="$t('auto.List.15.0da9b1')" value="users" />
            <el-option :label="$t('auto.List.16.6727a8')" value="rfm" />
            <el-option :label="$t('auto.List.17.361fdc')" value="agents" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.3.7aed24')" prop="chartType">
          <el-select v-model="form.chartType" style="width: 100%">
            <el-option :label="$t('auto.List.18.b339aa')" value="table" />
            <el-option :label="$t('auto.List.19.c9c9bd')" value="line" />
            <el-option :label="$t('auto.List.20.7c8ee6')" value="bar" />
            <el-option :label="$t('auto.List.21.d5e33d')" value="pie" />
            <el-option :label="$t('auto.List.22.af99db')" value="area" />
            <el-option :label="$t('auto.List.23.d87f21')" value="card" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.24.d55f7f')">
          <el-select v-model="form.dimensions" multiple style="width: 100%" :placeholder="$t('auto.List.25.cfa2c0')">
            <el-option v-for="d in DIM_OPTIONS" :key="d.value" :label="d.label" :value="d.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.26.8331f9')">
          <el-select v-model="form.metrics" multiple style="width: 100%" :placeholder="$t('auto.List.27.4837e0')">
            <el-option v-for="m in METRIC_OPTIONS" :key="m.value" :label="m.label" :value="m.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.28.5148fc')">
          <div class="filters">
            <div v-for="(f, idx) in form.filters" :key="idx" class="filter-row">
              <el-input v-model="f.field" :placeholder="$t('auto.List.29.9caecd')" style="width: 140px" />
              <el-select v-model="f.operator" style="width: 110px">
                <el-option :label="$t('auto.List.30.4c35bf')" value="eq" />
                <el-option :label="$t('auto.List.31.14a8af')" value="ne" />
                <el-option :label="$t('auto.List.32.2791dc')" value="gt" />
                <el-option :label="$t('auto.List.33.f09dc0')" value="lt" />
                <el-option :label="$t('auto.List.34.e13556')" value="like" />
              </el-select>
              <el-input v-model="f.value" :placeholder="$t('auto.List.35.fe7509')" style="width: 140px" />
              <el-button type="danger" text @click="removeFilter(idx)">{{ $t('auto.List.43.86048b') }}</el-button>
            </div>
            <el-button @click="addFilter">{{ $t('auto.List.44.3c0961') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.List.36.2432b5')">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.37.b6a098')">
          <el-switch v-model="form.isPublic" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.45.625fb2') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">{{ $t('auto.List.46.38cf16') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="viewDialogVisible" :title="`查看报表: ${currentReportName}`" width="1000px">
      <el-table :data="reportData" v-loading="viewLoading" stripe border max-height="500px">
        <el-table-column
          v-for="col in reportColumns"
          :key="col.key"
          :prop="col.key"
          :label="col.label"
          :min-width="col.minWidth"
          show-overflow-tooltip
        />
        <template #empty>
          <el-empty :description="$t('auto.List.38.58462e')" />
        </template>
      </el-table>
      <template #footer>
        <el-button @click="viewDialogVisible = false">{{ $t('auto.List.47.b15d91') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'
import { ref, onMounted, onActivated } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { toList } from '@/utils/list.js'
import {
  getCustomReports,
  createCustomReport,
  updateCustomReport,
  deleteCustomReport,
  exportCustomReport,
  queryReportData,
} from '@/api/customReport.js'

const DIM_MAP = {
  date: { field: 'date', label: t('auto.List.48.4ff1e7'), data_type: 'date', group_by: true },
  source: { field: 'source', label: t('auto.List.49.26ca20'), data_type: 'string', group_by: true },
  platform: { field: 'platform', label: t('auto.List.50.498128'), data_type: 'string', group_by: true },
  region: { field: 'region', label: t('auto.List.51.2560b3'), data_type: 'string', group_by: true },
  user_type: { field: 'user_type', label: t('auto.List.52.d340e9'), data_type: 'string', group_by: true },
}
const METRIC_MAP = {
  count: { field: 'count', label: t('auto.List.53.0bf60b'), agg_func: 'count', data_type: 'number' },
  amount: { field: 'amount', label: t('auto.List.54.4cf24a'), agg_func: 'sum', data_type: 'number' },
  conversion: { field: 'conversion', label: t('auto.List.55.715054'), agg_func: 'avg', data_type: 'percent' },
  click_rate: { field: 'click_rate', label: t('auto.List.56.e6867b'), agg_func: 'avg', data_type: 'percent' },
  open_rate: { field: 'open_rate', label: t('auto.List.57.eb9eb1'), agg_func: 'avg', data_type: 'percent' },
}
const DIM_OPTIONS = Object.entries(DIM_MAP).map(([value, v]) => ({ value, label: v.label }))
const METRIC_OPTIONS = Object.entries(METRIC_MAP).map(([value, v]) => ({ value, label: v.label }))
const DS_LABELS = {
  sessions: t('auto.List.11.987711'), messages: t('auto.List.12.878466'), orders: t('auto.List.13.69a3cf'), clues: t('auto.List.14.68262f'),
  users: t('auto.List.15.0da9b1'), rfm: 'RFM 数据', agents: t('auto.List.17.361fdc'),
}
const CT_LABELS = { table: t('auto.List.18.b339aa'), line: t('auto.List.19.c9c9bd'), bar: t('auto.List.20.7c8ee6'), pie: t('auto.List.21.d5e33d'), area: t('auto.List.22.af99db'), card: t('auto.List.23.d87f21') }

const { t } = useI18n()

const loading = ref(false)
const reports = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref(t('auto.List.58.025119'))
const submitting = ref(false)
const formRef = ref()
const viewDialogVisible = ref(false)
const viewLoading = ref(false)
const currentReportName = ref('')
const reportData = ref([])
const reportColumns = ref([])

const form = ref(emptyForm())
function emptyForm() {
  return {
    id: 0,
    name: '',
    dataSource: 'clues',
    chartType: 'table',
    dimensions: [],
    metrics: [],
    filters: [],
    description: '',
    isPublic: true,
  }
}

const formRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.List.10.948e17')), trigger: 'blur' }],
  dataSource: [{ required: true, message: i18n.global.t(t('auto.List.59.9acb96')), trigger: 'change' }],
  chartType: [{ required: true, message: i18n.global.t(t('auto.List.60.31ceb5')), trigger: 'change' }],
}

const dataSourceLabel = (v) => DS_LABELS[v] || v || '-'
const chartTypeLabel = (v) => CT_LABELS[v] || v || '-'
const parseDimensions = (raw) => safeParseArray(raw).map((d) => (typeof d === 'string' ? { field: d, label: d } : d)).filter(Boolean)
const parseMetrics = (raw) => safeParseArray(raw).map((m) => (typeof m === 'string' ? { field: m, label: m } : m)).filter(Boolean)
function safeParseArray(raw) {
  if (Array.isArray(raw)) return raw
  if (typeof raw !== 'string' || !raw) return []
  try { return JSON.parse(raw) } catch (e) { return [] }
}

const refreshData = async () => {
  loading.value = true
  try {
    const res = await getCustomReports()
    reports.value = toList(res?.data ?? res)
  } catch (e) {
    ElMessage.error(t('auto.List.61.c4fa9b') + (e && e.message ? e.message : ''))
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  form.value = emptyForm()
  dialogTitle.value = t('auto.List.58.025119')
  dialogVisible.value = true
}

const addFilter = () => form.value.filters.push({ field: '', operator: 'eq', value: '' })
const removeFilter = (idx) => form.value.filters.splice(idx, 1)

const editReport = (row) => {
  form.value = {
    id: row.id,
    name: row.name || '',
    dataSource: row.data_source || 'clues',
    chartType: row.chart_type || 'table',
    dimensions: parseDimensions(row.dimensions).map((d) => d.field),
    metrics: parseMetrics(row.metrics).map((m) => m.field),
    filters: safeParseArray(row.filters).map((f) => ({ field: f.field || '', operator: f.operator || 'eq', value: f.value != null ? String(f.value) : '' })),
    description: row.description || '',
    isPublic: !!row.is_public,
  }
  dialogTitle.value = t('auto.List.62.8ab69f')
  dialogVisible.value = true
}

const buildPayload = () => {
  const dimensions = form.value.dimensions.map((v) => DIM_MAP[v]).filter(Boolean)
  const metrics = form.value.metrics.map((v) => METRIC_MAP[v]).filter(Boolean)
  const filters = form.value.filters
    .filter((f) => f.field && f.value !== '' && f.value !== null)
    .map((f) => ({ field: f.field, operator: f.operator, value: f.value }))
  return {
    name: form.value.name,
    description: form.value.description,
    data_source: form.value.dataSource,
    chart_type: form.value.chartType,
    dimensions,
    metrics,
    filters,
    is_public: form.value.isPublic,
  }
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const payload = buildPayload()
      if (form.value.id) {
        await updateCustomReport(form.value.id, payload)
      } else {
        await createCustomReport(payload)
      }
      ElMessage.success(i18n.global.t(t('auto.List.63.33130f')))
      dialogVisible.value = false
      await refreshData()
    } catch (e) {
      ElMessage.error(i18n.global.t(t('auto.List.64.5fa802')) + '：' + (e && e.message ? e.message : ''))
    } finally {
      submitting.value = false
    }
  })
}

const viewReport = async (row) => {
  currentReportName.value = row.name
  viewDialogVisible.value = true
  viewLoading.value = true
  reportData.value = []
  reportColumns.value = []
  try {
    const res = await queryReportData(row.id)
    const payload = res?.data ?? res
    const rows = Array.isArray(payload) ? payload : (payload.data || [])
    if (rows.length) {
      const keys = Object.keys(rows[0])
      reportColumns.value = keys.map((k) => ({ key: k, label: k, minWidth: 120 }))
      reportData.value = rows
    } else if (payload && payload.dimensions) {
      reportColumns.value = (payload.dimensions || []).map((d) => ({ key: d, label: d, minWidth: 120 }))
    }
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.65.30f249')) + '：' + (e && e.message ? e.message : ''))
  } finally {
    viewLoading.value = false
  }
}

const exportReport = async (row) => {
  try {
    const res = await exportCustomReport(row.id)
    const payload = res?.data ?? res
    const rows = Array.isArray(payload) ? payload : (payload && payload.data) || []
    ElMessage.success(`导出成功，共 ${rows.length} 行`)
  } catch (e) {
    ElMessage.error(t('auto.List.66.bccb12') + (e && e.message ? e.message : ''))
  }
}

const deleteReport = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除报表 "${row.name}"？`, t('auto.List.67.e83a25'), { type: 'warning' })
    await deleteCustomReport(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.68.0007d1')))
    await refreshData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.69.acf066')))
  }
}

onMounted(() => refreshData())
onActivated(() => refreshData())
</script>

<style scoped lang="scss">
.custom-report-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.tag-item { margin-right: 6px; }
.muted { color: #c0c4cc; }
.filters { display: flex; flex-direction: column; gap: 8px; }
.filter-row { display: flex; align-items: center; gap: 8px; }
</style>
