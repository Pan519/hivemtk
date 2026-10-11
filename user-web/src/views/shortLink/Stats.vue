<template>
  <div class="short-link-stats-container">
    
    <div class="page-header">
      <h2>{{ $t('短链统计') }}</h2>
      <div class="action-buttons">
        <el-button type="primary" @click="handleAllStats">
          <el-icon><DataAnalysis /></el-icon>
          {{ $t('总体统计') }}
        </el-button>
      </div>
    </div>

    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="search-form-content">
        <el-form-item ::label="$t('auto.Stats.e8fb33e2')">
          <el-input v-model="searchForm.short_code" ::placeholder="$t('auto.Stats.0d80133a')" clearable />
        </el-form-item>
        <el-form-item ::label="$t('auto.Stats.7ad5dc73')">
          <el-input v-model="searchForm.original_url" ::placeholder="$t('auto.Stats.a20b8828')" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ $t('搜索') }}
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshRight /></el-icon>
            {{ $t('重置') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    
    <el-table :data="shortLinkList" border style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="short_code" ::label="$t('auto.Stats.e8fb33e2')" width="120" />
      <el-table-column prop="original_url" ::label="$t('auto.Stats.7ad5dc73')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="title" ::label="$t('auto.Stats.319eda89')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="click_count" ::label="$t('auto.Stats.b5a7040e')" width="120" />
      <el-table-column prop="status_str" ::label="$t('auto.Stats.1b857887')" width="100">
        <template #default="scope">
          <el-tag :type="getStatusType(scope.row)">
            {{ scope.row.status_str }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.Stats.2b6bc0f2')" width="200" fixed="right">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleStats(scope.row)">
            <el-icon><DataAnalysis /></el-icon>
            统计
          </el-button>
          <el-button type="success" size="small" @click="handleShare(scope.row)">
            <el-icon><Share /></el-icon>
            分享
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <div class="pagination-container">
      <el-pagination
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    
    <el-dialog
      v-model="statsDialogVisible"
      :title="$t('auto.Stats.aa8158f3')"
      width="80%"
      top="5vh"
    >
      <div v-loading="statsLoading" class="stats-content">
        <div class="stats-summary">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ currentStats.total_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.f3d321e3') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ currentStats.today_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.bee53f7b') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ currentStats.yesterday_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.10ce7cc2') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ currentStats.avg_daily_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.27144593') }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div class="stats-charts">
          <el-row :gutter="20">
            <el-col :span="16">
              <el-card>
                <div class="chart-title">{{ $t('auto.Stats.b5b95b2b') }}</div>
                <div ref="trendChartRef" class="chart-container"></div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card>
                <div class="chart-title">{{ $t('auto.Stats.c76b9b27') }}</div>
                <div ref="deviceChartRef" class="chart-container"></div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-dialog>

    
    <el-dialog
      v-model="allStatsDialogVisible"
      :title="$t('auto.Stats.480c698d')"
      width="80%"
      top="5vh"
    >
      <div v-loading="allStatsLoading" class="stats-content">
        <div class="stats-summary">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ allStats.total_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.f3d321e3') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ allStats.today_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.bee53f7b') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ allStats.total_links || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.3bbe5c0a') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ allStats.active_links || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.Stats.1aee54f8') }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div class="stats-charts">
          <el-row :gutter="20">
            <el-col :span="16">
              <el-card>
                <div class="chart-title">{{ $t('auto.Stats.b5b95b2b') }}</div>
                <div ref="allTrendChartRef" class="chart-container"></div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card>
                <div class="chart-title">{{ $t('auto.Stats.c76b9b27') }}</div>
                <div ref="allDeviceChartRef" class="chart-container"></div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-dialog>

    
    <el-dialog
      v-model="shareDialogVisible"
      :title="$t('auto.Stats.7bae0c6f')"
      width="500px"
    >
      <div v-loading="shareLoading" class="share-content">
        <el-form label-width="80px">
          <el-form-item :label="$t('auto.Stats.74cf91f3')">
            <el-input v-model="shareData.short_url" readonly>
              <template #append>
                <el-button @click="copyToClipboard(shareData.short_url)">{{ $t('auto.Stats.79d3abe9') }}</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('auto.Stats.22b03c02')">
            <div class="qr-code-container">
              <img v-if="shareData.qr_code" :src="shareData.qr_code" :alt="$t('auto.Stats.22b03c02')" />
              <div v-else class="qr-code-placeholder">{{ $t('auto.Stats.3db02883') }}</div>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { DataAnalysis, Share, Search, RefreshRight } from '@element-plus/icons-vue'
