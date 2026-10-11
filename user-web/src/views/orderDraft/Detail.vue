<template>
  <div class="order-draft-detail">
    <el-alert
      v-if="!me"
      class="mb-12"
      type="info"
      show-icon
      :closable="false"
      :title="$t('auto.Detail.1.d5e452')"
      :description="$t('auto.Detail.2.56d017')"
    />
    <el-alert
      v-if="loadError"
      class="mb-12"
      :type="loadError === 404 ? 'info' : 'warning'"
      show-icon
      :closable="false"
      ::title="$t('auto.Detail.3.6dd3ac')"
      ::description="$t('auto.Detail.4.a3d4e3')"
    />

    <el-card
      v-loading="loading"
      shadow="never"
    >
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button
              link
              type="primary"
              @click="goBack"
            >
              ← 返回列表
            </el-button>
            <span>{{ $t('auto.Detail.25.2ccc34') }}</span>
            <el-tag
              v-if="draft"
              size="small"
              :type="STATUS_TAGS[draft.status] || 'info'"
            >
              {{ STATUS_LABELS[draft.status] || draft.status }}
            </el-tag>
          </div>
          <div
            v-if="draft && draft.status === 'pending'"
            class="header-actions"
          >
            <el-button
              type="success"
              :disabled="!me"
              :loading="acting === 'confirm'"
              @click="doConfirm"
            >
              确认成单
            </el-button>
            <el-button
              type="danger"
              plain
              :disabled="!me"
              :loading="acting === 'cancel'"
              @click="doCancel"
            >
              取消草稿
            </el-button>
          </div>
        </div>
      </template>

      <template v-if="draft">
        <el-descriptions
          :column="2"
          border
          class="mb-12"
        >
          <el-descriptions-item :label="$t('auto.Detail.5.421db1')">
            <span class="mono">{{ draft.id }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.6.3fea7c')">
            {{ STATUS_LABELS[draft.status] || draft.status }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.7.ff0b20')">
            <span class="mono">{{ draft.customer_id }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.8.2fa2f0')">
            <span class="mono">{{ draft.owner_id || '—' }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.9.a01543')">
            {{ draft.product_name }}
            <el-tag
              v-if="draft.category"
              size="small"
              type="info"
              class="ml-6"
            >
              {{ draft.category }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.10.8f3599')">
            {{ pct(draft.confidence) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.11.0bf60b')">
            {{ draft.quantity }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.12.cd49be')">
            {{ money(draft.unit_price) }} × {{ draft.quantity }} = <b>{{ money(draft.total_amount) }}</b>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.13.26ca20')">
            {{ draft.source }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.14.ca1a5e')">
            <span class="mono">{{ draft.intent_id || '—' }}</span>
          </el-descriptions-item>
          <el-descriptions-item
            :span="2"
            :label="$t('auto.Detail.15.cbf24d')"
          >
            {{ draft.source_text || '—' }}
          </el-descriptions-item>
          <el-descriptions-item
            :span="2"
            :label="$t('auto.Detail.16.2432b5')"
          >
            {{ draft.note || '—' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.17.d9ac92')">
            {{ formatTime(draft.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Detail.18.851824')">
            <span :class="{ overdue: isExpiredSoon(draft) }">{{ formatTime(draft.expires_at) }}</span>
          </el-descriptions-item>
          <el-descriptions-item
            v-if="draft.order_id"
            :label="$t('auto.Detail.19.7b5da5')"
          >
            <span class="mono">{{ draft.order_id }}</span>
          </el-descriptions-item>
          <el-descriptions-item
            v-if="draft.confirmed_at"
            :label="$t('auto.Detail.20.df713d')"
          >
            {{ formatTime(draft.confirmed_at) }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="draft.cancel_reason"
            :span="2"
            :label="$t('auto.Detail.21.48ac07')"
          >
            {{ draft.cancel_reason }}
          </el-descriptions-item>
        </el-descriptions>

        <el-alert
          v-if="draft.status !== 'pending'"
          class="mb-12"
          type="info"
          show-icon
          :closable="false"
          :title="terminalHint(draft)"
        />

        <el-card
          v-if="draft.status === 'pending'"
          shadow="never"
        >
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Detail.26.3475fe') }}</span>
              <el-button
                type="primary"
                :disabled="!me || !dirty"
                :loading="acting === 'edit'"
                @click="doEdit"
              >
                保存修改
              </el-button>
            </div>
          </template>

          <el-form
            label-width="96px"
            @submit.prevent
          >
            <el-form-item :label="$t('auto.Detail.22.c6afa0')">
              <el-input
                v-model="form.product_name"
                maxlength="200"
                :placeholder="$t('auto.Detail.23.847e06')"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.Detail.11.0bf60b')">
              <el-input-number
                v-model="form.quantity"
                :min="1"
                :max="999999"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.Detail.24.da4abd')">
              <el-input-number
                v-model="form.unit_price"
                :min="0"
                :precision="2"
                :step="1"
              />
            </el-form-item>
            <el-form-item :label="$t('auto.Detail.16.2432b5')">
              <el-input
                v-model="form.note"
                type="textarea"
                :rows="3"
                maxlength="2000"
                show-word-limit
              />
            </el-form-item>
          </el-form>
          <p class="hint">
            只把你动过的字段发给后端（没给=不改）；数量至少 1、单价不能为负 —— 要废单请走上方「取消草稿」，不是把数量改成 0。
            保存成功后回读整条草稿，总额会在这里重算刷新。
          </p>
        </el-card>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orderDraftApi } from '@/api/orderDraft'
import { useUserStore } from '@/stores/user'

// 草稿详情 + 确认/取消/改价三个人工动作。深链 /dashboard/drafts/:id 就是
// 工作台聚合待办生成的那个 URL（sales_workbench.go），本页存在之前它是死链。
//
// 终态（已确认/已取消/已过期）是只读的：动作按钮整个消失而不是置灰 ——
// 置灰按钮会被读成「我权限不够」，而真原因是这条已经不需要人再处理了。

const STATUS_TAGS = { pending: 'warning', confirmed: 'success', cancelled: 'info', expired: 'danger' }
const STATUS_LABELS = { pending: t('auto.Detail.27.2a2772'), confirmed: t('auto.Detail.28.4113e7'), cancelled: t('auto.Detail.29.2111cc'), expired: t('auto.Detail.30.4d5ccd') }

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const me = computed(() => String(userStore.userInfo?.id ?? ''))

const { t } = useI18n()

const loading = ref(false)
const acting = ref('') // '' | 'confirm' | 'cancel' | 'edit'
const loadError = ref(0) // 0 无 / 404 不存在 / 500 读取失败
const draft = ref(null)
const form = reactive({ product_name: '', quantity: 1, unit_price: 0, note: '' })

const draftId = computed(() => String(route.params.id || ''))

// dirty 只比较动过的字段：后端对「没给」什么都不做，把原值也发过去
// 等于让「保存」按钮在无改动时也能点出一次无意义的写。
const dirty = computed(() => {
  const d = draft.value
  if (!d) return false
  return (
    form.product_name !== (d.product_name || '') ||
    form.quantity !== d.quantity ||
    form.unit_price !== d.unit_price ||
    form.note !== (d.note || '')
  )
})

function fillForm(d) {
  form.product_name = d.product_name || ''
  form.quantity = d.quantity || 1
  form.unit_price = Number(d.unit_price) || 0
  form.note = d.note || ''
}

function money(v) {
  const n = Number(v)
  return Number.isFinite(n) ? n.toFixed(2) : '—'
}

function pct(v) {
  const n = Number(v)
  return Number.isFinite(n) ? `${Math.round(n * 1000) / 10}%` : '—'
}

function formatTime(value) {
  if (!value) return '—'
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? String(value) : d.toLocaleString()
}

function isExpiredSoon(d) {
  const t = new Date(d.expires_at).getTime()
  if (Number.isNaN(t)) return false
  return t - Date.now() < 3600 * 1000
}

function terminalHint(d) {
  if (d.status === 'confirmed') return `这条已确认成单${d.order_id ? `（订单 ${d.order_id}）` : ''}，草稿侧不可再改。`
  if (d.status === 'cancelled') return `这条已取消${d.cancel_reason ? `：${d.cancel_reason}` : ''}，不可再操作。`
  if (d.status === 'expired') return t('auto.Detail.31.40d4fd')
  return ''
}

async function load() {
  if (!draftId.value) {
    loadError.value = 404
    return
  }
  loading.value = true
  loadError.value = 0
  try {
    const d = await orderDraftApi.get(draftId.value)
    draft.value = d
    fillForm(d)
  } catch (err) {
    draft.value = null
    if (err?.status === 404) {
      loadError.value = 404
    } else {
      loadError.value = 500
    }
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.push('/dashboard/drafts')
}

async function doConfirm() {
  try {
    await ElMessageBox.confirm(
      `确认把「${draft.value.product_name}」×${draft.value.quantity} 转成正式订单？确认后草稿不可再改。`,
      t('auto.Detail.32.9308e4'),
      { confirmButtonText: t('auto.Detail.32.9308e4'), cancelButtonText: t('auto.Detail.33.abe78b'), type: 'warning' }
    )
  } catch {
    return
  }
  acting.value = 'confirm'
  try {
    const result = await orderDraftApi.confirm(draftId.value)
    if (result?.order_provisional) {
      ElMessage.warning(
        `草稿已确认，但未拿到订单服务 ⇒ 订单号 ${result.order_id || '（空）'} 是本进程临时生成的，orders 表里没有这一行（需人工补建）`
      )
    } else {
      ElMessage.success(`已确认成单${result?.order_id ? `：订单 ${result.order_id}` : ''}`)
    }
    await load()
  } catch (err) {
    // 404/409：这条已被人处理或已过期 —— 回读一次让页面落到真实终态。
    if (err?.status === 404 || err?.status === 409) await load()
  } finally {
    acting.value = ''
  }
}

async function doCancel() {
  let reason
  try {
    ({ value: reason } = await ElMessageBox.prompt(
      `取消「${draft.value.product_name}」这张草稿，必须给理由`,
      t('auto.Detail.34.027c28'),
      {
        inputType: 'textarea',
        inputPlaceholder: t('auto.Detail.35.c5fe89'),
        inputValidator: (v) => (v && v.trim() ? true : t('auto.Detail.36.6210e8')),
        confirmButtonText: t('auto.Detail.34.027c28'),
        cancelButtonText: t('auto.Detail.37.5f4112')
      }
    ))
  } catch {
    return
  }
  acting.value = 'cancel'
  try {
    await orderDraftApi.cancel(draftId.value, reason.trim())
    ElMessage.success(t('auto.Detail.29.2111cc'))
    await load()
  } catch (err) {
    if (err?.status === 404 || err?.status === 409) await load()
  } finally {
    acting.value = ''
  }
}

async function doEdit() {
  const d = draft.value
  if (!d) return
  const updates = {}
  if (form.product_name !== (d.product_name || '')) updates.product_name = form.product_name.trim()
  if (form.quantity !== d.quantity) updates.quantity = form.quantity
  if (form.unit_price !== d.unit_price) updates.unit_price = form.unit_price
  if (form.note !== (d.note || '')) updates.note = form.note
  if (!Object.keys(updates).length) return

  acting.value = 'edit'
  try {
    const fresh = await orderDraftApi.edit(draftId.value, updates)
    draft.value = fresh
    fillForm(fresh)
    ElMessage.success(t('auto.Detail.38.20feaf'))
  } catch (err) {
    // 400（校验）由拦截器原样弹后端文案；409（同客户同产品 pending 撞车 / 已过态）回读落态。
    if (err?.status === 404 || err?.status === 409) await load()
  } finally {
    acting.value = ''
  }
}

onMounted(load)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}
.header-actions {
  display: flex;
  gap: 8px;
}
.mb-12 {
  margin-bottom: 12px;
}
.ml-6 {
  margin-left: 6px;
}
.overdue {
  color: #f56c6c;
}
.hint {
  color: #909399;
  font-size: 13px;
}
.mono {
  font-family: monospace;
  font-size: 12px;
}
</style>
