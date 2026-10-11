<template>
  <div class="geo-page">

    
    <el-card shadow="never" class="form-card">
      <template #header><span class="card-title">{{ $t('auto.Verification.21.4fd78a') }}</span></template>
      <el-form :model="form" label-width="92px">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="$t('auto.Verification.1.1e9faa')">
              <el-select
                v-model="form.article_id"
                :placeholder="$t('auto.Verification.2.1b3ab1')"
                filterable
                clearable
                style="width: 100%"
                @change="onArticleChange"
              >
                <el-option
                  v-for="a in articles"
                  :key="a.id"
                  :label="a.title || a.keyword"
                  :value="a.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="$t('auto.Verification.3.4e8713')">
              <el-input v-model="form.brand_name" :placeholder="$t('auto.Verification.4.a08eb5')" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="$t('auto.Verification.5.3a17e1')">
              <el-select
                v-model="form.models"
                multiple
                collapse-tags
                collapse-tags-tooltip
                :placeholder="$t('auto.Verification.6.9f701d')"
                style="width: 100%"
              >
                <el-option
                  v-for="m in modelOptions"
                  :key="m.value"
                  :label="m.label"
                  :value="m.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('auto.Verification.7.dd6023')">
          <el-input
            v-model="form.queries"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Verification.8.b5b294')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="verifying" @click="handleVerify">
            <el-icon><CircleCheck /></el-icon><span>{{ $t('auto.Verification.22.cd8992') }}</span>
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card v-if="results.length" shadow="never" class="result-card">
      <template #header><span class="card-title">{{ $t('auto.Verification.23.fc3c71') }}</span></template>
      <el-table v-loading="verifying" :data="results" stripe style="width: 100%">
        <el-table-column prop="model" :label="$t('auto.Verification.5.3a17e1')" width="140" />
        <el-table-column prop="query" :label="$t('auto.Verification.9.5dc99f')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="brand_name" :label="$t('auto.Verification.10.09307c')" width="120" />
        <el-table-column prop="brand_mentioned" :label="$t('auto.Verification.11.0a01f6')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.brand_mentioned ? 'success' : 'info'" size="small">
              {{ row.brand_mentioned ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="mention_count" :label="$t('auto.Verification.12.2ee012')" width="100" align="right" />
        <el-table-column prop="sentiment" :label="$t('auto.Verification.13.2b1c16')" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="sentimentType(row.sentiment)">{{ row.sentiment || '—' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="position" :label="$t('auto.Verification.14.1b2b4b')" width="130" />
      </el-table>
    </el-card>

    
    <el-card shadow="never" class="negative-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ $t('auto.Verification.24.53119f') }}</span>
          <el-tooltip ::content="$t('auto.Verification.15.f8761e')" :disabled="!!form.brand_name.trim()" placement="top">
            <el-button type="warning" :loading="monitoring" :disabled="!form.brand_name.trim()" @click="handleMonitorNegative">
              <el-icon><Warning /></el-icon><span>{{ $t('auto.Verification.25.2e00d1') }}</span>
            </el-button>
          </el-tooltip>
        </div>
      </template>
      <el-empty v-if="!negativeResults.length" :description="$t('auto.Verification.16.74e452')" :image-size="60" />
      <el-table v-else :data="negativeResults" stripe style="width: 100%">
        <el-table-column prop="query" :label="$t('auto.Verification.17.eb8454')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="mention_count" :label="$t('auto.Verification.12.2ee012')" width="100" align="right" />
        <el-table-column prop="negative_score" :label="$t('auto.Verification.18.5437e7')" width="100" align="right" />
        <el-table-column prop="risk_level" :label="$t('auto.Verification.19.25eb6b')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="riskType(row.risk_level)">{{ row.risk_level || '低' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="risk_description" :label="$t('auto.Verification.20.76f1c8')" min-width="220" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { CircleCheck, Warning } from '@element-plus/icons-vue'
import { geoApi } from '@/api/geo'
import { http } from '@/utils/request'

const { t } = useI18n()

const modelOptions = ref([])

const loadModelOptions = async () => {
  try {
    const res = await http.get('/api/llm/models')
    const list = Array.isArray(res) ? res : (res?.list || res?.items || res?.data || [])
    const options = list
      .filter((m) => {
        if (!m || m.enabled === false) return false
        if (m.vendor === 'local')
          return false;
        const url = (m.base_url || m.baseURL || '').toLowerCase()
        if (url.includes('127.0.0.1') || url.includes('localhost') || url.includes('0.0.0.0')) return false
        return true
      })
      .map((m) => ({
        label: `${m.display_name || m.name} (${m.model || m.name})`,
        value: m.name
      }))
    modelOptions.value = options
    if (options.length && (!form.models || !form.models.length)) {
      form.models = [options[0].value]
    }
  } catch (e) {
    console.warn(t('auto.Verification.26.a8d452'), e.message)
  }
}

const form = reactive({
  article_id: '',
  brand_name: '',
  models: [],
  queries: ''
})

const articles = ref([])
const verifying = ref(false)
const monitoring = ref(false)
const results = ref([])
const negativeResults = ref([])

const sentimentType = (s) => {
  const map = { 正面: 'success', positive: 'success', 中性: 'info', neutral: 'info', 负面: 'danger', negative: 'danger' }
  return map[s] || 'info'
}
const riskType = (r) => {
  const map = { 高: 'danger', 中: 'warning', 低: 'success' }
  return map[r] || 'success'
}

const loadArticles = async () => {
  try {
    const res = await geoApi.getArticleList({ page: 1, limit: 100 })
    articles.value = res?.list || res?.items || res || []
  } catch (e) {
    articles.value = []
  }
}

const onArticleChange = async (id) => {
  if (!id) return
  try {
    await geoApi.getArticleByID(id)
  } catch (e) {
    ElMessage.error(e.message || t('auto.Verification.27.aaff45'))
  }
}

const handleVerify = async () => {
  if (!form.queries.trim()) {
    ElMessage.warning(t('auto.Verification.28.2b216e'))
    return
  }
  if (!form.brand_name.trim()) {
    ElMessage.warning(t('auto.Verification.29.0e7774'))
    return
  }
  verifying.value = true
  results.value = []
  try {
    const queries = form.queries.split('\n').map((s) => s.trim()).filter(Boolean);
    for (const q of queries) {
      const res = await geoApi.verifyArticle({
        article_id: form.article_id || undefined,
        query: q,
        brand_name: form.brand_name,
        models: form.models
      })
      if (res) results.value.push(res)
    }
    ElMessage.success(t('auto.Verification.30.816ade'))
  } catch (e) {
    ElMessage.error(e.message || t('auto.Verification.31.e441b1'))
  } finally {
    verifying.value = false
  }
}

const handleMonitorNegative = async () => {
  if (!form.brand_name.trim()) {
    ElMessage.warning(t('auto.Verification.32.8d1157'))
    return
  }
  monitoring.value = true
  try {
    const res = await geoApi.monitorNegative({
      brand_name: form.brand_name
    })
    negativeResults.value = res?.queries || []
    ElMessage.success(t('auto.Verification.33.cf50ff'))
  } catch (e) {
    ElMessage.error(e.message || t('auto.Verification.34.e3cb70'))
  } finally {
    monitoring.value = false
  }
}

onMounted(async () => {
  await loadModelOptions()
  await loadArticles()
})
</script>

<style lang="scss" scoped>
.geo-page {
  padding: $spacing-lg 24px;
}
.page-header h2 {
  margin: 0 0 6px;
  font-size: $font-size-extra-large;
  font-weight: 700;
  color: $text-primary;
}
.page-header .sub {
  margin: 0 0 16px;
  color: $info-color;
  font-size: $font-size-small;
}
.form-card,
.result-card,
.negative-card {
  border: 1px solid $border-base;
  border-radius: 10px;
  margin-bottom: $spacing-md;
}
.card-title {
  font-weight: 600;
  color: $text-primary;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
