<template>
  <div class="ltc-config-page">
    <el-alert
      v-if="cfg.degraded"
      class="degraded-alert"
      type="error"
      :closable="false"
      show-icon
      :title="$t('auto.LTCConfig.bb2d62c4')"
    >
      <div class="alert-body">
        <p>
          底层原因：{{ cfg.degrade_reason || '（后端未给出）' }}
        </p>
        <p class="alert-note">
          这不是"运营把开关关了"。存储读不动时闸门按 fail-closed 拦着所有 LTC 路由，
          先把存储修好再回来判读这份表单。保存一次读不到内容的表单会把库里那一份盖掉。
        </p>
      </div>
    </el-alert>

    <el-row :gutter="16">
      <el-col
        :xs="24"
        :md="8"
      >
        <el-card
          shadow="never"
          class="master-card"
        >
          <template #header>
            <div class="card-head">
              <span>{{ $t('auto.LTCConfig.b5db5cef') }}</span>
              <el-tag
                :type="cfg.source === 'kv' ? 'success' : 'info'"
                size="small"
              >
                来源：{{ sourceText }}
              </el-tag>
            </div>
          </template>
          <div class="master-body">
            <el-switch
              v-model="cfg.enabled"
              :active-text="$t('auto.LTCConfig.84932056')"
              :inactive-text="$t('auto.LTCConfig.d58a55bc')"
            />
            <p class="hint">
              关掉时六个阶段全部拦在外面，无论各自那一项配的是什么。
              这一档与分阶段档是两道独立的锁。
            </p>
          </div>
        </el-card>
      </el-col>

      <el-col
        :xs="24"
        :md="16"
      >
        <el-card shadow="never">
          <template #header>
            <div class="card-head">
              <span>{{ $t('auto.LTCConfig.bf242e90') }}</span>
              <span class="head-sub">归管路由合计 {{ cfg.guarded_total }} 条</span>
            </div>
          </template>
          <el-table
            :data="stageRows"
            size="small"
          >
            <el-table-column
              :label="$t('auto.LTCConfig.d316330b')"
              width="120"
            >
              <template #default="{ row }">
                {{ stageLabel(row.stage) }}
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('auto.LTCConfig.a6beb974')"
              width="90"
            >
              <template #default="{ row }">
                <el-switch v-model="cfg.stages_enabled[row.stage]" />
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('auto.LTCConfig.1d84af3b')"
              width="150"
            >
              <template #default="{ row }">
                <el-tag
                  :type="row.active ? 'success' : 'danger'"
                  size="small"
                >
                  {{ reasonText(row.reason) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('auto.LTCConfig.cf63954d')"
              width="110"
            >
              <template #default="{ row }">
                {{ row.guarded }} 条
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.LTCConfig.36bf4729')">
              <template #default="{ row }">
                <code>{{ row.stage }}</code>
              </template>
            </el-table-column>
          </el-table>
          <p class="hint">
            打开某一项只影响它自己那一段。归管路由为 0 表示这个阶段还没有任何业务路由
            挂到闸门上 —— 此刻这一项配成开或配成关，现网行为一模一样。
          </p>
        </el-card>
      </el-col>
    </el-row>

    <el-card
      v-if="boundsReady"
      shadow="never"
      class="threshold-card"
    >
      <template #header>
        <div class="card-head">
          <span>{{ $t('auto.LTCConfig.e0de1c9a') }}</span>
          <span class="head-sub">{{ $t('auto.LTCConfig.fc10abed') }}</span>
        </div>
      </template>
      <el-form
        label-width="150px"
        label-position="left"
      >
        <el-row :gutter="16">
          <el-col
            v-for="key in thresholdKeys"
            :key="key"
            :xs="24"
            :md="12"
          >
            <el-form-item :label="thresholdLabel(key)">
              <el-input-number
                v-model="cfg.thresholds[key]"
                :min="inputMin(key)"
                :max="maxOf(key)"
                :step="stepOf(key)"
                :precision="key === 'lead_score' || key === 'discount_percent' ? 0 : 3"
                controls-position="right"
              />
              <div class="bound-text">
                <span>{{ bounds[key].desc }}</span>
                <span
                  v-if="!bounds[key].zero_legal"
                  class="zero-bad"
                >{{ $t('auto.LTCConfig.03be9251') }}</span>
                <span
                  v-else
                  class="zero-ok"
                >{{ $t('auto.LTCConfig.82c82bda') }}</span>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <el-card
      v-else
      shadow="never"
      class="threshold-card"
    >
      <template #header>
        <div class="card-head">
          <span>{{ $t('auto.LTCConfig.e0de1c9a') }}</span>
        </div>
      </template>
      <p class="hint">
        阈值输入框要等后端的取值边界到位才能渲染（边界与校验同源，写在页面里就是第二份事实）。
        点「重新读取」再试一次；一直读不到就是存储或权限的问题，别凭记忆填。
      </p>
    </el-card>

    <el-card
      v-if="cfg.reading_hints && cfg.reading_hints.length"
      shadow="never"
      class="hints-card"
    >
      <template #header>
        <div class="card-head">
          <span>{{ $t('auto.LTCConfig.04d6d7e9') }}</span>
        </div>
      </template>
      <ul class="hints">
        <li
          v-for="(h, i) in cfg.reading_hints"
          :key="i"
        >
          {{ h }}
        </li>
      </ul>
    </el-card>

    <el-alert
      v-if="saveError"
      class="save-error"
      type="error"
      :closable="false"
      show-icon
      :title="saveError"
    />
    <el-alert
      v-if="saveNotice.audit_error"
      class="save-error"
      type="warning"
      :closable="false"
      show-icon
      ::title="$t('auto.LTCConfig.865f2617')"
    />

    <div class="actions">
      <el-button
        type="primary"
        :loading="saving"
        @click="onSave"
      >
        保存整份策略
      </el-button>
      <el-button
        :loading="loading"
        @click="load"
      >
        重新读取
      </el-button>
      <span class="actions-meta">
        KV 键 <code>{{ cfg.kv_key }}</code> · 上限 {{ cfg.max_bytes }} 字节 ·
        缓存 TTL {{ cfg.cache_ttl_seconds }}s（保存后立即失效，不必等）
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getLTCConfig, saveLTCConfig } from '@/api/ltcConfig'

// 阶段键与中文名：顺序照后端 LTCKnownStages（销售流程顺序），
// 键名照 ltc.config 的 JSON 字面量 —— 这两处错一处，保存就会被后端拒收。
const STAGES = [
  { stage: 'opportunity', label: t('auto.LTCConfig.55290e55') },
  { stage: 'outreach', label: t('auto.LTCConfig.301dd54b') },
  { stage: 'quote', label: t('auto.LTCConfig.00bae3b9') },
  { stage: 'bill', label: t('auto.LTCConfig.a08c6e69') },
  { stage: 'payment', label: t('auto.LTCConfig.1755fc2a') },
  { stage: 'collection', label: t('auto.LTCConfig.b01bde53') }
]

const THRESHOLDS = [
  { key: 'lead_score', label: t('auto.LTCConfig.fb1e8b62') },
  { key: 'confidence', label: t('auto.LTCConfig.94fa8ea9') },
  { key: 'discount_percent', label: t('auto.LTCConfig.42cf7fb6') },
  { key: 'win_probability', label: t('auto.LTCConfig.7ff8686a') }
]

const REASON_TEXT = {
  active: t('auto.LTCConfig.45639999'),
  master_off: t('auto.LTCConfig.ea6e3329'),
  stage_off: t('auto.LTCConfig.9bff6a01'),
  degraded: t('auto.LTCConfig.3d8283fb'),
  unknown_stage: t('auto.LTCConfig.2a1bfefb')
}

const emptyCfg = () => ({
  enabled: false,
  stages_enabled: { opportunity: false, outreach: false, quote: false, bill: false, payment: false, collection: false },
  thresholds: { lead_score: 70, confidence: 0.8, discount_percent: 15, win_probability: 0.5 },
  stage_status: STAGES.map((s) => ({ stage: s.stage, active: false, reason: 'master_off' })),
  guarded_routes: {},
  guarded_total: 0,
  reading_hints: [],
  source: 'default',
  degraded: false,
  degrade_reason: '',
  kv_key: 'ltc.config',
  max_bytes: 8192,
  cache_ttl_seconds: 60,
  threshold_ranges: {},
  threshold_bounds: {}
})

const { t } = useI18n()

const cfg = reactive(emptyCfg())
const bounds = ref({})
const thresholdKeys = THRESHOLDS.map((t) => t.key)
const loading = ref(false)
const saving = ref(false)
const saveError = ref('')
const saveNotice = ref({})

const sourceText = computed(() => (cfg.source === 'kv' ? t('auto.LTCConfig.56f7a2fe') : t('auto.LTCConfig.c7987057')))

// 边界没到位之前不渲染阈值区：首帧 bounds 还是空对象，直接取 .max 会让整页白屏。
const boundsReady = computed(() => thresholdKeys.every((k) => !!bounds.value[k]))

const stageRows = computed(() =>
  (cfg.stage_status || []).map((row) => ({ ...row, guarded: cfg.guarded_routes?.[row.stage] ?? 0 }))
)

const stageLabel = (s) => STAGES.find((x) => x.stage === s)?.label || s
const reasonText = (r) => REASON_TEXT[r] || r
const thresholdLabel = (k) => THRESHOLDS.find((t) => t.key === k)?.label || k

// 输入边界取自后端同一份表：写在视图里的第二份范围，迟早有一天和校验各说各话。
const inputMin = (k) => (bounds.value[k]?.zero_legal ? 0 : bounds.value[k]?.min)
const maxOf = (k) => bounds.value[k]?.max
const stepOf = (k) => (k === 'lead_score' || k === 'discount_percent' ? 1 : 0.01)

async function load() {
  loading.value = true
  saveError.value = ''
  try {
    const data = await getLTCConfig()
    Object.assign(cfg, emptyCfg(), data)
    bounds.value = data.threshold_bounds || {}
  } catch (err) {
    saveError.value = t('auto.LTCConfig.4154d0e7') + (err?.message || t('auto.LTCConfig.974e7484'))
  } finally {
    loading.value = false
  }
}

function payload() {
  const stages = {}
  for (const s of STAGES) stages[s.stage] = !!cfg.stages_enabled[s.stage]
  const thresholds = {}
  for (const k of thresholdKeys) thresholds[k] = cfg.thresholds[k]
  return { enabled: !!cfg.enabled, stages_enabled: stages, thresholds }
}

async function onSave() {
  saving.value = true
  saveError.value = ''
  saveNotice.value = {}
  try {
    const res = await saveLTCConfig(payload())
    saveNotice.value = res || {}
    if (res?.effective) {
      Object.assign(cfg, emptyCfg(), res.effective)
      bounds.value = res.effective.threshold_bounds || {}
    }
    if (res?.persisted && !res?.audit_error) {
      ElMessage.success(`已保存（${res.stored_bytes} 字节，开启阶段：${(res.stages_on || []).join('、') || t('auto.LTCConfig.d81bb206')}）`)
    }
  } catch (err) {
    // 后端把拒绝原因拼到了句子级（「阶段名 qoute 不认识」「lead_score 给 0 等于拆闸」…），
    // 原样铺出来，不要换成一句"参数错误"让运营在六个阶段名里猜。
    saveError.value = err?.message || t('auto.LTCConfig.6de920b4')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.ltc-config-page {
  padding: 16px;
}
.ltc-config-page > .el-row,
.ltc-config-page > .el-card,
.ltc-config-page > .el-alert {
  margin-bottom: 16px;
}
.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.head-sub {
  font-size: 12px;
  color: #909399;
}
.master-body .hint,
.hint {
  margin: 12px 0 0;
  font-size: 12px;
  line-height: 1.7;
  color: #909399;
}
.bound-text {
  display: flex;
  flex-direction: column;
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
.zero-bad {
  color: #f56c6c;
}
.zero-ok {
  color: #67c23a;
}
.hints {
  margin: 0;
  padding-left: 18px;
  font-size: 13px;
  line-height: 1.9;
}
.actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
.actions-meta {
  font-size: 12px;
  color: #909399;
}
.alert-body p {
  margin: 4px 0;
}
.alert-note {
  font-size: 12px;
  color: #606266;
}
code {
  background: #f5f7fa;
  padding: 0 4px;
  border-radius: 3px;
}
</style>
