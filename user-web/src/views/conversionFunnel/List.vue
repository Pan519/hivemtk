<template>
  <div class="conversion-funnel">
    <el-card class="page-card">
      <div class="page-header">
        <div>
          <h2 class="page-title">{{ $t('auto.List.233d6b00') }}</h2>
          <p class="page-sub">{{ $t('auto.List.73619d48') }}</p>
        </div>
        <el-button type="primary" @click="loadAll">{{ $t('auto.List.694fc5ef') }}</el-button>
      </div>

      <el-row :gutter="16" class="summary-row">
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.fe5c7b05')" :value="summary.totalEnter" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.043c174c')" :value="summary.totalConvert" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.85117a42')" :value="summary.overallRate" suffix="%" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.f2896345')" :value="summary.avgLossRate" suffix="%" />
        </el-col>
      </el-row>

      <el-row :gutter="16" class="chart-row">
        <el-col :span="14">
          <el-card shadow="never" :header="$t('auto.List.d9d6b8f2')">
            <div ref="funnelChart" class="chart-box"></div>
          </el-card>
        </el-col>
        <el-col :span="10">
          <el-card shadow="never" :header="$t('auto.List.2c3348b2')">
            <el-table :data="stageTable" stripe>
              <el-table-column prop="name" :label="$t('auto.List.d316330b')" />
              <el-table-column prop="count" :label="$t('auto.List.0bf60b32')" />
              <el-table-column :label="$t('auto.List.989e84a4')">
                <template #default="{ row }">{{ (row.rate || 0).toFixed(1) }}%</template>
              </el-table-column>
              <el-table-column :label="$t('auto.List.f620ee32')">
                <template #default="{ row }">{{ (row.dropRate || 0).toFixed(1) }}%</template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>

      <el-card shadow="never" :header="$t('auto.List.46524ca7')" class="loss-card">
        <el-form inline>
          <el-form-item :label="$t('auto.List.dddb3a68')">
            <el-select v-model="selectedStage" :placeholder="$t('auto.List.dddb3a68')" @change="loadLoss">
              <el-option v-for="s in stageTable" :key="s.stage" :label="s.name" :value="s.stage" />
            </el-select>
          </el-form-item>
          <el-button type="primary" @click="loadLoss">{{ $t('auto.List.bee912d7') }}</el-button>
        </el-form>
        <el-descriptions v-if="loss" :column="3" border>
          <el-descriptions-item :label="$t('auto.List.d316330b')">{{ loss.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.0bf60b32')">{{ loss.count }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.989e84a4')">{{ (loss.rate || 0).toFixed(1) }}%</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.101693d7')">{{ loss.avgDuration.toFixed(1) }}</el-descriptions-item>
        </el-descriptions>
        <el-table v-if="loss && loss.topSources.length" :data="loss.topSources" stripe class="src-table">
          <el-table-column :label="$t('auto.List.92865962')">
            <template #default="{ row }">{{ getChannelLabel(row.source) }}</template>
          </el-table-column>
          <el-table-column prop="count" :label="$t('auto.List.0bf60b32')" />
        </el-table>
      </el-card>
    </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import ConversionFunnelApi from '@/api/conversionFunnel'
import { getChannelLabel } from '@/constants/channel'

// 按阶段名取计数，不按数组下标取。
//
// 这里原先读的是 stages[0] 与 stages[stages.length - 1]：T-P4-06 在末尾接上商机段之后，
// 「最后一段」从会话变成了商机，于是「转化量(会话)」这块 KPI 会把商机数当会话数显示，
// 端到端转化率同时从 访问→会话 悄悄变成 访问→商机 —— 标签没动、口径动了，
// 而且后端加段不会让前端任何一处变红。阶段名是后端的对外契约（词表与响应各有测试钉住），
// 下标不是。
function countByStage(stages, key) {
  const leg = stages.find((s) => s.stage === key)
  return leg ? leg.count || 0 : 0
}

export default {
  name: 'ConversionFunnelList',
  data() {
    return {
      report: null,
      stageTable: [],
      summary: { totalEnter: 0, totalConvert: 0, overallRate: 0, avgLossRate: 0 },
      selectedStage: 'clue',
      loss: null,
      funnelChart: null,
      chartInst: null
    }
  },
  mounted() {
    this.chartInst = safeInit(this.$refs.funnelChart)
    this.loadAll()
  },
  beforeUnmount() {
    this._destroyed = true
    if (this.chartInst) {
      this.chartInst.dispose()
      this.chartInst = null
    }
  },
  methods: {
    getChannelLabel,
    async loadAll() {
      const res = await ConversionFunnelApi.getFunnel({})
      if (this._destroyed) return
      this.report = res || {}
      const stages = Array.isArray(this.report.stages) ? this.report.stages : []
      this.stageTable = stages.map((s) => ({
        stage: s.stage,
        name: s.name,
        count: s.count,
        rate: s.rate || 0,
        dropRate: s.drop_rate || 0
      }))
      const enterCount = countByStage(stages, 'visit')
      const convertCount = countByStage(stages, 'session')
      this.summary = {
        totalEnter: enterCount,
        totalConvert: convertCount,
        overallRate: enterCount ? (convertCount / enterCount) * 100 : 0,
        avgLossRate: stages.length ? stages.reduce((a, s) => a + (s.drop_rate || 0), 0) / stages.length : 0
      }
      this.renderFunnel()
      this.loadLoss()
    },
    async loadLoss() {
      if (!this.selectedStage) return
      const res = await ConversionFunnelApi.getStageDetails(this.selectedStage, {})
      if (this._destroyed) return
      const d = res || {}
      this.loss = {
        name: d.name || this.selectedStage,
        count: d.count || 0,
        rate: d.rate || 0,
        avgDuration: d.avg_duration_seconds || 0,
        topSources: Array.isArray(d.top_sources) ? d.top_sources : []
      }
    },
    renderFunnel() {
      if (!this.chartInst) return
      this.chartInst.setOption({
        tooltip: { trigger: 'item', formatter: '{b}: {c}' },
        series: [
          {
            type: 'funnel',
            left: '10%',
            width: '80%',
            label: { show: true, position: 'inside', formatter: '{b}: {c}' },
            data: this.stageTable.map((s) => ({ name: s.name, value: s.count }))
          }
        ]
      })
    }
  }
}
</script>

<style scoped>
.page-card { padding: 8px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; }
.page-sub { margin: 4px 0 0; color: #909399; font-size: 13px; }
.summary-row { margin-bottom: 16px; }
.chart-row { margin-bottom: 16px; }
.chart-box { height: 320px; }
.loss-card { margin-top: 8px; }
.src-table { margin-top: 12px; }
</style>
