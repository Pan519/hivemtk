<template>
  <div class="my-assets">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.MyAssets.22.de43bf') }}</span>
          <div>
            <el-button type="primary" @click="$router.push('/asset-market')">{{ $t('auto.MyAssets.23.47a140') }}</el-button>
            <el-button @click="openCreate">{{ $t('auto.MyAssets.19.591538') }}</el-button>
            <el-button @click="$router.push('/asset-market/sync-log')">{{ $t('auto.MyAssets.24.e1be44') }}</el-button>
          </div>
        </div>
      </template>

      <el-form inline class="mb-12">
        <el-form-item :label="$t('auto.MyAssets.1.226b09')">
          <el-select v-model="filter.asset_type" clearable :placeholder="$t('auto.MyAssets.2.a8b0c2')" @change="fetchList">
            <el-option :label="$t('auto.MyAssets.3.5ce3a1')" value="agent_persona" />
            <el-option :label="$t('auto.MyAssets.4.655e22')" value="sales_script" />
            <el-option :label="$t('auto.MyAssets.5.ec6455')" value="ab_test_plan" />
            <el-option :label="$t('auto.MyAssets.6.228526')" value="marketing_workflow" />
            <el-option :label="$t('auto.MyAssets.7.112acd')" value="industry_sop" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.8.26ca20')">
          <el-select v-model="filter.source" clearable :placeholder="$t('auto.MyAssets.2.a8b0c2')" @change="fetchList">
            <el-option :label="$t('auto.MyAssets.9.fb102c')" value="purchased" />
            <el-option :label="$t('auto.MyAssets.10.e75b8f')" value="manual" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filter.keyword" :placeholder="$t('auto.MyAssets.11.aa68f0')" clearable @keyup.enter="fetchList" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">{{ $t('auto.MyAssets.25.bee912') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="name" :label="$t('auto.MyAssets.12.d7ec2d')" min-width="160" />
        <el-table-column prop="asset_type" :label="$t('auto.MyAssets.1.226b09')" width="140" />
        <el-table-column prop="industry" :label="$t('auto.MyAssets.13.6adafe')" width="90" />
        <el-table-column prop="version" :label="$t('auto.MyAssets.14.fe2df0')" width="90" />
        <el-table-column prop="source" :label="$t('auto.MyAssets.8.26ca20')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="getSourceTagType(row.source)">
              {{ getSourceLabel(row.source) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.MyAssets.15.3fea7c')" width="90">
          <template #default="{ row }">
            <el-tag :type="getEnabledTagType(row.is_active ? 1 : 0)" size="small">
              {{ row.is_active ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="use_count" :label="$t('auto.MyAssets.16.194df9')" width="100" />
        <el-table-column prop="synced_at" :label="$t('auto.MyAssets.17.88e76f')" width="170" />
        <el-table-column :label="$t('auto.MyAssets.18.2b6bc0')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.source === 'purchased' || row.source === 'synced'"
              link
              type="primary"
              @click="handleSync(row)"
            >{{ $t('auto.MyAssets.26.6a620e') }}</el-button>
            <el-button link type="primary" @click="viewDetail(row)">{{ $t('auto.MyAssets.27.607e7a') }}</el-button>
            <el-button
              v-if="row.source === 'purchased' || row.source === 'synced'"
              link
              type="warning"
              @click="handleReport(row)"
            >{{ $t('auto.MyAssets.28.28fcf9') }}</el-button>
            <el-button link @click="toggle(row)">{{ row.is_active ? '停用' : '启用' }}</el-button>
            <el-button link type="danger" @click="remove(row)">{{ $t('auto.MyAssets.29.2f4aad') }}</el-button>
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

    <el-dialog v-model="createVisible" :title="$t('auto.MyAssets.19.591538')" width="720px">
      <el-form label-width="100px">
        <el-form-item :label="$t('auto.MyAssets.12.d7ec2d')" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.1.226b09')" required>
          <el-select v-model="form.asset_type" style="width: 100%">
            <el-option :label="$t('auto.MyAssets.3.5ce3a1')" value="agent_persona" />
            <el-option :label="$t('auto.MyAssets.4.655e22')" value="sales_script" />
            <el-option :label="$t('auto.MyAssets.5.ec6455')" value="ab_test_plan" />
            <el-option :label="$t('auto.MyAssets.6.228526')" value="marketing_workflow" />
            <el-option :label="$t('auto.MyAssets.7.112acd')" value="industry_sop" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.13.6adafe')" required>
          <el-select v-model="form.industry" style="width: 100%">
            <el-option v-for="i in industries" :key="i" :label="i" :value="i" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.MyAssets.20.e77523')" required>
          <el-input v-model="form.dataText" type="textarea" :rows="12" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ $t('auto.MyAssets.30.625fb2') }}</el-button>
        <el-button type="primary" @click="submitCreate">{{ $t('auto.MyAssets.31.d9ac92') }}</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="detailVisible" :title="$t('auto.MyAssets.21.f508a4')" size="50%">
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
const form = ref({ name: '', asset_type: 'agent_persona', industry: t('auto.MyAssets.32.1b39e3'), dataText: '' })

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
  ElMessage.success(t('auto.MyAssets.33.52b85c'))
  fetchList()
}

const handleReport = async (row) => {
  await reportUsage({ asset_id: row.asset_id })
  ElMessage.success(t('auto.MyAssets.34.d4d5c1'))
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
  ElMessage.success(t('auto.MyAssets.35.33130f'))
  fetchList()
}

const remove = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除「${row.name}」？`, t('auto.MyAssets.29.2f4aad'), { type: 'warning' })
  } catch {
    return
  }
  await deleteLocalAsset(row.id)
  ElMessage.success(t('auto.MyAssets.36.5cc232'))
  fetchList()
}

const openCreate = () => {
  form.value = {
    name: '',
    asset_type: 'agent_persona',
    industry: t('auto.MyAssets.32.1b39e3'),
    dataText: JSON.stringify(
      {
        schema_version: '1.0',
        asset_type: 'agent_persona',
        name: t('auto.MyAssets.37.f037ec'),
        industry: t('auto.MyAssets.32.1b39e3'),
        version: '1.0.0',
        system_prompt: t('auto.MyAssets.38.1ee51d'),
        persona: { tone: t('auto.MyAssets.39.a7ad35'), expertise: [t('auto.MyAssets.40.a74ccc')] }
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
  ElMessage.success(t('auto.MyAssets.41.04a691'))
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
