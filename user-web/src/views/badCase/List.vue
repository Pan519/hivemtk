<template>
  <div class="bad-case">
    <el-alert
      v-if="unavailable"
      class="mb-12"
      type="warning"
      show-icon
      :closable="false"
      :title="$t('auto.List.1.2ed4cc')"
      :description="$t('auto.List.2.b79c1c')"
    />
    <el-alert
      v-if="!me"
      class="mb-12"
      type="info"
      show-icon
      :closable="false"
      :title="$t('auto.List.3.d5e452')"
      :description="$t('auto.List.4.86595b')"
    />

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.List.41.d6784e') }}</span>
          <div class="header-actions">
            <el-button
              :disabled="!me"
              @click="openCreate"
            >
              补录坏例
            </el-button>
            <el-tooltip
              :disabled="canExportSet"
              :content="$t('auto.List.5.2cbb4f')"
              placement="bottom"
            >
              <span>
                <el-button
                  type="primary"
                  :disabled="!canExportSet"
                  :loading="exporting"
                  @click="doExport"
                >
                  导出评测集
                </el-button>
              </span>
            </el-tooltip>
            <el-button
              :loading="loading"
              @click="refresh"
            >
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <div class="stats">
        <span class="stat-chip ratio">
          已判率<b>{{ ratioText }}</b>
        </span>
        <span
          v-for="s in statusOrder"
          :key="s"
          class="stat-chip"
        >
          {{ statusLabel(s) }} {{ statusCount(s) }}
        </span>
        <span
          v-for="l in layerOrder"
          :key="l"
          class="stat-chip"
        >
          {{ layerLabel(l) }} 归因 {{ layerCount(l) }}
        </span>
        <span
          v-if="unknownStatusText"
          class="stat-chip unknown"
        >{{ unknownStatusText }}</span>
      </div>

      <el-form
        inline
        class="mb-12"
      >
        <el-form-item :label="$t('auto.List.6.26ca20')">
          <el-select
            v-model="filter.source"
            class="filter-select"
            multiple
            clearable
            collapse-tags
            :placeholder="$t('auto.List.7.99b7a3')"
            @change="reload"
          >
            <el-option
              v-for="s in taxonomy.sources"
              :key="s"
              :label="sourceLabel(s)"
              :value="s"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.8.3fea7c')">
          <el-select
            v-model="filter.status"
            class="filter-select"
            multiple
            clearable
            collapse-tags
            :placeholder="$t('auto.List.9.1ea5af')"
            @change="reload"
          >
            <el-option
              v-for="s in taxonomy.statuses"
              :key="s"
              :label="statusLabel(s)"
              :value="s"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.10.edbe83')">
          <el-select
            v-model="filter.label"
            class="filter-select"
            multiple
            clearable
            collapse-tags
            :placeholder="$t('auto.List.11.a01449')"
            @change="reload"
          >
            <el-option
              v-for="l in taxonomy.labels"
              :key="l"
              :label="labelName(l)"
              :value="l"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.12.d8fc33')">
          <el-select
            v-model="filter.fix_layer"
            class="filter-select"
            clearable
            :placeholder="$t('auto.List.13.214840')"
            @change="reload"
          >
            <el-option
              v-for="l in taxonomy.fix_layers"
              :key="l"
              :label="layerLabel(l)"
              :value="l"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="list"
        stripe
      >
        <el-table-column
          :label="$t('auto.List.6.26ca20')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              size="small"
              :type="sourceTagType(row.source)"
            >
              {{ sourceLabel(row.source) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.14.498bd0')"
          min-width="260"
        >
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="openDetail(row)"
            >
              {{ row.query_text || row.session_id || row.id }}
            </el-button>
            <div
              v-if="row.mark_reason"
              class="reason"
            >
              {{ row.mark_reason }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.15.8f3599')"
          width="150"
        >
          <template #default="{ row }">
            <span :class="{ low: row.confidence < row.threshold }">
              {{ fixed(row.confidence) }} ／ 阈值 {{ fixed(row.threshold) }}
            </span>
            <div class="reason">
              检回 {{ row.retrieved_count }} 条
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.8.3fea7c')"
          width="100"
        >
          <template #default="{ row }">
            {{ statusLabel(row.status) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.16.ffda23')"
          width="200"
        >
          <template #default="{ row }">
            <div v-if="row.label">
              {{ labelName(row.label) }}
              <el-tag
                size="small"
                class="ml-6"
              >
                {{ layerLabel(row.fix_layer) }}
              </el-tag>
            </div>
            <span v-else>—</span>
            <div
              v-if="row.label_note"
              class="reason"
            >
              {{ row.label_note }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.17.badc15')"
          width="120"
        >
          <template #default="{ row }">
            {{ row.labeler_id || '—' }}
            <div
              v-if="row.labeled_at"
              class="reason"
            >
              {{ formatTime(row.labeled_at) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.18.2b6bc0')"
          width="180"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              v-for="act in rowActions(row, me)"
              :key="act"
              link
              :type="act === 'dismiss' ? 'danger' : 'primary'"
              @click="runAction(row, act)"
            >
              {{ actionLabel(act) }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="mt-12"
        layout="total, prev, pager, next"
        :current-page="filter.page"
        :page-size="filter.page_size"
        :total="total"
        @current-change="onPageChange"
      />
    </el-card>

    <el-drawer
      v-model="detail.visible"
      :title="$t('auto.List.19.dab392')"
      size="42%"
    >
      <el-skeleton
        v-if="detail.loading"
        :rows="5"
        animated
      />
      <el-alert
        v-else-if="detail.error"
        class="mb-12"
        type="error"
        show-icon
        :closable="false"
        :title="detail.error"
      />
      <pre
        v-else
        class="json-preview"
      >{{ detail.text }}</pre>
    </el-drawer>

    <el-drawer
      v-model="mark.visible"
      :title="$t('auto.List.20.0bc59d')"
      size="46%"
    >
      <el-alert
        v-if="mark.error"
        class="mb-12"
        type="error"
        show-icon
        :closable="false"
        :title="mark.error"
      />
      <el-descriptions
        :column="1"
        border
        class="mb-12"
      >
        <el-descriptions-item :label="$t('auto.List.21.9cacfc')">
          {{ mark.row?.query_text || '—' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.22.a98b82')">
          {{ mark.row?.answer_text || '—' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.23.49a47c')">
          {{ mark.row?.mark_reason || '—' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-form
        label-position="top"
        @submit.prevent
      >
        <el-form-item :label="$t('auto.List.10.edbe83')">
          <el-select
            v-model="mark.label"
            :placeholder="$t('auto.List.24.c9c71a')"
            class="full"
          >
            <el-option
              v-for="l in taxonomy.labels"
              :key="l"
              :label="`${labelName(l)}（${layerLabel(fixLayerOf(l, taxonomy))}）`"
              :value="l"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.25.cd64ff')">
          <el-input
            v-model="mark.note"
            type="textarea"
            :rows="3"
            maxlength="2000"
            show-word-limit
            :placeholder="$t('auto.List.26.fd3a48')"
          />
        </el-form-item>
      </el-form>

      <div class="decide-actions">
        <el-button
          type="primary"
          :disabled="!canSubmit"
          :loading="mark.submitting"
          @click="submitLabel"
        >
          提交判定
        </el-button>
      </div>
      <p
        v-if="!canSubmit && mark.label"
        class="hint"
      >
        判定依据必须写一句才能提交。
      </p>
    </el-drawer>

    <el-drawer
      v-model="create.visible"
      :title="$t('auto.List.27.045a60')"
      size="46%"
    >
      <el-alert
        v-if="create.error"
        class="mb-12"
        type="error"
        show-icon
        :closable="false"
        :title="create.error"
      />
      <el-form
        label-position="top"
        @submit.prevent
      >
        <el-form-item :label="$t('auto.List.28.a45b37')">
          <el-input
            v-model="create.form.session_id"
            :placeholder="$t('auto.List.29.ccfb57')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.30.2e4b80')">
          <el-input
            v-model="create.form.message_id"
            :placeholder="$t('auto.List.31.1e7b1d')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.32.764aca')">
          <el-input v-model="create.form.intent_type" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.21.9cacfc')">
          <el-input
            v-model="create.form.query_text"
            type="textarea"
            :rows="2"
            :placeholder="$t('auto.List.33.537b39')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.22.a98b82')">
          <el-input
            v-model="create.form.answer_text"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.List.34.14cdc5')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.35.1c3e6f')">
          <el-input
            v-model="create.form.reason"
            :placeholder="$t('auto.List.36.c6075e')"
          />
        </el-form-item>
      </el-form>
      <div class="decide-actions">
        <el-button
          type="primary"
          :loading="create.submitting"
          @click="submitCreate"
        >
          记进队列
        </el-button>
      </div>
    </el-drawer>

    <el-drawer
      v-model="exported.visible"
      :title="$t('auto.List.37.b85ec4')"
      size="42%"
    >
      <el-descriptions
        :column="1"
        border
        class="mb-12"
      >
        <el-descriptions-item :label="$t('auto.List.38.c0c78e')">
          {{ exported.set?.eval_set_id }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.39.a6f7f1')">
          {{ exported.set?.count }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.List.40.bdfd7a')">
          {{ formatTime(exported.set?.exported_at) }}
        </el-descriptions-item>
      </el-descriptions>
      <p class="hint">
        这 {{ exported.set?.count }} 条的状态已经从「已判定」改成「已导出」，不会再被下一次导出取走。
      </p>
      <div class="decide-actions">
        <el-button
          type="primary"
          @click="saveJson"
        >
          下载 JSON
        </el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { badCaseApi } from '@/api/badCase'
import { toList } from '@/utils/list'
import { useUserStore } from '@/stores/user'
import {
  OPEN_STATUSES,
  canExport,
  canSubmitLabel,
  fixLayerOf,
  hasDismissReason,
  rowActions
} from './actions'

// Bad Case 队列（G-2 / T-P8-03）。这一页是"答得不好的那些次"的账本：
// 谁被记进来了、判给了哪一层、判完有没有走出评测集。
//
// 与待办中心（views/approvalTask）不是一件事，也不共用数据出口：那一边读 /api/human-tasks
// （有截止时间的等人办的事），这一边读 /api/bad-cases（要归因的样本）。隔离由
// tests/unit/api_badCase.test.js 双向钉住。
//
// 值域表（来源/状态/类目/责任层）从 GET /taxonomy 取，前端不抄第二张表 ——
// 抄一份的后果是库里加一个类目后页面"选不到它"，而这条链的下游是评测集，
// 选不到等于统计里少一层，且少的那层没人知道。中文文案（下面那几个 LABELS 表）是
// **显示用**的：值域来自 taxonomy，命中就译、没命中就原样显示英文枚举值。

const SOURCE_LABELS = {
  low_confidence: t('auto.List.42.0c4305'),
  zero_hit: t('auto.List.43.8732b9'),
  manual: t('auto.List.44.9fcf7b')
}
const SOURCE_TAG_TYPES = {
  low_confidence: 'warning',
  zero_hit: 'danger',
  manual: 'info'
}
const STATUS_LABELS = {
  pending: t('auto.List.45.5a7f3f'),
  labeled: t('auto.List.46.f36e8a'),
  exported: t('auto.List.47.5c51c3'),
  dismissed: t('auto.List.48.50239f')
}
const LABEL_NAMES = {
  kb_missing: t('auto.List.49.27de62'),
  kb_stale: t('auto.List.50.2167b7'),
  retrieve_miss: t('auto.List.51.df4c6e'),
  retrieve_wrong: t('auto.List.52.5ae0bd'),
  generation_wrong: t('auto.List.53.e739d3'),
  generation_style: t('auto.List.54.626c87'),
  intent_misjudge: t('auto.List.55.811e60')
}
const LAYER_LABELS = {
  knowledge: t('auto.List.56.7e4333'),
  retrieval: t('auto.List.57.8bedce'),
  generation: t('auto.List.58.4dfe70'),
  intent: t('auto.List.59.ef84bd')
}
const ACTION_LABELS = { label: t('auto.List.60.567503'), dismiss: t('auto.List.61.bd9fcf') }

// STATS_POLL_MS：已判率是本卡北极星，值班要的是"打开时看到的是此刻的数"。
// 30s 与待办中心同档；更短只是把读接口打成压力源。
const STATS_POLL_MS = 30000

const userStore = useUserStore()
// 后端的 labeler_id 存的是 user id 的十进制字符串，这里统一成字符串，
// 否则 userInfo.id 是数字时"这条是不是我判的"会永远判否。
const me = computed(() => String(userStore.userInfo?.id ?? ''))

const { t } = useI18n()

const loading = ref(false)
const exporting = ref(false)
const unavailable = ref(false)
const list = ref([])
const total = ref(0)

const stats = reactive({ by_status: {}, by_fix_layer: {}, labeled_ratio: 0 })
// taxonomy 的初值是空数组而不是猜一份值域：没取到值域时下拉框是空的（可见的缺），
// 而不是默默替后端"补"了一份它并不认的清单。
const taxonomy = reactive({ sources: [], statuses: [], labels: [], fix_layers: [], label_layer: {} })

const filter = reactive({
  source: [],
  status: [...OPEN_STATUSES],
  label: [],
  fix_layer: '',
  page: 1,
  page_size: 20
})

const detail = reactive({ visible: false, loading: false, error: '', text: '' })
const mark = reactive({ visible: false, submitting: false, error: '', row: null, label: '', note: '' })
const create = reactive({
  visible: false,
  submitting: false,
  error: '',
  form: { session_id: '', message_id: '', intent_type: '', query_text: '', answer_text: '', reason: '' }
})
const exported = reactive({ visible: false, set: null })

const sourceLabel = (s) => SOURCE_LABELS[s] || s || t('auto.List.62.36cead')
const sourceTagType = (s) => SOURCE_TAG_TYPES[s] || 'info'
const statusLabel = (s) => STATUS_LABELS[s] || s || t('auto.List.63.903b25')
const labelName = (l) => LABEL_NAMES[l] || l
const layerLabel = (l) => LAYER_LABELS[l] || l || t('auto.List.64.23ac2b')
const actionLabel = (a) => ACTION_LABELS[a] || a
const statusCount = (s) => Number(stats.by_status?.[s]) || 0
const layerCount = (l) => Number(stats.by_fix_layer?.[l]) || 0

// 读数条里的状态档：值域表在位时按它给顺序（后端那份，前端不抄第二张表）；
// 没在位时退回 by_status 自己的键并原样显示 —— 两条路都不许把任何一档藏起来。
const statusOrder = computed(() =>
  taxonomy.statuses.length ? taxonomy.statuses : Object.keys(stats.by_status || {})
)
const layerOrder = computed(() => taxonomy.fix_layers || [])

const ratioText = computed(() => `${(stats.labeled_ratio * 100).toFixed(1)}%`)

// unknownStatusText：值域外的状态行**必须自己站出来**。
// 只印 by_status 里已知的几档，等于把脏数据从读数条上抹掉 ——
// 而"队列里有 3 条状态是 pending_"这种事实，恰恰是这一页最该看见的东西。
// 值域表还没回来时不作答（那句"值域外"此刻没有依据，说出来的可能是错的）。
const unknownStatusText = computed(() => {
  if (!taxonomy.statuses.length) return ''
  const known = new Set(taxonomy.statuses)
  const rows = Object.entries(stats.by_status || {}).filter(
    ([s, n]) => Number(n) > 0 && !known.has(s)
  )
  if (!rows.length) return ''
  return `值域外状态 ${rows.map(([s, n]) => `${s || '(空)'} ${n}`).join('、')}（不进上面任何一档）`
})

const canExportSet = computed(() => canExport(stats))
const canSubmit = computed(() => canSubmitLabel(mark.label, mark.note))

function formatTime(value) {
  if (!value) return '—'
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? String(value) : d.toLocaleString()
}

function fixed(v) {
  const n = Number(v)
  return Number.isFinite(n) ? n.toFixed(3) : '—'
}

function queryParams() {
  const q = { page: filter.page, page_size: filter.page_size }
  if (filter.source.length) q.source = filter.source
  // 空 status 不进串 ⇒ 由后端兜它的默认视图（只列 pending）。所以"把状态筛选清空"
  // 看到的是待判定那一本，不是全量；想同时看见已判定待导出的，得把那一档勾回来。
  if (filter.status.length) q.status = filter.status
  if (filter.label.length) q.label = filter.label
  if (filter.fix_layer) q.fix_layer = filter.fix_layer
  return q
}

async function fetchList() {
  loading.value = true
  try {
    const data = await badCaseApi.list(queryParams())
    list.value = toList(data)
    total.value = Number(data?.total) || 0
    unavailable.value = false
  } catch (err) {
    // 503 与"这次查询失败"要分开：前者是配置状态，页面要长期挂着说一句话；
    // 后者刷新就好。这里绝不清成 [] 假装读到了。
    // 400（筛选器里有值域外的值）也留在上一轮列表上：那是句业务结论，
    // 但把表格清空会让人以为"这个类目下一条都没有"。
    unavailable.value = err?.status === 503
  } finally {
    loading.value = false
  }
}

async function fetchStats() {
  try {
    const data = await badCaseApi.stats()
    if (!data) return
    stats.by_status = data.by_status || {}
    stats.by_fix_layer = data.by_fix_layer || {}
    stats.labeled_ratio = Number(data.labeled_ratio) || 0
    unavailable.value = false
  } catch (err) {
    // 静默失败：上一轮读数留着。轮询每 30s 一次，弹一句"服务异常"只会变成噪音。
    if (err?.status === 503) unavailable.value = true
  }
}

async function fetchTaxonomy() {
  try {
    const data = await badCaseApi.taxonomy()
    if (!data) return
    taxonomy.sources = Array.isArray(data.sources) ? data.sources : []
    taxonomy.statuses = Array.isArray(data.statuses) ? data.statuses : []
    taxonomy.labels = Array.isArray(data.labels) ? data.labels : []
    taxonomy.fix_layers = Array.isArray(data.fix_layers) ? data.fix_layers : []
    taxonomy.label_layer = data.label_layer && typeof data.label_layer === 'object' ? data.label_layer : {}
  } catch {
    /* 静默：下拉框退回空清单，页面上是"选不到"而不是"猜了一份" */
  }
}

function reload() {
  filter.page = 1
  fetchList()
}

function refresh() {
  fetchList()
  fetchStats()
}

function onPageChange(page) {
  filter.page = page
  fetchList()
}

async function openDetail(row) {
  detail.visible = true
  detail.loading = true
  detail.error = ''
  detail.text = ''
  try {
    // 点标题读的是**服务端那一行**，不是表格里这一份：判完再点一次要看见的是落库结果，
    // 否则"页面显示已判定"可能只是本地那行对象还没被覆盖。
    const data = await badCaseApi.get(row.id)
    detail.text = JSON.stringify(data ?? row, null, 2)
  } catch (err) {
    detail.error = err?.message || t('auto.List.65.f87e6e')
  } finally {
    detail.loading = false
  }
}

function openMark(row) {
  mark.visible = true
  mark.row = row
  mark.label = row.label || ''
  mark.note = row.label_note || ''
  mark.error = ''
}

async function submitLabel() {
  mark.submitting = true
  mark.error = ''
  try {
    await badCaseApi.label(mark.row.id, mark.label, mark.note)
    ElMessage.success(t('auto.List.66.086beb'))
    mark.visible = false
    refresh()
  } catch (err) {
    const status = Number(err?.status) || 0
    if (status === 409) {
      // 两个人同时开着这一条：后提交的那个必被 CAS 挡下。这句"当前状态不允许"要让人读到，
      // 并把队列刷回来 —— 别人的结论该出现在界面上，而不是让人再点一次替自己盖章。
      ElMessage.warning(err.message || t('auto.List.67.236864'))
      mark.visible = false
      refresh()
    } else {
      // 其余一律留在抽屉里说一句话：400 是"类目或依据这条腿没齐"（界面已经替人挡过一轮），
      // 401 是会话在这期间过期了，5xx 是底座的事 —— 三种都不该表现成
      // "点了没反应，刷新看看"，那会让人以为是自己手慢。
      mark.error = err.message || t('auto.List.68.a669be')
    }
  } finally {
    mark.submitting = false
  }
}

async function dismissRow(row) {
  let value
  try {
    ({ value } = await ElMessageBox.prompt(t('auto.List.69.be5cf5'), t('auto.List.70.895b38'), {
      inputType: 'textarea',
      inputPlaceholder: t('auto.List.71.6e1e4a'),
      inputValidator: (v) => (hasDismissReason(v) ? true : t('auto.List.72.9c365d')),
      confirmButtonText: t('auto.List.61.bd9fcf'),
      cancelButtonText: t('auto.List.73.5f4112')
    }))
  } catch {
    return // 点了返回
  }
  try {
    await badCaseApi.dismiss(row.id, value)
    ElMessage.success(t('auto.List.48.50239f'))
    refresh()
  } catch (err) {
    if (Number(err?.status) === 409) {
      // 已经被别人判掉了：撤销这条路此刻不存在。刷新后看到的是当前那一行，不是猜。
      ElMessage.warning(err.message || t('auto.List.74.ca87dc'))
      refresh()
    } else {
      ElMessage.error(err.message || t('auto.List.75.a8301a'))
    }
  }
}

function runAction(row, action) {
  if (action === 'label') openMark(row)
  else if (action === 'dismiss') dismissRow(row)
}

function openCreate() {
  create.visible = true
  create.error = ''
}

async function submitCreate() {
  create.submitting = true
  create.error = ''
  try {
    await badCaseApi.create({ ...create.form })
    ElMessage.success(t('auto.List.76.a6d2ff'))
    create.visible = false
    refresh()
  } catch (err) {
    // 409 = 同一个 message_id 已经记过（幂等键撞上）。那句"已存在"是结果而不是失败，
    // 但在这里必须让人看见 —— 否则坐席会以为补录没成功而再点一次。
    create.error = err?.message || t('auto.List.77.6af607')
  }
  create.submitting = false
}

async function doExport() {
  exporting.value = true
  try {
    // 筛选器里选了类目就按选的那几类导出，没选则导出全部"已判定未导出"的行。
    // limit 不给：后端对"没给"兜默认档，一次点击取空整池不是这一页该替人做的决定。
    const set = await badCaseApi.exportSet(filter.label, undefined)
    exported.set = set
    exported.visible = true
    refresh()
  } catch (err) {
    const status = Number(err?.status) || 0
    // 这一竖的调用是静默的（导出按钮点下去不该由网关弹一句"服务异常"），
    // 所以三种结果都得由这里说话 —— 尤其 409：那句"没有可导出的样本"是业务结论，
    // 说成失败会让人再点一次，说成没事会让人以为已经导走了。
    if (status === 409) ElMessage.warning(t('auto.List.78.428ab6'))
    else ElMessage.error(err.message || t('auto.List.79.dd51ab'))
  } finally {
    exporting.value = false
  }
}

function saveJson() {
  if (!exported.set) return
  const blob = new Blob([JSON.stringify(exported.set, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${exported.set.eval_set_id || 'eval-set'}.json`
  a.click()
  // revoke 在 click 之后立刻调用是安全的（下载已经交给浏览器），不留着句柄等于泄漏一个 Blob。
  URL.revokeObjectURL(url)
}

let pollTimer = null

onMounted(() => {
  fetchTaxonomy()
  refresh()
  pollTimer = setInterval(fetchStats, STATS_POLL_MS)
})

onBeforeUnmount(() => {
  if (pollTimer) clearInterval(pollTimer)
  pollTimer = null
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.mb-12 {
  margin-bottom: 12px;
}
.mt-12 {
  margin-top: 12px;
}
.ml-6 {
  margin-left: 6px;
}
.full {
  width: 100%;
}
.filter-select {
  width: 220px;
}
.stats {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  font-size: 13px;
}
.stat-chip {
  white-space: nowrap;
}
.stat-chip.ratio b {
  margin-left: 4px;
}
.stat-chip.unknown {
  color: #f56c6c;
}
.low {
  color: #f56c6c;
}
.reason {
  color: #909399;
  font-size: 12px;
}
.hint {
  color: #909399;
  font-size: 13px;
}
.decide-actions {
  display: flex;
  gap: 8px;
}
.json-preview {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 6px;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
