<template>
  <div class="asset-bundle-list">
    <el-card class="filter-bar" shadow="never">
      <el-form inline>
        <el-form-item :label="$t('auto.List.9699a50e')">
          <el-input v-model="filter.keyword" :placeholder="$t('auto.List.33957583')" clearable />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3fea7ca7')">
          <el-select v-model="filter.status" :placeholder="$t('auto.List.a8b0c204')" clearable>
            <el-option :label="$t('auto.List.22b4334f')" value="draft" />
            <el-option :label="$t('auto.List.7854b52a')" value="active" />
            <el-option :label="$t('auto.List.5c56a889')" value="inactive" />
            <el-option :label="$t('auto.List.2f51c18f')" value="archived" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.4705b884')">
          <el-select v-model="filter.scope" :placeholder="$t('auto.List.a8b0c204')" clearable>
            <el-option :label="$t('auto.List.3dc5185d')" value="private" />
            <el-option :label="$t('auto.List.fffd3ce7')" value="shared" />
            <el-option :label="$t('auto.List.f5c3795b')" value="official" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">{{ $t('auto.List.bee912d7') }}</el-button>
          <el-button type="success" @click="$router.push('/asset-bundle/playground')">{{ $t('auto.List.bef4c01d') }}</el-button>
          <el-button type="warning" @click="$router.push('/asset-bundle/merchant-new')">{{ $t('auto.List.0bc814df') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="hot-bar" shadow="never">
      <div class="hot-row">
        <span class="hot-title">{{ $t('auto.List.a07c2315') }}</span>
        <div v-if="enabledList.length" class="hot-list">
          <el-tag
            v-for="b in enabledList"
            :key="b.id"
            type="success"
            class="hot-tag"
            @click="$router.push(`/asset-bundle/merchant/${b.asset_id}`)"
          >
            {{ b.title }} ({{ b.asset_id }})
          </el-tag>
        </div>
        <span v-else class="hot-empty">{{ $t('auto.List.9198ae45') }}</span>
        <el-button size="small" link @click="fetchEnabled">{{ $t('auto.List.694fc5ef') }}</el-button>
      </div>
    </el-card>

    <el-table :data="list" v-loading="loading" border stripe>
      <el-table-column :label="$t('auto.List.41231061')" width="80">
        <template #default="{ row }">
          <el-image
            v-if="row.cover_image"
            :src="row.cover_image"
            :preview-src-list="[row.cover_image]"
            fit="cover"
            style="width: 64px; height: 48px; border-radius: 6px"
          />
          <div v-else class="cover-placeholder-bundle">{{ (row.title || '?').slice(0, 1) }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="asset_id" label="Asset ID" width="220" />
      <el-table-column prop="title" :label="$t('auto.List.32c65d8d')" min-width="200" />
      <el-table-column prop="author" :label="$t('auto.List.62cea749')" width="120" />
      <el-table-column prop="version" :label="$t('auto.List.fe2df04a')" width="80" />
      <el-table-column prop="industry" :label="$t('auto.List.6adafe27')" width="120" />
      <el-table-column :label="$t('auto.List.4705b884')" width="90">
        <template #default="{ row }">
          <el-tag size="small" :type="scopeTagType(row.scope)">{{ scopeLabel(row.scope) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.List.3fea7ca7')" width="90">
        <template #default="{ row }">
          <el-tag size="small" :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="use_count" :label="$t('auto.List.194df98d')" width="90" />
      <el-table-column prop="updated_at" :label="$t('auto.List.a001a226')" width="170">
        <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.List.2b6bc0f2')" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="primary" link @click="$router.push(`/asset-bundle/playground/${row.asset_id}`)">Playground</el-button>
          <el-button size="small" type="warning" link @click="$router.push(`/asset-bundle/merchant/${row.asset_id}`)">{{ $t('auto.List.9c878dab') }}</el-button>
          <el-button v-if="!hotEnabledIds.has(row.id)" size="small" type="success" link @click="handleEnable(row)">{{ $t('auto.List.44a43e51') }}</el-button>
          <el-button v-else size="small" type="info" link @click="handleDisable(row)">{{ $t('auto.List.0651f12c') }}</el-button>
          <el-button v-if="row.status === 'draft'" size="small" type="primary" link @click="handlePublish(row)">{{ $t('auto.List.83611abd') }}</el-button>
          <el-button size="small" type="danger" link @click="handleDelete(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="filter.page"
      v-model:page-size="filter.size"
      :total="total"
      layout="total, prev, pager, next, jumper"
      @current-change="fetchList"
      class="pager"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listBundles, listEnabledBundles, enableBundle, disableBundle,
  publishBundle, deleteBundle
} from '@/api/assetBundle'

const { t } = useI18n()

const filter = reactive({
  keyword: '',
  status: '',
  scope: '',
  page: 1,
  size: 20
})
const list = ref([])
const total = ref(0)
const loading = ref(false)
const enabledList = ref([])
const hotEnabledIds = ref(new Set());

const scopeLabel = (s) => ({ private: t('auto.List.3dc5185d'), shared: t('auto.List.fffd3ce7'), official: t('auto.List.f5c3795b') }[s] || s)
const scopeTagType = (s) => ({ private: 'info', shared: 'warning', official: 'success' }[s] || '')
const statusLabel = (s) => ({ draft: t('auto.List.22b4334f'), active: t('auto.List.7854b52a'), inactive: t('auto.List.5c56a889'), archived: t('auto.List.2f51c18f') }[s] || s)
const statusTagType = (s) => ({ draft: 'info', active: 'success', inactive: 'warning', archived: 'danger' }[s] || '')

const formatDate = (d) => {
  if (!d) return ''
  try {
    return new Date(d).toLocaleString('zh-CN', { hour12: false })
  } catch (e) {
    return d
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const resp = await listBundles({
      keyword: filter.keyword,
      status: filter.status,
      scope: filter.scope,
      page: filter.page,
      size: filter.size
    })
    const data = resp?.data || resp || {}
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    ElMessage.error(t('auto.List.06960ae4') + (e?.message || e))
  } finally {
    loading.value = false
  }
}

const fetchEnabled = async () => {
  try {
    const resp = await listEnabledBundles()
    const data = resp?.data || resp || {}
    const list = data.list || []
    enabledList.value = list
    hotEnabledIds.value = new Set(list.map(b => b.id))
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const handleEnable = async (row) => {
  try {
    await enableBundle(row.id)
    ElMessage.success(t('auto.List.8e38e05e'))
    fetchEnabled()
  } catch (e) {
    ElMessage.error(t('auto.List.1c78efaa') + (e?.message || e))
  }
}

const handleDisable = async (row) => {
  try {
    await disableBundle(row.id)
    ElMessage.success(t('auto.List.e491c037'))
    fetchEnabled()
  } catch (e) {
    ElMessage.error(t('auto.List.f8715c57') + (e?.message || e))
  }
}

const handlePublish = async (row) => {
  try {
    await ElMessageBox.confirm(`确认发布资产包「${row.title}」？`, t('auto.List.e83a256e'), { type: 'warning' })
    await publishBundle(row.id)
    row.status = 'active'
    ElMessage.success(t('auto.List.dca0c13b'))
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(t('auto.List.bd774adc') + (e?.message || e))
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除资产包「${row.title}」？此操作不可恢复。`, t('auto.List.653fc421'), { type: 'warning' })
    await deleteBundle(row.id)
    ElMessage.success(t('auto.List.5cc23262'))
    fetchList()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(t('auto.List.ad23f072') + (e?.message || e))
  }
}

onMounted(() => {
  fetchList()
  fetchEnabled()
})
</script>

<style scoped>
.asset-bundle-list {
  padding: 16px;
}
.filter-bar {
  margin-bottom: 12px;
}
.hot-bar {
  margin-bottom: 12px;
  background: linear-gradient(90deg, #fdf6ec 0%, #fff 100%);
}
.hot-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.hot-title {
  font-weight: 600;
  color: #e6a23c;
}
.hot-list {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.hot-tag {
  cursor: pointer;
}
.hot-empty {
  font-size: 12px;
  color: #909399;
}
.pager {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
.cover-placeholder-bundle {
  width: 64px;
  height: 48px;
  border-radius: 6px;
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
