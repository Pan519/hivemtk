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
        <el-table-column prop="name" ::label="$t('auto.List.ba6d3462')" min-width="180" />
        <el-table-column prop="data_source" :label="$t('auto.List.c11322c9')" width="120">
          <template #default="{ row }">
            <el-tag>{{ dataSourceLabel(row.data_source) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="chart_type" :label="$t('auto.List.7aed2408')" width="110">
          <template #default="{ row }">
            <el-tag type="info">{{ chartTypeLabel(row.chart_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.f29c543f')" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="d in parseDimensions(row.dimensions)" :key="d.field" class="tag-item" size="small">{{ d.label }}</el-tag>
            <span v-if="!parseDimensions(row.dimensions).length" class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.7e687515')" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="m in parseMetrics(row.metrics)" :key="m.field" class="tag-item" size="small" type="warning">{{ m.label }}</el-tag>
            <span v-if="!parseMetrics(row.metrics).length" class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_by" :label="$t('auto.List.95a43eaa')" width="90" />
        <el-table-column prop="created_at" :label="$t('auto.List.eca37cb0')" min-width="170" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="300" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewReport(row)">{{ $t('auto.List.607e7a4f') }}</el-button>
            <el-button link type="primary" @click="exportReport(row)">{{ $t('auto.List.55405ea6') }}</el-button>
            <el-button link type="primary" @click="editReport(row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button link type="danger" @click="deleteReport(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="820px">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="110px">
        <el-form-item :label="$t('auto.List.a2e06320')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.List.948e172e')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.c11322c9')" prop="dataSource">
          <el-select v-model="form.dataSource" style="width: 100%">
            <el-option :label="$t('auto.List.98771135')" value="sessions" />
            <el-option :label="$t('auto.List.87846692')" value="messages" />
            <el-option :label="$t('auto.List.69a3cf32')" value="orders" />
            <el-option :label="$t('auto.List.68262f4c')" value="clues" />
            <el-option :label="$t('auto.List.0da9b129')" value="users" />
            <el-option :label="$t('auto.List.6727a802')" value="rfm" />
            <el-option :label="$t('auto.List.361fdc0c')" value="agents" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.7aed2408')" prop="chartType">
          <el-select v-model="form.chartType" style="width: 100%">
            <el-option :label="$t('auto.List.b339aa87')" value="table" />
            <el-option :label="$t('auto.List.c9c9bd87')" value="line" />
            <el-option :label="$t('auto.List.7c8ee6e3')" value="bar" />
            <el-option :label="$t('auto.List.d5e33db7')" value="pie" />
            <el-option :label="$t('auto.List.af99db99')" value="area" />
            <el-option :label="$t('auto.List.d87f215d')" value="card" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.d55f7fff')">
          <el-select v-model="form.dimensions" multiple style="width: 100%" :placeholder="$t('auto.List.cfa2c082')">
            <el-option v-for="d in DIM_OPTIONS" :key="d.value" :label="d.label" :value="d.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.8331f98e')">
          <el-select v-model="form.metrics" multiple style="width: 100%" :placeholder="$t('auto.List.4837e0e6')">
            <el-option v-for="m in METRIC_OPTIONS" :key="m.value" :label="m.label" :value="m.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.5148fcc0')">
          <div class="filters">
            <div v-for="(f, idx) in form.filters" :key="idx" class="filter-row">
              <el-input v-model="f.field" :placeholder="$t('auto.List.9caecd93')" style="width: 140px" />
              <el-select v-model="f.operator" style="width: 110px">
                <el-option :label="$t('auto.List.4c35bf2e')" value="eq" />
                <el-option :label="$t('auto.List.14a8af58')" value="ne" />
                <el-option :label="$t('auto.List.2791dca3')" value="gt" />
                <el-option :label="$t('auto.List.f09dc0d1')" value="lt" />
                <el-option :label="$t('auto.List.e13556bb')" value="like" />
              </el-select>
              <el-input v-model="f.value" :placeholder="$t('auto.List.fe7509e0')" style="width: 140px" />
              <el-button type="danger" text @click="removeFilter(idx)">{{ $t('auto.List.86048b4f') }}</el-button>
            </div>
            <el-button @click="addFilter">{{ $t('auto.List.3c09616c') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.List.2432b575')">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.b6a09844')">
          <el-switch v-model="form.isPublic" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">{{ $t('auto.List.38cf16f2') }}</el-button>
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
          <el-empty :description="$t('auto.List.58462e00')" />
        </template>
      </el-table>
      <template #footer>
        <el-button @click="viewDialogVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
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
  date: { field: 'date', label: t('auto.List.4ff1e74e'), data_type: 'date', group_by: true },
  source: { field: 'source', label: t('auto.List.26ca20b1'), data_type: 'string', group_by: true },
  platform: { field: 'platform', label: t('auto.List.4981287f'), data_type: 'string', group_by: true },
  region: { field: 'region', label: t('auto.List.2560b304'), data_type: 'string', group_by: true },
  user_type: { field: 'user_type', label: t('auto.List.d340e9a6'), data_type: 'string', group_by: true },
}
const METRIC_MAP = {
  count: { field: 'count', label: t('auto.List.0bf60b32'), agg_func: 'count', data_type: 'number' },
  amount: { field: 'amount', label: t('auto.List.4cf24aba'), agg_func: 'sum', data_type: 'number' },
  conversion: { field: 'conversion', label: t('auto.List.71505437'), agg_func: 'avg', data_type: 'percent' },
  click_rate: { field: 'click_rate', label: t('auto.List.e6867bb9'), agg_func: 'avg', data_type: 'percent' },
  open_rate: { field: 'open_rate', label: t('auto.List.eb9eb101'), agg_func: 'avg', data_type: 'percent' },
}
const DIM_OPTIONS = Object.entries(DIM_MAP).map(([value, v]) => ({ value, label: v.label }))
const METRIC_OPTIONS = Object.entries(METRIC_MAP).map(([value, v]) => ({ value, label: v.label }))
const DS_LABELS = {
  sessions: t('auto.List.98771135'), messages: t('auto.List.87846692'), orders: t('auto.List.69a3cf32'), clues: t('auto.List.68262f4c'),
  users: t('auto.List.0da9b129'), rfm: 'RFM 数据', agents: t('auto.List.361fdc0c'),
}
const CT_LABELS = { table: t('auto.List.b339aa87'), line: t('auto.List.c9c9bd87'), bar: t('auto.List.7c8ee6e3'), pie: t('auto.List.d5e33db7'), area: t('auto.List.af99db99'), card: t('auto.List.d87f215d') }

const { t } = useI18n()

const loading = ref(false)
const reports = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref(t('auto.List.02511943'))
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
  name: [{ required: true, message: i18n.global.t(t('auto.List.948e172e')), trigger: 'blur' }],
  dataSource: [{ required: true, message: i18n.global.t(t('auto.List.9acb966d')), trigger: 'change' }],
  chartType: [{ required: true, message: i18n.global.t(t('auto.List.31ceb5d2')), trigger: 'change' }],
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
    ElMessage.error(t('auto.List.c4fa9b74') + (e && e.message ? e.message : ''))
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  form.value = emptyForm()
  dialogTitle.value = t('auto.List.02511943')
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
  dialogTitle.value = t('auto.List.8ab69ff6')
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
      ElMessage.success(i18n.global.t(t('auto.List.33130f5c')))
      dialogVisible.value = false
      await refreshData()
    } catch (e) {
      ElMessage.error(i18n.global.t(t('auto.List.5fa802be')) + '：' + (e && e.message ? e.message : ''))
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
    ElMessage.error(i18n.global.t(t('auto.List.30f249a5')) + '：' + (e && e.message ? e.message : ''))
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
    ElMessage.error(t('auto.List.bccb121f') + (e && e.message ? e.message : ''))
  }
}

const deleteReport = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除报表 "${row.name}"？`, t('auto.List.e83a256e'), { type: 'warning' })
    await deleteCustomReport(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
    await refreshData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
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
