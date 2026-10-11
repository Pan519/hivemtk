<template>
  <div class="conflicts-page">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Conflicts.ca03e24e') }}</span>
          <div class="header-actions">
            <el-tag v-if="pendingCount > 0" type="danger" effect="dark">待处理 {{ pendingCount }}</el-tag>
            <el-button type="primary" size="small" :loading="loading" @click="fetchData">{{ $t('auto.Conflicts.694fc5ef') }}</el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="pendingCount > 0"
        type="warning"
        :closable="false"
        show-icon
        :title="$t('auto.Conflicts.3a25041f')"
        :description="$t('auto.Conflicts.e2ff182b')"
        style="margin-bottom: 16px"
      />

      <el-table
        v-loading="loading"
        :data="conflicts"
        border
        stripe
        :empty-text="$t('auto.Conflicts.cbd5e3aa')"
        style="width: 100%"
      >
        <el-table-column prop="identity_type" :label="$t('auto.Conflicts.dbe8853d')" width="130">
          <template #default="{ row }">
            <el-tag :type="typeTag(row.identity_type)">{{ channelLabel(row.identity_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="identity_value" :label="$t('auto.Conflicts.af77ecf3')" min-width="160" />
        <el-table-column :label="$t('auto.Conflicts.6d717038')" min-width="200">
          <template #default="{ row }">
            <el-tag
              v-for="(cid, idx) in row.customer_ids"
              :key="cid"
              :type="idx === 0 ? 'success' : 'info'"
              size="small"
              style="margin-right: 6px"
            >
              {{ short(cid) }}
              <span v-if="idx === 0">{{ $t('auto.Conflicts.11c01014') }}</span>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.Conflicts.3fea7ca7')" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === 'resolved' ? 'success' : 'danger'">
              {{ row.status === 'resolved' ? '已解决' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.Conflicts.2fe5f3d8')" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.Conflicts.2b6bc0f2')" width="240" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status !== 'resolved'"
              type="danger"
              size="small"
              :loading="row._merging"
              @click="handleMerge(row)"
            >{{ $t('auto.Conflicts.e5d8537a') }}</el-button>
            <el-button
              v-if="row.status !== 'resolved'"
              type="info"
              size="small"
              :loading="row._ignoring"
              @click="handleIgnore(row)"
            >{{ $t('auto.Conflicts.08aecfb4') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      
      <el-dialog v-model="mergeDialog" :title="$t('auto.Conflicts.0ca23eab')" width="520px">
        <div v-if="selected">
          <p>{{ $t('auto.Conflicts.c5348675') }}<b>{{ channelLabel(selected.identity_type) }}</b> = <b>{{ selected.identity_value }}</b></p>
          <p>{{ $t('auto.Conflicts.33f9c5b1') }}</p>
          <el-table :data="mergeRows" border size="small">
            <el-table-column prop="role" :label="$t('auto.Conflicts.464f3d4e')" width="90" />
            <el-table-column prop="id" :label="$t('auto.Conflicts.3f61dec5')" />
          </el-table>
          <el-alert type="info" :closable="false" style="margin-top: 12px"
            >{{ $t('auto.Conflicts.c592be53') }}</el-alert>
        </div>
        <template #footer>
          <el-button @click="mergeDialog = false">{{ $t('auto.Conflicts.625fb26b') }}</el-button>
          <el-button type="danger" :loading="submitting" @click="confirmMerge">{{ $t('auto.Conflicts.8ba6082a') }}</el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { listConflicts, mergeOneID, resolveConflict } from '@/api/oneid'
import { ElMessage, ElMessageBox } from 'element-plus'

const { t } = useI18n()

const conflicts = ref([])
const loading = ref(false)
const mergeDialog = ref(false)
const submitting = ref(false)
const selected = ref(null)

const pendingCount = computed(() => conflicts.value.filter(c => c.status !== 'resolved').length)
const mergeRows = computed(() => {
  if (!selected.value) return []
  const ids = selected.value.customer_ids || []
  return ids.map((id, idx) => ({ role: idx === 0 ? t('auto.Conflicts.b31d4768') : t('auto.Conflicts.1fc6cbbd'), id }))
})

function channelLabel(t) {
  const map = {
    phone: t('auto.Conflicts.8098e2b4'), email: t('auto.Conflicts.3bc5e602'), wechat_open_id: t('auto.Conflicts.cfbf6f4c'),
    douyin_open_id: t('auto.Conflicts.277a984f'), xiaohongshu_id: t('auto.Conflicts.ae26bb04')
  }
  return map[t] || t
}
function typeTag() { return 'warning' }
function short(id) { return id ? id.slice(0, 8) : '-' }
function formatTime(t) {
  if (!t) return '-'
  const d = new Date(t)
  return isNaN(d) ? String(t) : d.toLocaleString()
}

async function fetchData() {
  loading.value = true
  try {
    const res = await listConflicts({ status: 'pending' })
    // 拦截器已拆包：conflicts 直接回 {list,total,...}，按信封取 .data 只会读到 undefined
    const data = res?.data?.data || res?.data || res?.list || []
    conflicts.value = Array.isArray(data) ? data.map(c => ({ ...c, customer_ids: c.customer_ids || [] })) : []
  } catch (e) {
    ElMessage.error(t('auto.Conflicts.913829cf') + (e?.message || e))
  } finally {
    loading.value = false
  }
}

function handleMerge(row) {
  selected.value = row
  mergeDialog.value = true
}

async function confirmMerge() {
  const row = selected.value
  if (!row || !row.customer_ids || row.customer_ids.length < 2) return
  const [primary, ...secondary] = row.customer_ids
  try {
    await ElMessageBox.confirm(
      `即将把以下客户合并为主客户 ${primary}：\n从属客户 ${secondary[0]} 将被永久删除，其会话与事件归属迁移至主客户。\n该操作不可撤销，确认继续？`,
      t('auto.Conflicts.f605d2a8'),
      { type: 'warning', confirmButtonText: t('auto.Conflicts.8ba6082a'), cancelButtonText: t('auto.Conflicts.625fb26b') }
    )
  } catch {
    return
  }
  submitting.value = true
  row._merging = true
  try {
    await mergeOneID({ primary_id: primary, secondary_id: secondary[0] })
    await resolveConflict(row.conflict_id, { action: 'merge', primary_id: primary, secondary_id: secondary[0] })
    ElMessage.success(t('auto.Conflicts.3018960b'))
    mergeDialog.value = false
    await fetchData()
  } catch (e) {
    ElMessage.error(t('auto.Conflicts.3f9aa174') + (e?.message || e))
  } finally {
    row._merging = false
    submitting.value = false
  }
}

async function handleIgnore(row) {
  try {
    await ElMessageBox.confirm(t('auto.Conflicts.90f4e850'), t('auto.Conflicts.08aecfb4'), { type: 'warning' })
  } catch {
    return
  }
  row._ignoring = true
  try {
    await resolveConflict(row.conflict_id, { action: 'ignore' })
    ElMessage.success(t('auto.Conflicts.17f2ba0c'))
    await fetchData()
  } catch (e) {
    ElMessage.error(t('auto.Conflicts.79585f92') + (e?.message || e))
  } finally {
    row._ignoring = false
  }
}

onMounted(fetchData)
</script>

<style scoped>
.conflicts-page { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; align-items: center; gap: 12px; }
</style>
