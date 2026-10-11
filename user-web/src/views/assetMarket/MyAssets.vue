<template>
  <div class="my-assets">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.MyAssets.de43bfaf') }}</span>
          <div>
            <el-button type="primary" @click="$router.push('/asset-market')">{{ $t('auto.MyAssets.47a140b5') }}</el-button>
            <el-button @click="openCreate">{{ $t('auto.MyAssets.59153818') }}</el-button>
            <el-button @click="$router.push('/asset-market/sync-log')">{{ $t('auto.MyAssets.e1be449f') }}</el-button>
          </div>
        </div>
      </template>

      <el-form inline class="mb-12">
        <el-form-item :label="$t('auto.MyAssets.226b0912')">
          <el-select v-model="filter.asset_type" clearable :placeholder="$t('auto.MyAssets.a8b0c204')" @change="fetchList">
            <el-option :label="$t('auto.MyAssets.5ce3a11b')" value="agent_persona" />
            <el-option :label="$t('auto.MyAssets.655e227e')" value="sales_script" />
            <el-option :label="$t('auto.MyAssets.ec64559a')" value="ab_test_plan" />
            <el-option :label="$t('auto.MyAssets.228526dc')" value="marketing_workflow" />
            <el-option :label="$t('auto.MyAssets.112acd54')" value="industry_sop" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.26ca20b1')">
          <el-select v-model="filter.source" clearable :placeholder="$t('auto.MyAssets.a8b0c204')" @change="fetchList">
            <el-option :label="$t('auto.MyAssets.fb102c8f')" value="purchased" />
            <el-option :label="$t('auto.MyAssets.e75b8f8b')" value="manual" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filter.keyword" :placeholder="$t('auto.MyAssets.aa68f080')" clearable @keyup.enter="fetchList" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">{{ $t('auto.MyAssets.bee912d7') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="name" :label="$t('auto.MyAssets.d7ec2d3f')" min-width="160" />
        <el-table-column prop="asset_type" :label="$t('auto.MyAssets.226b0912')" width="140" />
        <el-table-column prop="industry" :label="$t('auto.MyAssets.6adafe27')" width="90" />
        <el-table-column prop="version" :label="$t('auto.MyAssets.fe2df04a')" width="90" />
        <el-table-column prop="source" :label="$t('auto.MyAssets.26ca20b1')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getSourceTagType(row.source)">
              {{ getSourceLabel(row.source) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.MyAssets.3fea7ca7')" width="90">
          <template #default="{ row }">
            <el-tag :type="getEnabledTagType(row.is_active ? 1 : 0)" size="small">
              {{ row.is_active ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="use_count" :label="$t('auto.MyAssets.194df98d')" width="100" />
        <el-table-column prop="synced_at" :label="$t('auto.MyAssets.88e76fc4')" width="170" />
        <el-table-column :label="$t('auto.MyAssets.2b6bc0f2')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.source === 'purchased' || row.source === 'synced'"
              link
              type="primary"
              @click="handleSync(row)"
            >{{ $t('auto.MyAssets.6a620e3c') }}</el-button>
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.MyAssets.607e7a4f') }}</el-button>
            <el-button
              v-if="row.source === 'purchased' || row.source === 'synced'"
              link
              type="warning"
              @click="handleReport(row)"
            >{{ $t('auto.MyAssets.28fcf95e') }}</el-button>
            <el-button link @click="toggle(row)">{{ row.is_active ? '停用' : '启用' }}</el-button>
            <el-button link type="danger" @click="remove(row)">{{ $t('auto.MyAssets.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="mt-12"
        v-model:current-page="filter.page"
        :page-size="filter.size"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="fetchList"
      />
    </el-card>

    <el-dialog v-model="createVisible" :title="$t('auto.MyAssets.59153818')" width="720px">
      <el-form label-width="100px">
        <el-form-item :label="$t('auto.MyAssets.d7ec2d3f')" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.226b0912')" required>
          <el-select v-model="form.asset_type" style="width: 100%">
            <el-option :label="$t('auto.MyAssets.5ce3a11b')" value="agent_persona" />
            <el-option :label="$t('auto.MyAssets.655e227e')" value="sales_script" />
            <el-option :label="$t('auto.MyAssets.ec64559a')" value="ab_test_plan" />
            <el-option :label="$t('auto.MyAssets.228526dc')" value="marketing_workflow" />
            <el-option :label="$t('auto.MyAssets.112acd54')" value="industry_sop" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.6adafe27')" required>
          <el-select v-model="form.industry" style="width: 100%">
            <el-option v-for="i in industries" :key="i" :label="i" :value="i" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.e77523d0')" required>
          <el-input v-model="form.dataText" type="textarea" :rows="12" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ $t('auto.MyAssets.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitCreate">{{ $t('auto.MyAssets.d9ac9228') }}</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="detailVisible" :title="$t('auto.MyAssets.f508a46b')" size="50%">
      <pre class="json-preview">{{ detailText }}</pre>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  listLocalAssets,
  getLocalAsset,
  createLocalAsset,
  deleteLocalAsset,
  toggleLocalAsset,
  syncAsset,
  reportUsage
} from '@/api/assetMarket'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSourceLabel, getSourceTagType } from '@/constants/source';
import { getEnabledTagType } from '@/constants/enabled'
import { INDUSTRY_VALUES } from '@/constants/industry'

const industries = INDUSTRY_VALUES
const { t } = useI18n()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const filter = ref({ asset_type: '', source: '', keyword: '', page: 1, size: 20 })
const createVisible = ref(false)
const detailVisible = ref(false)
const detailText = ref('')
const form = ref({ name: '', asset_type: 'agent_persona', industry: t('auto.MyAssets.1b39e345'), dataText: '' })

const sourceLabel = (s) => getSourceLabel(s)

const fetchList = async () => {
  loading.value = true
  try {
    const resp = await listLocalAssets(filter.value)
    const data = resp?.data || resp || {}
    list.value = Array.isArray(data?.list) ? data.list : (Array.isArray(data) ? data : [])
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const handleSync = async (row) => {
  await syncAsset({ asset_id: row.asset_id })
  ElMessage.success(t('auto.MyAssets.52b85c91'))
  fetchList()
}

const handleReport = async (row) => {
  await reportUsage({ asset_id: row.asset_id })
  ElMessage.success(t('auto.MyAssets.d4d5c12f'))
  fetchList()
}

const viewDetail = async (row) => {
  const resp = await getLocalAsset(row.id)
  const data = resp?.data || resp
  detailText.value = JSON.stringify(data, null, 2)
  detailVisible.value = true
}

const toggle = async (row) => {
  await toggleLocalAsset(row.id, !row.is_active)
  ElMessage.success(t('auto.MyAssets.33130f5c'))
  fetchList()
}

const remove = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除「${row.name}」？`, t('auto.MyAssets.2f4aaddd'), { type: 'warning' })
  } catch {
    return
  }
  await deleteLocalAsset(row.id)
  ElMessage.success(t('auto.MyAssets.5cc23262'))
  fetchList()
}

const openCreate = () => {
  form.value = {
    name: '',
    asset_type: 'agent_persona',
    industry: t('auto.MyAssets.1b39e345'),
    dataText: JSON.stringify(
      {
        schema_version: '1.0',
        asset_type: 'agent_persona',
        name: t('auto.MyAssets.f037ec10'),
        industry: t('auto.MyAssets.1b39e345'),
        version: '1.0.0',
        system_prompt: t('auto.MyAssets.1ee51d9e'),
        persona: { tone: t('auto.MyAssets.a7ad35ac'), expertise: [t('auto.MyAssets.a74ccc6e')] }
      },
      null,
      2
    )
  }
  createVisible.value = true
}

const submitCreate = async () => {
  let data
  try {
    data = JSON.parse(form.value.dataText)
  } catch {
    ElMessage.error('JSON 格式错误')
    return
  }
  data.asset_type = form.value.asset_type
  data.industry = form.value.industry
  data.name = form.value.name || data.name
  await createLocalAsset({
    name: form.value.name,
    asset_type: form.value.asset_type,
    industry: form.value.industry,
    data
  })
  ElMessage.success(t('auto.MyAssets.04a691b3'))
  createVisible.value = false
  fetchList()
}

onMounted(fetchList)
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
.json-preview {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 6px;
  font-size: 12px;
  white-space: pre-wrap;
}
</style>
