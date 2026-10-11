<template>
  <div class="knowledge-statistics">
    
    <el-card class="filter-card">
      <div class="filter-bar">
        <el-select v-model="filter.product_id" ::placeholder="$t('auto.KnowledgeStatistics.1.c63099')" clearable style="width: 220px" @change="loadAll">
          <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
        </el-select>
        <el-select v-model="filter.days" style="width: 120px" @change="loadAll">
          <el-option :label="`近 ${filter.days} 天`" :value="filter.days" v-for="d in [7, 30, 90]" :key="d" />
        </el-select>
        <el-button :icon="Refresh" @click="loadAll">{{ $t('刷新') }}</el-button>
      </div>
    </el-card>

    
    <el-row :gutter="16" class="overview-row">
      <el-col :span="6">
        <el-card class="metric-card">
          <div class="metric-icon" style="background: #ecf5ff; color: #4F46E5">
            <el-icon><Document /></el-icon>
          </div>
          <div class="metric-info">
            <div class="metric-value">{{ overview.total_documents || 0 }}</div>
            <div class="metric-label">{{ $t('文档总数') }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="metric-card">
          <div class="metric-icon" style="background: #f0f9eb; color: #10B981">
            <el-icon><Files /></el-icon>
          </div>
          <div class="metric-info">
            <div class="metric-value">{{ overview.total_chunks || 0 }}</div>
            <div class="metric-label">{{ $t('分段总数') }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="metric-card">
          <div class="metric-icon" style="background: #fdf6ec; color: #F59E0B">
            <el-icon><Coin /></el-icon>
          </div>
          <div class="metric-info">
            <div class="metric-value">{{ formatNumber(overview.total_tokens) }}</div>
            <div class="metric-label">{{ $t('auto.KnowledgeStatistics.20.b615f5') }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="metric-card">
          <div class="metric-icon" style="background: #fef0f0; color: #EF4444">
            <el-icon><Search /></el-icon>
          </div>
          <div class="metric-info">
            <div class="metric-value">{{ overview.total_searches || 0 }}</div>
            <div class="metric-label">{{ $t('总检索次数') }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>{{ $t('索引健康度') }}</span>
          </template>
          <div class="health-block">
            <div class="health-item">
              <div class="health-num success">{{ overview.index_health?.indexed_docs || 0 }}</div>
              <div class="health-label">{{ $t('auto.KnowledgeStatistics.21.b0d887') }}</div>
            </div>
            <div class="health-item">
              <div class="health-num warning">{{ overview.index_health?.processing_docs || 0 }}</div>
              <div class="health-label">{{ $t('auto.KnowledgeStatistics.22.5d459d') }}</div>
            </div>
            <div class="health-item">
              <div class="health-num info">{{ overview.index_health?.pending_docs || 0 }}</div>
              <div class="health-label">{{ $t('auto.KnowledgeStatistics.23.047109') }}</div>
            </div>
            <div class="health-item">
              <div class="health-num danger">{{ overview.index_health?.failed_docs || 0 }}</div>
              <div class="health-label">{{ $t('auto.KnowledgeStatistics.24.acd5cb') }}</div>
            </div>
          </div>
          <el-progress :percentage="Math.round((overview.index_health?.index_rate || 0) * 100)" :stroke-width="14" status="success" />
          <div class="health-rate">就绪率 {{ formatPercent(overview.index_health?.index_rate) }}</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>{{ $t('auto.KnowledgeStatistics.25.15332d') }}</span>
          </template>
          <div v-for="(count, type) in overview.source_type_breakdown" :key="type" class="source-bar">
            <div class="source-label">{{ sourceTypeLabel(type) }}</div>
            <el-progress :percentage="calcPercent(count, overview.total_documents)" :format="() => count" :color="sourceTypeColor(type)" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card class="chart-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.KnowledgeStatistics.26.9c7434') }}</span>
          <el-radio-group v-model="docStatsTab" size="small">
            <el-radio-button label="trend">{{ $t('auto.KnowledgeStatistics.27.444247') }}</el-radio-button>
            <el-radio-button label="source">{{ $t('auto.KnowledgeStatistics.28.26a888') }}</el-radio-button>
            <el-radio-button label="category">{{ $t('auto.KnowledgeStatistics.29.cefbd4') }}</el-radio-button>
            <el-radio-button label="top">{{ $t('auto.KnowledgeStatistics.30.de2cb4') }}</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <div v-if="docStatsTab === 'trend'" class="trend-chart">
        <div v-for="item in documentStats.import_trend || []" :key="item.day" class="trend-bar">
          <div class="trend-day">{{ item.day.substring(5) }}</div>
          <div class="trend-stack">
            <div class="trend-success" :style="{ height: calcHeight(item.count, maxTrend) + 'px' }" :title="`成功 ${item.count - item.failed}`">
              <span v-if="item.count - item.failed > 0">{{ item.count - item.failed }}</span>
            </div>
            <div class="trend-failed" :style="{ height: calcHeight(item.failed, maxTrend) + 'px' }" :title="`失败 ${item.failed}`">
              <span v-if="item.failed > 0">{{ item.failed }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="docStatsTab === 'source'" class="pie-list">
        <div v-for="item in documentStats.source_type_pie || []" :key="item.type" class="pie-row">
          <span class="pie-label">{{ sourceTypeLabel(item.type) }}</span>
          <el-progress :percentage="calcPercent(item.count, overview.total_documents)" :color="sourceTypeColor(item.type)" :format="() => item.count" />
        </div>
      </div>

      <div v-else-if="docStatsTab === 'category'" class="pie-list">
        <div v-for="item in (documentStats.category_pie || []).slice(0, 10)" :key="item.category" class="pie-row">
          <span class="pie-label">{{ item.category }}</span>
          <el-progress :percentage="calcPercent(item.count, overview.total_documents)" :format="() => item.count" />
        </div>
      </div>

      <div v-else-if="docStatsTab === 'top'">
        <el-empty :description="$t('auto.KnowledgeStatistics.2.cbdf35')" />
      </div>
    </el-card>

    
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>{{ $t('auto.KnowledgeStatistics.31.27fcb5') }}</span>
          </template>
          <div class="trend-chart">
            <div v-for="item in searchStats.search_trend || []" :key="item.day" class="trend-bar">
              <div class="trend-day">{{ item.day.substring(5) }}</div>
              <div class="trend-stack">
                <div class="trend-success" :style="{ height: calcHeight(item.count, maxSearchTrend) + 'px' }">
                  <span v-if="item.count > 0">{{ item.count }}</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>{{ $t('auto.KnowledgeStatistics.32.1ae1d7') }}</span>
          </template>
          <div class="histogram">
            <div v-for="bucket in searchStats.score_histogram || []" :key="bucket.range" class="hist-bar">
              <div class="hist-value">{{ bucket.count }}</div>
              <div class="hist-rect" :style="{ height: calcHeight(bucket.count, maxScore) + 'px' }"></div>
              <div class="hist-label">{{ bucket.range }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-row :gutter="16" class="chart-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>{{ $t('auto.KnowledgeStatistics.33.ad40db') }}</span>
          </template>
          <el-table :data="searchStats.hot_queries || []" size="small" max-height="400">
            <el-table-column type="index" label="#" width="50" />
            <el-table-column prop="query" :label="$t('auto.KnowledgeStatistics.3.650e7c')" show-overflow-tooltip />
            <el-table-column prop="count" :label="$t('auto.KnowledgeStatistics.4.f965fe')" width="80" />
            <el-table-column :label="$t('auto.KnowledgeStatistics.5.678735')" width="100">
              <template #default="{ row }">{{ formatPercent(row.hit_rate) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>{{ $t('auto.KnowledgeStatistics.34.e15268') }}</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.6.38a402')">{{ searchStats.quality_stats?.total_searches || 0 }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.7.2f0aac')">{{ searchStats.quality_stats?.hit_count || 0 }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.8.208785')">{{ searchStats.quality_stats?.no_result_count || 0 }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.5.678735')">{{ formatPercent(searchStats.quality_stats?.hit_rate) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.9.356987')">{{ (searchStats.quality_stats?.avg_score || 0).toFixed(3) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.10.da21e8')">{{ (searchStats.quality_stats?.max_score || 0).toFixed(3) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.11.e31d50')">{{ (searchStats.quality_stats?.min_score || 0).toFixed(3) }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KnowledgeStatistics.12.0d652d')">{{ (searchStats.quality_stats?.avg_latency_ms || 0).toFixed(0) }} ms</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card class="chart-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.KnowledgeStatistics.35.c4129f') }}</span>
          <div>
            <span class="header-stat">{{ $t('auto.KnowledgeStatistics.36.050d87') }}<el-tag :type="importStats.success_rate >= 0.9 ? 'success' : 'warning'">{{ formatPercent(importStats.success_rate) }}</el-tag></span>
            <span class="header-stat">{{ $t('auto.KnowledgeStatistics.37.b278e1') }}<el-tag>{{ (importStats.avg_duration_ms || 0).toFixed(0) }} ms</el-tag></span>
          </div>
        </div>
      </template>
      <el-table :data="importStats.recent_logs || []" size="small" max-height="400">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="source_type" :label="$t('auto.KnowledgeStatistics.13.26ca20')" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ sourceTypeLabel(row.source_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="batch_no" :label="$t('auto.KnowledgeStatistics.14.51767c')" width="160" show-overflow-tooltip />
        <el-table-column prop="operator" :label="$t('auto.KnowledgeStatistics.15.f9ac4b')" width="100" />
        <el-table-column prop="status" :label="$t('auto.KnowledgeStatistics.16.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="getPassFailStatusTagType(row.status)" size="small">
              {{ getPassFailStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration_ms" :label="$t('auto.KnowledgeStatistics.17.7f62c0')" width="100" />
        <el-table-column prop="error_detail" :label="$t('auto.KnowledgeStatistics.18.4604d5')" show-overflow-tooltip>
          <template #default="{ row }">
            <el-text v-if="row.error_detail" type="danger" size="small">{{ row.error_detail }}</el-text>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.KnowledgeStatistics.19.19fcb9')" width="170">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card class="chart-card">
      <template #header>
        <span>{{ $t('auto.KnowledgeStatistics.38.d735a1') }}</span>
      </template>
      <el-row :gutter="16" class="openapi-stats">
        <el-col :span="6">
          <div class="openapi-stat-item">
            <div class="openapi-num">{{ openapiStats.total_sources || 0 }}</div>
            <div class="openapi-label">{{ $t('auto.KnowledgeStatistics.39.e67b65') }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="openapi-stat-item">
            <div class="openapi-num success">{{ openapiStats.enabled_sources || 0 }}</div>
            <div class="openapi-label">{{ $t('auto.KnowledgeStatistics.40.53ace4') }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="openapi-stat-item">
            <div class="openapi-num danger">{{ openapiStats.failed_sources || 0 }}</div>
            <div class="openapi-label">{{ $t('auto.KnowledgeStatistics.41.d61036') }}</div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="openapi-stat-item">
            <div class="openapi-num primary">{{ formatNumber(openapiStats.total_synced) }}</div>
            <div class="openapi-label">{{ $t('auto.KnowledgeStatistics.42.a53731') }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Document, Files, Coin, Search, Refresh } from '@element-plus/icons-vue'
import { knowledgeAPI } from '@/api/knowledge'
import { ragProductConfigAPI } from '@/api/ragProductConfig'
import { getSourceLabel } from '@/constants/source';
import { PASS_FAIL_STATUS, getStatusLabel, getStatusTagType } from '@/constants/status';

const getPassFailStatusLabel = (s) => getStatusLabel(s, PASS_FAIL_STATUS);
const getPassFailStatusTagType = (s) => getStatusTagType(s, PASS_FAIL_STATUS)

const { t } = useI18n()

const productList = ref([])
const overview = ref({})
const documentStats = ref({})
const searchStats = ref({})
const importStats = ref({})
const openapiStats = ref({})

const filter = reactive({
  product_id: '',
  days: 30
})

const docStatsTab = ref('trend')

const loadAll = async () => {
  await Promise.all([loadProducts(), loadOverview(), loadDocumentStats(), loadSearchStats(), loadImportStats(), loadOpenAPIStats()])
}

const loadProducts = async () => {
  try {
    const res = await ragProductConfigAPI.listProducts()
    if (Array.isArray(res)) {
      productList.value = res
    } else if (Array.isArray(res?.items)) {
      productList.value = res.items
    }
  } catch (e) {
    console.error(t('auto.KnowledgeStatistics.43.4c46f3'), e)
  }
}

const loadOverview = async () => {
  try {
    const res = await knowledgeAPI.getOverviewStats({ product_id: filter.product_id })
    if (res) overview.value = res
  } catch (e) {
    console.error(t('auto.KnowledgeStatistics.44.3af525'), e)
  }
}

const loadDocumentStats = async () => {
  try {
    const res = await knowledgeAPI.getDocumentStats({ product_id: filter.product_id, days: filter.days })
    if (res) documentStats.value = res
  } catch (e) {
    console.error(t('auto.KnowledgeStatistics.45.854510'), e)
  }
}

const loadSearchStats = async () => {
  try {
    const res = await knowledgeAPI.getSearchStats({ product_id: filter.product_id, days: filter.days })
    if (res) searchStats.value = res
  } catch (e) {
    console.error(t('auto.KnowledgeStatistics.46.6f1114'), e)
  }
}

const loadImportStats = async () => {
  try {
    const res = await knowledgeAPI.getImportStats({ product_id: filter.product_id, days: filter.days })
    if (res) importStats.value = res
  } catch (e) {
    console.error(t('auto.KnowledgeStatistics.47.7a1dec'), e)
  }
}

const loadOpenAPIStats = async () => {
  try {
    const res = await knowledgeAPI.getOpenAPIStats({ product_id: filter.product_id })
    if (res) openapiStats.value = res
  } catch (e) {
    console.error(t('auto.KnowledgeStatistics.48.4b2c67'), e)
  }
}

const maxTrend = computed(() => {
  const list = documentStats.value.import_trend || []
  return Math.max(1, ...list.map(i => i.count || 0))
})

const maxSearchTrend = computed(() => {
  const list = searchStats.value.search_trend || []
  return Math.max(1, ...list.map(i => i.count || 0))
})

const maxScore = computed(() => {
  const list = searchStats.value.score_histogram || []
  return Math.max(1, ...list.map(i => i.count || 0))
})

const SOURCE_COLOR_MAP = {
  upload: '#4F46E5', text: '#10B981', url: '#F59E0B', openapi: '#EF4444', batch: '#909399'
};
const sourceTypeLabel = (t) => getSourceLabel(t)
const sourceTypeColor = (t) => SOURCE_COLOR_MAP[t] || '#4F46E5'

const calcPercent = (val, total) => {
  if (!total || total <= 0) return 0
  return Math.round((val / total) * 100)
}

const calcHeight = (val, max) => {
  if (!max || max <= 0) return 0
  return Math.max(2, Math.round((val / max) * 100))
}

const formatNumber = (n) => n == null ? '-' : Number(n).toLocaleString()
const formatPercent = (n) => n == null ? '-' : (n * 100).toFixed(1) + '%'
const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(() => {
  loadAll()
})
</script>

<style scoped lang="scss">
.knowledge-statistics {
  .filter-card {
    margin-bottom: 16px;
  }
  .filter-bar {
    display: flex;
    gap: 12px;
    align-items: center;
  }
  .overview-row, .chart-row {
    margin-bottom: 16px;
  }
  .chart-card {
    margin-bottom: 16px;
  }
  .metric-card {
    display: flex;
    align-items: center;
    padding: 16px;
    .metric-icon {
      width: 48px;
      height: 48px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 22px;
      margin-right: 12px;
    }
    .metric-info {
      .metric-value {
        font-size: 22px;
        font-weight: 600;
        color: #303133;
      }
      .metric-label {
        font-size: 13px;
        color: #909399;
        margin-top: 4px;
      }
    }
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .header-stat {
      margin-left: 16px;
      font-size: 13px;
      color: #606266;
    }
  }
  .health-block {
    display: flex;
    justify-content: space-around;
    margin-bottom: 16px;
    .health-item {
      text-align: center;
      .health-num {
        font-size: 26px;
        font-weight: 600;
        &.success { color: #10B981; }
        &.warning { color: #F59E0B; }
        &.info { color: #909399; }
        &.danger { color: #EF4444; }
      }
      .health-label {
        font-size: 12px;
        color: #909399;
        margin-top: 4px;
      }
    }
  }
  .health-rate {
    text-align: center;
    margin-top: 8px;
    font-size: 13px;
    color: #606266;
  }
  .source-bar, .pie-row {
    margin-bottom: 12px;
    .source-label, .pie-label {
      font-size: 13px;
      color: #303133;
      margin-bottom: 4px;
      display: block;
    }
  }
  .trend-chart {
    display: flex;
    align-items: flex-end;
    height: 160px;
    gap: 4px;
    padding: 0 8px;
    .trend-bar {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 4px;
      .trend-day {
        font-size: 10px;
        color: #909399;
        writing-mode: vertical-rl;
        height: 16px;
      }
      .trend-stack {
        width: 100%;
        max-width: 24px;
        display: flex;
        flex-direction: column;
        justify-content: flex-end;
        .trend-success {
          background: #10B981;
          border-radius: 2px 2px 0 0;
          display: flex;
          align-items: flex-end;
          justify-content: center;
          color: #fff;
          font-size: 10px;
          min-height: 2px;
          padding-bottom: 2px;
        }
        .trend-failed {
          background: #EF4444;
          border-radius: 2px 2px 0 0;
          display: flex;
          align-items: flex-end;
          justify-content: center;
          color: #fff;
          font-size: 10px;
          min-height: 0;
        }
      }
    }
  }
  .histogram {
    display: flex;
    align-items: flex-end;
    justify-content: space-around;
    height: 160px;
    .hist-bar {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 4px;
      .hist-value {
        font-size: 12px;
        color: #303133;
      }
      .hist-rect {
        background: linear-gradient(180deg, #4F46E5 0%, #79bbff 100%);
        width: 36px;
        border-radius: 4px 4px 0 0;
      }
      .hist-label {
        font-size: 11px;
        color: #909399;
      }
    }
  }
  .openapi-stats {
    text-align: center;
    .openapi-stat-item {
      padding: 16px 0;
      .openapi-num {
        font-size: 28px;
        font-weight: 600;
        color: #303133;
        &.success { color: #10B981; }
        &.danger { color: #EF4444; }
        &.primary { color: #4F46E5; }
      }
      .openapi-label {
        font-size: 13px;
        color: #909399;
        margin-top: 4px;
      }
    }
  }
}
</style>
