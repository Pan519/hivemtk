<template>
  <div class="humanize-panel">
    <el-card class="header-card">
      <div class="header-content">
        <h2>{{ $t('auto.Panel.23.c7fd80') }}</h2>
        <p class="subtitle">{{ $t('auto.Panel.24.6db90a') }}</p>
      </div>
    </el-card>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('auto.Panel.25.83a7b7') }}</div>
            <div class="stat-value">{{ stats.totalScored }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('平均拟人度') }}</div>
            <div class="stat-value" :style="{ color: avgScoreColor }">{{ stats.avgScore }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('低质样本') }}</div>
            <div class="stat-value" style="color: #EF4444">{{ stats.lowQualityCount }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('销冠基线') }}</div>
            <div class="stat-value" style="color: #10B981">{{ stats.baselineCount }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-tabs v-model="activeTab" class="content-tabs">
      
      <el-tab-pane ::label="$t('auto.Panel.1.55ce19')" name="scores">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Panel.26.9ce7be') }}</span>
              <el-button @click="loadScores">
                <el-icon><Refresh /></el-icon>
                {{ $t('刷新') }}
              </el-button>
            </div>
          </template>
          <el-table :data="scores" v-loading="scoresLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="session_id" :label="$t('auto.Panel.2.427bc5')" width="160" />
            <el-table-column prop="rule_score" :label="$t('auto.Panel.3.564c79')" width="100">
              <template #default="{ row }">
                <el-tag :type="confTagType(row.rule_score)">{{ (row.rule_score || 0).toFixed(3) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="llm_score" :label="$t('auto.Panel.4.1db9ce')" width="100">
              <template #default="{ row }">
                {{ (row.llm_score || 0).toFixed(3) }}
              </template>
            </el-table-column>
            <el-table-column prop="final_score" :label="$t('auto.Panel.5.00dc6d')" width="100">
              <template #default="{ row }">
                <el-tag :type="confTagType(row.final_score)" effect="dark">{{ (row.final_score || 0).toFixed(3) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="regenerated" :label="$t('auto.Panel.6.458f9d')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.regenerated ? 'warning' : 'info'">{{ row.regenerated ? '是' : '否' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="issues" :label="$t('auto.Panel.7.bf457a')" min-width="180">
              <template #default="{ row }">
                <el-tag v-for="tag in (row.issues || [])" :key="tag" size="small" type="danger" style="margin-right: 4px">{{ tag }}</el-tag>
                <span v-if="!row.issues || row.issues.length === 0" class="muted">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('auto.Panel.8.362e80')" width="170">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            :current-page="scorePage"
            :page-size="scorePageSize"
            :total="scoreTotal"
            :page-sizes="[20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @current-change="loadScores"
            @size-change="loadScores"
            style="margin-top: 16px; text-align: right"
          />
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.Panel.9.9045f6')" name="baselines">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Panel.27.c2b41f') }}</span>
              <el-button @click="loadBaselines">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          <el-table :data="baselines" v-loading="baselinesLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="scenario" :label="$t('auto.Panel.10.c93165')" width="160" />
            <el-table-column prop="staff_id" :label="$t('auto.Panel.11.83bc18')" width="100" />
            <el-table-column prop="avg_reward" :label="$t('auto.Panel.12.6f224d')" width="100">
              <template #default="{ row }">
                <el-tag :type="confTagType(row.avg_reward)">{{ (row.avg_reward || 0).toFixed(3) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="dialogue_count" :label="$t('auto.Panel.13.e97204')" width="100" />
            <el-table-column prop="scripts" :label="$t('auto.Panel.14.c3ba6f')" width="100">
              <template #default="{ row }">
                {{ (row.scripts || []).length }}
              </template>
            </el-table-column>
            <el-table-column prop="updated_at" :label="$t('auto.Panel.15.a001a2')" width="170">
              <template #default="{ row }">
                {{ formatTime(row.updated_at) }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Panel.16.2b6bc0')" width="120" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="viewBaselineScripts(row)">{{ $t('auto.Panel.28.fd52fc') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.Panel.17.efe753')" name="lowQuality">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Panel.29.658a9a') }}</span>
              <el-button @click="loadLowQuality">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          <el-table :data="lowQuality" v-loading="lowQualityLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="session_id" :label="$t('auto.Panel.2.427bc5')" width="160" />
            <el-table-column prop="score" :label="$t('auto.Panel.18.108db1')" width="100">
              <template #default="{ row }">
                <el-tag type="danger">{{ (row.score || 0).toFixed(3) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ai_reply" :label="$t('auto.Panel.19.ae8297')" min-width="320" show-overflow-tooltip />
            <el-table-column prop="issues" :label="$t('auto.Panel.7.bf457a')" min-width="180">
              <template #default="{ row }">
                <el-tag v-for="tag in (row.issues || [])" :key="tag" size="small" type="danger" style="margin-right: 4px">{{ tag }}</el-tag>
                <span v-if="!row.issues || row.issues.length === 0" class="muted">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('auto.Panel.20.319a3f')" width="170">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            :current-page="lowQualityPage"
            :page-size="lowQualityPageSize"
            :total="lowQualityTotal"
            :page-sizes="[20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @current-change="loadLowQuality"
            @size-change="loadLowQuality"
            style="margin-top: 16px; text-align: right"
          />
        </el-card>
      </el-tab-pane>
    </el-tabs>

    
    <el-dialog v-model="scriptsDialogVisible" :title="$t('auto.Panel.21.06f066')" width="720px">
      <div v-if="currentBaseline">
        <p><strong>{{ $t('auto.Panel.30.ae9ded') }}</strong>{{ currentBaseline.scenario }} · <strong>{{ $t('auto.Panel.31.bf597b') }}</strong>{{ currentBaseline.staff_id }} · <strong>{{ $t('auto.Panel.32.3e050d') }}</strong>{{ (currentBaseline.avg_reward || 0).toFixed(3) }}</p>
        <el-divider />
        <div v-for="(s, i) in (currentBaseline.scripts || [])" :key="i" class="script-item">
          <h4>{{ s.title || '话术 ' + (i + 1) }}</h4>
          <p class="script-content">{{ s.content }}</p>
          <p v-if="s.tags" class="muted">标签：{{ (s.tags || []).join('、') }}</p>
          <el-divider v-if="i < (currentBaseline.scripts || []).length - 1" />
        </div>
        <el-empty v-if="!currentBaseline.scripts || currentBaseline.scripts.length === 0" :description="$t('auto.Panel.22.bd8901')" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  getHumanizeScores,
  getHumanizeScoreStats,
  getChampionBaselines,
  getLowQualitySamples
} from '@/api/tuning'

const { t } = useI18n()

const activeTab = ref('scores')

const stats = ref({ totalScored: 0, avgScore: '0.000', lowQualityCount: 0, baselineCount: 0 });
const avgScoreColor = computed(() => {
  const v = parseFloat(stats.value.avgScore)
  if (v >= 0.7) return '#10B981'
  if (v >= 0.5) return '#F59E0B'
  return '#EF4444'
})

const scores = ref([]);
const scorePage = ref(1)
const scorePageSize = ref(20)
const scoreTotal = ref(0)
const scoresLoading = ref(false)
async function loadScores() {
  scoresLoading.value = true
  try {
    const res = await getHumanizeScores({ page: scorePage.value, page_size: scorePageSize.value })
    scores.value = res?.list || res?.data?.list || []
    scoreTotal.value = res?.total || res?.data?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.Panel.33.85937d') + (e?.message || e))
  } finally {
    scoresLoading.value = false
  }
}

const baselines = ref([]);
const baselinesLoading = ref(false)
async function loadBaselines() {
  baselinesLoading.value = true
  try {
    const res = await getChampionBaselines({ page: 1, page_size: 50 })
    baselines.value = res?.list || res?.data?.list || []
  } catch (e) {
    ElMessage.error(t('auto.Panel.34.30eeda') + (e?.message || e))
  } finally {
    baselinesLoading.value = false
  }
}

const lowQuality = ref([]);
const lowQualityPage = ref(1)
const lowQualityPageSize = ref(20)
const lowQualityTotal = ref(0)
const lowQualityLoading = ref(false)
async function loadLowQuality() {
  lowQualityLoading.value = true
  try {
    const res = await getLowQualitySamples({ page: lowQualityPage.value, page_size: lowQualityPageSize.value })
    lowQuality.value = res?.list || res?.data?.list || []
    lowQualityTotal.value = res?.total || res?.data?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.Panel.35.2fe87d') + (e?.message || e))
  } finally {
    lowQualityLoading.value = false
  }
}

async function loadStats() {
  try {
    const res = await getHumanizeScoreStats({ range: '24h' })
    const data = res?.data || res || {}
    stats.value = {
      totalScored: data.total || 0,
      avgScore: (data.avg_score || 0).toFixed(3),
      lowQualityCount: data.low_quality_count || 0,
      baselineCount: baselines.value.length || 0
    }
  } catch (e) {
    stats.value = { totalScored: 0, avgScore: '0.000', lowQualityCount: 0, baselineCount: 0 }
  }
}

const scriptsDialogVisible = ref(false);
const currentBaseline = ref(null)
function viewBaselineScripts(row) {
  currentBaseline.value = row
  scriptsDialogVisible.value = true
}

function confTagType(v) {
  if (v >= 0.7) return 'success'
  if (v >= 0.5) return 'warning'
  return 'danger'
}
function formatTime(t) {
  if (!t) return '-'
  try {
    return new Date(t).toLocaleString('zh-CN', { hour12: false })
  } catch {
    return t
  }
}

onMounted(async () => {
  await Promise.all([loadScores(), loadBaselines(), loadLowQuality()])
  await loadStats()
})
</script>

<style scoped>
.humanize-panel { padding: 16px; }
.header-card { margin-bottom: 16px; }
.header-content h2 { margin: 0 0 4px 0; font-size: 20px; }
.subtitle { color: #909399; margin: 0; font-size: 13px; }
.stats-row { margin-bottom: 16px; }
.stat-item { text-align: center; }
.stat-label { color: #909399; font-size: 12px; margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: 600; color: #303133; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.script-item h4 { margin: 0 0 8px 0; color: #303133; }
.script-content { background: #f5f7fa; padding: 12px; border-radius: 4px; line-height: 1.6; }
.muted { color: #909399; font-size: 12px; }
</style>
