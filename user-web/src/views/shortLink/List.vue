<template>
  <div class="short-link-container">
    
    <div class="page-header">
      <h2>{{ $t('短链管理') }}</h2>
      <div class="action-buttons">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          {{ $t('添加短链') }}
        </el-button>
      </div>
    </div>

    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="search-form-content">
        <el-form-item ::label="$t('auto.List.1.e8fb33')">
          <el-input v-model="searchForm.short_code" ::placeholder="$t('auto.List.2.0d8013')" clearable />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.3.7ad5dc')">
          <el-input v-model="searchForm.original_url" ::placeholder="$t('auto.List.4.a20b88')" clearable />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.5.1b8578')">
          <el-select v-model="searchForm.status" ::placeholder="$t('auto.List.6.333291')" clearable>
            <el-option ::label="$t('auto.List.7.62f11e')" value="1" />
            <el-option ::label="$t('auto.List.8.27ea4f')" value="2" />
          </el-select>
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

    
    <PageState v-if="error" state="error" :error-text="error" @retry="fetchShortLinkList" />
    <el-table v-if="!error" :data="shortLinkList" border style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="short_code" ::label="$t('auto.List.1.e8fb33')" width="120" />
      <el-table-column prop="original_url" ::label="$t('auto.List.3.7ad5dc')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="title" ::label="$t('auto.List.9.319eda')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="description" ::label="$t('auto.List.10.d1ef7e')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="click_count" ::label="$t('auto.List.11.20d518')" width="100" />
      <el-table-column prop="status_str" ::label="$t('auto.List.5.1b8578')" width="100">
        <template #default="scope">
          <el-tag :type="getStatusType(scope.row)">
            {{ scope.row.status_str }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="expire_time" :label="$t('auto.List.12.1fa23f')" width="180">
        <template #default="scope">
          {{ scope.row.expire_time || '永不过期' }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="$t('auto.List.13.eca37c')" width="180" />
      <el-table-column :label="$t('auto.List.14.2b6bc0')" width="300" fixed="right">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">
            编辑
          </el-button>
          <el-button type="info" size="small" @click="handleStats(scope.row)">
            <el-icon><DataAnalysis /></el-icon>
            统计
          </el-button>
          <el-button type="success" size="small" @click="handleShare(scope.row)">
            <el-icon><Share /></el-icon>
            分享
          </el-button>
          <el-button type="danger" size="small" @click="handleDelete(scope.row)">
            删除
          </el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty :description="$t('auto.List.15.21efd8')" />
      </template>
    </el-table>

    
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item :label="$t('auto.List.16.2fdc27')" prop="short_code">
          <el-input v-model="form.short_code" :placeholder="$t('auto.List.17.c5bfc0')">
            <template #append>
              <el-button @click="generateShortCode" :loading="generating">{{ $t('auto.List.35.4dfe70') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('auto.List.18.5e3436')" prop="original_url">
          <el-input v-model="form.original_url" :placeholder="$t('auto.List.19.c2c347')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.20.32c65d')" prop="title">
          <el-input v-model="form.title" :placeholder="$t('auto.List.21.96641a')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.22.3bdd08')" prop="description">
          <el-input v-model="form.description" type="textarea" :placeholder="$t('auto.List.23.11956a')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.24.300e01')" prop="password">
          <el-input v-model="form.password" type="password" :placeholder="$t('auto.List.25.b2340c')" show-password />
        </el-form-item>
        <el-form-item :label="$t('auto.List.12.1fa23f')" prop="expire_time">
          <el-date-picker
            v-model="form.expire_time"
            type="datetime"
            :placeholder="$t('auto.List.26.71342e')"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.27.3fea7c')" prop="status" v-if="isEdit">
          <el-select v-model="form.status" :placeholder="$t('auto.List.28.e1c965')">
            <el-option :label="$t('auto.List.29.fd6e80')" :value="1" />
            <el-option :label="$t('auto.List.30.710ad0')" :value="2" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.36.625fb2') }}</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="statsDialogVisible"
      :title="$t('auto.List.31.aa8158')"
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
                  <div class="stats-label">{{ $t('auto.List.37.f3d321') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ currentStats.today_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.List.38.bee53f') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ currentStats.yesterday_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.List.39.10ce7c') }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stats-card">
                <div class="stats-item">
                  <div class="stats-value">{{ currentStats.avg_daily_clicks || 0 }}</div>
                  <div class="stats-label">{{ $t('auto.List.40.271445') }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div class="stats-charts">
          <el-row :gutter="20">
            <el-col :span="16">
              <el-card>
                <div class="chart-title">{{ $t('auto.List.41.b5b95b') }}</div>
                <div ref="trendChartRef" class="chart-container"></div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card>
                <div class="chart-title">{{ $t('auto.List.42.c76b9b') }}</div>
                <div ref="deviceChartRef" class="chart-container"></div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-dialog>

    
    <el-dialog
      v-model="shareDialogVisible"
      :title="$t('auto.List.32.7bae0c')"
      width="500px"
    >
      <div v-loading="shareLoading" class="share-content">
        <el-form label-width="80px">
          <el-form-item :label="$t('auto.List.33.74cf91')">
            <el-input v-model="shareData.short_url" readonly>
              <template #append>
                <el-button @click="copyToClipboard(shareData.short_url)">{{ $t('auto.List.43.79d3ab') }}</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('auto.List.34.22b03c')">
            <div class="qr-code-container">
              <img v-if="shareData.qr_code" :src="shareData.qr_code" :alt="$t('auto.List.34.22b03c')" />
              <div v-else class="qr-code-placeholder">{{ $t('auto.List.44.3db028') }}</div>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, RefreshRight, DataAnalysis, Share } from '@element-plus/icons-vue'
import { shortLinkApi } from '@/api/shortLink'
import PageState from '@/components/PageState.vue'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'

const { t } = useI18n()

const loading = ref(false);
const error = ref('')
const generating = ref(false)
const statsLoading = ref(false)
const shareLoading = ref(false)
const shortLinkList = ref([])
const dialogVisible = ref(false)
const statsDialogVisible = ref(false)
const shareDialogVisible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const currentStats = ref({})
const shareData = ref({})

const trendChartRef = ref(null);
const deviceChartRef = ref(null)

let trendChart = null;
let deviceChart = null

const searchForm = reactive({
  short_code: '',
  original_url: '',
  status: ''
});

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const form = reactive({
  id: null,
  short_code: '',
  original_url: '',
  title: '',
  description: '',
  password: '',
  expire_time: null,
  status: 1
});

const rules = {
  short_code: [
    { required: true, message: i18n.global.t(t('auto.List.17.c5bfc0')), trigger: 'blur' },
    { min: 4, max: 20, message: i18n.global.t(t('auto.List.45.f5a9a7')), trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9]+$/, message: i18n.global.t(t('auto.List.46.af8fdf')), trigger: 'blur' }
  ],
  original_url: [
    { required: true, message: i18n.global.t(t('auto.List.19.c2c347')), trigger: 'blur' },
    { type: 'url', message: i18n.global.t(t('auto.List.47.a7f6ca')), trigger: 'blur' }
  ]
};

const dialogTitle = computed(() => {
  return isEdit.value ? t('auto.List.48.f344cb') : t('auto.List.49.b2abd5')
});

onMounted(() => {
  fetchShortLinkList()
});

const fetchShortLinkList = async () => {
  loading.value = true
  error.value = ''
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      short_code: searchForm.short_code,
      original_url: searchForm.original_url,
      status: searchForm.status ? parseInt(searchForm.status) : undefined
    }
    const res = await shortLinkApi.getList(params)
    if (res && res.list) {
      shortLinkList.value = res.list
      pagination.total = res.total
    } else {
      error.value = i18n.global.t(t('auto.List.50.049f09'))
    }
  } catch (error) {
    error.value = i18n.global.t(t('auto.List.50.049f09'))
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
  searchForm.status = ''
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

const handleAdd = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  form.id = row.id
  form.short_code = row.short_code
  form.original_url = row.original_url
  form.title = row.title
  form.description = row.description
  form.password = row.password
  form.expire_time = row.expire_time
  form.status = row.status
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除短链 "${row.short_code}" 吗？`,
    t('auto.List.51.02d981'),
    {
      confirmButtonText: t('auto.List.52.38cf16'),
      cancelButtonText: t('auto.List.36.625fb2'),
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await shortLinkApi.delete(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.53.0007d1')));
      fetchShortLinkList()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.54.acf066')))
      console.error(error)
    }
  }).catch((e) => {
    if (e !== 'cancel' && e !== 'close')
      throw e;
  })
}

const generateShortCode = async () => {
  generating.value = true
  try {
    const res = await shortLinkApi.generateShortCode({ length: 6 })
    form.short_code = res.short_code;
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.55.055b5a')))
    console.error(error)
  } finally {
    generating.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        let res
        if (isEdit.value) {
          res = await shortLinkApi.update(form)
        } else {
          res = await shortLinkApi.create(form)
        }

        ElMessage.success(isEdit.value ? t('auto.List.56.55aa63') : t('auto.List.57.04a691'));
        dialogVisible.value = false
        fetchShortLinkList()
      } catch (error) {
        ElMessage.error(isEdit.value ? t('auto.List.58.930442') : t('auto.List.59.a88928'))
        console.error(error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDialogClose = () => {
  resetForm()
}

const resetForm = () => {
  form.id = null
  form.short_code = ''
  form.original_url = ''
  form.title = ''
  form.description = ''
  form.password = ''
  form.expire_time = null
  form.status = 1
  if (formRef.value) {
    formRef.value.resetFields()
  }
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
    ElMessage.error(i18n.global.t(t('auto.List.60.063dc7')))
    console.error(error)
  } finally {
    statsLoading.value = false
  }
}

const handleShare = async (row) => {
  shareDialogVisible.value = true
  shareLoading.value = true
  
  try {
    const res = await shortLinkApi.share(row.id)
    shareData.value = res;
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.61.0df764')))
    console.error(error)
  } finally {
    shareLoading.value = false
  }
}

const copyToClipboard = (text) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success(i18n.global.t(t('auto.List.62.20a495')))
    }).catch(() => {
      ElMessage.error(i18n.global.t(t('auto.List.63.5154ae')))
    })
  } else {
    const textarea = document.createElement('textarea');
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success(i18n.global.t(t('auto.List.62.20a495')))
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
      data: [t('auto.List.64.4461a6')]
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
        name: t('auto.List.64.4461a6'),
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
        name: t('auto.List.65.02f41e'),
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
</script>

<style lang="scss" scoped>
.short-link-container {
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