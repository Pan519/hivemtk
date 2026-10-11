<template>
  <div class="approval-task">
    <el-alert
      v-if="unavailable"
      class="mb-12"
      type="warning"
      show-icon
      :closable="false"
      :title="$t('auto.List.fdbee843')"
      :description="$t('auto.List.ad96af55')"
    />
    <el-alert
      v-if="!me"
      class="mb-12"
      type="info"
      show-icon
      :closable="false"
      :title="$t('auto.List.d5e452f2')"
      :description="$t('auto.List.978472f6')"
    />

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.List.ba75064b') }}</span>
          <el-button
            :loading="loading"
            @click="refresh"
          >
            刷新
          </el-button>
        </div>
      </template>

      <div class="counts">
        <span
          v-for="row in counts.by_kind"
          :key="row.kind"
          class="count-chip"
        >
          {{ kindLabel(row.kind) }} 开放 {{ row.open }} ／<b :class="{ overdue: row.overdue > 0 }">逾期 {{ row.overdue }}</b>
        </span>
        <span
          v-if="counts.unknown_kind_open"
          class="count-chip"
        >
          未知类型开放 {{ counts.unknown_kind_open }}（不进上面任何一档的逾期读数）
        </span>
        <span
          v-if="counts.at"
          class="counts-at"
        >读数时刻 {{ formatTime(counts.at) }}</span>
      </div>

      <el-form
        inline
        class="mb-12"
      >
        <el-form-item :label="$t('auto.List.226b0912')">
          <el-select
            v-model="filter.kind"
            class="filter-select"
            multiple
            clearable
            collapse-tags
            :placeholder="$t('auto.List.c079e7d5')"
            @change="reload"
          >
            <el-option
              v-for="k in KINDS"
              :key="k"
              :label="kindLabel(k)"
              :value="k"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.3fea7ca7')">
          <el-select
            v-model="filter.status"
            class="filter-select"
            multiple
            clearable
            collapse-tags
            :placeholder="$t('auto.List.5f17fb92')"
            @change="reload"
          >
            <el-option
              v-for="s in STATUSES"
              :key="s"
              :label="statusLabel(s)"
              :value="s"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.2296bd83')">
          <el-select
            v-model="filter.assignee"
            class="filter-select"
            @change="reload"
          >
            <el-option
              :label="$t('auto.List.04d85c53')"
              value="all"
            />
            <el-option
              :label="$t('auto.List.b671b7db')"
              value="me"
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
          :label="$t('auto.List.226b0912')"
          width="110"
        >
          <template #default="{ row }">
            <el-tag
              size="small"
              :type="kindTagType(row.kind)"
            >
              {{ kindLabel(row.kind) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.32c65d8d')"
          min-width="240"
        >
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="openDetail(row)"
            >
              {{ row.title || row.subject_id }}
            </el-button>
            <div
              v-if="row.reason"
              class="reason"
            >
              {{ row.reason }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.3fea7ca7')"
          width="100"
        >
          <template #default="{ row }">
            {{ statusLabel(row.status) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.b5d396a8')"
          width="100"
        >
          <template #default="{ row }">
            {{ row.assignee_user_id || '—' }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.cd8fbf6f')"
          width="170"
        >
          <template #default="{ row }">
            <span :class="{ overdue: isOverdue(row) }">{{ slaText(row) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.List.2b6bc0f2')"
          width="240"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              v-for="act in rowActions(row, me)"
              :key="act"
              link
              :type="act === 'cancel' ? 'danger' : 'primary'"
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
      :title="$t('auto.List.a152e111')"
      size="42%"
    >
      <pre class="json-preview">{{ detail.text }}</pre>
    </el-drawer>

    <el-drawer
      v-model="decide.visible"
      :title="$t('auto.List.e722a3f5')"
      size="46%"
    >
      <el-alert
        v-if="decide.error"
        class="mb-12"
        type="error"
        show-icon
        :closable="false"
        :title="decide.error"
      />
      <el-skeleton
        v-if="decide.loading"
        :rows="5"
        animated
      />

      <template v-else>
        <el-descriptions
          :column="1"
          border
          class="mb-12"
        >
          <el-descriptions-item :label="$t('auto.List.b60ec8b0')">
            {{ decide.task?.title || '—' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.84a2c6a2')">
            {{ decide.approvalId }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.66914536')">
            {{ decide.approval?.policy_key || '—' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.5249b6de')">
            {{ subjectText(decide.approval) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.d352aeb6')">
            {{ decide.approval?.status || '—' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.cd8fbf6f')">
            {{ formatTime(decide.approval?.expires_at) }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="decide.approval?.decided_by"
            :label="$t('auto.List.d43c7ac9')"
          >
            {{ decide.approval.decided_by }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="decide.approval?.decision_note"
            :label="$t('auto.List.fe9803f2')"
          >
            {{ decide.approval.decision_note }}
          </el-descriptions-item>
        </el-descriptions>

        <el-form
          v-if="verdicts.length"
          label-position="top"
        >
          <el-form-item :label="$t('auto.List.fe9803f2')">
            <el-input
              v-model="decide.note"
              type="textarea"
              :rows="3"
              maxlength="2000"
              show-word-limit
              :placeholder="$t('auto.List.850b2750')"
            />
          </el-form-item>
        </el-form>
        <p
          v-else-if="decide.approval"
          class="hint"
        >
          <!-- 两句话不能合成一句：「时限已过」时这条既没落定也没被清扫收口，
               而处置方向是去修流程的 TTL；括号里那句说的是另一件事。 -->
          {{ decideHint }}
        </p>

        <div class="decide-actions">
          <el-button
            v-for="v in verdicts"
            :key="v"
            :type="v === 'approved' ? 'success' : 'danger'"
            :loading="decide.submitting"
            :disabled="v === 'rejected' && !hasRejectReason(decide.note)"
            @click="submitDecide(v)"
          >
            {{ verdictLabel(v) }}
          </el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import i18n from '@/i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { humanTaskApi } from '@/api/humanTask'
import { approvalApi } from '@/api/approval'
import { toList } from '@/utils/list'
import { useUserStore } from '@/stores/user'
import {
  KIND_APPROVAL,
  KIND_COLLECTION,
  KIND_HANDOFF,
  hasRejectReason,
  rowActions,
  slaField,
  verdictBlockReason,
  verdictButtons
} from './actions'

// 待办中心（N-9）。三类待办一张表：会话转人工 / 待审批 / 催收升级。
//
// 与坐席收件箱（views/inbox）的关系是**分离视图、不分离数据出口**之外的另一半：
// 收件箱读的是会话（/api/inbox/conversations），这一页读的是"等人办的事"
// （/api/human-tasks），两边谁也不转调谁 —— 隔离由 tests/unit/api_humanTask.test.js 双向钉住。
//
// 通知不在这里发：全站通知出口只有一个（views/Notifications.vue），待办页再挂一个铃铛
// 就会有两套已读状态。这一页要的是"打开时看到的是此刻的数"，所以未读数走轮询而不是推送。

const KINDS = [KIND_HANDOFF, KIND_APPROVAL, KIND_COLLECTION]
const STATUSES = ['pending', 'claimed', 'done', 'cancelled']
const OPEN_STATUSES = ['pending', 'claimed']

const KIND_LABELS = {
  [KIND_HANDOFF]: i18n.global.t('auto.List.e8c16d26'),
  [KIND_APPROVAL]: i18n.global.t('auto.List.b0bf01a4'),
  [KIND_COLLECTION]: i18n.global.t('auto.List.7f180af3')
}
const KIND_TAG_TYPES = {
  [KIND_HANDOFF]: 'warning',
  [KIND_APPROVAL]: 'danger',
  [KIND_COLLECTION]: 'info'
}
const STATUS_LABELS = {
  pending: i18n.global.t('auto.List.047109de'),
  claimed: i18n.global.t('auto.List.b32d7c1e'),
  done: i18n.global.t('auto.List.fad5222c'),
  cancelled: i18n.global.t('auto.List.50239f4f')
}
const ACTION_LABELS = {
  claim: i18n.global.t('auto.List.54b931f9'),
  release: i18n.global.t('auto.List.cbdcabad'),
  complete: i18n.global.t('auto.List.769d88e4'),
  cancel: i18n.global.t('auto.List.bd9fcf46'),
  decide: i18n.global.t('auto.List.7756f1bd')
}
const VERDICT_LABELS = { approved: i18n.global.t('auto.List.7a712b9a'), rejected: i18n.global.t('auto.List.32525478') }

// COUNTS_POLL_MS 未读数轮询间隔。30s 是"值班换班时看得见新账"的量级；
// 更短只是把读接口打成压力源，而这一页的正确性不靠推送（推送归 Notifications.vue）。
const COUNTS_POLL_MS = 30000

const userStore = useUserStore()
// 后端的 operator / assignee_user_id 存的是 user id 的十进制字符串；这里统一成字符串，
// 否则 userInfo.id 是数字时 rowActions 的"是不是我认领的"会永远判否。
const me = computed(() => String(userStore.userInfo?.id ?? ''))

const { t } = useI18n()

const loading = ref(false)
const unavailable = ref(false)
const list = ref([])
const total = ref(0)
const counts = reactive({ at: '', total_open: 0, by_kind: [], unknown_kind_open: 0 })
const filter = reactive({ kind: [], status: [...OPEN_STATUSES], assignee: 'all', page: 1, page_size: 20 })

const detail = reactive({ visible: false, text: '' })
const decide = reactive({
  visible: false,
  loading: false,
  submitting: false,
  error: '',
  task: null,
  approvalId: '',
  approval: null,
  note: ''
})

const kindLabel = (k) => KIND_LABELS[k] || k || t('auto.List.485d34dc')
const kindTagType = (k) => KIND_TAG_TYPES[k] || 'info'
const statusLabel = (s) => STATUS_LABELS[s] || s
const actionLabel = (a) => ACTION_LABELS[a] || a
const verdictLabel = (v) => VERDICT_LABELS[v] || v
const subjectText = (a) => (a ? `${a.subject_type}/${a.subject_id}` : '—')

const verdicts = computed(() => verdictButtons(decide.approval))
// 没有按钮时要说得出为什么没有：只把按钮收起来，读的人看到的是"这条点不动"，
// 而他会去查的方向（权限？数据坏了？）和真原因（时限已过）往往不是同一个。
const decideHint = computed(() =>
  verdictBlockReason(decide.approval) || t('auto.List.005f2b2e')
)

function formatTime(value) {
  if (!value) return '—'
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? String(value) : d.toLocaleString()
}

// isOverdue 只判"这一行自己那一档的截止"过没过。
// 判据时刻用后端 counts.at 而不是本地 Date.now()：值班机器的钟快十分钟，
// 满屏就都是逾期，而那是个看不出来的错。拿不到 at 时退回本地钟。
function isOverdue(row) {
  const due = dueAt(row)
  if (!due) return false
  const now = counts.at ? new Date(counts.at).getTime() : Date.now()
  return due < now
}

function dueAt(row) {
  const field = slaField(row?.kind)
  const raw = field ? row[field] : ''
  if (!raw) return 0
  const t = new Date(raw).getTime()
  return Number.isNaN(t) ? 0 : t
}

function slaText(row) {
  const field = slaField(row.kind)
  if (!field) return '—'
  const due = dueAt(row)
  if (!due) return t('auto.List.a534e79f')
  if (row.status === 'done' || row.status === 'cancelled') return formatTime(row[field])
  return `${formatTime(row[field])} 前`
}

function queryParams() {
  const q = { page: filter.page, page_size: filter.page_size }
  if (filter.kind.length) q.kind = filter.kind
  if (filter.status.length) q.status = filter.status
  // 'all' 不进查询串：后端把空 assignee 读成"不过滤"，把 'all' 读成"查一个叫 all 的人"。
  if (filter.assignee === 'me') q.assignee = 'me'
  return q
}

async function fetchList() {
  loading.value = true
  try {
    const data = await humanTaskApi.list(queryParams())
    list.value = toList(data)
    total.value = Number(data?.total) || 0
    unavailable.value = false
  } catch (err) {
    // 503 与"这次查询失败"要分开：前者是配置状态，页面要长期挂着说一句话；
    // 后者刷新就好。空列表两者都不是答案，所以这里绝不清成 []假装读到了。
    unavailable.value = err?.status === 503
  } finally {
    loading.value = false
  }
}

async function fetchCounts() {
  try {
    const data = await humanTaskApi.counts()
    if (!data) return
    counts.at = data.at || ''
    counts.total_open = Number(data.total_open) || 0
    counts.by_kind = Array.isArray(data.by_kind) ? data.by_kind : []
    counts.unknown_kind_open = Number(data.unknown_kind_open) || 0
    unavailable.value = false
  } catch (err) {
    // 静默失败：上一轮读数留着。轮询每 30s 一次，弹一句"服务异常"只会变成噪音。
    if (err?.status === 503) unavailable.value = true
  }
}

function reload() {
  filter.page = 1
  fetchList()
}

function refresh() {
  fetchList()
  fetchCounts()
}

function onPageChange(page) {
  filter.page = page
  fetchList()
}

function openDetail(row) {
  detail.text = JSON.stringify(row, null, 2)
  detail.visible = true
}

async function openDecide(row) {
  decide.visible = true
  decide.task = row
  decide.approvalId = row.subject_id || ''
  decide.approval = null
  decide.note = ''
  decide.error = ''
  if (!decide.approvalId) {
    decide.error = t('auto.List.8f071b58')
    return
  }
  decide.loading = true
  try {
    decide.approval = await approvalApi.get(decide.approvalId)
  } catch (err) {
    decide.error = err?.message || t('auto.List.22c5faea')
  } finally {
    decide.loading = false
  }
}

async function submitDecide(verdict) {
  decide.submitting = true
  try {
    decide.approval = await approvalApi.decide(decide.approvalId, verdict, decide.note)
    ElMessage.success(verdict === 'approved' ? t('auto.List.58464f9c') : t('auto.List.dbf36ddc'))
    decide.visible = false
    refresh()
  } catch (err) {
    if (err?.status === 409) {
      // 两个人同时开着这条详情：后端在 409 的响应体里回的就是当前那一行，
      // 直接把它显示出来 —— 这比"刷新一下看看"准，因为它不会把别人的结论冲成新读取。
      decide.approval = err.response?.data?.data || decide.approval
      ElMessage.warning(err.message || t('auto.List.12ec8373'))
      refresh()
    }
  } finally {
    decide.submitting = false
  }
}

// 认领/释放/完成/撤销的调用点。这里不逐条写成功提示之外的分支：
// 409/403 的文案服务端已经给得很准（"这条刚被别人完成了"），
// 前端再包一句同义改写只会让两边对不上。
const ACTION_FN = {
  claim: (id) => humanTaskApi.claim(id),
  release: (id) => humanTaskApi.release(id),
  complete: (id) => humanTaskApi.complete(id)
}

async function runAction(row, action) {
  if (action === 'decide') {
    await openDecide(row)
    return
  }
  if (action === 'cancel') {
    await cancelTask(row)
    return
  }
  try {
    await ACTION_FN[action](row.id)
    ElMessage.success(`${actionLabel(action)}成功`)
    refresh()
  } catch (err) {
    if (err?.status === 409 || err?.status === 403) refresh()
  }
}

async function cancelTask(row) {
  let value
  try {
    ({ value } = await ElMessageBox.prompt(t('auto.List.df10c505'), t('auto.List.b622ae1c'), {
      inputType: 'textarea',
      inputPlaceholder: t('auto.List.c9823f5f'),
      inputValidator: (v) => (hasRejectReason(v) ? true : t('auto.List.9c365dd0')),
      confirmButtonText: t('auto.List.bd9fcf46'),
      cancelButtonText: t('auto.List.5f411223')
    }))
  } catch {
    return // 点了返回
  }
  try {
    await humanTaskApi.cancel(row.id, value)
    ElMessage.success(t('auto.List.50239f4f'))
    refresh()
  } catch (err) {
    if (err?.status === 409 || err?.status === 403) refresh()
  }
}

let pollTimer = null

onMounted(() => {
  refresh()
  pollTimer = setInterval(fetchCounts, COUNTS_POLL_MS)
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
.mb-12 {
  margin-bottom: 12px;
}
.mt-12 {
  margin-top: 12px;
}
.filter-select {
  width: 220px;
}
.counts {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  font-size: 13px;
}
.count-chip {
  white-space: nowrap;
}
.counts-at {
  color: #909399;
  font-size: 12px;
}
.overdue {
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
