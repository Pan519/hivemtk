<template>
  <div class="geo-page">

    <div class="p-4" style="display:flex; gap:16px; min-height:calc(100vh - 180px)">
    
    <el-card style="width:260px; flex-shrink:0">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">{{ $t('auto.WorkflowEditor.19.8d6dc4') }}</span>
          <el-button size="small" type="primary" @click="onNewWorkflow">{{ $t('auto.WorkflowEditor.20.26bb84') }}</el-button>
        </div>
      </template>
      <div class="space-y-2">
        <div
          v-for="w in workflows"
          :key="w.id || w.name"
          class="workflow-item"
          role="button"
          tabindex="0"
          :class="{ active: currentWorkflow && (currentWorkflow.id === w.id || currentWorkflow.name === w.name) }"
          @click="selectWorkflow(w)"
          @keydown.enter.prevent="selectWorkflow(w)"
          @keydown.space.prevent="selectWorkflow(w)"
        >
          <div class="font-medium">{{ w.name }}</div>
          <div class="text-xs text-gray-500">{{ w.description || `${(w.steps || []).length} 个步骤` }}</div>
        </div>
      </div>
    </el-card>

    
    <div style="flex:1; display:flex; flex-direction:column; gap:16px">
      
      <el-card style="flex:1">
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="font-bold">工作流：{{ currentWorkflow?.name || '未选择' }}</span>
              <el-input v-model="wfName" :placeholder="$t('auto.WorkflowEditor.1.ec7361')" size="small" style="width:200px" v-if="currentWorkflow" />
            </div>
            <div class="flex items-center gap-2">
              <el-button size="small" @click="onSaveWorkflow" :disabled="!currentWorkflow">{{ $t('auto.WorkflowEditor.21.be5fbb') }}</el-button>
              <el-button size="small" type="danger" plain @click="onDeleteWorkflow" :disabled="!currentWorkflow?.id">{{ $t('auto.WorkflowEditor.22.2f4aad') }}</el-button>
              <el-button size="small" type="primary" @click="onRunNow" :disabled="!currentWorkflow" :loading="running">Run Now</el-button>
            </div>
          </div>
        </template>

        
        <div class="steps-container">
          <div
            v-for="(step, idx) in steps"
            :key="idx"
            class="step-block"
            :class="'step-' + step.type"
          >
            <div class="step-index">{{ idx + 1 }}</div>
            <div class="step-body">
              <div class="step-title">
                <el-tag v-if="stepMeta(step.type, 'kind') === 'control'" size="small" type="info" effect="plain" style="margin-right:6px">{{ stepMeta(step.type, 'label') }}</el-tag>
                <!-- 兜底不能用空串：stepMeta 对 STEP_TYPES 之外的陈旧 step.type 返回 undefined，
                     `undefined || ''` 会把空串传给 el-tag 的 type，触发 prop 校验告警并回落到
                     default(primary)。合法集合为 primary/success/info/warning/danger，中性用 info。 -->
                <el-tag v-else size="small" :type="stepMeta(step.type, 'color') || 'info'" effect="light" style="margin-right:6px">{{ stepMeta(step.type, 'label') }}</el-tag>
                <el-input v-model="step.name" size="small" :placeholder="$t('auto.WorkflowEditor.2.9a2536')" style="width:200px" />
              </div>
              <div class="step-config">
                
                <template v-if="step.type === 'trigger'">
                  <el-tag size="small" type="info">{{ $t('auto.WorkflowEditor.23.8d293e') }}</el-tag>
                </template>
                <template v-else-if="step.type === 'end'">
                  <el-tag size="small" type="info">{{ $t('auto.WorkflowEditor.24.ab65b4') }}</el-tag>
                </template>
                <template v-else-if="step.type === 'decision'">
                  <div class="flex gap-2 items-center">
                    <el-input v-model="step.condition" size="small" :placeholder="$t('auto.WorkflowEditor.3.07c216')" style="flex:1" />
                    <el-input v-model="step.jump_to" size="small" :placeholder="$t('auto.WorkflowEditor.4.1f4333')" style="width:160px" />
                  </div>
                </template>
                <template v-else-if="step.type === 'action'">
                  <el-tag size="small" type="info">{{ $t('auto.WorkflowEditor.25.57553c') }}</el-tag>
                </template>

                
                <template v-else-if="step.type === 'content_generate'">
                  <div class="flex gap-2 flex-wrap">
                    <el-input v-model="step.topic" size="small" :placeholder="$t('auto.WorkflowEditor.5.f9f840')" style="flex:1;min-width:140px" />
                    <el-input v-model="step.keyword" size="small" :placeholder="$t('auto.WorkflowEditor.6.011565')" style="width:160px" />
                    <el-input v-model="step.brand" size="small" :placeholder="$t('auto.WorkflowEditor.7.fc4a05')" style="width:140px" />
                    <el-input v-model="step.platform" size="small" :placeholder="$t('auto.WorkflowEditor.8.0b4228')" style="width:140px" />
                  </div>
                  <el-input v-model="step.advantages" size="small" :placeholder="$t('auto.WorkflowEditor.9.805268')" class="mt-2" />
                </template>
                <template v-else-if="step.type === 'content_score'">
                  <div class="flex gap-2 items-center">
                    <span class="text-xs text-gray-500">{{ $t('auto.WorkflowEditor.26.521a84') }}</span>
                    <el-input-number v-model="step.min_score" :min="0" :max="100" size="small" style="width:100px" />
                    <el-input v-model="step.brand" size="small" :placeholder="$t('auto.WorkflowEditor.7.fc4a05')" style="width:160px" />
                  </div>
                </template>
                <template v-else-if="step.type === 'eeat_enhance'">
                  <el-input v-model="step.brand" size="small" :placeholder="$t('auto.WorkflowEditor.7.fc4a05')" style="width:200px" />
                  <span class="text-xs text-gray-500 ml-2">{{ $t('auto.WorkflowEditor.27.c33daa') }}</span>
                </template>
                <template v-else-if="step.type === 'fact_density_enhance'">
                  <div class="flex gap-2 items-center">
                    <span class="text-xs text-gray-500">{{ $t('auto.WorkflowEditor.28.f02e11') }}</span>
                    <el-slider v-model="step.target_density" :min="0.1" :max="1.0" :step="0.05" style="flex:1" />
                  </div>
                </template>
                <template v-else-if="step.type === 'verify'">
                  <div class="flex gap-2 flex-wrap">
                    <el-input v-model="step.brand" size="small" :placeholder="$t('auto.WorkflowEditor.10.2284d7')" style="width:180px" />
                  </div>
                  <div class="text-xs text-gray-500 mt-1">{{ $t('auto.WorkflowEditor.29.d1a41a') }}</div>
                </template>
                <template v-else-if="step.type === 'query_probe'">
                  <div class="flex gap-2 flex-wrap">
                    <el-input v-model="step.keyword" size="small" :placeholder="$t('auto.WorkflowEditor.11.48a10a')" style="flex:1;min-width:160px" />
                    <el-select v-model="step.engine" size="small" :placeholder="$t('auto.WorkflowEditor.12.6c1b19')" style="width:140px">
                      <el-option label="Qwen" value="qwen" />
                      <el-option label="DeepSeek" value="deepseek" />
                      <el-option label="SenseNova" value="sensenova" />
                      <el-option label="Doubao" value="doubao" />
                    </el-select>
                  </div>
                </template>
                <template v-else-if="step.type === 'source_attribution'">
                  <el-tag size="small" type="info">{{ $t('auto.WorkflowEditor.30.c65f13') }}</el-tag>
                </template>
                <template v-else-if="step.type === 'content_gap_fill'">
                  <el-input v-model="step.keyword" size="small" :placeholder="$t('auto.WorkflowEditor.13.cddadc')" style="width:240px" />
                </template>
                <template v-else-if="step.type === 'capture_lead'">
                  <el-tag size="small" type="info">{{ $t('auto.WorkflowEditor.31.ece49a') }}</el-tag>
                </template>
                <template v-else>
                  <div class="text-xs text-gray-400">未知步骤类型：{{ step.type }}</div>
                </template>
              </div>
            </div>
            <div class="step-actions">
              <el-button size="small" circle @click="moveStep(idx, -1)" :disabled="idx === 0">↑</el-button>
              <el-button size="small" circle @click="moveStep(idx, 1)" :disabled="idx === steps.length - 1">↓</el-button>
              <el-button size="small" circle type="danger" plain @click="removeStep(idx)">×</el-button>
            </div>
          </div>

          <div class="mt-4">
            <div class="text-xs text-gray-500 mb-1">{{ $t('auto.WorkflowEditor.32.7945f0') }}</div>
            <div class="flex gap-2 flex-wrap mb-2">
              <template v-for="(t, k) in STEP_TYPES" :key="'c-'+k">
                <el-button v-if="t.kind==='control'" size="small" plain @click="addStep(k)">+ {{ t.label }}</el-button>
              </template>
            </div>
            <div class="text-xs text-gray-500 mb-1">{{ $t('auto.WorkflowEditor.33.b8348b') }}</div>
            <div class="flex gap-2 flex-wrap">
              <template v-for="(t, k) in STEP_TYPES" :key="'e-'+k">
                <el-button v-if="t.kind==='executor'" size="small" plain @click="addStep(k)">+ {{ t.label }}</el-button>
              </template>
            </div>
          </div>
        </div>
      </el-card>

      
      <el-card>
        <template #header><span class="font-bold">{{ $t('auto.WorkflowEditor.34.0a8e24') }}</span></template>
        <el-table :data="executions" v-loading="execLoading" size="small">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="workflow_name" label="Workflow" width="160" />
          <el-table-column prop="status" :label="$t('auto.WorkflowEditor.14.3fea7c')" width="120">
            <template #default="{ row }">
              <el-tag :type="execStatusType(row.status)" size="small">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="current_step" :label="$t('auto.WorkflowEditor.15.59cecb')" width="140" />
          <el-table-column prop="progress" :label="$t('auto.WorkflowEditor.16.c7bff7')" width="160">
            <template #default="{ row }">
              <el-progress :percentage="row.progress || 0" size="small" />
            </template>
          </el-table-column>
          <el-table-column prop="started_at" :label="$t('auto.WorkflowEditor.17.592c59')" width="180" />
          <el-table-column prop="finished_at" :label="$t('auto.WorkflowEditor.18.f78277')" width="180" />
        </el-table>
      </el-card>
    </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { listWorkflows, saveWorkflow, deleteWorkflow, runWorkflow, listWorkflowExecutions } from '@/api/geoEntity.js'

