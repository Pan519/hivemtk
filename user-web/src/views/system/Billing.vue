<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Billing.da6b9b7b') }}</span>
        </div>
      </template>

      <el-tabs v-model="tab">
        <!-- 计费方案 CRUD -->
        <el-tab-pane :label="$t('auto.Billing.a56bd8a7')" name="plans">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openPlanDialog()">{{ $t('auto.Billing.90e34b85') }}</el-button>
          </div>
          <el-table v-loading="plansLoading" :data="plans" border size="small">
            <el-table-column prop="code" :label="$t('auto.Billing.cc6c35a3')" width="110" />
            <el-table-column prop="name" :label="$t('auto.Billing.d7ec2d3f')" width="120" />
            <el-table-column :label="$t('auto.Billing.a1af448f')" width="90">
              <template #default="{ row }">{{ row.price_mode === 'per_seat' ? '按人头' : '固定' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.e45ac657')" width="100">
              <template #default="{ row }">¥{{ (row.price_cents / 100).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.30806dcb')" width="110">
              <template #default="{ row }">¥{{ (row.yearly_price_cents / 100).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.937bdd75')" width="90">
              <template #default="{ row }">{{ row.max_users || '不限' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.5d4a37e5')" width="100">
              <template #default="{ row }">{{ row.max_messages_monthly || '不限' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.3fea7ca7')" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                  {{ row.status === 'active' ? '在售' : '下架' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="remark" :label="$t('auto.Billing.2432b575')" min-width="160" show-overflow-tooltip />
            <el-table-column :label="$t('auto.Billing.2b6bc0f2')" width="150" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="openPlanDialog(row)">{{ $t('auto.Billing.95b351c8') }}</el-button>
                <el-button size="small" type="danger" @click="removePlan(row)">{{ $t('auto.Billing.2f4aaddd') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 团队订阅 -->
        <el-tab-pane :label="$t('auto.Billing.5c2f1147')" name="subs">
          <el-table v-loading="subsLoading" :data="subs" border size="small">
            <el-table-column :label="$t('auto.Billing.4675fbaf')" width="140">
              <template #default="{ row }">{{ teamName(row.team_id) }}</template>
            </el-table-column>
            <el-table-column prop="plan_code" :label="$t('auto.Billing.729de1d0')" width="110" />
            <el-table-column :label="$t('auto.Billing.3fea7ca7')" width="90">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                  {{ row.status === 'active' ? '生效中' : '已停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.308fdfa0')" width="160">
              <template #default="{ row }">{{ fmtTime(row.starts_at) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.a2bf2b41')" width="160">
              <template #default="{ row }">{{ fmtTime(row.ends_at) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.d7550fbb')" width="90">
              <template #default="{ row }">{{ row.auto_renew ? '是' : '否' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.2b6bc0f2')" width="220" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="openAssign(row)">{{ $t('auto.Billing.3bf8b198') }}</el-button>
                <el-button v-if="row.status === 'active'" size="small" type="warning" plain @click="cancelSub(row)">{{ $t('auto.Billing.5c56a889') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div style="margin-top: 12px">
            <el-button type="primary" size="small" @click="openAssign(null)">{{ $t('auto.Billing.b7da3b7c') }}</el-button>
          </div>
        </el-tab-pane>

        <!-- 账单记录 -->
        <el-tab-pane :label="$t('auto.Billing.931a6ffd')" name="records">
          <div style="margin-bottom: 12px; display: flex; gap: 8px; align-items: center">
            <el-select v-model="recordTeamId" :placeholder="$t('auto.Billing.c7f83c74')" clearable size="small" style="width: 180px" @change="loadRecords">
              <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
            </el-select>
            <el-select v-model="recordStatus" :placeholder="$t('auto.Billing.443483c9')" clearable size="small" style="width: 140px" @change="loadRecords">
              <el-option :label="$t('auto.Billing.fccab93a')" value="pending" />
              <el-option :label="$t('auto.Billing.b511e9bd')" value="paid" />
              <el-option :label="$t('auto.Billing.c4a4d9f7')" value="void" />
            </el-select>
            <el-button type="primary" size="small" @click="openRecordDialog">{{ $t('auto.Billing.2a7bf9da') }}</el-button>
          </div>
          <el-table v-loading="recordsLoading" :data="records" border size="small">
            <el-table-column :label="$t('auto.Billing.4675fbaf')" width="140">
              <template #default="{ row }">{{ teamName(row.team_id) }}</template>
            </el-table-column>
            <el-table-column prop="plan_code" :label="$t('auto.Billing.729de1d0')" width="110" />
            <el-table-column :label="$t('auto.Billing.4cf24aba')" width="110">
              <template #default="{ row }">¥{{ (row.amount_cents / 100).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="months" :label="$t('auto.Billing.479196cf')" width="70" />
            <el-table-column :label="$t('auto.Billing.3fea7ca7')" width="90">
              <template #default="{ row }">
                <el-tag :type="row.status === 'paid' ? 'success' : row.status === 'pending' ? 'warning' : 'info'" size="small">
                  {{ { pending: '应收', paid: '已收', void: '作废' }[row.status] }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Billing.05c5dc04')" width="160">
              <template #default="{ row }">{{ fmtTime(row.paid_at) }}</template>
            </el-table-column>
            <el-table-column prop="remark" :label="$t('auto.Billing.2432b575')" min-width="140" show-overflow-tooltip />
            <el-table-column :label="$t('auto.Billing.2b6bc0f2')" width="170" fixed="right">
              <template #default="{ row }">
                <el-button v-if="row.status === 'pending'" size="small" type="success" plain @click="markPaid(row)">{{ $t('auto.Billing.37943999') }}</el-button>
                <el-button v-if="row.status === 'pending'" size="small" type="warning" plain @click="voidRecord(row)">{{ $t('auto.Billing.c4a4d9f7') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-model:current-page="recordPage"
            :page-size="recordSize"
            :total="recordTotal"
            layout="total, prev, pager, next"
            small
            style="margin-top: 12px; justify-content: flex-end"
            @current-change="loadRecords"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 方案弹窗 -->
    <el-dialog v-model="planDialogVisible" ::title="$t('auto.Billing.2a3697fa')" width="480px">
      <el-form label-width="110px">
        <el-form-item :label="$t('auto.Billing.cc6c35a3')">
          <el-input v-model="planForm.code" :disabled="planEditing" placeholder="free/pro/..." maxlength="32" />
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.d7ec2d3f')">
          <el-input v-model="planForm.name" maxlength="60" />
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.a1af448f')">
          <el-radio-group v-model="planForm.price_mode">
            <el-radio value="flat">{{ $t('auto.Billing.468c5b11') }}</el-radio>
            <el-radio value="per_seat">{{ $t('auto.Billing.2190545d') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="planForm.price_mode === 'flat'" :label="$t('auto.Billing.03253c56')">
          <el-input-number v-model="planForm.priceYuan" :min="0" :precision="2" style="width: 160px" />
        </el-form-item>
        <el-form-item v-if="planForm.price_mode === 'per_seat'" :label="$t('auto.Billing.eb9c24d2')">
          <el-input-number v-model="planForm.seatYuan" :min="0" :precision="2" style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.7c3c701d')">
          <el-input-number v-model="planForm.yearlyYuan" :min="0" :precision="2" style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.937bdd75')">
          <el-input-number v-model="planForm.max_users" :min="0" style="width: 160px" />
          <span class="muted" style="margin-left: 8px">{{ $t('auto.Billing.768abb4d') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.5d4a37e5')">
          <el-input-number v-model="planForm.max_messages_monthly" :min="0" style="width: 160px" />
          <span class="muted" style="margin-left: 8px">{{ $t('auto.Billing.768abb4d') }}</span>
        </el-form-item>
        <el-form-item v-if="planEditing" :label="$t('auto.Billing.3fea7ca7')">
          <el-radio-group v-model="planForm.status">
            <el-radio value="active">{{ $t('auto.Billing.d31f8f1e') }}</el-radio>
            <el-radio value="retired">{{ $t('auto.Billing.d2379a56') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.2432b575')">
          <el-input v-model="planForm.remark" maxlength="200" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="planDialogVisible = false">{{ $t('auto.Billing.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="savePlan">{{ $t('auto.Billing.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    <!-- 指派/续期弹窗 -->
    <el-dialog v-model="assignVisible" :title="$t('auto.Billing.f896ef6b')" width="440px" append-to-body>
      <el-form label-width="110px">
        <el-form-item :label="$t('auto.Billing.4675fbaf')">
          <el-select v-model="assignForm.team_id" :placeholder="$t('auto.Billing.1e2e42c5')" style="width: 220px" :disabled="!!assignEditingTeam">
            <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.729de1d0')">
          <el-select v-model="assignForm.plan_id" style="width: 220px">
            <el-option v-for="p in activePlans" :key="p.id" :label="`${p.name}（¥${(p.price_cents / 100).toFixed(2)}/月）`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.479196cf')">
          <el-input-number v-model="assignForm.months" :min="1" :max="36" style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.d7550fbb')">
          <el-switch v-model="assignForm.auto_renew" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignVisible = false">{{ $t('auto.Billing.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="saveAssign">{{ $t('auto.Billing.e83a256e') }}</el-button>
      </template>
    </el-dialog>

    <!-- 应收账单弹窗 -->
    <el-dialog v-model="recordDialogVisible" :title="$t('auto.Billing.2a7bf9da')" width="440px" append-to-body>
      <el-form label-width="110px">
        <el-form-item :label="$t('auto.Billing.4675fbaf')">
          <el-select v-model="recordForm.team_id" style="width: 220px">
            <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.729de1d0')">
          <el-select v-model="recordForm.plan_id" style="width: 220px">
            <el-option v-for="p in plans" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.479196cf')">
          <el-input-number v-model="recordForm.months" :min="1" :max="36" style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.92c95a39')">
          <el-input-number v-model="recordForm.amountYuan" :min="0" :precision="2" style="width: 160px" />
          <span class="muted" style="margin-left: 8px">{{ $t('auto.Billing.edb1967f') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.Billing.2432b575')">
          <el-input v-model="recordForm.remark" maxlength="200" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="recordDialogVisible = false">{{ $t('auto.Billing.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="saveRecord">{{ $t('auto.Billing.d9ac9228') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { http } from '@/utils/request'
import { teamAdminApi } from '@/api/teamAdmin'

// 计费 API（独立成段：与团队管理面同源但路径不同）
const billingApi = {
  listPlans: () => http.get('/api/manage/billing/plans'),
  createPlan: (d) => http.post('/api/manage/billing/plans', d),
  updatePlan: (id, d) => http.put(`/api/manage/billing/plans/${id}`, d),
  removePlan: (id) => http.delete(`/api/manage/billing/plans/${id}`),
  listSubs: () => http.get('/api/manage/billing/subscriptions'),
  assign: (teamId, d) => http.put(`/api/manage/billing/subscriptions/${teamId}`, d),
  cancel: (teamId) => http.post(`/api/manage/billing/subscriptions/${teamId}/cancel`),
  listRecords: (params) => http.get('/api/manage/billing/records', { params }),
  createRecord: (d) => http.post('/api/manage/billing/records', d),
  payRecord: (id) => http.put(`/api/manage/billing/records/${id}/pay`),
  voidRecord: (id) => http.put(`/api/manage/billing/records/${id}/void`)
}

const { t } = useI18n()

const tab = ref('plans')
const teams = ref([])

// 方案
const plans = ref([])
const plansLoading = ref(false)
const planDialogVisible = ref(false)
const planEditing = ref(false)
const planEditingId = ref(0)
const saving = ref(false)
const planForm = ref({})

// 订阅
const subs = ref([])
const subsLoading = ref(false)
const assignVisible = ref(false)
const assignEditingTeam = ref(null)
const assignForm = ref({ team_id: null, plan_id: null, months: 1, auto_renew: false })

// 账单
const records = ref([])
const recordsLoading = ref(false)
const recordTotal = ref(0)
const recordPage = ref(1)
const recordSize = 20
const recordTeamId = ref(null)
const recordStatus = ref('')
const recordDialogVisible = ref(false)
const recordForm = ref({})

const activePlans = computed(() => plans.value.filter((p) => p.status === 'active'))
const teamName = (id) => teams.value.find((t) => t.id === id)?.name || `团队 #${id}`
const fmtTime = (v) => (v ? String(v).replace('T', ' ').slice(0, 19) : '—')

const loadPlans = async () => {
  plansLoading.value = true
  try {
    const res = await billingApi.listPlans()
    plans.value = (res?.data ?? res ?? {}).list ?? []
  } finally {
    plansLoading.value = false
  }
}

const loadSubs = async () => {
  subsLoading.value = true
  try {
    const res = await billingApi.listSubs()
    subs.value = (res?.data ?? res ?? {}).list ?? []
  } finally {
    subsLoading.value = false
  }
}

const loadRecords = async () => {
  recordsLoading.value = true
  try {
    const params = { page: recordPage.value, size: recordSize }
    if (recordTeamId.value) params.team_id = recordTeamId.value
    if (recordStatus.value) params.status = recordStatus.value
    const res = await billingApi.listRecords(params)
    const data = res?.data ?? res ?? {}
    records.value = data.list ?? []
    recordTotal.value = data.total ?? 0
  } finally {
    recordsLoading.value = false
  }
}

const loadTeams = async () => {
  const res = await teamAdminApi.list({ page: 1, size: 200 })
  teams.value = (res?.data ?? res ?? {}).list ?? []
}

const openPlanDialog = (row) => {
  planEditing.value = !!row
  planEditingId.value = row?.id || 0
  planForm.value = row
    ? { code: row.code, name: row.name, price_mode: row.price_mode || 'flat', priceYuan: row.price_cents / 100, seatYuan: (row.per_seat_cents || 0) / 100, yearlyYuan: row.yearly_price_cents / 100, max_users: row.max_users, max_messages_monthly: row.max_messages_monthly, status: row.status, remark: row.remark || '' }
    : { code: '', name: '', price_mode: 'flat', priceYuan: 0, seatYuan: 0, yearlyYuan: 0, max_users: 0, max_messages_monthly: 0, status: 'active', remark: '' }
  planDialogVisible.value = true
}

const savePlan = async () => {
  saving.value = true
  try {
    const payload = {
      code: planForm.value.code,
      name: planForm.value.name,
      price_mode: planForm.value.price_mode || 'flat',
      per_seat_cents: Math.round((planForm.value.seatYuan || 0) * 100),
      price_cents: Math.round((planForm.value.priceYuan || 0) * 100),
      yearly_price_cents: Math.round((planForm.value.yearlyYuan || 0) * 100),
      max_users: planForm.value.max_users || 0,
      max_messages_monthly: planForm.value.max_messages_monthly || 0,
      remark: planForm.value.remark
    }
    if (planEditing.value) {
      payload.status = planForm.value.status
      await billingApi.updatePlan(planEditingId.value, payload)
    } else {
      await billingApi.createPlan(payload)
    }
    ElMessage.success(t('auto.Billing.f8dfedcd'))
    planDialogVisible.value = false
    await loadPlans()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Billing.6de920b4'))
  } finally {
    saving.value = false
  }
}

const removePlan = async (row) => {
  await ElMessageBox.confirm(`确认删除方案「${row.name}」？被订阅引用的方案需先下架。`, t('auto.Billing.50eaf94d'), { type: 'warning' })
  try {
    await billingApi.removePlan(row.id)
    ElMessage.success(t('auto.Billing.5cc23262'))
    await loadPlans()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Billing.acf0664a'))
  }
}

const openAssign = (row) => {
  assignEditingTeam.value = row?.team_id || null
  assignForm.value = {
    team_id: row?.team_id || null,
    plan_id: row?.plan_id || null,
    months: 1,
    auto_renew: row?.auto_renew || false
  }
  assignVisible.value = true
}

const saveAssign = async () => {
  if (!assignForm.value.team_id || !assignForm.value.plan_id) {
    ElMessage.warning(t('auto.Billing.3c691082'))
    return
  }
  saving.value = true
  try {
    await billingApi.assign(assignForm.value.team_id, {
      plan_id: assignForm.value.plan_id,
      months: assignForm.value.months,
      auto_renew: assignForm.value.auto_renew
    })
    ElMessage.success(t('auto.Billing.6c8a95ba'))
    assignVisible.value = false
    await loadSubs()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Billing.67c95bbe'))
  } finally {
    saving.value = false
  }
}

const cancelSub = async (row) => {
  await ElMessageBox.confirm(`确认停用 ${teamName(row.team_id)} 的订阅？`, t('auto.Billing.e83a256e'), { type: 'warning' })
  await billingApi.cancel(row.team_id)
  ElMessage.success(t('auto.Billing.69b0f684'))
  await loadSubs()
}

const openRecordDialog = () => {
  recordForm.value = { team_id: recordTeamId.value || null, plan_id: null, months: 1, amountYuan: 0, remark: '' }
  recordDialogVisible.value = true
}

const saveRecord = async () => {
  if (!recordForm.value.team_id || !recordForm.value.plan_id) {
    ElMessage.warning(t('auto.Billing.3c691082'))
    return
  }
  saving.value = true
  try {
    await billingApi.createRecord({
      team_id: recordForm.value.team_id,
      plan_id: recordForm.value.plan_id,
      months: recordForm.value.months,
      amount_cents: Math.round((recordForm.value.amountYuan || 0) * 100),
      remark: recordForm.value.remark
    })
    ElMessage.success(t('auto.Billing.e1030f60'))
    recordDialogVisible.value = false
    await loadRecords()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Billing.a889286a'))
  } finally {
    saving.value = false
  }
}

const markPaid = async (row) => {
  await billingApi.payRecord(row.id)
  ElMessage.success(t('auto.Billing.a36166c2'))
  await loadRecords()
}

const voidRecord = async (row) => {
  await ElMessageBox.confirm(t('auto.Billing.6d74d1b7'), t('auto.Billing.e83a256e'), { type: 'warning' })
  await billingApi.voidRecord(row.id)
  ElMessage.success(t('auto.Billing.dcc961ae'))
  await loadRecords()
}

onMounted(() => {
  loadPlans()
  loadSubs()
  loadRecords()
  loadTeams()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.muted {
  color: #999;
  font-size: 12px;
}
</style>
