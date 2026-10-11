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
        <el-table-column prop="tool" :label="$t('auto.List.20dce2c6')" width="150" />
        <el-table-column prop="target" :label="$t('auto.List.73e82552')" min-width="200" />
        <el-table-column prop="totalCount" :label="$t('auto.List.9ed7d3ad')" width="100" />
        <el-table-column prop="successCount" :label="$t('auto.List.330363df')" width="100" />
        <el-table-column prop="failedCount" :label="$t('auto.List.acd5cb84')" width="100" />
        <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="$t('auto.List.70b3635a')" width="180" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
            <el-button link type="danger" v-if="row.status === 'running'" @click="cancelTask(row)">{{ $t('auto.List.625fb26b') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="currentTool?.name" width="800px">
      <el-form :model="form" label-width="120px">
        <el-form-item :label="$t('auto.List.de9cc3dd')">
          <el-input :value="currentTool?.name" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.List.a094e5b7')">
          <el-radio-group v-model="form.source">
            <el-radio label="all">{{ $t('auto.List.a8b0c204') }}</el-radio>
            <el-radio label="filter">{{ $t('auto.List.c2fe6253') }}</el-radio>
            <el-radio label="import">{{ $t('auto.List.8d9a071e') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.List.5148fcc0')" v-if="form.source === 'filter'">
          <el-input v-model="form.filter" type="textarea" :rows="3" placeholder='{"status": "active", "created_at": ">2024-01-01"}' />
        </el-form-item>
        <el-form-item :label="$t('auto.List.7499f450')" v-if="form.source === 'import'">
          <el-upload action="/api/batch/import" :show-file-list="false">
            <el-button>{{ $t('auto.List.fd7e0c99') }}</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item :label="$t('auto.List.716c8aba')">
          <el-input v-model="form.params" type="textarea" :rows="4" placeholder='{"action": "send_message", "template": "welcome"}' />
        </el-form-item>
        <el-form-item :label="$t('auto.List.874469c9')">
          <el-switch v-model="form.dryRun" />
          <span class="form-tip">{{ $t('auto.List.8e80c218') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button @click="preview">{{ $t('auto.List.645dbc55') }}</el-button>
        <el-button type="primary" @click="execute">{{ $t('auto.List.1a6aa24e') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="previewVisible" :title="$t('auto.List.c66881bc')" width="900px" :close-on-click-modal="false">
      <div v-loading="previewLoading">
        <div class="preview-summary" v-if="!previewLoading">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-statistic :title="$t('auto.List.625dafe9')" :value="previewSummary.total" />
            </el-col>
            <el-col :span="8">
              <el-statistic :title="$t('auto.List.6f4fa8ce')" :value="previewSummary.success">
                <template #suffix>
                  <el-icon color="#10B981"><component :is="'Check'" /></el-icon>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic :title="$t('auto.List.18fcd54b')" :value="previewSummary.failed">
                <template #suffix>
                  <el-icon color="#EF4444"><component :is="'Close'" /></el-icon>
                </template>
              </el-statistic>
            </el-col>
          </el-row>
        </div>

        <div class="preview-errors" v-if="previewSummary.errors.length > 0">
          <h4>{{ $t('auto.List.18e60ca5') }}</h4>
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
          <el-table-column type="index" :label="$t('auto.List.faaadc44')" width="60" />
          <el-table-column prop="target" :label="$t('auto.List.73e82552')" min-width="200" />
          <el-table-column prop="action" :label="$t('auto.List.2b6bc0f2')" width="150" />
          <el-table-column prop="status" :label="$t('auto.List.8136c57c')" width="120">
            <template #default="{ row }">
              <el-tag :type="getPreviewStatusTagType(row.status)" size="small">
                {{ getPreviewStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="message" :label="$t('auto.List.2432b575')" show-overflow-tooltip min-width="200" />
        </el-table>

        <el-empty v-if="!previewLoading && previewData.length === 0" :description="$t('auto.List.65fa2e1a')" />
      </div>

      <template #footer>
        <el-button @click="previewVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
        <el-button type="primary" @click="execute">{{ $t('auto.List.d575bfd2') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailDialogVisible" :title="$t('auto.List.5288a978')" width="800px">
      <div v-loading="detailLoading">
        <template v-if="taskDetail">
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('auto.List.3a3778f2')">{{ taskDetail.id }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.ac6116b2')">{{ taskDetail.tool }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.73e82552')">{{ taskDetail.target}}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.3fea7ca7')">
              <el-tag :type="getStatusType(taskDetail.status)">{{ getStatusText(taskDetail.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.70b3635a')">{{ taskDetail.createdAt }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.7cbadbbe')">{{ taskDetail.finishedAt}}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.8887e33b')">{{ taskDetail.totalCount }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.7586c67e')">
              <span style="color: #10B981">{{ taskDetail.successCount }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.21a0bbc2')">
              <span style="color: #EF4444">{{ taskDetail.failedCount }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.fea07017')">{{ taskDetail.skippedCount || 0 }}</el-descriptions-item>
          </el-descriptions>

          <div v-if="taskDetail.errorMessage" style="margin-top: 20px">
            <h4>{{ $t('auto.List.4604d502') }}</h4>
            <el-alert
              :title="taskDetail.errorMessage"
              type="error"
              show-icon
              :closable="false"
            />
          </div>

          <div v-if="taskDetail.errorDetails && taskDetail.errorDetails.length > 0" style="margin-top: 20px">
            <h4>{{ $t('auto.List.71f62364') }}</h4>
            <el-table :data="taskDetail.errorDetails" stripe max-height="300">
              <el-table-column type="index" :label="$t('auto.List.faaadc44')" width="60" />
              <el-table-column prop="target" :label="$t('auto.List.73e82552')" min-width="150" />
              <el-table-column prop="error" :label="$t('auto.List.2fc597a3')" min-width="250" show-overflow-tooltip />
              <el-table-column prop="time" :label="$t('auto.List.19fcb9eb')" width="170" />
            </el-table>
          </div>
        </template>
        <el-empty v-if="!detailLoading && !taskDetail" :description="$t('auto.List.e10f7697')" />
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
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

const getPreviewStatusLabel = (s) => getStatusLabel(s, PASS_FAIL_STATUS) === t('auto.List.330363df') ? t('auto.List.596149b3') : (getStatusLabel(s, PASS_FAIL_STATUS) || s || '-');
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
const getStatusText = (s) => ({ running: t('auto.List.46e38679'), success: t('auto.List.769d88e4'), failed: t('auto.List.acd5cb84'), partial: t('auto.List.983b6813') }[s] || s)

const loadTools = async () => {
  const res = await getBatchTools()
  tools.value = res || [
    { name: t('auto.List.3182c77a'), description: t('auto.List.17fb37fa'), color: '#4F46E5', icon: 'ChatLineRound', usedCount: 0 },
    { name: t('auto.List.be8c0ab0'), description: 'Excel/CSV 批量导入', color: '#10B981', icon: 'Document', usedCount: 0 },
    { name: t('auto.List.52140dcf'), description: t('auto.List.05e71378'), color: '#F59E0B', icon: 'Promotion', usedCount: 0 },
    { name: t('auto.List.b2879bce'), description: t('auto.List.9ceef428'), color: '#EF4444', icon: 'User', usedCount: 0 }
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
    ElMessage.error(t('auto.List.943ad7a2') + (error?.message || error))
    previewVisible.value = false
  } finally {
    previewLoading.value = false
  }
}

const execute = async () => {
  const loading = ElMessage({ message: i18n.global.t(t('auto.List.c380d742')), type: 'info', duration: 0 })
  try {
    await runBatch({ tool: currentTool.value.name, ...form.value })
    loading.close()
    ElMessage.success(i18n.global.t(t('auto.List.b29706e5')))
    dialogVisible.value = false
    previewVisible.value = false
    loadHistories()
  } catch (error) {
    loading.close()
    ElMessage.error(i18n.global.t(t('auto.List.1c83d797')))
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
    await ElMessageBox.confirm(t('auto.List.1fed3850'), t('auto.List.e83a256e'), { type: 'warning' })
    await cancelBatch(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.2111ccbb')))
    loadHistories()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.482e8539')))
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
