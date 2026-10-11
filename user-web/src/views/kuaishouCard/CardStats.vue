<template>
  <div class="kuaishou-card-stats-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>快手卡片统计 - {{ cardInfo.title ?? '加载中...' }}</span>
          <div class="header-actions">
            <el-button @click="goBack">{{ $t('返回') }}</el-button>
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              :range-separator="$t('auto.CardStats.1.981cbe')"
              :start-placeholder="$t('auto.CardStats.2.b44c0f')"
              :end-placeholder="$t('auto.CardStats.3.1d468b')"
              @change="handleDateChange"
            />
            <el-select v-model="groupBy" ::placeholder="$t('auto.CardStats.4.bf8db0')" @change="handleGroupByChange">
              <el-option ::label="$t('auto.CardStats.5.4eb3b0')" value="day" />
              <el-option ::label="$t('auto.CardStats.6.787456')" value="week" />
              <el-option ::label="$t('auto.CardStats.7.31019b')" value="month" />
            </el-select>
            <el-button type="primary" @click="refreshData">{{ $t('刷新') }}</el-button>
          </div>
        </div>
      </template>

      
      <div class="card-info-section" v-if="cardInfo.id">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-image 
              style="width: 100%; height: 200px" 
              :src="cardInfo.image_url" 
              :preview-src-list="[cardInfo.image_url]"
              fit="cover"
            />
          </el-col>
          <el-col :span="18">
            <el-descriptions :column="2" border>
              <el-descriptions-item :label="$t('auto.CardStats.8.883016')">{{ cardInfo.id }}</el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.9.32c65d')">{{ cardInfo.title }}</el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.10.3fea7c')">
                <el-tag :type="cardInfo.is_active ? 'success' : 'danger'">
                  {{ cardInfo.is_active ? '激活' : '禁用' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.11.eca37c')">{{ cardInfo.created_at }}</el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.12.3bdd08')" :span="2">{{ cardInfo.description }}</el-descriptions-item>
            </el-descriptions>
          </el-col>
        </el-row>
      </div>

      
      <div class="stats-overview">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatNumber(cardStats.view_count) }}</div>
                <div class="stats-label">{{ $t('auto.CardStats.17.7e6ff5') }}</div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      
      <div class="charts-section">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-card>
              <template #header>
                <div class="card-header">
                  <span>{{ $t('auto.CardStats.18.8205ea') }}</span>
                </div>
              </template>
              <div ref="visitTrendChart" class="chart-container"></div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      
      <div class="bottom-section">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-card>
              <template #header>
                <div class="card-header">
                  <span>{{ $t('auto.CardStats.19.73378a') }}</span>
                </div>
              </template>
              <el-table :data="cardStats.recentActivity.filter(item => item.action === 'view')" style="width: 100%">
                <el-table-column prop="username" :label="$t('auto.CardStats.13.1fd02a')" />
                <el-table-column prop="action" :label="$t('auto.CardStats.14.2b6bc0')">
                  <template #default="scope">
                    <el-tag :type="getActionType(scope.row.action)">
                      {{ getActionText(scope.row.action) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="ip_address" :label="$t('auto.CardStats.15.2dc910')" />
                <el-table-column prop="created_at" :label="$t('auto.CardStats.16.19fcb9')">
                  <template #default="scope">
                    {{ formatDate(scope.row.created_at) }}
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { getKuaishouCard, getKuaishouCardStats } from '@/api/kuaishouCard'

const route = useRoute();
const router = useRouter()

const { t } = useI18n()

const loading = ref(false);
const cardId = ref(route.params.id)
const cardInfo = ref({})
const dateRange = ref([])
const groupBy = ref('day')

const cardStats = reactive({
  cardId: 0,
  title: '',
  viewCount: 0,
  dailyStats: [],
  recentActivity: []
});

const visitTrendChart = ref(null);
let visitChartInstance = null

const fetchCardInfo = async () => {
  try {
    const res = await getKuaishouCard(cardId.value)
    cardInfo.value = res;
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.CardStats.20.eea2ec')))
    console.error(t('auto.CardStats.21.a4ec28'), error)
  }
};

const fetchCardStats = async () => {
  try {
    const params = {}
    
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = formatDate(dateRange.value[0])
      params.end_date = formatDate(dateRange.value[1])
    }
    
    params.group_by = groupBy.value;
    
    const res = await getKuaishouCardStats(cardId.value, params)
    Object.assign(cardStats, res);
    if (!Array.isArray(cardStats.dailyStats))
      cardStats.dailyStats = [];
    if (!Array.isArray(cardStats.recentActivity)) cardStats.recentActivity = []

    nextTick(() => {
      initCharts()
    });
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.CardStats.22.063dc7')))
    console.error(t('auto.CardStats.23.f43cb0'), error)
  }
};

const initCharts = () => {
  if (visitChartInstance) {
    visitChartInstance.dispose()
  }

  if (visitTrendChart.value) {
    visitChartInstance = safeInit(visitTrendChart.value)
    const visitOption = {
      tooltip: {
        trigger: 'axis'
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
        data: cardStats.dailyStats.map(item => item.date)
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: t('auto.CardStats.24.88bc4d'),
          type: 'line',
          smooth: true,
          data: cardStats.dailyStats.map(item => item.view_count),
          itemStyle: {
            color: '#2684FF'
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(38, 132, 255, 0.3)' },
              { offset: 1, color: 'rgba(38, 132, 255, 0.1)' }
            ])
          }
        }
      ]
    }
    visitChartInstance.setOption(visitOption)
  }

  window.addEventListener('resize', handleResize);
};

