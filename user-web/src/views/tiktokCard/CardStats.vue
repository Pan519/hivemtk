<template>
  <div class="tiktok-card-stats-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>TikTok卡片统计 - {{ cardInfo.title || '加载中...' }}</span>
          <div class="header-actions">
            <el-button @click="goBack">{{ $t('返回') }}</el-button>
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
              <el-descriptions-item :label="$t('auto.CardStats.8830162c')">{{ cardInfo.id }}</el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.32c65d8d')">{{ cardInfo.title }}</el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.3fea7ca7')">
                <el-tag :type="cardInfo.is_active ? 'success' : 'danger'">
                  {{ cardInfo.is_active ? '激活' : '禁用' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.eca37cb0')">{{ cardInfo.created_at }}</el-descriptions-item>
              <el-descriptions-item :label="$t('auto.CardStats.3bdd08ad')" :span="2">{{ cardInfo.description }}</el-descriptions-item>
            </el-descriptions>
          </el-col>
        </el-row>
      </div>

      
      <div class="stats-overview">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-card class="stats-card">
              <div class="stats-item">
                <div class="stats-value">{{ formatNumber(cardStats.viewCount) }}</div>
                <div class="stats-label">{{ $t('auto.CardStats.7e6ff5a0') }}</div>
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
                  <span>{{ $t('auto.CardStats.8205ea14') }}</span>
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
                  <span>{{ $t('auto.CardStats.73378a1f') }}</span>
                </div>
              </template>
              <el-table :data="cardStats.recentActivity.filter(item => item.action === 'view')" style="width: 100%">
                <el-table-column prop="username" :label="$t('auto.CardStats.1fd02a90')" />
                <el-table-column prop="action" :label="$t('auto.CardStats.2b6bc0f2')">
                  <template #default="scope">
                    <el-tag :type="getActionType(scope.row.action)">
                      {{ getActionText(scope.row.action) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="ip_address" :label="$t('auto.CardStats.2dc9105c')" />
                <el-table-column prop="created_at" :label="$t('auto.CardStats.19fcb9eb')">
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
import { getTikTokCard, getTikTokCardStats } from '@/api/tiktokCard'

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
    const res = await getTikTokCard(cardId.value)
    cardInfo.value = res;
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.CardStats.eea2ec18')))
    console.error(error)
  }
};

const fetchCardStats = async () => {
  loading.value = true
  try {
    const params = {
      groupBy: groupBy.value
    }
    
    if (dateRange.value && dateRange.value.length === 2) {
      params.startDate = formatDate(dateRange.value[0])
      params.endDate = formatDate(dateRange.value[1])
    }
    
    const res = await getTikTokCardStats(cardId.value, params)
    Object.assign(cardStats, res);

    nextTick(() => {
      updateCharts()
    });
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.CardStats.063dc79a')))
    console.error(error)
  } finally {
    loading.value = false
  }
};

const updateCharts = () => {
  if (visitChartInstance) {
    visitChartInstance.dispose()
  }
  
  visitChartInstance = safeInit(visitTrendChart.value)
  const visitOption = {
    title: {
      text: t('auto.CardStats.771eefa9'),
      left: 'center'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: cardStats.dailyStats.map(item => item.date)
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: t('auto.CardStats.88bc4d7a'),
        type: 'line',
        data: cardStats.dailyStats.map(item => item.view),
        smooth: true
      }
    ]
  }
  visitChartInstance.setOption(visitOption)
  
  window.addEventListener('resize', () => {
    if (visitChartInstance) visitChartInstance.resize()
  });
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
  router.go(-1)
};

const formatNumber = (num) => {
  if (!num) return '0'
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
};

const formatDate = (date) => {
  if (!date) return ''
  if (typeof date === 'string') return date
  
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  
  return `${year}-${month}-${day}`
};

// el-tag 的 type 只接受 primary/success/info/warning/danger，传 '' 会触发 Vue prop 校验告警
// 并回落到 default='primary'。中性动作用 'info'，与 kuaishouCard/CardStats.vue 既有写法一致。
const getActionType = (action) => {
  const actionMap = {
    view: 'info'
  }
  return actionMap[action] || 'info'
};

const getActionText = (action) => {
  const actionMap = {
    view: t('auto.CardStats.9c5c5cdb')
  }
  return actionMap[action] || action
};

onMounted(() => {
  const endDate = new Date();
  const startDate = new Date()
  startDate.setDate(endDate.getDate() - 7)
  
  dateRange.value = [startDate, endDate]
  
  fetchCardInfo()
  fetchCardStats()
});
</script>

<style scoped>
.tiktok-card-stats-container {
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
}

.card-info-section {
  margin-bottom: 20px;
}

.stats-overview {
  margin-bottom: 20px;
}

.charts-section {
  margin-bottom: 20px;
}

.bottom-section {
  margin-bottom: 20px;
}

.stats-card {
  text-align: center;
}

.stats-item {
  padding: 10px 0;
}

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

.chart-container {
  height: 300px;
  width: 100%;
}
</style>