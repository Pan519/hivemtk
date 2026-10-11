<template>
  <div class="geo-page">

    
    <el-card shadow="never" class="search-card">
      <el-form :model="query" label-width="92px" class="search-form">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item :label="$t('auto.KeywordMining.1.d83af8')">
              <el-input
                v-model="query.seed_words"
                :placeholder="$t('auto.KeywordMining.2.40bb36')"
                type="textarea"
                :rows="2"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item :label="$t('auto.KeywordMining.3.037701')">
              <el-select v-model="query.mode" style="width: 100%">
                <el-option :label="$t('auto.KeywordMining.4.e6b171')" value="llm" />
                <el-option :label="$t('auto.KeywordMining.5.8446cf')" value="combination" />
                <el-option :label="$t('auto.KeywordMining.6.29508c')" value="mixed" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="4" class="action-col">
            <el-button type="primary" :loading="mining" @click="handleMine">
              <el-icon><Search /></el-icon>
              <span>{{ $t('auto.KeywordMining.27.8b88a7') }}</span>
            </el-button>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :xs="24" :md="14">
            <el-form-item :label="$t('auto.KeywordMining.7.5b5734')">
              <el-checkbox-group v-model="query.engines">
                <el-checkbox value="baidu">{{ $t('auto.KeywordMining.28.a3f4a5') }}</el-checkbox>
                <el-checkbox value="bing">Bing</el-checkbox>
                <el-checkbox value="google">Google</el-checkbox>
                <el-checkbox value="360">360</el-checkbox>
                <el-checkbox value="sogou">{{ $t('auto.KeywordMining.29.10d247') }}</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="5">
            <div class="pipeline-actions">
              <el-button :loading="suggesting" @click="handleCrawlSuggest">
                <el-icon><Download /></el-icon>
                <span>{{ $t('auto.KeywordMining.30.fe2805') }}</span>
              </el-button>
              <el-button :loading="longtailing" @click="handleLongtail">
                <el-icon><MagicStick /></el-icon>
                <span>{{ $t('auto.KeywordMining.31.3aedd0') }}</span>
              </el-button>
            </div>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    
    <el-card shadow="never" class="funnel-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.KeywordMining.32.06c839') }}</span>
          <el-button link type="primary" :loading="funnelLoading" @click="loadFunnel">{{ $t('auto.KeywordMining.33.694fc5') }}</el-button>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col :xs="24" :md="14">
          <div class="funnel-layers">
            <div v-for="layer in funnelLayers" :key="layer.key" class="funnel-row">
              <span class="funnel-label">{{ layer.label }}</span>
              <div class="funnel-bar-wrap">
                <div class="funnel-bar" :style="{ width: layer.percent + '%', background: layer.color }" />
              </div>
              <span class="funnel-count">{{ layer.count }}</span>
            </div>
          </div>
          <div class="funnel-meta">
            <span>{{ $t('auto.KeywordMining.34.599b5a') }}<b>{{ funnel.total || 0 }}</b></span>
            <span v-for="stage in funnelStageList" :key="stage.name" class="funnel-stage-tag">
              {{ stageLabel(stage.name) }}: {{ stage.count }}
            </span>
          </div>
        </el-col>
        <el-col :xs="24" :md="10">
          <div class="funnel-side">
            <div class="funnel-side-title">{{ $t('auto.KeywordMining.35.1b65f2') }}</div>
            <div v-if="engineList.length" class="tag-wrap">
              <el-tag v-for="e in engineList" :key="e.name" size="small" effect="light">
                {{ engineLabel(e.name) }} · {{ e.count }}
              </el-tag>
            </div>
            <el-empty v-else :description="$t('auto.KeywordMining.8.7609ee')" :image-size="48" />
            <div class="funnel-side-title" style="margin-top: 12px">{{ $t('auto.KeywordMining.36.0ae513') }}</div>
            <div v-if="intentList.length" class="tag-wrap">
              <el-tag v-for="i in intentList" :key="i.name" size="small" type="success" effect="plain">
                {{ intentLabel(i.name) }} · {{ i.count }}
              </el-tag>
            </div>
            <el-empty v-else :description="$t('auto.KeywordMining.9.fee39c')" :image-size="48" />
          </div>
        </el-col>
      </el-row>
    </el-card>

    
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.KeywordMining.37.f3f5fd') }}</span>
          <div class="header-actions">
            <el-select v-model="listQuery.layer" :placeholder="$t('auto.KeywordMining.10.8f3acb')" clearable style="width: 130px" @change="handleLayerChange">
              <el-option :label="$t('auto.KeywordMining.1.d83af8')" value="seed" />
              <el-option :label="$t('auto.KeywordMining.11.a85b92')" value="related" />
              <el-option :label="$t('auto.KeywordMining.12.719b80')" value="suggest" />
              <el-option :label="$t('auto.KeywordMining.13.022f20')" value="longtail" />
            </el-select>
            <el-input
              v-model="listQuery.search"
              :placeholder="$t('auto.KeywordMining.14.9e415f')"
              clearable
              style="width: 240px"
              @keyup.enter="loadList"
              @clear="loadList"
            />
            <el-button @click="loadList">{{ $t('auto.KeywordMining.38.bee912') }}</el-button>
          </div>
        </div>
      </template>

      <el-table v-loading="loading" :data="tableData" stripe style="width: 100%">
        <el-table-column type="selection" width="44" />
        <el-table-column prop="keyword" :label="$t('auto.KeywordMining.15.9699a5')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="layer" :label="$t('auto.KeywordMining.16.6896da')" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="layerTagType(row.layer)">{{ layerLabel(row.layer) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="parent_keyword" :label="$t('auto.KeywordMining.17.514925')" width="130" show-overflow-tooltip>
          <template #default="{ row }">{{ row.parent_keyword || '—' }}</template>
        </el-table-column>
        <el-table-column prop="category" :label="$t('auto.KeywordMining.18.383a6d')" width="110" show-overflow-tooltip />
        <el-table-column prop="source" :label="$t('auto.KeywordMining.19.26ca20')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="sourceTagType(row.source)">{{ row.source || '—' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="search_volume" :label="$t('auto.KeywordMining.20.71cdef')" width="90" align="right" />
        <el-table-column prop="difficulty" :label="$t('auto.KeywordMining.21.3f256e')" width="80" align="center">
          <template #default="{ row }">
            <el-progress
              v-if="row.difficulty != null"
              :percentage="Number(row.difficulty) || 0"
              :stroke-width="6"
              :show-text="false"
            />
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column prop="query_intent" :label="$t('auto.KeywordMining.22.764aca')" width="90">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ intentLabel(row.query_intent || row.intent) || '—' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="suggest_count" :label="$t('auto.KeywordMining.23.96247f')" width="90" align="center">
          <template #default="{ row }">
            <span v-if="row.suggest_count > 0">{{ row.suggest_count }} 个引擎</span>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.KeywordMining.24.3fea7c')" width="80">
          <template #default="{ row }">
            <el-tag size="small" ::type="$t('auto.KeywordMining.25.550edc')">
              {{ row.status || '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.KeywordMining.26.2b6bc0')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" :loading="row._expanding" @click="handleExpand(row)">{{ $t('auto.KeywordMining.39.7cff93') }}</el-button>
            <el-button link type="primary" :loading="row._clustering" @click="handleCluster(row)">{{ $t('auto.KeywordMining.40.b39799') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ $t('auto.KeywordMining.41.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pager">
        <el-pagination
          v-model:current-page="listQuery.page"
          v-model:page-size="listQuery.limit"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Download, MagicStick } from '@element-plus/icons-vue'
import { geoApi } from '@/api/geo'

const { t } = useI18n()

const brandName = ref('')

const mining = ref(false)
const suggesting = ref(false)
const longtailing = ref(false)
const loading = ref(false)
const funnelLoading = ref(false)
const tableData = ref([])
const total = ref(0)
const funnel = ref({})

const query = reactive({
  seed_words: '',
  mode: 'llm',
  engines: ['baidu', 'bing', 'google']
})

const listQuery = reactive({
  search: '',
  layer: '',
  page: 1,
  limit: 20
})

const LAYER_META = {
  seed: { label: t('auto.KeywordMining.1.d83af8'), color: '#409eff' },
  related: { label: t('auto.KeywordMining.11.a85b92'), color: '#67c23a' },
  suggest: { label: t('auto.KeywordMining.12.719b80'), color: '#e6a23c' },
  longtail: { label: t('auto.KeywordMining.13.022f20'), color: '#f56c6c' }
}

const funnelLayers = computed(() => {
  const max = Math.max(
    funnel.value.seed_count || 0,
    funnel.value.related_count || 0,
    funnel.value.suggest_count || 0,
    funnel.value.longtail_count || 0,
    1
  )
  return [
    { key: 'seed', label: t('auto.KeywordMining.1.d83af8'), count: funnel.value.seed_count || 0 },
    { key: 'related', label: t('auto.KeywordMining.11.a85b92'), count: funnel.value.related_count || 0 },
    { key: 'suggest', label: t('auto.KeywordMining.12.719b80'), count: funnel.value.suggest_count || 0 },
    { key: 'longtail', label: t('auto.KeywordMining.13.022f20'), count: funnel.value.longtail_count || 0 }
  ].map(l => ({
    ...l,
    color: LAYER_META[l.key].color,
    percent: Math.max(Math.round((l.count / max) * 100), l.count > 0 ? 4 : 0)
  }))
})

const engineList = computed(() =>
  Object.entries(funnel.value.suggest_engines || {})
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
)

const intentList = computed(() =>
  Object.entries(funnel.value.intents || {})
    .filter(([name]) => name)
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
)

const funnelStageList = computed(() =>
  Object.entries(funnel.value.funnel_stages || {})
    .filter(([name]) => name)
    .map(([name, count]) => ({ name, count }))
)

const ENGINE_LABELS = { baidu: t('auto.KeywordMining.28.a3f4a5'), bing: 'Bing', google: 'Google', 360: '360', sogou: t('auto.KeywordMining.29.10d247') }
const INTENT_LABELS = {
  informational: t('auto.KeywordMining.42.d8c7e0'), navigational: t('auto.KeywordMining.43.056f2d'), transactional: t('auto.KeywordMining.44.2685c0'),
  commercial: t('auto.KeywordMining.45.837171'), 信息: t('auto.KeywordMining.42.d8c7e0'), 导航: t('auto.KeywordMining.43.056f2d'), 交易: t('auto.KeywordMining.44.2685c0'), 商业: t('auto.KeywordMining.45.837171')
}
const STAGE_LABELS = { awareness: t('auto.KeywordMining.46.2351d7'), consideration: t('auto.KeywordMining.47.42d166'), decision: t('auto.KeywordMining.48.5fd3a1') }

const engineLabel = (name) => ENGINE_LABELS[name] || name
const intentLabel = (name) => (name ? (INTENT_LABELS[name] || name) : '')
const stageLabel = (name) => STAGE_LABELS[name] || name
const layerLabel = (layer) => LAYER_META[layer || 'seed']?.label || layer || t('auto.KeywordMining.1.d83af8')

const layerTagType = (layer) => {
  const map = { seed: 'primary', related: 'success', suggest: 'warning', longtail: 'danger' }
  return map[layer] || 'primary'
}

const sourceTagType = (source) => {
  const map = { AI: 'primary', 词库: 'warning', 混合: 'success', suggest: 'warning', longtail: 'danger', 下拉词: 'warning', 长尾: 'danger' }
  return map[source] || 'info'
}

const parseSeeds = () => {
  const seeds = query.seed_words.split(/[\n\r，,]/).map(s => s.trim()).filter(Boolean)
  if (!seeds.length) {
    ElMessage.warning(t('auto.KeywordMining.49.a922b6'))
    return null
  }
  return seeds
}

const loadList = async () => {
  loading.value = true
  try {
    const params = { ...listQuery }
    if (!params.layer) delete params.layer
    const res = await geoApi.getKeywordList(params)
    tableData.value = res?.list || res?.items || res || []
    total.value = res?.total || tableData.value.length
  } catch (e) {
    ElMessage.error(e.message || t('auto.KeywordMining.50.e156e9'))
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const loadFunnel = async () => {
  funnelLoading.value = true
  try {
    funnel.value = (await geoApi.getKeywordFunnel()) || {}
  } catch (e) {
    ElMessage.error(e.message || t('auto.KeywordMining.51.cbf6c6'))
  } finally {
    funnelLoading.value = false
  }
}

const handleLayerChange = () => {
  listQuery.page = 1
  loadList()
}

const handleMine = async () => {
  if (!query.seed_words.trim()) {
    ElMessage.warning(t('auto.KeywordMining.49.a922b6'))
    return
  }
  mining.value = true
  try {
    const seedWords = query.seed_words.split(/[\n\r，,]/).map(s => s.trim()).filter(Boolean)
    const res = await geoApi.mineKeywords({
      seed_words: seedWords,
      mode: query.mode
    })
    ElMessage.success(`已挖掘 ${Array.isArray(res) ? res.length : (res?.total || 0)} 个关键词`)
    listQuery.page = 1
    await Promise.all([loadList(), loadFunnel()])
  } catch (e) {
    ElMessage.error(e.message || t('auto.KeywordMining.52.2b9d4f'))
  } finally {
    mining.value = false
  }
}

const handleCrawlSuggest = async () => {
  const seeds = parseSeeds()
  if (!seeds) return
  if (!query.engines.length) {
    ElMessage.warning(t('auto.KeywordMining.53.a77c26'))
    return
  }
  suggesting.value = true
  try {
    const res = await geoApi.crawlSuggest({ seeds, engines: query.engines })
    ElMessage.success(`下拉词抓取 ${res?.count || 0} 个，落库 ${res?.saved || 0} 个`)
    listQuery.page = 1
    await Promise.all([loadList(), loadFunnel()])
  } catch (e) {
    ElMessage.error(e.message || t('auto.KeywordMining.54.91b052'))
  } finally {
    suggesting.value = false
  }
}

const handleLongtail = async () => {
  const seeds = parseSeeds()
  if (!seeds) return
  longtailing.value = true
  try {
    const res = await geoApi.combineLongtail({ seeds, use_default_templates: true })
    ElMessage.success(`长尾词组合 ${res?.count || 0} 个，落库 ${res?.saved || 0} 个`)
    listQuery.page = 1
    await Promise.all([loadList(), loadFunnel()])
  } catch (e) {
    ElMessage.error(e.message || t('auto.KeywordMining.55.dd3912'))
  } finally {
    longtailing.value = false
  }
}

const handleExpand = async (row) => {
  row._expanding = true
  try {
    const res = await geoApi.semanticExpand({ keywords: [row.keyword], brand_name: brandName.value })
    const expanded = Array.isArray(res) ? res : (res?.keywords || [])
    ElMessage.success(`语义扩展完成，新增 ${expanded.length} 个近义 / 长尾词`)
    await Promise.all([loadList(), loadFunnel()])
  } catch (e) {
    ElMessage.error(e.message || t('auto.KeywordMining.56.3c31a0'))
  } finally {
    row._expanding = false
  }
}

const handleCluster = async (row) => {
  row._clustering = true
  try {
    await geoApi.topicCluster({ keywords: [row.keyword], brand_name: brandName.value })
    ElMessage.success(t('auto.KeywordMining.57.f814d8'))
    await loadList()
  } catch (e) {
    ElMessage.error(e.message || t('auto.KeywordMining.58.1ee42b'))
  } finally {
    row._clustering = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除关键词「${row.keyword}」？`, t('auto.KeywordMining.59.50eaf9'), {
    type: 'warning'
  })
    .then(async () => {
      try {
        await geoApi.deleteKeyword(row.id)
        ElMessage.success(t('auto.KeywordMining.60.5cc232'))
        await Promise.all([loadList(), loadFunnel()])
      } catch (e) {
        ElMessage.error(e.message || t('auto.KeywordMining.61.acf066'))
      }
    })
    .catch(() => {})
}

onMounted(async () => {
  loadList()
  loadFunnel()
  try {
    const cfg = await geoApi.getConfig()
    brandName.value = cfg?.brand_name || ''
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
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
.search-card,
.funnel-card,
.table-card {
  border: 1px solid $border-base;
  border-radius: 10px;
  margin-bottom: $spacing-md;
}
.action-col {
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  padding-bottom: $spacing-xs;
}
.pipeline-actions {
  display: flex;
  gap: $spacing-sm;
  padding-bottom: $spacing-xs;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: $spacing-md;
}
.header-actions {
  display: flex;
  gap: $spacing-sm;
}
.funnel-layers {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}
.funnel-row {
  display: flex;
  align-items: center;
  gap: $spacing-md;
}
.funnel-label {
  width: 56px;
  flex-shrink: 0;
  font-size: $font-size-small;
  color: $text-regular;
  text-align: right;
}
.funnel-bar-wrap {
  flex: 1;
  height: 18px;
  background: rgba(128, 128, 128, 0.08);
  border-radius: 9px;
  overflow: hidden;
}
.funnel-bar {
  height: 100%;
  border-radius: 9px;
  transition: width 0.4s ease;
  min-width: 0;
}
.funnel-count {
  width: 52px;
  flex-shrink: 0;
  font-weight: 600;
  color: $text-primary;
}
.funnel-meta {
  margin-top: $spacing-md;
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-md;
  font-size: $font-size-small;
  color: $info-color;
}
.funnel-stage-tag {
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(128, 128, 128, 0.08);
}
.funnel-side {
  padding-left: $spacing-md;
}
.funnel-side-title {
  font-size: $font-size-small;
  color: $info-color;
  margin-bottom: $spacing-xs;
}
.tag-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-xs;
}
.pager {
  margin-top: $spacing-md;
  display: flex;
  justify-content: flex-end;
}
</style>
