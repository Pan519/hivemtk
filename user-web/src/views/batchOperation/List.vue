<template>
  <div class="batch-operation-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('批量操作') }}</h2>
        <p class="subtitle">{{ $t('批量处理线索、客户、消息等数据') }}</p>
      </div>
    </el-card>

    <el-row :gutter="20">
      <el-col :span="8" v-for="tool in tools" :key="tool.name">
        <el-card class="tool-card" shadow="hover" @click="openTool(tool)">
          <div class="tool-icon">
            <el-icon :size="40" :color="tool.color"><component :is="tool.icon" /></el-icon>
          </div>
          <h3>{{ tool.name }}</h3>
          <p>{{ tool.description }}</p>
          <div class="tool-stats">
            <span>已执行: {{ tool.usedCount }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="history-card">
      <template #header>
        <span>{{ $t('操作历史') }}</span>
      </template>
      <el-table :data="histories" v-loading="loading" stripe>
        <el-table-column prop="tool" :label="$t('auto.List.1.20dce2')" width="150" />
        <el-table-column prop="target" :label="$t('auto.List.2.73e825')" min-width="200" />
        <el-table-column prop="totalCount" :label="$t('auto.List.3.9ed7d3')" width="100" />
        <el-table-column prop="successCount" :label="$t('auto.List.4.330363')" width="100" />
        <el-table-column prop="failedCount" :label="$t('auto.List.5.acd5cb')" width="100" />
        <el-table-column prop="status" :label="$t('auto.List.6.3fea7c')" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="$t('auto.List.7.70b363')" width="180" />
        <el-table-column :label="$t('auto.List.8.2b6bc0')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.List.34.f26225') }}</el-button>
            <el-button link type="danger" v-if="row.status === 'running'" @click="cancelTask(row)">{{ $t('auto.List.35.625fb2') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="currentTool?.name" width="800px">
      <el-form :model="form" label-width="120px">
        <el-form-item :label="$t('auto.List.9.de9cc3')">
          <el-input :value="currentTool?.name" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.List.10.a094e5')">
          <el-radio-group v-model="form.source">
            <el-radio label="all">{{ $t('auto.List.36.a8b0c2') }}</el-radio>
            <el-radio label="filter">{{ $t('auto.List.37.c2fe62') }}</el-radio>
            <el-radio label="import">{{ $t('auto.List.38.8d9a07') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.List.11.5148fc')" v-if="form.source === 'filter'">
          <el-input v-model="form.filter" type="textarea" :rows="3" placeholder='{"status": "active", "created_at": ">2024-01-01"}' />
        </el-form-item>
        <el-form-item :label="$t('auto.List.12.7499f4')" v-if="form.source === 'import'">
          <el-upload action="/api/batch/import" :show-file-list="false">
            <el-button>{{ $t('auto.List.39.fd7e0c') }}</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item :label="$t('auto.List.13.716c8a')">
          <el-input v-model="form.params" type="textarea" :rows="4" placeholder='{"action": "send_message", "template": "welcome"}' />
        </el-form-item>
        <el-form-item :label="$t('auto.List.14.874469')">
          <el-switch v-model="form.dryRun" />
          <span class="form-tip">{{ $t('auto.List.40.8e80c2') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.35.625fb2') }}</el-button>
        <el-button @click="preview">{{ $t('auto.List.41.645dbc') }}</el-button>
        <el-button type="primary" @click="execute">{{ $t('auto.List.42.1a6aa2') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="previewVisible" :title="$t('auto.List.15.c66881')" width="900px" :close-on-click-modal="false">
      <div v-loading="previewLoading">
        <div class="preview-summary" v-if="!previewLoading">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-statistic :title="$t('auto.List.16.625daf')" :value="previewSummary.total" />
            </el-col>
            <el-col :span="8">
              <el-statistic :title="$t('auto.List.17.6f4fa8')" :value="previewSummary.success">
                <template #suffix>
                  <el-icon color="#10B981"><component :is="'Check'" /></el-icon>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic :title="$t('auto.List.18.18fcd5')" :value="previewSummary.failed">
                <template #suffix>
                  <el-icon color="#EF4444"><component :is="'Close'" /></el-icon>
                </template>
              </el-statistic>
            </el-col>
          </el-row>
        </div>

        <div class="preview-errors" v-if="previewSummary.errors.length > 0">
          <h4>{{ $t('auto.List.43.18e60c') }}</h4>
          <el-alert
            v-for="(err, idx) in previewSummary.errors"
            :key="idx"
            :title="err"
            type="warning"
            show-icon
            :closable="false"
            style="margin-bottom: 8px"
          />
        </div>

        <el-table :data="previewData" style="margin-top: 20px" v-if="previewData.length > 0">
          <el-table-column type="index" :label="$t('auto.List.19.faaadc')" width="60" />
          <el-table-column prop="target" :label="$t('auto.List.2.73e825')" min-width="200" />
          <el-table-column prop="action" :label="$t('auto.List.8.2b6bc0')" width="150" />
          <el-table-column prop="status" :label="$t('auto.List.20.8136c5')" width="120">
            <template #default="{ row }">
              <el-tag :type="getPreviewStatusTagType(row.status)" size="small">
                {{ getPreviewStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="message" :label="$t('auto.List.21.2432b5')" show-overflow-tooltip min-width="200" />
        </el-table>

        <el-empty v-if="!previewLoading && previewData.length === 0" :description="$t('auto.List.22.65fa2e')" />
      </div>

      <template #footer>
        <el-button @click="previewVisible = false">{{ $t('auto.List.44.b15d91') }}</el-button>
        <el-button type="primary" @click="execute">{{ $t('auto.List.45.d575bf') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailDialogVisible" :title="$t('auto.List.23.5288a9')" width="800px">
      <div v-loading="detailLoading">
        <template v-if="taskDetail">
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('auto.List.24.3a3778')">{{ taskDetail.id }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.25.ac6116')">{{ taskDetail.tool }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.2.73e825')">{{ taskDetail.target}}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.6.3fea7c')">
              <el-tag :type="getStatusType(taskDetail.status)">{{ getStatusText(taskDetail.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.7.70b363')">{{ taskDetail.createdAt }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.26.7cbadb')">{{ taskDetail.finishedAt}}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.27.8887e3')">{{ taskDetail.totalCount }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.28.7586c6')">
              <span style="color: #10B981">{{ taskDetail.successCount }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.29.21a0bb')">
              <span style="color: #EF4444">{{ taskDetail.failedCount }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.30.fea070')">{{ taskDetail.skippedCount || 0 }}</el-descriptions-item>
          </el-descriptions>

          <div v-if="taskDetail.errorMessage" style="margin-top: 20px">
            <h4>{{ $t('auto.List.46.4604d5') }}</h4>
            <el-alert
              :title="taskDetail.errorMessage"
              type="error"
              show-icon
              :closable="false"
            />
          </div>

          <div v-if="taskDetail.errorDetails && taskDetail.errorDetails.length > 0" style="margin-top: 20px">
            <h4>{{ $t('auto.List.47.71f623') }}</h4>
            <el-table :data="taskDetail.errorDetails" stripe max-height="300">
              <el-table-column type="index" :label="$t('auto.List.19.faaadc')" width="60" />
              <el-table-column prop="target" :label="$t('auto.List.2.73e825')" min-width="150" />
              <el-table-column prop="error" :label="$t('auto.List.31.2fc597')" min-width="250" show-overflow-tooltip />
              <el-table-column prop="time" :label="$t('auto.List.32.19fcb9')" width="170" />
            </el-table>
          </div>
        </template>
        <el-empty v-if="!detailLoading && !taskDetail" :description="$t('auto.List.33.e10f76')" />
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">{{ $t('auto.List.44.b15d91') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, User, ChatLineRound, Promotion, Money, Document, Bell, ShoppingCart } from '@element-plus/icons-vue'
import { getBatchTools, runBatch, getBatchHistories, cancelBatch, previewBatch, getBatchDetail } from '@/api/batchOperation.js'
import { toList } from '@/utils/list.js'
import { PASS_FAIL_STATUS, getStatusLabel, getStatusTagType } from '@/constants/status';

const getPreviewStatusLabel = (s) => getStatusLabel(s, PASS_FAIL_STATUS) === t('auto.List.4.330363') ? t('auto.List.48.596149') : (getStatusLabel(s, PASS_FAIL_STATUS) || s || '-');
const getPreviewStatusTagType = (s) => getStatusTagType(s, PASS_FAIL_STATUS)

const { t } = useI18n()

const loading = ref(false)
const tools = ref([])
const histories = ref([])
const dialogVisible = ref(false)
const currentTool = ref(null)
const form = ref({ source: 'all', filter: '', params: '', dryRun: true })

const detailDialogVisible = ref(false)
const detailLoading = ref(false)
const taskDetail = ref(null)

const previewVisible = ref(false);
const previewLoading = ref(false)
const previewData = ref([])
const previewSummary = ref({ total: 0, success: 0, failed: 0, errors: [] })

const getStatusType = (s) => ({ running: 'warning', success: 'success', failed: 'danger', partial: 'warning' }[s])
const getStatusText = (s) => ({ running: t('auto.List.49.46e386'), success: t('auto.List.50.769d88'), failed: t('auto.List.5.acd5cb'), partial: t('auto.List.51.983b68') }[s] || s)

const loadTools = async () => {
  const res = await getBatchTools()
  tools.value = res || [
    { name: t('auto.List.52.3182c7'), description: t('auto.List.53.17fb37'), color: '#4F46E5', icon: 'ChatLineRound', usedCount: 0 },
    { name: t('auto.List.54.be8c0a'), description: 'Excel/CSV 批量导入', color: '#10B981', icon: 'Document', usedCount: 0 },
    { name: t('auto.List.55.52140d'), description: t('auto.List.56.05e713'), color: '#F59E0B', icon: 'Promotion', usedCount: 0 },
    { name: t('auto.List.57.b2879b'), description: t('auto.List.58.9ceef4'), color: '#EF4444', icon: 'User', usedCount: 0 }
  ]
}

const loadHistories = async () => {
  loading.value = true
  try {
    const res = await getBatchHistories()
    histories.value = toList(res)
  } finally {
    loading.value = false
  }
}

const openTool = (tool) => {
  currentTool.value = tool
  form.value = { source: 'all', filter: '', params: '', dryRun: true }
  dialogVisible.value = true
}

const preview = async () => {
  previewLoading.value = true
  previewVisible.value = true
  previewData.value = []
  previewSummary.value = { total: 0, success: 0, failed: 0, errors: [] }

  try {
    const res = await previewBatch({ tool: currentTool.value.name, ...form.value, dryRun: true })
    const data = res;

    previewData.value = data.items || data.records || [];
    previewSummary.value = {
      total: data.total || data.totalCount || previewData.value.length,
      success: data.successCount || data.success || 0,
      failed: data.failedCount || data.failed || 0,
      errors: data.errors || []
    }
  } catch (error) {
    ElMessage.error(t('auto.List.59.943ad7') + (error?.message || error))
    previewVisible.value = false
  } finally {
    previewLoading.value = false
  }
}

const execute = async () => {
  const loading = ElMessage({ message: i18n.global.t(t('auto.List.60.c380d7')), type: 'info', duration: 0 })
  try {
    await runBatch({ tool: currentTool.value.name, ...form.value })
    loading.close()
    ElMessage.success(i18n.global.t(t('auto.List.61.b29706')))
    dialogVisible.value = false
    previewVisible.value = false
    loadHistories()
  } catch (error) {
    loading.close()
    ElMessage.error(i18n.global.t(t('auto.List.62.1c83d7')))
  }
}

const viewDetail = async (row) => {
  detailDialogVisible.value = true
  detailLoading.value = true
  try {
    const res = await getBatchDetail(row.id)
    taskDetail.value = res || row
  } catch {
    taskDetail.value = row
  } finally {
    detailLoading.value = false
  }
}

const cancelTask = async (row) => {
  try {
    await ElMessageBox.confirm(t('auto.List.63.1fed38'), t('auto.List.64.e83a25'), { type: 'warning' })
    await cancelBatch(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.65.2111cc')))
    loadHistories()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.66.482e85')))
  }
}

onMounted(() => {
  loadTools()
  loadHistories()
})
</script>

<style scoped lang="scss">
.batch-operation-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.tool-card {
  cursor: pointer;
  text-align: center;
  margin-bottom: 20px;
  transition: transform 0.3s;
  &:hover { transform: translateY(-5px); }
  .tool-icon { margin-bottom: 15px; }
  h3 { margin: 0 0 10px 0; }
  p { color: #909399; font-size: 13px; }
  .tool-stats { color: #4F46E5; font-size: 12px; margin-top: 10px; }
}
.history-card { margin-top: 20px; }
.form-tip { color: #909399; font-size: 12px; margin-left: 10px; }
</style>
