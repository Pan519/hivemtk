<template>
  <div class="ai-productivity">
    <el-card class="page-card">
      <div class="page-header">
        <div>
          <h2 class="page-title">{{ $t('auto.List.063fc77e') }}</h2>
          <p class="page-sub">{{ $t('auto.List.c43ffda4') }}</p>
        </div>
        <el-button type="primary" @click="loadAll">{{ $t('auto.List.694fc5ef') }}</el-button>
      </div>

      <el-row :gutter="16" class="kpi-row">
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.7ed257f5')" :value="rep.total_conversations" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.cde02bb8')" :value="rep.ai_ratio" suffix="%" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.3871cd44')" :value="Number(Number(rep.avg_response_time || 0).toFixed(1))" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.71505437')" :value="rep.conversion_rate" suffix="%" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.9b925fc0')" :value="rep.ai_replies" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.58661dd9')" :value="rep.human_replies" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.List.2550ea21')" :value="rep.total_conversions" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="LLM Token" :value="rep.llm_tokens" />
        </el-col>
      </el-row>

      <el-card shadow="never" :header="$t('auto.List.3c34782b')" class="trend-card">
        <div ref="trendChart" class="chart-box"></div>
      </el-card>
    </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import i18n from '@/i18n'
const t = i18n.global.t
import { safeInit } from '@/utils/echarts'
import AIProductivityApi from '@/api/aiProductivity'

export default {
  name: 'AIProductivityList',
  data() {
    return {
      rep: {},
      trend: [],
      trendChart: null,
      chartInst: null
    }
  },
  mounted() {
    this.chartInst = safeInit(this.$refs.trendChart)
    this.loadAll()
  },
  beforeUnmount() {
    if (this.chartInst) this.chartInst.dispose()
  },
  methods: {
    async loadAll() {
      const [rep, trendRes] = await Promise.all([
        AIProductivityApi.getReport({}),
        AIProductivityApi.getTrend({ days: 30 })
      ])
      this.rep = rep || {}
      const trendData = trendRes && Array.isArray(trendRes.trend) ? trendRes.trend : (Array.isArray(trendRes) ? trendRes : [])
      this.trend = t
      this.renderTrend()
    },
    renderTrend() {
      if (!this.chartInst) return
      const dates = this.trend.map((d) => d.date)
      const conv = this.trend.map((d) => d.conversations)
      const ai = this.trend.map((d) => d.ai_replies)
      const orders = this.trend.map((d) => d.conversions)
      this.chartInst.setOption({
        tooltip: { trigger: 'axis' },
        legend: { data: [t('auto.List.4d4b9c2d'), 'AI 回复数', t('auto.List.d4268d7f')] },
        xAxis: { type: 'category', data: dates },
        yAxis: { type: 'value' },
        series: [
          { name: t('auto.List.4d4b9c2d'), type: 'line', data: conv },
          { name: 'AI 回复数', type: 'line', data: ai },
          { name: t('auto.List.d4268d7f'), type: 'bar', data: orders }
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
.kpi-row { margin-bottom: 16px; }
.trend-card { margin-top: 8px; }
.chart-box { height: 340px; }
</style>
