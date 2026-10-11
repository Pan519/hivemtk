<template>
  <div class="xianyu-card-detail-stats">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('闲鱼卡片详情统计') }}</span>
          <div class="header-actions">
            <el-button @click="goBack" icon="ArrowLeft">{{ $t('返回') }}</el-button>
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              :range-separator="$t('auto.CardStats.981cbe31')"
              :start-placeholder="$t('auto.CardStats.b44c0f33')"
              :end-placeholder="$t('auto.CardStats.1d468be9')"
              @change="handleDateChange"
            />
            <el-select v-model="groupBy" ::placeholder="$t('auto.CardStats.bf8db0fe')" @change="handleGroupByChange">
              <el-option ::label="$t('auto.CardStats.4eb3b00e')" value="day" />
              <el-option ::label="$t('auto.CardStats.7874566d')" value="week" />
              <el-option ::label="$t('auto.CardStats.31019b33')" value="month" />
            </el-select>
            <el-button type="primary" @click="refreshData">{{ $t('刷新') }}</el-button>
          </div>
        </div>
      </template>

      
      <div class="card-info" v-if="cardInfo">
        <el-descriptions :column="3" border>
          <el-descriptions-item :label="$t('auto.CardStats.3944f217')">{{ cardInfo.title }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.CardStats.3bdd08ad')">{{ cardInfo.description }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.CardStats.3fea7ca7')">
            <el-tag :type="cardInfo.is_active ? 'success' : 'danger'">
              {{ cardInfo.is_active ? '激活' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.CardStats.eca37cb0')">{{ cardInfo.created_at }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.CardStats.7e6ff5a0')">{{ cardInfo.view_count || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.CardStats.d28bf874')">{{ cardInfo.click_count || 0 }}</el-descriptions-item>
        </el-descriptions>
      </div>

      
      <div class="stats-overview">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatNumber(statsData.total_views) }}</div>
                <div class="stats-label">{{ $t('auto.CardStats.88bc4d7a') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatNumber(statsData.total_clicks) }}</div>
                <div class="stats-label">{{ $t('auto.CardStats.2c85f0c8') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatNumber(statsData.total_shares) }}</div>
                <div class="stats-label">{{ $t('auto.CardStats.34688e27') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatPercent(statsData.click_rate) }}</div>
                <div class="stats-label">{{ $t('auto.CardStats.e6867bb9') }}</div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      
      <div class="charts-container">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-card>
              <div class="chart-title">{{ $t('auto.CardStats.72995d12') }}</div>
              <div ref="viewsChart" class="chart"></div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card>
              <div class="chart-title">{{ $t('auto.CardStats.d0b94533') }}</div>
              <div ref="clicksChart" class="chart"></div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      
      <div class="detail-table">
        <el-card>
          <div class="table-title">{{ $t('auto.CardStats.331e29b7') }}</div>
          <el-table :data="detailData" style="width: 100%" v-loading="loading">
            <el-table-column prop="date" :label="$t('auto.CardStats.4ff1e74e')" width="120" />
            <el-table-column prop="views" :label="$t('auto.CardStats.88bc4d7a')" width="100" />
            <el-table-column prop="clicks" :label="$t('auto.CardStats.2c85f0c8')" width="100" />
            <el-table-column prop="shares" :label="$t('auto.CardStats.34688e27')" width="100" />
            <el-table-column prop="click_rate" :label="$t('auto.CardStats.e6867bb9')" width="100">
              <template #default="scope">
                {{ formatPercent(scope.row.click_rate) }}
              </template>
            </el-table-column>
            <el-table-column prop="unique_visitors" :label="$t('auto.CardStats.cf0848d0')" width="100" />
          </el-table>
          
          
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.page_size"
              :page-sizes="[10, 20, 50, 100]"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useRouter, useRoute } from 'vue-router'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { getXianyuCard, getXianyuCardStats } from '@/api/xianyuCard'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const loading = ref(false)

const cardId = ref(route.params.id);

const dateRange = ref([]);
const groupBy = ref('day')

const cardInfo = ref(null);

const statsData = reactive({
  total_views: 0,
  total_clicks: 0,
  total_shares: 0,
  click_rate: 0
});

const chartData = reactive({
  views: [],
  clicks: [],
  dates: []
})

const detailData = ref([])

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
});

