<template>
  <div class="geo-page">

    
    <el-card shadow="never" class="form-card">
      <template #header><span class="card-title">{{ $t('auto.Verification.4fd78a67') }}</span></template>
      <el-form :model="form" label-width="92px">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="$t('auto.Verification.1e9faaf9')">
              <el-select
                v-model="form.article_id"
                :placeholder="$t('auto.Verification.1b3ab189')"
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
            <el-form-item :label="$t('auto.Verification.4e871349')">
              <el-input v-model="form.brand_name" :placeholder="$t('auto.Verification.a08eb519')" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="$t('auto.Verification.3a17e15c')">
              <el-select
                v-model="form.models"
                multiple
                collapse-tags
                collapse-tags-tooltip
                :placeholder="$t('auto.Verification.9f701de3')"
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
        <el-form-item :label="$t('auto.Verification.dd60233d')">
          <el-input
            v-model="form.queries"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Verification.b5b294a6')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="verifying" @click="handleVerify">
            <el-icon><CircleCheck /></el-icon><span>{{ $t('auto.Verification.cd8992b6') }}</span>
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card v-if="results.length" shadow="never" class="result-card">
      <template #header><span class="card-title">{{ $t('auto.Verification.fc3c71a6') }}</span></template>
      <el-table v-loading="verifying" :data="results" stripe style="width: 100%">
        <el-table-column prop="model" :label="$t('auto.Verification.3a17e15c')" width="140" />
        <el-table-column prop="query" :label="$t('auto.Verification.5dc99f6e')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="brand_name" :label="$t('auto.Verification.09307ce8')" width="120" />
        <el-table-column prop="brand_mentioned" :label="$t('auto.Verification.0a01f66d')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.brand_mentioned ? 'success' : 'info'" size="small">
              {{ row.brand_mentioned ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="mention_count" :label="$t('auto.Verification.2ee0125f')" width="100" align="right" />
        <el-table-column prop="sentiment" :label="$t('auto.Verification.2b1c1688')" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="sentimentType(row.sentiment)">{{ row.sentiment || '—' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="position" :label="$t('auto.Verification.1b2b4b84')" width="130" />
      </el-table>
    </el-card>

    
    <el-card shadow="never" class="negative-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">{{ $t('auto.Verification.53119fb3') }}</span>
          <el-tooltip :content="form.brand_name.trim() ? '' : '请先填写品牌名称'" :disabled="!!form.brand_name.trim()" placement="top">
            <el-button type="warning" :loading="monitoring" :disabled="!form.brand_name.trim()" @click="handleMonitorNegative">
              <el-icon><Warning /></el-icon><span>{{ $t('auto.Verification.2e00d171') }}</span>
            </el-button>
          </el-tooltip>
        </div>
      </template>
      <el-empty v-if="!negativeResults.length" :description="$t('auto.Verification.74e45207')" :image-size="60" />
      <el-table v-else :data="negativeResults" stripe style="width: 100%">
        <el-table-column prop="query" :label="$t('auto.Verification.eb8454e9')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="mention_count" :label="$t('auto.Verification.2ee0125f')" width="100" align="right" />
        <el-table-column prop="negative_score" :label="$t('auto.Verification.5437e71b')" width="100" align="right" />
        <el-table-column prop="risk_level" :label="$t('auto.Verification.25eb6ba3')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="riskType(row.risk_level)">{{ row.risk_level || '低' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="risk_description" :label="$t('auto.Verification.76f1c846')" min-width="220" show-overflow-tooltip />
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
    console.warn(t('auto.Verification.a8d4525a'), e.message)
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
    ElMessage.error(e.message || t('auto.Verification.aaff4555'))
  }
}

const handleVerify = async () => {
  if (!form.queries.trim()) {
    ElMessage.warning(t('auto.Verification.2b216e37'))
    return
  }
  if (!form.brand_name.trim()) {
    ElMessage.warning(t('auto.Verification.0e777497'))
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
    ElMessage.success(t('auto.Verification.816ade08'))
  } catch (e) {
    ElMessage.error(e.message || t('auto.Verification.e441b11e'))
  } finally {
    verifying.value = false
  }
}

const handleMonitorNegative = async () => {
  if (!form.brand_name.trim()) {
    ElMessage.warning(t('auto.Verification.8d1157f9'))
    return
  }
  monitoring.value = true
  try {
    const res = await geoApi.monitorNegative({
      brand_name: form.brand_name
    })
    negativeResults.value = res?.queries || []
    ElMessage.success(t('auto.Verification.cf50ff00'))
  } catch (e) {
    ElMessage.error(e.message || t('auto.Verification.e3cb70be'))
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