const STEP_TYPES = {
  trigger:  { label: t('auto.WorkflowEditor.35.469672'),     kind: 'control', default: {} },
  action:   { label: t('auto.WorkflowEditor.36.f429d4'),   kind: 'control', default: {} },
  decision: { label: t('auto.WorkflowEditor.37.c71b12'),   kind: 'control', default: { condition: '', jump_to: '' } },
  end:      { label: t('auto.WorkflowEditor.38.642f96'),   kind: 'control', default: {} },

  content_generate:     { label: 'LLM 内容生成', kind: 'executor', color: 'primary', default: { topic: '', brand: 'HiveMTK', keyword: '', advantages: '', platform: '' } },
  content_score:        { label: t('auto.WorkflowEditor.39.18a442'), kind: 'executor', color: 'warning', default: { min_score: 70, brand: 'HiveMTK' } },
  eeat_enhance:         { label: 'EEAT 权威增强', kind: 'executor', color: 'success', default: { brand: 'HiveMTK' } },
  fact_density_enhance: { label: t('auto.WorkflowEditor.40.3225f3'), kind: 'executor', color: 'success', default: { target_density: 0.8 } },
  verify:               { label: t('auto.WorkflowEditor.41.f314af'), kind: 'executor', color: 'info', default: { brand: 'HiveMTK', queries: [] } },
  query_probe:          { label: t('auto.WorkflowEditor.42.169e12'), kind: 'executor', color: 'primary', default: { keyword: '', engine: 'qwen' } },
  source_attribution:   { label: t('auto.WorkflowEditor.43.f7bd71'), kind: 'executor', color: 'info', default: {} },
  content_gap_fill:     { label: t('auto.WorkflowEditor.44.17d15d'), kind: 'executor', color: 'warning', default: { keyword: '' } },
  capture_lead:         { label: t('auto.WorkflowEditor.45.1fe100'), kind: 'executor', color: 'danger', default: {} },
}