import { shortLinkApi } from '@/api/shortLink'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'

const { t } = useI18n()

const loading = ref(false);
const statsLoading = ref(false)
const allStatsLoading = ref(false)
const shareLoading = ref(false)
const shortLinkList = ref([])
const statsDialogVisible = ref(false)
const allStatsDialogVisible = ref(false)
const shareDialogVisible = ref(false)
const currentStats = ref({})
const allStats = ref({})
const shareData = ref({})

const trendChartRef = ref(null);
const deviceChartRef = ref(null)
const allTrendChartRef = ref(null)
const allDeviceChartRef = ref(null)

let trendChart = null;
let deviceChart = null
let allTrendChart = null
let allDeviceChart = null

const searchForm = reactive({
  short_code: '',
  original_url: ''
});

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

onMounted(() => {
  fetchShortLinkList()
});

const fetchShortLinkList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      short_code: searchForm.short_code,
      original_url: searchForm.original_url
    }
    const res = await shortLinkApi.getList(params)
    shortLinkList.value = res.list || [];
    pagination.total = res.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Stats.049f09dd')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const getStatusType = (row) => {
  if (row.status === 2)
    return 'danger';
  if (row.expire_time && new Date(row.expire_time) < new Date())
    return 'warning';
  return 'success';
}

const handleSearch = () => {
  pagination.page = 1
  fetchShortLinkList()
}

const resetSearch = () => {
  searchForm.short_code = ''
  searchForm.original_url = ''
  pagination.page = 1
  fetchShortLinkList()
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchShortLinkList()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchShortLinkList()
}

const handleStats = async (row) => {
  statsDialogVisible.value = true
  statsLoading.value = true
  
  try {
    const endDate = new Date();
    const startDate = new Date()
    startDate.setDate(endDate.getDate() - 6)
    
    const params = {
      start_date: formatDate(startDate),
      end_date: formatDate(endDate)
    }
    
    const res = await shortLinkApi.getStats(row.id, params)
    currentStats.value = res;

    await nextTick();
    initCharts(res)
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Stats.063dc79a')))
    console.error(error)
  } finally {
    statsLoading.value = false
  }
}

const handleAllStats = async () => {
  allStatsDialogVisible.value = true
  allStatsLoading.value = true
  
  try {
    const endDate = new Date();
    const startDate = new Date()
    startDate.setDate(endDate.getDate() - 6)
    
    const params = {
      start_date: formatDate(startDate),
      end_date: formatDate(endDate)
    }
    
    const res = await shortLinkApi.getAllStats(params)
    allStats.value = res;

    await nextTick();
    initAllCharts(res)
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Stats.063dc79a')))
    console.error(error)
  } finally {
    allStatsLoading.value = false
  }
}

const handleShare = async (row) => {
  shareDialogVisible.value = true
  shareLoading.value = true
  
  try {
    const res = await shortLinkApi.share(row.id)
    shareData.value = res;
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Stats.0df764a2')))
    console.error(error)
  } finally {
    shareLoading.value = false
  }
}

const copyToClipboard = (text) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success(i18n.global.t(t('auto.Stats.20a49536')))
    }).catch(() => {
      ElMessage.error(i18n.global.t(t('auto.Stats.5154ae17')))
    })
  } else {
    const textarea = document.createElement('textarea');
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success(i18n.global.t(t('auto.Stats.20a49536')))
  }
}

const formatDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const initCharts = (data) => {
  if (trendChart)
    trendChart.dispose();
  if (deviceChart) deviceChart.dispose()
  
  trendChart = safeInit(trendChartRef.value);
  const trendOption = {
    title: {
      show: false
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: [t('auto.Stats.4461a6d2')]
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.daily_stats ? data.daily_stats.map(item => item.date) : []
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: t('auto.Stats.4461a6d2'),
        type: 'line',
        stack: 'Total',
        smooth: true,
        data: data.daily_stats ? data.daily_stats.map(item => item.count) : []
      }
    ]
  }
  trendChart.setOption(trendOption)
  
  deviceChart = safeInit(deviceChartRef.value);
  const deviceOption = {
    title: {
      show: false
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 10
    },
    series: [
      {
        name: t('auto.Stats.02f41e23'),
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data.device_stats ? Object.keys(data.device_stats).map(key => ({
          value: data.device_stats[key],
          name: key
        })) : []
      }
    ]
  }
  deviceChart.setOption(deviceOption)
  
  window.addEventListener('resize', () => {
    trendChart.resize()
    deviceChart.resize()
  });
}

const initAllCharts = (data) => {
  if (allTrendChart)
    allTrendChart.dispose();
  if (allDeviceChart) allDeviceChart.dispose()
  
  allTrendChart = safeInit(allTrendChartRef.value);
  const trendOption = {
    title: {
      show: false
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: [t('auto.Stats.4461a6d2')]
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.daily_stats ? data.daily_stats.map(item => item.date) : []
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: t('auto.Stats.4461a6d2'),
        type: 'line',
        stack: 'Total',
        smooth: true,
        data: data.daily_stats ? data.daily_stats.map(item => item.count) : []
      }
    ]
  }
  allTrendChart.setOption(trendOption)
  
  allDeviceChart = safeInit(allDeviceChartRef.value);
  const deviceOption = {
    title: {
      show: false
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 10
    },
    series: [
      {
        name: t('auto.Stats.02f41e23'),
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data.device_stats ? Object.keys(data.device_stats).map(key => ({
          value: data.device_stats[key],
          name: key
        })) : []
      }
    ]
  }
  allDeviceChart.setOption(deviceOption)
  
  window.addEventListener('resize', () => {
    allTrendChart.resize()
    allDeviceChart.resize()
  });
}
</script>

<style lang="scss" scoped>
.short-link-stats-container {
  padding: 20px;
  
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      font-size: 24px;
      color: #303133;
    }
    
    .action-buttons {
      display: flex;
      gap: 10px;
    }
  }
  
  .search-form {
    margin-bottom: 20px;
    padding: 15px;
    background-color: #f5f7fa;
    border-radius: 4px;
    
    .search-form-content {
      margin: 0;
    }
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
  
  .stats-content {
    .stats-summary {
      margin-bottom: 20px;
      
      .stats-card {
        text-align: center;
        
        .stats-item {
          padding: 10px 0;
          
          .stats-value {
            font-size: 24px;
            font-weight: bold;
            color: #4F46E5;
            margin-bottom: 5px;
          }
          
          .stats-label {
            font-size: 14px;
            color: #909399;
          }
        }
      }
    }
    
    .stats-charts {
      .chart-title {
        font-size: 16px;
        font-weight: bold;
        margin-bottom: 10px;
        text-align: center;
      }
      
      .chart-container {
        height: 300px;
      }
    }
  }
  
  .share-content {
    .qr-code-container {
      display: flex;
      justify-content: center;
      align-items: center;
      padding: 10px;
      
      img {
        width: 200px;
        height: 200px;
      }
      
      .qr-code-placeholder {
        width: 200px;
        height: 200px;
        display: flex;
        justify-content: center;
        align-items: center;
        background-color: #f5f7fa;
        color: #909399;
      }
    }
  }
}
</style>