<template>
  <div class="churn-prediction-page">
    <el-card class="header-card">
      <div class="header-content">
        <h2>{{ $t('auto.List.e016323d') }}</h2>
        <p class="subtitle">{{ $t('auto.List.e4aa48c2') }}</p>
      </div>
      <div>
        <el-button type="primary" :loading="predicting" @click="runPrediction">
          <el-icon><Operation /></el-icon>
          运行预测
        </el-button>
        <el-button @click="refreshAll">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </el-card>

    <el-tabs v-model="activeTab">
      
      <el-tab-pane :label="$t('auto.List.a6440fee')" name="predictions">
        <el-row :gutter="20" class="stats-row">
          <el-col :span="6">
            <el-card>
              <div class="stat-item">
                <div class="stat-label">{{ $t('auto.List.28198f83') }}</div>
                <div class="stat-value" style="color: #EF4444">{{ stats.highRisk }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card>
              <div class="stat-item">
                <div class="stat-label">{{ $t('auto.List.36e164b6') }}</div>
                <div class="stat-value" style="color: #F59E0B">{{ stats.mediumRisk }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card>
              <div class="stat-item">
                <div class="stat-label">{{ $t('auto.List.d9c0b8a7') }}</div>
                <div class="stat-value" style="color: #10B981">{{ stats.lowRisk }}</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card>
              <div class="stat-item">
                <div class="stat-label">{{ $t('auto.List.beab8b64') }}</div>
                <div class="stat-value" style="color: #B91C1C">{{ stats.criticalRisk }}</div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.82b808f1') }}</span>
              <el-input v-model="searchKeyword" :placeholder="$t('auto.List.c3635c94')" clearable style="width: 200px" />
            </div>
          </template>
          <el-table :data="filteredUsers" v-loading="loadingPred" stripe>
            <el-table-column prop="user_id" :label="$t('auto.List.30acd20d')" width="120" />
            <el-table-column prop="churn_score" :label="$t('auto.List.8372f3ae')" width="160">
              <template #default="{ row }">
                <el-progress :percentage="Math.round(row.churn_score || 0)" :color="getRiskColor(row.churn_score)" />
              </template>
            </el-table-column>
            <el-table-column prop="churn_risk" :label="$t('auto.List.25eb6ba3')" width="100">
              <template #default="{ row }">
                <el-tag :type="getRiskLevelType(row.churn_risk)">{{ row.churn_risk }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="last_activity_at" :label="$t('auto.List.7f2846ef')" min-width="180" />
            <el-table-column prop="predicted_at" :label="$t('auto.List.e79e7b21')" min-width="180" />
            <el-table-column prop="risk_factors" :label="$t('auto.List.0bed4e95')" min-width="200" show-overflow-tooltip />
            <el-table-column :label="$t('auto.List.2b6bc0f2')" width="180" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="viewUserDetail(row)">{{ $t('auto.List.f26225bd') }}</el-button>
                <el-button link type="success" @click="intervene(row)">{{ $t('auto.List.ba996d64') }}</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.21efd88b')" />
            </template>
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.8c9bcbdd')" name="warnings">
        <el-card>
          <el-table :data="warnings" v-loading="loadingWarn" stripe>
            <el-table-column prop="user_id" :label="$t('auto.List.30acd20d')" width="120" />
            <el-table-column prop="warning_level" :label="$t('auto.List.64d58e05')" width="100">
              <template #default="{ row }">
                <el-tag :type="warnLevelType(row.warning_level)">{{ row.warning_level }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="warning_type" :label="$t('auto.List.226b0912')" width="140" />
            <el-table-column prop="description" :label="$t('auto.List.3bdd08ad')" min-width="200" show-overflow-tooltip />
            <el-table-column prop="suggestion" :label="$t('auto.List.6dca6fb4')" min-width="200" show-overflow-tooltip />
            <el-table-column :label="$t('auto.List.3fea7ca7')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.is_handled ? 'success' : 'danger'">
                  {{ row.is_handled ? '已处理' : '未处理' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.2b6bc0f2')" width="120" fixed="right">
              <template #default="{ row }">
                <el-button v-if="!row.is_handled" size="small" type="primary" @click="handleWarning(row)">
                  标记已处理
                </el-button>
                <span v-else class="handled-text">{{ $t('auto.List.5ad6056a') }}</span>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.f8006830')" />
            </template>
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.List.73cbcf30')" name="config">
        <el-card>
          <el-form :model="configForm" label-width="170px" style="max-width: 640px">
            <el-divider content-position="left">{{ $t('auto.List.fa2d940e') }}</el-divider>
            <el-form-item :label="$t('auto.List.82c8e82a')">
              <el-input-number v-model="configForm.inactive_days_weight" :step="0.05" :min="0" :max="1" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.3ff41641')">
              <el-input-number v-model="configForm.purchase_freq_weight" :step="0.05" :min="0" :max="1" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.5a4d65e1')">
              <el-input-number v-model="configForm.order_value_weight" :step="0.05" :min="0" :max="1" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.4f00858f')">
              <el-input-number v-model="configForm.engagement_weight" :step="0.05" :min="0" :max="1" />
            </el-form-item>
            <el-divider content-position="left">{{ $t('auto.List.35a388ea') }}</el-divider>
            <el-form-item :label="$t('auto.List.a12f2b56')">
              <el-input-number v-model="configForm.inactive_threshold" :min="1" :max="365" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.0d33fcbc')">
              <el-input-number v-model="configForm.purchase_threshold" :min="1" :max="365" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.8771512f')">
              <el-input-number v-model="configForm.high_risk_score" :min="0" :max="100" />
            </el-form-item>
            <el-form-item :label="$t('auto.List.44844987')">
              <el-input-number v-model="configForm.critical_risk_score" :min="0" :max="100" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="savingConfig" @click="saveConfig">{{ $t('auto.List.ed752639') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    
    <el-dialog v-model="interveneDialogVisible" :title="$t('auto.List.defbd81f')" width="500px">
      <el-form :model="interveneForm" label-width="100px">
        <el-form-item :label="$t('auto.List.66a3c7a8')">
          <el-select v-model="interveneForm.method" style="width: 100%">
            <el-option :label="$t('auto.List.53285e0c')" value="coupon" />
            <el-option :label="$t('auto.List.1f9b1f07')" value="sms" />
            <el-option :label="$t('auto.List.997fd210')" value="manual" />
            <el-option :label="$t('auto.List.9a4e5d49')" value="push" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.34f5ab79')">
          <el-input v-model="interveneForm.content" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="interveneDialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="submittingIntervene" @click="confirmIntervene">{{ $t('auto.List.d97d5953') }}</el-button>
      </template>
    </el-dialog>

    
    <!-- v-loading 不能直接挂在 el-dialog 上：el-dialog 走 teleport 渲染、根节点不是单个
         元素，Vue 会打 "Runtime directive used on component with non-element root node.
         The directives will not function as intended."，指令也实际不生效。挂到弹窗体内
         的普通 div 上才是 Element Plus 的正确用法。 -->
    <el-dialog v-model="detailVisible" :title="$t('auto.List.214890a9')" width="700px">
      <div v-loading="detailLoading" style="min-height: 80px">
        <template v-if="currentUser">
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('auto.List.30acd20d')">{{ currentUser.user_id }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.25eb6ba3')">
              <el-tag :type="getRiskLevelType(currentUser.churn_risk)">{{ currentUser.churn_risk }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.7f2846ef')">{{ currentUser.last_activity_at }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.e79e7b21')">{{ currentUser.predicted_at }}</el-descriptions-item>
          </el-descriptions>

          <el-divider content-position="left">{{ $t('auto.List.8372f3ae') }}</el-divider>
          <div style="padding: 0 20px">
            <el-progress
              :percentage="Math.round(currentUser.churn_score || 0)"
              :color="getRiskColor(currentUser.churn_score)"
              :stroke-width="20"
            />
          </div>

          <el-divider content-position="left">{{ $t('auto.List.ba2a87b8') }}</el-divider>
          <el-alert
            v-if="currentUser.risk_factors"
            :title="currentUser.risk_factors"
            type="warning"
            :closable="false"
            show-icon
          />

          <template v-if="currentUser.riskFactors && currentUser.riskFactors.length">
            <el-divider content-position="left">{{ $t('auto.List.914b48f8') }}</el-divider>
            <el-table :data="currentUser.riskFactors" border size="small">
              <el-table-column prop="factor" :label="$t('auto.List.a7c4560b')" min-width="180" />
              <el-table-column prop="impact" :label="$t('auto.List.f6590a76')" width="120">
                <template #default="{ row: r }">
                  <el-tag ::type="$t('auto.List.ef044656')" size="small">
                    {{ r.impact }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="description" :label="$t('auto.List.f411d0f1')" min-width="200" show-overflow-tooltip />
            </el-table>
          </template>
        </template>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onActivated } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Operation, Refresh } from '@element-plus/icons-vue'
import {
  getChurnPredictions,
  getChurnPrediction,
  getChurnWarnings,
  getChurnModelConfig,
  saveChurnModelConfig,
  runChurnPrediction,
  interveneUser,
  markWarningHandled,
} from '@/api/churnPrediction.js'
import { toList } from '@/utils/list.js'

const { t } = useI18n()

const activeTab = ref('predictions')
const loadingPred = ref(false)
const loadingWarn = ref(false)
const predicting = ref(false)
const savingConfig = ref(false)
const submittingIntervene = ref(false)

const searchKeyword = ref('')
const users = ref([])
const warnings = ref([])
const stats = reactive({ highRisk: 0, mediumRisk: 0, lowRisk: 0, criticalRisk: 0 })

const configForm = reactive({
  id: 0,
  inactive_days_weight: 0.3,
  purchase_freq_weight: 0.3,
  order_value_weight: 0.2,
  engagement_weight: 0.2,
  inactive_threshold: 30,
  purchase_threshold: 60,
  high_risk_score: 70,
  critical_risk_score: 85,
})

const interveneDialogVisible = ref(false)
const interveneForm = reactive({ userId: 0, method: 'coupon', content: '' })
const detailVisible = ref(false)
const detailLoading = ref(false)
const currentUser = ref(null)

const filteredUsers = computed(() => {
  if (!searchKeyword.value) return users.value
  return users.value.filter((u) => (u.user_id || '').includes(searchKeyword.value))
})

const getRiskColor = (score) => {
  if (score >= 80) return '#EF4444'
  if (score >= 50) return '#F59E0B'
  return '#10B981'
}
const getRiskLevelType = (level) => {
  const map = { high: 'danger', critical: 'danger', medium: 'warning', low: 'success' }
  return map[(level || '').toLowerCase()] || 'info'
}
const warnLevelType = (level) => {
  const map = { high: 'danger', critical: 'danger', medium: 'warning', low: 'success' }
  return map[(level || '').toLowerCase()] || 'info'
}

const computeStats = (list) => {
  const counts = { high: 0, medium: 0, low: 0, critical: 0 }
  for (const u of list) {
    const level = (u.churn_risk || '').toLowerCase()
    if (level === 'high') counts.high++
    else if (level === 'critical') counts.critical++
    else if (level === 'medium') counts.medium++
    else if (level === 'low') counts.low++
  }
  stats.highRisk = counts.high
  stats.mediumRisk = counts.medium
  stats.lowRisk = counts.low
  stats.criticalRisk = counts.critical
}

const loadPredictions = async () => {
  loadingPred.value = true
  try {
    const res = await getChurnPredictions({ page: 1, page_size: 100 })
    users.value = toList(res?.data ?? res)
    computeStats(users.value)
  } catch (e) {
    ElMessage.error(t('auto.List.d13651b4') + (e && e.message ? e.message : ''))
  } finally {
    loadingPred.value = false
  }
}

const loadWarnings = async () => {
  loadingWarn.value = true
  try {
    const res = await getChurnWarnings({ page: 1, page_size: 100 })
    warnings.value = toList(res?.data ?? res)
  } catch (e) {
    ElMessage.error(t('auto.List.84109eeb') + (e && e.message ? e.message : ''))
  } finally {
    loadingWarn.value = false
  }
}

const loadConfig = async () => {
  try {
    const res = await getChurnModelConfig()
    const cfg = res?.data ?? res
    if (cfg && typeof cfg === 'object') {
      Object.assign(configForm, {
        id: cfg.id || 0,
        inactive_days_weight: cfg.inactive_days_weight ?? 0.3,
        purchase_freq_weight: cfg.purchase_freq_weight ?? 0.3,
        order_value_weight: cfg.order_value_weight ?? 0.2,
        engagement_weight: cfg.engagement_weight ?? 0.2,
        inactive_threshold: cfg.inactive_threshold ?? 30,
        purchase_threshold: cfg.purchase_threshold ?? 60,
        high_risk_score: cfg.high_risk_score ?? 70,
        critical_risk_score: cfg.critical_risk_score ?? 85,
      })
    }
  } catch (e) { console.warn("[request] 后台调用失败(已忽略):", e) }
}

const refreshAll = async () => {
  await Promise.all([loadPredictions(), loadWarnings(), loadConfig()])
}

const runPrediction = async () => {
  predicting.value = true
  try {
    await runChurnPrediction()
    ElMessage.success(t('auto.List.38db5ccb'))
    await refreshAll()
  } catch (e) {
    ElMessage.error(t('auto.List.bc20e050') + (e && e.message ? e.message : ''))
  } finally {
    predicting.value = false
  }
}

const intervene = (row) => {
  interveneForm.userId = row.user_id
  interveneForm.method = 'coupon'
  interveneForm.content = ''
  interveneDialogVisible.value = true
}

const confirmIntervene = async () => {
  submittingIntervene.value = true
  try {
    await interveneUser({
      warning_id: interveneForm.userId,
      intervention_type: interveneForm.method,
      note: interveneForm.content,
    })
    ElMessage.success(t('auto.List.02178cf5'))
    interveneDialogVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.List.625b6a84') + (e && e.message ? e.message : ''))
  } finally {
    submittingIntervene.value = false
  }
}

const viewUserDetail = async (row) => {
  detailVisible.value = true
  detailLoading.value = true
  try {
    const res = await getChurnPrediction({ user_id: row.user_id })
    let data = res?.data ?? res ?? {}
    if (typeof data.risk_factors === 'string' && data.risk_factors) {
      try {
        data.risk_factors = JSON.parse(data.risk_factors)
      } catch (e) {
        data.risk_factors = []
      }
    }
    currentUser.value = data
  } catch {
    currentUser.value = row
  } finally {
    detailLoading.value = false
  }
}

const handleWarning = async (row) => {
  try {
    const { value } = await ElMessageBox.prompt(t('auto.List.5dfa6db6'), t('auto.List.4247e294'), {
      confirmButtonText: t('auto.List.38cf16f2'),
      cancelButtonText: t('auto.List.625fb26b'),
      inputType: 'textarea',
    })
    await markWarningHandled(row.id, { note: value || '' })
    ElMessage.success(t('auto.List.193512d5'))
    await loadWarnings()
  } catch (e) {
    if (e !== 'cancel' && e && e.message) {
      ElMessage.error(t('auto.List.79585f92') + e.message)
    }
  }
}

const saveConfig = async () => {
  savingConfig.value = true
  try {
    await saveChurnModelConfig({ ...configForm })
    ElMessage.success(t('auto.List.e9472910'))
  } catch (e) {
    ElMessage.error(t('auto.List.78d235cc') + (e && e.message ? e.message : ''))
  } finally {
    savingConfig.value = false
  }
}

onMounted(() => refreshAll())
onActivated(() => refreshAll())
</script>

<style scoped lang="scss">
.churn-prediction-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.stats-row {
  margin-bottom: 20px;
  .stat-item {
    text-align: center;
    .stat-label { color: #909399; font-size: 14px; margin-bottom: 10px; }
    .stat-value { font-size: 28px; font-weight: bold; }
  }
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.handled-text { color: #909399; }
</style>
