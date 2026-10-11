<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ isPlatform ? '团队管理（租户）' : '我的团队' }}</span>
          <el-button v-if="isPlatform" type="primary" size="small" @click="openCreate">{{ $t('auto.Teams.8e10eedb') }}</el-button>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="isPlatform
          ? '超级管理员视角：全量管理所有团队。成员调拨后该用户旧令牌立即作废，需重新登录。'
          : '团队管理员视角：仅本团队（服务端强制）。可管理成员、生成邀请码、查看订阅与账单。'"
        style="margin-bottom: 12px"
      />

      <!-- 团队管理员：本团队信息 + 订阅账单卡 -->
      <template v-if="!isPlatform">
        <el-descriptions v-if="myTeam" :column="3" border size="small" style="margin-bottom: 12px">
          <el-descriptions-item :label="$t('auto.Teams.6ce21a3c')">{{ myTeam.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Teams.ff22063c')">{{ myTeam.plan }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Teams.fc4a0583')">{{ myTeam.brand_name || '—' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Teams.2432b575')" :span="3">{{ myTeam.remark || '—' }}</el-descriptions-item>
        </el-descriptions>
        <el-descriptions v-if="myBilling" :column="3" border size="small" style="margin-bottom: 12px" :title="$t('auto.Teams.82aded13')">
          <el-descriptions-item :label="$t('auto.Teams.ce717abb')">{{ myBilling.subscription?.plan_code || '未订阅' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Teams.22e888c2')">{{ fmtTime(myBilling.subscription?.ends_at) }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Teams.d7550fbb')">{{ myBilling.subscription?.auto_renew ? '是' : '否' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Teams.4c8ac590')" :span="3">
            <template v-if="(myBilling.records || []).length">
              <el-tag
                v-for="r in myBilling.records.slice(0, 5)"
                :key="r.id"
                size="small"
                :type="r.status === 'paid' ? 'success' : r.status === 'pending' ? 'warning' : 'info'"
                style="margin-right: 6px"
              >¥{{ (r.amount_cents / 100).toFixed(2) }}（{{ { pending: '应收', paid: '已收', void: '作废' }[r.status] }}）</el-tag>
            </template>
            <span v-else class="muted">{{ $t('auto.Teams.a151e3de') }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </template>

      <el-table v-if="isPlatform" v-loading="loading" :data="rows" border size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" :label="$t('auto.Teams.6ce21a3c')" min-width="140" />
        <el-table-column :label="$t('auto.Teams.ff22063c')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="row.plan === 'enterprise' ? 'warning' : row.plan === 'pro' ? 'success' : 'info'">
              {{ row.plan }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Teams.3fea7ca7')" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Teams.81bf43d2')" width="100">
          <template #default="{ row }">
            <span v-if="row.owner_user_id">用户 #{{ row.owner_user_id }}</span>
            <span v-else class="muted">{{ $t('auto.Teams.7598f152') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Teams.eb9d06af')" width="90">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.member_count ?? 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" :label="$t('auto.Teams.2432b575')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('auto.Teams.eca37cb0')" width="160">
          <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.Teams.2b6bc0f2')" width="240" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" plain @click="openMembers(row)">{{ $t('auto.Teams.ab5dea29') }}</el-button>
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.Teams.95b351c8') }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.Teams.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        :page-size="size"
        :total="total"
        layout="total, prev, pager, next"
        small
        style="margin-top: 12px; justify-content: flex-end"
        @current-change="load"
      />
    </el-card>

    <!-- 新建/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" ::title="$t('auto.Teams.f9b0ed4d')" width="480px">
      <el-form label-width="100px">
        <el-form-item :label="$t('auto.Teams.6ce21a3c')">
          <el-input v-model="form.name" maxlength="100" :placeholder="$t('auto.Teams.6bc959ff')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Teams.ff22063c')">
          <el-radio-group v-model="form.plan">
            <el-radio value="free">free</el-radio>
            <el-radio value="pro">pro</el-radio>
            <el-radio value="enterprise">enterprise</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="editing" :label="$t('auto.Teams.3fea7ca7')">
          <el-radio-group v-model="form.status">
            <el-radio value="active">{{ $t('auto.Teams.7854b52a') }}</el-radio>
            <el-radio value="disabled">{{ $t('auto.Teams.5c56a889') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.Teams.81bf43d2')">
          <el-input v-model="form.owner_user_id" :placeholder="$t('auto.Teams.142d4e16')" style="width: 200px" />
        </el-form-item>
        <el-form-item :label="$t('auto.Teams.2432b575')">
          <el-input v-model="form.remark" maxlength="200" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.Teams.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="doSave">{{ $t('auto.Teams.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    <!-- 成员抽屉 -->
    <el-drawer v-model="membersVisible" :title="`成员管理 — ${current.name || ''}`" size="680px">
      <div style="margin-bottom: 12px; display: flex; gap: 8px">
        <el-input v-model="moveUserId" :placeholder="$t('auto.Teams.919d7a15')" style="width: 220px" size="small" />
        <el-button type="primary" size="small" :loading="moving" @click="doMove">{{ $t('auto.Teams.3b0ad53a') }}</el-button>
      </div>
      <div style="margin-bottom: 16px">
        <div style="margin-bottom: 8px; display: flex; align-items: center; gap: 8px">
          <span style="font-weight: 600">{{ $t('auto.Teams.fee0b675') }}</span>
          <el-button size="small" type="primary" plain @click="generateInvite">{{ $t('auto.Teams.1edb8b29') }}</el-button>
        </div>
        <el-table v-if="invites.length" :data="invites" border size="small">
          <el-table-column :label="$t('auto.Teams.fee0b675')" min-width="140">
            <template #default="{ row }"><code>{{ row.code }}</code></template>
          </el-table-column>
          <el-table-column :label="$t('auto.Teams.580f5201')" width="100">
            <template #default="{ row }">{{ row.used_count }}/{{ row.max_uses || '∞' }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.Teams.3fea7ca7')" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                {{ row.status === 'active' ? '可用' : '停用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Teams.2b6bc0f2')" width="100">
            <template #default="{ row }">
              <el-button v-if="row.status === 'active'" size="small" type="warning" plain @click="disableInvite(row)">{{ $t('auto.Teams.5c56a889') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="muted">{{ $t('auto.Teams.2ab7d3a2') }}</div>
      </div>
      <el-table v-loading="membersLoading" :data="members" border size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="username" :label="$t('auto.Teams.819767ad')" min-width="120" />
        <el-table-column :label="$t('auto.Teams.464f3d4e')" width="140">
          <template #default="{ row }">
            <el-tag v-if="row.role === 'team_admin'" type="warning" size="small">{{ $t('auto.Teams.00cd0a37') }}</el-tag>
            <el-tag v-else-if="row.role === 'admin'" type="danger" size="small">{{ $t('auto.Teams.fc79d47e') }}</el-tag>
            <el-tag v-else size="small">{{ row.role }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Teams.3fea7ca7')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled && row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.enabled && row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Teams.2b6bc0f2')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.role === 'user'"
              size="small"
              type="warning"
              plain
              @click="setRole(row, 'team_admin')"
            >{{ $t('auto.Teams.8fe415bc') }}</el-button>
            <el-button
              v-else-if="row.role === 'team_admin'"
              size="small"
              @click="setRole(row, 'user')"
            >{{ $t('auto.Teams.e7fda1d8') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { teamAdminApi, myTeamApi } from '@/api/teamAdmin'
import { http } from '@/utils/request'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const { t } = useI18n()

const rows = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const size = 20

const dialogVisible = ref(false)
const editing = ref(false)
const saving = ref(false)
const editingId = ref(0)
const form = ref({ name: '', plan: 'free', status: 'active', owner_user_id: '', remark: '' })

const membersVisible = ref(false)
const membersLoading = ref(false)
const members = ref([])
const current = ref({})
const moveUserId = ref('')
const moving = ref(false)

const fmtTime = (v) => (v ? String(v).replace('T', ' ').slice(0, 19) : '—')

// 多租户角色自适应（合并 MyTeam）：admin=平台全量；team_admin=仅本团队
const isPlatform = computed(() => userStore.role === 'admin')
const myTeam = ref(null)
const myBilling = ref(null)
const memberPage = ref(1)
const memberTotal = ref(0)
// 邀请码（管理员抽屉按团队 / 团队管理员本团队）
const invites = ref([])
const invitesLoading = ref(false)

const load = async () => {
  if (!isPlatform.value) {
    // 团队管理员：本团队详情 + 成员 + 订阅账单
    try {
      const res = await myTeamApi.detail()
      myTeam.value = (res?.data ?? res ?? {}).team ?? null
    } catch { myTeam.value = null }
    try {
      const res = await myTeamApi.members({ page: memberPage.value, size })
      const data = res?.data ?? res ?? {}
      members.value = data.list ?? []
      memberTotal.value = data.total ?? 0
    } catch { members.value = [] }
    try {
      const res = await http.get('/api/team/billing')
      myBilling.value = res?.data ?? res ?? null
    } catch { myBilling.value = null }
    await loadInvites()
    return
  }
  loading.value = true
  try {
    const res = await teamAdminApi.list({ page: page.value, size })
    const data = res?.data ?? res ?? {}
    rows.value = data.list ?? []
    total.value = data.total ?? 0
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editing.value = false
  editingId.value = 0
  form.value = { name: '', plan: 'free', status: 'active', owner_user_id: '', remark: '' }
  dialogVisible.value = true
}

const openEdit = (row) => {
  editing.value = true
  editingId.value = row.id
  form.value = {
    name: row.name,
    plan: row.plan || 'free',
    status: row.status || 'active',
    owner_user_id: row.owner_user_id ? String(row.owner_user_id) : '',
    remark: row.remark || ''
  }
  dialogVisible.value = true
}

const doSave = async () => {
  saving.value = true
  try {
    const payload = {
      name: form.value.name,
      plan: form.value.plan,
      remark: form.value.remark
    }
    if (form.value.owner_user_id) payload.owner_user_id = Number(form.value.owner_user_id)
    if (editing.value) {
      payload.status = form.value.status
      await teamAdminApi.update(editingId.value, payload)
    } else {
      await teamAdminApi.create(payload)
    }
    ElMessage.success(t('auto.Teams.f8dfedcd'))
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Teams.6de920b4'))
  } finally {
    saving.value = false
  }
}

const remove = async (row) => {
  await ElMessageBox.confirm(
    `确认删除团队「${row.name}」？仅空团队可删除（有成员时会被拒绝）；默认团队不可删除。`,
    t('auto.Teams.50eaf94d'),
    { type: 'warning' }
  )
  try {
    await teamAdminApi.remove(row.id)
    ElMessage.success(t('auto.Teams.5cc23262'))
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Teams.acf0664a'))
  }
}

const openMembers = (row) => {
  current.value = row
  membersVisible.value = true
  loadMembers()
  loadInvites()
}

const loadMembers = async () => {
  membersLoading.value = true
  try {
    const res = await teamAdminApi.members(current.value.id, { page: 1, size: 100 })
    const data = res?.data ?? res ?? {}
    members.value = data.list ?? []
  } finally {
    membersLoading.value = false
  }
}

// 邀请码：平台管理员按团队 / 团队管理员本团队
const loadInvites = async () => {
  invitesLoading.value = true
  try {
    const res = isPlatform.value && current.value.id
      ? await http.get(`/api/manage/teams/${current.value.id}/invite-codes`)
      : await http.get('/api/team/invite-codes')
    invites.value = (res?.data ?? res ?? {}).list ?? []
  } catch {
    invites.value = []
  } finally {
    invitesLoading.value = false
  }
}

const generateInvite = async () => {
  const teamID = isPlatform.value ? current.value.id : null
  if (isPlatform.value && !teamID) return
  const { value } = await ElMessageBox.prompt(
    t('auto.Teams.0e5c2c1d'), t('auto.Teams.1edb8b29'), {
      inputPlaceholder: '7',
      inputPattern: /^\d*$/,
      inputErrorMessage: t('auto.Teams.de66a286')
    }).catch(() => ({ value: null }))
  if (value === null) return
  const days = Number(value) || 0
  try {
    const res = isPlatform.value
      ? await http.post(`/api/manage/teams/${teamID}/invite-codes`, { expires_days: days })
      : await http.post('/api/team/invite-codes', { expires_days: days })
    const code = (res?.data ?? res ?? {})?.invite?.code
    ElMessage.success(`邀请码已生成：${code}`)
    await loadInvites()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Teams.7f7de8a2'))
  }
}

const disableInvite = async (row) => {
  try {
    if (isPlatform.value) {
      await http.post(`/api/manage/teams/${current.value.id}/invite-codes/${row.id}/disable`)
    } else {
      await http.post(`/api/team/invite-codes/${row.id}/disable`)
    }
    ElMessage.success(t('auto.Teams.69b0f684'))
    await loadInvites()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Teams.5fa802be'))
  }
}

const setRole = async (row, role) => {
  const label = role === 'team_admin' ? t('auto.Teams.8fe415bc') : t('auto.Teams.e1394dd3')
  await ElMessageBox.confirm(`确认将用户「${row.username}」${label}？`, t('auto.Teams.e83a256e'), { type: 'warning' })
  try {
    if (isPlatform.value) {
      await teamAdminApi.setMemberRole(current.value.id, row.id, role)
    } else {
      await myTeamApi.setMemberRole(row.id, role)
    }
    ElMessage.success(t('auto.Teams.1b92092d'))
    await loadMembers()
    if (!isPlatform.value) await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Teams.5fa802be'))
  }
}

const doMove = async () => {
  const uid = Number(moveUserId.value)
  if (!uid) {
    ElMessage.warning(t('auto.Teams.db670474'))
    return
  }
  moving.value = true
  try {
    await teamAdminApi.moveMember(current.value.id, uid)
    ElMessage.success(t('auto.Teams.eb55d76c'))
    moveUserId.value = ''
    await loadMembers()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Teams.c63bd074'))
  } finally {
    moving.value = false
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
.muted {
  color: #999;
  font-size: 12px;
}
</style>