const { t } = useI18n()

const workflows = ref([])
const currentWorkflow = ref(null)
const wfName = ref('')
const steps = ref([])
const running = ref(false)

const executions = ref([])
const execLoading = ref(false)

const stepMeta = (type, key) => {
  const meta = STEP_TYPES[type]
  if (!meta) return key === 'label' ? (type || t('auto.WorkflowEditor.46.1622dc')) : undefined
  return meta[key]
};

const selectWorkflow = (w) => {
  currentWorkflow.value = w
  wfName.value = w.name || ''
  steps.value = Array.isArray(w.steps) ? JSON.parse(JSON.stringify(w.steps)) : []
}

const onNewWorkflow = () => {
  currentWorkflow.value = { name: t('auto.WorkflowEditor.47.c3ad92'), description: '', steps: [] }
  wfName.value = t('auto.WorkflowEditor.47.c3ad92')
  steps.value = []
}

const addStep = (type) => {
  steps.value.push({ type, ...JSON.parse(JSON.stringify(STEP_TYPES[type].default)) })
}

const removeStep = (idx) => steps.value.splice(idx, 1)
const moveStep = (idx, dir) => {
  const target = idx + dir
  if (target < 0 || target >= steps.value.length) return
  const [item] = steps.value.splice(idx, 1)
  steps.value.splice(target, 0, item)
}

