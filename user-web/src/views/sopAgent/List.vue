<template>
  <div class="sop-agent-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.List.849893f1') }}</h2>
          <p class="subtitle">{{ $t('auto.List.8138a1e5') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="refreshAll" :loading="loading">
            <el-icon><Refresh /></el-icon>
            {{ $t('刷新') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <el-tabs v-model="activeMainTab" class="main-tabs" @tab-change="onMainTabChange">
      
      <el-tab-pane ::label="$t('auto.List.bdeab318')" name="manage">
        
        <el-row :gutter="16" class="stats-row" v-loading="statsLoading">
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ sopStats.total ?? 0 }}</div>
              <div class="stat-label">{{ $t('auto.List.bd430606') }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card stat-active">
              <div class="stat-value">{{ sopStats.active ?? 0 }}</div>
              <div class="stat-label">{{ $t('已激活') }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card stat-inactive">
              <div class="stat-value">{{ sopStats.inactive ?? 0 }}</div>
              <div class="stat-label">{{ $t('已停用') }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card stat-running">
              <div class="stat-value">{{ sopStats.running ?? 0 }}</div>
              <div class="stat-label">{{ $t('进行中执行') }}</div>
            </el-card>
          </el-col>
        </el-row>

        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.91ee8a4f') }}</span>
              <div class="header-controls">
                <el-input v-model="sopFilter.keyword" ::placeholder="$t('auto.List.dcba2b9d')" clearable size="small" style="width: 180px" @keyup.enter="onSopSearch" @clear="onSopSearch" />
                <el-select v-model="sopFilter.status" ::placeholder="$t('auto.List.1b857887')" clearable size="small" style="width: 120px" @change="onSopSearch">
                  <el-option ::label="$t('auto.List.82b015cb')" value="active" />
                  <el-option ::label="$t('auto.List.44a5154d')" value="inactive" />
                  <el-option ::label="$t('auto.List.7f46da48')" value="draft" />
                </el-select>
                <el-button size="small" type="primary" @click="onSopSearch">{{ $t('查询') }}</el-button>
                <el-button size="small" type="success" @click="openCreateDialog">
                  <el-icon><Plus /></el-icon> 创建 SOP
                </el-button>
              </div>
            </div>
          </template>
          <el-table :data="sopList" v-loading="sopLoading" stripe>
            <el-table-column prop="name" :label="$t('auto.List.d7ec2d3f')" min-width="160" show-overflow-tooltip />
            <el-table-column prop="description" :label="$t('auto.List.3bdd08ad')" min-width="200" show-overflow-tooltip />
            <el-table-column :label="$t('auto.List.3fea7ca7')" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.93cb163f')" width="90" align="center">
              <template #default="{ row }">{{ (row.nodes || []).length || row.node_count || 0 }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.34f19c98')" width="120">
              <template #default="{ row }">
                <el-tag v-if="row.trigger_intent" size="small" effect="plain">{{ row.trigger_intent }}</el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.eca37cb0')" width="160">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.2b6bc0f2')" width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="viewSopDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
                <el-button link type="warning" size="small" @click="openEditDialog(row)">{{ $t('auto.List.95b351c8') }}</el-button>
                <el-button v-if="row.status !== 'active'" link type="success" size="small" @click="toggleStatus(row, 'activate')">{{ $t('auto.List.83a991d7') }}</el-button>
                <el-button v-else link type="info" size="small" @click="toggleStatus(row, 'deactivate')">{{ $t('auto.List.5c56a889') }}</el-button>
                <el-button link type="danger" size="small" @click="removeSop(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.3bb9dbe8')" />
            </template>
          </el-table>
          <div class="pagination-container" v-if="sopPagination.total > 0">
            <el-pagination
              v-model:current-page="sopPagination.page"
              v-model:page-size="sopPagination.pageSize"
              :page-sizes="[10, 20, 50]"
              layout="total, prev, pager, next"
              :total="sopPagination.total"
              @current-change="loadSopList"
              @size-change="loadSopList"
            />
          </div>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.0a8e2412')" name="executions">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.6b5bc9be') }}</span>
              <div class="header-controls">
                <el-input v-model="execFilter.customer_id" :placeholder="$t('auto.List.3f61dec5')" clearable size="small" style="width: 160px" @keyup.enter="onExecSearch" @clear="onExecSearch" />
                <el-select v-model="execFilter.status" :placeholder="$t('auto.List.b56fe3d1')" clearable size="small" style="width: 140px" @change="onExecSearch">
                  <el-option :label="$t('auto.List.d679aea3')" value="running" />
                  <el-option :label="$t('auto.List.a2d930fd')" value="paused" />
                  <el-option :label="$t('auto.List.fad5222c')" value="completed" />
                  <el-option :label="$t('auto.List.2111ccbb')" value="cancelled" />
                  <el-option :label="$t('auto.List.acd5cb84')" value="failed" />
                </el-select>
                <el-button size="small" type="primary" @click="onExecSearch">{{ $t('auto.List.bee912d7') }}</el-button>
              </div>
            </div>
          </template>
          <el-table :data="executionList" v-loading="execLoading" stripe>
            <el-table-column :label="$t('auto.List.ff0b2077')" min-width="140">
              <template #default="{ row }">{{ row.customer_id || row.customer_name || '-' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.fad3e81a')" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">{{ row.sop_name || row.sop_id || '-' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.f9bad394')" width="140">
              <template #default="{ row }">
                <el-tag size="small" type="warning">{{ row.current_node || row.current_node_id || '-' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.3fea7ca7')" width="100">
              <template #default="{ row }">
                <el-tag :type="getExecStatusType(row.status)" size="small">{{ getExecStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.c7bff79d')" width="180">
              <template #default="{ row }">
                <el-progress :percentage="calcProgress(row)" :status="getProgressStatus(row.status)" />
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.592c5958')" width="160">
              <template #default="{ row }">{{ formatTime(row.started_at || row.created_at) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.2b6bc0f2')" width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="viewExecutionDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
                <el-button v-if="row.status === 'running'" link type="warning" size="small" @click="pauseExecution(row)">{{ $t('auto.List.8d63ef38') }}</el-button>
                <el-button v-if="row.status === 'paused'" link type="success" size="small" @click="resumeExecution(row)">{{ $t('auto.List.c7db6d4f') }}</el-button>
                <el-button v-if="['running', 'paused'].includes(row.status)" link type="danger" size="small" @click="cancelExecution(row)">{{ $t('auto.List.625fb26b') }}</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.124d09f1')" />
            </template>
          </el-table>
          <div class="pagination-container" v-if="execPagination.total > 0">
            <el-pagination
              v-model:current-page="execPagination.page"
              v-model:page-size="execPagination.pageSize"
              :page-sizes="[10, 20, 50]"
              layout="total, prev, pager, next"
              :total="execPagination.total"
              @current-change="loadExecutions"
              @size-change="loadExecutions"
            />
          </div>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.d44f774c')" name="match">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.d44f774c') }}</span>
              <span class="match-tip">{{ $t('auto.List.e6f7cab3') }}</span>
            </div>
          </template>
          <el-form label-width="90px">
            <el-form-item :label="$t('auto.List.ee9ff0a5')">
              <el-input
                v-model="matchIntent"
                :placeholder="$t('auto.List.04e307e2')"
                clearable
                @keyup.enter="runMatch"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="matchLoading" @click="runMatch">
                <el-icon><Search /></el-icon> 匹配 SOP
              </el-button>
              <el-button @click="quickFillIntent">{{ $t('auto.List.49e7f968') }}</el-button>
            </el-form-item>
          </el-form>

          <template v-if="matchResults.length">
            <el-divider content-position="left">匹配结果（{{ matchResults.length }} 个）</el-divider>
            <el-table :data="matchResults" stripe>
              <el-table-column prop="name" :label="$t('auto.List.fad3e81a')" min-width="160" show-overflow-tooltip />
              <el-table-column prop="description" :label="$t('auto.List.3bdd08ad')" min-width="200" show-overflow-tooltip />
              <el-table-column :label="$t('auto.List.fbdae97f')" width="160">
                <template #default="{ row }">
                  <el-progress :percentage="Math.round((row.score || row.match_score || 0) * 100)" :stroke-width="14" :text-inside="true" />
                </template>
              </el-table-column>
              <el-table-column :label="$t('auto.List.3fea7ca7')" width="100">
                <template #default="{ row }">
                  <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusLabel(row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="$t('auto.List.93cb163f')" width="90" align="center">
                <template #default="{ row }">{{ (row.nodes || []).length || row.node_count || 0 }}</template>
              </el-table-column>
              <el-table-column :label="$t('auto.List.2b6bc0f2')" width="160" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" size="small" @click="openExecuteDialog(row)">{{ $t('auto.List.1a6aa24e') }}</el-button>
                  <el-button link type="info" size="small" @click="viewSopDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
                </template>
              </el-table-column>
            </el-table>
          </template>
          <el-empty v-else-if="!matchLoading && hasSearched" :description="$t('auto.List.5cab5ed7')" :image-size="80" />
        </el-card>
      </el-tab-pane>
    </el-tabs>

    
    <el-dialog v-model="sopDialogVisible" :title="editingSop.id ? '编辑 SOP' : '创建 SOP'" width="720px" top="6vh">
      <el-form :model="sopForm" label-width="100px">
        <el-form-item :label="$t('auto.List.fad3e81a')" required>
          <el-input v-model="sopForm.name" :placeholder="$t('auto.List.74e27b6d')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3bdd08ad')">
          <el-input v-model="sopForm.description" type="textarea" :rows="2" :placeholder="$t('auto.List.e9b8b84c')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.34f19c98')">
          <el-input v-model="sopForm.trigger_intent" :placeholder="$t('auto.List.14fa9d5d')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3fea7ca7')">
          <el-select v-model="sopForm.status" style="width: 100%">
            <el-option :label="$t('auto.List.22b4334f')" value="draft" />
            <el-option :label="$t('auto.List.f6ebf8f5')" value="active" />
            <el-option :label="$t('auto.List.69b0f684')" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.d5840185')">
          <div class="nodes-editor">
            <div v-for="(node, idx) in sopForm.nodes" :key="idx" class="node-item">
              <el-input v-model="node.id" :placeholder="$t('auto.List.0458b995')" style="width: 110px" />
              <el-select v-model="node.type" :placeholder="$t('auto.List.bcafffc3')" style="width: 130px" filterable>
                <el-option v-for="t in sopNodeTypeOptions" :key="t.value" :label="t.label" :value="t.value" />
              </el-select>
              <el-input v-model="node.name" :placeholder="$t('auto.List.b1785ef0')" style="flex: 1" />
              <el-input v-model="node.action" :placeholder="$t('auto.List.d1d90491')" style="width: 140px" />
              <el-button link type="danger" @click="removeNode(idx)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <el-button type="primary" link @click="addNode">
              <el-icon><Plus /></el-icon> 添加节点
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="sopDialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="saveSop">{{ $t('auto.List.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="sopDetailVisible" :title="$t('auto.List.701911a6')" width="680px">
      <template v-if="currentSop">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.d7ec2d3f')">{{ currentSop.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.3fea7ca7')">
            <el-tag :type="getStatusType(currentSop.status)" size="small">{{ getStatusLabel(currentSop.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.34f19c98')">{{ currentSop.trigger_intent || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.93cb163f')">{{ (currentSop.nodes || []).length }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.3bdd08ad')" :span="2">{{ currentSop.description || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.eca37cb0')" :span="2">{{ formatTime(currentSop.created_at) }}</el-descriptions-item>
        </el-descriptions>
        <el-divider content-position="left">{{ $t('auto.List.ce4eca14') }}</el-divider>
        <el-steps :active="(currentSop.nodes || []).length" align-center v-if="(currentSop.nodes || []).length">
          <el-step v-for="(node, idx) in (currentSop.nodes || [])" :key="idx" :title="node.name || node.id" :description="node.action" />
        </el-steps>
        <el-empty v-else :description="$t('auto.List.f1729dc2')" :image-size="80" />
      </template>
    </el-dialog>

    
    <el-dialog v-model="execDetailVisible" :title="$t('auto.List.2eceed7d')" width="760px">
      <div v-loading="execDetailLoading" style="min-height: 120px">
      <template v-if="currentExecution">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.818e17e4')">{{ currentExecution.id || currentExecution.execution_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.ff0b2077')">{{ currentExecution.customer_id || currentExecution.customer_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.fad3e81a')">{{ currentExecution.sop_name || currentExecution.sop_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.3fea7ca7')">
            <el-tag :type="getExecStatusType(currentExecution.status)" size="small">{{ getExecStatusLabel(currentExecution.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.f9bad394')">{{ currentExecution.current_node || currentExecution.current_node_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.592c5958')">{{ formatTime(currentExecution.started_at || currentExecution.created_at) }}</el-descriptions-item>
        </el-descriptions>
        <el-divider content-position="left">{{ $t('auto.List.5ab41ba4') }}</el-divider>
        <el-steps :active="currentNodeIndex" align-center v-if="(currentExecution.nodes || currentExecution.steps || []).length">
          <el-step
            v-for="(node, idx) in (currentExecution.nodes || currentExecution.steps || [])"
            :key="idx"
            :title="node.name || node.node_id || `节点${idx + 1}`"
            :description="getNodeStepDesc(node)"
            :status="getNodeStepStatus(node, idx)"
          />
        </el-steps>
        <el-empty v-else :description="$t('auto.List.92229b1e')" :image-size="80" />
      </template>
      </div>
    </el-dialog>

    
    <el-dialog v-model="executeDialogVisible" :title="$t('auto.List.d9c3e220')" width="520px">
      <el-form :model="executeForm" label-width="100px">
        <el-form-item label="SOP">
          <el-input :model-value="executeForm.sop_name" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3f61dec5')" required>
          <el-input v-model="executeForm.customer_id" :placeholder="$t('auto.List.8b3ba498')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.50f198f0')">
          <el-input v-model="executeForm.context" type="textarea" :rows="3" :placeholder="$t('auto.List.f711ccef')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="executeDialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="executing" @click="confirmExecute">{{ $t('auto.List.00d76b2e') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus, Delete, Search } from '@element-plus/icons-vue'
import { sopApi } from '@/api/sopAgent.js'

const { t } = useI18n()

const loading = ref(false)
const activeMainTab = ref('manage')

const statsLoading = ref(false);
const sopStats = ref({})
const sopLoading = ref(false)
const sopList = ref([])
const sopFilter = ref({ keyword: '', status: '' })
const sopPagination = ref({ page: 1, pageSize: 20, total: 0 })

const sopDialogVisible = ref(false);
const saving = ref(false)
const editingSop = ref({})
const sopForm = ref({ name: '', description: '', trigger_intent: '', status: 'draft', nodes: [] })

const sopNodeTypeOptions = [
  { value: 'start', label: t('auto.List.e96a5bf0') },
  { value: 'end', label: t('auto.List.30ccd9cf') },
  { value: 'message', label: t('auto.List.83636b55') },
  { value: 'greeting', label: t('auto.List.6adcdb86') },
  { value: 'inquire', label: t('auto.List.9cbcd274') },
  { value: 'introduce', label: t('auto.List.fe6daec3') },
  { value: 'handle', label: t('auto.List.a6a6d46c') },
  { value: 'close', label: t('auto.List.4a4f0c70') },
  { value: 'invite', label: t('auto.List.c33a37fe') },
  { value: 'follow_up', label: t('auto.List.dbcfa189') },
  { value: 'activate', label: t('auto.List.8ffa5fc6') },
  { value: 'nurture', label: t('auto.List.fa82526d') },
  { value: 'condition', label: t('auto.List.fe4481c8') },
  { value: 'branch', label: t('auto.List.88724866') },
  { value: 'llm', label: 'LLM 决策 (llm)' },
  { value: 'ai_decide', label: 'AI 决策(旧) (ai_decide)' },
  { value: 'wait', label: t('auto.List.64f38136') },
  { value: 'action', label: t('auto.List.5ece43ff') },
  { value: 'send_offer', label: t('auto.List.14bcc550') }
];

const sopDetailVisible = ref(false);
const currentSop = ref(null)

const execLoading = ref(false);
const executionList = ref([])
const execFilter = ref({ customer_id: '', status: '' })
const execPagination = ref({ page: 1, pageSize: 20, total: 0 })

const execDetailVisible = ref(false)
const execDetailLoading = ref(false)
const currentExecution = ref(null)

const matchIntent = ref('');
const matchLoading = ref(false)
const matchResults = ref([])
const hasSearched = ref(false)

const executeDialogVisible = ref(false);
const executing = ref(false)
const executeForm = ref({ sop_id: '', sop_name: '', customer_id: '', context: '' })

const currentNodeIndex = computed(() => {
  if (!currentExecution.value) return 0
  const nodes = currentExecution.value.nodes || currentExecution.value.steps || []
  const current = currentExecution.value.current_node || currentExecution.value.current_node_id
  if (!current) return nodes.length
  const idx = nodes.findIndex(n => (n.id || n.node_id) === current)
  return idx >= 0 ? idx + 1 : nodes.length
})

const getStatusType = (status) => {
  if (status === 'active') return 'success'
  if (status === 'inactive') return 'info'
  if (status === 'draft') return 'warning'
  return 'info'
}

const getStatusLabel = (status) => {
  const map = { active: t('auto.List.f6ebf8f5'), inactive: t('auto.List.69b0f684'), draft: t('auto.List.22b4334f') }
  return map[status] || status || '-'
}

const getExecStatusType = (status) => {
  const map = { running: 'primary', paused: 'warning', completed: 'success', cancelled: 'info', failed: 'danger' }
  return map[status] || 'info'
}

const getExecStatusLabel = (status) => {
  const map = { running: t('auto.List.d679aea3'), paused: t('auto.List.a2d930fd'), completed: t('auto.List.fad5222c'), cancelled: t('auto.List.2111ccbb'), failed: t('auto.List.acd5cb84') }
  return map[status] || status || '-'
}

const getProgressStatus = (status) => {
  if (status === 'completed') return 'success'
  if (status === 'failed') return 'exception'
  if (status === 'cancelled') return 'exception'
  return undefined
}

const calcProgress = (row) => {
  const nodes = row.nodes || row.steps || []
  if (!nodes.length) return 0
  const current = row.current_node || row.current_node_id
  if (!current) return row.status === 'completed' ? 100 : 0
  const idx = nodes.findIndex(n => (n.id || n.node_id) === current)
  if (idx < 0) return row.status === 'completed' ? 100 : 0
  return Math.round(((idx + 1) / nodes.length) * 100)
}

const getNodeStepDesc = (node) => {
  if (!node) return ''
  const parts = []
  if (node.action) parts.push(`动作: ${node.action}`)
  if (node.status) parts.push(`状态: ${getExecStatusLabel(node.status)}`)
  if (node.executed_at || node.completed_at) parts.push(formatTime(node.executed_at || node.completed_at))
  return parts.join(' / ')
}

const getNodeStepStatus = (node, idx) => {
  if (!currentExecution.value) return 'wait'
  const nodeStatus = node.status
  if (nodeStatus === 'completed' || nodeStatus === 'success') return 'success'
  if (nodeStatus === 'failed' || nodeStatus === 'error') return 'error'
  if (nodeStatus === 'running' || nodeStatus === 'active') return 'process'
  const nodes = currentExecution.value.nodes || currentExecution.value.steps || [];
  const current = currentExecution.value.current_node || currentExecution.value.current_node_id
  const currentIdx = nodes.findIndex(n => (n.id || n.node_id) === current)
  if (currentIdx >= 0) {
    if (idx < currentIdx) return 'success'
    if (idx === currentIdx) return 'process'
    return 'wait'
  }
  if (currentExecution.value.status === 'completed') return 'success'
  return 'wait'
}

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadStats = async () => {
  statsLoading.value = true
  try {
    const res = await sopApi.getStats()
    sopStats.value = res || {}
  } catch (e) {
    ElMessage.error(t('auto.List.4f0b6e2e') + (e.message || t('auto.List.974e7484')))
  } finally {
    statsLoading.value = false
  }
};

const onSopSearch = () => {
  sopPagination.value.page = 1
  loadSopList()
}

const loadSopList = async () => {
  sopLoading.value = true
  try {
    const params = {
      page: sopPagination.value.page,
      page_size: sopPagination.value.pageSize
    }
    if (sopFilter.value.keyword) params.keyword = sopFilter.value.keyword
    if (sopFilter.value.status) params.status = sopFilter.value.status
    const res = await sopApi.list(params)
    sopList.value = res?.list || []
    sopPagination.value.total = res?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.List.6a7dfe50') + (e.message || t('auto.List.974e7484')))
  } finally {
    sopLoading.value = false
  }
}

const onExecSearch = () => {
  execPagination.value.page = 1
  loadExecutions()
}

const loadExecutions = async () => {
  execLoading.value = true
  try {
    const params = {
      page: execPagination.value.page,
      page_size: execPagination.value.pageSize
    }
    if (execFilter.value.customer_id) params.customer_id = execFilter.value.customer_id
    if (execFilter.value.status) params.status = execFilter.value.status
    const res = await sopApi.listExecutions(params)
    executionList.value = res?.list || []
    execPagination.value.total = res?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.List.f93d8b57') + (e.message || t('auto.List.974e7484')))
  } finally {
    execLoading.value = false
  }
}

const onMainTabChange = (tab) => {
  if (tab === 'manage') {
    loadStats()
    loadSopList()
  } else if (tab === 'executions') {
    loadExecutions()
  }
}

const refreshAll = async () => {
  loading.value = true
  try {
    if (activeMainTab.value === 'manage') {
      await Promise.all([loadStats(), loadSopList()])
    } else if (activeMainTab.value === 'executions') {
      await loadExecutions()
    }
  } finally {
    loading.value = false
  }
}

const resetSopForm = () => {
  sopForm.value = { name: '', description: '', trigger_intent: '', status: 'draft', nodes: [] }
  editingSop.value = {}
};

const addNode = () => {
  sopForm.value.nodes.push({ id: '', type: 'message', name: '', action: '' })
}

const removeNode = (idx) => {
  sopForm.value.nodes.splice(idx, 1)
}

const openCreateDialog = () => {
  resetSopForm()
  sopDialogVisible.value = true
}

const openEditDialog = async (row) => {
  resetSopForm()
  editingSop.value = row
  try {
    const res = await sopApi.get(row.id)
    const detail = res || row
    sopForm.value = {
      name: detail.name || '',
      description: detail.description || '',
      trigger_intent: detail.trigger_intent || '',
      status: detail.status || 'draft',
      nodes: (detail.nodes || []).map(n => ({
        id: n.id || '',
        type: n.type || 'message',
        name: n.name || '',
        action: n.action || ''
      }))
    }
    sopDialogVisible.value = true
  } catch (e) {
    ElMessage.error(t('auto.List.d172e6a3') + (e.message || t('auto.List.974e7484')))
  }
}

const saveSop = async () => {
  if (!sopForm.value.name || !sopForm.value.name.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.74e27b6d')))
    return
  }
  saving.value = true
  try {
    const data = {
      name: sopForm.value.name,
      description: sopForm.value.description,
      trigger_intent: sopForm.value.trigger_intent,
      status: sopForm.value.status,
      nodes: sopForm.value.nodes.filter(n => n.id || n.name).map(n => ({
        id: n.id || '',
        type: n.type || 'message',
        name: n.name || '',
        action: n.action || '',
        next: n.next || []
      }))
    }
    if (editingSop.value.id) {
      await sopApi.update(editingSop.value.id, data)
      ElMessage.success(i18n.global.t('SOP 更新成功'))
    } else {
      await sopApi.create(data)
      ElMessage.success(i18n.global.t('SOP 创建成功'))
    }
    sopDialogVisible.value = false
    loadSopList()
    loadStats()
  } catch (e) {
    ElMessage.error(t('auto.List.78d235cc') + (e.message || t('auto.List.974e7484')))
  } finally {
    saving.value = false
  }
}

const toggleStatus = async (row, action) => {
  try {
    if (action === 'activate') {
      await sopApi.activate(row.id)
      ElMessage.success(i18n.global.t('SOP 已激活'))
    } else {
      await sopApi.deactivate(row.id)
      ElMessage.success(i18n.global.t('SOP 已停用'))
    }
    loadSopList()
    loadStats()
  } catch (e) {
    ElMessage.error(t('auto.List.79585f92') + (e.message || t('auto.List.974e7484')))
  }
}

const removeSop = (row) => {
  ElMessageBox.confirm(`确认删除 SOP「${row.name}」吗？`, t('auto.List.50eaf94d'), {
    type: 'warning',
    confirmButtonText: t('auto.List.631cd220'),
    cancelButtonText: t('auto.List.625fb26b')
  }).then(async () => {
    try {
      await sopApi.remove(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
      loadSopList()
      loadStats()
    } catch (e) {
      ElMessage.error(t('auto.List.bc42bb8f') + (e.message || t('auto.List.974e7484')))
    }
  }).catch(() => {})
}

const viewSopDetail = async (row) => {
  try {
    const res = await sopApi.get(row.id)
    currentSop.value = res || row
    sopDetailVisible.value = true
  } catch (e) {
    currentSop.value = row
    sopDetailVisible.value = true
  }
}

const viewExecutionDetail = async (row) => {
  execDetailVisible.value = true
  execDetailLoading.value = true
  try {
    const id = row.id || row.execution_id
    const res = await sopApi.getExecution(id)
    currentExecution.value = res || row
  } catch (e) {
    ElMessage.error(t('auto.List.99d2cd5c') + (e.message || t('auto.List.974e7484')))
    currentExecution.value = row
  } finally {
    execDetailLoading.value = false
  }
};

const pauseExecution = async (row) => {
  try {
    await sopApi.pauseExecution(row.id || row.execution_id)
    ElMessage.success(i18n.global.t(t('auto.List.611f2522')))
    loadExecutions()
  } catch (e) {
    ElMessage.error(t('auto.List.f37e77dc') + (e.message || t('auto.List.974e7484')))
  }
}

const resumeExecution = async (row) => {
  try {
    await sopApi.resumeExecution(row.id || row.execution_id)
    ElMessage.success(i18n.global.t(t('auto.List.fd1b5088')))
    loadExecutions()
  } catch (e) {
    ElMessage.error(t('auto.List.5dc657e5') + (e.message || t('auto.List.974e7484')))
  }
}

const cancelExecution = (row) => {
  ElMessageBox.confirm(t('auto.List.6961a58b'), t('auto.List.ff748aa1'), {
    type: 'warning',
    confirmButtonText: t('auto.List.e8606af7'),
    cancelButtonText: t('auto.List.24fb8692')
  }).then(async () => {
    try {
      await sopApi.cancelExecution(row.id || row.execution_id)
      ElMessage.success(i18n.global.t(t('auto.List.14d5f1db')))
      loadExecutions()
    } catch (e) {
      ElMessage.error(t('auto.List.7b2cea68') + (e.message || t('auto.List.974e7484')))
    }
  }).catch(() => {})
}

const runMatch = async () => {
  if (!matchIntent.value || !matchIntent.value.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.e979c360')))
    return
  }
  matchLoading.value = true
  hasSearched.value = true
  matchResults.value = []
  try {
    const res = await sopApi.matchByIntent({ intent: matchIntent.value.trim() })
    matchResults.value = Array.isArray(res) ? res : (res?.list || [])
  } catch (e) {
    ElMessage.error(t('auto.List.a4ccd593') + (e.message || t('auto.List.974e7484')))
  } finally {
    matchLoading.value = false
  }
};

const quickFillIntent = () => {
  const intents = ['price_inquiry', 'purchase', 'objection_price', 'objection_trust', 'after_sale', 'churn']
  matchIntent.value = intents[Math.floor(Math.random() * intents.length)]
}

const openExecuteDialog = (row) => {
  executeForm.value = { sop_id: row.id, sop_name: row.name, customer_id: '', context: '' }
  executeDialogVisible.value = true
};

const confirmExecute = async () => {
  if (!executeForm.value.customer_id || !executeForm.value.customer_id.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.8b3ba498')))
    return
  }
  executing.value = true
  try {
    const data = {
      sop_id: executeForm.value.sop_id,
      customer_id: executeForm.value.customer_id
    }
    if (executeForm.value.context) {
      try {
        data.context = JSON.parse(executeForm.value.context)
      } catch {
        data.context = executeForm.value.context
      }
    }
    await sopApi.execute(data)
    ElMessage.success(i18n.global.t('SOP 已启动执行'))
    executeDialogVisible.value = false
    activeMainTab.value = 'executions'
    loadExecutions()
  } catch (e) {
    ElMessage.error(t('auto.List.95324f25') + (e.message || t('auto.List.974e7484')))
  } finally {
    executing.value = false
  }
}

loadStats();
loadSopList()
</script>

<style scoped lang="scss">
.sop-agent-page { padding: 20px; }

.header-card {
  margin-bottom: 16px;
  :deep(.el-card__body) { padding: 16px 20px; }
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    h2 { margin: 0 0 6px 0; font-size: 20px; }
    .subtitle { color: #909399; margin: 0; font-size: 13px; }
  }
}

.main-tabs { margin-top: 8px; }

.stats-row {
  margin-bottom: 16px;
  .stat-card {
    text-align: center;
    .stat-value { font-size: 28px; font-weight: bold; color: #303133; }
    .stat-label { color: #909399; font-size: 13px; margin-top: 6px; }
    &.stat-active .stat-value { color: #10B981; }
    &.stat-inactive .stat-value { color: #909399; }
    &.stat-running .stat-value { color: #4F46E5; }
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  .header-controls { display: flex; gap: 8px; align-items: center; }
  .match-tip { color: #909399; font-size: 12px; }
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.nodes-editor {
  width: 100%;
  .node-item {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
    align-items: center;
  }
}
</style>