const handleResize = () => {
  if (visitChartInstance) visitChartInstance.resize()
};

const handleDateChange = () => {
  fetchCardStats()
};

const handleGroupByChange = () => {
  fetchCardStats()
};

const refreshData = () => {
  fetchCardInfo()
  fetchCardStats()
};

const goBack = () => {
  router.push('/kuaishouCard')
};

const formatNumber = (num) => {
  if (!num) return '0'
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
};

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
};

// 'default' 不在 el-tag 的合法 type 里（primary/success/info/warning/danger），会触发 Vue prop
// 校验告警并回落到 default='primary'。未识别动作用中性的 'info'。
const getActionType = (action) => {
  const typeMap = {
    'view': 'info'
  }
  return typeMap[action] || 'info'
};

const getActionText = (action) => {
  const typeMap = {
    'view': t('auto.CardStats.25.9c5c5c')
  }
  return typeMap[action] || t('auto.CardStats.26.1622dc')
};

onMounted(async () => {
  const endDate = new Date();
  const startDate = new Date()
  startDate.setDate(startDate.getDate() - 30)
  dateRange.value = [startDate, endDate]
  
  await fetchCardInfo();
  await fetchCardStats()
});

import { onUnmounted } from 'vue';
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (visitChartInstance) visitChartInstance.dispose()
})
</script>

<style scoped>
.kuaishou-card-stats-container {
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

.card-info-section {
  margin-bottom: 30px;
}

.stats-overview {
  margin-bottom: 30px;
}

.stats-card {
  text-align: center;
}

.stats-item {
  padding: 20px;
}

.stats-value {
  font-size: 24px;
  font-weight: 600;
  color: #4F46E5;
  margin-bottom: 8px;
}

.stats-label {
  font-size: 14px;
  color: #909399;
}

.charts-section {
  margin-bottom: 30px;
}

.chart-container {
  height: 400px;
}

.bottom-section {
  margin-bottom: 30px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.el-tag {
  margin-right: 5px;
}

.el-tag:last-child {
  margin-right: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .header-actions {
    width: 100%;
    flex-wrap: wrap;
  }
}
</style>