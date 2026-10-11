<template>
  <div class="external-import-page">
    <el-alert
      ::title="$t('auto.ExternalImport.1.7c872e')"
      type="info"
      :closable="false"
      show-icon
    >
      <template #default>
        通过 API Token 鉴权,接入外部系统(飞书 / Notion / 钉钉 / 自有 CRM)的文档数据。
        支持两种模式:1) <strong>{{ $t('异步任务') }}</strong>{{ $t('auto.ExternalImport.22.d4c3de') }}<strong>{{ $t('同步模式') }}</strong>(sync=true)立即返回处理结果。
      </template>
    </el-alert>

    <el-row :gutter="16" style="margin-top: 16px">
      
      <el-col :span="14">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.ExternalImport.23.756331') }}</span>
              <el-tag :type="testMode ? 'warning' : 'success'" size="small">
                {{ testMode ? '同步模式' : '异步模式' }}
              </el-tag>
            </div>
          </template>
          <el-form :model="form" label-width="100px">
            <el-form-item :label="$t('auto.ExternalImport.2.c11322')">
              <el-radio-group v-model="form.source" @change="handleSourceChange">
                <el-radio-button label="custom">{{ $t('auto.ExternalImport.24.fb20b6') }}</el-radio-button>
                <el-radio-button label="feishu">{{ $t('auto.ExternalImport.25.7714e5') }}</el-radio-button>
                <el-radio-button label="notion">Notion</el-radio-button>
                <el-radio-button label="dingtalk">{{ $t('auto.ExternalImport.26.4a0e91') }}</el-radio-button>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="$t('auto.ExternalImport.3.a01543')" required>
              <el-select v-model="form.product_id" :placeholder="$t('auto.ExternalImport.4.cbf1cf')" style="width: 100%">
                <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
              </el-select>
            </el-form-item>
            <el-form-item v-if="form.source === 'feishu'" :label="$t('auto.ExternalImport.5.5fe7ed')">
              <el-input v-model="form.feishu_doc_id" :placeholder="$t('auto.ExternalImport.6.288d39')" />
            </el-form-item>
            <el-form-item v-if="form.source === 'notion'" label="Notion PageID">
              <el-input v-model="form.notion_page_id" :placeholder="$t('auto.ExternalImport.7.fec57c')" />
            </el-form-item>
            <el-form-item v-if="form.source === 'dingtalk'" :label="$t('auto.ExternalImport.8.6918fd')">
              <el-input v-model="form.dingtalk_doc_id" :placeholder="$t('auto.ExternalImport.9.c6e794')" />
            </el-form-item>
            <el-form-item v-if="form.source === 'custom'" :label="$t('auto.ExternalImport.10.c35ee6')" required>
              <el-input
                v-model="form.itemsJson"
                type="textarea"
                :rows="10"
                placeholder='JSON 数组,例如:&#10;[{"title": "FAQ1", "content": "...", "category": "售后", "tags": ["退款"]}]'
              />
            </el-form-item>
            <el-form-item>
              <el-checkbox v-model="testMode">{{ $t('auto.ExternalImport.27.d482e5') }}</el-checkbox>
              <el-button :icon="Document" @click="loadTemplate" style="margin-left: 8px">{{ $t('auto.ExternalImport.28.5806bb') }}</el-button>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Promotion" :loading="importing" :disabled="!canSubmit" @click="handleImport" style="width: 100%">
                提交导入
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        
        <el-card v-if="lastResponse" style="margin-top: 16px">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.ExternalImport.29.ac24ff') }}</span>
              <el-button link @click="lastResponse = null">{{ $t('auto.ExternalImport.30.b15d91') }}</el-button>
            </div>
          </template>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item :label="$t('auto.ExternalImport.11.7cb78f')">{{ lastResponse.job_no }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.ExternalImport.12.3fea7c')">
              <el-tag :type="statusType(lastResponse.status)">{{ statusLabel(lastResponse.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.ExternalImport.13.9ed7d3')">{{ lastResponse.total }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.ExternalImport.14.330363')">{{ lastResponse.accepted }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.ExternalImport.15.acd5cb')">{{ lastResponse.failed_items }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.ExternalImport.16.8b5a24')">
              <el-tag :type="lastResponse.async ? 'warning' : 'success'" size="small">
                {{ lastResponse.async ? '异步处理中' : '同步完成' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
          <div v-if="lastResponse.document_ids && lastResponse.document_ids.length" style="margin-top: 12px">
            <el-text size="small">{{ $t('auto.ExternalImport.31.a3bc87') }}</el-text>
            <el-tag v-for="id in lastResponse.document_ids" :key="id" size="small" style="margin-left: 4px">#{{ id }}</el-tag>
          </div>
          <el-alert
            v-if="lastResponse.errors && lastResponse.errors.length"
            type="error"
            :closable="false"
            style="margin-top: 12px"
          >
            <template #title>失败明细 ({{ lastResponse.errors.length }})</template>
            <ul style="margin: 0; padding-left: 20px">
              <li v-for="(e, idx) in lastResponse.errors" :key="idx">{{ e }}</li>
            </ul>
          </el-alert>
        </el-card>
      </el-col>

      
      <el-col :span="10">
        <el-card>
          <template #header>
            <span>API Token</span>
          </template>
          <el-form label-width="80px">
            <el-form-item label="Token">
              <el-input
                v-model="apiToken"
                type="password"
                :placeholder="$t('auto.ExternalImport.17.2457a9')"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-text size="small" type="info">{{ $t('auto.ExternalImport.32.6e8aa1') }}<code>X-Knowledge-Token</code>{{ $t('auto.ExternalImport.33.999851') }}</el-text>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card style="margin-top: 16px">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.ExternalImport.34.2762e0') }}</span>
              <el-button link :icon="Refresh" @click="loadJobs">{{ $t('auto.ExternalImport.35.694fc5') }}</el-button>
            </div>
          </template>
          <el-select v-model="jobFilterProduct" :placeholder="$t('auto.ExternalImport.18.6437ee')" clearable style="width: 100%; margin-bottom: 8px" @change="loadJobs">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
          <el-table :data="jobs" v-loading="jobsLoading" max-height="380" size="small">
            <el-table-column prop="job_no" :label="$t('auto.ExternalImport.11.7cb78f')" width="160" />
            <el-table-column prop="source" :label="$t('auto.ExternalImport.19.26ca20')" width="80">
              <template #default="{ row }">
                <el-tag size="small">{{ row.source }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.ExternalImport.12.3fea7c')" width="90">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="total_items" :label="$t('auto.ExternalImport.13.9ed7d3')" width="60" />
            <el-table-column :label="$t('auto.ExternalImport.20.b9ae34')" width="100">
              <template #default="{ row }">
                <span style="color: #10B981">{{ row.done_items || 0 }}</span>
                /
                <span style="color: #EF4444">{{ row.failed_items || 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('auto.ExternalImport.21.19fcb9')" width="150">
              <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, Document, Promotion } from '@element-plus/icons-vue'
import { knowledgeMerchantAPI } from '@/api/knowledgeMerchant'
import { ragProductConfigAPI } from '@/api/ragProductConfig'

const { t } = useI18n()

const importing = ref(false)
const jobsLoading = ref(false)
const productList = ref([])
const apiToken = ref('')
const jobFilterProduct = ref('')
const jobs = ref([])
const lastResponse = ref(null)

const form = reactive({
  source: 'custom',
  product_id: '',
  feishu_doc_id: '',
  notion_page_id: '',
  itemsJson: ''
})

const testMode = ref(false)

const canSubmit = computed(() => {
  if (!form.product_id) return false
  if (form.source === 'custom') return form.itemsJson.trim().length > 0
  if (form.source === 'feishu') return !!form.feishu_doc_id
  if (form.source === 'dingtalk') return !!form.dingtalk_doc_id
  if (form.source === 'notion') return !!form.notion_page_id
  return false
})

const loadProducts = async () => {
  try {
    const res = await ragProductConfigAPI.listProducts()
    if (Array.isArray(res)) productList.value = res
    else if (res?.items) productList.value = res.items
  } catch (e) {
    console.error(t('auto.ExternalImport.36.fcbbfe'), e)
  }
}

const handleSourceChange = () => {
  form.feishu_doc_id = ''
  form.dingtalk_doc_id = ''
  form.notion_page_id = ''
  form.itemsJson = ''
}

const loadTemplate = () => {
  form.itemsJson = JSON.stringify([
    { title: i18n.global.t(t('auto.ExternalImport.37.1a1063')), content: i18n.global.t(t('auto.ExternalImport.38.48b1d8')), category: t('auto.ExternalImport.39.59bd68'), tags: [t('auto.ExternalImport.40.44c198'), t('auto.ExternalImport.39.59bd68')] },
    { title: i18n.global.t(t('auto.ExternalImport.41.85bead')), content: i18n.global.t(t('auto.ExternalImport.42.8e40c9')), category: t('auto.ExternalImport.39.59bd68'), tags: [t('auto.ExternalImport.43.b909d2')] }
  ], null, 2)
}

const handleImport = async () => {
  if (!canSubmit.value) {
    ElMessage.warning(i18n.global.t(t('auto.ExternalImport.44.5b5301')))
    return
  }
  if (!apiToken.value) {
    ElMessage.warning(i18n.global.t(t('auto.ExternalImport.45.188b02')))
    return
  }
  let items
  if (form.source === 'custom') {
    try {
      items = JSON.parse(form.itemsJson)
      if (!Array.isArray(items)) throw new Error('items 必须是数组')
    } catch (e) {
      ElMessage.error('JSON 格式错误: ' + e.message)
      return
    }
  }
  importing.value = true
  try {
    const payload = {
      source: form.source,
      product_id: form.product_id,
      sync: testMode.value
    }
    if (form.source === 'feishu') {
      payload.feishu_doc_id = form.feishu_doc_id
    } else if (form.source === 'dingtalk') {
      payload.dingtalk_doc_id = form.dingtalk_doc_id
    } else if (form.source === 'notion') {
      payload.notion_page_id = form.notion_page_id
    } else {
      payload.items = items
    }
    const res = await knowledgeMerchantAPI.externalImport(payload, apiToken.value)
    if (res) {
      lastResponse.value = res
      ElMessage.success(res.async ? t('auto.ExternalImport.46.52c403') : t('auto.ExternalImport.47.8edcee'))
      loadJobs()
    }
  } catch (e) {
    ElMessage.error(t('auto.ExternalImport.48.931ce4') + (e.message || ''))
  } finally {
    importing.value = false
  }
}

const loadJobs = async () => {
  jobsLoading.value = true
  try {
    const res = await knowledgeMerchantAPI.listExternalJobs({
      product_id: jobFilterProduct.value,
      page: 1,
      page_size: 20
    })
    jobs.value = res?.items || []
  } catch (e) {
    console.error(t('auto.ExternalImport.49.7fb38b'), e)
  } finally {
    jobsLoading.value = false
  }
}

const statusType = (s) => ({ pending: 'info', running: 'warning', completed: 'success', failed: 'danger' }[s] || '')
const statusLabel = (s) => ({ pending: t('auto.ExternalImport.50.65dd9e'), running: t('auto.ExternalImport.51.5d459d'), completed: t('auto.ExternalImport.52.fad522'), failed: t('auto.ExternalImport.15.acd5cb') }[s] || s)
const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(async () => {
  await loadProducts()
  await loadJobs()
  apiToken.value = localStorage.getItem('knowledge_merchant_token') || '';
})

import { watch } from 'vue';
watch(apiToken, (val) => {
  if (val) localStorage.setItem('knowledge_merchant_token', val)
})
</script>

<style scoped lang="scss">
.external-import-page {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
