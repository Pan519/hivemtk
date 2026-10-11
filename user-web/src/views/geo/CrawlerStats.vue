<template>
  <div class="geo-page p-4">
    
    <el-row :gutter="12" class="mb-4">
      <el-col :span="6">
        <el-card shadow="hover" class="coverage-card">
          <div class="text-center py-2">
            <div class="text-xs text-gray-500 mb-1">{{ $t('auto.CrawlerStats.17.94079c') }}</div>
            <div class="text-3xl font-bold" :style="{ color: coverageColor }">{{ (coverageScore || 0).toFixed(1) }}%</div>
            <el-progress :percentage="Math.round(coverageScore || 0)" :color="coverageColor" :stroke-width="8" class="mt-2" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="3">
        <el-card shadow="hover"><el-statistic :title="$t('auto.CrawlerStats.1.c7810f')" :value="s.today_visits || 0" /></el-card>
      </el-col>
      <el-col :span="3">
        <el-card shadow="hover"><el-statistic title="HiveMTK" :value="s.hivemtk_visits || 0" /></el-card>
      </el-col>
      <el-col :span="3">
        <el-card shadow="hover"><el-statistic :title="$t('auto.CrawlerStats.2.f0d182')" :value="s.competitor_visits || 0" /></el-card>
      </el-col>
      <el-col :span="3">
        <el-card shadow="hover"><el-statistic :title="$t('auto.CrawlerStats.3.e9870d')" :value="s.active_keywords || 0" /></el-card>
      </el-col>
      <el-col :span="3">
        <el-card shadow="hover"><el-statistic :title="$t('auto.CrawlerStats.4.02b40a')" :value="s.active_engines || 0" /></el-card>
      </el-col>
      <el-col :span="3">
        <el-card shadow="hover"><el-statistic :title="$t('auto.CrawlerStats.5.7ff807')" :value="s.active_domains || 0" /></el-card>
      </el-col>
    </el-row>

    
    <el-card class="mb-4">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">{{ $t('auto.CrawlerStats.18.130885') }}</span>
          <div class="flex items-center gap-2">
            <el-button size="small" type="primary" @click="loadAll">{{ $t('auto.CrawlerStats.19.694fc5') }}</el-button>
            <el-button size="small" type="success" :loading="running" @click="triggerCrawl">{{ $t('auto.CrawlerStats.20.4a152a') }}</el-button>
          </div>
        </div>
      </template>
      <el-table :data="domainCompare" v-loading="loading" size="default" :default-sort="{prop:'visits', order:'descending'}">
        <el-table-column :label="$t('auto.CrawlerStats.6.a4dc00')" width="60">
          <template #default="{ $index }">
            <span class="font-bold">{{ $index + 1 }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.CrawlerStats.7.190980')" min-width="240">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <el-tag v-if="row.is_hivemtk" type="success" size="small" effect="dark" round>🔥 HiveMTK</el-tag>
              <el-tag v-else type="info" size="small" effect="plain" round>{{ $t('auto.CrawlerStats.21.9868d2') }}</el-tag>
              <span class="font-mono" :class="row.is_hivemtk ? 'font-bold' : ''">{{ row.domain }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="visits" :label="$t('auto.CrawlerStats.8.9ca20a')" width="120" sortable>
          <template #default="{ row }">
            <span class="font-mono font-bold text-lg">{{ row.visits }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="engines" :label="$t('auto.CrawlerStats.9.18d0f4')" width="120">
          <template #default="{ row }">
            <el-progress :percentage="Math.round(row.engines / 8 * 100)" :stroke-width="10" :color="row.engines >= 6 ? '#22c55e' : row.engines >= 4 ? '#eab308' : '#ef4444'" />
            <div class="text-xs text-gray-500 mt-1">{{ row.engines }} / 8</div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.CrawlerStats.10.a790ba')" width="200">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <el-progress :percentage="Math.round(row.share_pct)" :stroke-width="12" :color="row.is_hivemtk ? '#10b981' : '#6366f1'" style="flex:1" />
              <span class="font-mono font-bold w-14 text-right">{{ row.share_pct.toFixed(1) }}%</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="source_level" :label="$t('auto.CrawlerStats.11.b99f2a')" width="100">
          <template #default="{ row }">
            <el-tag :type="levelType(row.source_level)" size="small">{{ row.source_level || 'D' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <template #empty>
        <div class="py-8 text-gray-400">{{ $t('auto.CrawlerStats.22.569e80') }}</div>
      </template>
    </el-card>

    
    <el-card class="mb-4">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">{{ $t('auto.CrawlerStats.23.057067') }}</span>
          <el-input v-model="kwFilter" :placeholder="$t('auto.CrawlerStats.12.2f5112')" size="small" style="width:160px" clearable />
        </div>
      </template>
      <el-table :data="filteredKeywordStats" v-loading="loading" size="small" :default-sort="{prop:'visit_count', order:'descending'}">
        <el-table-column prop="keyword" :label="$t('auto.CrawlerStats.13.9699a5')" min-width="220">
          <template #default="{ row }"><span class="font-semibold">{{ row.keyword }}</span></template>
        </el-table-column>
        <el-table-column prop="engine" :label="$t('auto.CrawlerStats.4.02b40a')" width="160">
          <template #default="{ row }">
            <el-tag size="small" :type="engineTagType(row.engine)">{{ row.engine }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="visit_count" :label="$t('auto.CrawlerStats.14.f0be64')" width="100" sortable>
          <template #default="{ row }"><span class="font-mono font-bold">{{ row.visit_count }}</span></template>
        </el-table-column>
        <el-table-column :label="$t('auto.CrawlerStats.15.2e2566')" min-width="240">
          <template #default="{ row }">
            <div class="heat-bar" :style="{ width: heatWidth(row.visit_count) + '%', background: heatColor(row.visit_count) }"></div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">{{ $t('auto.CrawlerStats.24.00fa60') }}</span>
          <el-tag size="small" type="info">{{ $t('auto.CrawlerStats.25.95d706') }}</el-tag>
        </div>
      </template>
      <el-table :data="domainStats" v-loading="loading" size="small" :default-sort="{prop:'visit_count', order:'descending'}">
        <el-table-column prop="domain" :label="$t('auto.CrawlerStats.7.190980')" min-width="200">
          <template #default="{ row }"><span class="font-mono">{{ row.domain }}</span></template>
        </el-table-column>
        <el-table-column prop="engine" :label="$t('auto.CrawlerStats.16.6c1b19')" width="160" />
        <el-table-column prop="visit_count" :label="$t('auto.CrawlerStats.14.f0be64')" width="100" sortable />
        <el-table-column prop="source_level" :label="$t('auto.CrawlerStats.11.b99f2a')" width="100">
          <template #default="{ row }">
            <el-tag :type="levelType(row.source_level)" size="small">{{ row.source_level || 'D' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getCrawlerStats, runCrawler as runCrawlerApi } from '@/api/geoProbe.js'

const { t } = useI18n()

const loading = ref(false)
const running = ref(false)
const keywordStats = ref([])
const domainStats = ref([])
const domainCompare = ref([])
const coverageScore = ref(0)
const s = ref({})
const kwFilter = ref('')

const filteredKeywordStats = computed(() => {
  if (!kwFilter.value) return keywordStats.value
  return keywordStats.value.filter(k => k.keyword.includes(kwFilter.value))
})

const coverageColor = computed(() => {
  const v = coverageScore.value
  if (v >= 70) return '#10b981'
  if (v >= 50) return '#f59e0b'
  return '#ef4444'
})

const maxCount = computed(() => {
  const all = [...keywordStats.value, ...domainStats.value].map(d => d.visit_count || 0)
  return Math.max(1, ...all)
})
const heatWidth = (c) => Math.min(100, Math.round(((c || 0) / maxCount.value) * 100))
const heatColor = (c) => {
  const ratio = (c || 0) / maxCount.value
  if (ratio > 0.8) return '#ef4444'
  if (ratio > 0.6) return '#f59e0b'
  if (ratio > 0.3) return '#eab308'
  return '#22c55e'
}
const levelType = (l) => ({ A: 'success', B: 'primary', C: 'warning', D: 'info' }[l] || 'info')
const engineTagType = (e) => {
  if (e.includes('GPT')) return 'primary'
  if (e.includes('Claude')) return 'warning'
  if (e.includes('Perplexity')) return 'success'
  if (e.includes('Google')) return 'danger'
  return 'info'
}

const loadAll = async () => {
  loading.value = true
  try {
    const data = await getCrawlerStats()
    s.value = data?.summary || {}
    keywordStats.value = data?.keyword_stats || []
    domainStats.value = data?.domain_stats || []
    domainCompare.value = data?.domain_compare || []
    coverageScore.value = data?.coverage_score || 0
  } catch {
    keywordStats.value = []
    domainStats.value = []
    domainCompare.value = []
    coverageScore.value = 0
  }
  loading.value = false
}

const triggerCrawl = async () => {
  running.value = true
  try {
    const d = await runCrawlerApi()
    ElMessage.success(d?.message || t('auto.CrawlerStats.26.001bc6'))
    setTimeout(loadAll, 35000)
  } catch (e) {
    ElMessage.error(t('auto.CrawlerStats.27.acfa83') + (e?.data?.message || e?.message || e))
  }
  running.value = false
}

onMounted(loadAll)
</script>

<style lang="scss" scoped>
.coverage-card {
  background: linear-gradient(135deg, #f0fdf4 0%, #ecfdf5 100%);
  border-color: #86efac;
}
.heat-bar {
  height: 16px;
  border-radius: 4px;
  transition: width .4s ease;
}
</style>
