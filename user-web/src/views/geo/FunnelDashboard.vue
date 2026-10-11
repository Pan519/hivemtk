<template>
  <div class="funnel-dashboard">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.FunnelDashboard.3073eda2') }}</span>
          <el-button type="primary" size="small" @click="load">{{ $t('auto.FunnelDashboard.694fc5ef') }}</el-button>
        </div>
      </template>

      <div v-if="loading">{{ $t('auto.FunnelDashboard.26b5bd49') }}</div>
      <el-empty v-else-if="!stats" :description="$t('auto.FunnelDashboard.0c151559')" />

      <template v-else>
        <!-- 双漏斗 -->
        <el-row :gutter="16" class="funnel-row">
          <el-col v-for="(val, key) in funnel" :key="key" :span="4">
            <div class="funnel-step" :style="{ background: gradientColor(key) }">
              <div class="step-num">{{ val || 0 }}</div>
              <div class="step-label">{{ labelMap[key] }}</div>
            </div>
          </el-col>
        </el-row>

        <!-- 引擎对比 -->
        <el-row :gutter="16" class="stats-row">
          <el-col :span="12">
            <h4>{{ $t('auto.FunnelDashboard.c2fbb488') }}</h4>
            <el-table :data="engineIndexed" size="small" border>
              <el-table-column prop="engine" :label="$t('auto.FunnelDashboard.6c1b19dd')" width="120" />
              <el-table-column prop="count" :label="$t('auto.FunnelDashboard.ed138cfc')" />
              <el-table-column :label="$t('auto.FunnelDashboard.b0e0fc44')">
                <template #default="{ row }">
                  <el-progress :percentage="rate(row.count)" />
                </template>
              </el-table-column>
            </el-table>
          </el-col>
          <el-col :span="12">
            <h4>{{ $t('auto.FunnelDashboard.8b6f7af2') }}</h4>
            <el-table :data="engineCited" size="small" border>
              <el-table-column prop="engine" :label="$t('auto.FunnelDashboard.6c1b19dd')" width="120" />
              <el-table-column prop="count" :label="$t('auto.FunnelDashboard.373f31aa')" />
              <el-table-column :label="$t('auto.FunnelDashboard.db472b2e')">
                <template #default="{ row }">
                  <el-progress :percentage="rate(row.count)" />
                </template>
              </el-table-column>
            </el-table>
          </el-col>
        </el-row>

        <!-- 转化率 -->
        <el-row :gutter="16" class="rate-row">
          <el-col v-for="(val, key) in stats.conversion_rates" :key="key" :span="8">
            <el-statistic :title="labelMap[key]" :value="val" :precision="1" suffix="%" />
          </el-col>
        </el-row>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { geoApi } from '@/api/geo'

const { t } = useI18n()

const stats = ref(null)
const loading = ref(false)

const labelMap = {
  seed_count: t('auto.FunnelDashboard.d83af87b'),
  longtail_count: t('auto.FunnelDashboard.022f200b'),
  article_count: t('auto.FunnelDashboard.c75625dc'),
  deployed_count: t('auto.FunnelDashboard.d87522c3'),
  indexed: t('auto.FunnelDashboard.b585faf9'),
  ai_cited: t('auto.FunnelDashboard.d5c968f7'),
  longtail_to_article: t('auto.FunnelDashboard.dc1a7bdf'),
  article_to_deployed: t('auto.FunnelDashboard.33ecec4b')
}

const funnel = computed(() => {
  if (!stats.value) return {}
  return {
    seed_count: stats.value.seed_count + stats.value.longtail_count,
    article_count: stats.value.article_count,
    deployed_count: stats.value.deployed_count,
  }
})

const engineIndexed = computed(() => {
  const m = stats.value?.indexed_by_engine || {}
  return Object.entries(m).map(([engine, count]) => ({ engine, count }))
})
const engineCited = computed(() => {
  const m = stats.value?.ai_cited_by_engine || {}
  return Object.entries(m).map(([engine, count]) => ({ engine, count }))
})

function gradientColor(key) {
  const colors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399']
  const idx = { seed_count: 0, article_count: 1, deployed_count: 2 }[key] ?? 3
  return colors[idx]
}
function rate(n) {
  const base = stats.value?.deployed_count || 1
  return Math.round(n / base * 100)
}

async function load() {
  loading.value = true
  try {
    stats.value = await geoApi.getIndexFunnel()
  } catch (e) {
    console.warn(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.funnel-row { margin-bottom: 20px; }
.stats-row { margin-bottom: 20px; }
.rate-row { text-align: center; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.funnel-step {
  padding: 20px; border-radius: 8px; color: #fff; text-align: center;
  margin: 8px 4px;
}
.step-num { font-size: 28px; font-weight: bold; }
.step-label { font-size: 13px; opacity: 0.9; margin-top: 4px; }
</style>
