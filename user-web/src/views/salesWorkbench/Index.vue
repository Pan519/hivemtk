<template>
  <div class="sales-workbench">
    <div class="wb-header">
      <div class="wb-title">
        <h3>{{ $t('auto.Index.05f65226') }}</h3>
        <span class="wb-sub">
          {{ overview.name || '—' }}{{ overview.team ? ' · ' + overview.team : '' }}
          <el-tag
            v-if="overview.my_rank > 0"
            type="success"
            size="small"
            class="wb-rank"
          >
            本月第 {{ overview.my_rank }} 名
          </el-tag>
        </span>
      </div>
      <div class="wb-actions">
        <el-select
          v-model="days"
          style="width: 120px"
          @change="loadTeam"
        >
          <el-option
            v-for="d in dayOptions"
            :key="d"
            :value="d"
            :label="`近 ${d} 天`"
          />
        </el-select>
        <el-button
          :icon="Refresh"
          :loading="loading"
          @click="loadAll"
        >
          刷新
        </el-button>
      </div>
    </div>

    <el-alert
      v-if="error"
      :title="error"
      type="warning"
      show-icon
      :closable="false"
      class="wb-alert"
    />

    <!-- 今日 / 本月 / 关键指标 -->
    <el-row
      :gutter="12"
      class="wb-row"
    >
      <el-col
        :xs="24"
        :md="8"
      >
        <el-card
          shadow="never"
          class="wb-card"
        >
          <template #header>
            今日
          </template>
          <div class="wb-stats">
            <div class="wb-stat">
              <span class="num">{{ n(overview.today?.new_orders) }}</span><span class="lab">{{ $t('auto.Index.ea105393') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ money(overview.today?.new_revenue) }}</span><span class="lab">{{ $t('auto.Index.a3c21582') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.today?.follow_ups) }}</span><span class="lab">{{ $t('auto.Index.994dcf28') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.today?.conversions) }}</span><span class="lab">{{ $t('auto.Index.613b48bd') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.today?.conversion_rate) }}</span><span class="lab">{{ $t('auto.Index.71505437') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.today?.ai_deals) }}</span><span class="lab">{{ $t('auto.Index.124bec51') }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col
        :xs="24"
        :md="8"
      >
        <el-card
          shadow="never"
          class="wb-card"
        >
          <template #header>
            本月
          </template>
          <div class="wb-stats">
            <div class="wb-stat">
              <span class="num">{{ n(overview.month?.total_orders) }}</span><span class="lab">{{ $t('auto.Index.4c117f20') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ money(overview.month?.total_revenue) }}</span><span class="lab">{{ $t('auto.Index.a883e8d5') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.month?.follow_ups) }}</span><span class="lab">{{ $t('auto.Index.994dcf28') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.month?.conversions) }}</span><span class="lab">{{ $t('auto.Index.613b48bd') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.month?.conversion_rate) }}</span><span class="lab">{{ $t('auto.Index.71505437') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.month?.new_customers) }}</span><span class="lab">{{ $t('auto.Index.8759f4b4') }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col
        :xs="24"
        :md="8"
      >
        <el-card
          shadow="never"
          class="wb-card"
        >
          <template #header>
            关键指标（本月口径）
          </template>
          <div class="wb-stats">
            <div class="wb-stat">
              <span class="num">{{ money(overview.metrics?.avg_deal_amount) }}</span><span class="lab">{{ $t('auto.Index.53c7e7b4') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.metrics?.renewal_rate) }}</span><span class="lab">{{ $t('auto.Index.b1a24e45') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.metrics?.repurchase_rate) }}</span><span class="lab">{{ $t('auto.Index.106d4954') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.metrics?.churn_rate) }}</span><span class="lab">{{ $t('auto.Index.f620ee32') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.metrics?.active_customers) }}</span><span class="lab">{{ $t('auto.Index.023ff0d3') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.metrics?.dormant_customers) }}</span><span class="lab">{{ $t('auto.Index.5ad07edf') }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快链（A10：GetQuickActions 读口的消费区，URL 已全为真实落点） -->
    <el-card
      v-if="quickActions.length"
      shadow="never"
      class="wb-card wb-block"
    >
      <template #header>
        <span>{{ $t('auto.Index.9fdf0e64') }}</span>
        <span class="wb-hint">{{ $t('auto.Index.661cacac') }}</span>
      </template>
      <div class="wb-quick">
        <el-button
          v-for="a in quickActions"
          :key="a.id"
          class="wb-quick-btn"
          @click="goQuick(a)"
        >
          {{ a.title }}
        </el-button>
      </div>
    </el-card>

    <!-- 我的待办 -->
    <el-card
      shadow="never"
      class="wb-card wb-block"
    >
      <template #header>
        <span>{{ $t('auto.Index.536d2e8e') }}</span>
        <span class="wb-hint">共 {{ overview.todos?.length || 0 }} 条，点击可跳转对应落点</span>
      </template>
      <el-table
        v-if="overview.todos?.length"
        :data="overview.todos"
        size="small"
      >
        <el-table-column
          :label="$t('auto.Index.ee8ecb9e')"
          width="80"
        >
          <template #default="{ row }">
            <el-tag
              :type="priorityType(row.priority)"
              size="small"
            >
              {{ priorityLabel(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.Index.226b0912')"
          width="90"
        >
          <template #default="{ row }">
            <el-tag
              size="small"
              effect="plain"
            >
              {{ typeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="title"
          :label="$t('auto.Index.1bff2322')"
          min-width="160"
          show-overflow-tooltip
        />
        <el-table-column
          prop="description"
          :label="$t('auto.Index.f411d0f1')"
          min-width="200"
          show-overflow-tooltip
        />
        <el-table-column
          :label="$t('auto.Index.cd8fbf6f')"
          width="140"
        >
          <template #default="{ row }">
            {{ fmtTime(row.due_at) }}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('auto.Index.2b6bc0f2')"
          width="90"
        >
          <template #default="{ row }">
            <el-button
              v-if="row.url && row.url.startsWith('/')"
              link
              type="primary"
              size="small"
              @click="goTodo(row)"
            >
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty
        v-else
        :description="$t('auto.Index.44406eaf')"
        :image-size="60"
      />
    </el-card>

    <!-- 团队仪表盘（/team-dashboard 读口） -->
    <el-row
      :gutter="12"
      class="wb-row"
    >
      <el-col
        :xs="24"
        :md="14"
      >
        <el-card
          shadow="never"
          class="wb-card"
        >
          <template #header>
            <span>团队榜（近 {{ days }} 天）</span>
            <span class="wb-hint">共 {{ team.total_customers ?? 0 }} 客户进入旅程</span>
          </template>
          <el-table
            v-if="team.top_sales?.length"
            :data="team.top_sales"
            size="small"
          >
            <el-table-column
              prop="rank"
              label="#"
              width="50"
            />
            <el-table-column
              prop="name"
              :label="$t('auto.Index.7f0dc5e3')"
              min-width="90"
            />
            <el-table-column
              prop="team"
              :label="$t('auto.Index.4675fbaf')"
              width="90"
            />
            <el-table-column
              prop="total_orders"
              :label="$t('auto.Index.4c117f20')"
              width="70"
            />
            <el-table-column
              :label="$t('auto.Index.a883e8d5')"
              width="120"
            >
              <template #default="{ row }">
                {{ money(row.total_revenue) }}
              </template>
            </el-table-column>
            <el-table-column
              prop="total_follow_ups"
              :label="$t('auto.Index.994dcf28')"
              width="70"
            />
            <el-table-column
              :label="$t('auto.Index.71505437')"
              width="90"
            >
              <template #default="{ row }">
                {{ rate(row.conversion_rate) }}
              </template>
            </el-table-column>
            <el-table-column
              prop="active_days"
              :label="$t('auto.Index.9dc1f6f4')"
              width="80"
            />
          </el-table>
          <el-empty
            v-else
            :description="$t('auto.Index.3b80b3fb')"
            :image-size="60"
          />
        </el-card>
      </el-col>
      <el-col
        :xs="24"
        :md="10"
      >
        <el-card
          shadow="never"
          class="wb-card"
        >
          <template #header>
            转化漏斗（草稿/旅程内存态）
          </template>
          <template v-if="team.funnel?.stages?.length">
            <div class="wb-funnel">
              <div
                v-for="st in team.funnel.stages"
                :key="st.stage"
                class="wb-funnel-row"
              >
                <span class="st-lab">{{ st.label }}</span>
                <el-progress
                  :percentage="stagePct(st.stage_rate)"
                  :stroke-width="14"
                  class="st-bar"
                />
                <span class="st-meta">{{ st.customers }} 人 · 步转化 {{ rate(st.step_rate) }}</span>
              </div>
            </div>
            <div class="wb-funnel-foot">
              总进入 {{ team.funnel.total_entered }} · 成交 {{ team.funnel.total_won }} ·
              端到端 {{ rate(team.funnel.end_to_end_rate) }} · 平均停留 {{ (team.funnel.avg_dwell_days || 0).toFixed(1) }} 天
            </div>
          </template>
          <!-- journey 未装配时后端回 funnel=null（少一块 ≠ 500），按"未装配"措辞，不当成没数据 -->
          <el-empty
            v-else
            :description="$t('auto.Index.e3f17730')"
            :image-size="60"
          />
        </el-card>
      </el-col>
    </el-row>

    <!-- 销冠画像（/champion 独立读口）+ AI 产能 -->
    <el-row
      :gutter="12"
      class="wb-row"
    >
      <el-col
        :xs="24"
        :md="12"
      >
        <el-card
          shadow="never"
          class="wb-card"
        >
          <template #header>
            <span>销冠画像（近 {{ days }} 天）</span>
            <span class="wb-hint">均转化 {{ rate(champion.avg_conversion_rate) }} · 均单 {{ money(champion.avg_deal_amount) }}</span>
          </template>
          <div
            v-if="champion.common_tags?.length"
            class="wb-tags"
          >
            <el-tag
              v-for="t in champion.common_tags"
              :key="t"
              size="small"
            >
              {{ t }}
            </el-tag>
          </div>
          <div
            v-if="champion.insights?.length"
            class="wb-list"
          >
            <div
              v-for="(ins, i) in champion.insights"
              :key="i"
              class="wb-list-item"
            >
              · {{ ins }}
            </div>
          </div>
          <div
            v-if="champion.recommended_sops?.length"
            class="wb-list"
          >
            <div class="wb-list-head">
              推荐 SOP
            </div>
            <div
              v-for="s in champion.recommended_sops"
              :key="s"
              class="wb-list-item"
            >
              · {{ s }}
            </div>
          </div>
          <el-empty
            v-if="!champion.common_tags?.length && !champion.insights?.length"
            :description="$t('auto.Index.dfcbe644')"
            :image-size="60"
          />
        </el-card>
      </el-col>
      <el-col
        :xs="24"
        :md="12"
      >
        <el-card
          shadow="never"
          class="wb-card"
        >
          <template #header>
            <span>{{ $t('auto.Index.4d22acba') }}</span>
            <span class="wb-hint">{{ $t('auto.Index.ab516c10') }}</span>
          </template>
          <div class="wb-stats">
            <div class="wb-stat">
              <span class="num">{{ n(overview.ai_product?.total_ai_deals) }}</span><span class="lab">{{ $t('auto.Index.124bec51') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.ai_product?.ai_replied) }}</span><span class="lab">{{ $t('auto.Index.53ac193e') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.ai_product?.reply_rate) }}</span><span class="lab">{{ $t('auto.Index.0d078219') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.ai_product?.transfer_rate) }}</span><span class="lab">{{ $t('auto.Index.381706fc') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.ai_product?.ai_conversion_rate) }}</span><span class="lab">{{ $t('auto.Index.3da64002') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ rate(overview.ai_product?.productivity_gain) }}</span><span class="lab">{{ $t('auto.Index.b9fb7351') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ n(overview.ai_product?.total_cost_tokens) }}</span><span class="lab">{{ $t('auto.Index.c75aff11') }}</span>
            </div>
            <div class="wb-stat">
              <span class="num">{{ money(overview.ai_product?.solo_deal_amount) }}</span><span class="lab">{{ $t('auto.Index.3e968096') }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { salesWorkbenchApi } from '@/api/salesWorkbench'
import { useUserStore } from '@/stores/user'

// 销售工作台首页：消费 /api/sales-workbench 三读口（overview / team-dashboard /
// champion）。三条此前在 user-web 零引用（I6/I7 交付了后端半），本页是它们的
// 第一个前端消费方 —— 与 /sales-cockpit（AI 运维聚合）刻意分页，两套域不混。
//
// 错误口径对齐后端判据：503 = 服务未装配（app.InitSalesWorkbenchRuntime 没跑，
// 不是没数据）；400 = sales_id/days 参数错；漏斗 null = journey 未装配，单独措辞，
// 不把它说成"暂无数据"。rate 后端已 *100（0–100 口径），前端不再乘。
const router = useRouter()
const userStore = useUserStore()

const { t } = useI18n()

const loading = ref(false)
const error = ref('')
const days = ref(30)
const dayOptions = [7, 30, 90]

const overview = ref({})
const team = ref({})
const champion = ref({})
const quickActions = ref([])

const me = () => String(userStore.userInfo?.id ?? '')

const n = (v) => Number(v ?? 0).toLocaleString('zh-CN')
const money = (v) =>
  '¥' + Number(v ?? 0).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
const rate = (v) => Number(v ?? 0).toFixed(1) + '%'
const stagePct = (v) => Math.min(100, Math.max(0, Number(v ?? 0)))
const fmtTime = (t) => (t ? String(t).replace('T', ' ').slice(0, 16) : '—')

const priorityLabel = (p) => ({ 5: t('auto.Index.cb0ed5ef'), 4: t('auto.Index.4296d7d2'), 3: t('auto.Index.19ac677a') }[p] || t('auto.Index.35242cff'))
const priorityType = (p) => (p >= 5 ? 'danger' : p === 4 ? 'warning' : 'info')
const typeLabel = (t) => ({ draft: t('auto.Index.83861e99'), followup: t('auto.Index.994dcf28') }[t] || t || '—')

function errMsg(e, what) {
  const s = e?.status
  if (s === 503) return `${what}失败：销售工作台服务未装配（装配缺失 ≠ 没数据，请联系运维）`
  if (s === 400) return `${what}失败：参数不合法（sales_id/days 由后端校验）`
  return `${what}失败：${e?.message || e}`
}

async function loadOverview() {
  const id = me()
  if (!id) {
    error.value = t('auto.Index.689e238c')
    return
  }
  overview.value = await salesWorkbenchApi.overview(id)
}

async function loadTeam() {
  const [dash, champ] = await Promise.all([
    salesWorkbenchApi.teamDashboard(days.value),
    salesWorkbenchApi.champion(days.value)
  ])
  team.value = dash || {}
  champion.value = champ || {}
}

// 快链（A10）：GetQuickActions 此前零 HTTP 暴露，本页快链区是它的第一个
// 消费方。失败不单独吞：与三读口同 Promise.all，503 同措辞。
async function loadQuickActions() {
  const res = await salesWorkbenchApi.quickActions()
  quickActions.value = res?.list || []
}

async function loadAll() {
  loading.value = true
  error.value = ''
  try {
    await Promise.all([loadOverview(), loadTeam(), loadQuickActions()])
  } catch (e) {
    error.value = errMsg(e, t('auto.Index.1ccf798e'))
  } finally {
    loading.value = false
  }
}

// 统一的"先解析再跳"：落点在路由表里才 push，否则给提示不进 NotFound。
// 待办与快链共用 —— 快链 URL 已改指真实落点（A10），保留探活是防未来
// URL 再漂移时静默 404。
function goResolved(url, kind) {
  if (!url?.startsWith('/')) return
  const r = router.resolve(url)
  if (r.matched.length) {
    router.push(url)
  } else {
    ElMessage.warning(`该${kind}的落点页尚未交付`)
  }
}

function goTodo(row) {
  goResolved(row.url, t('auto.Index.b60ec8b0'))
}

function goQuick(action) {
  goResolved(action.url, t('auto.Index.9b4a031b'))
}

onMounted(loadAll)
</script>

<style scoped>
.sales-workbench {
  padding: 16px;
}
.wb-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.wb-title h3 {
  margin: 0;
}
.wb-sub {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.wb-actions {
  display: flex;
  gap: 8px;
}
.wb-alert {
  margin-bottom: 12px;
}
.wb-row {
  margin-bottom: 12px;
}
.wb-card {
  margin-bottom: 0;
}
.wb-block {
  margin-bottom: 12px;
}
.wb-hint {
  float: right;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
.wb-quick {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.wb-quick-btn {
  margin-left: 0;
}
.wb-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
.wb-stat {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.wb-stat .num {
  font-size: 16px;
  font-weight: 600;
}
.wb-stat .lab {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
.wb-funnel-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.wb-funnel-row .st-lab {
  width: 72px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.wb-funnel-row .st-bar {
  flex: 1;
}
.wb-funnel-row .st-meta {
  width: 128px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  text-align: right;
}
.wb-funnel-foot {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}
.wb-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 8px;
}
.wb-list {
  margin-bottom: 8px;
}
.wb-list-head {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}
.wb-list-item {
  font-size: 13px;
  line-height: 1.7;
}
</style>
