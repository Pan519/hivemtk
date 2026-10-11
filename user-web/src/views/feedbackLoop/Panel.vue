<template>
  <div class="feedback-loop-panel">
    <el-card class="header-card">
      <div class="header-content">
        <h2>{{ $t('反馈学习闭环') }}</h2>
        <p class="subtitle">{{ $t('auto.Panel.9868ef11') }}</p>
      </div>
    </el-card>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('auto.Panel.20db3a78') }}</div>
            <div class="stat-value">{{ stats.feedbackCount }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('销冠对话') }}</div>
            <div class="stat-value" style="color: #10B981">{{ stats.dialogueCount }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('auto.Panel.82146047') }}</div>
            <div class="stat-value" style="color: #4F46E5">{{ stats.candidateCount }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('auto.Panel.07bf9a67') }}</div>
            <div class="stat-value" style="color: #F59E0B">{{ stats.banditCount }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-tabs v-model="activeTab" class="content-tabs">
      
      <el-tab-pane ::label="$t('auto.Panel.9aa0ddbd')" name="events">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Panel.14de1a34') }}</span>
              <el-button @click="loadEvents">
                <el-icon><Refresh /></el-icon>
                {{ $t('刷新') }}
              </el-button>
            </div>
          </template>
          <el-row :gutter="12" style="margin-bottom: 12px">
            <el-col :span="6">
              <el-select v-model="eventTypeFilter" :placeholder="$t('auto.Panel.f974c846')" clearable style="width: 100%">
                <el-option :label="$t('auto.Panel.68f6d29e')" value="explicit" />
                <el-option :label="$t('auto.Panel.426f5c89')" value="implicit" />
                <el-option :label="$t('auto.Panel.4a27b5d5')" value="champion" />
              </el-select>
            </el-col>
            <el-col :span="6">
              <el-select v-model="signalKeyFilter" :placeholder="$t('auto.Panel.5fc98258')" clearable style="width: 100%">
                <el-option label="like" value="like" />
                <el-option label="dislike" value="dislike" />
                <el-option label="conversion" value="conversion" />
                <el-option label="reply_rate" value="reply_rate" />
                <el-option label="duration" value="duration" />
                <el-option label="transfer" value="transfer" />
                <el-option label="champion_mark" value="champion_mark" />
              </el-select>
            </el-col>
          </el-row>
          <el-table :data="events" v-loading="eventsLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="event_type" :label="$t('auto.Panel.226b0912')" width="100">
              <template #default="{ row }">
                <el-tag :type="eventTypeTagType(row.event_type)">{{ eventTypeLabel(row.event_type) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="signal_key" :label="$t('auto.Panel.eaaf95ff')" width="100" />
            <el-table-column prop="session_id" :label="$t('auto.Panel.9a834e9a')" width="160" />
            <el-table-column prop="reward" :label="$t('auto.Panel.667378a1')" width="100">
              <template #default="{ row }">
                <el-tag :type="confTagType(row.reward)">{{ (row.reward || 0).toFixed(3) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ai_reply" :label="$t('auto.Panel.ae829781')" min-width="280" show-overflow-tooltip />
            <el-table-column prop="customer_msg" :label="$t('auto.Panel.80ded089')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="created_at" :label="$t('auto.Panel.19fcb9eb')" width="170">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            :current-page="eventPage"
            :page-size="eventPageSize"
            :total="eventTotal"
            :page-sizes="[20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @current-change="loadEvents"
            @size-change="loadEvents"
            style="margin-top: 16px; text-align: right"
          />
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.Panel.971a4cbe')" name="dialogues">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Panel.d5736f39') }}</span>
              <el-button type="primary" @click="triggerDialogueAnalysis">
                <el-icon><MagicStick /></el-icon>
                立即触发分析
              </el-button>
            </div>
          </template>
          <el-table :data="dialogues" v-loading="dialoguesLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="scenario" :label="$t('auto.Panel.c931653c')" width="140" />
            <el-table-column prop="journey_stage" :label="$t('auto.Panel.439ef921')" width="120" />
            <el-table-column prop="staff_id" :label="$t('auto.Panel.c08c17c9')" width="100" />
            <el-table-column prop="customer_msg" :label="$t('auto.Panel.80ded089')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="champion_reply" :label="$t('auto.Panel.b8d94269')" min-width="280" show-overflow-tooltip />
            <el-table-column prop="reward" :label="$t('auto.Panel.c39a1ae0')" width="100">
              <template #default="{ row }">
                <el-tag :type="confTagType(row.reward)" effect="dark">{{ (row.reward || 0).toFixed(3) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="cluster_id" :label="$t('auto.Panel.88e7f33d')" width="100" />
            <el-table-column prop="created_at" :label="$t('auto.Panel.19fcb9eb')" width="170">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            :current-page="dialoguePage"
            :page-size="dialoguePageSize"
            :total="dialogueTotal"
            :page-sizes="[20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @current-change="loadDialogues"
            @size-change="loadDialogues"
            style="margin-top: 16px; text-align: right"
          />
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.Panel.8c68ef75')" name="prompts">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Panel.f40e5273') }}</span>
              <el-button @click="loadPrompts">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          <el-table :data="prompts" v-loading="promptsLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="sop_id" label="SOP" width="80" />
            <el-table-column prop="sop_node_id" :label="$t('auto.Panel.3bf3c0a8')" width="120" />
            <el-table-column prop="version" :label="$t('auto.Panel.fe2df04a')" width="100" />
            <el-table-column prop="title" :label="$t('auto.Panel.32c65d8d')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="status" :label="$t('auto.Panel.3fea7ca7')" width="100">
              <template #default="{ row }">
                <el-tag :type="promptStatusTagType(row.status)">{{ promptStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="generated_by" :label="$t('auto.Panel.26ca20b1')" width="80" />
            <el-table-column prop="created_at" :label="$t('auto.Panel.eca37cb0')" width="170">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Panel.2b6bc0f2')" width="180" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="viewPrompt(row)">{{ $t('auto.Panel.607e7a4f') }}</el-button>
                <el-button v-if="row.status === 'draft'" size="small" type="success" @click="approvePrompt(row)">{{ $t('auto.Panel.7a712b9a') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane label="Bandit A/B" name="bandit">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Panel.051b51c4') }}</span>
              <el-button @click="loadBandits">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          <el-table :data="bandits" v-loading="banditsLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="experiment_id" :label="$t('auto.Panel.1014aa99')" width="200" show-overflow-tooltip />
            <el-table-column prop="arm_key" :label="$t('auto.Panel.d55d9eb6')" width="100" />
            <el-table-column prop="status" :label="$t('auto.Panel.3fea7ca7')" width="100">
              <template #default="{ row }">
                <el-tag :type="banditStatusTagType(row.status)">{{ banditStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="alpha" label="α" width="80">
              <template #default="{ row }">{{ (row.alpha || 0).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="beta" label="β" width="80">
              <template #default="{ row }">{{ (row.beta || 0).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="successes" :label="$t('auto.Panel.330363df')" width="80" />
            <el-table-column prop="failures" :label="$t('auto.Panel.acd5cb84')" width="80" />
            <el-table-column prop="trials" :label="$t('auto.Panel.6cdece36')" width="80" />
            <el-table-column prop="current_traffic_pct" :label="$t('auto.Panel.8a83a364')" width="110">
              <template #default="{ row }">
                {{ (row.current_traffic_pct || 0).toFixed(1) }}%
              </template>
            </el-table-column>
            <el-table-column prop="posterior_best_prob" :label="$t('auto.Panel.29d5fc3f')" width="100">
              <template #default="{ row }">
                {{ (row.posterior_best_prob || 0).toFixed(3) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    
    <el-dialog v-model="promptDialogVisible" ::title="$t('auto.Panel.8dece260')" width="780px">
      <div v-if="currentPrompt">
        <p><strong>{{ $t('auto.Panel.2684c463') }}</strong>{{ currentPrompt.version }} · <strong>{{ $t('auto.Panel.ae9dede6') }}</strong>{{ currentPrompt.scenario }} · <strong>{{ $t('auto.Panel.bec98b4d') }}</strong>{{ promptStatusLabel(currentPrompt.status) }}</p>
        <el-divider />
        <h4>{{ $t('auto.Panel.277e62f5') }}</h4>
        <pre class="prompt-content">{{ currentPrompt.system_prompt }}</pre>
        <h4>{{ $t('auto.Panel.e6534649') }}</h4>
        <pre class="prompt-content">{{ currentPrompt.user_prompt_template }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getFeedbackEvents,
  getChampionDialogues,
  getPromptCandidates,
  updatePromptCandidateStatus,
  getBanditArms
} from '@/api/tuning'

const { t } = useI18n()

const activeTab = ref('events')

const stats = ref({ feedbackCount: 0, dialogueCount: 0, candidateCount: 0, banditCount: 0 });

const events = ref([]);
const eventPage = ref(1)
const eventPageSize = ref(20)
const eventTotal = ref(0)
const eventsLoading = ref(false)
const eventTypeFilter = ref('')
const signalKeyFilter = ref('')
async function loadEvents() {
  eventsLoading.value = true
  try {
    const params = { page: eventPage.value, page_size: eventPageSize.value }
    if (eventTypeFilter.value) params.event_type = eventTypeFilter.value
    if (signalKeyFilter.value) params.signal_key = signalKeyFilter.value
    const res = await getFeedbackEvents(params)
    events.value = res?.list || res?.data?.list || []
    eventTotal.value = res?.total || res?.data?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.Panel.dcc1e901') + (e?.message || e))
  } finally {
    eventsLoading.value = false
  }
}
watch([eventTypeFilter, signalKeyFilter], () => {
  eventPage.value = 1
  loadEvents()
})

const dialogues = ref([]);
const dialoguePage = ref(1)
const dialoguePageSize = ref(20)
const dialogueTotal = ref(0)
const dialoguesLoading = ref(false)
async function loadDialogues() {
  dialoguesLoading.value = true
  try {
    const res = await getChampionDialogues({ page: dialoguePage.value, page_size: dialoguePageSize.value })
    dialogues.value = res?.list || res?.data?.list || []
    dialogueTotal.value = res?.total || res?.data?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.Panel.29b9d4cf') + (e?.message || e))
  } finally {
    dialoguesLoading.value = false
  }
}
async function triggerDialogueAnalysis() {
  try {
    await ElMessageBox.confirm(t('auto.Panel.95fe77bd'), t('auto.Panel.02d9819d'), { type: 'info' })
    ElMessage.info(i18n.global.t(t('auto.Panel.2abfb4ae')));
  } catch (e) {
    if (e !== 'cancel' && e !== 'close') {
      ElMessage.error(t('auto.Panel.d195ec24') + (e?.message || e))
    }
  }
}

const prompts = ref([]);
const promptsLoading = ref(false)
async function loadPrompts() {
  promptsLoading.value = true
  try {
    const res = await getPromptCandidates({ page: 1, page_size: 50 })
    prompts.value = res?.list || res?.data?.list || []
  } catch (e) {
    ElMessage.error('Prompt 候选加载失败：' + (e?.message || e))
  } finally {
    promptsLoading.value = false
  }
}
const promptDialogVisible = ref(false)
const currentPrompt = ref(null)
function viewPrompt(row) {
  currentPrompt.value = row
  promptDialogVisible.value = true
}
async function approvePrompt(row) {
  try {
    await ElMessageBox.confirm(`确认批准 Prompt「${row.title}」？批准后将创建 Bandit A/B 实验。`, t('auto.Panel.02d9819d'), { type: 'warning' })
    await updatePromptCandidateStatus(row.id, { status: 'approved' })
    ElMessage.success(i18n.global.t(t('auto.Panel.337737f5')))
    await loadPrompts()
  } catch (e) {
    if (e !== 'cancel' && e !== 'close') {
      ElMessage.error(t('auto.Panel.32365752') + (e?.message || e))
    }
  }
}

const bandits = ref([]);
const banditsLoading = ref(false)
async function loadBandits() {
  banditsLoading.value = true
  try {
    const res = await getBanditArms({ page: 1, page_size: 100 })
    bandits.value = res?.list || res?.data?.list || []
  } catch (e) {
    ElMessage.error('Bandit 加载失败：' + (e?.message || e))
  } finally {
    banditsLoading.value = false
  }
}

function confTagType(v) {
  if (v >= 0.7) return 'success'
  if (v >= 0) return 'warning'
  return 'danger'
}
function eventTypeLabel(t) {
  return { explicit: t('auto.Panel.34a09b1f'), implicit: t('auto.Panel.00878ef5'), champion: t('auto.Panel.c08c17c9') }[t] || t
}
function eventTypeTagType(t) {
  return { explicit: 'primary', implicit: 'info', champion: 'success' }[t] || ''
}
function promptStatusLabel(s) {
  return { draft: t('auto.Panel.22b4334f'), approved: t('auto.Panel.337737f5'), active: t('auto.Panel.8d4a7497'), retired: t('auto.Panel.0a666759') }[s] || s
}
function promptStatusTagType(s) {
  return { draft: 'info', approved: 'success', active: 'success', retired: 'danger' }[s] || ''
}
function banditStatusLabel(s) {
  return { exploring: t('auto.Panel.67643ed9'), exploiting: t('auto.Panel.7f817d3d'), promoted: t('auto.Panel.c19943d6'), retired: t('auto.Panel.5b109a64') }[s] || s
}
function banditStatusTagType(s) {
  return { exploring: 'warning', exploiting: 'primary', promoted: 'success', retired: 'info' }[s] || ''
}
function formatTime(t) {
  if (!t) return '-'
  try { return new Date(t).toLocaleString('zh-CN', { hour12: false }) } catch { return t }
}

onMounted(async () => {
  await Promise.all([loadEvents(), loadDialogues(), loadPrompts(), loadBandits()])
  stats.value = {
    feedbackCount: eventTotal.value,
    dialogueCount: dialogueTotal.value,
    candidateCount: prompts.value.length,
    banditCount: bandits.value.length
  }
})
</script>

<style scoped>
.feedback-loop-panel { padding: 16px; }
.header-card { margin-bottom: 16px; }
.header-content h2 { margin: 0 0 4px 0; font-size: 20px; }
.subtitle { color: #909399; margin: 0; font-size: 13px; }
.stats-row { margin-bottom: 16px; }
.stat-item { text-align: center; }
.stat-label { color: #909399; font-size: 12px; margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: 600; color: #303133; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.prompt-content { background: #f5f7fa; padding: 12px; border-radius: 4px; line-height: 1.6; font-family: monospace; white-space: pre-wrap; word-break: break-word; }
</style>