let viewsChartInstance = null;
let clicksChartInstance = null

const viewsChart = ref(null)
const clicksChart = ref(null)

const formatNumber = (num) => {
  if (!num) return '0'
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + t('auto.CardStats.9d032066')
  }
  return num.toString()
};

const formatPercent = (rate) => {
  if (!rate) return '0%'
  return (rate * 100).toFixed(2) + '%'
};

const goBack = () => {
  router.back()
};

const initCharts = () => {
  if (!viewsChart.value || !clicksChart.value)
    return;
  viewsChartInstance = safeInit(viewsChart.value);
  const viewsOption = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: chartData.dates
    },
    yAxis: {
      type: 'value',
      name: t('auto.CardStats.88bc4d7a')
    },
    series: [{
      data: chartData.views,
      type: 'line',
      smooth: true,
      areaStyle: {
        opacity: 0.3
      },
      itemStyle: {
        color: '#FF6B35'
      }
    }]
  }
  viewsChartInstance.setOption(viewsOption)

  clicksChartInstance = safeInit(clicksChart.value);
  const clicksOption = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: chartData.dates
    },
    yAxis: {
      type: 'value',
      name: t('auto.CardStats.2c85f0c8')
    },
    series: [{
      data: chartData.clicks,
      type: 'line',
      smooth: true,
      areaStyle: {
        opacity: 0.3
      },
      itemStyle: {
        color: '#1890ff'
      }
    }]
  }
  clicksChartInstance.setOption(clicksOption)
};

const fetchCardInfo = async () => {
  try {
    const res = await getXianyuCard(cardId.value)
    cardInfo.value = res;
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.CardStats.eea2ec18')))
    console.error(error)
  }
};

const fetchStats = async () => {
  loading.value = true
  try {
    const params = {
      start_date: dateRange.value?.[0] || '',
      end_date: dateRange.value?.[1] || '',
      group_by: groupBy.value,
      page: pagination.page,
      page_size: pagination.page_size
    }
    
    const res = await getXianyuCardStats(cardId.value, params)
    Object.assign(statsData, res.stats || {});

    chartData.dates = res.chart?.dates || [];
    chartData.views = res.chart?.views || []
    chartData.clicks = res.chart?.clicks || []

    detailData.value = res.details?.list || [];
    pagination.total = res.details?.total || 0

    if (viewsChartInstance && clicksChartInstance) {
      initCharts()
    }
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.CardStats.063dc79a')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleDateChange = () => {
  pagination.page = 1
  fetchStats()
};

const handleGroupByChange = () => {
  pagination.page = 1
  fetchStats()
};

const refreshData = () => {
  fetchStats()
};

const handleSizeChange = (size) => {
  pagination.page_size = size
  fetchStats()
};

const handleCurrentChange = (page) => {
  pagination.page = page
  fetchStats()
};

const handleResize = () => {
  if (viewsChartInstance) {
    viewsChartInstance.resize()
  }
  if (clicksChartInstance) {
    clicksChartInstance.resize()
  }
};

onMounted(() => {
  const endDate = new Date();
  const startDate = new Date()
  startDate.setDate(startDate.getDate() - 30)
  dateRange.value = [startDate.toISOString().split('T')[0], endDate.toISOString().split('T')[0]]
  
  fetchCardInfo()
  fetchStats()
  
  setTimeout(() => {
    initCharts()
  }, 100);
  
  window.addEventListener('resize', handleResize);
})

onUnmounted(() => {
  if (viewsChartInstance) {
    viewsChartInstance.dispose()
  }
  if (clicksChartInstance) {
    clicksChartInstance.dispose()
  }
  
  window.removeEventListener('resize', handleResize);
})
</script>

<style scoped>
.xianyu-card-detail-stats {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.card-info {
  margin-bottom: 20px;
}

.stats-overview {
  margin-bottom: 20px;
}

.stats-card {
  text-align: center;
}

.stats-item {
  padding: 20px;
}

.stats-value {
  font-size: 32px;
  font-weight: bold;
  color: #FF6B35;
  margin-bottom: 10px;
}

.stats-label {
  font-size: 14px;
  color: #666;
}

.charts-container {
  margin-bottom: 20px;
}

.chart-title,
.table-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 15px;
}

.chart {
  width: 100%;
  height: 300px;
}

.detail-table {
  margin-top: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>