<template>
  <div class="whatsapp-matrix">
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('auto.BulkMatrix.14.03d4ea') }}</div>
          <div class="stat-value">{{ stats.sent }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card success">
          <div class="stat-label">{{ $t('auto.BulkMatrix.15.77642d') }}</div>
          <div class="stat-value">{{ stats.delivered }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card primary">
          <div class="stat-label">{{ $t('auto.BulkMatrix.16.8ad629') }}</div>
          <div class="stat-value">{{ stats.read }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card danger">
          <div class="stat-label">{{ $t('auto.BulkMatrix.17.acd5cb') }}</div>
          <div class="stat-value">{{ stats.failed }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.BulkMatrix.18.aec4c3') }}</span>
          <el-button type="primary" @click="showPreview = true">{{ $t('auto.BulkMatrix.19.94ce5f') }}</el-button>
        </div>
      </template>
      <el-form :model="form" label-width="120px">
        <el-form-item :label="$t('auto.BulkMatrix.1.adac24')">
          <el-select v-model="form.templateId" filterable>
            <el-option v-for="t in templates" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.BulkMatrix.2.723c9d')">
          <el-select v-model="form.audienceType">
            <el-option :label="$t('auto.BulkMatrix.3.41ae5f')" value="all" />
            <el-option :label="$t('auto.BulkMatrix.4.46ee3c')" value="segment" />
            <el-option :label="$t('auto.BulkMatrix.5.6682c6')" value="clue" />
            <el-option :label="$t('auto.BulkMatrix.6.4d9223')" value="csv" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.audienceType === 'segment'" :label="$t('auto.BulkMatrix.7.8ce5b5')">
          <el-select v-model="form.segmentId" filterable>
            <el-option v-for="s in segments" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.BulkMatrix.8.5a1914')">
          <el-input-number v-model="form.ratePerMinute" :min="1" :max="100" />
          <span style="margin-left: 8px; color: #94A3B8">{{ $t('auto.BulkMatrix.20.1cf0da') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.BulkMatrix.9.a05577')">
          <el-input v-model="form.variablesJson" type="textarea" :rows="4" placeholder='{"customer.name": "张三", "order.id": "ORD-001"}' />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="sending" @click="send">{{ $t('auto.BulkMatrix.21.48ad8a') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="progress-card">
      <template #header><span>{{ $t('auto.BulkMatrix.22.558452') }}</span></template>
      <el-progress :percentage="progress" :status="progressStatus" />
      <el-table :data="progressList" max-height="200">
        <el-table-column prop="phone" :label="$t('auto.BulkMatrix.10.9f9d36')" width="160" />
        <el-table-column prop="status" :label="$t('auto.BulkMatrix.11.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="messageId" :label="$t('auto.BulkMatrix.12.795e00')" width="200" />
        <el-table-column prop="error" :label="$t('auto.BulkMatrix.13.7030ff')" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus'
import { http } from '@/utils/request'

const { t } = useI18n()

const stats = ref({ sent: 0, delivered: 0, read: 0, failed: 0 })
const form = ref({
  templateId: null,
  audienceType: 'all',
  segmentId: null,
  ratePerMinute: 20,
  variablesJson: '{}'
})
const templates = ref([])
const segments = ref([])
const showPreview = ref(false)
const sending = ref(false)
const progress = ref(0)
const progressList = ref([])
const progressStatus = ref('')

async function load() {
  // 这两个端点的 data 不是裸数组：/api/whatsapp/templates 返回 { data: [...] }，
  // /api/user-segments 返回 { list: [...], total }。http 拦截器已解一层信封
  // （utils/request.js:155 return data.data），这里必须再解一层，
  // 否则 v-for 遍历的是对象的属性值，option 的 label/value 全 undefined，
  // 模板与分群下拉实际渲染空白且选不中（实测 12 条 ElOption type check 告警）。
  const [tplRes, segRes] = await Promise.all([
    http.get('/api/whatsapp/templates'),
    http.get('/api/user-segments')
  ])
  templates.value = tplRes?.data || []
  segments.value = segRes?.list || []
}

const progressStats = computed(() => {
  return { progress: progress.value, status: progressStatus.value }
})

function statusType(status) {
  return {
    sent: 'info',
    delivered: 'success',
    read: 'primary',
    failed: 'danger'
  }[status] || ''
}

async function send() {
  if (!form.value.templateId) return ElMessage.warning(t('auto.BulkMatrix.23.9af823'))
  let variables;
  try { variables = JSON.parse(form.value.variablesJson) } catch (e) {
    return ElMessage.error(t('auto.BulkMatrix.24.4e94fa'))
  }
  sending.value = true
  progressStatus.value = ''
  try {
    const job = await http.post('/api/whatsapp/bulk-send', {
      template_id: form.value.templateId,
      audience: { type: form.value.audienceType, segment_id: form.value.segmentId },
      rate_per_minute: form.value.ratePerMinute,
      variables
    })
    ElMessage.success(t('auto.BulkMatrix.25.e2486b') + job.id)
    pollProgress(job.id)
  } finally {
    sending.value = false
  }
}

async function pollProgress(jobId) {
  const timer = setInterval(async () => {
    try {
      const p = await http.get(`/api/whatsapp/jobs/${jobId}/progress`)
      progress.value = p.percentage
      progressList.value = p.recent || []
      stats.value = p.stats || stats.value
      if (p.status === 'done' || p.status === 'failed') {
        clearInterval(timer)
        progressStatus.value = p.status === 'done' ? 'success' : 'exception'
        ElMessage[p.status === 'done' ? 'success' : 'error'](`任务${p.status === 'done' ? t('auto.BulkMatrix.26.769d88') : t('auto.BulkMatrix.17.acd5cb')}`)
      }
    } catch (_) {
      clearInterval(timer)
    }
  }, 2000)
}

onMounted(load)
</script>

<style scoped>
.whatsapp-matrix { padding: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; padding: 8px; }
.stat-label { color: #64748B; font-size: 12px; }
.stat-value { font-size: 32px; font-weight: 700; margin: 8px 0; }
.stat-card.success .stat-value { color: #10B981; }
.stat-card.primary .stat-value { color: #4F46E5; }
.stat-card.danger .stat-value { color: #EF4444; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.progress-card { margin-top: 16px; }
</style>
