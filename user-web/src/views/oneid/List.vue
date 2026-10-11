<template>
  <div class="oneid-list-container">
    <el-card class="header-card">
      <div class="header-row">
        <div class="header-text">
          <h2>{{ $t('auto.List.18.81b236') }}</h2>
          <p class="subtitle">{{ $t('auto.List.19.75ad32') }}</p>
        </div>
        <div class="header-actions">
          <el-input
            v-model="keyword"
            ::placeholder="$t('auto.List.1.ec8f4c')"
            clearable
            style="width: 320px; margin-right: 12px"
            @keyup.enter="handleSearch"
          />
          <el-button type="primary" @click="handleSearch">{{ $t('搜索') }}</el-button>
          <el-button type="success" @click="resolveDialogVisible = true">
            <el-icon><Plus /></el-icon>解析/创建 OneID
          </el-button>
        </div>
      </div>
      <el-alert
        type="info"
        :closable="false"
        ::title="$t('auto.List.2.3fbe7b')"
      />
    </el-card>

    
    <el-row :gutter="16" class="stats-row" v-loading="statsLoading">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(79,70,229,.1); color: #4F46E5">
            <el-icon :size="22"><UserFilled /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.List.20.6aa66e') }}</div>
            <div class="stat-value">{{ stats.total || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(16,185,129,.1); color: #10B981">
            <el-icon :size="22"><Cellphone /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.List.21.660b48') }}</div>
            <div class="stat-value">{{ stats.withPhone || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(245,158,11,.1); color: #F59E0B">
            <el-icon :size="22"><Message /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.List.22.11b831') }}</div>
            <div class="stat-value">{{ stats.withEmail || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(99,102,241,.1); color: #6366F1">
            <el-icon :size="22"><Connection /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.List.23.37c52b') }}</div>
            <div class="stat-value">{{ stats.multiIdentity || 0 }}</div>
            <div class="stat-extra">占比 {{ stats.multiRate || 0 }}%</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="table-card">
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="unified_id" label="UnifiedID" min-width="200">
          <template #default="{ row }">
            <el-tag size="small">{{ row.unified_id }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone" :label="$t('auto.List.3.8098e2')" min-width="140" />
        <el-table-column prop="email" :label="$t('auto.List.4.3bc5e6')" min-width="200" />
        <el-table-column :label="$t('auto.List.5.1f63f9')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getIdentityCount(row) >= 2 ? 'success' : 'info'">
              {{ getIdentityCount(row) }} 个
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="rfm_score" :label="$t('auto.List.6.09834d')" width="100" />
        <el-table-column prop="churn_risk" :label="$t('auto.List.7.4b6cdc')" width="120">
          <template #default="{ row }">
            <el-tag :type="riskType(row.churn_risk)">{{ row.churn_risk }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.List.8.eca37c')" min-width="180" />
        <el-table-column :label="$t('auto.List.9.2b6bc0')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="showIdentities(row)">{{ $t('auto.List.24.52d523') }}</el-button>
            <el-button size="small" type="primary" @click="showLink(row)">{{ $t('auto.List.25.a8c3ad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next, jumper"
        @current-change="loadList"
        style="margin-top: 16px; text-align: right"
      />
    </el-card>

    
    <el-dialog v-model="identitiesDialogVisible" :title="$t('auto.List.10.0fb633')" width="640px">
      <div v-if="selectedCustomer">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="UnifiedID">
            <el-tag>{{ selectedCustomer.unified_id }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.3.8098e2')">
            {{ selectedCustomer.phone || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.4.3bc5e6')">
            {{ selectedCustomer.email || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.11.147f03')">
            {{ selectedCustomer.wechat_open_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.12.4b74f4')">
            {{ selectedCustomer.douyin_open_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.13.95daed')">
            {{ selectedCustomer.xiaohongshu_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.6.09834d')">
            {{ selectedCustomer.rfm_score }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.7.4b6cdc')">
            <el-tag :type="riskType(selectedCustomer.churn_risk)">
              {{ selectedCustomer.churn_risk }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    
    <el-dialog v-model="resolveDialogVisible" :title="$t('auto.List.14.fd23eb')" width="640px">
      <el-form :model="resolveForm" label-width="100px">
        <el-form-item :label="$t('auto.List.3.8098e2')">
          <el-input v-model="resolveForm.phone" :placeholder="$t('auto.List.15.d4769d')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.4.3bc5e6')">
          <el-input v-model="resolveForm.email" :placeholder="$t('auto.List.16.5e6cb0')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.11.147f03')">
          <el-input v-model="resolveForm.wechat_open_id" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.12.4b74f4')">
          <el-input v-model="resolveForm.douyin_open_id" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resolveDialogVisible = false">{{ $t('auto.List.26.625fb2') }}</el-button>
        <el-button type="primary" @click="handleResolve">{{ $t('auto.List.27.488672') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="linkDialogVisible" :title="$t('auto.List.17.8772fc')" width="640px">
      <el-alert
        v-if="linkForm.customerId"
        :title="`为客户 ${linkForm.customerId} 链接新身份`"
        type="info"
        :closable="false"
        style="margin-bottom: 16px"
      />
      <el-form :model="linkForm" label-width="100px">
        <el-form-item :label="$t('auto.List.3.8098e2')">
          <el-input v-model="linkForm.phone" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.4.3bc5e6')">
          <el-input v-model="linkForm.email" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.11.147f03')">
          <el-input v-model="linkForm.wechat_open_id" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.12.4b74f4')">
          <el-input v-model="linkForm.douyin_open_id" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="linkDialogVisible = false">{{ $t('auto.List.26.625fb2') }}</el-button>
        <el-button type="primary" @click="handleLink">{{ $t('auto.List.28.bfe68d') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, UserFilled, Cellphone, Message, Connection } from '@element-plus/icons-vue'
import { listOneID, resolveIdentity, getIdentityMappings, linkIdentity, getOneIDStats } from '@/api/oneid'

const { t } = useI18n()

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const loading = ref(false)
const statsLoading = ref(false)
const stats = ref({ total: 0, withPhone: 0, withEmail: 0, multiIdentity: 0, multiRate: 0 })
const identitiesDialogVisible = ref(false)
const resolveDialogVisible = ref(false)
const linkDialogVisible = ref(false)
const selectedCustomer = ref(null)
const resolveForm = ref({ phone: '', email: '', wechat_open_id: '', douyin_open_id: '' })
const linkForm = ref({ customerId: '', phone: '', email: '', wechat_open_id: '', douyin_open_id: '' })

const getIdentityCount = (row) => {
  let c = 0
  if (row.phone) c++
  if (row.email) c++
  if (row.wechat_open_id) c++
  if (row.douyin_open_id) c++
  if (row.xiaohongshu_id) c++
  return c
}

const riskType = (risk) => {
  if (risk === 'high') return 'danger'
  if (risk === 'medium') return 'warning'
  return 'success'
}

const loadStats = async () => {
  statsLoading.value = true
  try {
    const res = await getOneIDStats()
    const s = res || {}
    stats.value = {
      total: s.total ?? 0,
      withPhone: s.with_phone ?? s.withPhone ?? 0,
      withEmail: s.with_email ?? s.withEmail ?? 0,
      multiIdentity: s.multi_identity ?? s.multiIdentity ?? 0,
      multiRate: s.total ? Math.round(((s.multi_identity ?? s.multiIdentity ?? 0) * 1000) / s.total) / 10 : 0
    }
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ } finally {
    statsLoading.value = false
  }
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await listOneID({ page: page.value, page_size: pageSize.value, keyword: keyword.value })
    list.value = res?.list || [];
    total.value = res?.total || 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  loadList()
}

const showIdentities = async (row) => {
  try {
    const res = await getIdentityMappings(row.id)
    selectedCustomer.value = res?.customer
    identitiesDialogVisible.value = true
  } catch (e) {
    ElMessage.error(t('auto.List.29.5b7318') + e.message)
  }
}

const showLink = (row) => {
  linkForm.value = {
    customerId: row.id,
    phone: '',
    email: '',
    wechat_open_id: '',
    douyin_open_id: ''
  }
  linkDialogVisible.value = true
}

const handleResolve = async () => {
  try {
    await resolveIdentity(resolveForm.value)
    ElMessage.success(i18n.global.t(t('auto.List.30.53de01')))
    resolveDialogVisible.value = false
    loadList()
  } catch (e) {
    ElMessage.error(t('auto.List.31.4d1bf8') + e.message)
  }
}

  const handleLink = async () => {
    try {
      const { customerId, ...identifiers } = linkForm.value
      if (!customerId) {
        ElMessage.warning(i18n.global.t(t('auto.List.32.ba1dcc')))
        return
      }
      if (!identifiers.phone && !identifiers.email && !identifiers.wechat_open_id && !identifiers.douyin_open_id) {
        ElMessage.warning(i18n.global.t(t('auto.List.33.886f28')))
        return
      }
      await linkIdentity(customerId, identifiers)
      ElMessage.success(i18n.global.t(t('auto.List.34.a2391c')))
      linkDialogVisible.value = false
      loadList()
    } catch (e) {
      ElMessage.error(t('auto.List.35.15ecaa') + e.message)
    }
  }

onMounted(() => {
  loadList()
  loadStats()
})
</script>

<style scoped>
.oneid-list-container { padding: 0; }
.header-card { margin-bottom: 16px; }
.header-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; flex-wrap: wrap; gap: 12px; }
.header-text { flex: 1; min-width: 200px; }
.header-text h2 { margin: 0; font-size: 20px; font-weight: 600; }
.header-text .subtitle { margin: 4px 0 0; font-size: 13px; color: #909399; }
.header-actions { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.table-card { background: #fff; }

.stats-row { margin-bottom: 16px; }
.stats-row .stat-card { display: flex; align-items: center; gap: 12px; }
.stats-row .stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stats-row .stat-body { flex: 1; }
.stats-row .stat-label { font-size: 13px; color: #909399; }
.stats-row .stat-value { font-size: 22px; font-weight: 600; color: #303133; line-height: 1.2; margin-top: 4px; }
.stats-row .stat-extra { font-size: 12px; color: #6366F1; margin-top: 4px; }
</style>