const onSaveWorkflow = async () => {
  try {
    const payload = {
      id: currentWorkflow.value?.id,
      name: wfName.value,
      description: currentWorkflow.value?.description,
      steps: steps.value
    }
    const result = await saveWorkflow(payload)
    ElMessage.success(t('auto.WorkflowEditor.48.3b1083'))
    loadWorkflows()
    if (result?.id) currentWorkflow.value.id = result.id
  } catch (e) {
    ElMessage.error(t('auto.WorkflowEditor.49.78d235') + (e?.message || e))
  }
}

const onDeleteWorkflow = async () => {
  if (!currentWorkflow.value?.id) return
  try {
    await deleteWorkflow(currentWorkflow.value.id)
    ElMessage.success(t('auto.WorkflowEditor.50.5cc232'))
    currentWorkflow.value = null
    steps.value = []
    loadWorkflows()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.WorkflowEditor.51.acf066'))
  }
}

const onRunNow = async () => {
  if (!currentWorkflow.value) return
  running.value = true
  try {
    await runWorkflow(currentWorkflow.value.id || wfName.value, { steps: steps.value })
    ElMessage.success(t('auto.WorkflowEditor.52.6f8c8d'))
    loadExecutions()
  } catch (e) {
    ElMessage.error(t('auto.WorkflowEditor.53.95324f') + (e?.message || e))
  } finally {
    running.value = false
  }
}

const execStatusType = (s) => {
  if (s === 'running' || s === 'pending') return 'warning'
  if (s === 'success' || s === 'done') return 'success'
  if (s === 'failed' || s === 'error') return 'danger'
  return 'info'
}

const loadWorkflows = async () => {
  try {
    const data = await listWorkflows()
    workflows.value = Array.isArray(data) ? data : (data?.list || data?.items || [])
  } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
const loadExecutions = async () => {
  execLoading.value = true
  try {
    const data = await listWorkflowExecutions(currentWorkflow.value?.id)
    executions.value = Array.isArray(data) ? data : (data?.list || data?.items || [])
  } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
  execLoading.value = false
}

onMounted(() => { loadWorkflows(); loadExecutions() })
</script>

<style lang="scss" scoped>
.workflow-item {
  padding: $spacing-md;
  border-radius: 6px;
  cursor: pointer;
  background: var(--el-fill-color-light);
  transition: background .15s;
}
.workflow-item:hover { background: var(--el-fill-color); }
.workflow-item.active { background: var(--el-color-primary-light-9); border: 1px solid var(--el-color-primary-light-5); }
.step-block {
  display: flex;
  align-items: flex-start;
  gap: $spacing-md;
  padding: $spacing-md 16px;
  border-radius: 8px;
  margin-bottom: $spacing-md;
  background: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color-lighter);
}
.step-index {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--el-color-primary);
  color: $bg-color;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: $font-size-base;
  flex-shrink: 0;
}
.step-body { flex: 1; min-width: 0; }
.step-title { font-weight: 600; margin-bottom: $spacing-sm; }
.step-config .el-input, .step-config .el-select { width: 100%; }
.step-actions { display: flex; flex-direction: column; gap: $spacing-xs; }
.step-keywords { border-left: 3px solid $warning-color; }
.step-generate { border-left: 3px solid $primary-color; }
.step-optimize { border-left: 3px solid $success-color; }
.step-verify { border-left: 3px solid $primary-color; }
.step-publish { border-left: 3px solid $danger-color; }

/* 控制节点 */
.step-trigger  { border-left: 3px solid #909399; background: var(--el-fill-color-light); }
.step-action   { border-left: 3px solid #606266; background: var(--el-fill-color-light); }
.step-decision { border-left: 3px solid #e6a23c; background: var(--el-color-warning-light-9); }
.step-end      { border-left: 3px solid #67c23a; background: var(--el-color-success-light-9); }

/* 执行器节点 */
.step-content_generate     { border-left: 3px solid var(--el-color-primary); }
.step-content_score        { border-left: 3px solid var(--el-color-warning); }
.step-eeat_enhance         { border-left: 3px solid var(--el-color-success); }
.step-fact_density_enhance { border-left: 3px solid var(--el-color-success); }
.step-verify               { border-left: 3px solid var(--el-color-info); }
.step-query_probe          { border-left: 3px solid var(--el-color-primary); }
.step-source_attribution   { border-left: 3px solid var(--el-color-info); }
.step-content_gap_fill     { border-left: 3px solid var(--el-color-warning); }
.step-capture_lead         { border-left: 3px solid var(--el-color-danger); }
</style>
