<template>
  <div class="xianyu-card-stats-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('闲鱼卡片统计') }}</span>
          <div class="header-actions">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              :range-separator="$t('auto.Stats.1.981cbe')"
              :start-placeholder="$t('auto.Stats.2.b44c0f')"
              :end-placeholder="$t('auto.Stats.3.1d468b')"
              @change="handleDateChange"
            />
            <el-select v-model="groupBy" ::placeholder="$t('auto.Stats.4.bf8db0')" @change="handleGroupByChange">
              <el-option ::label="$t('auto.Stats.5.4eb3b0')" value="day" />
              <el-option ::label="$t('auto.Stats.6.787456')" value="week" />
              <el-option ::label="$t('auto.Stats.7.31019b')" value="month" />
            </el-select>
            <el-button type="primary" @click="refreshData">{{ $t('刷新') }}</el-button>
          </div>
        </div>
      </template>

      
      <div class="stats-overview">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ overallStats.total_cards || 0 }}</div>
                <div class="stats-label">{{ $t('auto.Stats.13.8135cb') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ overallStats.active_cards || 0 }}</div>
                <div class="stats-label">{{ $t('auto.Stats.14.56601a') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatNumber(overallStats.total_views) }}</div>
                <div class="stats-label">{{ $t('auto.Stats.15.7e6ff5') }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatNumber(overallStats.total_clicks) }}</div>
                <div class="stats-label">{{ $t('auto.Stats.16.d28bf8') }}</div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      
      <div class="charts-container">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-card>
              <div class="chart-title">{{ $t('auto.Stats.17.72995d') }}</div>
              <div ref="viewsChart" class="chart"></div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card>
              <div class="chart-title">{{ $t('auto.Stats.18.d0b945') }}</div>
              <div ref="clicksChart" class="chart"></div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      
      <div class="ranking-container">
        <el-card>
          <div class="card-title">{{ $t('auto.Stats.19.34e47f') }}</div>
          <el-table :data="cardRanking" style="width: 100%" v-loading="loading">
            <el-table-column prop="title" :label="$t('auto.Stats.8.3944f2')" show-overflow-tooltip />
            <el-table-column prop="view_count" :label="$t('auto.Stats.9.88bc4d')" width="100" />
            <el-table-column prop="click_count" :label="$t('auto.Stats.10.2c85f0')" width="100" />
            <el-table-column prop="click_rate" :label="$t('auto.Stats.11.e6867b')" width="100">
              <template #default="scope">
                {{ formatPercent(scope.row.click_rate) }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Stats.12.2b6bc0')" width="100">
              <template #default="scope">
                <el-button size="small" type="primary" @click="viewCardDetails(scope.row)">{{ $t('auto.Stats.20.f26225') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
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
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { getXianyuCardOverallStats } from '@/api/xianyuCard'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)

const dateRange = ref([]);
const groupBy = ref('day')

const overallStats = reactive({
  total_cards: 0,
  active_cards: 0,
  total_views: 0,
  total_clicks: 0,
  total_shares: 0
});

const chartData = reactive({
  views: [],
  clicks: [],
  dates: []
})

const cardRanking = ref([])

let viewsChartInstance = null;
let clicksChartInstance = null

const viewsChart = ref(null)
const clicksChart = ref(null)

const formatNumber = (num) => {
  if (!num) return '0'
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + t('auto.Stats.21.9d0320')
  }
  return num.toString()
};

const formatPercent = (rate) => {
  if (!rate) return '0%'
  return (rate * 100).toFixed(2) + '%'
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
      name: t('auto.Stats.9.88bc4d')
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
      name: t('auto.Stats.10.2c85f0')
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

const fetchStats = async () => {
  loading.value = true
  try {
    const params = {
      start_date: dateRange.value?.[0] || '',
      end_date: dateRange.value?.[1] || '',
      group_by: groupBy.value
    }
    
    const res = await getXianyuCardOverallStats(params)
    const data = res;

    Object.assign(overallStats, data.overall || {});

    chartData.dates = data.chart?.dates || [];
    chartData.views = data.chart?.views || []
    chartData.clicks = data.chart?.clicks || []

    cardRanking.value = data.ranking || [];

    if (viewsChartInstance && clicksChartInstance) {
      initCharts()
    }
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Stats.22.063dc7')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const viewCardDetails = (card) => {
  router.push(`/xianyu-card-stats/${card.id}`)
};

const handleDateChange = () => {
  fetchStats()
};

const handleGroupByChange = () => {
  fetchStats()
};

const refreshData = () => {
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
.xianyu-card-stats-container {
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

.chart-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 15px;
}

.chart {
  width: 100%;
  height: 300px;
}

.card-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 15px;
}

.ranking-container {
  margin-top: 20px;
}
</style>