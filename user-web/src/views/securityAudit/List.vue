<template>
  <div class="audit-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('安全审计') }}</h2>
        <p class="subtitle">{{ $t('系统安全检查、风险评估、异常告警追踪') }}</p>
      </div>
      <div>
        <el-button @click="loadData">
          <el-icon><Refresh /></el-icon>
          {{ $t('刷新') }}
        </el-button>
        <el-button type="primary" :loading="running" @click="runAudit">
          <el-icon><VideoPlay /></el-icon>
          {{ $t('立即审计') }}
        </el-button>
      </div>
    </el-card>

    <el-row :gutter="20" class="stat-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('审计总数') }}</div>
            <div class="stat-value">{{ pagination.total }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('高风险') }}</div>
            <div class="stat-value danger">{{ stats.high }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('中风险') }}</div>
            <div class="stat-value warning">{{ stats.medium }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('低风险') }}</div>
            <div class="stat-value success">{{ stats.low }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <div class="filter-bar">
        <el-input v-model="searchKeyword" ::placeholder="$t('auto.List.1.f5230c')" clearable style="width: 220px" />
        <el-select v-model="filterRisk" ::placeholder="$t('auto.List.2.4a708d')" clearable style="width: 150px">
          <el-option ::label="$t('auto.List.3.b85811')" value="high" />
          <el-option ::label="$t('auto.List.4.bd9256')" value="medium" />
          <el-option ::label="$t('auto.List.5.c81e7d')" value="low" />
          <el-option ::label="$t('auto.List.6.60a04b')" value="none" />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          :range-separator="$t('auto.List.7.981cbe')"
          :start-placeholder="$t('auto.List.8.b44c0f')"
          :end-placeholder="$t('auto.List.9.1d468b')"
          value-format="YYYY-MM-DD"
        />
        <el-button @click="loadData">
          <el-icon><Search /></el-icon>
          {{ $t('查询') }}
        </el-button>
      </div>

      <el-table :data="filteredAudits" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="audit_name" ::label="$t('auto.List.10.3ff797')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="risk_level" ::label="$t('auto.List.2.4a708d')" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="getRiskType(row.risk_level)" size="small">
              {{ getRiskText(row.risk_level) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_checks" :label="$t('auto.List.11.b71f3e')" width="100" align="center" />
        <el-table-column prop="passed" :label="$t('auto.List.12.23c1f3')" width="90" align="center">
          <template #default="{ row }">
            <span class="pass-count">{{ row.passed || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="failed" :label="$t('auto.List.13.acd5cb')" width="90" align="center">
          <template #default="{ row }">
            <span class="fail-count">{{ row.failed || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="warnings" :label="$t('auto.List.14.900c70')" width="90" align="center">
          <template #default="{ row }">
            <span class="warn-count">{{ row.warnings || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="score" :label="$t('auto.List.15.108db1')" width="100" align="center">
          <template #default="{ row }">
            <span :class="getScoreClass(row.score)">{{ row.score ?? '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="started_at" :label="$t('auto.List.16.592c59')" min-width="170">
          <template #default="{ row }">
            {{ formatTime(row.started_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="finished_at" :label="$t('auto.List.17.7cbadb')" min-width="170">
          <template #default="{ row }">
            {{ formatTime(row.finished_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.18.2b6bc0')" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.List.30.f26225') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.19.79c058')" />
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

    
    <el-dialog v-model="detailVisible" :title="$t('auto.List.20.3ffbfd')" width="780px">
      <div v-if="detailRecord">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.21.9211b5')">{{ detailRecord.id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.22.366355')">{{ detailRecord.audit_name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.23.25eb6b')">
            <el-tag :type="getRiskType(detailRecord.risk_level)" size="small">
              {{ getRiskText(detailRecord.risk_level) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.15.108db1')">
            <span :class="getScoreClass(detailRecord.score)">{{ detailRecord.score ?? '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.24.85bfd1')">{{ detailRecord.total_checks || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.25.bf98ae')">
            <span class="pass-count">{{ detailRecord.passed || 0 }}</span> /
            <span class="fail-count">{{ detailRecord.failed || 0 }}</span> /
            <span class="warn-count">{{ detailRecord.warnings || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.16.592c59')">{{ formatTime(detailRecord.started_at) }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.17.7cbadb')">{{ formatTime(detailRecord.finished_at) }}</el-descriptions-item>
        </el-descriptions>

        <h3 class="section-title">{{ $t('auto.List.31.600adc') }}</h3>
        <el-table :data="detailItems" v-loading="detailLoading" max-height="360" border>
          <el-table-column prop="name" :label="$t('auto.List.11.b71f3e')" min-width="180" />
          <el-table-column prop="category" :label="$t('auto.List.26.383a6d')" width="120" />
          <el-table-column prop="level" :label="$t('auto.List.27.e78e4b')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getRiskType(row.level)" size="small">
                {{ getRiskText(row.level) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="result" :label="$t('auto.List.28.5ad7f5')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getResultType(row.result)" size="small">
                {{ getResultText(row.result) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="message" :label="$t('auto.List.29.f411d0')" min-width="220" show-overflow-tooltip />
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, VideoPlay, Search } from '@element-plus/icons-vue'
import {
  getSecurityAuditList,
  runSecurityAudit,
  getSecurityAuditDetail
} from '@/api/securityAudit.js'

const { t } = useI18n()

const loading = ref(false)
const running = ref(false)
const detailLoading = ref(false)
const audits = ref([])
const pagination = reactive({ page: 1, size: 20, total: 0 })
const searchKeyword = ref('')
const filterRisk = ref('')
const dateRange = ref([])

const detailVisible = ref(false)
const detailRecord = ref(null)
const detailItems = ref([])

const stats = computed(() => {
  const result = { high: 0, medium: 0, low: 0, none: 0 }
  audits.value.forEach(a => {
    if (a.risk_level === 'high') result.high++
    else if (a.risk_level === 'medium') result.medium++
    else if (a.risk_level === 'low') result.low++
    else result.none++
  })
  return result
})

const filteredAudits = computed(() => {
  let result = audits.value
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(a => a.audit_name?.toLowerCase().includes(kw))
  }
  if (filterRisk.value) result = result.filter(a => a.risk_level === filterRisk.value)
  return result
})

const getRiskText = (level) => {
  const map = { high: t('auto.List.32.4433e7'), medium: t('auto.List.33.8a8992'), low: t('auto.List.34.b56976'), none: t('auto.List.35.e6044a') }
  return map[level] || level || t('auto.List.36.1622dc')
}
const getRiskType = (level) => {
  const map = { high: 'danger', medium: 'warning', low: 'success', none: 'info' }
  return map[level] || 'info'
}
const getResultText = (result) => {
  const map = { pass: t('auto.List.12.23c1f3'), fail: t('auto.List.13.acd5cb'), warning: t('auto.List.14.900c70') }
  return map[result] || result || '-'
}
const getResultType = (result) => {
  const map = { pass: 'success', fail: 'danger', warning: 'warning' }
  return map[result] || 'info'
}
const getScoreClass = (score) => {
  if (score == null) return ''
  if (score >= 90) return 'score-good'
  if (score >= 70) return 'score-mid'
  return 'score-bad'
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
    const res = await getSecurityAuditList({
      page: pagination.page,
      page_size: pagination.size
    })
    const data = res || {}
    const list = data.list || data || []
    audits.value = list
    pagination.total = data.total || list.length || 0
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.37.ea1926')))
    audits.value = []
  } finally {
    loading.value = false
  }
}

const runAudit = async () => {
  running.value = true
  try {
    await runSecurityAudit({ audit_name: `manual_audit_${Date.now()}` })
    ElMessage.success(i18n.global.t(t('auto.List.38.d704a7')))
    await loadData()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.List.39.ede84e'))
  } finally {
    running.value = false
  }
}

const viewDetail = async (row) => {
  detailRecord.value = row
  detailItems.value = []
  detailVisible.value = true
  detailLoading.value = true
  try {
    const res = await getSecurityAuditDetail(row.id)
    const data = res.data || res
    if (data) {
      detailRecord.value = data
      detailItems.value = data.items || data.checks || data.findings || data.results || []
    }
  } catch (e) {
    ElMessage.warning(i18n.global.t(t('auto.List.40.d839b5')));
  } finally {
    detailLoading.value = false
  }
}

onMounted(() => loadData())
</script>

<style scoped lang="scss">
.audit-page { padding: 20px; }
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
    .stat-value.danger { color: #EF4444; }
    .stat-value.warning { color: #F59E0B; }
    .stat-value.success { color: #10B981; }
  }
}
.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
  flex-wrap: wrap;
}
.section-title { margin: 16px 0 12px; font-size: 16px; font-weight: 600; }
.pass-count { color: #10B981; font-weight: 600; }
.fail-count { color: #EF4444; font-weight: 600; }
.warn-count { color: #F59E0B; font-weight: 600; }
.score-good { color: #10B981; font-weight: 700; font-size: 16px; }
.score-mid { color: #F59E0B; font-weight: 700; font-size: 16px; }
.score-bad { color: #EF4444; font-weight: 700; font-size: 16px; }
</style>
